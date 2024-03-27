package chess

/*
Chess Evaluation:
Material count
Piece mobility with move bonus for different pieces (including moves protecting ally pieces, and not including moves to controlled squares) low queen bonus
Piece locations piece-square tables (dp for knight 2 moves potential squares) tables for opening and endgame, then interpolate
Penalty isolated pawns worth less than chained pawns
Penalty for attacked squares close to king
Penalty for lots of mobility if king were a queen
Bonus for attacking close to own king
Bonus for pinning pieces to more valuable pieces
Bonus for queen-rook, queen-bishop, bishop-bishop, rook-rook combos
*/

/*
Responsible for:
- evaluating a board and returning a score
*/
func newSimpleEvaluator(b *SimpleBoard, p *SimplePlayerCollection) *SimpleEvaluator {
    return &SimpleEvaluator{
        b: b,
        p: p,

        material: make([]int, p.getPlayers()),
        mobility: make([]int, p.getPlayers()),
    }
}

type SimpleEvaluator struct {
    b *SimpleBoard
    p *SimplePlayerCollection

    material []int
    mobility []int
}

func (e *SimpleEvaluator) eval(score []int) {
    for player := range score {
        score[player] = 0
    }

    gameOver := e.p.getGameOver()
    if gameOver {
        winner := e.p.getWinner()
        if winner < 0 {
            return
        }

        for player := range score {
            score[player] = -100000
        }
        score[winner] = 100000

        return
    }

    e.evalMaterial(score)
    e.evalMobility(score)

    // Piece position comparison (piece-square tables)
    // we need: the locations of each piece by player

    // Mobility comparison
    // we need: the moves each piece can make including attacking ally pieces
}

func (e *SimpleEvaluator) evalMaterial(score []int) {
    totalMaterial := 0
    for color := range score {
        e.material[color] = 0
    }

    pieces := e.b.pieces
    for y := 0; y < e.b.y; y++ {
        for x := 0; x < e.b.x; x++ {
            piece := pieces[y][x]
            if piece == nil {
                continue
            }

            value := piece.value()
            e.material[piece.color] += value
            totalMaterial += value
        }
    }

    for color := range score {
        percentage := int(float64(e.material[color]) / float64(totalMaterial) * 100) * 10 // weighted 10 times
        score[color] += percentage
    }
}

func (e *SimpleEvaluator) evalMobility(score []int) {
    totalMobility := 0
    for color := range score {
        e.mobility[color] = 0
    }

    for color := range score {
        value := e.b.moves[color].count - e.b.queenMoveCount[color] - e.b.kingMoveCount[color]
        e.mobility[color] = value
        totalMobility += value
    }

    for color := range score {
        percentage := int(float64(e.mobility[color]) / float64(totalMobility) * 100) // weighted 1 time
        score[color] += percentage
    }
}

