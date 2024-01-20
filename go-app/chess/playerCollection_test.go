package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getNext(t *testing.T) {
    s, err := newSimplePlayerCollection(
        []*Player{
            {"white", true},
            {"black", true},
            {"blue", true},
            {"red", true},
        },
    )
    assert.Nil(t, err)
    playerColor, err := s.getCurrent()
    assert.Equal(t, "white", playerColor)
    player, err := s.getNext()
    assert.Nil(t, err)
    assert.Equal(t, "black", player.color)

    s, err = newSimplePlayerCollection(
        []*Player{
            {"white", true},
            {"black", false},
            {"blue", true},
            {"red", false},
        },
    )
    assert.Nil(t, err)
    playerColor, err = s.getCurrent()
    assert.Equal(t, "white", playerColor)
    player, err = s.getNext()
    assert.Nil(t, err)
    assert.Equal(t, "blue", player.color)

    s, err = newSimplePlayerCollection(
        []*Player{
            {"white", false},
            {"black", false},
            {"blue", false},
            {"red", false},
        },
    )
    assert.Nil(t, err)
    playerColor, err = s.getCurrent()
    assert.Equal(t, "white", playerColor)
    player, err = s.getNext()
    assert.Nil(t, err)
    assert.Equal(t, "white", player.color)

    s, err = newSimplePlayerCollection(
        []*Player{
            {"white", false},
            {"black", false},
            {"blue", false},
            {"red", true},
        },
    )
    assert.Nil(t, err)
    playerColor, err = s.getCurrent()
    assert.Equal(t, "white", playerColor)
    player, err = s.getNext()
    assert.Nil(t, err)
    assert.Equal(t, "red", player.color)
}

