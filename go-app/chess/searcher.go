package chess

import (
    "math"
    "fmt"
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
func newSimpleSearcher(b *SimpleBoard, p *SimplePlayerCollection, stop chan bool) *SimpleSearcher {
    return &SimpleSearcher{
        b: b,
        p: p,
        e: newSimpleEvaluator(b, p),

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
    transpositionMapLevels []map[uint64][]int

    maxDepth int
    moveKey MoveKey

    stop chan bool
    stopReached bool
}

func (s *SimpleSearcher) searchWithMinimax(maxDepth int) (MoveKey, error) {
    s.players = s.p.getPlayers()
    s.maxDepth = maxDepth
    s.moveKey = MoveKey{-1, -1, -1, -1, ""}

    s.scoreLevels = make([][]int, maxDepth+1)
    s.transitionLevels = make([]PlayerTransition, maxDepth+1)
    s.moveLevels = make([]Array1000[FastMove], maxDepth+1)
    s.captureMoveLevels = make([]Array1000[FastMove], maxDepth+1)
    s.transpositionMapLevels = make([]map[uint64][]int, maxDepth+1)
    for i := 0; i < maxDepth+1; i++ {
        s.scoreLevels[i] = make([]int, s.players)
        s.transitionLevels[i] = PlayerTransition{}
        s.moveLevels[i] = Array1000[FastMove]{}
        s.captureMoveLevels[i] = Array1000[FastMove]{}
        s.transpositionMapLevels[i] = map[uint64][]int{}
    }

    s.b.CalculateMoves()
    s.minimax(0)

    if s.moveKey.XTo == -1 || s.moveKey.YTo == -1 || s.moveKey.XFrom == -1 || s.moveKey.YFrom == -1 {
        return s.moveKey, fmt.Errorf("No move found")
    }
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

    if _, ok := s.transpositionMapLevels[depth][hash]; ok {
        score := s.transpositionMapLevels[depth][hash]

        for i := 0; i < s.players; i++ {
            s.scoreLevels[depth][i] = score[i]
        }

        return
    }

    if depth >= s.maxDepth || s.p.getGameOver() {
        s.e.eval(s.scoreLevels[depth])

        newScore := make([]int, s.players)
        for i := 0; i < s.players; i++ {
            newScore[i] = s.scoreLevels[depth][i]
        }
        s.transpositionMapLevels[depth][hash] = newScore

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
        s.b.CalculateMoves()
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
    s.transpositionMapLevels[depth][hash] = newScore
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

        if s.liftScore(s.scoreLevels[depth+1][color]) > s.scoreLevels[depth][color] {
            found = true
            for i := 0; i < len(s.scoreLevels[depth]); i++ {
                s.scoreLevels[depth][i] = s.liftScore(s.scoreLevels[depth+1][i])
            }

            if depth <= 0 {
                s.moveKey.XTo = move.toLocation.x
                s.moveKey.YTo = move.toLocation.y
                s.moveKey.XFrom = move.fromLocation.x
                s.moveKey.YFrom = move.fromLocation.y
                s.moveKey.Promotion = ""
            }
        }
    }

    return found
}

func (s *SimpleSearcher) liftScore(score int) int {
    if score > 0 {
        return score - 1
    } else if score < 0 {
        return score + 1
    }

    return score
}

func newParallelSearcher(b *SimpleBoard, p *SimplePlayerCollection, stop chan bool) *ParallelSearcher {
    return &ParallelSearcher{
        b: b,
        p: p,

        stop: stop,
        stopReached: false,
    }
}

type ParallelSearcher struct {
    b *SimpleBoard
    p *SimplePlayerCollection

    stops []chan bool
    result chan *MoveKeyWithScore

    maxDepth int
    moveKey MoveKey

    stop chan bool
    stopReached bool
}

func (s *ParallelSearcher) searchWithMinimax(maxDepth int) (MoveKey, error) {
    s.maxDepth = maxDepth
    s.moveKey = MoveKey{-1, -1, -1, -1, ""}

    s.b.CalculateMoves()
    currentPlayer := s.p.getCurrent()
    moveCount := 0
    moveCount += s.b.captureMoves[currentPlayer].count
    moveCount += s.b.moves[currentPlayer].count

    s.result = make(chan *MoveKeyWithScore, moveCount)
    s.stops = make([]chan bool, moveCount)
    for i := 0; i < moveCount; i++ {
        s.stops[i] = make(chan bool)
    }

    for i := 0; i < moveCount; i++ {
        boardCopy, err := s.b.Copy()
        if err != nil {
            return s.moveKey, err
        }

        playerCollectionCopy, err := s.p.Copy()
        if err != nil {
            return s.moveKey, err
        }

        go minimaxWrapper(boardCopy, playerCollectionCopy, s.stops[i], s.result, s.maxDepth-1, i)
    }

    counter := 0
    score := math.MinInt
    loop := true
    for loop {
        select {
        case <-s.stop:
            for i := 0; i < moveCount; i++ {
                s.stops[i] <- true
            }

            loop = false
        case moveKeyWithScorePtr := <-s.result:
            if moveKeyWithScorePtr != nil {
                if moveKeyWithScorePtr.Score > score {
                    score = moveKeyWithScorePtr.Score
                    s.moveKey.XTo = moveKeyWithScorePtr.XTo
                    s.moveKey.YTo = moveKeyWithScorePtr.YTo
                    s.moveKey.XFrom = moveKeyWithScorePtr.XFrom
                    s.moveKey.YFrom = moveKeyWithScorePtr.YFrom
                    s.moveKey.Promotion = moveKeyWithScorePtr.Promotion
                }
            } else {
                fmt.Println("No move found")
            }

            counter++
            loop = counter < moveCount
        }
    }

    if s.moveKey.XTo == -1 || s.moveKey.YTo == -1 || s.moveKey.XFrom == -1 || s.moveKey.YFrom == -1 {
        return s.moveKey, fmt.Errorf("No move found")
    }
    return s.moveKey, nil
}

func minimaxWrapper(b *SimpleBoard, p *SimplePlayerCollection, stop chan bool, result chan *MoveKeyWithScore, depth int, i int) {
    b.CalculateMoves()
    currentPlayer := p.getCurrent()

    var transition PlayerTransition
    var move FastMove
    captureMoves := &b.captureMoves[currentPlayer]
    moves := &b.moves[currentPlayer]

    if i >= captureMoves.count {
        move = moves.array[i-captureMoves.count]
    } else {
        move = captureMoves.array[i]
    }

    move.execute()
    b.CalculateMoves()
    if b.Check(currentPlayer) {
        result <- nil
        return
    }
    createPlayerTransition(b, p, false, false, &transition)
    transition.execute()

    searcher := newSimpleSearcher(b, p, stop)

    searcher.searchWithMinimax(depth)
    result <- &MoveKeyWithScore{
        XTo: move.toLocation.x,
        YTo: move.toLocation.y,
        XFrom: move.fromLocation.x,
        YFrom: move.fromLocation.y,
        Promotion: "",
        Score: searcher.scoreLevels[0][currentPlayer],
    }
}

