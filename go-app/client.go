package main

import (
	"fmt"
	"time"
    "encoding/json"

	"github.com/gorilla/websocket"
)

const writeWait = 10 * time.Second
const pongWait = 60 * time.Second
const pingPeriod = (pongWait * 9) / 10
const maxMessageSize = 512

type ClientMoveData struct {
    moveData MoveData
    client *Client
}

type MoveData struct {
    XFrom int
    YFrom int
    XTo int
    YTo int
}

type BoardData struct {
    squares [][]SquareData
    turn string
}

type SquareData struct {
    color string
    piece string
}

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
            fmt.Println("error reading message")
            break
        }
        c.handleMessage(message)
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

func (c *Client) handleMessage(message []byte) {
    fmt.Println("received message", message)

    var jsonMessage map[string]json.RawMessage
    err := json.Unmarshal(message, &jsonMessage)
    if err != nil {
        fmt.Println("error unmarshalling message")
        return
    }

    var messageTitle string
    err = json.Unmarshal(jsonMessage["title"], &messageTitle)
    if err != nil {
        fmt.Println("error unmarshalling message title")
        return
    }

    var messageData json.RawMessage
    err = json.Unmarshal(jsonMessage["data"], &messageData)
    if err != nil {
        fmt.Println("error unmarshalling message data")
        return
    }

    if messageTitle == "move" {
        c.handleMoveMessage(messageData)
    } else {
        fmt.Println("unknown message type")
    }
}

func (c *Client) handleMoveMessage(messageData json.RawMessage) {
    fmt.Println("handling move message", messageData)

    var moveData MoveData
    err := json.Unmarshal(messageData, &moveData)
    if err != nil {
        fmt.Println("error unmarshalling move data")
        return
    }

    fmt.Println("printing move data", moveData)

    c.hub.move <- &ClientMoveData{
        moveData: moveData,
        client: c,
    }
}

