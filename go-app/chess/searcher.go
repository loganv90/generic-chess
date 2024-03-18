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
    b *SimpleBoard
    p *SimplePlayerCollection
    e *SimpleEvaluator
    transpositionMap map[string]MoveKeyAndScore
    minimaxCalls int
}

func (s *SimpleSearcher) search() (MoveKey, error) {
    // we should start by implementing adapting minimax to an arbitrary number of players
    // we can just do recursive search and pass around a single game object while execuing and undoing moves
    // first we need to make a copy of the game object

    s.b.CalculateMoves()

    _, move, ok, err := s.minimax(4)
    if err != nil {
        return MoveKey{}, err
    }
    
    if !ok {
        return MoveKey{}, fmt.Errorf("no move found")
    }

    promotionString := ""
    if move.promotionIndex >= 0 {
        promotionString = piece_names[move.promotionIndex]
    }

    return MoveKey{
        XFrom: move.fromLocation.x,
        YFrom: move.fromLocation.y,
        XTo: move.toLocation.x,
        YTo: move.toLocation.y,
        Promotion: promotionString,
    }, nil
}

func (s *SimpleSearcher) minimax(depth int) ([]int, FastMove, bool, error) {
    s.minimaxCalls++

    gameOver := s.p.getGameOver()
    if depth == 0 || gameOver {
        score, err := s.e.eval()
        if err != nil {
            panic(err)
        }

        return score, FastMove{}, false, nil
    }

    players := s.p.getPlayers()
    currentPlayer := s.p.getCurrent()
    if currentPlayer < 0 || currentPlayer >= players {
        panic(fmt.Errorf("invalid player"))
    }

    movesPointer := s.b.MovesOfColor(currentPlayer)
    if movesPointer == nil {
        panic(fmt.Errorf("invalid player"))
    }
    moves := *movesPointer

    inCheck := s.b.Check(currentPlayer)

    found := false
    var bestMove FastMove
    var transition PlayerTransition
    bestScore := make([]int, players)
    bestScore[currentPlayer] = -1000000

    for i := 0; i < moves.count; i++ {
        move := moves.array[i]
        if move.allyDefense {
            continue
        }

        move.execute()

        s.b.CalculateMoves()

        if s.b.Check(currentPlayer) {
            move.undo()
            continue
        }

        createPlayerTransition(s.b, s.p, false, false, &transition)

        transition.execute()

        score, _, _, err := s.minimax(depth-1)
        if err != nil {
            panic(err)
        }

        move.undo()

        transition.undo()

        if score[currentPlayer] > bestScore[currentPlayer] {
            bestScore = score
            bestMove = move
            found = true
        }
    }

    if !found {
        // stalemate
        if !inCheck {
            return make([]int, players), FastMove{}, false, nil // TODO this is probably a problem
        }

        // checkmate
        createPlayerTransition(s.b, s.p, true, false, &transition)

        transition.execute()

        score, _, _, err := s.minimax(depth)
        if err != nil {
            panic(err)
        }

        transition.undo()

        return score, FastMove{}, false, nil
    }

    return bestScore, bestMove, true, nil
}

