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
	getPiece(location Point) (Piece, bool)
	setPiece(location Point, piece Piece) bool
    disableLocation(location Point) error
    getVulnerables(color string) ([]Point, error) // if these locations are attacked, the player is in check
    setVulnerables(color string, locations []Point) error
	getEnPassant(color string) (*EnPassant, error) // if these locations are attacked, a piece is captured en passant
	setEnPassant(color string, enPassant *EnPassant) error
    possibleEnPassant(color string, location Point) ([]*EnPassant, error)
    clearEnPassant(color string) error

    // these are for the playerTransition
    disablePieces(color string, disable bool) error

    // these are for the bot
    getPieceLocations() map[string][]Point

    // these are for the game
    CalculateMoves() error // calcutes moves for every color
    CalculateMovesPartial(move Move) error // recalculates moves affectedy by a move
    ValidMoves(fromLocation Point) ([]Move, error) // returns moves from a location
    AvailableMoves(color string) ([]Move, error) // returns moves for a color
    LegalMoves(color string) ([]Move, error)
    Move(fromLocation Point, toLocation Point, promotion string) (Move, error)
    State() *BoardData
    Check(color string) bool
    Checkmate(color string) bool
    Stalemate(color string) bool
    Mate(color string) (bool, bool, error)

	Print() string
    Copy() (Board, error) 
    UniqueString() string
}

func newSimpleBoard(size Point) (*SimpleBoard, error) {
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
        kingLocationMap: map[string]Point{},
        pieceLocationsMap: map[string][]Point{},
		enPassantMap: map[string]*EnPassant{},
        vulnerablesMap: map[string][]Point{},
        fromToToToMoveMap: map[Point]map[Point]Move{},
        toToFromToMoveMap: map[Point]map[Point]Move{},
        checkMap: map[string]bool{},
        checkmateMap: map[string]bool{},
        stalemateMap: map[string]bool{},
        test: false,
	}, nil
}

type SimpleBoard struct {
    size Point
	pieces [][]Piece
    disabledLocations map[Point]bool
    kingLocationMap map[string]Point
    pieceLocationsMap map[string][]Point
	enPassantMap map[string]*EnPassant
    vulnerablesMap map[string][]Point
    fromToToToMoveMap map[Point]map[Point]Move
    toToFromToMoveMap map[Point]map[Point]Move
    checkMap map[string]bool
    checkmateMap map[string]bool
    stalemateMap map[string]bool
    test bool
}

func (s *SimpleBoard) disablePieces(color string, disable bool) error {
    for _, pieceLocation := range s.pieceLocationsMap[color] {
        piece, ok := s.getPiece(pieceLocation)
        if !ok {
            continue
        }
        piece.setDisabled(disable)
    }
    return nil
}

func (s *SimpleBoard) getPiece(location Point) (Piece, bool) {
    if s.pointOutOfBounds(location) {
        return nil, false
    }

	return s.pieces[location.y][location.x], true
}

func (s *SimpleBoard) setPiece(location Point, p Piece) bool {
    if s.pointOutOfBounds(location) {
        return false
    }

    oldPiece, ok := s.getPiece(location)
    if ok && oldPiece != nil {
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
                pieceLocations[len(pieceLocations)-1] = Point{}
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
            s.pieceLocationsMap[p.getColor()] = []Point{location}
        }

        if _, ok := p.(*King); ok {
            s.kingLocationMap[p.getColor()] = location
        }
    }

	s.pieces[location.y][location.x] = p
    return true
}

func (s *SimpleBoard) disableLocation(location Point) error {
    s.disabledLocations[location] = true
    return nil
}

func (s *SimpleBoard) getVulnerables(color string) ([]Point, error) {
    vulnerables, okVulnerables := s.vulnerablesMap[color]
    kingLocation, okKingLocation := s.kingLocationMap[color]

    if okVulnerables && okKingLocation {
        return append(vulnerables, kingLocation), nil
    } else if okVulnerables {
        return vulnerables, nil
    } else if okKingLocation {
        return []Point{kingLocation}, nil
    } else {
        return []Point{}, nil
    }
}

