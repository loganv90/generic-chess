package main

import (
    "encoding/json"
    "fmt"
    "go-app/chess"
)

type MoveData struct {
    XFrom int
    YFrom int
    XTo int
    YTo int
}

type ViewData struct {
    X int
    Y int
}

type Hub struct {
    clients map[*Client]bool
    broadcast chan []byte
    register chan *Client
    unregister chan *Client
    send chan *ClientMessage
    capacity int
    game chess.Game
}

func newHub() *Hub {
    game, err := chess.NewSimpleGame()
    if err != nil {
        panic(err)
    }

    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        send:       make(chan *ClientMessage),
        capacity:   2,
        game:       game,
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            if len(h.clients) <= 0 {
                return
            }
        case clientMessage := <-h.send:
            h.handleMessage(
                clientMessage.client,
                clientMessage.message,
            )
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (h *Hub) full() bool {
    return len(h.clients) >= h.capacity
}

func (h *Hub) broadcastMessage(message []byte) {
    for client := range h.clients {
        select {
        case client.send <- message:
        default:
            close(client.send)
            delete(h.clients, client)
        }
    }
}

func (h *Hub) handleMessage(c *Client, message []byte) {
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
        h.handleMoveMessage(messageData)
    } else if messageTitle == "view" {
        h.handleViewMessage(messageData)
    } else {
        fmt.Println("unknown message type")
    }
}

func (h *Hub) handleMoveMessage(messageData json.RawMessage) {
    var moveData MoveData
    err := json.Unmarshal(messageData, &moveData)
    if err != nil {
        fmt.Println("error unmarshalling move data")
        return
    }

    fmt.Println("printing move data", moveData)

    state, err := h.game.Execute(
        moveData.XFrom,
        moveData.YFrom,
        moveData.XTo,
        moveData.YTo,
    )
    if err != nil {
        fmt.Println("error executing move")
        return
    }

    message, err := h.createStateMessage(state)
    if err != nil {
        fmt.Println("error creating state message")
        return
    }

    h.broadcastMessage(message)
    fmt.Println(h.game.Print())
    fmt.Println(state)
}

func (h *Hub) handleViewMessage(messageData json.RawMessage) {
    var viewData ViewData
    err := json.Unmarshal(messageData, &viewData)
    if err != nil {
        fmt.Println("error unmarshalling view data")
        return
    }

    fmt.Println("printing view data", viewData)
}

func (h *Hub) createStateMessage(status *chess.State) ([]byte, error) {
    message, err := json.Marshal(status)
    if err != nil {
        fmt.Println("error marshalling state")
        return nil, err
    }

    return message, nil
}
