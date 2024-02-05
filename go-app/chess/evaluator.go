package chess

/*
Chess Evaluation:
Material count
Piece mobility with move bonus for different pieces (including moves protecting ally pieces, and not including moves to controlled squares) low queen bonus
Piece locations piece-square tables (dp for knight 2 moves potential squares) tables for opening and endgame, then interpolate
Penalty isolated pawns worth less than chained pawns
Penalty for attacked squares close to king
Penalty for lots of mobility if king were a queen
Bonus for attacking close to own king
Bonus for pinning pieces to more valuable pieces
Bonus for queen-rook, queen-bishop, bishop-bishop, rook-rook combos
*/

/*
Responsible for:
- evaluating a board and returning a score
*/
type Evaluator interface {
    eval() int
}

func newSimpleEvaluator(b Board, p PlayerCollection) (*SimpleEvaluator, error) {
    return &SimpleEvaluator{
        b: b,
        p: p,
    }, nil
}

type SimpleEvaluator struct {
    b Board
    p PlayerCollection
}

func (e *SimpleEvaluator) eval() (int, error) {
    gameOver, err := e.p.getGameOver()
    if err != nil {
        return 0, err
    }

    winner, err := e.p.getWinner()
    if err != nil {
        return 0, err
    }

    current, err := e.p.getCurrent()
    if err != nil {
        return 0, err
    }

    if gameOver {
        if winner == "" {
            return 0, nil
        } else if winner == current {
            return 100000, nil
        } else {
            return -100000, nil
        }
    }

    // Material comparison
    // we need: the locations of each piece by player

    // Piece position comparison (piece-square tables)
    // we need: the locations of each piece by player

    // Mobility comparison
    // we need: the moves each piece can make including attacking ally pieces

    return 1, nil
}

