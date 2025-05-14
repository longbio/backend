package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/longbio/backend/internal/handlers"
)

func AuthRoutes(r *gin.RouterGroup) {
	g := r.Group("auth")

	g.POST("verification-code/send", handlers.SendVerificationEmail)
	g.POST("verification-code/verify", handlers.VerifyEmail)

	g.POST("refresh", handlers.RefreshToken)
}
