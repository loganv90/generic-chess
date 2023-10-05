package chess

type piece interface {
	getMoved() piece
}

type pawn struct {
	moved bool
	xDir  int
	yDir  int
}

func (p *pawn) getMoved() piece {
	return &pawn{
		moved: true,
		xDir:  p.xDir,
		yDir:  p.yDir,
	}
}

type knight struct{}

func (k *knight) getMoved() piece {
	return &knight{}
}

type bishop struct{}

func (b *bishop) getMoved() piece {
	return &bishop{}
}

type rook struct{}

func (r *rook) getMoved() piece {
	return &rook{}
}

type queen struct{}

func (q *queen) getMoved() piece {
	return &queen{}
}

type king struct{}

func (k *king) getMoved() piece {
	return &king{}
}
