package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/global"
	"net/http"
)

func LogoutGetHandler() gin.HandlerFunc { // get the logout page
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(global.Userkey, "") // this will mark the session as "written" and hopefully remove the username
		session.Set("session", "")
		session.Clear()
		session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0
		err := session.Save()
		if err != nil {
			fmt.Println("error saving session", err)
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}
