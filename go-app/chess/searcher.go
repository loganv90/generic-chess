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
    search(color string) (*MoveKey, error)
}

func newSimpleSearcher(g Game) (*SimpleSearcher, error) {
    return &SimpleSearcher{
        g: g,
    }, nil
}

type SimpleSearcher struct {
    g Game
}

func (s *SimpleSearcher) search(color string) (*MoveKey, error) {
    _, err := s.g.Copy()
    if err != nil {
        return nil, err
    }

    // we should start by implementing adapting minimax to an arbitrary number of players
    // we can just do recursive search and pass around a single game object while execuing and undoing moves
    // first we need to make a copy of the game object
    return nil, nil
}

