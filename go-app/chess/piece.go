package chess

import "fmt"

type direction struct {
	dx int
	dy int
}

type allegiant struct {
	color string
}

func (a *allegiant) getColor() string {
	return a.color
}

type piece interface {
	getColor() string
	movedCopy() piece
	moves(board, int, int) []move
	print() string
}

func movedOutOfBounds(
	xTo int,
	yTo int,
	b board,
) bool {
	return xTo < 0 ||
		xTo >= b.xLen() ||
		yTo < 0 ||
		yTo >= b.yLen()
}

func movedToPromotionSquare(
	xTo int,
	yTo int,
	xFrom int,
	yFrom int,
	b board,
) bool {
	return (xTo == 0 && xFrom != 0) ||
		(xTo == b.xLen()-1 && xFrom != b.xLen()-1) ||
		(yTo == 0 && yFrom != 0) ||
		(yTo == b.yLen()-1 && yFrom != b.yLen()-1)
}

func addDirection(
	xFrom int,
	yFrom int,
	b board,
	moves *[]move,
	dx int,
	dy int,
	color string,
) {
	xCurrent := xFrom + dx
	yCurrent := yFrom + dy

	for !movedOutOfBounds(xCurrent, yCurrent, b) {
		p, err := b.getPiece(xCurrent, yCurrent)
		if err != nil {
			break
		}

		if p == nil {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
		} else if p.getColor() != color {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
			break
		} else {
			break
		}

		xCurrent += dx
		yCurrent += dy
	}
}

func addSimple(
	xFrom int,
	yFrom int,
	b board,
	moves *[]move,
	dx int,
	dy int,
	color string,
) {
	xTo := xFrom + dx
	yTo := yFrom + dy

	if movedOutOfBounds(xTo, yTo, b) {
		return
	}

	p, err := b.getPiece(xTo, yTo)
	if err != nil {
		return
	}

	if p == nil {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	} else if p.getColor() != color {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	}
}

func newPawn(color string, moved bool, xDir int, yDir int) (*pawn, error) {
	var forward1 *direction
	var forward2 *direction
	var captures []*direction

	if xDir == 1 || xDir == -1 {
		forward1 = &direction{xDir, 0}
		forward2 = &direction{xDir * 2, 0}
		captures = []*direction{
			{xDir, 1},
			{xDir, -1},
		}
	} else if yDir == 1 || yDir == -1 {
		forward1 = &direction{0, yDir}
		forward2 = &direction{0, yDir * 2}
		captures = []*direction{
			{1, yDir},
			{-1, yDir},
		}
	} else {
		return nil, fmt.Errorf("invalid direction")
	}

	return &pawn{
		allegiant{color},
		moved,
		forward1,
		forward2,
		captures,
	}, nil
}

type pawn struct {
	allegiant
	moved    bool
	forward1 *direction
	forward2 *direction
	captures []*direction
}

func (a *pawn) print() string {
	return "P"
}

func (a *pawn) movedCopy() piece {
	return &pawn{
		allegiant{a.color},
		true,
		a.forward1,
		a.forward2,
		a.captures,
	}
}

func (a *pawn) moves(b board, xFrom int, yFrom int) []move {
	moves := &[]move{}
	a.addForward(b, xFrom, yFrom, moves)
	a.addCaptures(b, xFrom, yFrom, moves)
	return *moves
}

func (a *pawn) addForward(b board, xFrom int, yFrom int, moves *[]move) {
	xTo1 := xFrom + a.forward1.dx
	yTo1 := yFrom + a.forward1.dy

	if piece, err := b.getPiece(xTo1, yTo1); err != nil || piece != nil {
		return
	} else if movedToPromotionSquare(xTo1, yTo1, xFrom, yFrom, b) {
		promotionMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xTo1, yTo1)
		if err == nil {
			*moves = append(*moves, promotionMove)
		}
	} else {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xTo1, yTo1)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	}

	if a.moved {
		return
	}

	xTo2 := xFrom + a.forward2.dx
	yTo2 := yFrom + a.forward2.dy

	if piece, err := b.getPiece(xTo2, yTo2); err != nil || piece != nil {
		return
	} else if movedToPromotionSquare(xTo2, yTo2, xFrom, yFrom, b) {
		promotionMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xTo2, yTo2)
		if err == nil {
			*moves = append(*moves, promotionMove)
		}
	} else {
		simpleMove, err := moveFactoryInstance.newRevealEnPassantMove(b, xFrom, yFrom, xTo2, yTo2, xTo1, yTo1)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	}
}

func (a *pawn) addCaptures(b board, xFrom int, yFrom int, moves *[]move) {
	for _, capture := range a.captures {
		xTo := xFrom + capture.dx
		yTo := yFrom + capture.dy

		if piece, err := b.getPiece(xTo, yTo); err != nil {
			continue
		} else if len(b.possibleEnPassants(a.color, xTo, yTo)) > 0 {
			captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(b, xFrom, yFrom, xTo, yTo)
			if err == nil {
				*moves = append(*moves, captureEnPassantMove)
			}
		} else if piece != nil && piece.getColor() != a.color {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
		}
	}
}

