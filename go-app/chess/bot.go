package chess

/*
Responsible for:
- keeping track of the state of the game the bot is playing
*/
type Bot interface {
    FindMove() (*MoveKey, error)
}

func NewSimpleBot(game Game) (Bot, error) {
    return &SimpleBot{
        game: game,
    }, nil
}

type SimpleBot struct {
    game Game
}

func (b *SimpleBot) FindMove() (*MoveKey, error) {
    // we should probably clone the game here before doing the engine stuff

    moveKey := &MoveKey{
        XFrom: 4,
        YFrom: 1,
        XTo: 4,
        YTo: 3,
        Promotion: "",
    } 

    return moveKey, nil
}

