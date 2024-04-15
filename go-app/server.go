package main

import (
	"fmt"
	"net/http"
    "math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

func createRoomId(hubs map[string]*Hub) (string, error) {
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
            return "", fmt.Errorf("couldn't create game")
        }
    }
    return roomId, nil
}

func startClient(c *gin.Context, hub *Hub, conn *websocket.Conn) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        fmt.Println("error upgrading connection")
        return
    }

    playerClient, err := newPlayerClient(hub, conn)
    if err != nil {
        fmt.Println("error creating client")
        return
    }

    playerClient.hub.register <- playerClient

    go playerClient.writeLoop()
    go playerClient.readLoop()
}

func main() {
    var hubs = make(map[string]*Hub)

	router := gin.Default()
    router.GET("/ws/twobot", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newTwoPlayerHubWithBot()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/smalltwobot", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newSmallTwoPlayerHubWithBot()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/fourbot", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newFourPlayerHubWithBot()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/smallfourbot", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newSmallFourPlayerHubWithBot()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/two", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newTwoPlayerHub()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/smalltwo", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newSmallTwoPlayerHub()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/four", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newFourPlayerHub()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/smallfour", func(c *gin.Context) {
        roomId, err := createRoomId(hubs)
        if err != nil {
            fmt.Println("couldn't create game")
            return
        }
        hub := newSmallFourPlayerHub()
        hubs[roomId] = hub
        go func() {
            hub.run()
            fmt.Println("shutting down hub")
            delete(hubs, roomId)
        }()
        fmt.Println("new connection, new game, gameId: ", roomId)

        startClient(c, hub, nil)
    })
    router.GET("/ws/join/:gameId", func(c *gin.Context) {
        hub, ok := hubs[c.Param("gameId")]
        if !ok {
            fmt.Println("couldn't find game")
            return
        }
        if hub.full() {
            fmt.Println("game is full")
            return
        }
        fmt.Println("new connection, existing game, gameId: ", c.Param("gameId"))

        startClient(c, hub, nil)
	})
	router.Run(":8080")
}

