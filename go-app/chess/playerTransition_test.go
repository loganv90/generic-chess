package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_IncrementalTransition_inCheckmate(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)

    var transition PlayerTransition
    createPlayerTransition(b, p, true, false, &transition)
    assert.Equal(t, transition.eliminated, true)

    transition.execute()
    assert.Equal(t, p.getWinner(), black)
    assert.Equal(t, p.getGameOver(), true)
    assert.Equal(t, p.getCurrent(), black)

    transition.undo()
    assert.Equal(t, p.getWinner(), -1)
    assert.Equal(t, p.getGameOver(), false)
    assert.Equal(t, p.getCurrent(), white)
}

func Test_IncrementalTransition_noCheckmate(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)

    var transition PlayerTransition
    createPlayerTransition(b, p, false, false, &transition)
    assert.Equal(t, transition.eliminated, false)

    transition.execute()
    assert.Equal(t, p.getWinner(), -1)
    assert.Equal(t, p.getGameOver(), false)
    assert.Equal(t, p.getCurrent(), black)

    transition.undo()
    assert.Equal(t, p.getWinner(), -1)
    assert.Equal(t, p.getGameOver(), false)
    assert.Equal(t, p.getCurrent(), white)
}

