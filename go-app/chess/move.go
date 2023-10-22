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

type enPassantCapture struct {
	enPassant     *enPassant
	capturedPiece piece
}

func getMoveFromSlice(moves []move, xTo int, yTo int) move {
	for _, m := range moves {
		if action := m.getAction(); action.xTo == xTo && action.yTo == yTo {
			return m
		}
	}

	return nil
}

var moveFactoryInstance = moveFactory(&concreteMoveFactory{})

type moveFactory interface {
	newSimpleMove(b board, xFrom int, yFrom int, xTo int, yTo int) (*simpleMove, error)
	newRevealEnPassantMove(b board, xFrom int, yFrom int, xTo int, yTo int, xTarget int, yTarget int) (*revealEnPassantMove, error)
	newCaptureEnPassantMove(b board, xFrom int, yFrom int, xTo int, yTo int) (*captureEnPassantMove, error)
	newCastleMove(b board, xFrom int, yFrom int, xTo int, yTo int, xToKing int, yToKing int, xToRook int, yToRook int) (*castleMove, error)
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

	en, err := b.getEnPassant(piece.getColor())
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
		en,
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

func (f *concreteMoveFactory) newCaptureEnPassantMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
) (*captureEnPassantMove, error) {
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

	encs := []*enPassantCapture{}

	for _, enPassant := range b.possibleEnPassants(piece.getColor(), xTo, yTo) {
		capturedPiece, err := b.getPiece(enPassant.xPiece, enPassant.yPiece)
		if err != nil {
			return nil, err
		}

		encs = append(encs, &enPassantCapture{
			enPassant,
			capturedPiece,
		})
	}

	return &captureEnPassantMove{
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
		encs,
	}, nil
}

func (f *concreteMoveFactory) newCastleMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
	xToKing int,
	yToKing int,
	xToRook int,
	yToRook int,
) (*castleMove, error) {
	king, err := b.getPiece(xFrom, yFrom)
	if err != nil {
		return nil, err
	}

	newKing := king.movedCopy()

	rook, err := b.getPiece(xTo, yTo)
	if err != nil {
		return nil, err
	}

	newRook := rook.movedCopy()

	en, err := b.getEnPassant(king.getColor())
	if err != nil {
		return nil, err
	}

	return &castleMove{
		action{
			b:     b,
			xFrom: xFrom,
			yFrom: yFrom,
			xTo:   xTo,
			yTo:   yTo,
		},
		king,
		newKing,
		xToKing,
		yToKing,
		rook,
		newRook,
		xToRook,
		yToRook,
		en,
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
	en            *enPassant
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

	s.b.setEnPassant(s.piece.getColor(), s.en)
	s.b.decrement()

	return nil
}

type revealEnPassantMove struct {
	action
	piece         piece
	newPiece      piece
	capturedPiece piece
	en            *enPassant
	newEn         *enPassant
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

	r.b.setEnPassant(r.piece.getColor(), r.newEn)
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

	r.b.setEnPassant(r.piece.getColor(), r.en)
	r.b.decrement()

	return nil
}

type captureEnPassantMove struct {
	action
	piece         piece
	newPiece      piece
	capturedPiece piece
	en            *enPassant
	encs          []*enPassantCapture
}

func (c *captureEnPassantMove) execute() error {
	err := c.b.setPiece(c.xFrom, c.yFrom, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xTo, c.yTo, c.newPiece)
	if err != nil {
		return err
	}

	for _, enc := range c.encs {
		err = c.b.setPiece(enc.enPassant.xPiece, enc.enPassant.yPiece, nil)
		if err != nil {
			return err
		}
	}

	c.b.clrEnPassant(c.piece.getColor())
	c.b.increment()

	return nil
}

func (c *captureEnPassantMove) undo() error {
	err := c.b.setPiece(c.xFrom, c.yFrom, c.piece)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xTo, c.yTo, c.capturedPiece)
	if err != nil {
		return err
	}

	for _, enc := range c.encs {
		err = c.b.setPiece(enc.enPassant.xPiece, enc.enPassant.yPiece, enc.capturedPiece)
		if err != nil {
			return err
		}
	}

	c.b.setEnPassant(c.piece.getColor(), c.en)
	c.b.decrement()

	return nil
}

type castleMove struct {
	action
	king    piece
	newKing piece
	xToKing int
	yToKing int
	rook    piece
	newRook piece
	xToRook int
	yToRook int
	en      *enPassant
}

func (c *castleMove) execute() error {
	err := c.b.setPiece(c.xFrom, c.yFrom, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xTo, c.yTo, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xToRook, c.yToRook, c.newRook)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xToKing, c.yToKing, c.newKing)
	if err != nil {
		return err
	}

	c.b.clrEnPassant(c.king.getColor())
	c.b.increment()

	return nil

}

func (c *castleMove) undo() error {
	err := c.b.setPiece(c.xFrom, c.yFrom, c.king)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xTo, c.yTo, c.rook)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xToRook, c.yToRook, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.xToKing, c.yToKing, nil)
	if err != nil {
		return err
	}

	c.b.setEnPassant(c.king.getColor(), c.en)
	c.b.decrement()

	return nil
}
