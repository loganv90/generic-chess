package chess

import (
    "fmt"
)

func createSimpleMove(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
) (FastMove, error) {
	piece, ok := b.getPiece(fromLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at fromLocation")
	}

	newPiece, err := piece.copy()
    if err != nil {
        return FastMove{}, err
    }

	capturedPiece, ok := b.getPiece(toLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at toLocation")
	}

	enPassant, err := b.getEnPassant(color)
	if err != nil {
		return FastMove{}, err
	}

    vulnerables, err := b.getVulnerables(color)
    if err != nil {
        return FastMove{}, err
    }

    // TODO fix
    vstart := Point{-1, -1}
    vend := Point{-1, -1}
    if len(vulnerables) > 0 {
        vminx := vulnerables[0].x
        vminy := vulnerables[0].y
        vmaxx := vulnerables[0].x
        vmaxy := vulnerables[0].y
        for _, v := range vulnerables {
            vminx = min(vminx, v.x)
            vminy = min(vminy, v.y)
            vmaxx = max(vmaxx, v.x)
            vmaxy = max(vmaxy, v.y)
        }
        vstart = Point{vminx, vminy}
        vend = Point{vmaxx, vmaxy}
    }

    newPieceAndLocations := Array3[pieceAndLocation]{}
    newPieceAndLocations.append(pieceAndLocation{piece: newPiece, location: toLocation})

    oldPieceAndLocations := Array3[pieceAndLocation]{}
    oldPieceAndLocations.append(pieceAndLocation{piece: piece, location: fromLocation})
    oldPieceAndLocations.append(pieceAndLocation{piece: capturedPiece, location: toLocation})

    newEnPassantTarget := Point{-1, -1}
    newEnPassantRisk := Point{-1, -1}

    oldEnPassantTarget := enPassant.target
    oldEnPassantRisk := enPassant.pieceLocation

    newVulnerableStart := Point{-1, -1}
    newVulnerableEnd := Point{-1, -1}

    oldVulnerableStart := vstart
    oldVulnerableEnd := vend

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: color,
        allyDefense: false,
        newPieceAndLocation: newPieceAndLocations,
        oldPieceAndLocation: oldPieceAndLocations,
        newEnPassantTarget: newEnPassantTarget,
        oldEnPassantTarget: oldEnPassantTarget,
        newEnPassantRisk: newEnPassantRisk,
        oldEnPassantRisk: oldEnPassantRisk,
        newVulnerableStart: newVulnerableStart,
        oldVulnerableStart: oldVulnerableStart,
        newVulnerableEnd: newVulnerableEnd,
        oldVulnerableEnd: oldVulnerableEnd,
    }, nil
}

func createRevealEnPassantMove(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    target Point,
) (FastMove, error) {
    m, err := createSimpleMove(b, color, fromLocation, toLocation)
    if err != nil {
        return m, err
    }

    m.newEnPassantTarget = target
    m.oldEnPassantRisk = toLocation

    return m, nil
}

func createCaptureEnPassantMove(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
) (FastMove, error) {
    m, err := createSimpleMove(b, color, fromLocation, toLocation)
    if err != nil {
        return m, err
    }

    // TODO fix
    possibleEnPassant, err := b.possibleEnPassant(color, toLocation)
    if err == nil {
        for _, enPassant := range possibleEnPassant {
            capturedPiece, ok := b.getPiece(enPassant.pieceLocation)
            if !ok {
                return m, fmt.Errorf("no piece at enPassant.pieceLocation")
            }

            m.newPieceAndLocation.append(pieceAndLocation{piece: nil, location: enPassant.pieceLocation})
            m.oldPieceAndLocation.append(pieceAndLocation{piece: capturedPiece, location: enPassant.pieceLocation})
        }
    }

    return m, nil
}

type FastMove struct {
    b Board
    fromLocation Point
    toLocation Point
    color int
    allyDefense bool

    newPieceAndLocation Array3[pieceAndLocation]
    oldPieceAndLocation Array3[pieceAndLocation]

    newEnPassantTarget Point
    oldEnPassantTarget Point

    newEnPassantRisk Point
    oldEnPassantRisk Point

    newVulnerableStart Point
    oldVulnerableStart Point

    newVulnerableEnd Point
    oldVulnerableEnd Point
}

