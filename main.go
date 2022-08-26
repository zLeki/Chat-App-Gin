package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/global"
	"github.com/zLeki/Chat-App-Gin/web/backend/routes"
	"github.com/zLeki/Chat-App-Gin/web/middleware"
	"io/ioutil"
)

func main() {
	router := gin.Default()
	HandleRoutes(router)
	public := router.Group("/")
	routes.PublicRoutes(public)
	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)
	err := router.Run(global.PORT)
	if err != nil {
		panic(err)
	}
}
func HandleRoutes(g *gin.Engine) {
	var filenames []string
	items, _ := ioutil.ReadDir("web/templates")
	for _, item := range items {
		if item.IsDir() {
			subdir, _ := ioutil.ReadDir("web/templates/" + item.Name())
			for _, item2 := range subdir {
				filenames = append(filenames, "web/templates/"+item.Name()+"/"+item2.Name())
			}
		}
	}
	g.LoadHTMLFiles(filenames...)
	g.Use(sessions.Sessions("session", cookie.NewStore(global.Secret)))
	g.Static("/web", "./web")
}
