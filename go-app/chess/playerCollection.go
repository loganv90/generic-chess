package chess

import (
	"fmt"
)

/*
Responsible for:
- keeping track of the players in the game
*/
type PlayerCollection interface {
    increment() (string, error)
    decrement() (string, error)
    eliminate(color string) error
    restore(color string) error
    getCurrent() (string, error)
    setCurrent(color string) error
    getWinner() (string, error)
    setWinner(color string) error
}

func newSimplePlayerCollection(players []*Player) (*SimplePlayerCollection, error) {
    if len(players) <= 1 {
        return nil, fmt.Errorf("not enough players")
    }

    playerMap := map[string]int{}
    for i, p := range players {
        if _, ok := playerMap[p.color]; ok {
            return nil, fmt.Errorf("duplicate player color")
        }

        playerMap[p.color] = i
    }

	return &SimplePlayerCollection{
        players: players,
        playerMap: playerMap,
        currentPlayer: 0,
        winningPlayer: -1,
	}, nil
}

type SimplePlayerCollection struct {
    players []*Player
    playerMap map[string]int
    currentPlayer int
    winningPlayer int
}

func (s *SimplePlayerCollection) increment() (string, error) {
    currentPlayer := s.currentPlayer

    for {
        s.incrementOnce()

        if s.currentPlayer == currentPlayer {
            break
        }

        if s.players[s.currentPlayer].alive {
            break
        }
    }

    return s.players[s.currentPlayer].color, nil
}

func (s *SimplePlayerCollection) incrementOnce() {
    s.currentPlayer = (s.currentPlayer + 1) % len(s.players)
    if s.currentPlayer < 0 {
        s.currentPlayer = len(s.players) - 1
    }
}

func (s *SimplePlayerCollection) decrement() (string, error) {
    currentPlayer := s.currentPlayer

    for {
        s.decrementOnce()

        if s.currentPlayer == currentPlayer {
            break
        }

        if s.players[s.currentPlayer].alive {
            break
        }
    }

    return s.players[s.currentPlayer].color, nil
}

func (s *SimplePlayerCollection) decrementOnce() {
    s.currentPlayer = (s.currentPlayer - 1) % len(s.players)
    if s.currentPlayer < 0 {
        s.currentPlayer = len(s.players) - 1
    }
}

func (s *SimplePlayerCollection) eliminate(color string) error {
    if i, ok := s.playerMap[color]; ok {
        s.players[i].alive = false
        return nil
    } 

    return fmt.Errorf("player not found")
}

func (s *SimplePlayerCollection) restore(color string) error {
    if i, ok := s.playerMap[color]; ok {
        s.players[i].alive = true
        return nil
    } 

    return fmt.Errorf("player not found")
}

func (s *SimplePlayerCollection) getCurrent() (string, error) {
    if s.currentPlayer < 0 || s.currentPlayer >= len(s.players) {
        return "", fmt.Errorf("no current player")
    }

    return s.players[s.currentPlayer].color, nil
}

func (s *SimplePlayerCollection) setCurrent(color string) error {
    if i, ok := s.playerMap[color]; ok {
        s.currentPlayer = i
        return nil
    }

    return fmt.Errorf("player not found")
}

func (s *SimplePlayerCollection) getWinner() (string, error) {
    if s.winningPlayer < 0 || s.winningPlayer >= len(s.players) {
        return "", fmt.Errorf("no winning player")
    }

    return s.players[s.winningPlayer].color, nil
}

func (s *SimplePlayerCollection) setWinner(color string) error {
    if i, ok := s.playerMap[color]; ok {
        s.winningPlayer = i
        return nil
    }

    return fmt.Errorf("player not found")
}

