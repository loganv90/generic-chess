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
		index: -1,
	}, nil
}

type Invoker interface {
	execute(m FastMove, p PlayerTransition) error
    executeHalf(p PlayerTransition) error
	undo() error
	redo() error
    Copy() (Invoker, error)
}

type SimpleInvoker struct {
	history []Command
	index int
}

func (s *SimpleInvoker) execute(m FastMove, p PlayerTransition) error {
    err := m.execute()
    if err != nil {
        return err
    }

    err = p.execute()
    if err != nil {
        return err
    }

	s.history = append(s.history[:s.index+1], Command{m, p, true})
    s.index++

	return nil
}

func (s *SimpleInvoker) executeHalf(p PlayerTransition) error {
    err := p.execute()
    if err != nil {
        return err
    }

	s.history = append(s.history[:s.index+1], Command{FastMove{}, p, false})
    s.index++

	return nil
}

func (s *SimpleInvoker) undo() error {
	if s.index < 0 {
		return fmt.Errorf("no moves to undo")
	}
    commandToUndo := s.history[s.index]

    for !commandToUndo.fullMove {
        err := s.undoHelper()
        if err != nil {
            return err
        }

        if s.index < 0 {
            return fmt.Errorf("no moves to undo")
        }
        commandToUndo = s.history[s.index]
    }

    err := s.undoHelper()
    if err != nil {
        return err
    }

    return nil
}

func (s *SimpleInvoker) undoHelper() error {
	if s.index < 0 {
		return fmt.Errorf("no moves to undo")
	}

    command := s.history[s.index]

    if command.fullMove {
        err := command.m.undo()
        if err != nil {
            return err
        }
    }

    err := command.p.undo()
    if err != nil {
        return err
    }

	s.index--

	return nil
}

func (s *SimpleInvoker) redo() error {
    err := s.redoHelper()
    if err != nil {
        return err
    }

	if s.index+1 > len(s.history)-1 {
		return nil
	}
    commandToRedo := s.history[s.index+1]

    for !commandToRedo.fullMove {
        err := s.redoHelper()
        if err != nil {
            return err
        }

        if s.index+1 > len(s.history)-1 {
            return nil
        }
        commandToRedo = s.history[s.index+1]
    }

    return nil
}

func (s *SimpleInvoker) redoHelper() error {
	if s.index+1 > len(s.history)-1 {
		return fmt.Errorf("no moves to redo")
	}

    command := s.history[s.index+1]

    if command.fullMove {
        err := command.m.execute()
        if err != nil {
            return err
        }
    }

    err := command.p.execute()
    if err != nil {
        return err
    }

	s.index++

	return nil
}

func (s *SimpleInvoker) Copy() (Invoker, error) {
	return &SimpleInvoker{
		history: []Command{},
		index: -1,
	}, nil
}

