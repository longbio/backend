package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/longbio/backend/internal/auth"
	"github.com/longbio/backend/internal/database"
	"github.com/longbio/backend/internal/messages"
	"github.com/longbio/backend/internal/models"
	"github.com/longbio/backend/internal/responses"
)

func GetCurrentUser(c *gin.Context) {
	userId, err := auth.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ApiResponse{Status: http.StatusUnauthorized, Message: messages.GeneralUnauthorized, Data: nil})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ApiResponse{Status: http.StatusNotFound, Message: messages.UserNotFound, Data: nil})
		return
	}

	c.JSON(http.StatusOK, responses.ApiResponse{
		Status:  http.StatusOK,
		Message: messages.GeneralSuccess,
		Data: map[string]any{
			"id":           user.ID,
			"hasPet":       user.HasPet,
			"fullName":     user.FullName,
			"doesExercise": user.DoesExercise,
		},
	})
}

func UpdateCurrentUser(c *gin.Context) {
	userId, err := auth.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ApiResponse{Status: http.StatusUnauthorized, Message: messages.GeneralUnauthorized, Data: nil})
		return
	}

	var request struct {
		FullName     string `json:"fullName"`
		HasPet       bool   `json:"hasPet"`
		DoesExercise bool   `json:"doesExercise"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.ApiResponse{Status: http.StatusBadRequest, Message: messages.GeneralFailed, Data: nil})
		return
	}

	var user models.User

	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ApiResponse{Status: http.StatusNotFound, Message: messages.UserNotFound, Data: nil})
		return
	}

	updates := map[string]any{}
	updates["full_name"] = strings.TrimSpace(request.FullName)

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ApiResponse{Status: http.StatusInternalServerError, Message: messages.GeneralFailed, Data: nil})
		return
	}

	c.JSON(http.StatusOK, responses.ApiResponse{Status: http.StatusOK, Message: messages.GeneralSuccess, Data: user})
}
