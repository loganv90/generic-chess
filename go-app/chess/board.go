package chess

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Point struct {
    x int
    y int
}

func (p *Point) equals(other *Point) bool {
    return p.x == other.x && p.y == other.y
}

func (p *Point) add(other *Point) *Point {
    return &Point{p.x + other.x, p.y + other.y}
}

type EnPassant struct {
    target *Point
    pieceLocation *Point
}

type Board interface {
	getPiece(location *Point) (Piece, error)
	setPiece(location *Point, piece Piece) error
	getEnPassant(color string) (*EnPassant, error)
	setEnPassant(color string, enPassant *EnPassant)
	clrEnPassant(color string)
    possibleEnPassants(color string, target *Point) []*EnPassant
	moves(location *Point) []Move
    increment()
    decrement()
    xLen() int
	yLen() int
	print() string
    turn() string
    squares() [][]*SquareData
    checkmate(color string) bool
    check(color string) bool
    pointOutOfBounds(p *Point) bool
    pointOnPromotionSquare(p *Point) bool
}

func newSimpleBoard(players []string, fen string) (*SimpleBoard, error) {
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

	return &SimpleBoard{
		currentPlayer: 0,
		players:       players,
		enPassantMap:  map[string]*EnPassant{},
		pieces:        pieces,
		halfMoveClock: halfMoveClock,
		fullMoveClock: fullMoveClock,
	}, nil
}

func createPiecesFromFen(fenRows string) [][]Piece {
	fenRowsSplit := strings.Split(fenRows, "/")
	pieceRows := [][]Piece{}

	for _, row := range fenRowsSplit {
		pieces := []Piece{}

		for _, char := range row {
			if unicode.IsDigit(char) {
				for i := 0; i < int(char-'0'); i++ {
					pieces = append(pieces, nil)
				}
			} else {
				pieces = append(pieces, createPieceFromChar(char))
			}
		}

		pieceRows = append(pieceRows, pieces)
	}

	return pieceRows
}

func createPieceFromChar(char rune) Piece {
	switch char {
	case 'r':
		if rook, err := newRook("black", false); err != nil {
			return nil
		} else {
			return rook
		}
	case 'n':
		if knight, err := newKnight("black"); err != nil {
			return nil
		} else {
			return knight
		}
	case 'b':
		if bishop, err := newBishop("black"); err != nil {
			return nil
		} else {
			return bishop
		}
	case 'q':
		if queen, err := newQueen("black"); err != nil {
			return nil
		} else {
			return queen
		}
	case 'k':
		if king, err := newKing("black", false, 0, 1); err != nil {
			return nil
		} else {
			return king
		}
	case 'p':
		if pawn, err := newPawn("black", false, 0, 1); err != nil {
			return nil
		} else {
			return pawn
		}
	case 'R':
		if rook, err := newRook("white", false); err != nil {
			return nil
		} else {
			return rook
		}
	case 'N':
		if knight, err := newKnight("white"); err != nil {
			return nil
		} else {
			return knight
		}
	case 'B':
		if bishop, err := newBishop("white"); err != nil {
			return nil
		} else {
			return bishop
		}
	case 'Q':
		if queen, err := newQueen("white"); err != nil {
			return nil
		} else {
			return queen
		}
	case 'K':
		if king, err := newKing("white", false, 0, -1); err != nil {
			return nil
		} else {
			return king
		}
	case 'P':
		if pawn, err := newPawn("white", false, 0, -1); err != nil {
			return nil
		} else {
			return pawn
		}
	default:
		return nil
	}
}

type SimpleBoard struct {
	currentPlayer int
	players       []string
	enPassantMap  map[string]*EnPassant
	pieces        [][]Piece
	halfMoveClock int
	fullMoveClock int
}

func (s *SimpleBoard) pointOutOfBounds(p *Point) bool {
    if p.y < 0 || p.y >= len(s.pieces) || p.x < 0 || p.x >= len(s.pieces[p.y]) {
        return true
    }
    return false
}

func (s *SimpleBoard) pointOnPromotionSquare(p *Point) bool {
    if p.y == 0 || p.y == len(s.pieces)-1 || p.x == 0 || p.x == len(s.pieces[p.y])-1 {
        return true
    }
    return false
}

func (s *SimpleBoard) getPiece(location *Point) (Piece, error) {
    if s.pointOutOfBounds(location) {
        return nil, fmt.Errorf("point out of bounds")
    }

	return s.pieces[location.y][location.x], nil
}

func (s *SimpleBoard) setPiece(location *Point, p Piece) error {
    if s.pointOutOfBounds(location) {
        return fmt.Errorf("point out of bounds")
    }

	s.pieces[location.y][location.x] = p
	return nil
}