func (s *SimpleBoard) setVulnerables(color string, locations []Point) error {
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

func (s *SimpleBoard) possibleEnPassant(color string, target Point) ([]*EnPassant, error) {
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

func (s *SimpleBoard) ValidMoves(fromLocation Point) ([]Move, error) {
    moves := []Move{}
    
    for _, move := range s.fromToToToMoveMap[fromLocation] {
        if _, ok := move.(*AllyDefenseMove); ok {
            continue
        }
        moves = append(moves, move)
    }

    return moves, nil
}

func (s *SimpleBoard) Move(fromLocation Point, toLocation Point, promotion string) (Move, error) {
    toToMoveMap, ok := s.toToFromToMoveMap[toLocation]
    if !ok {
        return nil, fmt.Errorf("move not possible")
    }

    move, ok := toToMoveMap[fromLocation]
    if !ok {
        return nil, fmt.Errorf("move not possible")
    }

    return move, nil
}

func (s *SimpleBoard) AvailableMoves(color string) ([]Move, error) {
    moves := []Move{}

    for _, pieceLocation := range s.pieceLocationsMap[color] {
        for _, move := range s.fromToToToMoveMap[pieceLocation] { 
            if _, ok := move.(*AllyDefenseMove); ok {
                continue
            }

            moves = append(moves, move)
        }
    }

    return moves, nil
}

// TODO add 3 move repetition and 50 move rule
// TODO add rule to allow checks and only lose on king capture
// TODO add rule to check for checkmate and stalemate on all players after every move
func (s *SimpleBoard) CalculateMoves() error {
    if s.test {
        return nil
    }

    colorToMoveCountMap := map[string]int{}
    s.fromToToToMoveMap = map[Point]map[Point]Move{}
    s.toToFromToMoveMap = map[Point]map[Point]Move{}

    for _, pieceLocations := range s.pieceLocationsMap {
        for _, fromLocation := range pieceLocations {
            piece, ok := s.getPiece(fromLocation)
            if piece == nil || !ok {
                continue
            }

            moves := piece.moves(s, fromLocation)
            for _, move := range moves {
                action := move.getAction()
                color := move.getNewPiece().getColor()

                if _, ok := s.fromToToToMoveMap[action.fromLocation]; !ok {
                    s.fromToToToMoveMap[action.fromLocation] = map[Point]Move{}
                }
                if _, ok := s.toToFromToMoveMap[action.toLocation]; !ok {
                    s.toToFromToMoveMap[action.toLocation] = map[Point]Move{}
                }
                if _, ok := colorToMoveCountMap[color]; !ok {
                    colorToMoveCountMap[color] = 0
                }

                s.fromToToToMoveMap[action.fromLocation][action.toLocation] = move
                s.toToFromToMoveMap[action.toLocation][action.fromLocation] = move
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

            fromLocationsToExclude := boardCopy.getFromLocationsGivenToLocations([]Point{action.toLocation, action.fromLocation})
            toLocationsToInclude := boardCopy.calculateToLocationsGivenFromLocations(color, fromLocationsToExclude)
            if boardCopy.isInCheck(color, fromLocationsToExclude, toLocationsToInclude) {
                delete(s.fromToToToMoveMap[action.fromLocation], action.toLocation)
                delete(s.toToFromToMoveMap[action.toLocation], action.fromLocation)
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

func (s *SimpleBoard) CalculateMovesPartial(move Move) error {
    // TODO implement dynamic move calculations based on previous move
    // TODO implement delayed check calculation system for bot
    action := move.getAction()
    relevantLocations := []Point{action.fromLocation, action.toLocation}

    fromLocationsToRecalculate := s.getFromLocationsGivenToLocations(relevantLocations)

    for fromLocation := range fromLocationsToRecalculate {
        s.fromToToToMoveMap[fromLocation] = map[Point]Move{}
        for toLocation := range s.toToFromToMoveMap[fromLocation] {
            delete(s.toToFromToMoveMap[toLocation], fromLocation)
        }
    }

    for fromLocation := range fromLocationsToRecalculate {
        piece, ok := s.getPiece(fromLocation)
        if piece == nil || !ok {
            continue
        }

        moves := piece.moves(s, fromLocation)
        for _, move := range moves {
            action := move.getAction()

            if _, ok := s.fromToToToMoveMap[action.fromLocation]; !ok {
                s.fromToToToMoveMap[action.fromLocation] = map[Point]Move{}
            }
            if _, ok := s.toToFromToMoveMap[action.toLocation]; !ok {
                s.toToFromToMoveMap[action.toLocation] = map[Point]Move{}
            }

            s.fromToToToMoveMap[action.fromLocation][action.toLocation] = move
            s.toToFromToMoveMap[action.toLocation][action.fromLocation] = move
        }
    }

    return nil
}

func (s *SimpleBoard) calculateToLocationsGivenFromLocations(color string, fromLocations map[Point]bool) map[Point]bool {
    toLocations := map[Point]bool{}

    for fromLocation := range fromLocations {
        piece, ok := s.getPiece(fromLocation)
        if piece == nil || !ok || piece.getColor() == color {
            continue
        }

        moves := piece.moves(s, fromLocation)
        for _, move := range moves {
            action := move.getAction()
            toLocations[action.toLocation] = true
        }
    }

    return toLocations
}

func (s *SimpleBoard) getFromLocationsGivenToLocations(toLocations []Point) map[Point]bool {
    fromLocations := map[Point]bool{}

    for _, toLocation := range toLocations {
        for fromLocation := range s.toToFromToMoveMap[toLocation] {
            fromLocations[fromLocation] = true
        }
    }

    return fromLocations
}

func (s *SimpleBoard) isLocationAttacked(color string, fromLocationsToExclude map[Point]bool, toLocationsToInclude map[Point]bool, location Point) bool {
    if _, ok2 := toLocationsToInclude[location]; ok2 {
        return true
    }

    if fromToMoveMap, ok2 := s.toToFromToMoveMap[location]; ok2 {
        for fromLocation := range fromToMoveMap {
            if _, ok3 := fromLocationsToExclude[fromLocation]; ok3 {
                continue
            }

            piece, ok := s.getPiece(fromLocation)
            if piece == nil || !ok {
                continue
            }

            if piece.getColor() != color {
                return true
            }
        }
    }
    
    return false
}

func (s *SimpleBoard) isInCheck(color string, fromLocationsToExclude map[Point]bool, toLocationsToInclude map[Point]bool) bool {
    // fromLocationsToExclude: contains the locations of the pieces that are attacking the toLocation or fromLocation of the recent move
    // toLocationsToInclude: contains the updated attacked locations of the pieces that were attacking the toLocation or fromLocation of the recent move

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
            if s.pointOutOfBounds(Point{x, y}) {
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
                disabled := piece.getDisabled()
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

func (s *SimpleBoard) Mate(color string) (bool, bool, error) {
    legalMoves, err := s.LegalMoves(color)
    if err != nil {
        return false, false, err
    }

    if len(legalMoves) > 0 {
        return false, false, nil
    }

    if s.Check2(color) {
        return true, false, nil
    }

    return false, true, nil
}

func (s *SimpleBoard) LegalMoves(color string) ([]Move, error) {
    boardCopy, err := s.copy()
    if err != nil {
        return []Move{}, err
    }

    pieceLocations, ok := s.pieceLocationsMap[color]
    if !ok {
        return []Move{}, fmt.Errorf("color not found")
    }

    legalMoves := []Move{}
    for _, fromLocation := range pieceLocations {
        toToMoveMap, ok := s.fromToToToMoveMap[fromLocation]
        if !ok {
            continue
        }

        for _, move := range toToMoveMap {
            // TODO copy the move so we're not editing the real moves
            if _, ok := move.(*AllyDefenseMove); ok {
                continue
            }

            action := move.getAction()
            action.b = boardCopy

            err = move.execute()
            if err != nil {
                return []Move{}, err
            }

            fromLocationsToExclude := boardCopy.getFromLocationsGivenToLocations([]Point{action.toLocation, action.fromLocation})
            toLocationsToInclude := boardCopy.calculateToLocationsGivenFromLocations(color, fromLocationsToExclude)
            if !boardCopy.isInCheck(color, fromLocationsToExclude, toLocationsToInclude) {
                legalMoves = append(legalMoves, move)
            }

            err = move.undo()
            if err != nil {
                return []Move{}, err
            }

            action.b = s
        }
    }

    return legalMoves, nil
}

func (s *SimpleBoard) Check2(color string) bool {
    // TODO this shouldn't return an error
    vulnerableLocations, err := s.getVulnerables(color)
    if err != nil {
        return false
    }

    for _, vulnerableLocation := range vulnerableLocations {
        if fromToMoveMap, ok := s.toToFromToMoveMap[vulnerableLocation]; ok {
            for fromLocation := range fromToMoveMap {
                piece, ok := s.getPiece(fromLocation)
                if piece == nil || !ok {
                    continue
                }

                if piece.getColor() != color {
                    return true
                }
            }
        }
    }

    return false
}

func (s *SimpleBoard) getPieceLocations() map[string][]Point {
    return s.pieceLocationsMap
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
                simpleBoard.setPiece(Point{x, y}, piece.copy())
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

func (s *SimpleBoard) Copy() (Board, error) {
    board, err := s.copy()
    return board, err
}

func (s *SimpleBoard) UniqueString() string {
    builder := strings.Builder{}

    counter := 0
    for y, row := range s.pieces {
        for x := range row {
            piece := s.pieces[y][x]

            if piece == nil {
                counter += 1
                continue
            }

            if counter > 0 {
                builder.WriteString(fmt.Sprintf("%d", counter))
                counter = 0
            }

            if piece.getDisabled() {
                builder.WriteString("d")
                continue
            }

            builder.WriteString(piece.print())
            builder.WriteString(piece.getColor())
            if !piece.getMoved() {
                builder.WriteString("m")
            }
        }
    }

    return builder.String()
}

func (s *SimpleBoard) pointOutOfBounds(p Point) bool {
    _, ok := s.disabledLocations[p]
    return ok || p.y < 0 || p.y >= s.size.y || p.x < 0 || p.x >= s.size.x
}

func newFastBoard(size Point) (*FastBoard, error) {
    if size.x <= 0 || size.y <= 0 {
        return nil, fmt.Errorf("invalid board size")
    }

    pieces := make([][]Piece, size.y)
    for p := range pieces {
        pieces[p] = make([]Piece, size.x)
    }

	return &FastBoard{
        size: size,
		pieces: pieces,
        disabledLocations: map[Point]bool{},
        kingLocationMap: map[string]Point{},
        pieceLocationsMap: map[string][]Point{},
		enPassantMap: map[string]*EnPassant{},
        vulnerablesMap: map[string][]Point{},
        fromToToToMoveMap: map[Point]map[Point]Move{},
        toToFromToMoveMap: map[Point]map[Point]Move{},
	}, nil
}

type FastBoard struct {
    size Point
	pieces [][]Piece
    disabledLocations map[Point]bool
    kingLocationMap map[string]Point
    pieceLocationsMap map[string][]Point
	enPassantMap map[string]*EnPassant
    vulnerablesMap map[string][]Point
    fromToToToMoveMap map[Point]map[Point]Move
    toToFromToMoveMap map[Point]map[Point]Move
}

func (b *FastBoard) disablePieces(color string, disable bool) error {
    for _, pieceLocation := range b.pieceLocationsMap[color] {
        piece, ok := b.getPiece(pieceLocation)
        if !ok {
            continue
        }
        piece.setDisabled(disable)
    }
    return nil
}

func (b *FastBoard) getPiece(location Point) (Piece, bool) {
    if b.pointOutOfBounds(location) {
        return nil, false
    }

	return b.pieces[location.y][location.x], true
}

func (b *FastBoard) setPiece(location Point, p Piece) bool {
    if b.pointOutOfBounds(location) {
        return false
    }

    oldPiece, ok := b.getPiece(location)
    if ok && oldPiece != nil {
        if pieceLocations, ok := b.pieceLocationsMap[oldPiece.getColor()]; ok {
            removeIndex := -1
            for i, pieceLocation := range pieceLocations {
                if pieceLocation.equals(location) {
                    removeIndex = i
                    break
                }
            }
            if removeIndex != -1 {
                pieceLocations[removeIndex] = pieceLocations[len(pieceLocations)-1]
                pieceLocations[len(pieceLocations)-1] = Point{}
                pieceLocations = pieceLocations[:len(pieceLocations)-1]
                b.pieceLocationsMap[oldPiece.getColor()] = pieceLocations
            }
        }
    }

    if p != nil {
        if pieceLocations, ok := b.pieceLocationsMap[p.getColor()]; ok {
            pieceLocations = append(pieceLocations, location)
            b.pieceLocationsMap[p.getColor()] = pieceLocations
        } else {
            b.pieceLocationsMap[p.getColor()] = []Point{location}
        }

        if _, ok := p.(*King); ok {
            b.kingLocationMap[p.getColor()] = location
        }
    }

	b.pieces[location.y][location.x] = p
    return true
}

func (b *FastBoard) disableLocation(location Point) error {
    b.disabledLocations[location] = true
    return nil
}

func (b *FastBoard) getVulnerables(color string) ([]Point, error) {
    vulnerables, okVulnerables := b.vulnerablesMap[color]
    kingLocation, okKingLocation := b.kingLocationMap[color]

    if okVulnerables && okKingLocation {
        return append(vulnerables, kingLocation), nil
    } else if okVulnerables {
        return vulnerables, nil
    } else if okKingLocation {
        return []Point{kingLocation}, nil
    } else {
        return []Point{}, nil
    }
}

func (b *FastBoard) setVulnerables(color string, locations []Point) error {
    b.vulnerablesMap[color] = locations
    return nil
}

func (b *FastBoard) getEnPassant(color string) (*EnPassant, error) {
	en, ok := b.enPassantMap[color]
	if !ok {
		return nil, nil
	}

	return en, nil
}

func (b *FastBoard) setEnPassant(color string, enPassant *EnPassant) error {
	b.enPassantMap[color] = enPassant
    return nil
}

func (b *FastBoard) possibleEnPassant(color string, target Point) ([]*EnPassant, error) {
    enPassants := []*EnPassant{}

	for k, v := range b.enPassantMap {
        if v == nil {
            continue
        }
		if k != color && target.equals(v.target) {
            enPassants = append(enPassants, v)
		}
	}

    return enPassants, nil
}

func (b *FastBoard) clearEnPassant(color string) error {
    delete(b.enPassantMap, color)
    return nil
}

func (b *FastBoard) ValidMoves(fromLocation Point) ([]Move, error) {
    panic("not implemented")
}

func (b *FastBoard) Move(fromLocation Point, toLocation Point, promotion string) (Move, error) {
    panic("not implemented")
}

func (b *FastBoard) AvailableMoves(color string) ([]Move, error) {
    moves := []Move{}

    for _, pieceLocation := range b.pieceLocationsMap[color] {
        for _, move := range b.fromToToToMoveMap[pieceLocation] { 
            if _, ok := move.(*AllyDefenseMove); ok {
                continue
            }

            moves = append(moves, move)
        }
    }

    return moves, nil
}

func (b *FastBoard) CalculateMoves() error {
    b.fromToToToMoveMap = map[Point]map[Point]Move{}
    b.toToFromToMoveMap = map[Point]map[Point]Move{}

    for _, pieceLocations := range b.pieceLocationsMap {
        for _, fromLocation := range pieceLocations {
            piece, ok := b.getPiece(fromLocation)
            if piece == nil || !ok {
                continue
            }

            moves := piece.moves(b, fromLocation)
            for _, move := range moves {
                action := move.getAction()

                if _, ok := b.fromToToToMoveMap[action.fromLocation]; !ok {
                    b.fromToToToMoveMap[action.fromLocation] = map[Point]Move{}
                }
                if _, ok := b.toToFromToMoveMap[action.toLocation]; !ok {
                    b.toToFromToMoveMap[action.toLocation] = map[Point]Move{}
                }

                b.fromToToToMoveMap[action.fromLocation][action.toLocation] = move
                b.toToFromToMoveMap[action.toLocation][action.fromLocation] = move
            }
        }
    }

    return nil
}

func (b *FastBoard) LegalMoves(color string) ([]Move, error) {
    panic("not implemented")
}

func (b *FastBoard) CalculateMovesPartial(move Move) error {
    panic("not implemented")
}

func (b *FastBoard) Print() string {
	var builder strings.Builder
	var cellWidth int = 12

	for y, row := range b.pieces {
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*b.size.x-1)))
		for x := range row {
			builder.WriteString(fmt.Sprintf("|%s%2dx ", strings.Repeat(" ", cellWidth-4), x))
		}
		builder.WriteString("|\n")
		for x, piece := range row {
            if b.pointOutOfBounds(Point{x, y}) {
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
	builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*b.size.x-1)))

	return builder.String()
}

func (b *FastBoard) State() *BoardData {
    panic("not implemented")
}

func (b *FastBoard) Check(color string) bool {
    // TODO this shouldn't return an error
    vulnerableLocations, err := b.getVulnerables(color)
    if err != nil {
        return false
    }

    for _, vulnerableLocation := range vulnerableLocations {
        if fromToMoveMap, ok := b.toToFromToMoveMap[vulnerableLocation]; ok {
            for fromLocation := range fromToMoveMap {
                piece, ok := b.getPiece(fromLocation)
                if piece == nil || !ok {
                    continue
                }

                if piece.getColor() != color {
                    return true
                }
            }
        }
    }

    return false
}

func (b *FastBoard) Checkmate(color string) bool {
    panic("not implemented")
}

func (b *FastBoard) Stalemate(color string) bool {
    panic("not implemented")
}

func (b *FastBoard) getPieceLocations() map[string][]Point {
    return b.pieceLocationsMap
}

func (b *FastBoard) Copy() (Board, error) {
    panic("not implemented")
}

func (b *FastBoard) UniqueString() string {
    panic("not implemented")
}

func (b *FastBoard) pointOutOfBounds(p Point) bool {
    _, ok := b.disabledLocations[p]
    return ok || p.y < 0 || p.y >= b.size.y || p.x < 0 || p.x >= b.size.x
}

func (b *FastBoard) Mate(color string) (bool, bool, error) {
    panic("not implemented")
}

