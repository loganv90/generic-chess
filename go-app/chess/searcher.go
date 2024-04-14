package chess

import (
    "math"
)

/*
Chess Search:
Alpha beta search (only works for 2 player games)
Iterative deepening
Transposition table (repeated and colour flipped positions)
Keep track of killer moves (moves that caused cut off) either by [piece][to] or [from][to]
Quiescence search (do all captures on last capture/promotion square)
Move ordering (Hashed moves, winning captures/promotions, equal captures, killer moves, non-captures, losing captures) (only works with alpha beta)
Create a hash function to uniquely identify the chess position
q2k2q1/2nqn2b/1n1P1n1b/2rnr2Q/1NQ1QN1Q/3Q3B/2RQR2B/Q2K2Q1 w - - this position causes engines to explode
*/

/*
Responsible for:
- searching for moves given the current state of the game
*/
func newSimpleSearcher(g Game, stop chan bool) *SimpleSearcher {
    board := g.getBoard()
    playerCollection := g.getPlayerCollection()

    return &SimpleSearcher{
        b: board,
        p: playerCollection,
        e: newSimpleEvaluator(board, playerCollection),

        players: playerCollection.getPlayers(),

        scoreLevels: [][]int{},
        transitionLevels: []PlayerTransition{},
        moveLevels: []Array1000[FastMove]{},
        captureMoveLevels: []Array1000[FastMove]{},

        transpositionMap: map[uint64][]int{},

        stop: stop,
        stopReached: false,
    }
}

type SimpleSearcher struct {
    b *SimpleBoard
    p *SimplePlayerCollection
    e *SimpleEvaluator

    players int

    scoreLevels [][]int
    transitionLevels []PlayerTransition
    moveLevels []Array1000[FastMove]
    captureMoveLevels []Array1000[FastMove]

    maxDepth int
    moveKey MoveKey
    eliminated []bool

    transpositionMap map[uint64][]int

    stop chan bool
    stopReached bool
}

func (s *SimpleSearcher) searchWithMinimax(maxDepth int) (MoveKey, error) {
    s.maxDepth = maxDepth

    s.scoreLevels = make([][]int, maxDepth+1)
    s.transitionLevels = make([]PlayerTransition, maxDepth+1)
    s.moveLevels = make([]Array1000[FastMove], maxDepth+1)
    s.captureMoveLevels = make([]Array1000[FastMove], maxDepth+1)
    for i := 0; i < maxDepth+1; i++ {
        s.scoreLevels[i] = make([]int, s.players)
        s.transitionLevels[i] = PlayerTransition{}
        s.moveLevels[i] = Array1000[FastMove]{}
        s.captureMoveLevels[i] = Array1000[FastMove]{}
    }

    s.eliminated = make([]bool, s.players)
    s.transpositionMap = map[uint64][]int{}

    s.b.CalculateMoves()
    s.minimax(0)

    return s.moveKey, nil
}

