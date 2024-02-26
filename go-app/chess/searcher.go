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

    b := s.g.getBoard()
    p := s.g.getPlayerCollection()

    _, move, err := s.minimax(b, p, 3)
    if err != nil {
        return MoveKey{}, err
    }

    action := move.getAction()

    return MoveKey{
        XFrom: action.fromLocation.x,
        YFrom: action.fromLocation.y,
        XTo: action.toLocation.x,
        YTo: action.toLocation.y,
        Promotion: "Q",
    }, nil
}

func (s *SimpleSearcher) minimax(b Board, p PlayerCollection, depth int) (map[string]int, Move, error) {
    gameOver, err := p.getGameOver()
    if err != nil {
        return nil, nil, err
    }

    if depth == 0 || gameOver {
        // TODO do not make a new evaluator every time
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

    moves, err := b.MovesOfColor(currentPlayer)
    if err != nil {
        return nil, nil, err
    }

    var bestMove Move
    bestScore := map[string]int{currentPlayer: -1000000}

    for _, move := range moves {
        if _, ok := move.(*AllyDefenseMove); ok {
            continue
        }

        if promotionMove, ok := move.(*PromotionMove); ok {
            // TODO do not make a new queen every time
            promotionMove.setPromotionPiece(newQueen(currentPlayer))
        }

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

        transition, err := p.GetTransition(b, false, false)
        if err != nil {
            return nil, nil, err
        }

        err = transition.execute()
        if err != nil {
            return nil, nil, err
        }

        score, _, err := s.minimax(b, p, depth-1)
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

        err = transition.undo()
        if err != nil {
            return nil, nil, err
        }
    }

    if bestMove == nil {
        // stalemate
        if !b.Check(currentPlayer) {
            return map[string]int{}, nil, nil
        }

        // checkmate
        transition, err := p.GetTransition(b, true, false)
        if err != nil {
            return nil, nil, err
        }

        err = transition.execute()
        if err != nil {
            return nil, nil, err
        }

        score, _, err := s.minimax(b, p, depth-1)
        if err != nil {
            return nil, nil, err
        }

        err = transition.undo()
        if err != nil {
            return nil, nil, err
        }

        return score, nil, nil
    }

    return bestScore, bestMove, nil
}

