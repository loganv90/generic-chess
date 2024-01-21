package chess

import "fmt"

/*
Responsible for:
- keeping track of the board, playerCollection, and the invoker
*/
type Game interface {
    // these are for the hub
	Execute(xFrom int, yFrom int, xTo int, yTo int, promotion string) error // called when a player tries to make a move
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

    p, err := createSimplePlayerCollectionWithDefaultPlayers()
    if err != nil {
        return nil, err
    }

	i, err := invokerFactoryInstance.newSimpleInvoker()
	if err != nil {
		return nil, err
	}

	return &SimpleGame{
		b: b,
        p: p,
		i: i,
	}, nil
}

func NewSimpleFourPlayerGame() (Game, error) {
    b, err := createSimpleFourPlayerBoardWithDefaultPieceLocations()
    if err != nil {
        return nil, err
    }

    p, err := createSimpleFourPlayerPlayerCollectionWithDefaultPlayers()
    if err != nil {
        return nil, err
    }

    i, err := invokerFactoryInstance.newSimpleInvoker()
    if err != nil {
        return nil, err
    }

    return &SimpleGame{
        b: b,
        p: p,
        i: i,
    }, nil
}

type SimpleGame struct {
	b Board
    p PlayerCollection
	i Invoker
}

func (s *SimpleGame) State() (*BoardData, error) {
    boardData := s.b.State()

    currentPlayer, err := s.p.getCurrent()
    if err != nil {
        return nil, err
    }

    boardData.CurrentPlayer = currentPlayer

    winningPlayer, err := s.p.getWinner()
    if err != nil {
        return nil, err
    }

    boardData.WinningPlayer = winningPlayer

    return boardData, nil
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int, promotion string) error {
    fromLocation := &Point{xFrom, yFrom}
    toLocation := &Point{xTo, yTo}

    moves, err := s.b.ValidMoves(fromLocation)
    if err != nil {
        return err
    }

    move := getMoveFromSlice(moves, toLocation, promotion)
    if move == nil {
        return fmt.Errorf("move not possible")
    }

    err = s.i.execute(move, s.b, s.p)
    if err != nil {
        return err
    }

    return nil
}

func getMoveFromSlice(moves []Move, toLocation *Point, promotion string) Move {
	for _, m := range moves {
        actionToLocation := m.getAction().toLocation
        if actionToLocation.equals(toLocation) {
            if promotionMove, ok := m.(*PromotionMove); ok {
                if _, ok := promotionMove.promotionPiece.(*Queen); ok && promotion == "Q" {
                    return m
                } else if _, ok := promotionMove.promotionPiece.(*Rook); ok && promotion == "R" {
                    return m
                } else if _, ok := promotionMove.promotionPiece.(*Bishop); ok && promotion == "B" {
                    return m
                } else if _, ok := promotionMove.promotionPiece.(*Knight); ok && promotion == "N" {
                    return m
                }
            } else {
                return m
            }
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

    currentPlayer, err := s.p.getCurrent()
    if err == nil && currentPlayer == piece.getColor() {
        moves, err := s.b.ValidMoves(location)
        if err != nil {
            return &PieceState{
                X: x,
                Y: y,
                Moves: []*MoveData{},
                Turn: false,
            }, nil
        }

        moveSet := make(map[MoveData]bool)
        moveDatas := make([]*MoveData, 0)
        for _, m := range moves {
            _, ok := m.(*PromotionMove)
            moveData := MoveData{
                X: m.getAction().toLocation.x,
                Y: m.getAction().toLocation.y,
                P: ok,
            }

            if _, ok := moveSet[moveData]; !ok {
                moveSet[moveData] = true
                moveDatas = append(moveDatas, &moveData)
            }
        }

        return &PieceState{
            X: x,
            Y: y,
            Moves: moveDatas,
            Turn: true,
        }, nil
    } else {
        _, _ = s.b.PotentialMoves(location)
        return &PieceState{
            X: x,
            Y: y,
            Moves: []*MoveData{},
            Turn: false,
        }, nil
    }
}

func (s *SimpleGame) Undo() error {
    err := s.i.undo()
    if err != nil {
        return err
    }

    return nil
}

func (s *SimpleGame) Redo() error {
    err := s.i.redo()
    if err != nil {
        return err
    }

    return nil
}

func (s *SimpleGame) Print() string {
	return s.b.Print()
}

