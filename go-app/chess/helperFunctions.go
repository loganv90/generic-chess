package chess

func createSimpleBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    black := 1
    white := 0

    simpleBoard, err := newSimpleBoard(8, 8, 2)
    if err != nil {
        return nil, err
    }

    simpleBoard.setPiece(simpleBoard.getIndex(0, 0), simpleBoard.getAllPiece(black, ROOK))
    simpleBoard.setPiece(simpleBoard.getIndex(1, 0), simpleBoard.getAllPiece(black, KNIGHT))
    simpleBoard.setPiece(simpleBoard.getIndex(2, 0), simpleBoard.getAllPiece(black, BISHOP))
    simpleBoard.setPiece(simpleBoard.getIndex(3, 0), simpleBoard.getAllPiece(black, QUEEN))
    simpleBoard.setPiece(simpleBoard.getIndex(4, 0), simpleBoard.getAllPiece(black, KING_D))
    simpleBoard.setPiece(simpleBoard.getIndex(5, 0), simpleBoard.getAllPiece(black, BISHOP))
    simpleBoard.setPiece(simpleBoard.getIndex(6, 0), simpleBoard.getAllPiece(black, KNIGHT))
    simpleBoard.setPiece(simpleBoard.getIndex(7, 0), simpleBoard.getAllPiece(black, ROOK))

    simpleBoard.setPiece(simpleBoard.getIndex(0, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(1, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(2, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(3, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(4, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(5, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(6, 1), simpleBoard.getAllPiece(black, PAWN_D))
    simpleBoard.setPiece(simpleBoard.getIndex(7, 1), simpleBoard.getAllPiece(black, PAWN_D))

    simpleBoard.setPiece(simpleBoard.getIndex(0, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(1, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(2, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(3, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(4, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(5, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(6, 6), simpleBoard.getAllPiece(white, PAWN_U))
    simpleBoard.setPiece(simpleBoard.getIndex(7, 6), simpleBoard.getAllPiece(white, PAWN_U))

    simpleBoard.setPiece(simpleBoard.getIndex(0, 7), simpleBoard.getAllPiece(white, ROOK))
    simpleBoard.setPiece(simpleBoard.getIndex(1, 7), simpleBoard.getAllPiece(white, KNIGHT))
    simpleBoard.setPiece(simpleBoard.getIndex(2, 7), simpleBoard.getAllPiece(white, BISHOP))
    simpleBoard.setPiece(simpleBoard.getIndex(3, 7), simpleBoard.getAllPiece(white, QUEEN))
    simpleBoard.setPiece(simpleBoard.getIndex(4, 7), simpleBoard.getAllPiece(white, KING_U))
    simpleBoard.setPiece(simpleBoard.getIndex(5, 7), simpleBoard.getAllPiece(white, BISHOP))
    simpleBoard.setPiece(simpleBoard.getIndex(6, 7), simpleBoard.getAllPiece(white, KNIGHT))
    simpleBoard.setPiece(simpleBoard.getIndex(7, 7), simpleBoard.getAllPiece(white, ROOK))

    simpleBoard.CalculateMoves()
    
    return simpleBoard, nil
}

func createSimpleFourPlayerBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    black := 2
    white := 0
    red := 1
    blue := 3

    simpleBoard, err := newSimpleBoard(14, 14, 4)
    if err != nil {
        return nil, err
    }

    simpleBoard.setPiece(simpleBoard.getIndex(3, 0), simpleBoard.getAllPiece(black, ROOK)) 
    simpleBoard.setPiece(simpleBoard.getIndex(4, 0), simpleBoard.getAllPiece(black, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(5, 0), simpleBoard.getAllPiece(black, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(6, 0), simpleBoard.getAllPiece(black, QUEEN)) 
    simpleBoard.setPiece(simpleBoard.getIndex(7, 0), simpleBoard.getAllPiece(black, KING_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(8, 0), simpleBoard.getAllPiece(black, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(9, 0), simpleBoard.getAllPiece(black, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(10, 0), simpleBoard.getAllPiece(black, ROOK)) 

    simpleBoard.setPiece(simpleBoard.getIndex(3, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(4, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(5, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(6, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(7, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(8, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(9, 1), simpleBoard.getAllPiece(black, PAWN_D)) 
    simpleBoard.setPiece(simpleBoard.getIndex(10, 1), simpleBoard.getAllPiece(black, PAWN_D)) 

    simpleBoard.setPiece(simpleBoard.getIndex(3, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(4, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(5, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(6, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(7, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(8, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(9, 12), simpleBoard.getAllPiece(white, PAWN_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(10, 12), simpleBoard.getAllPiece(white, PAWN_U)) 

    simpleBoard.setPiece(simpleBoard.getIndex(3, 13), simpleBoard.getAllPiece(white, ROOK)) 
    simpleBoard.setPiece(simpleBoard.getIndex(4, 13), simpleBoard.getAllPiece(white, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(5, 13), simpleBoard.getAllPiece(white, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(6, 13), simpleBoard.getAllPiece(white, QUEEN)) 
    simpleBoard.setPiece(simpleBoard.getIndex(7, 13), simpleBoard.getAllPiece(white, KING_U)) 
    simpleBoard.setPiece(simpleBoard.getIndex(8, 13), simpleBoard.getAllPiece(white, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(9, 13), simpleBoard.getAllPiece(white, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(10, 13), simpleBoard.getAllPiece(white, ROOK)) 

    simpleBoard.setPiece(simpleBoard.getIndex(0, 3), simpleBoard.getAllPiece(red, ROOK)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 4), simpleBoard.getAllPiece(red, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 5), simpleBoard.getAllPiece(red, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 6), simpleBoard.getAllPiece(red, QUEEN)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 7), simpleBoard.getAllPiece(red, KING_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 8), simpleBoard.getAllPiece(red, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 9), simpleBoard.getAllPiece(red, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(0, 10), simpleBoard.getAllPiece(red, ROOK)) 

    simpleBoard.setPiece(simpleBoard.getIndex(1, 3), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 4), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 5), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 6), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 7), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 8), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 9), simpleBoard.getAllPiece(red, PAWN_R)) 
    simpleBoard.setPiece(simpleBoard.getIndex(1, 10), simpleBoard.getAllPiece(red, PAWN_R)) 

    simpleBoard.setPiece(simpleBoard.getIndex(12, 3), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 4), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 5), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 6), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 7), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 8), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 9), simpleBoard.getAllPiece(blue, PAWN_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(12, 10), simpleBoard.getAllPiece(blue, PAWN_L)) 

    simpleBoard.setPiece(simpleBoard.getIndex(13, 3), simpleBoard.getAllPiece(blue, ROOK)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 4), simpleBoard.getAllPiece(blue, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 5), simpleBoard.getAllPiece(blue, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 6), simpleBoard.getAllPiece(blue, QUEEN)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 7), simpleBoard.getAllPiece(blue, KING_L)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 8), simpleBoard.getAllPiece(blue, BISHOP)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 9), simpleBoard.getAllPiece(blue, KNIGHT)) 
    simpleBoard.setPiece(simpleBoard.getIndex(13, 10), simpleBoard.getAllPiece(blue, ROOK)) 

    simpleBoard.disableLocation(&Point{0, 0})
    simpleBoard.disableLocation(&Point{0, 1})
    simpleBoard.disableLocation(&Point{0, 2})
    simpleBoard.disableLocation(&Point{1, 0})
    simpleBoard.disableLocation(&Point{1, 1})
    simpleBoard.disableLocation(&Point{1, 2})
    simpleBoard.disableLocation(&Point{2, 0})
    simpleBoard.disableLocation(&Point{2, 1})
    simpleBoard.disableLocation(&Point{2, 2})

    simpleBoard.disableLocation(&Point{0, 11})
    simpleBoard.disableLocation(&Point{0, 12})
    simpleBoard.disableLocation(&Point{0, 13})
    simpleBoard.disableLocation(&Point{1, 11})
    simpleBoard.disableLocation(&Point{1, 12})
    simpleBoard.disableLocation(&Point{1, 13})
    simpleBoard.disableLocation(&Point{2, 11})
    simpleBoard.disableLocation(&Point{2, 12})
    simpleBoard.disableLocation(&Point{2, 13})

    simpleBoard.disableLocation(&Point{11, 0})
    simpleBoard.disableLocation(&Point{11, 1})
    simpleBoard.disableLocation(&Point{11, 2})
    simpleBoard.disableLocation(&Point{12, 0})
    simpleBoard.disableLocation(&Point{12, 1})
    simpleBoard.disableLocation(&Point{12, 2})
    simpleBoard.disableLocation(&Point{13, 0})
    simpleBoard.disableLocation(&Point{13, 1})
    simpleBoard.disableLocation(&Point{13, 2})

    simpleBoard.disableLocation(&Point{11, 11})
    simpleBoard.disableLocation(&Point{11, 12})
    simpleBoard.disableLocation(&Point{11, 13})
    simpleBoard.disableLocation(&Point{12, 11})
    simpleBoard.disableLocation(&Point{12, 12})
    simpleBoard.disableLocation(&Point{12, 13})
    simpleBoard.disableLocation(&Point{13, 11})
    simpleBoard.disableLocation(&Point{13, 12})
    simpleBoard.disableLocation(&Point{13, 13})

    simpleBoard.CalculateMoves()

    return simpleBoard, nil
}

func createSimplePlayerCollectionWithDefaultPlayers() (*SimplePlayerCollection, error) {
    return newSimplePlayerCollection(2)
}

func createSimpleFourPlayerPlayerCollectionWithDefaultPlayers() (*SimplePlayerCollection, error) {
    return newSimplePlayerCollection(4)
}

