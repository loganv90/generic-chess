package chess

type piece interface {
	getColor() string
	movedCopy() piece
	moves(board, int, int) []move
}

type allegiant struct {
	color string
}

func (a *allegiant) getColor() string {
	return a.color
}

type pawn struct {
	allegiant
	moved bool
	xDir  int
	yDir  int
}

func (p *pawn) movedCopy() piece {
	return &pawn{
		allegiant{p.color},
		true,
		p.xDir,
		p.yDir,
	}
}

func (p *pawn) moves(b board, x int, y int) []move {
	return []move{}
}

type knight struct {
	allegiant
}

func (k *knight) movedCopy() piece {
	return &knight{
		allegiant{k.color},
	}
}

func (k *knight) moves(b board, x int, y int) []move {
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
