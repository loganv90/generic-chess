package chess

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type enPassant struct {
	xTarget int
	yTarget int
	xPiece  int
	yPiece  int
}

type board interface {
	getPiece(x int, y int) (piece, error)
	setPiece(x int, y int, piece piece) error
	getEnPassant(color string) (*enPassant, error)
	setEnPassant(color string, enPassant *enPassant)
	clrEnPassant(color string)
	possibleEnPassants(color string, xTarget int, yTarget int) []*enPassant
	moves(x int, y int) []move
	increment()
	decrement()
	xLen() int
	yLen() int
}

func newSimpleBoard(players []string, fen string) (*simpleBoard, error) {
	if len(players) < 2 {
		return nil, fmt.Errorf("not enough players")
	}

	fenSplit := strings.Split(fen, " ")
	if len(fenSplit) != 6 {
		return nil, fmt.Errorf("invalid fen")
	}

	pieces := createPiecesFromFen(fenSplit[0])

	halfMoveClock, err := strconv.Atoi(fenSplit[4])
	if err != nil {
		return nil, fmt.Errorf("invalid half move clock")
	}

	fullMoveClock, err := strconv.Atoi(fenSplit[5])
	if err != nil {
		return nil, fmt.Errorf("invalid full move clock")
	}

	return &simpleBoard{
		currentPlayer: 0,
		players:       players,
		enPassantMap:  map[string]*enPassant{},
		pieces:        pieces,
		halfMoveClock: halfMoveClock,
		fullMoveClock: fullMoveClock,
	}, nil
}

func createPiecesFromFen(fenRows string) [][]piece {
	f := &concreteMoveFactory{}
	fenRowsSplit := strings.Split(fenRows, "/")
	pieceRows := [][]piece{}

	for _, row := range fenRowsSplit {
		pieces := []piece{}

		for _, char := range row {
			if unicode.IsDigit(char) {
				for i := 0; i < int(char); i++ {
					pieces = append(pieces, nil)
				}
			} else {
				pieces = append(pieces, createPieceFromChar(f, char))
			}
		}

		pieceRows = append(pieceRows, pieces)
	}

	return pieceRows
}

func createPieceFromChar(f moveFactory, char rune) piece {
	switch char {
	case 'r':
		return &rook{}
	case 'n':
		return &knight{}
	case 'b':
		return &bishop{}
	case 'q':
		return &queen{}
	case 'k':
		return &king{}
	case 'p':
		if pawn, err := newPawn(f, "black", false, 0, -1); err != nil {
			return nil
		} else {
			return pawn
		}
	case 'R':
		return &rook{}
	case 'N':
		return &knight{}
	case 'B':
		return &bishop{}
	case 'Q':
		return &queen{}
	case 'K':
		return &king{}
	case 'P':
		if pawn, err := newPawn(f, "white", false, 0, -1); err != nil {
			return nil
		} else {
			return pawn
		}
	default:
		return nil
	}
}

type simpleBoard struct {
	currentPlayer int
	players       []string
	enPassantMap  map[string]*enPassant
	pieces        [][]piece
	halfMoveClock int
	fullMoveClock int
}

func (s *simpleBoard) getPiece(x int, y int) (piece, error) {
	if y < 0 || y >= len(s.pieces) || x < 0 || x >= len(s.pieces[y]) {
		return nil, fmt.Errorf("coordinates out of bounds")
	}

	return s.pieces[y][x], nil
}

func (s *simpleBoard) setPiece(x int, y int, p piece) error {
	if y < 0 || y >= len(s.pieces) || x < 0 || x >= len(s.pieces[y]) {
		return fmt.Errorf("coordinates out of bounds")
	}

	s.pieces[y][x] = p

	return nil
}

func (s *simpleBoard) getEnPassant(color string) (*enPassant, error) {
	en, ok := s.enPassantMap[color]
	if !ok {
		return en, fmt.Errorf("en passant not found")
	}

	return en, nil
}

func (s *simpleBoard) setEnPassant(color string, enPassant *enPassant) {
	s.enPassantMap[color] = enPassant
}

func (s *simpleBoard) clrEnPassant(color string) {
	delete(s.enPassantMap, color)
}

func (s *simpleBoard) possibleEnPassants(color string, xTarget int, yTarget int) []*enPassant {
	ens := []*enPassant{}

	for k, v := range s.enPassantMap {
		if k != color && v.xTarget == xTarget && v.yTarget == yTarget {
			ens = append(ens, v)
		}
	}

	return ens
}

func (s *simpleBoard) moves(x int, y int) []move {
	piece := s.pieces[y][x]

	if piece != nil {
		return piece.moves(s, x, y)
	}

	return []move{}
}

func (s *simpleBoard) increment() {
	s.currentPlayer = (s.currentPlayer + 1) % len(s.players)
	s.halfMoveClock++
	if s.currentPlayer == 0 {
		s.fullMoveClock++
	}
}

func (s *simpleBoard) decrement() {
	s.currentPlayer = (s.currentPlayer - 1) % len(s.players)
	s.halfMoveClock--
	if s.currentPlayer == len(s.players)-1 {
		s.fullMoveClock--
	}
}

func (s *simpleBoard) xLen() int {
	return len(s.pieces[0])
}

func (s *simpleBoard) yLen() int {
	return len(s.pieces)
}
