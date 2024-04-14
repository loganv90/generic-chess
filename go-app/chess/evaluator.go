package chess

import (
    "math"
    // "fmt"
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
func newSimpleEvaluator(b *SimpleBoard, p *SimplePlayerCollection) *SimpleEvaluator {
    return &SimpleEvaluator{
        b: b,
        p: p,

        totalMaterial: 0,
        totalMobility: 0,
        totalPosition: 0,

        material: make([]int, p.getPlayers()),
        mobility: make([]int, p.getPlayers()),
        position: make([]int, p.getPlayers()),
    }
}

type SimpleEvaluator struct {
    b *SimpleBoard
    p *SimplePlayerCollection

    totalMaterial int
    totalMobility int
    totalPosition int

    material []int
    mobility []int
    position []int
}

func (e *SimpleEvaluator) eval(score []int) {
    if e.p.getGameOver() {
        winner := e.p.getWinner()

        if winner < 0 {
            for player := range score {
                score[player] = 0
            }
        } else {
            for player := range score {
                score[player] = math.MinInt
            }
            score[winner] = math.MaxInt
        }

        /*
        if e.b.getPiece(e.b.getIndex(0, 3)) == e.b.getAllPiece(1, QUEEN) && e.b.getPiece(e.b.getIndex(0, 0)) == e.b.getAllPiece(0, KING_U_M) {
            fmt.Println(e.p.getCurrent())
            fmt.Println(score)
            fmt.Println(e.b.Print())
        }
        */

        return
    }

    for player := range score {
        score[player] = 0
        e.totalMaterial = 0
        e.totalMobility = 0
        e.totalPosition = 0
        e.material[player] = 0
        e.mobility[player] = 0
        e.position[player] = 0
    }

    e.evalMaterial()
    e.evalPosition()
    e.evalMobility()

    for color := range score {
        if !e.p.playersAlive[color] {
            score[color] = math.MinInt
            continue
        }

        percentage := 0

        if !e.b.Check(color) {
            percentage += 10000 * 4 // weighted 2 times, bonus for not being in check
        }

        percentage += int(
            float64(e.material[color] + e.position[color]) / 
            float64(e.totalMaterial + e.totalPosition) * 10000,
        ) * 10 // weighted 10 times

        percentage += int(
            float64(e.mobility[color]) /
            float64(e.totalMobility) * 10000,
        ) * 2 // weighted 2 times

        score[color] = percentage
    }

    /*
    if e.b.getPiece(e.b.getIndex(0, 3)) == e.b.getAllPiece(1, QUEEN) && e.b.getPiece(e.b.getIndex(0, 0)) == e.b.getAllPiece(0, KING_U_M) {
        fmt.Println(e.p.getCurrent())
        fmt.Println(score)
        fmt.Println(e.b.Print())
    }
    */
}

func (e *SimpleEvaluator) evalMaterial() {
    pieces := e.b.pieces
    for y := 0; y < e.b.y; y++ {
        for x := 0; x < e.b.x; x++ {
            piece := pieces[y][x]
            if piece == nil {
                continue
            }

            value := piece.value()
            e.material[piece.color] += value
            e.totalMaterial += value
        }
    }
}

func (e *SimpleEvaluator) evalPosition() {
    pieces := e.b.pieces
    for y := 0; y < e.b.y; y++ {
        for x := 0; x < e.b.x; x++ {
            piece := pieces[y][x]
            if piece == nil {
                continue
            }

            if piece.isPawn() {
                value := e.b.pieceSquareTables[piece.index][y][x]
                e.position[piece.color] += value
                e.totalPosition += value
            }
        }
    }
}

func (e *SimpleEvaluator) evalMobility() {
    for color := range e.mobility {
        value := e.b.moves[color].count + e.b.captureMoves[color].count + e.b.defenseMoves[color].count - e.b.queenMoveCount[color] - e.b.kingMoveCount[color]
        e.mobility[color] = value
        e.totalMobility += value
    }
}

