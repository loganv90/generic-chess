package chess

import (
    "fmt"
)

type Action struct {
	b Board
    fromLocation Point
    toLocation Point
}

func (a Action) getAction() Action {
	return a
}

type EnPassantCapture struct {
	enPassant EnPassant
	capturedPiece Piece
}

var moveFactoryInstance = MoveFactory(&ConcreteMoveFactory{})

type MoveFactory interface {
	newSimpleMove(b Board, fromLocation Point, toLocation Point) (*SimpleMove, error)
	newRevealEnPassantMove(b Board, fromLocation Point, toLocation Point, target Point) (*RevealEnPassantMove, error)
	newCaptureEnPassantMove(b Board, fromLocation Point, toLocation Point) (*CaptureEnPassantMove, error)
	newCastleMove(b Board, fromLocation Point, toLocation Point, toKingLocation Point, toRookLocation Point, newVulnerable Vulnerable) (*CastleMove, error)
    newPromotionMove(move Move) (*PromotionMove, error)
    newAllyDefenseMove(b Board, fromLocation Point, toLocation Point) (*AllyDefenseMove, error)
}

type ConcreteMoveFactory struct{}

func (f *ConcreteMoveFactory) newSimpleMove(b Board, fromLocation Point, toLocation Point) (*SimpleMove, error) {
	piece, ok := b.getPiece(fromLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at fromLocation")
	}

	newPiece, err := piece.copy()
    if err != nil {
        return nil, err
    }

	capturedPiece, ok := b.getPiece(toLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at toLocation")
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

    vulnerable, err := b.getVulnerable(piece.getColor())
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
        vulnerable,
	}, nil
}

func (f *ConcreteMoveFactory) newRevealEnPassantMove(b Board, fromLocation Point, toLocation Point, target Point) (*RevealEnPassantMove, error) {
	piece, ok := b.getPiece(fromLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at fromLocation")
	}

	newPiece, err := piece.copy()
    if err != nil {
        return nil, err
    }

	capturedPiece, ok := b.getPiece(toLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at toLocation")
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}

    vulnerable, err := b.getVulnerable(piece.getColor())
    if err != nil {
        return nil, err
    }

	newEn := EnPassant{
        target: target,
        risk: toLocation,
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
        vulnerable,
	}, nil
}

func (f *ConcreteMoveFactory) newCaptureEnPassantMove(b Board, fromLocation Point, toLocation Point) (*CaptureEnPassantMove, error) {
	piece, ok := b.getPiece(fromLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at fromLocation")
	}

	newPiece, err := piece.copy()
    if err != nil {
        return nil, err
    }

	capturedPiece, ok := b.getPiece(toLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at toLocation")
	}

	en, err := b.getEnPassant(piece.getColor())
	if err != nil {
		return nil, err
	}
    
    vulnerable, err := b.getVulnerable(piece.getColor())
    if err != nil {
        return nil, err
    }

	encs := []EnPassantCapture{}
    possibleEnPassant, err := b.possibleEnPassant(piece.getColor(), toLocation)
    if err == nil {
        for _, enPassant := range possibleEnPassant {
            capturedPiece, ok := b.getPiece(enPassant.risk)
            if !ok {
                return nil, fmt.Errorf("no piece at enPassant.pieceLocation")
            }

            encs = append(encs, EnPassantCapture{
                enPassant,
                capturedPiece,
            })
        }
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
        vulnerable,
	}, nil
}

func (f *ConcreteMoveFactory) newCastleMove(b Board, fromLocation Point, toLocation Point, toKingLocation Point, toRookLocation Point, newVulnerable Vulnerable) (*CastleMove, error) {
	king, ok := b.getPiece(fromLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at fromLocation")
	}

	newKing, err := king.copy()
    if err != nil {
        return nil, err
    }

	rook, ok := b.getPiece(toLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at toLocation")
	}

	newRook, err := rook.copy()
    if err != nil {
        return nil, err
    }

	en, err := b.getEnPassant(king.getColor())
	if err != nil {
		return nil, err
	}

    vulnerable, err := b.getVulnerable(king.getColor())
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
        vulnerable,
        newVulnerable,
	}, nil
}

func (f *ConcreteMoveFactory) newPromotionMove(move Move) (*PromotionMove, error) {
    action := move.getAction()

    promotionMove := &PromotionMove{
        Action: action,
        baseMove: move,
        promotionPiece: nil,
    }

    return promotionMove, nil
}

