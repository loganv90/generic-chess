package chess

import (
	"testing"
    "strings"

	"github.com/stretchr/testify/assert"
)

func Test_getNextAndRemaining(t *testing.T) {
    white := 0
    black := 1
    blue := 2
    red := 3

    // all alive
    s, err := newSimplePlayerCollection(4)
    assert.Nil(t, err)

    playerColor := s.getCurrent()
    assert.Equal(t, white, playerColor)

    next, remaining := s.getNextAndRemaining()
    assert.Equal(t, black, next)
    assert.Equal(t, 4, remaining)

    // 0 and 2 alive
    s, err = newSimplePlayerCollection(4)
    assert.Nil(t, err)
    s.eliminate(black)
    s.eliminate(red)

    playerColor = s.getCurrent()
    assert.Equal(t, white, playerColor)

    next, remaining = s.getNextAndRemaining()
    assert.Equal(t, blue, next)
    assert.Equal(t, 2, remaining)

    // none alive
    s, err = newSimplePlayerCollection(4)
    assert.Nil(t, err)
    s.eliminate(white)
    s.eliminate(black)
    s.eliminate(blue)
    s.eliminate(red)

    playerColor = s.getCurrent()
    assert.Equal(t, white, playerColor)

    next, remaining = s.getNextAndRemaining()
    assert.Equal(t, 0, next)
    assert.Equal(t, 0, remaining)

    // 3 alive
    s, err = newSimplePlayerCollection(4)
    assert.Nil(t, err)
    s.eliminate(white)
    s.eliminate(black)
    s.eliminate(blue)
        
    playerColor = s.getCurrent()
    assert.Equal(t, white, playerColor)

    next, remaining = s.getNextAndRemaining()
    assert.Equal(t, red, next)
    assert.Equal(t, 1, remaining)
}

func Test_PlayerCollectionUniqueString_Default(t *testing.T) {
    builder := strings.Builder{}

    p, err := createSimplePlayerCollectionWithDefaultPlayers()
    assert.Nil(t, err)

    p.UniqueString(&builder)
    expected := "01-0"
    assert.Equal(t, expected, builder.String())
}

