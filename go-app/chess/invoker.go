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
	execute(m Move, p PlayerTransition) error
	undo() error
	redo() error
    Copy() (Invoker, error)
}

type SimpleInvoker struct {
	history []Command
	index int
}

func (s *SimpleInvoker) execute(m Move, p PlayerTransition) error {
    fullMove := true

    if m != nil {
        err := m.execute()
        if err != nil {
            return err
        }
    } else {
        fullMove = false
    }

    if p != nil {
        err := p.execute()
        if err != nil {
            return err
        }
    } else {
        fullMove = false
    }

	s.history = append(s.history[:s.index+1], Command{m, p, fullMove})
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

    if command.m != nil {
        err := command.m.undo()
        if err != nil {
            return err
        }
    }

    if command.p != nil {
        err := command.p.undo()
        if err != nil {
            return err
        }
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

    if command.m != nil {
        err := command.m.execute()
        if err != nil {
            return err
        }
    }

    if command.p != nil {
        err := command.p.execute()
        if err != nil {
            return err
        }
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

