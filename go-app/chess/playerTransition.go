package chess

var playerTransitionFactoryInstance = PlayerTransitionFactory(&ConcretePlayerTransitionFactory{})

type PlayerTransitionFactory interface {
	newIncrementalTransitionAsPlayerTransition(b Board, p PlayerCollection) (PlayerTransition, error)
    newIncrementalTransition(b Board, p PlayerCollection) (*IncrementalTransition, error)
}

type ConcretePlayerTransitionFactory struct{}

func (f *ConcretePlayerTransitionFactory) newIncrementalTransitionAsPlayerTransition(b Board, p PlayerCollection) (PlayerTransition, error) {
    return f.newIncrementalTransition(b, p)
}

func (f *ConcretePlayerTransitionFactory) newIncrementalTransition(b Board, p PlayerCollection) (*IncrementalTransition, error) {
    oldCurrent, err := p.getCurrent()
    if err != nil {
        return nil, err
    }

    oldWinner, err := p.getWinner()
    if err != nil {
        return nil, err
    }

    oldGameOver, err := p.getGameOver()
    if err != nil {
        return nil, err
    }

    newCurrent := ""
    newWinner := ""
    newGameOver := false
    eliminated := []string{}
    for {
        newPlayer, err := p.getNext()
        if err != nil {
            return nil, err
        }

        if newPlayer.color == oldCurrent {
            newCurrent = newPlayer.color
            newWinner = newPlayer.color
            newGameOver = true
            break
        }

        err = b.CalculateMoves()
        if err != nil {
            return nil, err
        }

        if b.Checkmate(newPlayer.color) {
            eliminated = append(eliminated, newPlayer.color)

            err = p.setCurrent(newPlayer.color)
            if err != nil {
                return nil, err
            }

            continue
        } else if b.Stalemate(newPlayer.color) {
            newCurrent = newPlayer.color
            newGameOver = true
            break
        } else {
            newCurrent = newPlayer.color
            break
        }
    }
    
    err = p.setCurrent(oldCurrent)
    if err != nil {
        return nil, err
    }

    return &IncrementalTransition{
        p: p,
        b: b,
        oldCurrent: oldCurrent,
        newCurrent: newCurrent,
        newWinner: newWinner,
        oldWinner: oldWinner,
        eliminated: eliminated,
        oldGameOver: oldGameOver,
        newGameOver: newGameOver,
    }, nil
}

type PlayerTransition interface {
	execute() error
	undo() error
}

type IncrementalTransition struct {
    p PlayerCollection
    b Board
    oldCurrent string
    newCurrent string
    oldWinner string
    newWinner string
    eliminated []string
    oldGameOver bool
    newGameOver bool
}

func (s *IncrementalTransition) execute() error {
    err := s.p.setCurrent(s.newCurrent)
    if err != nil {
        return err
    }

    err = s.p.setWinner(s.newWinner)
    if err != nil {
        return err
    }

    err = s.p.setGameOver(s.newGameOver)
    if err != nil {
        return err
    }

    for _, color := range s.eliminated {
        err = s.p.eliminate(color)
        if err != nil {
            return err
        }
        err = s.b.disablePieces(color, true)
        if err != nil {
            return err
        }
    }

    err = s.b.CalculateMoves()
    if err != nil {
        return err
    }

    return nil
}

func (s *IncrementalTransition) undo() error {
    err := s.p.setCurrent(s.oldCurrent)
    if err != nil {
        return err
    }

    err = s.p.setWinner(s.oldWinner)
    if err != nil {
        return err
    }

    err = s.p.setGameOver(s.oldGameOver)
    if err != nil {
        return err
    }

    for _, color := range s.eliminated {
        err = s.p.restore(color)
        if err != nil {
            return err
        }
        err = s.b.disablePieces(color, false)
        if err != nil {
            return err
        }
    }

    err = s.b.CalculateMoves()
    if err != nil {
        return err
    }

    return nil
}

