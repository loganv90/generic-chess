package chess

import (
    "fmt"
)

const (
    EMPTY = 0
    PAWN_R = 1
    PAWN_L = 2
    PAWN_D = 3
    PAWN_U = 4
    PAWN_R_M = 5
    PAWN_L_M = 6
    PAWN_D_M = 7
    PAWN_U_M = 8
    KNIGHT = 9
    BISHOP = 10
    ROOK = 11
    ROOK_M = 12
    QUEEN = 13
    KING_R = 14
    KING_L = 15
    KING_D = 16
    KING_U = 17
    KING_R_M = 18
    KING_L_M = 19
    KING_D_M = 20
    KING_U_M = 21
)

var pieceFactoryInstance = newPieceFactory()

func newPieceFactory() *PieceFactory {
    var pieceInstances = make([][]Piece, 4)
    for color := range pieceInstances {
        pieceInstances[color] = generateRowForPieceReference(color)
    }

    return &PieceFactory{
        pieceInstances: pieceInstances,
    }
}

func generateRowForPieceReference(color int) []Piece {
    return []Piece{
        nil,
        newPawn(color, false, 1, 0),
        newPawn(color, false, -1, 0),
        newPawn(color, false, 0, 1),
        newPawn(color, false, 0, -1),
        newPawn(color, true, 1, 0),
        newPawn(color, true, -1, 0),
        newPawn(color, true, 0, 1),
        newPawn(color, true, 0, -1),
        newKnight(color),
        newBishop(color),
        newRook(color, false),
        newRook(color, true),
        newQueen(color),
        newKing(color, false, 1, 0),
        newKing(color, false, -1, 0),
        newKing(color, false, 0, 1),
        newKing(color, false, 0, -1),
        newKing(color, true, 1, 0),
        newKing(color, true, -1, 0),
        newKing(color, true, 0, 1),
        newKing(color, true, 0, -1),
    }
}

type PieceFactory struct {
    pieceInstances [][]Piece
}

func (r *PieceFactory) get(color int, pieceType int) Piece {
    return r.pieceInstances[color][pieceType]
}

type CastleDirection struct {
	Point
	kingOffset Point
	rookOffset Point
}

type Allegiant struct {
	color int
}

func (a *Allegiant) getColor() int {
	return a.color
}

type Piece interface {
	getColor() int
    getValue() int
	copy() (Piece, error) // returns the moved version of the piece
    getMoved() bool
	moves(Board, Point) []Move
	print() string
}

func addDirection(
    fromLocation Point,
	b Board,
	moves *[]Move,
    direction Point,
	color int,
) {
    currentLocation := fromLocation.add(direction)

    for {
        p, ok := b.getPiece(currentLocation)
        if !ok {
            break
        }

		if p == nil {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, currentLocation)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
		} else if p.getColor() != color {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, currentLocation)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
			break
		} else {
            allyDefenseMove, err := moveFactoryInstance.newAllyDefenseMove(b, fromLocation, currentLocation)
            if err == nil {
                *moves = append(*moves, allyDefenseMove)
            }
			break
		}

        currentLocation = currentLocation.add(direction)
	}
}

func addSimple(
    fromLocation Point,
	b Board,
	moves *[]Move,
    direction Point,
	color int,
) {
    toLocation := fromLocation.add(direction)

	p, ok := b.getPiece(toLocation)
	if !ok {
		return
	}

	if p == nil {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	} else if p.getColor() != color {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	} else {
        allyDefenseMove, err := moveFactoryInstance.newAllyDefenseMove(b, fromLocation, toLocation)
        if err == nil {
            *moves = append(*moves, allyDefenseMove)
        }
    }
}

func newPawn(color int, moved bool, xDir int, yDir int) *Pawn {
    var forward1 Point
    var forward2 Point
    var captures []Point

	if xDir == 1 || xDir == -1 {
		forward1 = Point{xDir, 0}
		forward2 = Point{xDir * 2, 0}
		captures = []Point{
			{xDir, 1},
			{xDir, -1},
		}
	} else if yDir == 1 || yDir == -1 {
		forward1 = Point{0, yDir}
		forward2 = Point{0, yDir * 2}
		captures = []Point{
			{1, yDir},
			{-1, yDir},
		}
	} else {
        panic("invalid direction")
	}

	return &Pawn{
		Allegiant{color},
		moved,
		forward1,
		forward2,
		captures,
        Point{xDir, yDir},
	}
}

type Pawn struct {
	Allegiant
	moved bool
	forward1 Point
	forward2 Point
	captures []Point
    direction Point
}

