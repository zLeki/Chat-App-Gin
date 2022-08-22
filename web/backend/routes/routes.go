package routes

import (
	"github.com/gin-gonic/gin"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/")
}
func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/chat-rooms")
}
