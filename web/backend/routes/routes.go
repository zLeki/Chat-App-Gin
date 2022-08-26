package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/web/backend/controller"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", controller.IndexGetHandler())
	g.POST("/", controller.IndexPostHandler())
}
func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/chat", controller.ChatGetHandler())
}
