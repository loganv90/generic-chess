package chess

import (
	"fmt"
    "math/rand"
)

/*
Responsible for:
- keeping track of the players in the game
*/
func newSimplePlayerCollection(numberOfPlayers int) (*SimplePlayerCollection, error) {
    if numberOfPlayers <= 0 {
        return nil, fmt.Errorf("not enough players")
    }

    playersAlive := make([]bool, numberOfPlayers)
    for i := range playersAlive {
        playersAlive[i] = true
    }

    zobristCurrentPlayer := make([]uint64, numberOfPlayers)
    zobristPlayerAlive := make([]uint64, numberOfPlayers)
    for i := 0; i < numberOfPlayers; i++ {
        zobristCurrentPlayer[i] = rand.Uint64()
        zobristPlayerAlive[i] = rand.Uint64()
    }

	return &SimplePlayerCollection{
        players: numberOfPlayers,
        playersAlive: playersAlive,
        currentPlayer: 0,
        winningPlayer: -1,
        gameOver: false,

        zobristCurrentPlayer: zobristCurrentPlayer,
        zobristPlayerAlive: zobristPlayerAlive,
	}, nil
}

type SimplePlayerCollection struct {
    players int
    playersAlive []bool
    currentPlayer int
    winningPlayer int
    gameOver bool

    zobristCurrentPlayer []uint64
    zobristPlayerAlive []uint64
}

func (s *SimplePlayerCollection) colorOutOfBounds(color int) bool {
    return color < 0 || color >= s.players
}

func (s *SimplePlayerCollection) eliminate(color int) {
    if s.colorOutOfBounds(color) {
        return
    }

    s.playersAlive[color] = false
}

func (s *SimplePlayerCollection) restore(color int) {
    if s.colorOutOfBounds(color) {
        return
    }

    s.playersAlive[color] = true
}

func (s *SimplePlayerCollection) getCurrent() int {
    if s.colorOutOfBounds(s.currentPlayer) {
        return -1
    }

    return s.currentPlayer
}

func (s *SimplePlayerCollection) setCurrent(color int) {
    if s.colorOutOfBounds(color) {
        return
    }

    s.currentPlayer = color
}

func (s *SimplePlayerCollection) getWinner() int {
    return s.winningPlayer
}

func (s *SimplePlayerCollection) setWinner(color int) {
    s.winningPlayer = color
}

func (s *SimplePlayerCollection) getGameOver() bool {
    return s.gameOver
}

func (s *SimplePlayerCollection) setGameOver(gameOver bool) {
    s.gameOver = gameOver
}

func (s *SimplePlayerCollection) getNextAndRemaining() (int, int) {
    currentPlayer := s.currentPlayer
    for {
        currentPlayer = s.incrementOnce(currentPlayer)

        if s.playersAlive[currentPlayer] {
            break
        }

        if s.currentPlayer == currentPlayer {
            break
        }
    }

    remaining := 0
    for _, alive := range s.playersAlive {
        if alive {
            remaining++
        }
    }

    return currentPlayer, remaining
}

func (s *SimplePlayerCollection) incrementOnce(start int) int {
    end := (start + 1) % s.players
    if end < 0 {
        end = s.players - 1
    }
    return end
}

func (s *SimplePlayerCollection) getPlayers() int {
    return s.players
}

func (s *SimplePlayerCollection) Copy() (*SimplePlayerCollection, error) {
    simplePlayerCollection, err := newSimplePlayerCollection(s.players)
    if err != nil {
        return nil, err
    }

    for color, alive := range s.playersAlive {
        simplePlayerCollection.playersAlive[color] = alive
    }
    simplePlayerCollection.currentPlayer = s.currentPlayer
    simplePlayerCollection.winningPlayer = s.winningPlayer
    simplePlayerCollection.gameOver = s.gameOver

    return simplePlayerCollection, nil
}

func (s *SimplePlayerCollection) ZobristHash() uint64 {
    hash := uint64(0)

    hash ^= s.zobristCurrentPlayer[s.currentPlayer]
    for color := 0; color < s.players; color++ {
        if s.playersAlive[color] {
            hash ^= s.zobristPlayerAlive[color]
        }
    }

    return hash
}

