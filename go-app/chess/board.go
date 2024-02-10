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
    disableLocation(location *Point) error
    getVulnerables(color string) ([]*Point, error) // if these locations are attacked, the player is in check
    setVulnerables(color string, locations []*Point) error
	getEnPassant(color string) (*EnPassant, error) // if these locations are attacked, a piece is captured en passant
	setEnPassant(color string, enPassant *EnPassant) error
    possibleEnPassant(color string, location *Point) ([]*EnPassant, error)
    clearEnPassant(color string) error

    // these are for the playerTransition
    disablePieces(color string, disable bool) error

    // these are for the game
    CalculateMoves() error // calcutes moves for every color
    ValidMoves(fromLocation *Point) ([]Move, error) // returns moves from a location
    AvailableMoves(color string) ([]*MoveKey, error) // returns moves for a color
    Size() *Point
	Print() string
    State() *BoardData
    Check(color string) bool
    Checkmate(color string) bool
    Stalemate(color string) bool
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
        disabledLocations: map[Point]bool{},
        kingLocationMap: map[string]*Point{},
        pieceLocationsMap: map[string][]*Point{},
		enPassantMap: map[string]*EnPassant{},
        vulnerablesMap: map[string][]*Point{},
        fromToToToMoveMap: map[Point]map[Point]Move{},
        toToFromToMoveMap: map[Point]map[Point]Move{},
        checkMap: map[string]bool{},
        checkmateMap: map[string]bool{},
        stalemateMap: map[string]bool{},
	}, nil
}

type SimpleBoard struct {
    size *Point
	pieces [][]Piece
    disabledLocations map[Point]bool
    kingLocationMap map[string]*Point
    pieceLocationsMap map[string][]*Point
	enPassantMap map[string]*EnPassant
    vulnerablesMap map[string][]*Point
    fromToToToMoveMap map[Point]map[Point]Move
    toToFromToMoveMap map[Point]map[Point]Move
    checkMap map[string]bool
    checkmateMap map[string]bool
    stalemateMap map[string]bool
}

func (s *SimpleBoard) disablePieces(color string, disable bool) error {
    for _, pieceLocation := range s.pieceLocationsMap[color] {
        piece, err := s.getPiece(pieceLocation)
        if err != nil {
            continue
        }
        piece.setDisabled(disable)
    }
    return nil
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
    if err == nil && oldPiece != nil {
        if pieceLocations, ok := s.pieceLocationsMap[oldPiece.getColor()]; ok {
            removeIndex := -1
            for i, pieceLocation := range pieceLocations {
                if pieceLocation.equals(location) {
                    removeIndex = i
                    break
                }
            }
            if removeIndex != -1 {
                pieceLocations[removeIndex] = pieceLocations[len(pieceLocations)-1]
                pieceLocations[len(pieceLocations)-1] = nil
                pieceLocations = pieceLocations[:len(pieceLocations)-1]
                s.pieceLocationsMap[oldPiece.getColor()] = pieceLocations
            }
        }
    }

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

	s.pieces[location.y][location.x] = p
    return nil
}

