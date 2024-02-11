package chess

import (
    "fmt"
)

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

var moveFactoryInstance = MoveFactory(&ConcreteMoveFactory{})

type MoveFactory interface {
	newSimpleMove(b Board, fromLocation *Point, toLocation *Point) (*SimpleMove, error)
	newRevealEnPassantMove(b Board, fromLocation *Point, toLocation *Point, target *Point) (*RevealEnPassantMove, error)
	newCaptureEnPassantMove(b Board, fromLocation *Point, toLocation *Point) (*CaptureEnPassantMove, error)
	newCastleMove(b Board, fromLocation *Point, toLocation *Point, toKingLocation *Point, toRookLocation *Point, newVulnerables []*Point) (*CastleMove, error)
    newPromotionMoves(move Move, promotionPieces []Piece) ([]*PromotionMove, error)
    newAllyDefenseMove(b Board, fromLocation *Point, toLocation *Point) (*AllyDefenseMove, error)
}

type ConcreteMoveFactory struct{}

func (f *ConcreteMoveFactory) newSimpleMove(b Board, fromLocation *Point, toLocation *Point) (*SimpleMove, error) {
	piece, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newPiece := piece.copy()
    newPiece.setMoved()

	capturedPiece, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

    vulnerables, err := b.getVulnerables(piece.getColor())
    if err != nil {
        return nil, err
    }

	return &SimpleMove{
		&Action{
			b: b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
        vulnerables,
	}, nil
}

func (f *ConcreteMoveFactory) newRevealEnPassantMove(b Board, fromLocation *Point, toLocation *Point, target *Point) (*RevealEnPassantMove, error) {
	piece, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newPiece := piece.copy()
    newPiece.setMoved()

	capturedPiece, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

    vulnerables, err := b.getVulnerables(piece.getColor())
    if err != nil {
        return nil, err
    }

	newEn := &EnPassant{
        target: target,
        pieceLocation: toLocation,
	}

	return &RevealEnPassantMove{
		&Action{
			b:     b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
		newEn,
        vulnerables,
	}, nil
}

func (f *ConcreteMoveFactory) newCaptureEnPassantMove(b Board, fromLocation *Point, toLocation *Point) (*CaptureEnPassantMove, error) {
	piece, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newPiece := piece.copy()
    newPiece.setMoved()

	capturedPiece, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}
    
    vulnerables, err := b.getVulnerables(piece.getColor())
    if err != nil {
        return nil, err
    }

	encs := []*EnPassantCapture{}
    possibleEnPassant, err := b.possibleEnPassant(piece.getColor(), toLocation)
    if err == nil {
        for _, enPassant := range possibleEnPassant {
            capturedPiece, err := b.getPiece(enPassant.pieceLocation)
            if err != nil {
                return nil, err
            }

            encs = append(encs, &EnPassantCapture{
                enPassant,
                capturedPiece,
            })
        }
    }

	return &CaptureEnPassantMove{
		&Action{
			b:     b,
            fromLocation: fromLocation,
            toLocation: toLocation,
		},
		piece,
		newPiece,
		capturedPiece,
		en,
		encs,
        vulnerables,
	}, nil
}

func (f *ConcreteMoveFactory) newCastleMove(b Board, fromLocation *Point, toLocation *Point, toKingLocation *Point, toRookLocation *Point, newVulnerables []*Point) (*CastleMove, error) {
	king, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

	newKing := king.copy()
    newKing.setMoved()

	rook, err := b.getPiece(toLocation)
	if err != nil {
		return nil, err
	}

	newRook := rook.copy()
    newRook.setMoved()

	en, err := b.getEnPassant(king.getColor())
	if err != nil {
		return nil, err
	}

    vulnerables, err := b.getVulnerables(king.getColor())
    if err != nil {
        return nil, err
    }

	return &CastleMove{
		&Action{
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
        vulnerables,
        newVulnerables,
	}, nil
}

func (f *ConcreteMoveFactory) newPromotionMoves(move Move, promotionPieces []Piece) ([]*PromotionMove, error) {
    action := move.getAction()

    promotionMoves := []*PromotionMove{}
    for _, promotionPiece := range promotionPieces {
        promotionMove := &PromotionMove{
            Action: action,
            baseMove: move,
            promotionPiece: promotionPiece,
        }

        promotionMoves = append(promotionMoves, promotionMove)
    }

    return promotionMoves, nil
}

func (m *ConcreteMoveFactory) newAllyDefenseMove(b Board, fromLocation *Point, toLocation *Point) (*AllyDefenseMove, error) {
	p, err := b.getPiece(fromLocation)
	if err != nil {
		return nil, err
	}

    return &AllyDefenseMove{
        &Action{
            b: b,
            fromLocation: fromLocation,
            toLocation: toLocation,
        },
        p,
    }, nil
}

type Move interface {
	execute() error
	undo() error
	getAction() *Action
    getNewPiece() Piece
}

type SimpleMove struct {
    *Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            *EnPassant
    vulnerables   []*Point
}

func (s *SimpleMove) getNewPiece() Piece {
    return s.newPiece
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

    s.b.setVulnerables(s.piece.getColor(), []*Point{})
	s.b.clearEnPassant(s.piece.getColor())

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

    s.b.setVulnerables(s.piece.getColor(), s.vulnerables)
	s.b.setEnPassant(s.piece.getColor(), s.en)

	return nil
}

type RevealEnPassantMove struct {
	*Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            *EnPassant
	newEn         *EnPassant
    vulnerables   []*Point
}

func (r *RevealEnPassantMove) getNewPiece() Piece {
    return r.newPiece
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

    r.b.setVulnerables(r.piece.getColor(), []*Point{})
	r.b.setEnPassant(r.piece.getColor(), r.newEn)

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

    r.b.setVulnerables(r.piece.getColor(), r.vulnerables)
	r.b.setEnPassant(r.piece.getColor(), r.en)

	return nil
}

type CaptureEnPassantMove struct {
	*Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            *EnPassant
	encs          []*EnPassantCapture
    vulnerables   []*Point
}

func (c *CaptureEnPassantMove) getNewPiece() Piece {
    return c.newPiece
}

func (c *CaptureEnPassantMove) execute() error {
	err := c.b.setPiece(c.fromLocation, nil)
	if err != nil {
		return err
	}

	for _, enc := range c.encs {
        err = c.b.setPiece(enc.enPassant.pieceLocation, nil)
		if err != nil {
			return err
		}
	}

	err = c.b.setPiece(c.toLocation, c.newPiece)
	if err != nil {
		return err
	}

    c.b.setVulnerables(c.piece.getColor(), []*Point{})
	c.b.clearEnPassant(c.piece.getColor())

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

    c.b.setVulnerables(c.piece.getColor(), c.vulnerables)
	c.b.setEnPassant(c.piece.getColor(), c.en)

	return nil
}

type CastleMove struct {
	*Action
	king    Piece
	newKing Piece
    toKingLocation *Point
	rook    Piece
	newRook Piece
    toRookLocation *Point
	en      *EnPassant
    vulnerables   []*Point
    newVulnerables []*Point
}

func (c *CastleMove) getNewPiece() Piece {
    return c.newKing
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

    c.b.setVulnerables(c.king.getColor(), c.newVulnerables)
	c.b.clearEnPassant(c.king.getColor())

	return nil
}

func (c *CastleMove) undo() error {
    err := c.b.setPiece(c.toRookLocation, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toKingLocation, nil)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.fromLocation, c.king)
	if err != nil {
		return err
	}

	err = c.b.setPiece(c.toLocation, c.rook)
	if err != nil {
		return err
	}

    c.b.setVulnerables(c.king.getColor(), c.vulnerables)
	c.b.setEnPassant(c.king.getColor(), c.en)

	return nil
}

type PromotionMove struct {
    *Action
    baseMove Move
    promotionPiece Piece
}

func (p *PromotionMove) getNewPiece() Piece {
    return p.promotionPiece
}

func (p *PromotionMove) execute() error {
    p.baseMove.execute()
    p.b.setPiece(p.toLocation, p.promotionPiece)

	return nil
}

func (p *PromotionMove) undo() error {
    p.b.setPiece(p.toLocation, p.baseMove.getNewPiece())
    p.baseMove.undo()
    
	return nil
}

type AllyDefenseMove struct {
    *Action
    piece Piece
}

func (m *AllyDefenseMove) getNewPiece() Piece {
    return m.piece
}

func (m *AllyDefenseMove) execute() error {
	return fmt.Errorf("AllyDefenseMove cannot be executed")
}

func (m *AllyDefenseMove) undo() error {
	return fmt.Errorf("AllyDefenseMove cannot be undone")
}

