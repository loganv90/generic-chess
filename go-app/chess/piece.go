package chess

import "fmt"

type CastleDirection struct {
	Point
	kingOffset *Point
	rookOffset *Point
}

type Allegiant struct {
	color string
}

func (a *Allegiant) getColor() string {
	return a.color
}

type Piece interface {
	getColor() string
	movedCopy() Piece
	moves(Board, *Point) []Move
	print() string
}

func addDirection(
    fromLocation *Point,
	b Board,
	moves *[]Move,
    direction *Point,
	color string,
) {
    currentLocation := fromLocation.add(direction)

    for {
        p, err := b.getPiece(currentLocation)
        if err != nil {
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
			break
		}

        currentLocation = currentLocation.add(direction)
	}
}

func addSimple(
    fromLocation *Point,
	b Board,
	moves *[]Move,
    direction *Point,
	color string,
) {
    toLocation := fromLocation.add(direction)

	p, err := b.getPiece(toLocation)
	if err != nil {
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
	}
}

func newPawn(color string, moved bool, xDir int, yDir int) (*Pawn, error) {
    var forward1 *Point
    var forward2 *Point
    var captures []*Point

	if xDir == 1 || xDir == -1 {
		forward1 = &Point{xDir, 0}
		forward2 = &Point{xDir * 2, 0}
		captures = []*Point{
			{xDir, 1},
			{xDir, -1},
		}
	} else if yDir == 1 || yDir == -1 {
		forward1 = &Point{0, yDir}
		forward2 = &Point{0, yDir * 2}
		captures = []*Point{
			{1, yDir},
			{-1, yDir},
		}
	} else {
		return nil, fmt.Errorf("invalid direction")
	}

	return &Pawn{
		Allegiant{color},
		moved,
		forward1,
		forward2,
		captures,
	}, nil
}

type Pawn struct {
	Allegiant
	moved    bool
	forward1 *Point
	forward2 *Point
	captures []*Point
}

func (a *Pawn) print() string {
	return "P"
}

func (a *Pawn) movedCopy() Piece {
	return &Pawn{
		Allegiant{a.color},
		true,
		a.forward1,
		a.forward2,
		a.captures,
	}
}

func (a *Pawn) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	a.addForward(b, fromLocation, moves)
	a.addCaptures(b, fromLocation, moves)
	return *moves
}

func (a *Pawn) addForward(b Board, fromLocation *Point, moves *[]Move) {
    to1Location := fromLocation.add(a.forward1)

	if piece, err := b.getPiece(to1Location); err != nil || piece != nil {
		return
    } else if b.pointOnPromotionSquare(to1Location) {
		promotionMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, to1Location)
		if err == nil {
			*moves = append(*moves, promotionMove)
		}
	} else {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, to1Location)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	}

	if a.moved {
		return
	}

    to2Location := fromLocation.add(a.forward2)

	if piece, err := b.getPiece(to2Location); err != nil || piece != nil {
		return
    } else if b.pointOnPromotionSquare(to2Location) {
		promotionMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, to2Location)
		if err == nil {
			*moves = append(*moves, promotionMove)
		}
	} else {
		simpleMove, err := moveFactoryInstance.newRevealEnPassantMove(b, fromLocation, to2Location, to1Location)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	}
}

func (a *Pawn) addCaptures(b Board, fromLocation *Point, moves *[]Move) {
	for _, capture := range a.captures {
        toLocation := fromLocation.add(capture)

		if piece, err := b.getPiece(toLocation); err != nil {
			continue
		} else if len(b.possibleEnPassants(a.color, toLocation)) > 0 {
			captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(b, fromLocation, toLocation)
			if err == nil {
				*moves = append(*moves, captureEnPassantMove)
			}
		} else if piece != nil && piece.getColor() != a.color {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
		}
	}
}

var knightSimples = []*Point{
	{1, 2},
	{-1, 2},
	{2, 1},
	{-2, 1},
	{1, -2},
	{-1, -2},
	{2, -1},
	{-2, -1},
}

func newKnight(color string) (*Knight, error) {
	return &Knight{
		Allegiant{color},
	}, nil
}

type Knight struct {
	Allegiant
}

func (n *Knight) print() string {
	return "N"
}

func (n *Knight) movedCopy() Piece {
	return &Knight{
		Allegiant{n.color},
	}
}

func (n *Knight) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	n.addSimples(b, fromLocation, moves)
	return *moves
}

func (n *Knight) addSimples(b Board, fromLocation *Point, moves *[]Move) {
	for _, simple := range knightSimples {
		addSimple(fromLocation, b, moves, simple, n.color)
	}
}

var bishopDirections = []*Point{
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newBishop(color string) (*Bishop, error) {
	return &Bishop{
		Allegiant{color},
	}, nil
}

type Bishop struct {
	Allegiant
}

func (s *Bishop) print() string {
	return "B"
}

func (s *Bishop) movedCopy() Piece {
	return &Bishop{
		Allegiant{s.color},
	}
}

func (s *Bishop) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	s.addDirections(b, fromLocation, moves)
	return *moves
}

func (s *Bishop) addDirections(b Board, fromLocation *Point, moves *[]Move) {
	for _, direction := range bishopDirections {
		addDirection(fromLocation, b, moves, direction, s.color)
	}
}

var rookDirections = []*Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func newRook(color string, moved bool) (*Rook, error) {
	return &Rook{
		Allegiant{color},
		moved,
	}, nil
}

type Rook struct {
	Allegiant
	moved bool
}

func (r *Rook) print() string {
	return "R"
}

func (r *Rook) movedCopy() Piece {
	return &Rook{
		Allegiant{r.color},
		true,
	}
}