func (m *ConcreteMoveFactory) newAllyDefenseMove(b Board, fromLocation Point, toLocation Point) (*AllyDefenseMove, error) {
	p, ok := b.getPiece(fromLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at fromLocation")
	}

    return &AllyDefenseMove{
        Action{
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
	getAction() Action
    getNewPiece() Piece
}

type SimpleMove struct {
    Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            EnPassant
    vulnerable    Vulnerable
}

func (s *SimpleMove) getNewPiece() Piece {
    return s.newPiece
}

func (s *SimpleMove) execute() error {
	ok := s.b.setPiece(s.fromLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = s.b.setPiece(s.toLocation, s.newPiece)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    s.b.setVulnerable(s.piece.getColor(), Vulnerable{Point{-1, -1}, Point{-1, -1}})
	s.b.setEnPassant(s.piece.getColor(), EnPassant{Point{-1, -1}, Point{-1, -1}})

	return nil
}

func (s *SimpleMove) undo() error {
	ok := s.b.setPiece(s.fromLocation, s.piece)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = s.b.setPiece(s.toLocation, s.capturedPiece)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    s.b.setVulnerable(s.piece.getColor(), s.vulnerable)
	s.b.setEnPassant(s.piece.getColor(), s.en)

	return nil
}

type RevealEnPassantMove struct {
	Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            EnPassant
	newEn         EnPassant
    vulnerable    Vulnerable
}

func (r *RevealEnPassantMove) getNewPiece() Piece {
    return r.newPiece
}

func (r *RevealEnPassantMove) execute() error {
	ok := r.b.setPiece(r.fromLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = r.b.setPiece(r.toLocation, r.newPiece)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    r.b.setVulnerable(r.piece.getColor(), Vulnerable{Point{-1, -1}, Point{-1, -1}})
	r.b.setEnPassant(r.piece.getColor(), r.newEn)

	return nil
}

func (r *RevealEnPassantMove) undo() error {
	ok := r.b.setPiece(r.fromLocation, r.piece)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = r.b.setPiece(r.toLocation, r.capturedPiece)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    r.b.setVulnerable(r.piece.getColor(), r.vulnerable)
	r.b.setEnPassant(r.piece.getColor(), r.en)

	return nil
}

type CaptureEnPassantMove struct {
	Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            EnPassant
	encs          []EnPassantCapture
    vulnerable    Vulnerable
}

func (c *CaptureEnPassantMove) getNewPiece() Piece {
    return c.newPiece
}

func (c *CaptureEnPassantMove) execute() error {
	ok := c.b.setPiece(c.fromLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	for _, enc := range c.encs {
        ok = c.b.setPiece(enc.enPassant.risk, nil)
		if !ok {
			return fmt.Errorf("no piece at enPassant.pieceLocation")
		}
	}

	ok = c.b.setPiece(c.toLocation, c.newPiece)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    c.b.setVulnerable(c.piece.getColor(), Vulnerable{Point{-1, -1}, Point{-1, -1}})
	c.b.setEnPassant(c.piece.getColor(), EnPassant{Point{-1, -1}, Point{-1, -1}})

	return nil
}

func (c *CaptureEnPassantMove) undo() error {
	ok := c.b.setPiece(c.fromLocation, c.piece)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = c.b.setPiece(c.toLocation, c.capturedPiece)
    if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

	for _, enc := range c.encs {
        ok = c.b.setPiece(enc.enPassant.risk, enc.capturedPiece)
		if !ok {
			return fmt.Errorf("no piece at enPassant.pieceLocation")
		}
	}

    c.b.setVulnerable(c.piece.getColor(), c.vulnerable)
	c.b.setEnPassant(c.piece.getColor(), c.en)

	return nil
}

type CastleMove struct {
	Action
	king    Piece
	newKing Piece
    toKingLocation Point
	rook    Piece
	newRook Piece
    toRookLocation Point
	en      EnPassant
    vulnerable Vulnerable
    newVulnerable Vulnerable
}

func (c *CastleMove) getNewPiece() Piece {
    return c.newKing
}

func (c *CastleMove) execute() error {
	ok := c.b.setPiece(c.fromLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = c.b.setPiece(c.toLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

	ok = c.b.setPiece(c.toRookLocation, c.newRook)
	if !ok {
		return fmt.Errorf("no piece at toRookLocation")
	}

	ok = c.b.setPiece(c.toKingLocation, c.newKing)
	if !ok {
		return fmt.Errorf("no piece at toKingLocation")
	}

    c.b.setVulnerable(c.king.getColor(), c.newVulnerable)
	c.b.setEnPassant(c.king.getColor(), EnPassant{Point{-1, -1}, Point{-1, -1}})

	return nil
}

func (c *CastleMove) undo() error {
    ok := c.b.setPiece(c.toRookLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at toRookLocation")
	}

	ok = c.b.setPiece(c.toKingLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at toKingLocation")
	}

	ok = c.b.setPiece(c.fromLocation, c.king)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	ok = c.b.setPiece(c.toLocation, c.rook)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    c.b.setVulnerable(c.king.getColor(), c.vulnerable)
	c.b.setEnPassant(c.king.getColor(), c.en)

	return nil
}

type PromotionMove struct {
    Action
    baseMove Move
    promotionPiece Piece
}

func (p *PromotionMove) getNewPiece() Piece {
    return p.baseMove.getNewPiece()
}

func (p *PromotionMove) setPromotionPiece(piece Piece) error {
    p.promotionPiece = piece

    return nil
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
    Action
    piece Piece
}

func (m *AllyDefenseMove) getNewPiece() Piece {
    return m.piece
}

func (m *AllyDefenseMove) execute() error {
	return nil
}

func (m *AllyDefenseMove) undo() error {
	return nil
}

