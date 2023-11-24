package chess

import "fmt"

type BoardState struct {
    Squares [][]*SquareData
    Turn string
	Check bool
	Mate  bool
}

type SquareData struct {
    C string
    P string
}

type PieceState struct {
    Moves []*MoveData
    Turn bool
}

type MoveData struct {
    X int
    Y int
}

type Game interface {
	Execute(xFrom int, yFrom int, xTo int, yTo int) (*BoardState, error)
    View(x int, y int) (*PieceState, error)
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
	b Board
	i Invoker
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int) (*BoardState, error) {
    fromLocation := &Point{xFrom, yFrom}
    toLocation := &Point{xTo, yTo}

	moves := s.b.moves(fromLocation)
	move := getMoveFromSlice(moves, toLocation)
	if move == nil {
		return nil, fmt.Errorf("move not possible")
	}

	err := s.i.execute(move)
	if err != nil {
		return nil, err
	}

    _, check, checkmate, stalemate := s.b.allMoves(s.b.turn())

	return &BoardState{
        Squares: s.b.squares(),
        Turn: s.b.turn(),
		Check: check,
		Mate:  checkmate || stalemate,
	}, nil
}

func (s *SimpleGame) View(x int, y int) (*PieceState, error) {
    location := &Point{x, y}

    piece, err := s.b.getPiece(location)
    if err != nil || piece == nil {
        return &PieceState{
            Moves: []*MoveData{},
            Turn: false,
        }, nil
    }

    moves := s.b.moves(location)
    moveDatas := make([]*MoveData, len(moves))
    for i, m := range moves {
        moveDatas[i] = &MoveData{
            X: m.getAction().toLocation.x,
            Y: m.getAction().toLocation.y,
        }
    }

    return &PieceState{
        Moves: moveDatas,
        Turn: s.b.turn() == piece.getColor(),
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

