package Routes

import (
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Controllers"
	"github.com/gin-gonic/gin"
)


func RegisterProductRoutes(router *gin.Engine) {
	router.GET("/product/:productId", Controllers.GetProduct)
	router.GET("/products", Controllers.ListProducts)
	router.POST("/product", Controllers.CreateProduct)
	router.PUT("/product/:productId", Controllers.UpdateProduct)
	router.DELETE("/product/:productId", Controllers.DeleteProduct)
}