func (a *Pawn) print() string {
	return "P"
}

func (a *Pawn) copy() (Piece, error) {
    if a.direction.x == 1 {
        return pieceFactoryInstance.get(a.color, PAWN_R), nil
    } else if a.direction.x == -1 {
        return pieceFactoryInstance.get(a.color, PAWN_L), nil
    } else if a.direction.y == 1 {
        return pieceFactoryInstance.get(a.color, PAWN_D), nil
    } else if a.direction.y == -1 {
        return pieceFactoryInstance.get(a.color, PAWN_U), nil
    } 

    return nil, fmt.Errorf("invalid direction")
}

func (a *Pawn) getMoved() bool {
    return a.moved
}

func (a *Pawn) getValue() int {
    return 1000
}

func (a *Pawn) moves(b Board, fromLocation Point) []Move {
	moves := &[]Move{}
	a.addForward(b, fromLocation, moves)
	a.addCaptures(b, fromLocation, moves)
	return *moves
}

func (a *Pawn) nextLocationInvalid(b Board, toLocation Point) bool {
    nextLocation := toLocation.add(a.direction)
    _, ok := b.getPiece(nextLocation)
    return !ok
}

func (a *Pawn) addForward(b Board, fromLocation Point, moves *[]Move) {
    to1Location := fromLocation.add(a.forward1)
	if piece, ok := b.getPiece(to1Location); !ok || piece != nil { // if the square is invalid or occupied
		return
	} else {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, to1Location)
        if err != nil {
            return
        } else if a.nextLocationInvalid(b, to1Location) {
            if promotionMove, err := moveFactoryInstance.newPromotionMove(simpleMove); err == nil {
                *moves = append(*moves, promotionMove)
            }
        } else {
			*moves = append(*moves, simpleMove)
		}
	}

	if a.moved {
		return
	}

    to2Location := fromLocation.add(a.forward2)
	if piece, ok := b.getPiece(to2Location); !ok || piece != nil { // if the square is invalid or occupied
		return
	} else {
		revealEnPassantMove, err := moveFactoryInstance.newRevealEnPassantMove(b, fromLocation, to2Location, to1Location)
        if err != nil {
            return
        } else if a.nextLocationInvalid(b, to2Location) {
            if promotionMove, err := moveFactoryInstance.newPromotionMove(revealEnPassantMove); err == nil {
                *moves = append(*moves, promotionMove)
            }
        } else {
			*moves = append(*moves, revealEnPassantMove)
		}
	}
}

func (a *Pawn) addCaptures(b Board, fromLocation Point, moves *[]Move) {
	for _, capture := range a.captures {
        toLocation := fromLocation.add(capture)

		if piece, ok := b.getPiece(toLocation); !ok { // if the square is invalid
			continue
        } else if ens, err := b.possibleEnPassant(a.color, toLocation); err == nil && len(ens) > 0 { // if the square is an en passant target
			captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(b, fromLocation, toLocation)
            if err != nil {
                continue
            } else if a.nextLocationInvalid(b, toLocation) {
                if promotionMove, err := moveFactoryInstance.newPromotionMove(captureEnPassantMove); err == nil {
                    *moves = append(*moves, promotionMove)
                }
            } else {
                *moves = append(*moves, captureEnPassantMove)
            }
		} else if piece != nil && piece.getColor() != a.color { // if the square is occupied by an enemy piece
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
            if err != nil {
                continue
            } else if a.nextLocationInvalid(b, toLocation) {
                if promotionMove, err := moveFactoryInstance.newPromotionMove(simpleMove); err == nil {
                    *moves = append(*moves, promotionMove)
                }
            } else {
                *moves = append(*moves, simpleMove)
            }
		} else if piece != nil {
            allyDefenseMove, err := moveFactoryInstance.newAllyDefenseMove(b, fromLocation, toLocation)
            if err == nil {
                *moves = append(*moves, allyDefenseMove)
            }
        }
	}
}

var knightSimples = []Point{
	{1, 2},
	{-1, 2},
	{2, 1},
	{-2, 1},
	{1, -2},
	{-1, -2},
	{2, -1},
	{-2, -1},
}

func newKnight(color int) *Knight {
	return &Knight{
		Allegiant{color},
	}
}

type Knight struct {
	Allegiant
}

func (n *Knight) print() string {
	return "N"
}

func (n *Knight) copy() (Piece, error) {
    return pieceFactoryInstance.get(n.color, KNIGHT), nil
}

func (n *Knight) getMoved() bool {
    return true
}

func (n *Knight) getValue() int {
    return 3000
}