func (s *SimpleBoard) getEnPassant(color string) (*EnPassant, error) {
	en, ok := s.enPassantMap[color]
	if !ok {
		return nil, nil
	}

	return en, nil
}

func (s *SimpleBoard) setEnPassant(color string, enPassant *EnPassant) {
	s.enPassantMap[color] = enPassant
}

func (s *SimpleBoard) clrEnPassant(color string) {
	delete(s.enPassantMap, color)
}

func (s *SimpleBoard) possibleEnPassants(color string, target *Point) []*EnPassant {
	ens := []*EnPassant{}

	for k, v := range s.enPassantMap {
		if k != color && target.equals(v.target) {
			ens = append(ens, v)
		}
	}

	return ens
}

func (s *SimpleBoard) check(color string) bool {
    kingLocation := &Point{-1, -1}
    attackerLocations := []*Point{}

    for y, row := range s.pieces {
        for x, piece := range row {
            if piece == nil {
                continue
            }

            if piece.getColor() == color {
			    if _, ok := piece.(*King); ok {
                    kingLocation = &Point{x, y}
                }
            } else {
                attackerLocations = append(attackerLocations, &Point{x, y})
            }
        }
    }

    if kingLocation == (&Point{-1, -1}) {
        return false
    }

    for _, attackerLocation := range attackerLocations {
        attackerPiece := s.pieces[attackerLocation.y][attackerLocation.x]
        if attackerPiece == nil {
            continue
        }

        attackerMoves := attackerPiece.moves(s, attackerLocation)
        for _, attackerMove := range attackerMoves {
            if attackerMove.getAction().toLocation == kingLocation {
                return true
            }
        }
    }


    return true
}

func (s *SimpleBoard) checkmate(color string) bool {
    return false
}

func (s *SimpleBoard) moves(location *Point) []Move {
	piece := s.pieces[location.y][location.x]

	if piece != nil {
		return piece.moves(s, location)
	}

	return []Move{}
}

func (s *SimpleBoard) increment() {
	s.currentPlayer = (s.currentPlayer + 1) % len(s.players)
	s.halfMoveClock++
	if s.currentPlayer == 0 {
		s.fullMoveClock++
	}
}

func (s *SimpleBoard) decrement() {
	s.currentPlayer = (s.currentPlayer - 1) % len(s.players)
	s.halfMoveClock--
	if s.currentPlayer == len(s.players)-1 {
		s.fullMoveClock--
	}
}

func (s *SimpleBoard) xLen() int {
    if s.yLen() <= 0 {
        return 0
    }
	return len(s.pieces[0])
}

func (s *SimpleBoard) yLen() int {
	return len(s.pieces)
}

func (s *SimpleBoard) print() string {
	var builder strings.Builder
	var cellWidth int = 12

	builder.WriteString(fmt.Sprintf("Player: %s\n", s.players[s.currentPlayer]))
	builder.WriteString(fmt.Sprintf("Check:  %t\n", false))
	builder.WriteString(fmt.Sprintf("Mate:   %t\n", false))
	for y, row := range s.pieces {
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*s.xLen()-1)))
		for x := range row {
			builder.WriteString(fmt.Sprintf("|%s%2dx ", strings.Repeat(" ", cellWidth-4), x))
		}
		builder.WriteString("|\n")
		for _, piece := range row {
			if piece == nil {
				builder.WriteString(fmt.Sprintf("|%s", strings.Repeat(" ", cellWidth)))
			} else {
				p := piece.print()
				if len(p) > 1 {
					p = p[:1]
				}

				pColor := piece.getColor()
				if len(pColor) > 8 {
					pColor = pColor[:8]
				}

				builder.WriteString(fmt.Sprintf("| %-1s %-8s ", p, pColor))
			}
		}
		builder.WriteString("|\n")
		for range row {
			builder.WriteString(fmt.Sprintf("|%s%2dy ", strings.Repeat(" ", cellWidth-4), y))
		}
		builder.WriteString("|\n")
	}
	builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*s.xLen()-1)))

	return builder.String()
}

func (s *SimpleBoard) turn() string {
    return s.players[s.currentPlayer]
}

func (s *SimpleBoard) squares() [][]*SquareData {
    squares := [][]*SquareData{}

    for y, row := range s.pieces {
        squaresRow := []*SquareData{}
        for x := range row {
            piece := s.pieces[y][x]
            if piece != nil {
                squaresRow = append(squaresRow, &SquareData{
                    Piece: piece.print(),
                    Color: piece.getColor(),
                })
            } else {
                squaresRow = append(squaresRow, &SquareData{})
            }
        }
        squares = append(squares, squaresRow)
    }

    return squares
}

