package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/global"
	"github.com/zLeki/Chat-App-Gin/web/backend/helpers"
	"net/http"
	"strconv"
	"time"
)

func ChatAddPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var convId int
		session := sessions.Default(c)
		usera := session.Get(global.Userkey)
		if usera == nil {
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		add := c.PostForm(helpers.Sanitize("add"))
		db := helpers.GetDB()
		defer db.Close()
		var Users []User
		rows, err := db.Query("SELECT * FROM users WHERE email = ?", usera.(string))
		if err != nil {
			fmt.Println("select all from email", err)
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		defer rows.Close()
		for rows.Next() {
			var user User
			err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Conversations, &user.Pfp)
			if err != nil {
				fmt.Println("scan all from email", err)
				return
			}
			Users = append(Users, user)
		}

		err = rows.Close()
		if err != nil {
			fmt.Println("close rows from email", err)
			return
		}
		rows2, err := db.Query("SELECT * FROM users WHERE name = ?", add)
		if err != nil {
			fmt.Println("select all from username", err)
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		err = db.Close()
		if err != nil {
			fmt.Println("close db", err)
			return
		}
		for rows2.Next() {
			var user User

			err := rows2.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Conversations, &user.Pfp)
			if err != nil {
				fmt.Println("scan all from username", err)
				return
			}
			if user.Email == usera.(string) {

				Error("You can't add yourself", errors.New("adding attempt to yourself prohibited"), usera, c)
				return
			}
			Users = append(Users, user)
		}
		err = rows2.Close()
		if err != nil {
			fmt.Println("close rows2", err)
			return
		}
		if len(Users) < 2 {
			Error("No user found with that name", errors.New("less then 2 users ChatAddPostHandler.go"), usera, c)
			return
		}
		i := 0

		var convPPL ConvPPLStruct
		convPPL.People = append(convPPL.People, struct {
			Name string `json:"name"`
			Pfp  string `json:"pfp"`
		}{Name: Users[0].Name, Pfp: Users[0].Pfp})

		convPPL.People = append(convPPL.People, struct {
			Name string `json:"name"`
			Pfp  string `json:"pfp"`
		}{Name: Users[1].Name, Pfp: Users[1].Pfp})
		convPPLJson, err := json.Marshal(convPPL)

		db1 := helpers.GetDB()
		var messages Messages
		messages.Messages = append(messages.Messages, struct {
			Text   string `json:"text"`
			Sender string `json:"sender"`
			Time   string `json:"time"`
		}{Text: "Hi 👋", Sender: Users[0].Name, Time: time.Now().String()})
		messageJson, err := json.Marshal(messages)

		for i < len(Users) {
			db2 := helpers.GetDB()
			fmt.Println(Users[i].Id)
			var conv Conv
			err := json.Unmarshal([]byte(Users[i].Conversations), &conv)
			if err != nil {
				fmt.Println("json decode", err, Users[i].Conversations)
				c.Redirect(http.StatusMovedPermanently, "/chat/home")
				return
			}
			if i == 0 {
				conv.Conversations = append(conv.Conversations, struct {
					Id         string `json:"id"`
					Name       string `json:"name"`
					LastSpoken string `json:"LastSpoken"`
					Pfp        string `json:"pfp"`
				}{Id: strconv.Itoa(Users[1].Id), Name: add, LastSpoken: time.Now().String(), Pfp: Users[1].Pfp})
			} else if i == 1 {
				conv.Conversations = append(conv.Conversations, struct {
					Id         string `json:"id"`
					Name       string `json:"name"`
					LastSpoken string `json:"LastSpoken"`
					Pfp        string `json:"pfp"`
				}{Id: strconv.Itoa(Users[0].Id), Name: Users[0].Name, LastSpoken: time.Now().String(), Pfp: Users[0].Pfp})
			}
			convJson, err := json.Marshal(conv)

			if err != nil {
				fmt.Println("marshal", err)
				c.Redirect(http.StatusMovedPermanently, "/chat/home")
				return
			}

			if i == 0 {
				_, err = db2.Exec("UPDATE users SET conversations = ? WHERE email = ?", string(convJson), usera.(string))
			} else if i == 1 {
				_, err = db2.Exec("UPDATE users SET conversations = ? WHERE name = ?", string(convJson), add)
			}

			if err != nil {
				fmt.Println("update", err, string(convJson))
				c.Redirect(http.StatusMovedPermanently, "/chat/home")
				return
			}
			db2.Close()
			convId += Users[i].Id
			i += 1
		}
		_, err = db1.Exec("INSERT INTO conversations (id, people, created, messages) VALUES (?, ?, ?, ?)", convId, convPPLJson, time.Now().String(), messageJson)
		if err != nil {
			Error("This conversation already exists!", err, usera, c)
			fmt.Println("insert", err)
			return
		}
		db1.Close()
		c.Redirect(http.StatusMovedPermanently, "/chat/conversation?id="+strconv.Itoa(convId))
	}
}

type ConvPPLStruct struct {
	People []struct {
		Name string `json:"name"`
		Pfp  string `json:"pfp"`
	}
}
type Messages struct {
	Messages []Message
}
type Message struct {
	Text   string `json:"text"`
	Sender string `json:"sender"`
	Time   string `json:"time"`
}
type ReceviedMsg struct {
	Sender       string `json:"sender"`
	Message      string `json:"message"`
	Conversation int    `json:"conversation"`
}
type User struct {
	Id            int
	Name          string
	Email         string
	Password      string
	Conversations string
	Pfp           string
}
