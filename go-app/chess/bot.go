package chess

import (
    "time"
    "fmt"
)

/*
Responsible for:
- keeping track of the state of the game the bot is playing
*/
type Bot interface {
    FindMoveIterativeDeepening() (MoveKey, error)
}

func NewSimpleBot(game Game) (Bot, error) {
    return &SimpleBot{
        game: game,
        depthStart: 2,
        depthLimit: 20,
        timeLimitSeconds: 5,
    }, nil
}

type SimpleBot struct {
    game Game
    depthStart int
    depthLimit int
    timeLimitSeconds int
}

func (b *SimpleBot) FindMoveIterativeDeepening() (MoveKey, error) {
    result := make(chan *MoveKey)
    stop := make(chan bool)

    moveKey := MoveKey{}
    err := fmt.Errorf("No move found")
    endTime := time.Now().Add(time.Duration(b.timeLimitSeconds) * time.Second)

    for depth := b.depthStart; depth <= b.depthLimit; depth++ {
        go b.findMove(depth, result, stop)

        select {
        case <-time.After(time.Until(endTime)):
            stop <- true
            return moveKey, err
        case moveKeyPtr := <-result:
            if moveKeyPtr != nil {
                moveKey = *moveKeyPtr
                err = nil
            }
            fmt.Println("depth reached: ", depth)
        }
    }

    return moveKey, err
}

func (b *SimpleBot) findMove(depth int, result chan *MoveKey, stop chan bool) {
    gameCopy, err := b.game.Copy()
    if err != nil {
        result <- nil
        return
    }

    searcher := newSimpleSearcher(gameCopy, stop)

    moveKey, err := searcher.searchWithMinimax(depth)
    if err != nil {
        result <- nil
        return
    }

    result <- &moveKey
}

