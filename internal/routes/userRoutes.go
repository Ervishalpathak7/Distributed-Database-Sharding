package userRoutes

import (
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/internal/controllers"
	"github.com/gin-gonic/gin"
)


func RegisterRoutes(router *gin.Engine) {
	// User routes
	router.GET("/user/:userId", userControllers.CreateUser)
	router.POST("/user", userControllers.CreateUser)
	router.PUT("user/:userId", userControllers.UpdateUser)
	router.DELETE("user/:userId", userControllers.DeleteUser)
}
