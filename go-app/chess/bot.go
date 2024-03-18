package chess

/*
Responsible for:
- keeping track of the state of the game the bot is playing
*/
type Bot interface {
    FindMove() (MoveKey, error)
}

func NewSimpleBot(game Game) (Bot, error) {
    return &SimpleBot{
        game: game,
    }, nil
}

type SimpleBot struct {
    game Game
}

func (b *SimpleBot) FindMove() (MoveKey, error) {
    gameCopy, err := b.game.Copy()
    if err != nil {
        return MoveKey{}, err
    }

    searcher := newSimpleSearcher(gameCopy)

    moveKey, err := searcher.search()

    return moveKey, nil
}

