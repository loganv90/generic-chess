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
type Evaluator interface {
    eval(color string) int
}

func newSimpleEvaluator(b Board, p PlayerCollection) (*SimpleEvaluator, error) {
    return &SimpleEvaluator{
        b: b,
        p: p,
    }, nil
}

type SimpleEvaluator struct {
    b Board
    p PlayerCollection
}

func (e *SimpleEvaluator) eval(color string) (int, error) {
    gameOver, err := e.p.getGameOver()
    if err != nil {
        return 0, err
    }

    winner, err := e.p.getWinner()
    if err != nil {
        return 0, err
    }

    if gameOver {
        if winner == "" {
            return 0, nil
        } else if winner == color {
            return 100000, nil
        } else {
            return -100000, nil
        }
    }

    pieceLocations := e.b.getPieceLocations()

    materialScore, err := e.evalMaterial(color, pieceLocations)
    if err != nil {
        return 0, err
    }

    // Piece position comparison (piece-square tables)
    // we need: the locations of each piece by player

    // Mobility comparison
    // we need: the moves each piece can make including attacking ally pieces

    return materialScore, nil
}

func (e *SimpleEvaluator) evalMaterial(ourColor string, pieceLocations map[string][]*Point) (int, error) {
    // idea is to compare our material to the leading player's material

    leadingMaterial := 0
    ourMaterial := 0
    totalMaterial := 0

    for color, locations := range pieceLocations {
        materialCount := 0
        for _, location := range locations {
            piece, err := e.b.getPiece(location)
            if err != nil || piece == nil {
                return 0, err
            }
            materialCount += piece.getValue()
        }

        if ourColor == color {
            ourMaterial = materialCount
        } else {
            leadingMaterial = max(leadingMaterial, materialCount)
        }
        totalMaterial += materialCount
    }

    materialDifference := ourMaterial - leadingMaterial

    // if we're leading in material, we incentivize low total material
    // if we're behind in material, we incentivize high total material
    // totalMaterialFactor := 1000 * materialDifference / totalMaterial // more + if winning, more - if losing
    // return materialDifference + totalMaterialFactor, nil

    return materialDifference, nil
}

