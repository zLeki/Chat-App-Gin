package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zLeki/Chat-App-Gin/web/backend/controller"
	"github.com/zLeki/Chat-App-Gin/web/backend/controller/chat"
	"github.com/zLeki/Chat-App-Gin/web/backend/helpers"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var usermap = make(map[int][]*websocket.Conn)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", controller.IndexGetHandler())
	g.POST("/", controller.IndexPostHandler())
	g.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	g.GET("/logout", controller.LogoutGetHandler())
}
func PrivateRoutes(g *gin.RouterGroup) {
	chatGroup := g.Group("/chat")
	{
		chatGroup.GET("/home", chat.ChatGetHandler())
		chatGroup.GET("/conversation", chat.DMGetHandler())
		chatGroup.POST("/add", chat.ChatAddPostHandler())
	}
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Failed to read message: %+v", err)
			break
		}

		if strings.Contains(string(msg), "getMessages") {
			//getMessages[{conversation: 5}]
			a, err := strconv.Atoi(strings.Split(strings.Split(string(msg), " ")[1], "}")[0])
			if err != nil {
				fmt.Printf("Failed to convert string to int: %+v", err)
				continue
			}
			db := helpers.GetDB()
			if err != nil {
				fmt.Printf("Failed to build db: %+v", err)
				continue
			}
			usermap[a] = append(usermap[a], conn)
			var messages string
			err = db.QueryRow("SELECT messages FROM conversations WHERE id = ?", a).Scan(&messages)
			if err != nil {
				fmt.Printf("Failed to select message: %+v", err)
				continue
			}
			var messg chat.Messages
			err = json.Unmarshal([]byte(messages), &messg)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %+v", err)
				continue
			}
			for _, v := range messg.Messages {
				err := conn.WriteMessage(websocket.TextMessage, []byte(`{"sender": "`+v.Sender+`","message": "`+v.Text+`", "time": "`+v.Time+`"}`))
				if err != nil {
					fmt.Println("Failed to write message: %+v", err)
					return
				}
			}
		}
		if strings.Contains(string(msg), "write") {
			//write[{"sender": "Ill", "message": "test", "conversation": 5}]
			var mesg chat.ReceviedMsg
			a := strings.Split(strings.Split(string(msg), "[")[1], "]")[0]
			err := json.Unmarshal([]byte(a), &mesg)
			if err != nil {
				fmt.Printf("WRITE Failed to unmarshal message: %+v", err)
				continue
			}
			var messages string
			db := helpers.GetDB()
			defer db.Close()
			err = db.QueryRow("SELECT messages FROM conversations WHERE id = ?", mesg.Conversation).Scan(&messages)
			if err != nil {
				fmt.Printf("Failed to select message: %+v", err)
				continue
			}
			var amessages chat.Messages
			err = json.Unmarshal([]byte(messages), &amessages)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %+v", err)
				continue
			}
			var newMessage chat.Message
			newMessage.Sender = mesg.Sender
			newMessage.Text = mesg.Message
			newMessage.Time = time.Now().Format("2006-01-02 15:04:05")
			amessages.Messages = append(amessages.Messages, newMessage)
			b, err := json.Marshal(amessages)
			if err != nil {
				fmt.Printf("Failed to marshal message: %+v", err)
				continue
			}
			_, err = db.Exec("UPDATE conversations SET messages = ? WHERE id = ?", string(b), mesg.Conversation)
			if err != nil {
				fmt.Printf("Failed to update message: %+v", err)
				continue
			}
			for _, v := range usermap[mesg.Conversation] {
				err = v.WriteMessage(websocket.TextMessage, []byte(`{"sender": "`+mesg.Sender+`","message": "`+mesg.Message+`", "time": "`+newMessage.Time+`"}`))
				if err != nil {
					//delete item

					continue
				}
			}

		}
	}
}
