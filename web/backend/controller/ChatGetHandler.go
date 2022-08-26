package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/global"
	"github.com/zLeki/Chat-App-Gin/web/backend/helpers"
	"net/http"
)

func ChatGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(global.Userkey)
		if user == nil {
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		id := c.Param("id")
		if id == "" {

			var conver string
			db := helpers.GetDB()
			err := db.QueryRow("SELECT conversations FROM users WHERE email = ?", user.(string)).Scan(&conver)
			if err != nil {
				fmt.Println(err)
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			}
			var conv Conv
			err = json.Unmarshal([]byte(conver), &conv)
			if err != nil {
				fmt.Println(err)
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			}
			c.HTML(http.StatusOK, "chat.html", gin.H{
				"Conversations": conv.Conversations,
				"content":       "",
				"success":       "",
				"user":          user,
			})
		}
	}
}

type Conv struct {
	Conversations []struct {
		Name       string `json:"name"`
		LastSpoken string `json:"LastSpoken"`
		Pfp        string `json:"pfp"`
	} `json:"conversations"`
}
