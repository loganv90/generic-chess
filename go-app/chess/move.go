package chess

import "github.com/stretchr/testify/mock"

type action struct {
	b     board
	xFrom int
	yFrom int
	xTo   int
	yTo   int
}

type mockMove struct {
	mock.Mock
}

func (m *mockMove) execute() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockMove) undo() error {
	args := m.Called()
	return args.Error(0)
}

var moveFactoryInstance = moveFactory(&concreteMoveFactory{})

type moveFactory interface {
	newSimpleMove(b board, xFrom int, yFrom int, xTo int, yTo int) (*simpleMove, error)
}

type concreteMoveFactory struct{}

func (f *concreteMoveFactory) newSimpleMove(b board, xFrom int, yFrom int, xTo int, yTo int) (*simpleMove, error) {
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

type move interface {
	execute() error
	undo() error
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
