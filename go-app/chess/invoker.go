package chess

import "fmt"

var invokerFactoryInstance = InvokerFactory(&ConcreteInvokerFactory{})

type InvokerFactory interface {
	newSimpleInvoker() (*SimpleInvoker, error)
}

type ConcreteInvokerFactory struct{}

func (f *ConcreteInvokerFactory) newSimpleInvoker() (*SimpleInvoker, error) {
	return &SimpleInvoker{
		history: []MoveAndPlayerTransition{},
		index:   0,
	}, nil
}

type Invoker interface {
	execute(m Move, p PlayerTransition) error
	undo() error
	redo() error
}

type SimpleInvoker struct {
	history []MoveAndPlayerTransition
	index   int
}

func (s *SimpleInvoker) execute(m Move, p PlayerTransition) error {
	err := m.execute()
	if err != nil {
		return err
	}

    err = p.execute()
    if err != nil {
        return err
    }

	s.history = append(s.history[:s.index], MoveAndPlayerTransition{m, p})
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

	err = s.history[s.index].p.execute()
	if err != nil {
		return err
	}

	s.index++

	return nil
}
