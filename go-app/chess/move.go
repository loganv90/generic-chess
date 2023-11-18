package chess

type Action struct {
	b Board
    fromLocation *Point
    toLocation *Point
}

func (a *Action) getAction() *Action {
	return a
}

type EnPassantCapture struct {
	enPassant     *EnPassant
	capturedPiece Piece
}

func getMoveFromSlice(moves []Move, toLocation *Point) Move {
	for _, m := range moves {
        actionToLocation := m.getAction().toLocation
        if actionToLocation.equals(toLocation) {
			return m
		}
	}

	return nil
}

var moveFactoryInstance = MoveFactory(&ConcreteMoveFactory{})

type MoveFactory interface {
	newSimpleMove(b Board, fromLocation *Point, toLocation *Point) (*SimpleMove, error)
	newRevealEnPassantMove(b Board, fromLocation *Point, toLocation *Point, target *Point) (*RevealEnPassantMove, error)
	newCaptureEnPassantMove(b Board, fromLocation *Point, toLocation *Point) (*CaptureEnPassantMove, error)
	newCastleMove(b Board, fromLocation *Point, toLocation *Point, toKingLocation *Point, toRookLocation *Point) (*CastleMove, error)
}

type ConcreteMoveFactory struct{}

func (f *ConcreteMoveFactory) newSimpleMove(b Board, fromLocation *Point, toLocation *Point) (*SimpleMove, error) {
	piece, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newPiece := piece.movedCopy()

	capturedPiece, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

	return &SimpleMove{
		Action{
			b: b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
	}, nil
}

func (f *ConcreteMoveFactory) newRevealEnPassantMove(b Board, fromLocation *Point, toLocation *Point, target *Point) (*RevealEnPassantMove, error) {
	piece, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newPiece := piece.movedCopy()

	capturedPiece, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

	newEn := &EnPassant{
        target: target,
        pieceLocation: toLocation,
	}

	return &RevealEnPassantMove{
		Action{
			b:     b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
		newEn,
	}, nil
}

func (f *ConcreteMoveFactory) newCaptureEnPassantMove(b Board, fromLocation *Point, toLocation *Point) (*CaptureEnPassantMove, error) {
	piece, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newPiece := piece.movedCopy()

	capturedPiece, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

	encs := []*EnPassantCapture{}

	for _, enPassant := range b.possibleEnPassants(piece.getColor(), toLocation) {
		capturedPiece, err := b.getPiece(enPassant.pieceLocation)
		if err != nil {
			return nil, err
		}

		encs = append(encs, &EnPassantCapture{
			enPassant,
			capturedPiece,
		})
	}

	return &CaptureEnPassantMove{
		Action{
			b:     b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
		encs,
	}, nil
}

func (f *ConcreteMoveFactory) newCastleMove(b Board, fromLocation *Point, toLocation *Point, toKingLocation *Point, toRookLocation *Point) (*CastleMove, error) {
	king, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newKing := king.movedCopy()

	rook, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	newRook := rook.movedCopy()

	en, err := b.getEnPassant(king.getColor())
	if err != nil {
		return nil, err
	}

	return &CastleMove{
		Action{
			b: b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		king,
		newKing,
        toKingLocation,
		rook,
		newRook,
        toRookLocation,
		en,
	}, nil
}

type Move interface {
	execute() error
	undo() error
	getAction() *Action
}

type SimpleMove struct {
	Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            *EnPassant
}

func (s *SimpleMove) execute() error {
	err := s.b.setPiece(s.fromLocation, nil)
	if err != nil {
		return err
	}

	err = s.b.setPiece(s.toLocation, s.newPiece)
	if err != nil {
		return err
	}

	s.b.clrEnPassant(s.piece.getColor())
	s.b.increment()

	return nil
}

func (s *SimpleMove) undo() error {
	err := s.b.setPiece(s.fromLocation, s.piece)
	if err != nil {
		return err
	}

	err = s.b.setPiece(s.toLocation, s.capturedPiece)
	if err != nil {
		return err
	}

	s.b.setEnPassant(s.piece.getColor(), s.en)
	s.b.decrement()

	return nil
}

type RevealEnPassantMove struct {
	Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            *EnPassant
	newEn         *EnPassant
}

func (r *RevealEnPassantMove) execute() error {
	err := r.b.setPiece(r.fromLocation, nil)
	if err != nil {
		return err
	}

	err = r.b.setPiece(r.toLocation, r.newPiece)
	if err != nil {
		return err
	}

	r.b.setEnPassant(r.piece.getColor(), r.newEn)
	r.b.increment()

	return nil
}

func (r *RevealEnPassantMove) undo() error {
	err := r.b.setPiece(r.fromLocation, r.piece)
	if err != nil {
		return err
	}

	err = r.b.setPiece(r.toLocation, r.capturedPiece)
	if err != nil {
		return err
	}

	r.b.setEnPassant(r.piece.getColor(), r.en)
	r.b.decrement()

	return nil
}

type CaptureEnPassantMove struct {
	Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            *EnPassant
	encs          []*EnPassantCapture
}

func (c *CaptureEnPassantMove) execute() error {
	err := c.b.setPiece(c.fromLocation, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toLocation, c.newPiece)
	if err != nil {
		return err
	}

	for _, enc := range c.encs {
        err = c.b.setPiece(enc.enPassant.pieceLocation, nil)
		if err != nil {
			return err
		}
	}

	c.b.clrEnPassant(c.piece.getColor())
	c.b.increment()

	return nil
}

func (c *CaptureEnPassantMove) undo() error {
	err := c.b.setPiece(c.fromLocation, c.piece)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toLocation, c.capturedPiece)
	if err != nil {
		return err
	}

	for _, enc := range c.encs {
        err = c.b.setPiece(enc.enPassant.pieceLocation, enc.capturedPiece)
		if err != nil {
			return err
		}
	}

	c.b.setEnPassant(c.piece.getColor(), c.en)
	c.b.decrement()

	return nil
}

type CastleMove struct {
	Action
	king    Piece
	newKing Piece
    toKingLocation *Point
	rook    Piece
	newRook Piece
    toRookLocation *Point
	en      *EnPassant
}

func (c *CastleMove) execute() error {
	err := c.b.setPiece(c.fromLocation, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toLocation, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toRookLocation, c.newRook)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toKingLocation, c.newKing)
	if err != nil {
		return err
	}

	c.b.clrEnPassant(c.king.getColor())
	c.b.increment()

	return nil

}

func (c *CastleMove) undo() error {
	err := c.b.setPiece(c.fromLocation, c.king)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toLocation, c.rook)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toRookLocation, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toKingLocation, nil)
	if err != nil {
		return err
	}

	c.b.setEnPassant(c.king.getColor(), c.en)
	c.b.decrement()

	return nil
}