func (s *SimpleSearcher) minimax(depth int) {
    select {
    case <-s.stop:
        s.stopReached = true
    default:
    }
    if s.stopReached {
        return
    }

    hash := s.b.ZobristHash() ^ s.p.ZobristHash()

    if _, ok := s.transpositionMap[hash]; ok {
        score := s.transpositionMap[hash]
        for i := 0; i < s.players; i++ {
            s.scoreLevels[depth][i] = score[i]
        }

        return
    }

    // the problem is that we don't check whether or not they're in check here
    // we can also save moves to the transpositionMap with unknown depth
    if depth >= s.maxDepth || s.p.getGameOver() {
        s.evaluate(depth)

        newScore := make([]int, s.players)
        for i := 0; i < s.players; i++ {
            newScore[i] = s.scoreLevels[depth][i]
        }
        s.transpositionMap[hash] = newScore

        return
    }

    currentPlayer := s.p.getCurrent()
    s.copyMoves(depth, currentPlayer)

    for i := 0; i < s.players; i++ {
        s.scoreLevels[depth][i] = math.MinInt
    }

    transition := &s.transitionLevels[depth]
    found1 := s.recurse(depth, currentPlayer, &s.captureMoveLevels[depth], transition)
    found2 := s.recurse(depth, currentPlayer, &s.moveLevels[depth], transition)

    if !found1 && !found2 {
        if depth <= 0 {
            panic("no moves in this position")
        }

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
    for i := 0; i < s.players; i++ {
        newScore[i] = s.scoreLevels[depth][i]
    }
    s.transpositionMap[hash] = newScore
}

func (s *SimpleSearcher) copyMoves(depth int, color int) {
    moves := &s.b.moves[color]
    captureMoves := &s.b.captureMoves[color]

    s.moveLevels[depth].clear()
    s.captureMoveLevels[depth].clear()

    s.moveLevels[depth].count = moves.count
    for i := 0; i < moves.count; i++ {
        s.moveLevels[depth].array[i] = moves.array[i]
    }

    s.captureMoveLevels[depth].count = captureMoves.count
    for i := 0; i < captureMoves.count; i++ {
        s.captureMoveLevels[depth].array[i] = captureMoves.array[i]
    }
}

func (s *SimpleSearcher) recurse(depth int, color int, moves *Array1000[FastMove], transition *PlayerTransition) bool {
    found := false

    for i := 0; i < moves.count; i++ {
        move := &moves.array[i]

        move.execute()

        s.b.CalculateMoves()
        if s.b.Check(color) {
            move.undo()
            continue
        }

        createPlayerTransition(s.b, s.p, false, false, transition)

        transition.execute()
        s.minimax(depth+1)
        transition.undo()

        move.undo()

        if s.scoreLevels[depth+1][color] > s.scoreLevels[depth][color] {
            found = true
            for i := 0; i < len(s.scoreLevels[depth]); i++ {
                s.scoreLevels[depth][i] = s.scoreLevels[depth+1][i]
            }

            if depth <= 0 {
                s.moveKey.XTo = move.toLocation.x
                s.moveKey.YTo = move.toLocation.y
                s.moveKey.XFrom = move.fromLocation.x
                s.moveKey.YFrom = move.fromLocation.y
            }
        }
    }

    return found
}

func (s *SimpleSearcher) canMove(color int, moves *Array1000[FastMove]) bool {
    found := false

    for i := 0; i < moves.count; i++ {
        move := &moves.array[i]

        move.execute()

        s.b.CalculateMoves()
        if s.b.Check(color) {
            move.undo()
            continue
        }

        found = true
        move.undo()
    }

    return found
}

func (s *SimpleSearcher) evaluate(depth int) {
    /*
    for i := 0; i < s.players; i++ {
        s.eliminated[i] = false
        if !s.p.playersAlive[i] {
            continue
        }
        if !s.b.Check(i) {
            continue
        }

        s.copyMoves(depth, i)
        if s.canMove(i, &s.captureMoveLevels[i]) {
            continue
        }
        if s.canMove(i, &s.moveLevels[i]) {
            continue
        }

        s.eliminated[i] = true
        s.p.eliminate(i)
    }

    lastAlive := -1
    numAlive := 0
    prevGameOver := s.p.getGameOver()
    prevWinner := s.p.getWinner()
    for i := 0; i < s.players; i++ {
        if s.p.playersAlive[i] {
            lastAlive = i
            numAlive++
        }
    }

    if numAlive == 1 {
        s.p.setGameOver(true)
        s.p.setWinner(lastAlive)
    }
    */

    s.e.eval(s.scoreLevels[depth])

    /*
    if numAlive == 1 {
        s.p.setGameOver(prevGameOver)
        s.p.setWinner(prevWinner)
    }

    for i := 0; i < s.players; i++ {
        if s.eliminated[i] {
            s.p.restore(i)
        }
    }
    */
}

