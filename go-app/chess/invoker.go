package chess

import "fmt"

var invokerFactoryInstance = InvokerFactory(&ConcreteInvokerFactory{})

type InvokerFactory interface {
	newSimpleInvoker() (*SimpleInvoker, error)
}

type ConcreteInvokerFactory struct{}

func (f *ConcreteInvokerFactory) newSimpleInvoker() (*SimpleInvoker, error) {
	return &SimpleInvoker{
		history: []Move{},
		index:   0,
	}, nil
}

type Invoker interface {
	execute(m Move) error
	undo() error
	redo() error
}

type SimpleInvoker struct {
	history []Move
	index   int
}

func (s *SimpleInvoker) execute(m Move) error {
	err := m.execute()
	if err != nil {
		return err
	}

	s.history = append(s.history[:s.index], m)
	s.index++

	return nil
}

func (s *SimpleInvoker) undo() error {
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

func (s *SimpleInvoker) redo() error {
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
