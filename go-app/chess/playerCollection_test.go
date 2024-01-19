package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_incrementAndDecrement(t *testing.T) {
    s, err := newSimplePlayerCollection(
        []*Player{
            {"white", true},
            {"black", true},
            {"blue", true},
            {"red", true},
        },
    )
    assert.Nil(t, err)
    assert.Equal(t, 0, s.currentPlayer)
    s.increment()
    assert.Equal(t, 1, s.currentPlayer)
    s.decrement()
    assert.Equal(t, 0, s.currentPlayer)

    s, err = newSimplePlayerCollection(
        []*Player{
            {"white", true},
            {"black", false},
            {"blue", true},
            {"red", false},
        },
    )
    assert.Nil(t, err)
    assert.Equal(t, 0, s.currentPlayer)
    s.increment()
    assert.Equal(t, 2, s.currentPlayer)
    s.decrement()
    assert.Equal(t, 0, s.currentPlayer)

    s, err = newSimplePlayerCollection(
        []*Player{
            {"white", false},
            {"black", false},
            {"blue", false},
            {"red", false},
        },
    )
    assert.Nil(t, err)
    assert.Equal(t, 0, s.currentPlayer)
    s.increment()
    assert.Equal(t, 0, s.currentPlayer)
    s.decrement()
    assert.Equal(t, 0, s.currentPlayer)

    s, err = newSimplePlayerCollection(
        []*Player{
            {"white", false},
            {"black", false},
            {"blue", false},
            {"red", true},
        },
    )
    assert.Nil(t, err)
    assert.Equal(t, 0, s.currentPlayer)
    s.increment()
    assert.Equal(t, 3, s.currentPlayer)
    s.decrement()
    assert.Equal(t, 3, s.currentPlayer)
}

