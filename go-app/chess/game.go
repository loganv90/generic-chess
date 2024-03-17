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

    currentPlayer := s.p.getCurrent()
    boardData.CurrentPlayer = currentPlayer

    winningPlayer := s.p.getWinner()
    boardData.WinningPlayer = winningPlayer

    gameOver := s.p.getGameOver()
    boardData.GameOver = gameOver

    return boardData, nil
}

func (s *SimpleGame) Execute(xFrom int, yFrom int, xTo int, yTo int, promotion string) error {
    fromLocation := s.b.getIndex(xFrom, yFrom)
    if fromLocation == nil {
        return fmt.Errorf("invalid from location")
    }

    toLocation := s.b.getIndex(xTo, yTo)
    if toLocation == nil {
        return fmt.Errorf("invalid to location")
    }

    gameOver := s.p.getGameOver()
    if gameOver {
        return fmt.Errorf("game is over")
    }

    moves, err := s.b.LegalMovesOfLocation(fromLocation)
    if err != nil {
        return err
    }

    found := false
    var move FastMove
    for _, m := range moves {
        promotionString := ""
        if m.promotionIndex >= 0 {
            promotionString = piece_names[m.promotionIndex]
        }

        if m.fromLocation == fromLocation && m.toLocation == toLocation && promotionString == promotion {
            move = m
            found = true
            break
        }

        if m.fromLocation == fromLocation && m.toLocation == toLocation && promotionString == "" {
            move = m
            found = true
        }
    }

    if !found {
        return fmt.Errorf("invalid move")
    }

    if move.allyDefense {
        return fmt.Errorf("AllyDefenseMove not possible")
    }

    transition := PlayerTransition{}
    createPlayerTransition(s.b, s.p, false, false, &transition)

    err = s.i.execute(move, transition)
    if err != nil {
        return err
    }

    s.b.CalculateMoves()

    for {
        currentPlayer := s.p.getCurrent()

        checkmate, stalemate, err := s.b.CheckmateAndStalemate(currentPlayer)
        if err != nil {
            return err
        }

        if checkmate {
            createPlayerTransition(s.b, s.p, true, false, &transition)

            err = s.i.executeHalf(transition)
            if err != nil {
                return err
            }

            continue
        } else if stalemate {
            createPlayerTransition(s.b, s.p, false, true, &transition)

            err = s.i.executeHalf(transition)
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

    piece := s.b.getPiece(&location)
    if piece == nil {
        return &PieceState{
            X: x,
            Y: y,
            Moves: []*MoveData{},
            Turn: false,
        }, nil
    }

    gameOver := s.p.getGameOver()
    if gameOver {
        return &PieceState{
            X: x,
            Y: y,
            Moves: []*MoveData{},
            Turn: false,
        }, nil
    }

    currentPlayer := s.p.getCurrent()
    if currentPlayer == piece.color {
        moves, err := s.b.LegalMovesOfLocation(&location)
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
            moveData := MoveData{
                X: m.toLocation.x,
                Y: m.toLocation.y,
                P: m.promotionIndex >= 0,
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
            XFrom: move.fromLocation.x,
            YFrom: move.fromLocation.y,
            XTo: move.toLocation.x,
            YTo: move.toLocation.y,
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

    s.b.CalculateMoves()

    return nil
}

func (s *SimpleGame) Redo() error {
    err := s.i.redo()
    if err != nil {
        return err
    }

    s.b.CalculateMoves()

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

