package chess

import "fmt"

type State struct {
    Squares [][]SquareData
    Turn string
	Check bool
	Mate  bool
}

type SquareData struct {
    Color string
    Piece string
}

type Game interface {
	Execute(xFrom int, yFrom int, xTo int, yTo int) (*State, error)
	Undo() error
	Redo() error
	Print() string
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

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int) (*State, error) {
	moves := s.b.moves(xFrom, yFrom)
	move := getMoveFromSlice(moves, xTo, yTo)
	if move == nil {
		return nil, fmt.Errorf("move not possible")
	}

	err := s.i.execute(move)
	if err != nil {
		return nil, err
	}

	return &State{
        Squares: s.b.squares(),
        Turn: s.b.turn(),
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

func (s *SimpleGame) Print() string {
	return s.b.print()
}

