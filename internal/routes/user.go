package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/longbio/backend/internal/handlers"
)

func UserRoutes(r *gin.RouterGroup) {
	v1 := r.Group("users")

	v1.GET("me", handlers.GetCurrentUser)
	v1.PATCH("me", handlers.UpdateCurrentUser)
}
