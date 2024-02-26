package chess

var playerTransitionFactoryInstance = PlayerTransitionFactory(&ConcretePlayerTransitionFactory{})

type PlayerTransitionFactory interface {
	newIncrementalTransitionAsPlayerTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (PlayerTransition, error)
    newIncrementalTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (*IncrementalTransition, error)
}

type ConcretePlayerTransitionFactory struct{}

func (f *ConcretePlayerTransitionFactory) newIncrementalTransitionAsPlayerTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (PlayerTransition, error) {
    return f.newIncrementalTransition(b, p, inCheckmate, inStalemate)
}

func (f *ConcretePlayerTransitionFactory) newIncrementalTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (*IncrementalTransition, error) {
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

    nextPlayers, err := p.getNext()
    if err != nil {
        return nil, err
    }

    var newCurrent string
    var newWinner string
    var newGameOver bool

    if inStalemate {
        newCurrent = oldCurrent
        newWinner = ""
        newGameOver = true
    } else if len(nextPlayers) < 1 {
        newCurrent = oldCurrent
        newWinner = oldWinner
        newGameOver = true
    } else if len(nextPlayers) == 1 {
        newCurrent = nextPlayers[0].color
        newWinner = nextPlayers[0].color
        newGameOver = true
    } else if len(nextPlayers) == 2 {
        if inCheckmate {
            newCurrent = nextPlayers[0].color
            newWinner = nextPlayers[0].color
            newGameOver = true
        } else {
            newCurrent = nextPlayers[0].color
            newWinner = ""
            newGameOver = false
        }
    } else {
        newCurrent = nextPlayers[0].color
        newWinner = ""
        newGameOver = false
    }

    return &IncrementalTransition{
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
    oldGameOver bool
    newGameOver bool
    eliminated bool
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

    if !s.eliminated {
        return nil
    }

    err = s.p.eliminate(s.oldCurrent)
    if err != nil {
        return err
    }

    err = s.b.disablePieces(s.oldCurrent, true)
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

    if !s.eliminated {
        return nil
    }

    err = s.p.restore(s.oldCurrent)
    if err != nil {
        return err
    }

    err = s.b.disablePieces(s.oldCurrent, false)
    if err != nil {
        return err
    }

    return nil
}

