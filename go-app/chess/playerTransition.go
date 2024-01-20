package chess

var playerTransitionFactoryInstance = PlayerTransitionFactory(&ConcretePlayerTransitionFactory{})

type PlayerTransitionFactory interface {
	newIncrementalTransition(b Board, p PlayerCollection) (*IncrementalTransition, error)
}

type ConcretePlayerTransitionFactory struct{}

func (f *ConcretePlayerTransitionFactory) newIncrementalTransition(b Board, p PlayerCollection) (*IncrementalTransition, error) {
    oldColor, err := p.getCurrent()
    if err != nil {
        return nil, err
    }

    newColor := ""
    winner := ""
    eliminatedColors := []string{}
    for {
        newPlayer, err := p.getNext()
        if err != nil {
            return nil, err
        }

        if newPlayer.color == oldColor {
            newColor = newPlayer.color
            winner = newPlayer.color
            break
        }

        b.CalculateMoves(newPlayer.color)

        if b.Checkmate() {
            eliminatedColors = append(eliminatedColors, newPlayer.color)
            continue
        } else if b.Stalemate() {
            newColor = newPlayer.color
            winner = "draw"
            break
        } else {
            newColor = newPlayer.color
            break
        }
    }

    return &IncrementalTransition{
        p: p,
        b: b,
        oldColor: oldColor,
        newColor: newColor,
        winner: winner,
        eliminatedColors: eliminatedColors,
    }, nil
}

type PlayerTransition interface {
	execute() error
	undo() error
}

type IncrementalTransition struct {
    p PlayerCollection
    b Board
    oldColor string
    newColor string
    winner string
    eliminatedColors []string
}

func (s *IncrementalTransition) execute() error {
    s.p.setCurrent(s.newColor)

    for _, color := range s.eliminatedColors {
        s.p.eliminate(color)
    }

    s.b.CalculateMoves(s.newColor)

    return nil
}

func (s *IncrementalTransition) undo() error {
    s.p.setCurrent(s.oldColor)

    for _, color := range s.eliminatedColors {
        s.p.restore(color)
    }

    s.b.CalculateMoves(s.oldColor)

    return nil
}

