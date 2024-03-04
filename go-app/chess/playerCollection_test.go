package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockPlayerCollection struct {
	mock.Mock
}

func (m *MockPlayerCollection) eliminate(color int) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockPlayerCollection) restore(color int) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockPlayerCollection) getCurrent() (int, bool) {
    args := m.Called()
    return args.Int(0), args.Bool(1)
}

func (m *MockPlayerCollection) setCurrent(color int) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockPlayerCollection) getWinner() (int, bool) {
    args := m.Called()
    return args.Int(0), args.Bool(1)
}

func (m *MockPlayerCollection) setWinner(color int) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockPlayerCollection) getGameOver() (bool, error) {
    args := m.Called()
    return args.Bool(0), args.Error(1)
}

func (m *MockPlayerCollection) setGameOver(gameOver bool) error {
    args := m.Called(gameOver)
    return args.Error(0)
}

func (m *MockPlayerCollection) getNextAndRemaining() (int, int, error) {
    args := m.Called()
    return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockPlayerCollection) getPlayers() int {
    args := m.Called()
    return args.Int(0)
}

func (m *MockPlayerCollection) Copy() (PlayerCollection, error) {
    args := m.Called()
    return args.Get(0).(PlayerCollection), args.Error(1)
}

func (m *MockPlayerCollection) GetTransition(b Board, inCheckmate bool, inStalemate bool) (PlayerTransition, error) {
    args := m.Called(b, inCheckmate, inStalemate)
    return args.Get(0).(PlayerTransition), args.Error(1)
}

func Test_getNextAndRemaining(t *testing.T) {
    white := 0
    black := 1
    blue := 2
    red := 3

    // all alive
    s, err := newSimplePlayerCollection(4)
    assert.Nil(t, err)

    playerColor, ok := s.getCurrent()
    assert.Equal(t, true, ok)
    assert.Equal(t, white, playerColor)

    next, remaining, err := s.getNextAndRemaining()
    assert.Nil(t, err)
    assert.Equal(t, black, next)
    assert.Equal(t, 4, remaining)

    // 0 and 2 alive
    s, err = newSimplePlayerCollection(4)
    assert.Nil(t, err)
    s.eliminate(black)
    s.eliminate(red)

    playerColor, ok = s.getCurrent()
    assert.Equal(t, true, ok)
    assert.Equal(t, white, playerColor)

    next, remaining, err = s.getNextAndRemaining()
    assert.Nil(t, err)
    assert.Equal(t, blue, next)
    assert.Equal(t, 2, remaining)

    // none alive
    s, err = newSimplePlayerCollection(4)
    assert.Nil(t, err)
    s.eliminate(white)
    s.eliminate(black)
    s.eliminate(blue)
    s.eliminate(red)

    playerColor, ok = s.getCurrent()
    assert.Equal(t, true, ok)
    assert.Equal(t, white, playerColor)

    next, remaining, err = s.getNextAndRemaining()
    assert.Nil(t, err)
    assert.Equal(t, 0, next)
    assert.Equal(t, 0, remaining)

    // 3 alive
    s, err = newSimplePlayerCollection(4)
    assert.Nil(t, err)
    s.eliminate(white)
    s.eliminate(black)
    s.eliminate(blue)
        
    playerColor, ok = s.getCurrent()
    assert.Equal(t, true, ok)
    assert.Equal(t, white, playerColor)

    next, remaining, err = s.getNextAndRemaining()
    assert.Nil(t, err)
    assert.Equal(t, red, next)
    assert.Equal(t, 1, remaining)
}

