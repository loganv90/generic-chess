package chess

import (
    "strings"
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

        transpositionMap: map[string][]int{},
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

    transpositionMap map[string][]int
}

func (s *SimpleSearcher) searchWithMinimax(depth int) (MoveKey, error) {
    levels := depth + 1

    s.scoreLevels = make([][]int, levels)
    s.transitionLevels = make([]PlayerTransition, levels)
    s.moveLevels = make([]Array1000[FastMove], levels)
    for i := 0; i < levels; i++ {
        s.scoreLevels[i] = make([]int, s.players)
        s.transitionLevels[i] = PlayerTransition{}
        s.moveLevels[i] = Array1000[FastMove]{}
    }
    s.transpositionMap = map[string][]int{}

    s.b.CalculateMoves()
    moveKey := MoveKey{}
    s.minimaxFirstCall(depth, &moveKey)
    return moveKey, nil
}

func (s *SimpleSearcher) minimaxFirstCall(depth int, moveKey *MoveKey) {
    if depth <= 0 || s.p.getGameOver() {
        panic("no moves in this position")
    }

    found := false
    currentPlayer := s.p.getCurrent()
    transition := &s.transitionLevels[depth]

    s.moveLevels[depth].clear()
    s.b.MovesOfColor(currentPlayer, &s.moveLevels[depth])

    for i := 0; i < len(s.scoreLevels[depth]); i++ {
        s.scoreLevels[depth][i] = -1000000
    }
    
    for i := 0; i < s.moveLevels[depth].count; i++ {
        move := &s.moveLevels[depth].array[i]

        move.execute()

        s.b.CalculateMoves()
        if s.b.Check(currentPlayer) {
            move.undo()
            continue
        }

        createPlayerTransition(s.b, s.p, false, false, transition)

        transition.execute()
        s.minimax(depth-1)
        transition.undo()

        move.undo()

        if s.scoreLevels[depth-1][currentPlayer] > s.scoreLevels[depth][currentPlayer] {
            for i := 0; i < len(s.scoreLevels[depth]); i++ {
                s.scoreLevels[depth][i] = s.scoreLevels[depth-1][i]
            }
            moveKey.XTo = move.toLocation.x
            moveKey.YTo = move.toLocation.y
            moveKey.XFrom = move.fromLocation.x
            moveKey.YFrom = move.fromLocation.y
            found = true
        }
    }

    if !found {
        panic("no moves in this position")
    }
}

func (s *SimpleSearcher) minimax(depth int) {
    s.minimaxCalls++

    builder := strings.Builder{}
    s.b.UniqueString(&builder)
    builder.WriteString("-")
    s.p.UniqueString(&builder)
    uniqueString := builder.String()

    if _, ok := s.transpositionMap[uniqueString]; ok {
        score := s.transpositionMap[uniqueString]
        for i := 0; i < len(s.scoreLevels[depth]); i++ {
            s.scoreLevels[depth][i] = score[i]
        }
        return
    }

    if depth <= 0 || s.p.getGameOver() {
        s.e.eval(s.scoreLevels[depth])
        return
    }

    found := false
    currentPlayer := s.p.getCurrent()
    transition := &s.transitionLevels[depth]

    s.moveLevels[depth].clear()
    s.b.MovesOfColor(currentPlayer, &s.moveLevels[depth])

    for i := 0; i < len(s.scoreLevels[depth]); i++ {
        s.scoreLevels[depth][i] = -1000000
    }
    
    for i := 0; i < s.moveLevels[depth].count; i++ {
        move := &s.moveLevels[depth].array[i]

        move.execute()

        s.b.CalculateMoves()
        if s.b.Check(currentPlayer) {
            move.undo()
            continue
        }

        createPlayerTransition(s.b, s.p, false, false, transition)

        transition.execute()
        s.minimax(depth-1)
        transition.undo()

        move.undo()

        if s.scoreLevels[depth-1][currentPlayer] > s.scoreLevels[depth][currentPlayer] {
            for i := 0; i < len(s.scoreLevels[depth]); i++ {
                s.scoreLevels[depth][i] = s.scoreLevels[depth-1][i]
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
        s.minimax(depth)
        transition.undo()
    }

    newScore := make([]int, s.players)
    for i := 0; i < len(s.scoreLevels[depth]); i++ {
        newScore[i] = s.scoreLevels[depth][i]
    }
    s.transpositionMap[uniqueString] = newScore
}