var knightSimples = []*direction{
	{1, 2},
	{-1, 2},
	{2, 1},
	{-2, 1},
	{1, -2},
	{-1, -2},
	{2, -1},
	{-2, -1},
}

func newKnight(color string) (*knight, error) {
	return &knight{
		allegiant{color},
	}, nil
}

type knight struct {
	allegiant
}

func (n *knight) print() string {
	return "N"
}

func (n *knight) movedCopy() piece {
	return &knight{
		allegiant{n.color},
	}
}

func (n *knight) moves(b board, xFrom int, yFrom int) []move {
	moves := &[]move{}
	n.addSimples(b, xFrom, yFrom, moves)
	return *moves
}

func (n *knight) addSimples(b board, xFrom int, yFrom int, moves *[]move) {
	for _, simple := range knightSimples {
		addSimple(xFrom, yFrom, b, moves, simple.dx, simple.dy, n.color)
	}
}

var bishopDirections = []*direction{
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newBishop(color string) (*bishop, error) {
	return &bishop{
		allegiant{color},
	}, nil
}

type bishop struct {
	allegiant
}

func (s *bishop) print() string {
	return "B"
}

func (s *bishop) movedCopy() piece {
	return &bishop{
		allegiant{s.color},
	}
}

func (s *bishop) moves(b board, xFrom int, yFrom int) []move {
	moves := &[]move{}
	s.addDirections(b, xFrom, yFrom, moves)
	return *moves
}

func (s *bishop) addDirections(b board, xFrom int, yFrom int, moves *[]move) {
	for _, direction := range bishopDirections {
		addDirection(xFrom, yFrom, b, moves, direction.dx, direction.dy, s.color)
	}
}

var rookDirections = []*direction{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func newRook(color string, moved bool) (*rook, error) {
	return &rook{
		allegiant{color},
		moved,
	}, nil
}

type rook struct {
	allegiant
	moved bool
}

func (r *rook) print() string {
	return "R"
}

func (r *rook) movedCopy() piece {
	return &rook{
		allegiant{r.color},
		true,
	}
}

func (r *rook) moves(b board, x int, y int) []move {
	moves := &[]move{}
	r.addDirections(b, x, y, moves)
	return *moves
}

func (r *rook) addDirections(b board, xFrom int, yFrom int, moves *[]move) {
	for _, direction := range rookDirections {
		addDirection(xFrom, yFrom, b, moves, direction.dx, direction.dy, r.color)
	}
}

var queenDirections = []*direction{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newQueen(color string) (*queen, error) {
	return &queen{
		allegiant{color},
	}, nil
}

type queen struct {
	allegiant
}

func (q *queen) print() string {
	return "Q"
}

func (q *queen) movedCopy() piece {
	return &queen{
		allegiant{q.color},
	}
}

func (q *queen) moves(b board, x int, y int) []move {
	moves := &[]move{}
	q.addDirections(b, x, y, moves)
	return *moves
}

func (q *queen) addDirections(b board, xFrom int, yFrom int, moves *[]move) {
	for _, direction := range queenDirections {
		addDirection(xFrom, yFrom, b, moves, direction.dx, direction.dy, q.color)
	}
}

var kingSimples = []*direction{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newKing(color string, moved bool, xDir int, yDir int) (*king, error) {
	var castles []*direction

	if xDir == 1 || xDir == -1 {
		castles = []*direction{
			{0, 1},
			{0, -1},
		}
	} else if yDir == 1 || yDir == -1 {
		castles = []*direction{
			{1, 0},
			{-1, 0},
		}
	} else {
		return nil, fmt.Errorf("invalid direction")
	}

	return &king{
		allegiant{color},
		moved,
		castles,
	}, nil
}

type king struct {
	allegiant
	moved   bool
	castles []*direction
}

func (k *king) print() string {
	return "K"
}

func (k *king) movedCopy() piece {
	return &king{
		allegiant{k.color},
		true,
		k.castles,
	}
}

func (k *king) moves(b board, x int, y int) []move {
	moves := &[]move{}
	k.addSimples(b, x, y, moves)
	k.addCastles(b, x, y, moves)
	return *moves
}

func (k *king) addSimples(b board, xFrom int, yFrom int, moves *[]move) {
	for _, simple := range kingSimples {
		addSimple(xFrom, yFrom, b, moves, simple.dx, simple.dy, k.color)
	}
}

func (k *king) addCastles(b board, xFrom int, yFrom int, moves *[]move) {
	if k.moved {
		return
	}

	for _, castle := range k.castles {
		xCurrent := xFrom + castle.dx
		yCurrent := yFrom + castle.dy

		for !movedOutOfBounds(xCurrent, yCurrent, b) {
			if piece, err := b.getPiece(xCurrent, yCurrent); err != nil {
				break
			} else if piece == nil {
				xCurrent += castle.dx
				yCurrent += castle.dy
				continue
			} else if rook, ok := piece.(*rook); ok && !rook.moved && rook.color == k.color {
				castleMove, err := moveFactoryInstance.newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
				if err == nil {
					*moves = append(*moves, castleMove)
				}

				xCurrent += castle.dx
				yCurrent += castle.dy
				continue
			} else {
				break
			}
		}
	}
}