func (n *Knight) moves(b Board, fromLocation Point) []Move {
	moves := &[]Move{}
	n.addSimples(b, fromLocation, moves)
	return *moves
}

func (n *Knight) addSimples(b Board, fromLocation Point, moves *[]Move) {
	for _, simple := range knightSimples {
		addSimple(fromLocation, b, moves, simple, n.color)
	}
}

var bishopDirections = []Point{
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newBishop(color int) *Bishop {
	return &Bishop{
		Allegiant{color},
	}
}

type Bishop struct {
	Allegiant
}

func (s *Bishop) print() string {
	return "B"
}

func (s *Bishop) copy() (Piece, error) {
    return pieceFactoryInstance.get(s.color, BISHOP), nil
}

func (s *Bishop) getMoved() bool {
    return true
}

func (s *Bishop) getValue() int {
    return 3000
}

func (s *Bishop) moves(b Board, fromLocation Point) []Move {
	moves := &[]Move{}
	s.addDirections(b, fromLocation, moves)
	return *moves
}

func (s *Bishop) addDirections(b Board, fromLocation Point, moves *[]Move) {
	for _, direction := range bishopDirections {
		addDirection(fromLocation, b, moves, direction, s.color)
	}
}

var rookDirections = []Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func newRook(color int, moved bool) *Rook {
	return &Rook{
		Allegiant{color},
		moved,
	}
}

type Rook struct {
	Allegiant
	moved bool
}

func (r *Rook) print() string {
	return "R"
}

func (r *Rook) copy() (Piece, error) {
    return pieceFactoryInstance.get(r.color, ROOK_M), nil
}

func (r *Rook) getMoved() bool {
    return r.moved
}

func (r *Rook) getValue() int {
    return 5000
}

func (r *Rook) moves(b Board, fromLocation Point) []Move {
	moves := &[]Move{}
	r.addDirections(b, fromLocation, moves)
	return *moves
}

func (r *Rook) addDirections(b Board, fromLocation Point, moves *[]Move) {
	for _, direction := range rookDirections {
		addDirection(fromLocation, b, moves, direction, r.color)
	}
}

var queenDirections = []Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newQueen(color int) *Queen {
	return &Queen{
		Allegiant{color},
	}
}

type Queen struct {
	Allegiant
}

func (q *Queen) print() string {
	return "Q"
}

func (q *Queen) copy() (Piece, error) {
    return pieceFactoryInstance.get(q.color, QUEEN), nil
}

func (q *Queen) getMoved() bool {
    return true
}

func (q *Queen) getValue() int {
    return 9000
}

func (q *Queen) moves(b Board, fromLocation Point) []Move {
	moves := &[]Move{}
	q.addDirections(b, fromLocation, moves)
	return *moves
}

func (q *Queen) addDirections(b Board, fromLocation Point, moves *[]Move) {
	for _, direction := range queenDirections {
		addDirection(fromLocation, b, moves, direction, q.color)
	}
}

var kingSimples = []Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newKing(color int, moved bool, xDir int, yDir int) *King {
	var castles []*CastleDirection

	if xDir == 1 || xDir == -1 {
		castles = []*CastleDirection{
			{Point{0, 1}, Point{0, -1}, Point{0, -2}},
			{Point{0, -1}, Point{0, 2}, Point{0, 3}},
		}
	} else if yDir == 1 || yDir == -1 {
		castles = []*CastleDirection{
			{Point{1, 0}, Point{-1, 0}, Point{-2, 0}},
			{Point{-1, 0}, Point{2, 0}, Point{3, 0}},
		}
	} else {
        panic("invalid direction")
	}

	return &King{
		Allegiant{color},
		moved,
		castles,
        Point{xDir, yDir},
	}
}

type King struct {
	Allegiant
	moved   bool
	castles []*CastleDirection
    direction Point
}

func (k *King) print() string {
	return "K"
}

func (k *King) copy() (Piece, error) {
    if k.direction.x == 1 {
        return pieceFactoryInstance.get(k.color, KING_R), nil
    } else if k.direction.x == -1 {
        return pieceFactoryInstance.get(k.color, KING_L), nil
    } else if k.direction.y == 1 {
        return pieceFactoryInstance.get(k.color, KING_D), nil
    } else if k.direction.y == -1 {
        return pieceFactoryInstance.get(k.color, KING_U), nil
    }

    return nil, fmt.Errorf("invalid direction")
}

func (k *King) getMoved() bool {
    return k.moved
}

func (k *King) getValue() int {
    return 0
}

func (k *King) moves(b Board, fromLocation Point) []Move {
	moves := &[]Move{}
	k.addSimples(b, fromLocation, moves)
	k.addCastles(b, fromLocation, moves)
	return *moves
}

func (k *King) addSimples(b Board, fromLocation Point, moves *[]Move) {
	for _, simple := range kingSimples {
		addSimple(fromLocation, b, moves, simple, k.color)
	}
}

func (k *King) findRookForCastle(b Board, fromLocation Point, direction Point) (Point, bool) {
    currentLocation := fromLocation.add(direction)

    for {
        piece, ok := b.getPiece(currentLocation)
        if !ok {
            return Point{}, false
        }

        if piece == nil {
            currentLocation = currentLocation.add(direction)
            continue
        }

        if rook, ok := piece.(*Rook); !ok || rook.moved || rook.color != k.color {
            return Point{}, false
        }

        return currentLocation, true
    }
}

func (k *King) findEdgeForCastle(b Board, fromLocation Point, direction Point) (Point, error) {
    previousLocation := fromLocation
    currentLocation := previousLocation.add(direction)

    for {
        _, ok := b.getPiece(currentLocation)
        if !ok {
            return previousLocation, nil
        }

        previousLocation = currentLocation
        currentLocation = previousLocation.add(direction)
    }
}

func (k *King) addCastles(b Board, fromLocation Point, moves *[]Move) {
	if k.moved {
		return
	}

	for _, castle := range k.castles {
        fromRookLocation, ok := k.findRookForCastle(b, fromLocation, castle.Point)
        if !ok {
            continue
        }

        edgeLocation, err := k.findEdgeForCastle(b, fromRookLocation, castle.Point)
        if err != nil {
            continue
        }

        toLocation := edgeLocation.add(castle.kingOffset)
        toRookLocation := edgeLocation.add(castle.rookOffset)

        xCheckedMin := min(fromLocation.x, fromRookLocation.x)
        xCheckedMax := max(fromLocation.x, fromRookLocation.x)
        yCheckedMin := min(fromLocation.y, fromRookLocation.y)
        yCheckedMax := max(fromLocation.y, fromRookLocation.y)

        xToMin := min(toLocation.x, toRookLocation.x)
        xToMax := max(toLocation.x, toRookLocation.x)
        yToMin := min(toLocation.y, toRookLocation.y)
        yToMax := max(toLocation.y, toRookLocation.y)

        clear := true
        for x := xCheckedMin - 1; x >= xToMin && clear; x-- {
            if piece, ok := b.getPiece(Point{x, fromLocation.y}); !ok || piece != nil {
                clear = false
                break
            }
        }
        for y := yCheckedMin - 1; y >= yToMin && clear; y-- {
            if piece, ok := b.getPiece(Point{fromLocation.x, y}); !ok || piece != nil {
                clear = false
                break
            }
        }
        for x := xCheckedMax + 1; x <= xToMax && clear; x++ {
            if piece, ok := b.getPiece(Point{x, fromLocation.y}); !ok || piece != nil {
                clear = false
                break
            }
        }
        for y := yCheckedMax + 1; y <= yToMax && clear; y++ {
            if piece, ok := b.getPiece(Point{fromLocation.y, y}); !ok || piece != nil {
                clear = false
                break
            }
        }
        if !clear {
            continue
        }

        vulnerableLocations := []Point{}
        if toLocation.x > fromLocation.x {
            for x := fromLocation.x + 1; x < toLocation.x; x++ {
                location := Point{x, fromLocation.y}
                vulnerableLocations = append(vulnerableLocations, location)
            }
        } else if toLocation.x < fromLocation.x {
            for x := fromLocation.x - 1; x > toLocation.x; x-- {
                location := Point{x, fromLocation.y}
                vulnerableLocations = append(vulnerableLocations, location)
            }
        } else if toLocation.y > fromLocation.y {
            for y := fromLocation.y + 1; y < toLocation.y; y++ {
                location := Point{fromLocation.x, y}
                vulnerableLocations = append(vulnerableLocations, location)
            }
        } else if toLocation.y < fromLocation.y {
            for y := fromLocation.y - 1; y > toLocation.y; y-- {
                location := Point{fromLocation.x, y}
                vulnerableLocations = append(vulnerableLocations, location)
            }
        }

		castleMove, err := moveFactoryInstance.newCastleMove(
            b,
            fromLocation,
            fromRookLocation,
            toLocation,
            toRookLocation,
            vulnerableLocations,
        )
		if err == nil {
			*moves = append(*moves, castleMove)
		}
	}
}
