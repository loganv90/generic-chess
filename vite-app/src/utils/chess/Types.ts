type PieceMove = {
    x: number,
    y: number,
    options: PieceMoveOptions,
}

type PieceMoveOptions = {
    direction?: true,
    noCapture?: true,
    mustCapture?: true,
    canPromote?: true,
    canCaptureEnPassant?: true,
    mustCastle?: true,
    mustRevealEnPassant?: {x: number, y: number},
}

type BoardMove = {
    xFrom: number,
    yFrom: number,
    xTo: number,
    yTo: number,
    options: BoardMoveOptions,
}

type BoardMoveOptions = {
    promote?: true,
    captureEnPassant?: true,
    revealEnPassant?: {x: number, y: number},
}

type EnPassantMap = {
    [key: string]: EnPassant | null,
}

type EnPassant = {
    x: number,
    y: number,
    xPiece: number,
    yPiece: number,
}

export type { PieceMove, PieceMoveOptions, BoardMove, BoardMoveOptions, EnPassantMap, EnPassant }