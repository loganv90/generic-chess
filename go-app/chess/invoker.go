package chess

import "fmt"

var invokerFactoryInstance = InvokerFactory(&ConcreteInvokerFactory{})

type InvokerFactory interface {
	newSimpleInvoker() (*SimpleInvoker, error)
}

type ConcreteInvokerFactory struct{}

func (f *ConcreteInvokerFactory) newSimpleInvoker() (*SimpleInvoker, error) {
	return &SimpleInvoker{
		history: []Command{},
		index:   0,
	}, nil
}

type Invoker interface {
	execute(m Move, b Board, p PlayerCollection) error
	undo() error
	redo() error
    Copy() (Invoker, error)
}

type SimpleInvoker struct {
	history []Command
	index   int
}

func (s *SimpleInvoker) execute(m Move, b Board, p PlayerCollection) error {
	err := m.execute()
	if err != nil {
		return err
	}

    err = b.CalculateMoves()
    if err != nil {
        return err
    }

    t, err := playerTransitionFactoryInstance.newIncrementalTransitionAsPlayerTransition(b, p)
    if err != nil {
        return err
    }

    err = t.execute()
    if err != nil {
        return err
    }

	s.history = append(s.history[:s.index], Command{m, t, b})
	s.index++

	return nil
}

func (s *SimpleInvoker) undo() error {
	if s.index <= 0 {
		return fmt.Errorf("no moves to undo")
	}

	err := s.history[s.index-1].m.undo()
	if err != nil {
		return err
	}

    err = s.history[s.index-1].b.CalculateMoves()
    if err != nil {
        return err
    }

	err = s.history[s.index-1].p.undo()
	if err != nil {
		return err
	}

	s.index--

	return nil
}

func (s *SimpleInvoker) redo() error {
	if s.index >= len(s.history) {
		return fmt.Errorf("no moves to redo")
	}

	err := s.history[s.index].m.execute()
	if err != nil {
		return err
	}

    err = s.history[s.index].b.CalculateMoves()
    if err != nil {
        return err
    }

	err = s.history[s.index].p.execute()
	if err != nil {
		return err
	}

	s.index++

	return nil
}

func (s *SimpleInvoker) Copy() (Invoker, error) {
	return &SimpleInvoker{
		history: []Command{},
		index:   0,
	}, nil
}

