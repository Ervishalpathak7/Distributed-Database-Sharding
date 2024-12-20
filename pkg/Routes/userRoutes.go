package Routes

import (
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Controllers"
	"github.com/gin-gonic/gin"
)


func RegisterUserRoutes(router *gin.Engine) {
	// User routes
	router.GET("/user/:userId", Controllers.CreateUser)
	router.POST("/user", Controllers.CreateUser)
	router.PUT("user/:userId", Controllers.UpdateUser)
	router.DELETE("user/:userId", Controllers.DeleteUser)
}
