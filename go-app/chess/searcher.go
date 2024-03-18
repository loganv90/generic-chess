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
func newSimpleSearcher(g Game) *SimpleSearcher {
    board := g.getBoard()
    playerCollection := g.getPlayerCollection()

    return &SimpleSearcher{
        b: board,
        p: playerCollection,
        e: newSimpleEvaluator(board, playerCollection),

        players: playerCollection.getPlayers(),
        scores: [][]int{},
        transpositionMap: map[string]MoveKeyAndScore{},
        minimaxCalls: 0,
    }
}

type SimpleSearcher struct {
    b *SimpleBoard
    p *SimplePlayerCollection
    e *SimpleEvaluator

    players int
    scores [][]int
    transpositionMap map[string]MoveKeyAndScore
    minimaxCalls int
}

func (s *SimpleSearcher) search(depth int) (MoveKey, error) {
    s.scores = make([][]int, depth + 1)
    for i := 0; i < depth + 1; i++ {
        s.scores[i] = make([]int, s.players)
    }

    moveKey := MoveKey{}

    s.b.CalculateMoves()
    s.minimax(depth, &moveKey)

    return moveKey, nil
}

func (s *SimpleSearcher) minimax(depth int, moveKey *MoveKey) {
    s.minimaxCalls++

    gameOver := s.p.getGameOver()
    if depth == 0 || gameOver {
        s.e.eval(s.scores[depth])
        return
    }

    found := false
    transition := PlayerTransition{}
    currentPlayer := s.p.getCurrent()
    movesPointer := s.b.MovesOfColor(currentPlayer)
    moves := *movesPointer
    
    for i := 0; i < len(s.scores[depth]); i++ {
        s.scores[depth][i] = -1000000
    }

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

        s.minimax(depth-1, nil)

        transition.undo()

        move.undo()

        if s.scores[depth-1][currentPlayer] > s.scores[depth][currentPlayer] {
            for i := 0; i < len(s.scores[depth]); i++ {
                s.scores[depth][i] = s.scores[depth-1][i]
            }
            if moveKey != nil {
                moveKey.XTo = move.toLocation.x
                moveKey.YTo = move.toLocation.y
                moveKey.XFrom = move.fromLocation.x
                moveKey.YFrom = move.fromLocation.y
            }
            found = true
        }
    }

    if !found {
        if s.b.Check(currentPlayer) {
            createPlayerTransition(s.b, s.p, true, false, &transition)
        } else {
            createPlayerTransition(s.b, s.p, false, true, &transition)
        }

        transition.execute()

        s.minimax(depth, moveKey)

        transition.undo()
    }
}

