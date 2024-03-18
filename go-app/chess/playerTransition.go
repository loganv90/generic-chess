package chess

func createPlayerTransition(b *SimpleBoard, p *SimplePlayerCollection, inCheckmate bool, inStalemate bool, t *PlayerTransition) {
    oldCurrent := p.getCurrent()
    oldWinner := p.getWinner()
    oldGameOver := p.getGameOver()
    next, remaining := p.getNextAndRemaining()

    var newCurrent int
    var newWinner int
    var newGameOver bool

    if inStalemate {
        newCurrent = oldCurrent
        newWinner = -1
        newGameOver = true
    } else if remaining < 1 {
        newCurrent = oldCurrent
        newWinner = oldWinner
        newGameOver = true
    } else if remaining == 1 {
        newCurrent = next
        newWinner = next
        newGameOver = true
    } else if remaining == 2 {
        if inCheckmate {
            newCurrent = next
            newWinner = next
            newGameOver = true
        } else {
            newCurrent = next
            newWinner = -1
            newGameOver = false
        }
    } else {
        newCurrent = next
        newWinner = -1
        newGameOver = false
    }

    t.p = p
    t.b = b
    t.oldCurrent = oldCurrent
    t.newCurrent = newCurrent
    t.oldWinner = oldWinner
    t.newWinner = newWinner
    t.oldGameOver = oldGameOver
    t.newGameOver = newGameOver
    t.eliminated = inCheckmate
}

type PlayerTransition struct {
    p *SimplePlayerCollection
    b *SimpleBoard
    oldCurrent int
    newCurrent int
    oldWinner int
    newWinner int
    oldGameOver bool
    newGameOver bool
    eliminated bool
}

func (s *PlayerTransition) execute() {
    s.p.setCurrent(s.newCurrent)
    s.p.setWinner(s.newWinner)
    s.p.setGameOver(s.newGameOver)

    if !s.eliminated {
        return
    }

    s.p.eliminate(s.oldCurrent)
    s.b.disablePieces(s.oldCurrent, true)
}

func (s *PlayerTransition) undo() {
    s.p.setCurrent(s.oldCurrent)
    s.p.setWinner(s.oldWinner)
    s.p.setGameOver(s.oldGameOver)

    if !s.eliminated {
        return
    }

    s.p.restore(s.oldCurrent)
    s.b.disablePieces(s.oldCurrent, false)
}

