package chess

import (
	"fmt"
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

type PieceLocations struct {
    ownPieceLocations []*Point
    enemyPieceLocations []*Point
}

type Board interface {
	getPiece(location *Point) (Piece, error)
	setPiece(location *Point, piece Piece) error
    getKingLocation(color string) (*Point, error)
    setKingLocation(color string, location *Point)
    getVulnerables(color string) ([]*Point, error)
    setVulnerables(color string, locations []*Point)
	getEnPassant(color string) (*EnPassant, error)
	setEnPassant(color string, enPassant *EnPassant)
	clrEnPassant(color string)
    possibleEnPassants(color string, target *Point) []*EnPassant
	moves(location *Point) []Move
    allMoves(color string) ([]Move, bool, bool, bool)
    increment()
    decrement()
    xLen() int
	yLen() int
	print() string
    turn() string
    squares() [][]*SquareData
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

    kingLocationMap := map[string]*Point{}
    for y, row := range pieces {
        for x, piece := range row {
            if piece != nil {
                if _, ok := piece.(*King); ok {
                    kingLocationMap[piece.getColor()] = &Point{x, y}
                }
            }
        }
    }

	return &SimpleBoard{
		currentPlayer: 0,
		players:       players,
		pieces:        pieces,
        kingLocationMap: kingLocationMap,
		enPassantMap:  map[string]*EnPassant{},
        vulnerablesMap: map[string][]*Point{},
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
	pieces        [][]Piece
    kingLocationMap map[string]*Point
	enPassantMap  map[string]*EnPassant
    vulnerablesMap map[string][]*Point // locations that should not be attacked by enemy pieces (for castling)
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

// TODO - keep track of pieces by color and keep them updated in the get/set piece methods
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

    if _, ok := p.(*King); ok {
        s.setKingLocation(p.getColor(), location)
    }

	s.pieces[location.y][location.x] = p

	return nil
}

func (s *SimpleBoard) getKingLocation(color string) (*Point, error) {
    kingLocation, ok := s.kingLocationMap[color]
    if !ok {
        return nil, fmt.Errorf("king location not found")
    }

    return kingLocation, nil
}

func (s *SimpleBoard) setKingLocation(color string, location *Point) {
    s.kingLocationMap[color] = location
}

func (s *SimpleBoard) getVulnerables(color string) ([]*Point, error) {
    vulnerables, ok := s.vulnerablesMap[color]
    if !ok {
        return nil, nil
    }

    return vulnerables, nil
}

func (s *SimpleBoard) setVulnerables(color string, locations []*Point) {
    s.vulnerablesMap[color] = locations
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

func (s *SimpleBoard) moves(location *Point) []Move {
	piece := s.pieces[location.y][location.x]

    if piece != nil {
        return piece.moves(s, location)
    }

    return []Move{}
}

// TODO - get the pieces from maps instead of iterating over the board
func (s *SimpleBoard) allMoves(color string) ([]Move, bool, bool, bool) {
    moves := []Move{}
    pieceLocations := s.getPieceLocations(color)
    check := s.isInCheck(color, pieceLocations.enemyPieceLocations)

    for _, ownPieceLocation := range pieceLocations.ownPieceLocations {
        piece, err := s.getPiece(ownPieceLocation)
        if piece == nil || err != nil {
            continue
        }

        pieceMoves := piece.moves(s, ownPieceLocation)
        for _, move := range pieceMoves {
            boardCopy := s.copy()
            move.getAction().b = boardCopy
            move.execute()

            if boardCopy.isInCheck(color, pieceLocations.enemyPieceLocations) {
                continue
            }

            moves = append(moves, move)
        }
    }

    checkmate := check && len(moves) == 0
    stalemate := !check && len(moves) == 0

    return moves, check, checkmate, stalemate
}

func (s *SimpleBoard) isInCheck(color string, ememyPieceLocations []*Point) bool {
    vulnerableLocations, err := s.getVulnerables(color)
    if err == nil {
        for _, vulnerableLocation := range vulnerableLocations {
            if s.isSquareAttacked(vulnerableLocation, ememyPieceLocations) {
                return true
            }
        }
    }

    kingLocation, err := s.getKingLocation(color)
    if err == nil {
        if s.isSquareAttacked(kingLocation, ememyPieceLocations) {
            return true
        }
    }

    return false
}

func (s *SimpleBoard) isSquareAttacked(squareLocation *Point, pieceLocations []*Point) bool {
    for _, pieceLocation := range pieceLocations {
        piece, err := s.getPiece(pieceLocation)
        if piece == nil || err != nil {
            continue
        }

        moves := piece.moves(s, pieceLocation)
        for _, move := range moves {
            if move.getAction().toLocation.equals(squareLocation) {
                return true
            }
        }
    }

    return false
}

func (s *SimpleBoard) getPieceLocations(color string) *PieceLocations {
    ownPieceLocations := []*Point{}
    enemyPieceLocations := []*Point{}

    for y, row := range s.pieces {
        for x, piece := range row {
            if piece == nil {
                continue
            }

            if piece.getColor() == color {
                ownPieceLocations = append(ownPieceLocations, &Point{x, y})
            } else {
                enemyPieceLocations = append(enemyPieceLocations, &Point{x, y})
            }
        }
    }

    return &PieceLocations{
        ownPieceLocations: ownPieceLocations,
        enemyPieceLocations: enemyPieceLocations,
    }
}

func (s *SimpleBoard) increment() {
	s.currentPlayer = (s.currentPlayer + 1) % len(s.players)
}

func (s *SimpleBoard) decrement() {
	s.currentPlayer = (s.currentPlayer - 1) % len(s.players)
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

	builder.WriteString(fmt.Sprintf("Player: %s\n", s.turn()))
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
                    P: piece.print(),
                    C: piece.getColor(),
                })
            } else {
                squaresRow = append(squaresRow, &SquareData{})
            }
        }
        squares = append(squares, squaresRow)
    }

    return squares
}

func (s *SimpleBoard) copy() *SimpleBoard {
    pieces := [][]Piece{}
    for _, row := range s.pieces {
        piecesRow := []Piece{}
        for _, piece := range row {
            piecesRow = append(piecesRow, piece)
        }
        pieces = append(pieces, piecesRow)
    }

    enPassantMap := map[string]*EnPassant{}
    for k, v := range s.enPassantMap {
        enPassantMap[k] = v
    }

    players := []string{}
    for _, player := range s.players {
        players = append(players, player)
    }

    return &SimpleBoard{
        currentPlayer: s.currentPlayer,
        players: players,
        pieces: pieces,
        kingLocationMap: s.kingLocationMap,
        enPassantMap: enPassantMap,
        vulnerablesMap: s.vulnerablesMap,
    }
}

