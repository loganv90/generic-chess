package chess

/*
Responsible for:
- keeping track of the state of the game the bot is playing
*/
type Bot interface {
    FindMove() (*MoveKey, error)
}

func NewSimpleBot(game Game) (Bot, error) {
    return &SimpleBot{}, nil
}

type SimpleBot struct {
}

func (b *SimpleBot) FindMove() (*MoveKey, error) {
    moveKey := &MoveKey{
        XFrom: 4,
        YFrom: 1,
        XTo: 4,
        YTo: 3,
        Promotion: "",
    } 

    return moveKey, nil
}

