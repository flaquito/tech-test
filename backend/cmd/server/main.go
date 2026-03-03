package main

import (
	"database/sql"
	"log"
	"os"

	"example.com/backend/rest"
	"example.com/backend/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	baseUrl := "http://localhost:8080"
	db, err := sql.Open("sqlite3", "./data/app.db")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll("data", 0755)

	r := gin.Default()

	r.Use(cors.Default())
	r.Static("/images", "./data/images")

	r.GET("/posts", rest.GetPostsHandler(db, baseUrl))
	r.POST("/uploads", rest.PostUploadsHandler(db, "data/images", baseUrl))
	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleWS(c.Writer, c.Request)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
