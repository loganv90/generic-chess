package chess

import (
	"fmt"
	"strings"
)

/*
Responsible for:
- keeping track of the pieces on the board
- keeping track of availalbe moves for all pieces
*/
type Board interface {
    // these are for the move
	getPiece(location *Point) (Piece, error)
	setPiece(location *Point, piece Piece) error
    getVulnerables(color string) ([]*Point, error) // if these locations are attacked, the player is in check
    setVulnerables(color string, locations []*Point) error
	getEnPassant(color string) (*EnPassant, error) // if these locations are attacked, a piece is captured en passant
	setEnPassant(color string, enPassant *EnPassant) error
    possibleEnPassant(color string, location *Point) ([]*EnPassant, error)
    clearEnPassant(color string) error

    // these are for the game
    PotentialMoves(fromLocation *Point) ([]Move, error) // returns moves for a piece without considering other pieces
    ValidMoves(fromLocation *Point) ([]Move, error) // returns moves for a piece using the moveMap
    CalculateMoves(color string) error // calcutes moves assuming it is the color's turn and sets moveMap and state
    Size() *Point
	Print() string
    State() ([][]*SquareData, bool, bool, bool)

    // misc
    pointOutOfBounds(p *Point) bool
    pointOnPromotionSquare(p *Point) bool
}

func newSimpleBoard(size *Point) (*SimpleBoard, error) {
    if size.x <= 0 || size.y <= 0 {
        return nil, fmt.Errorf("invalid board size")
    }

    pieces := make([][]Piece, size.y)
    for p := range pieces {
        pieces[p] = make([]Piece, size.x)
    }

	return &SimpleBoard{
        size: size,
		pieces: pieces,
        kingLocationMap: map[string]*Point{},
        pieceLocationsMap: map[string][]*Point{},
		enPassantMap: map[string]*EnPassant{},
        vulnerablesMap: map[string][]*Point{},
        moveMap: map[Point][]Move{},
        check: false,
        checkmate: false,
        stalemate: false,
	}, nil
}

type SimpleBoard struct {
    size *Point
	pieces [][]Piece
    kingLocationMap map[string]*Point
    pieceLocationsMap map[string][]*Point
	enPassantMap map[string]*EnPassant
    vulnerablesMap map[string][]*Point
    moveMap map[Point][]Move
    check bool
    checkmate bool
    stalemate bool
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

    oldPiece, err := s.getPiece(location)
    if err != nil && oldPiece != nil {
        if pieceLocations, ok := s.pieceLocationsMap[oldPiece.getColor()]; ok {
            removeIndex := -1
            for i, pieceLocation := range pieceLocations {
                if pieceLocation.equals(location) {
                    removeIndex = i
                    break
                }
            }
            if removeIndex == -1 {
                return nil
            }

            pieceLocations[removeIndex] = pieceLocations[len(pieceLocations)-1]
            pieceLocations[len(pieceLocations)-1] = nil
            pieceLocations = pieceLocations[:len(pieceLocations)-1]
            s.pieceLocationsMap[oldPiece.getColor()] = pieceLocations
        }
    }

	s.pieces[location.y][location.x] = p
    if p != nil {
        if pieceLocations, ok := s.pieceLocationsMap[p.getColor()]; ok {
            pieceLocations = append(pieceLocations, location)
            s.pieceLocationsMap[p.getColor()] = pieceLocations
        } else {
            s.pieceLocationsMap[p.getColor()] = []*Point{location}
        }

        if _, ok := p.(*King); ok {
            s.kingLocationMap[p.getColor()] = location
        }
    }

    return nil
}

func (s *SimpleBoard) getVulnerables(color string) ([]*Point, error) {
    vulnerables, okVulnerables := s.vulnerablesMap[color]
    kingLocation, okKingLocation := s.kingLocationMap[color]

    if okVulnerables && okKingLocation {
        return append(vulnerables, kingLocation), nil
    } else if okVulnerables {
        return vulnerables, nil
    } else if okKingLocation {
        return []*Point{kingLocation}, nil
    } else {
        return []*Point{}, nil
    }
}

func (s *SimpleBoard) setVulnerables(color string, locations []*Point) error {
    s.vulnerablesMap[color] = locations
    return nil
}

func (s *SimpleBoard) getEnPassant(color string) (*EnPassant, error) {
	en, ok := s.enPassantMap[color]
	if !ok {
		return nil, nil
	}

	return en, nil
}

func (s *SimpleBoard) setEnPassant(color string, enPassant *EnPassant) error {
	s.enPassantMap[color] = enPassant
    return nil
}

