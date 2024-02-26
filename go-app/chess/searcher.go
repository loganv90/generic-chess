package chess

/*
Chess Search:
Alpha beta search
Iterative deepening
Transposition table (repeated and colour flipped positions)
Keep track of killer moves (moves that caused cut off) either by [piece][to] or [from][to]
Quiescence search (do all captures on last capture/promotion square)
Move ordering (Hashed moves, winning captures/promotions, equal captures, killer moves, non-captures, losing captures)
Create a hash function to uniquely identify the chess position
q2k2q1/2nqn2b/1n1P1n1b/2rnr2Q/1NQ1QN1Q/3Q3B/2RQR2B/Q2K2Q1 w - - this position causes engines to explode
*/

/*
Responsible for:
- searching for moves given the current state of the game
*/
type Searcher interface {
    search() (MoveKey, error)
}

func newSimpleSearcher(g Game) (*SimpleSearcher, error) {
    return &SimpleSearcher{
        g: g,
        transpositionMap: map[string]MoveKeyAndScore{},
        minimaxCalls: 0,
    }, nil
}

type SimpleSearcher struct {
    g Game
    transpositionMap map[string]MoveKeyAndScore
    minimaxCalls int
}

func (s *SimpleSearcher) search() (MoveKey, error) {
    // we should start by implementing adapting minimax to an arbitrary number of players
    // we can just do recursive search and pass around a single game object while execuing and undoing moves
    // first we need to make a copy of the game object

    _, moveKey, err := s.minimax(s.g, 4)
    if err != nil {
        return MoveKey{}, err
    }

    return moveKey, nil
}

func (s *SimpleSearcher) minimax(game Game, depth int) (map[string]int, MoveKey, error) {
    s.minimaxCalls++

    state, err := game.State()
    if err != nil {
        return nil, MoveKey{}, err
    }

    board := game.getBoard()
    playerCollection := game.getPlayerCollection()

    uniqueString := board.UniqueString()
    if moveKeyAndScore, ok := s.transpositionMap[uniqueString]; ok {
        return moveKeyAndScore.score, moveKeyAndScore.moveKey, nil
    }

    if depth == 0 || state.GameOver {
        evaluator, err := newSimpleEvaluator(board, playerCollection)
        if err != nil {
            return nil, MoveKey{}, err
        }

        score, err := evaluator.eval()
        if err != nil {
            return nil, MoveKey{}, err
        }

        return score, MoveKey{}, nil
    }

    moves, err := game.Moves(state.CurrentPlayer)
    if err != nil {
        return nil, MoveKey{}, err
    }

    bestScore := map[string]int{state.CurrentPlayer: -1000000}
    bestMove := MoveKey{}
    for _, move := range moves {
        err := game.Execute(move.XFrom, move.YFrom, move.XTo, move.YTo, move.Promotion)
        if err != nil {
            return nil, MoveKey{}, err
        }

        score, moveKey, err := s.minimax(game, depth-1)
        if err != nil {
            return nil, MoveKey{}, err
        }

        s.transpositionMap[uniqueString] = MoveKeyAndScore{
            moveKey: moveKey,
            score: score,
        }

        if score[state.CurrentPlayer] > bestScore[state.CurrentPlayer] {
            bestScore = score
            bestMove = move
        }

        err = game.Undo()
        if err != nil {
            return nil, MoveKey{}, err
        }
    }

    return bestScore, bestMove, nil
}

func (s *SimpleSearcher) minimax2(b Board, p PlayerCollection, depth int) (map[string]int, Move, error) {
    if depth == 0 {
        evaluator, err := newSimpleEvaluator(b, p)
        if err != nil {
            return nil, nil, err
        }

        score, err := evaluator.eval()
        if err != nil {
            return nil, nil, err
        }

        return score, nil, nil
    }

    currentPlayer, err := p.getCurrent()
    if err != nil {
        return nil, nil, err
    }

    moves, err := b.AvailableMoves(currentPlayer)
    if err != nil {
        return nil, nil, err
    }

    bestScore := map[string]int{currentPlayer: -1000000}
    var bestMove Move
    playersEliminated := []string{}
    canMove := false

    for _, move := range moves {
        err := move.execute()
        if err != nil {
            return nil, nil, err
        }

        err = b.CalculateMoves()
        if err != nil {
            return nil, nil, err
        }

        if b.Check(currentPlayer) {
            err := move.undo()
            if err != nil {
                return nil, nil, err
            }
            continue
        }

        canMove = true

        score, _, err := s.minimax2(b, p, depth-1)
        if err != nil {
            return nil, nil, err
        }

        if score[currentPlayer] > bestScore[currentPlayer] {
            bestScore = score
            bestMove = move
        }

        err = move.undo()
        if err != nil {
            return nil, nil, err
        }
    }

    if !canMove {
        playersEliminated = append(playersEliminated, currentPlayer)
    }

    return bestScore, bestMove, nil
}

