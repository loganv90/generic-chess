package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_Bot(t *testing.T) {
    game, err := NewSimpleSmallGame()
    assert.Nil(t, err)

    bot, err := NewSimpleBot(game, 2, 5)
    assert.Nil(t, err)

    moveKey, err := bot.FindMoveIterativeDeepening()
    assert.Nil(t, err)
    assert.Equal(t, MoveKey{3, 3, 3, 2, ""}, moveKey)

    bot, err = NewSimpleBot(game, 5, 5)
    assert.Nil(t, err)

    moveKey, err = bot.FindMoveIterativeDeepening()
    assert.Nil(t, err)
    assert.Equal(t, MoveKey{2, 3, 2, 2, ""}, moveKey)
}