func (s *SimpleBoard) possibleEnPassant(color string, target *Point) ([]*EnPassant, error) {
    enPassants := []*EnPassant{}

	for k, v := range s.enPassantMap {
		if k != color && target.equals(v.target) {
            enPassants = append(enPassants, v)
		}
	}

    return enPassants, nil
}

func (s *SimpleBoard) clearEnPassant(color string) error {
    delete(s.enPassantMap, color)
    return nil
}

// TODO - we want the moves assuming empty board
func (s *SimpleBoard) PotentialMoves(fromLocation *Point) ([]Move, error) {
    return []Move{}, nil
}

func (s *SimpleBoard) ValidMoves(fromLocation *Point) ([]Move, error) {
    moves, ok := s.moveMap[*fromLocation]
    if !ok {
        return []Move{}, nil
    }

    return moves, nil
}

func (s *SimpleBoard) CalculateMoves(color string) error {
    s.moveMap = map[Point][]Move{}
    s.check = false
    s.checkmate = false
    s.stalemate = false

    ownPieceLocations, ok := s.pieceLocationsMap[color]
    if !ok {
        return fmt.Errorf("no pieces found for color %s", color)
    }

    enemyPieceLocations := []*Point{}
    for k, v := range s.pieceLocationsMap {
        if k != color {
            enemyPieceLocations = append(enemyPieceLocations, v...)
        }
    }

    s.check = s.isInCheck(color, enemyPieceLocations)

    for _, ownPieceLocation := range ownPieceLocations {
        piece, err := s.getPiece(ownPieceLocation)
        if piece == nil || err != nil {
            continue
        }
        
        moves := piece.moves(s, ownPieceLocation)
        for _, move := range moves {
            boardCopy, err := s.copy()
            if err != nil {
                continue
            }
            move.getAction().b = boardCopy

            move.execute()
            if boardCopy.isInCheck(color, enemyPieceLocations) {
                continue
            }

            move.getAction().b = s
            s.moveMap[*ownPieceLocation] = append(s.moveMap[*ownPieceLocation], move)
        }
    }

    s.checkmate = s.check && len(s.moveMap) == 0
    s.stalemate = !s.check && len(s.moveMap) == 0

    return nil
}

func (s *SimpleBoard) isInCheck(color string, ememyPieceLocations []*Point) bool {
    vulnerableLocations, err := s.getVulnerables(color)
    if err == nil {
        for _, vulnerableLocation := range vulnerableLocations {
            if s.isLocationAttacked(vulnerableLocation, ememyPieceLocations) {
                return true
            }
        }
    }

    return false
}

func (s *SimpleBoard) isLocationAttacked(location *Point, pieceLocations []*Point) bool {
    for _, pieceLocation := range pieceLocations {
        piece, err := s.getPiece(pieceLocation)
        if piece == nil || err != nil {
            continue
        }

        moves := piece.moves(s, pieceLocation)
        for _, move := range moves {
            if move.getAction().toLocation.equals(location) {
                return true
            }
        }
    }

    return false
}

func (s *SimpleBoard) Size() *Point {
    return s.size
}

func (s *SimpleBoard) Print() string {
	var builder strings.Builder
	var cellWidth int = 12

	for y, row := range s.pieces {
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*s.size.x-1)))
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
	builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*s.size.x-1)))

	return builder.String()
}

func (s *SimpleBoard) State() ([][]*SquareData, bool, bool, bool) {
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

    return squares, s.check, s.checkmate, s.stalemate
}

func (s *SimpleBoard) copy() (*SimpleBoard, error) {
    simpleBoard, err := newSimpleBoard(s.size)
    if err != nil {
        return nil, err
    }

    for y, row := range s.pieces {
        for x := range row {
            piece := s.pieces[y][x]
            if piece != nil {
                simpleBoard.setPiece(&Point{x, y}, piece.copy())
            }
        }
    }

    for k, v := range s.kingLocationMap {
        simpleBoard.kingLocationMap[k] = v
    }

    for k, v := range s.enPassantMap {
        simpleBoard.enPassantMap[k] = v
    }

    for k, v := range s.vulnerablesMap {
        simpleBoard.vulnerablesMap[k] = v
    }

    return simpleBoard, nil
}

func (s *SimpleBoard) pointOutOfBounds(p *Point) bool {
    return p.y < 0 || p.y >= s.size.y || p.x < 0 || p.x >= s.size.x
}

func (s *SimpleBoard) pointOnPromotionSquare(p *Point) bool {
    return false
}

