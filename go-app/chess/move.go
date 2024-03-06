package chess

import (
    "fmt"
    "sync"
)

var simpleMovePool = sync.Pool {
    New: func() interface{} {
        return &SimpleMove{}
    },
}

var revealEnPassantMovePool = sync.Pool {
    New: func() interface{} {
        return &RevealEnPassantMove{}
    },
}

var captureEnPassantMovePool = sync.Pool {
    New: func() interface{} {
        return &CaptureEnPassantMove{}
    },
}

var castleMovePool = sync.Pool {
    New: func() interface{} {
        return &CastleMove{}
    },
}

var promotionMovePool = sync.Pool {
    New: func() interface{} {
        return &PromotionMove{}
    },
}

var allyDefenseMovePool = sync.Pool {
    New: func() interface{} {
        return &AllyDefenseMove{}
    },
}

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
	newCastleMove(b Board, fromLocation Point, toLocation Point, toKingLocation Point, toRookLocation Point, newVulnerables []Point) (*CastleMove, error)
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

    vulnerables, err := b.getVulnerables(piece.getColor())
    if err != nil {
        return nil, err
    }

    simpleMove := simpleMovePool.Get().(*SimpleMove)
    simpleMove.b = b
    simpleMove.fromLocation = fromLocation
    simpleMove.toLocation = toLocation
    simpleMove.piece = piece
    simpleMove.newPiece = newPiece
    simpleMove.capturedPiece = capturedPiece
    simpleMove.en = en
    simpleMove.vulnerables = vulnerables

    return simpleMove, nil
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

    vulnerables, err := b.getVulnerables(piece.getColor())
    if err != nil {
        return nil, err
    }

	newEn := EnPassant{
        target: target,
        pieceLocation: toLocation,
	}

    revealEnPassantMove := revealEnPassantMovePool.Get().(*RevealEnPassantMove)
    revealEnPassantMove.b = b
    revealEnPassantMove.fromLocation = fromLocation
    revealEnPassantMove.toLocation = toLocation
    revealEnPassantMove.piece = piece
    revealEnPassantMove.newPiece = newPiece
    revealEnPassantMove.capturedPiece = capturedPiece
    revealEnPassantMove.en = en
    revealEnPassantMove.newEn = newEn
    revealEnPassantMove.vulnerables = vulnerables

    return revealEnPassantMove, nil
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
    
    vulnerables, err := b.getVulnerables(piece.getColor())
    if err != nil {
        return nil, err
    }

	encs := []EnPassantCapture{}
    possibleEnPassant, err := b.possibleEnPassant(piece.getColor(), toLocation)
    if err == nil {
        for _, enPassant := range possibleEnPassant {
            capturedPiece, ok := b.getPiece(enPassant.pieceLocation)
            if !ok {
                return nil, fmt.Errorf("no piece at enPassant.pieceLocation")
            }

            encs = append(encs, EnPassantCapture{
                enPassant,
                capturedPiece,
            })
        }
    }

    captureEnPassantMove := captureEnPassantMovePool.Get().(*CaptureEnPassantMove)
    captureEnPassantMove.b = b
    captureEnPassantMove.fromLocation = fromLocation
    captureEnPassantMove.toLocation = toLocation
    captureEnPassantMove.piece = piece
    captureEnPassantMove.newPiece = newPiece
    captureEnPassantMove.capturedPiece = capturedPiece
    captureEnPassantMove.en = en
    captureEnPassantMove.encs = encs
    captureEnPassantMove.vulnerables = vulnerables

    return captureEnPassantMove, nil
}

