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
    e, err := newSimpleEvaluator(g.getBoard(), g.getPlayerCollection())
    if err != nil {
        return nil, err
    }

    return &SimpleSearcher{
        b: g.getBoard(),
        p: g.getPlayerCollection(),
        e: e,
        transpositionMap: map[string]MoveKeyAndScore{},
        minimaxCalls: 0,
    }, nil
}

type SimpleSearcher struct {
    b Board
    p PlayerCollection
    e Evaluator
    transpositionMap map[string]MoveKeyAndScore
    minimaxCalls int
}

func (s *SimpleSearcher) search() (MoveKey, error) {
    // we should start by implementing adapting minimax to an arbitrary number of players
    // we can just do recursive search and pass around a single game object while execuing and undoing moves
    // first we need to make a copy of the game object

    _, move, err := s.minimax(4)
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

func (s *SimpleSearcher) minimax(depth int) (map[string]int, Move, error) {
    s.minimaxCalls++

    gameOver, err := s.p.getGameOver()
    if err != nil {
        panic(err)
    }

    if depth == 0 || gameOver {
        score, err := s.e.eval()
        if err != nil {
            panic(err)
        }

        return score, nil, nil
    }

    currentPlayer, err := s.p.getCurrent()
    if err != nil {
        panic(err)
    }

    moves, err := s.b.MovesOfColor(currentPlayer)
    if err != nil {
        panic(err)
    }

    inCheck := s.b.Check(currentPlayer)

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
            panic(err)
        }

        err = s.b.CalculateMoves()
        if err != nil {
            panic(err)
        }

        if s.b.Check(currentPlayer) {
            err := move.undo()
            if err != nil {
                panic(err)
            }
            continue
        }

        transition, err := s.p.GetTransition(s.b, false, false)
        if err != nil {
            panic(err)
        }

        err = transition.execute()
        if err != nil {
            panic(err)
        }

        score, _, err := s.minimax(depth-1)
        if err != nil {
            panic(err)
        }

        err = move.undo()
        if err != nil {
            panic(err)
        }

        err = transition.undo()
        if err != nil {
            panic(err)
        }

        if score[currentPlayer] > bestScore[currentPlayer] {
            bestScore = score
            bestMove = move
        }
    }

    if bestMove == nil {
        // stalemate
        if !inCheck {
            return map[string]int{currentPlayer: 0}, nil, nil
        }

        // checkmate
        transition, err := s.p.GetTransition(s.b, true, false)
        if err != nil {
            panic(err)
        }

        err = transition.execute()
        if err != nil {
            panic(err)
        }

        score, _, err := s.minimax(depth)
        if err != nil {
            panic(err)
        }

        err = transition.undo()
        if err != nil {
            panic(err)
        }

        return score, nil, nil
    }

    return bestScore, bestMove, nil
}

