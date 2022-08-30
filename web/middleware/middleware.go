package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		log.Println("User not logged in")
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
		return
	}
	c.Next()
}
