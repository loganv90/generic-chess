package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockPlayerCollection struct {
	mock.Mock
}

func (m *MockPlayerCollection) getNext() (*Player, error) {
    args := m.Called()

    if args.Get(0) == nil {
        return nil, args.Error(1)
    } else {
        return args.Get(0).(*Player), args.Error(1)
    }
}

func (m *MockPlayerCollection) getPlayerColors() []string {
    args := m.Called()
    return args.Get(0).([]string)
}

func (m *MockPlayerCollection) eliminate(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockPlayerCollection) restore(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockPlayerCollection) getCurrent() (string, error) {
    args := m.Called()
    return args.String(0), args.Error(1)
}

func (m *MockPlayerCollection) setCurrent(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockPlayerCollection) getWinner() (string, error) {
    args := m.Called()
    return args.String(0), args.Error(1)
}

func (m *MockPlayerCollection) setWinner(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockPlayerCollection) getGameOver() (bool, error) {
    args := m.Called()
    return args.Bool(0), args.Error(1)
}

func (m *MockPlayerCollection) setGameOver(gameOver bool) error {
    args := m.Called(gameOver)
    return args.Error(0)
}

func (m *MockPlayerCollection) Copy() (PlayerCollection, error) {
    args := m.Called()
    return args.Get(0).(PlayerCollection), args.Error(1)
}

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