func (f *ConcreteMoveFactory) newCastleMove(b Board, fromLocation Point, toLocation Point, toKingLocation Point, toRookLocation Point, newVulnerables []Point) (*CastleMove, error) {
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

    vulnerables, err := b.getVulnerables(king.getColor())
    if err != nil {
        return nil, err
    }

    castleMove := castleMovePool.Get().(*CastleMove)
    castleMove.b = b
    castleMove.fromLocation = fromLocation
    castleMove.toLocation = toLocation
    castleMove.king = king
    castleMove.newKing = newKing
    castleMove.toKingLocation = toKingLocation
    castleMove.rook = rook
    castleMove.newRook = newRook
    castleMove.toRookLocation = toRookLocation
    castleMove.en = en
    castleMove.vulnerables = vulnerables
    castleMove.newVulnerables = newVulnerables

    return castleMove, nil
}

func (f *ConcreteMoveFactory) newPromotionMove(move Move) (*PromotionMove, error) {
    action := move.getAction()

    promotionMove := promotionMovePool.Get().(*PromotionMove)
    promotionMove.b = action.b
    promotionMove.fromLocation = action.fromLocation
    promotionMove.toLocation = action.toLocation
    promotionMove.baseMove = move
    promotionMove.promotionPiece = nil

    return promotionMove, nil
}

func (m *ConcreteMoveFactory) newAllyDefenseMove(b Board, fromLocation Point, toLocation Point) (*AllyDefenseMove, error) {
	p, ok := b.getPiece(fromLocation)
	if !ok {
		return nil, fmt.Errorf("no piece at fromLocation")
	}

    allyDefenseMove := allyDefenseMovePool.Get().(*AllyDefenseMove)
    allyDefenseMove.b = b
    allyDefenseMove.fromLocation = fromLocation
    allyDefenseMove.toLocation = toLocation
    allyDefenseMove.piece = p

    return allyDefenseMove, nil
}

type Move interface {
	execute() error
	undo() error
	getAction() Action
    getNewPiece() Piece
    copy() Move
    putInPool()
}

type SimpleMove struct {
    Action
	piece         Piece
	newPiece      Piece
	capturedPiece Piece
	en            EnPassant
    vulnerables   []Point
}

func (s *SimpleMove) getNewPiece() Piece {
    return s.newPiece
}

func (s *SimpleMove) putInPool() {
    simpleMovePool.Put(s)
}

func (s *SimpleMove) copy() Move {
    simpleMove := simpleMovePool.Get().(*SimpleMove)
    simpleMove.b = s.b
    simpleMove.fromLocation = s.fromLocation
    simpleMove.toLocation = s.toLocation
    simpleMove.piece = s.piece
    simpleMove.newPiece = s.newPiece
    simpleMove.capturedPiece = s.capturedPiece
    simpleMove.en = s.en
    simpleMove.vulnerables = s.vulnerables

    return simpleMove
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

    s.b.setVulnerables(s.piece.getColor(), []Point{})
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

    s.b.setVulnerables(s.piece.getColor(), s.vulnerables)
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
    vulnerables   []Point
}

func (r *RevealEnPassantMove) getNewPiece() Piece {
    return r.newPiece
}

func (r *RevealEnPassantMove) putInPool() {
    revealEnPassantMovePool.Put(r)
}

func (r *RevealEnPassantMove) copy() Move {
    revealEnPassantMove := revealEnPassantMovePool.Get().(*RevealEnPassantMove)
    revealEnPassantMove.b = r.b
    revealEnPassantMove.fromLocation = r.fromLocation
    revealEnPassantMove.toLocation = r.toLocation
    revealEnPassantMove.piece = r.piece
    revealEnPassantMove.newPiece = r.newPiece
    revealEnPassantMove.capturedPiece = r.capturedPiece
    revealEnPassantMove.en = r.en
    revealEnPassantMove.newEn = r.newEn
    revealEnPassantMove.vulnerables = r.vulnerables

    return revealEnPassantMove
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

    r.b.setVulnerables(r.piece.getColor(), []Point{})
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

    r.b.setVulnerables(r.piece.getColor(), r.vulnerables)
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
    vulnerables   []Point
}

func (c *CaptureEnPassantMove) getNewPiece() Piece {
    return c.newPiece
}

