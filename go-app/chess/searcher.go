package chess

import (
    "fmt"
)

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

    err := s.b.CalculateMoves()
    if err != nil {
        return MoveKey{}, err
    }

    _, move, ok, err := s.minimax(4)
    if err != nil {
        return MoveKey{}, err
    }
    
    if !ok {
        return MoveKey{}, fmt.Errorf("no move found")
    }

    return MoveKey{
        XFrom: move.fromLocation.x,
        YFrom: move.fromLocation.y,
        XTo: move.toLocation.x,
        YTo: move.toLocation.y,
        Promotion: move.promotion,
    }, nil
}

func (s *SimpleSearcher) minimax(depth int) ([]int, FastMove, bool, error) {
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

        return score, FastMove{}, false, nil
    }

    players := s.p.getPlayers()
    currentPlayer, _ := s.p.getCurrent()
    if currentPlayer < 0 || currentPlayer >= players {
        panic(fmt.Errorf("invalid player"))
    }

    movesPointer, err := s.b.MovesOfColor(currentPlayer)
    moves := *movesPointer
    if err != nil {
        panic(err)
    }

    inCheck := s.b.Check(currentPlayer)

    found := false
    var bestMove FastMove
    bestScore := make([]int, players)
    bestScore[currentPlayer] = -1000000

    for i := 0; i < moves.count; i++ {
        move := moves.array[i]
        if move.allyDefense {
            continue
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

        score, _, _, err := s.minimax(depth-1)
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
            found = true
        }
    }

    if !found {
        // stalemate
        if !inCheck {
            return make([]int, players), FastMove{}, false, nil
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

        score, _, _, err := s.minimax(depth)
        if err != nil {
            panic(err)
        }

        err = transition.undo()
        if err != nil {
            panic(err)
        }

        return score, FastMove{}, false, nil
    }

    return bestScore, bestMove, true, nil
}

