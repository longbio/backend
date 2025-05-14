package handlers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/longbio/backend/internal/auth"
	"github.com/longbio/backend/internal/database"
	"github.com/longbio/backend/internal/messages"
	"github.com/longbio/backend/internal/models"
	"github.com/longbio/backend/internal/responses"
	"github.com/longbio/backend/internal/services"
	"gorm.io/gorm"
)

func SendVerificationEmail(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.ApiResponse{Status: http.StatusBadRequest, Message: messages.GeneralBadRequest, Data: nil})
		return
	}

	// Generate a 5-digit OTP
	rand.NewSource(time.Now().UnixNano())
	otp := rand.Intn(90000) + 10000

	if err := database.RedisClient.Set(context.Background(), "email:"+request.Email, otp, 24*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ApiResponse{Status: http.StatusInternalServerError, Message: messages.GeneralFailed, Data: nil})
		return
	}

	go services.SendVerificationCodeViaEmail(request.Email, strconv.Itoa(otp))

	c.JSON(http.StatusOK, responses.ApiResponse{Status: http.StatusOK, Message: messages.GeneralSuccess, Data: nil})
}

func VerifyEmail(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: messages.GeneralFailed,
			Data:    nil,
		})
		return
	}

	storedOTP, err := database.RedisClient.Get(context.Background(), "email:"+request.Email).Result()
	if err != nil || storedOTP != request.Code {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.ApiResponse{Status: http.StatusBadRequest, Message: messages.InvalidVerificationCode, Data: nil})
		return
	}

	database.RedisClient.Del(context.Background(), "email:"+request.Email)

	isNewUser := false

	var user models.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.User{
				Email:             request.Email,
				Gender:            "male",
				MaritalStatus:     "single",
				EducationalStatus: "none",
			}

			if err := database.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, responses.ApiResponse{Status: http.StatusInternalServerError, Message: messages.GeneralFailed, Data: nil})
				return
			}

			isNewUser = true
		} else {
			c.JSON(http.StatusInternalServerError, responses.ApiResponse{Status: http.StatusInternalServerError, Message: messages.GeneralFailed, Data: nil})
			return
		}
	}

	accessToken, err := auth.GenerateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ApiResponse{Status: http.StatusInternalServerError, Message: messages.GeneralFailed, Data: nil})
		return
	}

	refreshToken := uuid.New().String()
	err = database.RedisClient.Set(context.Background(), "refresh:"+refreshToken, user.ID, 30*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ApiResponse{Status: http.StatusInternalServerError, Message: messages.GeneralFailed, Data: nil})
		return
	}

	c.JSON(http.StatusOK, responses.ApiResponse{
		Status:  http.StatusOK,
		Message: messages.GeneralSuccess,
		Data: map[string]any{
			"status":    "success",
			"isNewUser": isNewUser,
			"tokens":    map[string]any{"accessToken": accessToken, "refreshToken": refreshToken},
		},
	})
}

func RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: messages.GeneralFailed,
			Data:    nil,
		})
		return
	}

	userIDStr, err := database.RedisClient.Get(context.Background(), "refresh:"+request.RefreshToken).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ApiResponse{
			Status:  http.StatusUnauthorized,
			Message: messages.InvalidRefreshToken,
			Data:    nil,
		})
		return
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: messages.GeneralFailed,
			Data:    nil,
		})
		return
	}

	newAccessToken, err := auth.GenerateJWTToken(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: messages.GeneralFailed,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, responses.ApiResponse{
		Status:  http.StatusOK,
		Message: messages.GeneralSuccess,
		Data: map[string]interface{}{
			"status": "success",
			"tokens": map[string]interface{}{
				"accessToken":  newAccessToken,
				"refreshToken": request.RefreshToken,
			},
		},
	})
}
