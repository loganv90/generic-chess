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
	moves []move,
	dx int,
	dy int,
	color string,
	f moveFactory,
) {
	xCurrent := xFrom + dx
	yCurrent := yFrom + dy

	for !movedOutOfBounds(xCurrent, yCurrent, b) {
		p, err := b.getPiece(xCurrent, yCurrent)
		if err != nil {
			break
		}

		if p == nil {
			simpleMove, err := f.newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
			if err == nil {
				moves = append(moves, simpleMove)
			}
		} else if p.getColor() != color {
			simpleMove, err := f.newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
			if err == nil {
				moves = append(moves, simpleMove)
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
	moves []move,
	dx int,
	dy int,
	color string,
	f moveFactory,
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
		simpleMove, err := f.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
		if err == nil {
			moves = append(moves, simpleMove)
		}
	} else if p.getColor() != color {
		simpleMove, err := f.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
		if err == nil {
			moves = append(moves, simpleMove)
		}
	}
}

func newPawn(f moveFactory, color string, moved bool, xDir int, yDir int) (*pawn, error) {
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
		f,
	}, nil
}

type pawn struct {
	allegiant
	moved    bool
	forward1 *direction
	forward2 *direction
	captures []*direction
	f        moveFactory
}

func (a *pawn) movedCopy() piece {
	return &pawn{
		allegiant{a.color},
		true,
		a.forward1,
		a.forward2,
		a.captures,
		a.f,
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
		promotionMove, err := a.f.newSimpleMove(b, xFrom, yFrom, xTo1, yTo1)
		if err == nil {
			*moves = append(*moves, promotionMove)
		}
	} else {
		simpleMove, err := a.f.newSimpleMove(b, xFrom, yFrom, xTo1, yTo1)
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
		promotionMove, err := a.f.newSimpleMove(b, xFrom, yFrom, xTo2, yTo2)
		if err == nil {
			*moves = append(*moves, promotionMove)
		}
	} else {
		simpleMove, err := a.f.newSimpleMove(b, xFrom, yFrom, xTo2, yTo2)
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
			captureEnPassantMove, err := a.f.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
			if err == nil {
				*moves = append(*moves, captureEnPassantMove)
			}
		} else if piece != nil && piece.getColor() != a.color {
			simpleMove, err := a.f.newSimpleMove(b, xFrom, yFrom, xTo, yTo)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
		}
	}
}

type knight struct {
	allegiant
}

func (n *knight) movedCopy() piece {
	return &knight{
		allegiant{n.color},
	}
}

func (n *knight) moves(b board, x int, y int) []move {
	return []move{}
}

type bishop struct {
	allegiant
}

func (s *bishop) movedCopy() piece {
	return &bishop{
		allegiant{s.color},
	}
}

func (s *bishop) moves(b board, x int, y int) []move {
	return []move{}
}

type rook struct {
	allegiant
}

func (r *rook) movedCopy() piece {
	return &rook{
		allegiant{r.color},
	}
}

func (r *rook) moves(b board, x int, y int) []move {
	return []move{}
}

type queen struct {
	allegiant
}

func (q *queen) movedCopy() piece {
	return &queen{
		allegiant{q.color},
	}
}

func (q *queen) moves(b board, x int, y int) []move {
	return []move{}
}

type king struct {
	allegiant
}

func (k *king) movedCopy() piece {
	return &king{
		allegiant{k.color},
	}
}

func (k *king) moves(b board, x int, y int) []move {
	return []move{}
}
