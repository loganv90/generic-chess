package chess

import "fmt"

var invokerFactoryInstance = invokerFactory(&concreteInvokerFactory{})

type invokerFactory interface {
	newSimpleInvoker() (*simpleInvoker, error)
}

type concreteInvokerFactory struct{}

func (f *concreteInvokerFactory) newSimpleInvoker() (*simpleInvoker, error) {
	return &simpleInvoker{
		history: []move{},
		index:   0,
	}, nil
}

type invoker interface {
	execute(m move) error
	undo() error
	redo() error
}

type simpleInvoker struct {
	history []move
	index   int
}

func (s *simpleInvoker) execute(m move) error {
	err := m.execute()
	if err != nil {
		return err
	}

	s.history = append(s.history[:s.index], m)
	s.index++

	return nil
}

func (s *simpleInvoker) undo() error {
	if s.index <= 0 {
		return fmt.Errorf("no moves to undo")
	}

	err := s.history[s.index-1].undo()
	if err != nil {
		return err
	}

	s.index--

	return nil
}

func (s *simpleInvoker) redo() error {
	if s.index >= len(s.history) {
		return fmt.Errorf("no moves to redo")
	}

	err := s.history[s.index].execute()
	if err != nil {
		return err
	}

	s.index++

	return nil
}
