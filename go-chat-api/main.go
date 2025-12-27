package main

import (
	"fmt"
	"go-chat-api/config"
	"go-chat-api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

func main() {
	r := gin.Default()
	config.ConnectDB()

	go handleMessages()

	r.GET("/ws", handleConnections)

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.GET("/messages", func(c *gin.Context) {
		var messages []models.Message
		config.DB.Find(&messages)
		c.JSON(http.StatusOK, messages)
	})

	fmt.Println("Server berjalan di port :8080")
	r.Run(":8080")
}

func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true
	fmt.Println("Klien baru terhubung!")

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		config.DB.Create(&msg)

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
