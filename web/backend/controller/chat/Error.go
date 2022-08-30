package chat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(error string, errr error, user interface{}, c *gin.Context) {
	chats, err := GetChats(user)
	if err != nil {
		fmt.Println("select all from username", err)
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.HTML(http.StatusInternalServerError, "home.html", gin.H{
		"content":       error,
		"Pfp":           GetImage(user.(string)),
		"error":         errr.Error(),
		"Conversations": chats.Conversations,
	})
}