func (c *CaptureEnPassantMove) putInPool() {
    captureEnPassantMovePool.Put(c)
}

func (c *CaptureEnPassantMove) copy() Move {
    captureEnPassantMove := captureEnPassantMovePool.Get().(*CaptureEnPassantMove)
    captureEnPassantMove.b = c.b
    captureEnPassantMove.fromLocation = c.fromLocation
    captureEnPassantMove.toLocation = c.toLocation
    captureEnPassantMove.piece = c.piece
    captureEnPassantMove.newPiece = c.newPiece
    captureEnPassantMove.capturedPiece = c.capturedPiece
    captureEnPassantMove.en = c.en
    captureEnPassantMove.encs = c.encs
    captureEnPassantMove.vulnerables = c.vulnerables

    return captureEnPassantMove
}

func (c *CaptureEnPassantMove) execute() error {
	ok := c.b.setPiece(c.fromLocation, nil)
	if !ok {
		return fmt.Errorf("no piece at fromLocation")
	}

	for _, enc := range c.encs {
        ok = c.b.setPiece(enc.enPassant.pieceLocation, nil)
		if !ok {
			return fmt.Errorf("no piece at enPassant.pieceLocation")
		}
	}

	ok = c.b.setPiece(c.toLocation, c.newPiece)
	if !ok {
		return fmt.Errorf("no piece at toLocation")
	}

    c.b.setVulnerables(c.piece.getColor(), []Point{})
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
        ok = c.b.setPiece(enc.enPassant.pieceLocation, enc.capturedPiece)
		if !ok {
			return fmt.Errorf("no piece at enPassant.pieceLocation")
		}
	}

    c.b.setVulnerables(c.piece.getColor(), c.vulnerables)
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
    vulnerables   []Point
    newVulnerables []Point
}

func (c *CastleMove) getNewPiece() Piece {
    return c.newKing
}

func (c *CastleMove) putInPool() {
    castleMovePool.Put(c)
}

func (c *CastleMove) copy() Move {
    castleMove := castleMovePool.Get().(*CastleMove)
    castleMove.b = c.b
    castleMove.fromLocation = c.fromLocation
    castleMove.toLocation = c.toLocation
    castleMove.king = c.king
    castleMove.newKing = c.newKing
    castleMove.toKingLocation = c.toKingLocation
    castleMove.rook = c.rook
    castleMove.newRook = c.newRook
    castleMove.toRookLocation = c.toRookLocation
    castleMove.en = c.en
    castleMove.vulnerables = c.vulnerables
    castleMove.newVulnerables = c.newVulnerables

    return castleMove
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

    c.b.setVulnerables(c.king.getColor(), c.newVulnerables)
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

    c.b.setVulnerables(c.king.getColor(), c.vulnerables)
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

func (p *PromotionMove) putInPool() {
    p.baseMove.putInPool()
    promotionMovePool.Put(p)
}

func (p *PromotionMove) copy() Move {
    promotionMove := promotionMovePool.Get().(*PromotionMove)
    promotionMove.b = p.b
    promotionMove.fromLocation = p.fromLocation
    promotionMove.toLocation = p.toLocation
    promotionMove.baseMove = p.baseMove.copy()
    promotionMove.promotionPiece = p.promotionPiece

    return promotionMove
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

func (m *AllyDefenseMove) putInPool() {
    allyDefenseMovePool.Put(m)
}

func (m *AllyDefenseMove) copy() Move {
    allyDefenseMove := allyDefenseMovePool.Get().(*AllyDefenseMove)
    allyDefenseMove.b = m.b
    allyDefenseMove.fromLocation = m.fromLocation
    allyDefenseMove.toLocation = m.toLocation
    allyDefenseMove.piece = m.piece

    return allyDefenseMove
}

func (m *AllyDefenseMove) execute() error {
	return nil
}

func (m *AllyDefenseMove) undo() error {
	return nil
}

