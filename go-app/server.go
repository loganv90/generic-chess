package main

import (
    "bytes"
	"fmt"
	"net/http"
	"time"
    "math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const writeWait = 10 * time.Second
const pongWait = 60 * time.Second
const pingPeriod = (pongWait * 9) / 10
const maxMessageSize = 512
const roomIdLength = 4

var newline = []byte{'\n'}
var space = []byte{' '}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")
func randomString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func main() {
    var hubs = make(map[string]*Hub)

	router := gin.Default()
    router.GET("/ws", func(c *gin.Context) {
        var roomId string
        attempts := 0
        attemptsLimit := 10
        for {
            roomId = randomString(roomIdLength)
            if _, ok := hubs[roomId]; !ok {
                break
            }
            attempts++
            if attempts > attemptsLimit {
                fmt.Println("couldn't create game")
                return
            }
        }
        fmt.Println("new connection, new game, gameId: ", roomId)

        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            return
        }

        hub := newHub()
        hubs[roomId] = hub
        go hub.run()

        client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
        client.hub.register <- client

        go client.readLoop()
        go client.writeLoop()
    })
    router.GET("/ws/:gameId", func(c *gin.Context) {
        hub, ok := hubs[c.Param("gameId")]
        if !ok {
            fmt.Println("couldn't find game")
            return
        }
        fmt.Println("new connection, existing game, gameId: ", c.Param("gameId"))

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

        client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
        client.hub.register <- client

        go client.readLoop()
        go client.writeLoop()
	})
	router.Run(":8080")
}

func (c *Client) readLoop() {
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    defer func() {
        fmt.Println("closing connection")
        c.hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        c.hub.broadcast <- message
        fmt.Println("received message")
    }
}

func (c *Client) writeLoop() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        c.conn.Close()
    }()
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write(newline)
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
