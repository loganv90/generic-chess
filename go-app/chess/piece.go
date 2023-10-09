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

func movedOutOfBounds(xTo int, yTo int, b board) bool {
	return xTo < 0 ||
		xTo >= b.xLen() ||
		yTo < 0 ||
		yTo >= b.yLen()
}

func movedToPromotionSquare(xTo int, yTo int, xFrom int, yFrom int, b board) bool {
	return (xTo == 0 && xFrom != 0) ||
		(xTo == b.xLen()-1 && xFrom != b.xLen()-1) ||
		(yTo == 0 && yFrom != 0) ||
		(yTo == b.yLen()-1 && yFrom != b.yLen()-1)
}

func addDirection(xFrom int, yFrom int, b board, moves []move, dx int, dy int, color string) {
	xCurrent := xFrom + dx
	yCurrent := yFrom + dy

	for !movedOutOfBounds(xCurrent, yCurrent, b) {
		p, err := b.getPiece(xCurrent, yCurrent)
		if err != nil {
			break
		}

		if p == nil {
			simpleMove, err := newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
			if err == nil {
				moves = append(moves, simpleMove)
			}
		} else if p.getColor() != color {
			simpleMove, err := newSimpleMove(b, xFrom, yFrom, xCurrent, yCurrent)
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

func addSimple(xFrom int, yFrom int, b board, moves []move, dx int, dy int, color string) {
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
		simpleMove, err := newSimpleMove(b, xFrom, yFrom, xTo, yTo)
		if err == nil {
			moves = append(moves, simpleMove)
		}
	} else if p.getColor() != color {
		simpleMove, err := newSimpleMove(b, xFrom, yFrom, xTo, yTo)
		if err == nil {
			moves = append(moves, simpleMove)
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
			&direction{xDir, 1},
			&direction{xDir, -1},
		}
	} else if yDir == 1 || yDir == -1 {
		forward1 = &direction{0, yDir}
		forward2 = &direction{0, yDir * 2}
		captures = []*direction{
			&direction{1, yDir},
			&direction{-1, yDir},
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

func (a *pawn) movedCopy() piece {
	return &pawn{
		allegiant{a.color},
		true,
		a.forward1,
		a.forward2,
		a.captures,
	}
}

func (a *pawn) moves(b board, x int, y int) []move {
	return []move{}
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
