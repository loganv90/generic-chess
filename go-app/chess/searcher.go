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
func newSimpleSearcher(g Game) *SimpleSearcher {
    board := g.getBoard()
    playerCollection := g.getPlayerCollection()

    return &SimpleSearcher{
        b: board,
        p: playerCollection,
        e: newSimpleEvaluator(board, playerCollection),

        players: playerCollection.getPlayers(),
        minimaxCalls: 0,

        scoreLevels: [][]int{},
        transitionLevels: []PlayerTransition{},
        moveLevels: []Array1000[FastMove]{},

        transpositionMap: map[string]MoveKeyAndScore{},
    }
}

type SimpleSearcher struct {
    b *SimpleBoard
    p *SimplePlayerCollection
    e *SimpleEvaluator

    players int
    minimaxCalls int

    scoreLevels [][]int
    transitionLevels []PlayerTransition
    moveLevels []Array1000[FastMove]

    transpositionMap map[string]MoveKeyAndScore
}

func (s *SimpleSearcher) search(depth int) (MoveKey, error) {
    levels := depth + 1

    s.scoreLevels = make([][]int, levels)
    s.transitionLevels = make([]PlayerTransition, levels)
    s.moveLevels = make([]Array1000[FastMove], levels)
    for i := 0; i < levels; i++ {
        s.scoreLevels[i] = make([]int, s.players)
        s.transitionLevels[i] = PlayerTransition{}
        s.moveLevels[i] = Array1000[FastMove]{}
    }

    moveKey := MoveKey{}

    s.b.CalculateMoves()
    s.minimax(depth, &moveKey)

    fmt.Println(s.b.Print())

    return moveKey, nil
}

func (s *SimpleSearcher) minimax(depth int, moveKey *MoveKey) {
    s.minimaxCalls++

    gameOver := s.p.getGameOver()
    if depth == 0 || gameOver {
        s.e.eval(s.scoreLevels[depth])
        return
    }

    found := false
    currentPlayer := s.p.getCurrent()
    moves := &s.moveLevels[depth]
    transition := &s.transitionLevels[depth]

    moves.clear()
    s.b.MovesOfColor(currentPlayer, moves)

    for i := 0; i < len(s.scoreLevels[depth]); i++ {
        s.scoreLevels[depth][i] = -1000000
    }
    
    for i := 0; i < moves.count; i++ {
        move := &moves.array[i]
        if move.allyDefense {
            continue
        }


        move.execute()
        s.b.CalculateMovesDynamic(move)
        //s.b.CalculateMoves()


        if s.b.Check(currentPlayer) {
            move.undo()
            s.b.CalculateMovesDynamic(move)

            continue
        }

        createPlayerTransition(s.b, s.p, false, false, transition)

        transition.execute()

        s.minimax(depth-1, nil)

        transition.undo()

        move.undo()
        s.b.CalculateMovesDynamic(move)


        if s.scoreLevels[depth-1][currentPlayer] > s.scoreLevels[depth][currentPlayer] {
            for i := 0; i < len(s.scoreLevels[depth]); i++ {
                s.scoreLevels[depth][i] = s.scoreLevels[depth-1][i]
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
            createPlayerTransition(s.b, s.p, true, false, transition)
        } else {
            createPlayerTransition(s.b, s.p, false, true, transition)
        }

        transition.execute()

        s.minimax(depth, moveKey)

        transition.undo()
    }
}

