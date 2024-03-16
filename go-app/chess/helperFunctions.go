package chess

func setPiece(board Board, point Point, piece Piece) {
    p := board.getPiece(&point)
    if p == nil {
        panic("what")
    }

    p.color = piece.color
    p.index = piece.index
}

func setEnPassant(board Board, color int, enPassant EnPassant) {
    e := board.getEnPassant(color)
    if e == nil {
        panic("what")
    }

    e.target = enPassant.target
    e.risk = enPassant.risk
}

func createSimpleBoardWithDefaultPieceLocations() (*SimpleBoard, error) {
    black := 1
    white := 0

    simpleBoard, err := newSimpleBoard(Point{8, 8}, 2)
    if err != nil {
        return nil, err
    }

    setPiece(simpleBoard, Point{0, 0}, Piece{black, ROOK})
    setPiece(simpleBoard, Point{1, 0}, Piece{black, KNIGHT})
    setPiece(simpleBoard, Point{2, 0}, Piece{black, BISHOP})
    setPiece(simpleBoard, Point{3, 0}, Piece{black, QUEEN})
    setPiece(simpleBoard, Point{4, 0}, Piece{black, KING_D})
    setPiece(simpleBoard, Point{5, 0}, Piece{black, BISHOP})
    setPiece(simpleBoard, Point{6, 0}, Piece{black, KNIGHT})
    setPiece(simpleBoard, Point{7, 0}, Piece{black, ROOK})

    setPiece(simpleBoard, Point{0, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{1, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{2, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{3, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{4, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{5, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{6, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{7, 1}, Piece{black, PAWN_D})

    setPiece(simpleBoard, Point{0, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{1, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{2, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{3, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{4, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{5, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{6, 6}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{7, 6}, Piece{white, PAWN_U})

    setPiece(simpleBoard, Point{0, 7}, Piece{white, ROOK})
    setPiece(simpleBoard, Point{1, 7}, Piece{white, KNIGHT})
    setPiece(simpleBoard, Point{2, 7}, Piece{white, BISHOP})
    setPiece(simpleBoard, Point{3, 7}, Piece{white, QUEEN})
    setPiece(simpleBoard, Point{4, 7}, Piece{white, KING_U})
    setPiece(simpleBoard, Point{5, 7}, Piece{white, BISHOP})
    setPiece(simpleBoard, Point{6, 7}, Piece{white, KNIGHT})
    setPiece(simpleBoard, Point{7, 7}, Piece{white, ROOK})

    simpleBoard.CalculateMoves()
    
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

    setPiece(simpleBoard, Point{3, 0}, Piece{black, ROOK})
    setPiece(simpleBoard, Point{4, 0}, Piece{black, KNIGHT})
    setPiece(simpleBoard, Point{5, 0}, Piece{black, BISHOP})
    setPiece(simpleBoard, Point{6, 0}, Piece{black, QUEEN})
    setPiece(simpleBoard, Point{7, 0}, Piece{black, KING_D})
    setPiece(simpleBoard, Point{8, 0}, Piece{black, BISHOP})
    setPiece(simpleBoard, Point{9, 0}, Piece{black, KNIGHT})
    setPiece(simpleBoard, Point{10, 0}, Piece{black, ROOK})

    setPiece(simpleBoard, Point{3, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{4, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{5, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{6, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{7, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{8, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{9, 1}, Piece{black, PAWN_D})
    setPiece(simpleBoard, Point{10, 1}, Piece{black, PAWN_D})

    setPiece(simpleBoard, Point{3, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{4, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{5, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{6, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{7, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{8, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{9, 12}, Piece{white, PAWN_U})
    setPiece(simpleBoard, Point{10, 12}, Piece{white, PAWN_U})

    setPiece(simpleBoard, Point{3, 13}, Piece{white, ROOK})
    setPiece(simpleBoard, Point{4, 13}, Piece{white, KNIGHT})
    setPiece(simpleBoard, Point{5, 13}, Piece{white, BISHOP})
    setPiece(simpleBoard, Point{6, 13}, Piece{white, QUEEN})
    setPiece(simpleBoard, Point{7, 13}, Piece{white, KING_U})
    setPiece(simpleBoard, Point{8, 13}, Piece{white, BISHOP})
    setPiece(simpleBoard, Point{9, 13}, Piece{white, KNIGHT})
    setPiece(simpleBoard, Point{10, 13}, Piece{white, ROOK})

    setPiece(simpleBoard, Point{0, 3}, Piece{red, ROOK})
    setPiece(simpleBoard, Point{0, 4}, Piece{red, KNIGHT})
    setPiece(simpleBoard, Point{0, 5}, Piece{red, BISHOP})
    setPiece(simpleBoard, Point{0, 6}, Piece{red, QUEEN})
    setPiece(simpleBoard, Point{0, 7}, Piece{red, KING_R})
    setPiece(simpleBoard, Point{0, 8}, Piece{red, BISHOP})
    setPiece(simpleBoard, Point{0, 9}, Piece{red, KNIGHT})
    setPiece(simpleBoard, Point{0, 10}, Piece{red, ROOK})

    setPiece(simpleBoard, Point{1, 3}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 4}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 5}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 6}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 7}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 8}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 9}, Piece{red, PAWN_R})
    setPiece(simpleBoard, Point{1, 10}, Piece{red, PAWN_R})

    setPiece(simpleBoard, Point{12, 3}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 4}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 5}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 6}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 7}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 8}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 9}, Piece{blue, PAWN_L})
    setPiece(simpleBoard, Point{12, 10}, Piece{blue, PAWN_L})

    setPiece(simpleBoard, Point{13, 3}, Piece{blue, ROOK})
    setPiece(simpleBoard, Point{13, 4}, Piece{blue, KNIGHT})
    setPiece(simpleBoard, Point{13, 5}, Piece{blue, BISHOP})
    setPiece(simpleBoard, Point{13, 6}, Piece{blue, QUEEN})
    setPiece(simpleBoard, Point{13, 7}, Piece{blue, KING_L})
    setPiece(simpleBoard, Point{13, 8}, Piece{blue, BISHOP})
    setPiece(simpleBoard, Point{13, 9}, Piece{blue, KNIGHT})
    setPiece(simpleBoard, Point{13, 10}, Piece{blue, ROOK})

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

