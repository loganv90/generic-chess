package main

import (
    "fmt"
    "go-app/chess"
)

type Hub struct {
    clients map[*Client]bool
    broadcast chan []byte
    register chan *Client
    unregister chan *Client
    move chan *ClientMoveData
    capacity int
    game chess.Game
}

func newHub() *Hub {
    game, err := chess.NewSimpleGame()
    if err != nil {
        panic(err)
    }

    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
        move:       make(chan *ClientMoveData),
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
        case clientMoveData := <-h.move:
            h.game.Execute(
                clientMoveData.moveData.XFrom,
                clientMoveData.moveData.YFrom,
                clientMoveData.moveData.XTo,
                clientMoveData.moveData.YTo,
            )
            fmt.Println(h.game.Print())
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
