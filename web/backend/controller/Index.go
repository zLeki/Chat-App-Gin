package controller

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zLeki/Chat-App-Gin/global"
	"github.com/zLeki/Chat-App-Gin/web/backend/helpers"
	"golang.org/x/crypto/bcrypt"
)

func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "",
			"siteKey": "6Lcxl8caAAAAALgD74KC0BkTCUnP9sDfotd_02ot",
			"success": "",
		})
	}
}
func IndexPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(global.Userkey)
		if user != nil {
			c.Redirect(http.StatusMovedPermanently, "/chat")
		}
		tos := helpers.Sanitize(c.PostForm("terms_and_cons"))
		log.Println("tos", tos)
		email_us := helpers.Sanitize(c.PostForm("emauil_us"))
		pass_us := helpers.Sanitize(c.PostForm("pass_us"))
		conf_pass_us := helpers.Sanitize(c.PostForm("conf_pass_us"))
		if conf_pass_us == "" {

			if email_us == "" || pass_us == "" {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"content": "These fields cannot be empty!",
				})
				return
			}
			db := helpers.GetDB()
			defer db.Close()
			db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, password TEXT)")
			var name string
			var email string
			var pass string
			rows, err := db.Query("SELECT name, email, password FROM users WHERE email = ?", email_us)
			for rows.Next() {
				err := rows.Scan(&name, &email, &pass)
				if err != nil {
					log.Println(err)
					return
				}
			}
			if !CheckPasswordHash(pass_us, pass) {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"content": "Password is incorrect",
				})
				return
			}

			if err != nil {
				log.Println(err)
				c.HTML(http.StatusOK, "index.html", gin.H{
					"content": "User doesn't exists",
				})
				return
			}
			session.Set(global.Userkey, email_us)
			session.Save()
			c.Redirect(http.StatusMovedPermanently, "/chat")
		} else {
			if tos != "on" {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"content": "You must accept the terms and conditions",
				})
				return
			}
			if conf_pass_us != pass_us {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"content": "Passwords don't match",
				})
				return
			}
			db := helpers.GetDB()
			defer db.Close()
			db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, password TEXT)")
			var name string
			db.QueryRow("SELECT name FROM users WHERE email = ?", email_us).Scan(&name)
			if name != "" {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"content": "User already exists",
				})
				return
			}
			q, err := db.Prepare("INSERT INTO users (name, email, password) VALUES(?, ?, ?)")
			if err != nil {
				log.Println(err)
				return
			}
			_, err = q.Exec(name, email_us, HashPassword(pass_us))
			if err != nil {
				log.Println(err)
				return
			}
			session.Set(global.Userkey, email_us)
			session.Save()
			c.Redirect(http.StatusMovedPermanently, "/chat")

		}

	}
}
func HashPassword(password string) string { // hash the password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
