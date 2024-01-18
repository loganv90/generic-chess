package chess

import "fmt"

/*
Responsible for:
- keeping track of the board and the invoker to execute moves
- keeping track of the players in the game
*/
type Game interface {
    // these are for the hub
	Execute(xFrom int, yFrom int, xTo int, yTo int, promotion string) error // called when a player tries to make a move
    State() (*BoardData, error) // called to get the game state
    View(xFrom int, yFrom int) (*PieceState, error) // show valid moves to current player and show all moves to others
    Player(color string) (*Player, error) // get player by color
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
        winningPlayer: -1,
        players: []*Player{
            {
                color: "white",
                alive: true,
            },
            {
                color: "black",
                alive: true,
            },
        },
	}, nil
}

func NewSimpleFourPlayerGame() (Game, error) {
    b, err := createSimpleFourPlayerBoardWithDefaultPieceLocations()
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
        winningPlayer: -1,
        players: []*Player{
            {
                color: "white",
                alive: true,
            },
            {
                color: "red",
                alive: true,
            },
            {
                color: "black",
                alive: true,
            },
            {
                color: "blue",
                alive: true,
            },
        },
    }, nil
}

// TODO we might want to make an addition to the invoker class such that we can store player state
type SimpleGame struct {
	b Board
	i Invoker
    currentPlayer int
    winningPlayer int
    players []*Player
}

func (s *SimpleGame) State() (*BoardData, error) {
    boardData := s.b.State()

    currentPlayer, err := s.getCurrentPlayer()
    if err == nil {
        boardData.CurrentPlayer = currentPlayer.color
    }

    winningPlayer, err := s.getWinningPlayer()
    if err == nil {
        boardData.WinningPlayer = winningPlayer.color
    }

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

    err = s.i.execute(move)
    if err != nil {
        return err
    }

    s.increment()
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

    currentPlayer, err := s.getCurrentPlayer()
    if err == nil && currentPlayer.color == piece.getColor() {
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

func (s *SimpleGame) Player(color string) (*Player, error) {
    for _, p := range s.players {
        if p.color == color {
            return p, nil
        }
    }

    return nil, fmt.Errorf("player not found")
}

func (s *SimpleGame) getCurrentPlayer() (*Player, error) {
    if s.currentPlayer < 0 || s.currentPlayer >= len(s.players) {
        return nil, fmt.Errorf("no current player")
    }

    return s.players[s.currentPlayer], nil
}

func (s *SimpleGame) getWinningPlayer() (*Player, error) {
    if s.winningPlayer < 0 || s.winningPlayer >= len(s.players) {
        return nil, fmt.Errorf("no winning player")
    }

    return s.players[s.winningPlayer], nil
}

func (s *SimpleGame) Undo() error {
    err := s.i.undo()
    if err != nil {
        return err
    }

    s.decrement()
    return nil
}

func (s *SimpleGame) Redo() error {
    err := s.i.redo()
    if err != nil {
        return err
    }

    s.increment()
    return nil
}

func (s *SimpleGame) incrementPlayer() string {
    s.currentPlayer = (s.currentPlayer + 1) % len(s.players)
    if s.currentPlayer < 0 {
        s.currentPlayer = len(s.players) - 1
    }
    return s.players[s.currentPlayer].color
}

func (s *SimpleGame) decrementPlayer() string {
    s.currentPlayer = (s.currentPlayer - 1) % len(s.players)
    if s.currentPlayer < 0 {
        s.currentPlayer = len(s.players) - 1
    }
    return s.players[s.currentPlayer].color
}

// TODO we're going to have to determine the winning player in the increment/decrement stuff
func (s *SimpleGame) decrement() {
    if s.b.Checkmate() || s.b.Stalemate() {
        s.players[s.currentPlayer].alive = true
    }

    currentPlayerColor := s.decrementPlayer()
    s.b.CalculateMoves(currentPlayerColor)
}

func (s *SimpleGame) increment() {
    currentPlayerColor := s.incrementPlayer()
    s.b.CalculateMoves(currentPlayerColor)

    if s.b.Checkmate() || s.b.Stalemate() {
        s.players[s.currentPlayer].alive = false
    }
}

func (s *SimpleGame) Print() string {
	return s.b.Print()
}

