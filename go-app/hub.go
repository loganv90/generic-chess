package main

import (
    "encoding/json"
    "fmt"
    "go-app/chess"
)

type Message struct {
    Type string
    Data json.RawMessage
}

type MoveData struct {
    XFrom int
    YFrom int
    XTo int
    YTo int
    Promotion string
}

type ViewData struct {
    X int
    Y int
}

type Hub struct {
    clients map[Client]bool
    register chan Client
    unregister chan Client
    send chan *ClientMessage
    capacity int
    game chess.Game
}

func newTwoPlayerHub() *Hub {
    game, err := chess.NewSimpleGame()
    if err != nil {
        panic(err)
    }

    hub := &Hub{
        clients:    make(map[Client]bool),
        register:   make(chan Client),
        unregister: make(chan Client),
        send:       make(chan *ClientMessage),
        capacity:   2,
        game:       game,
    }

    // TODO create another function for hubs with bots
    // TODO assign the player the colors that the bot does not have
    // TODO make undo/redo work with bots by undoing/redoing until the player's turn
    botClient, err := newBotClient(hub, game, []string{"black"})
    if err != nil {
        panic(err)
    }

    hub.handleClientJoin(botClient)

    go botClient.run()

    return hub
}

func newFourPlayerHub() *Hub {
    game, err := chess.NewSimpleFourPlayerGame()
    if err != nil {
        panic(err)
    }

    return &Hub{
        clients:    make(map[Client]bool),
        register:   make(chan Client),
        unregister: make(chan Client),
        send:       make(chan *ClientMessage),
        capacity:   4,
        game:       game,
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.handleClientJoin(client)
            if h.empty() {
                h.close()
                return
            }
        case client := <-h.unregister:
            h.handleClientLeave(client)
            if h.empty() {
                h.close()
                return
            }
        case clientMessage := <-h.send:
            h.handleMessage(
                clientMessage.client,
                clientMessage.message,
            )
        }
    }
}

func (h *Hub) full() bool {
    return len(h.clients) >= h.capacity
}

func (h *Hub) empty() bool {
    for client := range h.clients {
        if _, ok := client.(*PlayerClient); ok {
            return false
        }
    }
    return true
}

func (h *Hub) close() {
    for client := range h.clients {
        h.handleClientLeave(client)
    }
}

func (h *Hub) broadcastMessage(message []byte) {
    for client := range h.clients {
        err := client.sendMessage(message)
        if err != nil {
            fmt.Println("error sending message")
            return
        }
    }
}

func (h *Hub) handleClientJoin(c Client) {
    h.clients[c] = true

    message, err := h.createBoardStateMessage()
    if err != nil {
        fmt.Println("error creating state message")
        return
    }

    err = c.sendMessage(message)
    if err != nil {
        fmt.Println("error sending state message")
        return
    }
}

func (h *Hub) handleClientLeave(c Client) {
    if _, ok := h.clients[c]; ok {
        c.close()
        delete(h.clients, c)
    }
}

func (h *Hub) handleMessage(c Client, unmarshalledMessage []byte) {
    var message *Message
    err := json.Unmarshal(unmarshalledMessage, &message)
    if err != nil {
        fmt.Println("error unmarshalling message")
        return
    }

    if message.Type == "move" {
        h.handleMoveMessage(message.Data)
    } else if message.Type == "view" {
        h.handleViewMessage(message.Data)
    } else if message.Type == "undo" {
        h.handleUndoMessage()
    } else if message.Type == "redo" {
        h.handleRedoMessage()
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

    err = h.game.Execute(
        moveData.XFrom,
        moveData.YFrom,
        moveData.XTo,
        moveData.YTo,
        moveData.Promotion,
    )
    if err != nil {
        fmt.Println("error executing move")
        return
    }

    message, err := h.createBoardStateMessage()
    if err != nil {
        fmt.Println("error creating state message")
        return
    }

    h.broadcastMessage(message)
}

func (h *Hub) handleViewMessage(messageData json.RawMessage) {
    var viewData ViewData
    err := json.Unmarshal(messageData, &viewData)
    if err != nil {
        fmt.Println("error unmarshalling view data")
        return
    }

    pieceState, err := h.game.View(
        viewData.X,
        viewData.Y,
    )
    if err != nil {
        fmt.Println("error viewing piece")
        return
    }

    message, err := h.createPieceStateMessage(pieceState)
    if err != nil {
        fmt.Println("error creating state message")
        return
    }

    h.broadcastMessage(message)
}

func (h *Hub) handleUndoMessage() {
    err := h.game.Undo()
    if err != nil {
        fmt.Println("error undoing move")
        return
    }
        
    message, err := h.createBoardStateMessage()
    if err != nil {
        fmt.Println("error creating state message")
        return
    }

    h.broadcastMessage(message)
}

func (h *Hub) handleRedoMessage() {
    err := h.game.Redo()
    if err != nil {
        fmt.Println("error redoing move")
        return
    }
        
    message, err := h.createBoardStateMessage()
    if err != nil {
        fmt.Println("error creating state message")
        return
    }

    h.broadcastMessage(message)
}

func (h *Hub) createBoardStateMessage() ([]byte, error) {
    state, err := h.game.State()
    if err != nil {
        fmt.Println("error getting board state")
        return nil, err
    }

    marshalledState, err := json.Marshal(state)
    if err != nil {
        fmt.Println("error marshalling board state")
        return nil, err
    }

    message := Message{
        Type: "BoardState",
        Data: marshalledState,
    }

    marshalledMessage, err := json.Marshal(message)
    if err != nil {
        fmt.Println("error marshalling board state")
        return nil, err
    }

    return marshalledMessage, nil
}

func (h *Hub) createPieceStateMessage(state *chess.PieceState) ([]byte, error) {
    marshalledState, err := json.Marshal(state)
    if err != nil {
        fmt.Println("error marshalling piece state")
        return nil, err
    }

    message := Message{
        Type: "PieceState",
        Data: marshalledState,
    }

    marshalledMessage, err := json.Marshal(message)
    if err != nil {
        fmt.Println("error marshalling piece state")
        return nil, err
    }

    return marshalledMessage, nil
}

