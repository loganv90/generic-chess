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
    }
}

type SimpleEvaluator struct {
    b *SimpleBoard
    p *SimplePlayerCollection

    material []int
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

    pieceLocations := e.b.getPieceLocations()
    e.evalMaterial(pieceLocations, score)

    // Piece position comparison (piece-square tables)
    // we need: the locations of each piece by player

    // Mobility comparison
    // we need: the moves each piece can make including attacking ally pieces
}

func (e *SimpleEvaluator) evalMaterial(pieceLocations []Array100[*Point], score []int) {
    totalMaterial := 0

    for color, locations := range pieceLocations {
        materialCount := 0

        for i := 0; i < locations.count; i++ {
            location := locations.array[i]

            piece := e.b.getPiece(location)
            if piece == nil {
                continue
            }

            materialCount += piece.value()
        }

        totalMaterial += materialCount
        e.material[color] = materialCount
    }

    for color := range score {
        percentage := int(float64(e.material[color]) / float64(totalMaterial) * 100)
        score[color] += percentage
    }
}

