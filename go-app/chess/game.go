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
    Moves(color int) ([]MoveKey, error) // get all valid moves
	Undo() error
	Redo() error
	Print() string
    Copy() (Game, error)

    getBoard() Board
    getPlayerCollection() PlayerCollection
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

    currentPlayer, _ := s.p.getCurrent()

    boardData.CurrentPlayer = currentPlayer

    winningPlayer, _ := s.p.getWinner()

    boardData.WinningPlayer = winningPlayer

    gameOver, err := s.p.getGameOver()
    if err != nil {
        return nil, err
    }

    boardData.GameOver = gameOver

    return boardData, nil
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int, promotion string) error {
    fromLocation := Point{xFrom, yFrom}
    toLocation := Point{xTo, yTo}

    gameOver, err := s.p.getGameOver()
    if err != nil || gameOver {
        return fmt.Errorf("game is over")
    }

    var move Move

    moves, err := s.b.LegalMovesOfLocation(fromLocation)
    if err != nil {
        return err
    }

    for _, m := range moves {
        action := m.getAction()
        if action.fromLocation == fromLocation && action.toLocation == toLocation {
            move = m
            break
        }
    }

    if move == nil {
        return fmt.Errorf("invalid move")
    }

    if _, ok := move.(*AllyDefenseMove); ok {
        return fmt.Errorf("AllyDefenseMove not possible")
    }

    if promotionMove, ok := move.(*PromotionMove); ok {
        color := promotionMove.getNewPiece().getColor()

        var promotionPiece Piece
        if promotion == "Q" {
            promotionPiece = newQueen(color)
        } else if promotion == "R" {
            promotionPiece = newRook(color, true)
        } else if promotion == "B" {
            promotionPiece = newBishop(color)
        } else if promotion == "N" {
            promotionPiece = newKnight(color)
        } else {
            return fmt.Errorf("invalid promotion piece")
        }

        promotionMove.setPromotionPiece(promotionPiece)
    }

    transition, err := s.p.GetTransition(s.b, false, false)
    if err != nil {
        return err
    }

    err = s.i.execute(move, transition)
    if err != nil {
        return err
    }

    err = s.b.CalculateMoves()
    if err != nil {
        return err
    }

    for {
        currentPlayer, _ := s.p.getCurrent()

        checkmate, stalemate, err := s.b.CheckmateAndStalemate(currentPlayer)
        if err != nil {
            return err
        }

        if checkmate {
            transition, err := s.p.GetTransition(s.b, true, false)
            if err != nil {
                return err
            }

            err = s.i.execute(nil, transition)
            if err != nil {
                return err
            }

            continue
        } else if stalemate {
            transition, err := s.p.GetTransition(s.b, false, true)
            if err != nil {
                return err
            }

            err = s.i.execute(nil, transition)
            if err != nil {
                return err
            }
        }

        break
    }

    return nil
}

func (s *SimpleGame) View(x int, y int) (*PieceState, error) {
    location := Point{x, y}

    piece, ok := s.b.getPiece(location)
    if !ok || piece == nil {
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

    currentPlayer, _ := s.p.getCurrent()
    if currentPlayer == piece.getColor() {
        moves, err := s.b.LegalMovesOfLocation(location)
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

func (s *SimpleGame) Moves(color int) ([]MoveKey, error) {
    moveKeys := make([]MoveKey, 0)
    
    moves, err := s.b.LegalMovesOfColor(color)
    if err != nil {
        return nil, err
    }

    for _, move := range moves {
        moveKeys = append(moveKeys, MoveKey{
            XFrom: move.getAction().fromLocation.x,
            YFrom: move.getAction().fromLocation.y,
            XTo: move.getAction().toLocation.x,
            YTo: move.getAction().toLocation.y,
            Promotion: "",
        })
    }

    return moveKeys, nil
}

func (s *SimpleGame) Undo() error {
    err := s.i.undo()
    if err != nil {
        return err
    }

    err = s.b.CalculateMoves()
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

    err = s.b.CalculateMoves()
    if err != nil {
        return err
    }

    return nil
}

func (s *SimpleGame) Print() string {
	return s.b.Print()
}

func (s *SimpleGame) Copy() (Game, error) {
    newBoard, err := s.b.Copy()
    if err != nil {
        return nil, err
    }

    newPlayerCollection, err := s.p.Copy()
    if err != nil {
        return nil, err
    }

    newInvoker, err := s.i.Copy()
    if err != nil {
        return nil, err
    }

    return &SimpleGame{
        b: newBoard,
        p: newPlayerCollection,
        i: newInvoker,
    }, nil
}

func (s *SimpleGame) getBoard() Board {
    return s.b
}

func (s *SimpleGame) getPlayerCollection() PlayerCollection {
    return s.p
}

