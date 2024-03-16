package chess

import (
    "fmt"
)

func createPlayerTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (PlayerTransition, error) {
    oldCurrent, _ := p.getCurrent()

    oldWinner, _ := p.getWinner()

    oldGameOver, err := p.getGameOver()
    if err != nil {
        return PlayerTransition{}, err
    }

    next, remaining, err := p.getNextAndRemaining()
    if err != nil {
        return PlayerTransition{}, err
    }

    var newCurrent int
    var newWinner int
    var newGameOver bool

    if inStalemate {
        newCurrent = oldCurrent
        newWinner = -1
        newGameOver = true
    } else if remaining < 1 {
        newCurrent = oldCurrent
        newWinner = oldWinner
        newGameOver = true
    } else if remaining == 1 {
        newCurrent = next
        newWinner = next
        newGameOver = true
    } else if remaining == 2 {
        if inCheckmate {
            newCurrent = next
            newWinner = next
            newGameOver = true
        } else {
            newCurrent = next
            newWinner = -1
            newGameOver = false
        }
    } else {
        newCurrent = next
        newWinner = -1
        newGameOver = false
    }

    return PlayerTransition{
        p: p,
        b: b,
        oldCurrent: oldCurrent,
        newCurrent: newCurrent,
        newWinner: newWinner,
        oldWinner: oldWinner,
        oldGameOver: oldGameOver,
        newGameOver: newGameOver,
        eliminated: inCheckmate,
    }, nil
}

type PlayerTransition struct {
    p PlayerCollection
    b Board
    oldCurrent int
    newCurrent int
    oldWinner int
    newWinner int
    oldGameOver bool
    newGameOver bool
    eliminated bool
}

func (s *PlayerTransition) execute() error {
    ok := s.p.setCurrent(s.newCurrent)
    if !ok {
        return fmt.Errorf("invalid color")
    }

    ok = s.p.setWinner(s.newWinner)
    if !ok {
        return fmt.Errorf("invalid color")
    }

    err := s.p.setGameOver(s.newGameOver)
    if err != nil {
        return err
    }

    if !s.eliminated {
        return nil
    }

    err = s.p.eliminate(s.oldCurrent)
    if err != nil {
        return err
    }

    s.b.disablePieces(s.oldCurrent, true)

    return nil
}

func (s *PlayerTransition) undo() error {
    ok := s.p.setCurrent(s.oldCurrent)
    if !ok {
        return fmt.Errorf("invalid color")
    }

    ok = s.p.setWinner(s.oldWinner)
    if !ok {
        return fmt.Errorf("invalid color")
    }

    err := s.p.setGameOver(s.oldGameOver)
    if err != nil {
        return err
    }

    if !s.eliminated {
        return nil
    }

    err = s.p.restore(s.oldCurrent)
    if err != nil {
        return err
    }

    s.b.disablePieces(s.oldCurrent, false)

    return nil
}

