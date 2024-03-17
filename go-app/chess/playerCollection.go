package chess

import (
	"fmt"
)

/*
Responsible for:
- keeping track of the players in the game
*/
type PlayerCollection interface {
    eliminate(color int)
    restore(color int)
    getCurrent() int
    setCurrent(color int)
    getWinner() int
    setWinner(color int)
    getGameOver() bool
    setGameOver(gameOver bool)

    getNextAndRemaining() (int, int)
    getPlayers() int

    Copy() (PlayerCollection, error)
}

func newSimplePlayerCollection(numberOfPlayers int) (*SimplePlayerCollection, error) {
    if numberOfPlayers <= 0 {
        return nil, fmt.Errorf("not enough players")
    }

    playersAlive := make([]bool, numberOfPlayers)
    for i := range playersAlive {
        playersAlive[i] = true
    }

	return &SimplePlayerCollection{
        players: numberOfPlayers,
        playersAlive: playersAlive,
        currentPlayer: 0,
        winningPlayer: -1,
        gameOver: false,
	}, nil
}

type SimplePlayerCollection struct {
    players int
    playersAlive []bool
    currentPlayer int
    winningPlayer int
    gameOver bool
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

func (s *SimplePlayerCollection) Copy() (PlayerCollection, error) {
    playersAlive := make([]bool, s.players)
    for color, alive := range s.playersAlive {
        playersAlive[color] = alive
    }

    return &SimplePlayerCollection{
        players: s.players,
        playersAlive: playersAlive,
        currentPlayer: s.currentPlayer,
        winningPlayer: s.winningPlayer,
        gameOver: s.gameOver,
    }, nil
}