func (s *SimpleBoard) disableLocation(location *Point) error {
    s.disabledLocations[*location] = true
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
        if v == nil {
            continue
        }
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

func (s *SimpleBoard) ValidMoves(fromLocation *Point) ([]Move, error) {
    moves := []Move{}
    
    for _, move := range s.fromToToToMoveMap[*fromLocation] {
        if _, ok := move.(*AllyDefenseMove); ok {
            continue
        }
        moves = append(moves, move)
    }

    return moves, nil
}

func (s *SimpleBoard) AvailableMoves(color string) ([]*MoveKey, error) {
    moveKeys := []*MoveKey{}

    for _, pieceLocation := range s.pieceLocationsMap[color] {
        for _, move := range s.fromToToToMoveMap[*pieceLocation] { 
            if _, ok := move.(*AllyDefenseMove); ok {
                continue
            }

            action := move.getAction()
            promotion := ""
            if promotionMove, ok := move.(*PromotionMove); ok {
                promotion = promotionMove.promotionPiece.print()
            }
            moveKeys = append(moveKeys, &MoveKey{
                action.fromLocation.x,
                action.fromLocation.y,
                action.toLocation.x,
                action.toLocation.y,
                promotion,
            })
        }
    }
    return moveKeys, nil
}

// TODO add 3 move repetition and 50 move rule
// TODO add rule to allow checks and only lose on king capture
// TODO add rule to check for checkmate and stalemate on all players after every move
func (s *SimpleBoard) CalculateMoves() error {
    colorToMoveCountMap := map[string]int{}
    s.fromToToToMoveMap = map[Point]map[Point]Move{}
    s.toToFromToMoveMap = map[Point]map[Point]Move{}

    for _, pieceLocations := range s.pieceLocationsMap {
        for _, fromLocation := range pieceLocations {
            piece, err := s.getPiece(fromLocation)
            if piece == nil || err != nil {
                continue
            }

            moves := piece.moves(s, fromLocation)
            for _, move := range moves {
                action := move.getAction()
                color := move.getNewPiece().getColor()

                if _, ok := s.fromToToToMoveMap[*action.fromLocation]; !ok {
                    s.fromToToToMoveMap[*action.fromLocation] = map[Point]Move{}
                }
                if _, ok := s.toToFromToMoveMap[*action.toLocation]; !ok {
                    s.toToFromToMoveMap[*action.toLocation] = map[Point]Move{}
                }
                if _, ok := colorToMoveCountMap[color]; !ok {
                    colorToMoveCountMap[color] = 0
                }

                s.fromToToToMoveMap[*action.fromLocation][*action.toLocation] = move
                s.toToFromToMoveMap[*action.toLocation][*action.fromLocation] = move
                colorToMoveCountMap[color] += 1
            }
        }
    }

    boardCopy, err := s.copy()
    if err != nil {
        return err
    }

    for _, toLocations := range s.fromToToToMoveMap {
        for _, move := range toLocations {
            color := move.getNewPiece().getColor()
            if _, ok := move.(*AllyDefenseMove); ok {
                colorToMoveCountMap[color] -= 1
                continue
            }

            action := move.getAction()
            action.b = boardCopy

            err = move.execute()
            if err != nil {
                return err
            }

            fromLocationsToExclude := boardCopy.getFromLocationsGivenToLocations([]*Point{action.toLocation, action.fromLocation})
            toLocationsToInclude := boardCopy.calculateToLocationsGivenFromLocations(color, fromLocationsToExclude)
            if boardCopy.isInCheck(color, fromLocationsToExclude, toLocationsToInclude) {
                delete(s.fromToToToMoveMap[*action.fromLocation], *action.toLocation)
                delete(s.toToFromToMoveMap[*action.toLocation], *action.fromLocation)
                colorToMoveCountMap[color] -= 1
            }

            err = move.undo()
            if err != nil {
                return err
            }

            action.b = s
        }
    }

    for color, moveCount := range colorToMoveCountMap {
        s.checkMap[color] = s.isInCheck(color, map[Point]bool{}, map[Point]bool{})
        s.checkmateMap[color] = s.checkMap[color] && moveCount <= 0
        s.stalemateMap[color] = !s.checkMap[color] && moveCount <= 0
    }
    
    return nil
}

func (s *SimpleBoard) calculateToLocationsGivenFromLocations(color string, fromLocations map[Point]bool) map[Point]bool {
    toLocations := map[Point]bool{}

    for fromLocation := range fromLocations {
        piece, err := s.getPiece(&fromLocation)
        if piece == nil || err != nil || piece.getColor() == color {
            continue
        }

        moves := piece.moves(s, &fromLocation)
        for _, move := range moves {
            action := move.getAction()
            toLocations[*action.toLocation] = true
        }
    }

    return toLocations
}

func (s *SimpleBoard) getFromLocationsGivenToLocations(toLocations []*Point) map[Point]bool {
    fromLocations := map[Point]bool{}

    for _, toLocation := range toLocations {
        for fromLocation := range s.toToFromToMoveMap[*toLocation] {
            fromLocations[fromLocation] = true
        }
    }

    return fromLocations
}

func (s *SimpleBoard) isLocationAttacked(color string, fromLocationsToExclude map[Point]bool, toLocationsToInclude map[Point]bool, location *Point) bool {
    if _, ok2 := toLocationsToInclude[*location]; ok2 {
        return true
    }

    if fromToMoveMap, ok2 := s.toToFromToMoveMap[*location]; ok2 {
        for fromLocation, move := range fromToMoveMap {
            if _, ok3 := fromLocationsToExclude[fromLocation]; ok3 {
                continue
            }

            if move.getNewPiece().getColor() != color {
                return true
            }
        }
    }
    
    return false
}

func (s *SimpleBoard) isInCheck(color string, fromLocationsToExclude map[Point]bool, toLocationsToInclude map[Point]bool) bool {
    if kingLocation, ok1 := s.kingLocationMap[color]; ok1 {
        if s.isLocationAttacked(color, fromLocationsToExclude, toLocationsToInclude, kingLocation) {
            return true
        }
    }

    if vulnerableLocations, ok1 := s.vulnerablesMap[color]; ok1 {
        for _, vulnerableLocation := range vulnerableLocations {
            if s.isLocationAttacked(color, fromLocationsToExclude, toLocationsToInclude, vulnerableLocation) {
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
		for x, piece := range row {
            if s.pointOutOfBounds(&Point{x, y}) {
                builder.WriteString(fmt.Sprintf("|%s", strings.Repeat("X", cellWidth)))
            } else if piece == nil {
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

func (s *SimpleBoard) State() *BoardData {
    pieces := []*PieceData{}

    for y, row := range s.pieces {
        for x := range row {
            piece := s.pieces[y][x]
            if piece != nil {
                disabled, _ := piece.getDisabled()
                pieces = append(pieces, &PieceData{
                    T: piece.print(),
                    C: piece.getColor(),
                    X: x,
                    Y: y,
                    D: disabled,
                })
            }
        }
    }

    disabled := []*DisabledData{}
    for disabledLocation := range s.disabledLocations {
        disabled = append(disabled, &DisabledData{
            X: disabledLocation.x,
            Y: disabledLocation.y,
        })
    }

    return &BoardData{
        XSize: s.size.x,
        YSize: s.size.y,
        Disabled: disabled,
        Pieces: pieces,
    }
}

func (s *SimpleBoard) Check(color string) bool {
    if check, ok := s.checkMap[color]; ok {
        return check
    }
    return false
}

func (s *SimpleBoard) Checkmate(color string) bool {
    if checkmate, ok := s.checkmateMap[color]; ok {
        return checkmate
    }
    return false
}

func (s *SimpleBoard) Stalemate(color string) bool {
    if stalemate, ok := s.stalemateMap[color]; ok {
        return stalemate
    }
    return false
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

    for k, v := range s.disabledLocations {
        simpleBoard.disabledLocations[k] = v
    }

    for k, v := range s.enPassantMap {
        simpleBoard.enPassantMap[k] = v
    }

    for k, v := range s.fromToToToMoveMap {
        for k2, v2 := range v {
            if _, ok := simpleBoard.fromToToToMoveMap[k]; !ok {
                simpleBoard.fromToToToMoveMap[k] = map[Point]Move{}
            }
            simpleBoard.fromToToToMoveMap[k][k2] = v2
        }
    }

    for k, v := range s.toToFromToMoveMap {
        for k2, v2 := range v {
            if _, ok := simpleBoard.toToFromToMoveMap[k]; !ok {
                simpleBoard.toToFromToMoveMap[k] = map[Point]Move{}
            }
            simpleBoard.toToFromToMoveMap[k][k2] = v2
        }
    }

    return simpleBoard, nil
}

func (s *SimpleBoard) pointOutOfBounds(p *Point) bool {
    _, ok := s.disabledLocations[*p]
    return ok || p.y < 0 || p.y >= s.size.y || p.x < 0 || p.x >= s.size.x
}

