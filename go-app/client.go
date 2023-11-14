package main

import (
    "bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
    hub *Hub
    conn *websocket.Conn
    send chan []byte
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

