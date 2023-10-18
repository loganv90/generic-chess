package chess

type action struct {
	b     board
	xFrom int
	yFrom int
	xTo   int
	yTo   int
}

func (a *action) getAction() *action {
	return a
}

var moveFactoryInstance = moveFactory(&concreteMoveFactory{})

type moveFactory interface {
	newSimpleMove(b board, xFrom int, yFrom int, xTo int, yTo int) (*simpleMove, error)
	newRevealEnPassantMove(b board, xFrom int, yFrom int, xTo int, yTo int, xTarget int, yTarget int) (*revealEnPassantMove, error)
}

type concreteMoveFactory struct{}

func (f *concreteMoveFactory) newSimpleMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
) (*simpleMove, error) {
	piece, err := b.getPiece(xFrom, yFrom)
	if err != nil {
		return nil, err
	}

	newPiece := piece.movedCopy()

	capturedPiece, err := b.getPiece(xTo, yTo)
	if err != nil {
		return nil, err
	}

	enPassant, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

	return &simpleMove{
		action{
			b:     b,
			xFrom: xFrom,
			yFrom: yFrom,
			xTo:   xTo,
			yTo:   yTo,
		},
		piece,
		newPiece,
		capturedPiece,
		enPassant,
	}, nil
}

func (f *concreteMoveFactory) newRevealEnPassantMove(b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
	xTarget int,
	yTarget int,
) (*revealEnPassantMove, error) {
	piece, err := b.getPiece(xFrom, yFrom)
	if err != nil {
		return nil, err
	}

	newPiece := piece.movedCopy()

	capturedPiece, err := b.getPiece(xTo, yTo)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

	newEn := &enPassant{
		xTarget: xTarget,
		yTarget: yTarget,
		xPiece:  xTo,
		yPiece:  yTo,
	}

	return &revealEnPassantMove{
		action{
			b:     b,
			xFrom: xFrom,
			yFrom: yFrom,
			xTo:   xTo,
			yTo:   yTo,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
		newEn,
	}, nil
}

type move interface {
	execute() error
	undo() error
	getAction() *action
}

type simpleMove struct {
	action
	piece         piece
	newPiece      piece
	capturedPiece piece
	enPassant     *enPassant
}

func (s *simpleMove) execute() error {
	err := s.b.setPiece(s.xFrom, s.yFrom, nil)
	if err != nil {
		return err
	}

	err = s.b.setPiece(s.xTo, s.yTo, s.newPiece)
	if err != nil {
		return err
	}

	s.b.clrEnPassant(s.piece.getColor())
	s.b.increment()

	return nil
}

func (s *simpleMove) undo() error {
	err := s.b.setPiece(s.xFrom, s.yFrom, s.piece)
	if err != nil {
		return err
	}

	err = s.b.setPiece(s.xTo, s.yTo, s.capturedPiece)
	if err != nil {
		return err
	}

	s.b.setEnPassant(s.piece.getColor(), s.enPassant)
	s.b.decrement()

	return nil
}

type revealEnPassantMove struct {
	action
	piece         piece
	newPiece      piece
	capturedPiece piece
	enPassant     *enPassant
	newEnPassant  *enPassant
}

func (r *revealEnPassantMove) execute() error {
	err := r.b.setPiece(r.xFrom, r.yFrom, nil)
	if err != nil {
		return err
	}

	err = r.b.setPiece(r.xTo, r.yTo, r.newPiece)
	if err != nil {
		return err
	}

	r.b.setEnPassant(r.piece.getColor(), r.newEnPassant)
	r.b.increment()

	return nil
}

func (r *revealEnPassantMove) undo() error {
	err := r.b.setPiece(r.xFrom, r.yFrom, r.piece)
	if err != nil {
		return err
	}

	err = r.b.setPiece(r.xTo, r.yTo, r.capturedPiece)
	if err != nil {
		return err
	}

	r.b.setEnPassant(r.piece.getColor(), r.enPassant)
	r.b.decrement()

	return nil
}

func getMoveFromSlice(moves []move, xTo int, yTo int) move {
	for _, m := range moves {
		if action := m.getAction(); action.xTo == xTo && action.yTo == yTo {
			return m
		}
	}

	return nil
}
