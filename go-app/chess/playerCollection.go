package chess

import (
	"fmt"
)

/*
Responsible for:
- keeping track of the players in the game
*/
type PlayerCollection interface {
    getNext() ([]*Player, error)
    getPlayerColors() []string
    eliminate(color string) error
    restore(color string) error
    getCurrent() (string, error)
    setCurrent(color string) error
    getWinner() (string, error)
    setWinner(color string) error
    getGameOver() (bool, error)
    setGameOver(gameOver bool) error

    GetTransition(b Board, inCheckmate bool, inStalemate bool) (PlayerTransition, error)
    Copy() (PlayerCollection, error)
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
    gameOver bool
}

func (s *SimplePlayerCollection) getNext() ([]*Player, error) {
    next := []*Player{}
    currentPlayer := s.currentPlayer

    for {
        currentPlayer = s.incrementOnce(currentPlayer)

        if s.players[currentPlayer].alive {
            next = append(next, s.players[currentPlayer])
        }

        if s.currentPlayer == currentPlayer {
            break
        }
    }

    return next, nil
}

func (s *SimplePlayerCollection) getPlayerColors() []string {
    colors := []string{}
    for _, p := range s.players {
        colors = append(colors, p.color)
    }
    return colors
}

func (s *SimplePlayerCollection) incrementOnce(start int) int {
    end := (start + 1) % len(s.players)
    if end < 0 {
        end = len(s.players) - 1
    }
    return end
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
        return "", nil
    }

    return s.players[s.winningPlayer].color, nil
}

func (s *SimplePlayerCollection) setWinner(color string) error {
    if color == "" {
        s.winningPlayer = -1
        return nil
    }

    if i, ok := s.playerMap[color]; ok {
        s.winningPlayer = i
        return nil
    }

    return fmt.Errorf("player not found")
}

func (s *SimplePlayerCollection) getGameOver() (bool, error) {
    return s.gameOver, nil
}

func (s *SimplePlayerCollection) setGameOver(gameOver bool) error {
    s.gameOver = gameOver
    return nil
}

func (s *SimplePlayerCollection) GetTransition(b Board, inCheckmate bool, inStalemate bool) (PlayerTransition, error) {
    return playerTransitionFactoryInstance.newIncrementalTransition(b, s, inCheckmate, inStalemate)
}

func (s *SimplePlayerCollection) Copy() (PlayerCollection, error) {
    players := make([]*Player, len(s.players))
    for i, p := range s.players {
        players[i] = &Player{
            color: p.color,
            alive: p.alive,
        }
    }

    playerMap := map[string]int{}
    for k, v := range s.playerMap {
        playerMap[k] = v
    }

    return &SimplePlayerCollection{
        players: players,
        playerMap: playerMap,
        currentPlayer: s.currentPlayer,
        winningPlayer: s.winningPlayer,
        gameOver: s.gameOver,
    }, nil
}

