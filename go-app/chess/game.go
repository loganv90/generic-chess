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
    View(xFrom int, yFrom int) (*PieceState, error) // show valid moves of piece
    Moves() ([]*MoveKey, error) // get all valid moves
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


    gameOver, err := s.p.getGameOver()
    if err != nil {
        return nil, err
    }

    boardData.GameOver = gameOver

    return boardData, nil
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int, promotion string) error {
    fromLocation := &Point{xFrom, yFrom}
    toLocation := &Point{xTo, yTo}

    gameOver, err := s.p.getGameOver()
    if err != nil || gameOver {
        return fmt.Errorf("game is over")
    }

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
                if p, ok := promotionMove.promotionPiece.(*Queen); ok && promotion == p.print() {
                    return m
                } else if p, ok := promotionMove.promotionPiece.(*Rook); ok && promotion == p.print() {
                    return m
                } else if p, ok := promotionMove.promotionPiece.(*Bishop); ok && promotion == p.print() {
                    return m
                } else if p, ok := promotionMove.promotionPiece.(*Knight); ok && promotion == p.print() {
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

    gameOver, err := s.p.getGameOver()
    if err != nil || gameOver {
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
        return &PieceState{
            X: x,
            Y: y,
            Moves: []*MoveData{},
            Turn: false,
        }, nil
    }
}

func (s *SimpleGame) Moves() ([]*MoveKey, error) {
    // TODO find a new way to do this
    return s.b.AvailableMoves("white")
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

