package chess

import (
    "math"
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
    players := p.getPlayers()

    attackedSquares := make([][][]bool, players)
    for color := 0; color < players; color++ {
        attackedSquares[color] = make([][]bool, b.y)
        for y := 0; y < b.y; y++ {
            attackedSquares[color][y] = make([]bool, b.x)
        }
    }

    return &SimpleEvaluator{
        b: b,
        p: p,

        totalMaterial: 0,
        totalMobility: 0,
        totalPosition: 0,
        totalKingSafety: 0,

        material: make([]int, players),
        mobility: make([]int, players),
        position: make([]int, players),
        kingSafety: make([]int, players),
        pressure: make([]int, players),

        players: players,
        attackedSquares: attackedSquares,
    }
}

type SimpleEvaluator struct {
    b *SimpleBoard
    p *SimplePlayerCollection

    totalMaterial int
    totalMobility int
    totalPosition int
    totalKingSafety int
    totalPressure int

    material []int
    mobility []int
    position []int
    kingSafety []int
    pressure []int

    players int
    attackedSquares [][][]bool // [color][y][x]
}

func (e *SimpleEvaluator) eval(score []int) {
    if e.p.getGameOver() {
        winner := e.p.getWinner()

        if winner < 0 {
            for color := range score {
                score[color] = 0
            }
        } else {
            for color := range score {
                score[color] = math.MinInt
            }
            score[winner] = math.MaxInt
        }

        return
    }

    for color := 0; color < e.players; color++ {
        score[color] = 0
        e.totalMaterial = 0
        e.totalMobility = 0
        e.totalPosition = 0
        e.material[color] = 0
        e.mobility[color] = 0
        e.position[color] = 0
        e.kingSafety[color] = 0
        e.pressure[color] = 0
    }

    e.evalMaterial()
    e.evalPosition()
    e.evalMobility()
    e.evalKingSafety()

    for color := 0; color < e.players; color++ {
        if !e.p.playersAlive[color] {
            score[color] = math.MinInt
            continue
        }

        percentage := 0

        percentage += int(
            float64(e.material[color] + e.position[color]) / 
            float64(e.totalMaterial + e.totalPosition) * 10000,
        ) * 10 // weighted 10 times

        percentage += int(
            float64(e.mobility[color]) /
            float64(e.totalMobility) * 10000,
        ) * 1 // weighted 1 times

        percentage += int(
            float64(e.kingSafety[color]) /
            float64(e.totalKingSafety) * 10000,
        ) * 1 // weighted 1 times

        percentage -= int(
            float64(e.pressure[color]) /
            float64(e.totalPressure) * 10000,
        ) * 2 // weighted 2 times

        score[color] = percentage
    }
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
    for color := 0; color < e.players; color++ {
        value := e.b.moves[color].count + e.b.captureMoves[color].count + e.b.defenseMoves[color].count - e.b.queenMoveCount[color] - e.b.kingMoveCount[color]
        e.mobility[color] = value
        e.totalMobility += value
    }
}

func (e *SimpleEvaluator) evalKingSafety() {
    e.totalKingSafety = e.players
    e.totalPressure = 0

    for color := 0; color < e.players; color++ {
        e.kingSafety[color] = 1
        e.pressure[color] = 0
        for y := 0; y < e.b.y; y++ {
            for x := 0; x < e.b.x; x++ {
                e.attackedSquares[color][y][x] = false
            }
        }
    }

    for color := 0; color < e.players; color++ {
        captureMoves := e.b.captureMoves[color]
        for i := 0; i < captureMoves.count; i++ {
            to := captureMoves.array[i].toLocation
            for c := 0; c < e.players; c++ {
                if c != color {
                    e.attackedSquares[c][to.y][to.x] = true
                }
            }

            capturedPiece := captureMoves.array[i].oldPiece.array[1]
            if capturedPiece != nil {
                e.pressure[capturedPiece.color] += capturedPiece.value()
                e.totalPressure += capturedPiece.value()
            }
        }

        moves := e.b.moves[color]
        for i := 0; i < moves.count; i++ {
            to := moves.array[i].toLocation
            for c := 0; c < e.players; c++ {
                if c != color {
                    e.attackedSquares[c][to.y][to.x] = true
                }
            }
        }
    }

    for color := 0; color < e.players; color++ {
        king := e.b.kingLocations[color]

        if e.attackedSquares[color][king.y][king.x] {
            e.kingSafety[color] -= 1
            e.totalKingSafety -= 1
        } else {
            e.kingSafety[color] += 1
            e.totalKingSafety += 1
            continue
        }

        for _, direction := range queen_directions {
            to := e.b.addIndex(direction, king)
            if to == nil {
                continue
            }

            if !e.attackedSquares[color][to.y][to.x] {
                e.kingSafety[color] += 1
                e.totalKingSafety += 1
                break
            }
        }
    }
}

