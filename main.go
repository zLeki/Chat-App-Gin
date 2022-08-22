package main

import (
	"blocksuite-webbackend/globals"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/web/backend/routes"
	"github.com/zLeki/Chat-App-Gin/web/middleware"
	"io/ioutil"
	"net/http"
)

func main() {
	router := gin.Default()
	HandleRoutes(router)
	public := router.Group("/")
	routes.PrivateRoutes(public)
	router.NoRoute(func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/") })
	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)
	err := router.Run(port)
	if err != nil {
		panic(err)
	}
}
func HandleRoutes(g *gin.Engine) {
	var filenames []string
	items, _ := ioutil.ReadDir("webui/templates")
	for _, item := range items {
		if item.IsDir() {
			subdir, _ := ioutil.ReadDir("webui/templates/" + item.Name())
			for _, item2 := range subdir {
				filenames = append(filenames, "webui/templates/"+item.Name()+"/"+item2.Name())
			}
		}
	}
	g.LoadHTMLFiles(filenames...)
	g.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))
	g.Static("/static", "./static")
}
