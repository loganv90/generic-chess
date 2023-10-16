package chess

import "fmt"

type Status struct {
	Check bool
	Mate  bool
}

type Game interface {
	Execute(xFrom int, yFrom int, xTo int, yTo int) (*Status, error)
	Undo() error
	Redo() error
	print() string
}

func NewSimpleGame() (Game, error) {
	b, err := newSimpleBoard(
		[]string{"white", "black"},
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	)
	if err != nil {
		return nil, err
	}

	i, err := invokerFactoryInstance.newSimpleInvoker()
	if err != nil {
		return nil, err
	}

	return &SimpleGame{
		b: b,
		i: i,
	}, nil
}

type SimpleGame struct {
	b board
	i invoker
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int) (*Status, error) {
	moves := s.b.moves(xFrom, yFrom)
	move := getMoveFromSlice(moves, xTo, yTo)
	if move == nil {
		return &Status{}, fmt.Errorf("move not possible")
	}

	err := s.i.execute(move)
	if err != nil {
		return &Status{}, err
	}

	return &Status{
		Check: false,
		Mate:  false,
	}, nil
}

func (s *SimpleGame) Undo() error {
	return s.i.undo()
}

func (s *SimpleGame) Redo() error {
	return s.i.redo()
}

func (s *SimpleGame) print() string {
	return s.b.print()
}