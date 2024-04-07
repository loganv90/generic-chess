package main

import (
    "bytes"
	"fmt"
	"time"

    "encoding/json"
	"github.com/gorilla/websocket"
    "go-app/chess"
)

const writeWait = 10 * time.Second
const pongWait = 60 * time.Second
const pingPeriod = (pongWait * 9) / 10
const maxMessageSize = 1024 

type ClientMessage struct {
    message []byte
    client Client
}

type Client interface {
    sendMessage(message []byte) error
    close() error
}

func newPlayerClient(hub *Hub, conn *websocket.Conn) (*PlayerClient, error) {
    playerClient := &PlayerClient{
        hub: hub,
        conn: conn,
        send: make(chan []byte, 256),
    }

    return playerClient, nil
}

type PlayerClient struct {
    hub *Hub
    conn *websocket.Conn
    send chan []byte
}

func (c *PlayerClient) sendMessage(message []byte) error {
    select {
    case c.send <- message:
    default:
        c.hub.unregister <- c
    }
    return nil
}

func (c *PlayerClient) close() error {
    close(c.send)
    return nil
}

func (c *PlayerClient) readLoop() {
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    defer func() {
        fmt.Println("ending readLoop")
        c.hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            fmt.Println("error reading message")
            break
        }
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

        c.hub.send <- &ClientMessage{
            message: message,
            client: c,
        }
    }
}

func (c *PlayerClient) writeLoop() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        fmt.Println("ending writeLoop")
        c.hub.unregister <- c
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

func newBotClient(hub *Hub, game chess.Game) (*BotClient, error) {
    bot, err := chess.NewSimpleBot(game)
    if err != nil {
        return nil, err
    }

    botClient := &BotClient{
        hub: hub,
        bot: bot,
        send: make(chan []byte, 256),
    }

    return botClient, nil
}

type BotClient struct {
    hub *Hub
    bot chess.Bot
    send chan []byte
}

func (c *BotClient) sendMessage(message []byte) error {
    select {
    case c.send <- message:
    default:
        c.hub.unregister <- c
    }
    return nil
}

func (c *BotClient) close() error {
    close(c.send)
    return nil
}

func (c *BotClient) run() {
    defer func() {
        fmt.Println("ending run")
        c.hub.unregister <- c
    }()
    for {
        select {
        case unmarshalledMessage, ok := <-c.send:
            if !ok {
                return
            }

            var message *Message
            err := json.Unmarshal(unmarshalledMessage, &message)
            if err != nil {
                continue
            }

            if message.Type != "BoardState" {
                continue
            }

            var boardData *chess.BoardData
            err = json.Unmarshal(message.Data, &boardData)
            if err != nil {
                continue
            }

            for _, color := range c.hub.botColors {
                if !boardData.GameOver && color == boardData.CurrentPlayer {
                    moveKey, err := c.bot.FindMoveIterativeDeepening()
                    if err != nil {
                        fmt.Println("error finding move")
                        break
                    }

                    marshalledMoveKey, err := json.Marshal(moveKey)
                    if err != nil {
                        fmt.Println("error marshalling moveKey")
                        break
                    }

                    message := Message{
                        Type: "move",
                        Data: marshalledMoveKey,
                    }

                    marshalledMessage, err := json.Marshal(message)
                    if err != nil {
                        fmt.Println("error marshalling moveKey")
                        break
                    }

                    c.hub.send <- &ClientMessage{
                        message: marshalledMessage,
                        client: c,
                    }

                    break
                }
            }
        }
    }
}

