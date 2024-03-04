package chess

import (
    "fmt"
)

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
    eval() ([]int, error)
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

func (e *SimpleEvaluator) eval() ([]int, error) {
    gameOver, err := e.p.getGameOver()
    if err != nil {
        return nil, err
    }

    players := e.p.getPlayers()
    scores := make([]int, players)

    if gameOver {
        winner, _ := e.p.getWinner()

        if winner < 0 || winner >= players {
            return scores, nil
        } else {
            for player := range scores {
                scores[player] = -100000
            }
            scores[winner] = 100000

            return scores, nil
        }
    }

    pieceLocations := e.b.getPieceLocations()
    materialScores, err := e.evalMaterial(pieceLocations)
    if err != nil {
        return nil, err
    }

    // Piece position comparison (piece-square tables)
    // we need: the locations of each piece by player

    // Mobility comparison
    // we need: the moves each piece can make including attacking ally pieces

    for player := range scores {
        scores[player] += materialScores[player]
    }

    return scores, nil
}

func (e *SimpleEvaluator) evalMaterial(pieceLocations [][]Point) ([]int, error) {
    // idea is to compare our material to the leading player's material

    leadingMaterial := 0
    material := make([]int, len(pieceLocations))
    scores := make([]int, len(pieceLocations))

    for color, locations := range pieceLocations {
        materialCount := 0
        for _, location := range locations {
            piece, ok := e.b.getPiece(location)
            if !ok || piece == nil {
                return nil, fmt.Errorf("no piece at location")
            }
            materialCount += piece.getValue()
        }

        material[color] = materialCount
        leadingMaterial = max(leadingMaterial, materialCount)
    }

    for color, materialCount := range material {
        scores[color] = materialCount - leadingMaterial
    }

    // if we're leading in material, we incentivize low total material
    // if we're behind in material, we incentivize high total material
    // totalMaterialFactor := 1000 * materialDifference / totalMaterial // more + if winning, more - if losing
    // return materialDifference + totalMaterialFactor, nil

    return scores, nil
}

