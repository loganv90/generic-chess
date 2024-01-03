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

func createSimpleFourPlayerBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    color1 := "black"
    color2 := "white"
    color3 := "red"
    color4 := "blue"

    simpleBoard, err := newSimpleBoard(&Point{14, 14})
    if err != nil {
        return nil, err
    }

    simpleBoard.setPiece(&Point{3, 0}, newRook(color1, false))
    simpleBoard.setPiece(&Point{4, 0}, newKnight(color1))
    simpleBoard.setPiece(&Point{5, 0}, newBishop(color1))
    simpleBoard.setPiece(&Point{6, 0}, newQueen(color1))
    simpleBoard.setPiece(&Point{7, 0}, newKing(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{8, 0}, newBishop(color1))
    simpleBoard.setPiece(&Point{9, 0}, newKnight(color1))
    simpleBoard.setPiece(&Point{10, 0}, newRook(color1, false))

    simpleBoard.setPiece(&Point{3, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{4, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{5, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{6, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{7, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{8, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{9, 1}, newPawn(color1, false, 0, 1))
    simpleBoard.setPiece(&Point{10, 1}, newPawn(color1, false, 0, 1))

    simpleBoard.setPiece(&Point{3, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{4, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{5, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{6, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{7, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{8, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{9, 12}, newPawn(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{10, 12}, newPawn(color2, false, 0, -1))

    simpleBoard.setPiece(&Point{3, 13}, newRook(color2, false))
    simpleBoard.setPiece(&Point{4, 13}, newKnight(color2))
    simpleBoard.setPiece(&Point{5, 13}, newBishop(color2))
    simpleBoard.setPiece(&Point{6, 13}, newQueen(color2))
    simpleBoard.setPiece(&Point{7, 13}, newKing(color2, false, 0, -1))
    simpleBoard.setPiece(&Point{8, 13}, newBishop(color2))
    simpleBoard.setPiece(&Point{9, 13}, newKnight(color2))
    simpleBoard.setPiece(&Point{10, 13}, newRook(color2, false))

    simpleBoard.setPiece(&Point{0, 3}, newRook(color3, false))
    simpleBoard.setPiece(&Point{0, 4}, newKnight(color3))
    simpleBoard.setPiece(&Point{0, 5}, newBishop(color3))
    simpleBoard.setPiece(&Point{0, 6}, newQueen(color3))
    simpleBoard.setPiece(&Point{0, 7}, newKing(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{0, 8}, newBishop(color3))
    simpleBoard.setPiece(&Point{0, 9}, newKnight(color3))
    simpleBoard.setPiece(&Point{0, 10}, newRook(color3, false))

    simpleBoard.setPiece(&Point{1, 3}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 4}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 5}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 6}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 7}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 8}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 9}, newPawn(color3, false, 1, 0))
    simpleBoard.setPiece(&Point{1, 10}, newPawn(color3, false, 1, 0))

    simpleBoard.setPiece(&Point{12, 3}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 4}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 5}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 6}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 7}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 8}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 9}, newPawn(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{12, 10}, newPawn(color4, false, -1, 0))

    simpleBoard.setPiece(&Point{13, 3}, newRook(color4, false))
    simpleBoard.setPiece(&Point{13, 4}, newKnight(color4))
    simpleBoard.setPiece(&Point{13, 5}, newBishop(color4))
    simpleBoard.setPiece(&Point{13, 6}, newQueen(color4))
    simpleBoard.setPiece(&Point{13, 7}, newKing(color4, false, -1, 0))
    simpleBoard.setPiece(&Point{13, 8}, newBishop(color4))
    simpleBoard.setPiece(&Point{13, 9}, newKnight(color4))
    simpleBoard.setPiece(&Point{13, 10}, newRook(color4, false))

    return simpleBoard, nil
}

