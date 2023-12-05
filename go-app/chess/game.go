package chess

import "fmt"

/*
Responsible for:
- keeping track of the board and the invoker to execute moves
- keeping track of the players in the game
*/
type Game interface {
    // these are for the hub
	Execute(xFrom int, yFrom int, xTo int, yTo int) error // called when a player tries to make a move
    State() (*BoardData, error) // called to get the game state
    View(xFrom int, yFrom int) (*PieceState, error) // show valid moves to current player and show all moves to others
	Undo() error
	Redo() error
	Print() string
}

func NewSimpleGame() (Game, error) {
    b, err := createSimpleBoardWithDefaultPieceLocations()
	if err != nil {
		return nil, err
	}
    b.CalculateMoves("white")

	i, err := invokerFactoryInstance.newSimpleInvoker()
	if err != nil {
		return nil, err
	}

	return &SimpleGame{
		b: b,
		i: i,
        currentPlayer: 0,
        players: []string{"white", "black"},
	}, nil
}

type SimpleGame struct {
	b Board
	i Invoker
    currentPlayer int
    players []string
}

func (s *SimpleGame) State() (*BoardData, error) {
    return s.b.State(), nil
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int) error {
    fromLocation := &Point{xFrom, yFrom}
    toLocation := &Point{xTo, yTo}

    moves, err := s.b.ValidMoves(fromLocation)
    if err != nil {
        return err
    }

    move := getMoveFromSlice(moves, toLocation)
    if move == nil {
        return fmt.Errorf("move not possible")
    }

    err = s.i.execute(move)
    if err != nil {
        return err
    }

    s.increment()
    return nil
}

func getMoveFromSlice(moves []Move, toLocation *Point) Move {
	for _, m := range moves {
        actionToLocation := m.getAction().toLocation
        if actionToLocation.equals(toLocation) {
			return m
		}
	}

	return nil
}

func (s *SimpleGame) View(x int, y int) (*PieceState, error) {
    location := &Point{x, y}

    piece, err := s.b.getPiece(location)
    if err != nil || piece == nil {
        return &PieceState{
            X: x,
            Y: y,
            Moves: []*MoveData{},
            Turn: false,
        }, nil
    }

    if s.turn() == piece.getColor() {
        moves, err := s.b.ValidMoves(location)
        if err != nil {
            return &PieceState{
                X: x,
                Y: y,
                Moves: []*MoveData{},
                Turn: false,
            }, nil
        }

        moveDatas := make([]*MoveData, len(moves))
        for i, m := range moves {
            moveDatas[i] = &MoveData{
                X: m.getAction().toLocation.x,
                Y: m.getAction().toLocation.y,
            }
        }

        return &PieceState{
            X: x,
            Y: y,
            Moves: moveDatas,
            Turn: true,
        }, nil
    } else {
        moves, err := s.b.PotentialMoves(location)
        if err != nil {
            return &PieceState{
                X: x,
                Y: y,
                Moves: []*MoveData{},
                Turn: false,
            }, nil
        }

        moveDatas := make([]*MoveData, len(moves))
        for i, m := range moves {
            moveDatas[i] = &MoveData{
                X: m.getAction().toLocation.x,
                Y: m.getAction().toLocation.y,
            }
        }

        return &PieceState{
            X: x,
            Y: y,
            Moves: moveDatas,
            Turn: false,
        }, nil
    }
}

func (s *SimpleGame) turn() string {
    return s.players[s.currentPlayer]
}

func (s *SimpleGame) Undo() error {
    err := s.i.undo()
    if err != nil {
        return err
    }

    s.decrement()
    return nil
}

func (s *SimpleGame) decrement() {
    s.currentPlayer = (s.currentPlayer - 1) % len(s.players)
    s.b.CalculateMoves(s.turn())
}

func (s *SimpleGame) Redo() error {
    err := s.i.redo()
    if err != nil {
        return err
    }

    s.increment()
    return nil
}

func (s *SimpleGame) increment() {
    s.currentPlayer = (s.currentPlayer + 1) % len(s.players)
    s.b.CalculateMoves(s.turn())
}

func (s *SimpleGame) Print() string {
	return s.b.Print()
}

