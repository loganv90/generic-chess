package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func main() {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		fmt.Println("new connection")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Println("closing connection")
			return nil
		})

        go readLoop(conn)
        go writeLoop(conn)
	})
	router.Run(":8080")
}

func readLoop(conn *websocket.Conn) {
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
        fmt.Println("received message")
    }
    conn.Close()
}

func writeLoop(conn *websocket.Conn) {
    for {
        err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, client!"))
        if err != nil {
            break
        }
        time.Sleep(time.Second)
        fmt.Println("sending message")
    }
}