func (r *Rook) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	r.addDirections(b, fromLocation, moves)
	return *moves
}

func (r *Rook) addDirections(b Board, fromLocation *Point, moves *[]Move) {
	for _, direction := range rookDirections {
		addDirection(fromLocation, b, moves, direction, r.color)
	}
}

var queenDirections = []*Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newQueen(color string) (*Queen, error) {
	return &Queen{
		Allegiant{color},
	}, nil
}

type Queen struct {
	Allegiant
}

func (q *Queen) print() string {
	return "Q"
}

func (q *Queen) movedCopy() Piece {
	return &Queen{
		Allegiant{q.color},
	}
}

func (q *Queen) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	q.addDirections(b, fromLocation, moves)
	return *moves
}

func (q *Queen) addDirections(b Board, fromLocation *Point, moves *[]Move) {
	for _, direction := range queenDirections {
		addDirection(fromLocation, b, moves, direction, q.color)
	}
}

var kingSimples = []*Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newKing(color string, moved bool, xDir int, yDir int) (*King, error) {
	var castles []*CastleDirection

	if xDir == 1 || xDir == -1 {
		castles = []*CastleDirection{
			{Point{0, 1}, &Point{0, -2}, &Point{0, -3}},
			{Point{0, -1}, &Point{0, 1}, &Point{0, 2}},
		}
	} else if yDir == 1 || yDir == -1 {
		castles = []*CastleDirection{
			{Point{1, 0}, &Point{-1, 0}, &Point{-2, 0}},
			{Point{-1, 0}, &Point{2, 0}, &Point{3, 0}},
		}
	} else {
		return nil, fmt.Errorf("invalid direction")
	}

	return &King{
		Allegiant{color},
		moved,
		castles,
	}, nil
}

type King struct {
	Allegiant
	moved   bool
	castles []*CastleDirection
}

func (k *King) print() string {
	return "K"
}

func (k *King) movedCopy() Piece {
	return &King{
		Allegiant{k.color},
		true,
		k.castles,
	}
}

func (k *King) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	k.addSimples(b, fromLocation, moves)
	k.addCastles(b, fromLocation, moves)
	return *moves
}

func (k *King) addSimples(b Board, fromLocation *Point, moves *[]Move) {
	for _, simple := range kingSimples {
		addSimple(fromLocation, b, moves, simple, k.color)
	}
}

func (k *King) addCastles(b Board, fromLocation *Point, moves *[]Move) {
	if k.moved {
		return
	}

	for _, castle := range k.castles {
        currentLocation := fromLocation.add(&castle.Point)

		rookFound := false
		for {
			piece, err := b.getPiece(currentLocation)
			if err != nil {
				break
			}

			if piece == nil {
                currentLocation = currentLocation.add(&castle.Point)
				continue
			}

			rook, ok := piece.(*Rook)
			if !ok || rook.moved || rook.color != k.color {
				break
			}

			rookFound = true
			break
		}
		if !rookFound {
			continue
		}

		xToKing := fromLocation.x
		xToRook := currentLocation.x
		yToKing := fromLocation.y
		yToRook := currentLocation.y
		if castle.kingOffset.x > 0 {
			xToKing = castle.kingOffset.x
		} else if castle.kingOffset.x < 0 {
			xToKing = b.xLen() - 1 + castle.kingOffset.x
		}
		if castle.kingOffset.y > 0 {
			yToKing = castle.kingOffset.y
		} else if castle.kingOffset.y < 0 {
			yToKing = b.yLen() - 1 + castle.kingOffset.y
		}
		if castle.rookOffset.x > 0 {
			xToRook = castle.rookOffset.x
		} else if castle.rookOffset.x < 0 {
			xToRook = b.xLen() - 1 + castle.rookOffset.x
		}
		if castle.rookOffset.y > 0 {
			yToRook = castle.rookOffset.y
		} else if castle.rookOffset.y < 0 {
			yToRook = b.yLen() - 1 + castle.rookOffset.y
		}
        if xToKing < 0 || xToKing >= b.xLen() || yToKing < 0 || yToKing >= b.yLen() {
            continue
        }
        if xToRook < 0 || xToRook >= b.xLen() || yToRook < 0 || yToRook >= b.yLen() {
            continue
        }

        xCheckedMin := min(fromLocation.x, currentLocation.x)
        xCheckedMax := max(fromLocation.x, currentLocation.x)
        yCheckedMin := min(fromLocation.y, currentLocation.y)
        yCheckedMax := max(fromLocation.y, currentLocation.y)

        xToMin := min(xToKing, xToRook)
        xToMax := max(xToKing, xToRook)
        yToMin := min(yToKing, yToRook)
        yToMax := max(yToKing, yToRook)

        clear := true
        for x := xCheckedMin; x > xToMin && clear; x-- {
            if piece, err := b.getPiece(&Point{x, fromLocation.y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        for y := yCheckedMin; y > yToMin && clear; y-- {
            if piece, err := b.getPiece(&Point{fromLocation.x, y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        for x := xCheckedMax; x < xToMax && clear; x++ {
            if piece, err := b.getPiece(&Point{x, fromLocation.y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        for y := yCheckedMax; y < yToMax && clear; y++ {
            if piece, err := b.getPiece(&Point{fromLocation.y, y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        if !clear {
            continue
        }

		castleMove, err := moveFactoryInstance.newCastleMove(b, fromLocation, currentLocation, &Point{xToKing, yToKing}, &Point{xToRook, yToRook})
		if err == nil {
			*moves = append(*moves, castleMove)
		}
	}
}

