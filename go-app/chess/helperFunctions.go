package chess

func createSimpleBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    color1 := "black"
    color2 := "white"

    simpleBoard, err := newSimpleBoard(&Point{8, 8})
    if err != nil {
        return nil, err
    }

    simpleBoard.setPiece(&Point{0, 0}, newRook(color1, false))
    simpleBoard.setPiece(&Point{1, 0}, newKnight(color1))
    simpleBoard.setPiece(&Point{2, 0}, newBishop(color1))
    simpleBoard.setPiece(&Point{3, 0}, newQueen(color1))
    simpleBoard.setPiece(&Point{4, 0}, newKing(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{5, 0}, newBishop(color1))
    simpleBoard.setPiece(&Point{6, 0}, newKnight(color1))
    simpleBoard.setPiece(&Point{7, 0}, newRook(color1, false))

    simpleBoard.setPiece(&Point{0, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{1, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{2, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{3, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{4, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{5, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{6, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{7, 1}, newPawn(color1, false, 0, 1))

    simpleBoard.setPiece(&Point{0, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{1, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{2, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{3, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{4, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{5, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{6, 6}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{7, 6}, newPawn(color2, false, 0, -1))

    simpleBoard.setPiece(&Point{0, 7}, newRook(color2, false))
    simpleBoard.setPiece(&Point{1, 7}, newKnight(color2))
    simpleBoard.setPiece(&Point{2, 7}, newBishop(color2))
    simpleBoard.setPiece(&Point{3, 7}, newQueen(color2))
    simpleBoard.setPiece(&Point{4, 7}, newKing(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{5, 7}, newBishop(color2))
    simpleBoard.setPiece(&Point{6, 7}, newKnight(color2))
    simpleBoard.setPiece(&Point{7, 7}, newRook(color2, false))

    return simpleBoard, nil
}
