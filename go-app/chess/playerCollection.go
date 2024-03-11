package chess

import (
	"fmt"
)

/*
Responsible for:
- keeping track of the players in the game
*/
type PlayerCollection interface {
    eliminate(color int) error
    restore(color int) error
    getCurrent() (int, bool)
    setCurrent(color int) bool
    getWinner() (int, bool)
    setWinner(color int) bool
    getGameOver() (bool, error)
    setGameOver(gameOver bool) error

    getNextAndRemaining() (int, int, error)
    getPlayers() int

    GetTransition(b Board, inCheckmate bool, inStalemate bool) (PlayerTransition, error)
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

func (s *SimplePlayerCollection) eliminate(color int) error {
    if s.colorOutOfBounds(color) {
        return fmt.Errorf("invalid color")
    }

    s.playersAlive[color] = false
    return nil
}

func (s *SimplePlayerCollection) restore(color int) error {
    if s.colorOutOfBounds(color) {
        return fmt.Errorf("invalid color")
    }

    s.playersAlive[color] = true
    return nil
}

func (s *SimplePlayerCollection) getCurrent() (int, bool) {
    if s.colorOutOfBounds(s.currentPlayer) {
        return -1, false
    }

    return s.currentPlayer, true
}

func (s *SimplePlayerCollection) setCurrent(color int) bool {
    if s.colorOutOfBounds(color) {
        return false
    }

    s.currentPlayer = color
    return true
}

func (s *SimplePlayerCollection) getWinner() (int, bool) {
    if s.colorOutOfBounds(s.winningPlayer) {
        return -1, true
    }

    return s.winningPlayer, true
}

func (s *SimplePlayerCollection) setWinner(color int) bool {
    if s.colorOutOfBounds(color) {
        s.winningPlayer = -1
        return true
    }

    s.winningPlayer = color
    return true
}

func (s *SimplePlayerCollection) getGameOver() (bool, error) {
    return s.gameOver, nil
}

func (s *SimplePlayerCollection) setGameOver(gameOver bool) error {
    s.gameOver = gameOver
    return nil
}

func (s *SimplePlayerCollection) getNextAndRemaining() (int, int, error) {
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

    return currentPlayer, remaining, nil
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

func (s *SimplePlayerCollection) GetTransition(b Board, inCheckmate bool, inStalemate bool) (PlayerTransition, error) {
    return createPlayerTransition(b, s, inCheckmate, inStalemate)
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

func (s *SimplePlayerCollection) colorOutOfBounds(color int) bool {
    return color < 0 || color >= s.players
}

