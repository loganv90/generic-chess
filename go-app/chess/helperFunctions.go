package chess

func createSimpleBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    black := 1
    white := 0

    simpleBoard, err := newSimpleBoard(Point{8, 8}, 2)
    if err != nil {
        return nil, err
    }

    simpleBoard.setPiece(Point{0, 0}, newRook(black, false))
    simpleBoard.setPiece(Point{1, 0}, newKnight(black))
    simpleBoard.setPiece(Point{2, 0}, newBishop(black))
    simpleBoard.setPiece(Point{3, 0}, newQueen(black))
    simpleBoard.setPiece(Point{4, 0}, newKing(black, false, 0, 1))
    simpleBoard.setPiece(Point{5, 0}, newBishop(black))
    simpleBoard.setPiece(Point{6, 0}, newKnight(black))
    simpleBoard.setPiece(Point{7, 0}, newRook(black, false))

    simpleBoard.setPiece(Point{0, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{1, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{2, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{3, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{4, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{5, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{6, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{7, 1}, newPawn(black, false, 0, 1))

    simpleBoard.setPiece(Point{0, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{1, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{2, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{3, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{4, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{5, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{6, 6}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{7, 6}, newPawn(white, false, 0, -1))

    simpleBoard.setPiece(Point{0, 7}, newRook(white, false))
    simpleBoard.setPiece(Point{1, 7}, newKnight(white))
    simpleBoard.setPiece(Point{2, 7}, newBishop(white))
    simpleBoard.setPiece(Point{3, 7}, newQueen(white))
    simpleBoard.setPiece(Point{4, 7}, newKing(white, false, 0, -1))
    simpleBoard.setPiece(Point{5, 7}, newBishop(white))
    simpleBoard.setPiece(Point{6, 7}, newKnight(white))
    simpleBoard.setPiece(Point{7, 7}, newRook(white, false))

    err = simpleBoard.CalculateMoves()
    if err != nil {
        return nil, err
    }
    
    return simpleBoard, nil
}

func createSimpleFourPlayerBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    black := 2
    white := 0
    red := 1
    blue := 3

    simpleBoard, err := newSimpleBoard(Point{14, 14}, 4)
    if err != nil {
        return nil, err
    }

    simpleBoard.setPiece(Point{3, 0}, newRook(black, false))
    simpleBoard.setPiece(Point{4, 0}, newKnight(black))
    simpleBoard.setPiece(Point{5, 0}, newBishop(black))
    simpleBoard.setPiece(Point{6, 0}, newQueen(black))
    simpleBoard.setPiece(Point{7, 0}, newKing(black, false, 0, 1))
    simpleBoard.setPiece(Point{8, 0}, newBishop(black))
    simpleBoard.setPiece(Point{9, 0}, newKnight(black))
    simpleBoard.setPiece(Point{10, 0}, newRook(black, false))

    simpleBoard.setPiece(Point{3, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{4, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{5, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{6, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{7, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{8, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{9, 1}, newPawn(black, false, 0, 1))
    simpleBoard.setPiece(Point{10, 1}, newPawn(black, false, 0, 1))

    simpleBoard.setPiece(Point{3, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{4, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{5, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{6, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{7, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{8, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{9, 12}, newPawn(white, false, 0, -1))
    simpleBoard.setPiece(Point{10, 12}, newPawn(white, false, 0, -1))

    simpleBoard.setPiece(Point{3, 13}, newRook(white, false))
    simpleBoard.setPiece(Point{4, 13}, newKnight(white))
    simpleBoard.setPiece(Point{5, 13}, newBishop(white))
    simpleBoard.setPiece(Point{6, 13}, newQueen(white))
    simpleBoard.setPiece(Point{7, 13}, newKing(white, false, 0, -1))
    simpleBoard.setPiece(Point{8, 13}, newBishop(white))
    simpleBoard.setPiece(Point{9, 13}, newKnight(white))
    simpleBoard.setPiece(Point{10, 13}, newRook(white, false))

    simpleBoard.setPiece(Point{0, 3}, newRook(red, false))
    simpleBoard.setPiece(Point{0, 4}, newKnight(red))
    simpleBoard.setPiece(Point{0, 5}, newBishop(red))
    simpleBoard.setPiece(Point{0, 6}, newQueen(red))
    simpleBoard.setPiece(Point{0, 7}, newKing(red, false, 1, 0))
    simpleBoard.setPiece(Point{0, 8}, newBishop(red))
    simpleBoard.setPiece(Point{0, 9}, newKnight(red))
    simpleBoard.setPiece(Point{0, 10}, newRook(red, false))

    simpleBoard.setPiece(Point{1, 3}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 4}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 5}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 6}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 7}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 8}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 9}, newPawn(red, false, 1, 0))
    simpleBoard.setPiece(Point{1, 10}, newPawn(red, false, 1, 0))

    simpleBoard.setPiece(Point{12, 3}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 4}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 5}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 6}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 7}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 8}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 9}, newPawn(blue, false, -1, 0))
    simpleBoard.setPiece(Point{12, 10}, newPawn(blue, false, -1, 0))

    simpleBoard.setPiece(Point{13, 3}, newRook(blue, false))
    simpleBoard.setPiece(Point{13, 4}, newKnight(blue))
    simpleBoard.setPiece(Point{13, 5}, newBishop(blue))
    simpleBoard.setPiece(Point{13, 6}, newQueen(blue))
    simpleBoard.setPiece(Point{13, 7}, newKing(blue, false, -1, 0))
    simpleBoard.setPiece(Point{13, 8}, newBishop(blue))
    simpleBoard.setPiece(Point{13, 9}, newKnight(blue))
    simpleBoard.setPiece(Point{13, 10}, newRook(blue, false))

    simpleBoard.disableLocation(Point{0, 0})
    simpleBoard.disableLocation(Point{0, 1})
    simpleBoard.disableLocation(Point{0, 2})
    simpleBoard.disableLocation(Point{1, 0})
    simpleBoard.disableLocation(Point{1, 1})
    simpleBoard.disableLocation(Point{1, 2})
    simpleBoard.disableLocation(Point{2, 0})
    simpleBoard.disableLocation(Point{2, 1})
    simpleBoard.disableLocation(Point{2, 2})

    simpleBoard.disableLocation(Point{0, 11})
    simpleBoard.disableLocation(Point{0, 12})
    simpleBoard.disableLocation(Point{0, 13})
    simpleBoard.disableLocation(Point{1, 11})
    simpleBoard.disableLocation(Point{1, 12})
    simpleBoard.disableLocation(Point{1, 13})
    simpleBoard.disableLocation(Point{2, 11})
    simpleBoard.disableLocation(Point{2, 12})
    simpleBoard.disableLocation(Point{2, 13})

    simpleBoard.disableLocation(Point{11, 0})
    simpleBoard.disableLocation(Point{11, 1})
    simpleBoard.disableLocation(Point{11, 2})
    simpleBoard.disableLocation(Point{12, 0})
    simpleBoard.disableLocation(Point{12, 1})
    simpleBoard.disableLocation(Point{12, 2})
    simpleBoard.disableLocation(Point{13, 0})
    simpleBoard.disableLocation(Point{13, 1})
    simpleBoard.disableLocation(Point{13, 2})

    simpleBoard.disableLocation(Point{11, 11})
    simpleBoard.disableLocation(Point{11, 12})
    simpleBoard.disableLocation(Point{11, 13})
    simpleBoard.disableLocation(Point{12, 11})
    simpleBoard.disableLocation(Point{12, 12})
    simpleBoard.disableLocation(Point{12, 13})
    simpleBoard.disableLocation(Point{13, 11})
    simpleBoard.disableLocation(Point{13, 12})
    simpleBoard.disableLocation(Point{13, 13})

    err = simpleBoard.CalculateMoves()
    if err != nil {
        return nil, err
    }

    return simpleBoard, nil
}

func createSimplePlayerCollectionWithDefaultPlayers() (*SimplePlayerCollection, error) {
    return newSimplePlayerCollection(2)
}

func createSimpleFourPlayerPlayerCollectionWithDefaultPlayers() (*SimplePlayerCollection, error) {
    return newSimplePlayerCollection(4)
}

