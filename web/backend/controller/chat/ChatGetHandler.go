package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/global"
	"github.com/zLeki/Chat-App-Gin/web/backend/helpers"
	"net/http"
	"strconv"
	"strings"
)

func DMGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(global.Userkey)
		if user == nil {
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		var channelID int
		id, err := strconv.Atoi(helpers.Sanitize(c.Query("id")))
		namea := helpers.Sanitize(c.Query("name"))
		fmt.Println("ID = ", id, "NAME = ", namea)
		if namea != "" {
			fmt.Println(err)
			if err != nil {
				if !strings.Contains(err.Error(), "invalid syntax") {
					Error("Something went wrong", err, user.(string), c)
					fmt.Println("syntax", err)
				}
			}

			fmt.Println("getting id", err)
			var tempid int
			db := helpers.GetDB()
			defer db.Close()
			err := db.QueryRow("SELECT id FROM users WHERE email = ?", user.(string)).Scan(&tempid)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "home.html", gin.H{
					"content": "Something went wrong",
				})
				fmt.Println("email", err)
				return
			}

			channelID += tempid
			tempid = 0
			err = db.QueryRow("SELECT id FROM users WHERE name = ?", namea).Scan(&tempid)
			if err != nil {
				Error("Something went wrong", err, user.(string), c)
				fmt.Println("name", err, c.Query("id"))
				return
			}
			channelID += tempid
			c.Redirect(http.StatusMovedPermanently, "/chat/conversation?id="+strconv.Itoa(channelID))
			return

		}
		chats, err := GetChats(user)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "conversation.html", gin.H{
				"content": "Failed to pull conversations",
			})
			fmt.Println("get chats", err)
			return
		}
		var conv ConvPPLStruct
		var convjson string
		db := helpers.GetDB()
		defer db.Close()
		err = db.QueryRow("SELECT people from conversations where id = ?", uint8(id)).Scan(&convjson)
		if err != nil {
			Error("This conversation does not exist 😧"+strconv.Itoa(channelID), err, user, c)
			fmt.Println("select people", err)
			return
		}
		err = json.Unmarshal([]byte(convjson), &conv)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "home.html", gin.H{
				"content": "Something went wrong",
				"Pfp":     GetImage(user.(string)),
			})
			fmt.Println("unmarshal", err)
			return
		}
		var name string
		err = db.QueryRow("SELECT name FROM users WHERE email = ?", user.(string)).Scan(&name)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "home.html", gin.H{
				"Pfp":     GetImage(user.(string)),
				"content": "Something went wrong",
			})
			fmt.Println("get name", err)
			return
		}
		var msgs string
		err = db.QueryRow("SELECT messages FROM conversations WHERE id = ?", uint8(id)).Scan(&msgs)
		if err != nil {
			Error("Something went wrong", err, user.(string), c)
			fmt.Println("select messages", err)
			return
		}
		var Message Messages
		err = json.Unmarshal([]byte(msgs), &Message)
		if err != nil {
			Error("Something went wrong", err, user.(string), c)
			fmt.Println("unmarshal messages", err)
			return
		}
		//sort slice by time
		messages := make(map[string]string)
		for _, v := range Message.Messages {
			messages[v.Sender] = v.Text

		}
		var a int
		for i, v := range conv.People {
			if v.Name != name {
				a = i
			}
		}
		fmt.Println(a, conv.People[a], conv.People, namea, name)
		for _, v := range conv.People {
			if v.Name == name {
				c.HTML(http.StatusOK, "conversation.html", gin.H{
					"Messages":      messages,
					"People":        conv.People[a],
					"Username":      name,
					"Pfp":           GetImage(user.(string)),
					"Conversations": chats.Conversations,
					"id":            id,
				})
				return
			}
		}
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Pfp":           GetImage(user.(string)),
			"content":       "Sorry, you cannot access this conversation",
			"Conversations": chats.Conversations,
		})
	}
}
func GetImage(user string) string {
	var pfp string
	db := helpers.GetDB()
	defer db.Close()
	err := db.QueryRow("SELECT pfp FROM users WHERE email = ?", user).Scan(&pfp)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return pfp
}
func GetChats(user interface{}) (*Conv, error) {
	var conver string
	db := helpers.GetDB()
	defer db.Close()
	err := db.QueryRow("SELECT conversations FROM users WHERE email = ?", user.(string)).Scan(&conver)
	if err != nil {
		return nil, err
	}
	var conv Conv
	err = json.Unmarshal([]byte(conver), &conv)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}
func ChatGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(global.Userkey)
		if user == nil {
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		chats, err := GetChats(user)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "home.html", gin.H{
				"Conversations": chats.Conversations,
				"content":       "Something went wrong",
				"Pfp":           GetImage(user.(string)),
			})
			fmt.Println(err)
			return
		}
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Conversations": chats.Conversations,
			"content":       "",
			"Pfp":           GetImage(user.(string)),
			"success":       "",
			"user":          user,
		})
	}
}

type Conv struct {
	Conversations []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		LastSpoken string `json:"LastSpoken"`
		Pfp        string `json:"pfp"`
	} `json:"conversations"`
}
