package chess

type BoardData struct {
    XSize int
    YSize int
    Pieces []*PieceData
    Disabled []*DisabledData
    CurrentPlayer int
    WinningPlayer int
    GameOver bool
    Check bool
    Checkmate bool
    Stalemate bool
}

type Command struct {
    m Move
    p PlayerTransition
    fullMove bool
}

type MoveKeyAndScore struct {
    moveKey MoveKey
    score map[int]int
}

type PieceData struct {
    T string // Type
    C int // Color
    X int // X position
    Y int // Y position
    D bool // Disabled
}

type DisabledData struct {
    X int
    Y int
}

type PieceState struct {
    X int
    Y int
    Moves []*MoveData
    Turn bool
}

type MoveData struct {
    X int // X position
    Y int // Y position
    P bool // is promotion move
}

type MoveKey struct {
    XFrom int
    YFrom int
    XTo int
    YTo int
    Promotion string
}

type Player struct {
    color int
    alive bool
}

type Point struct {
    x int
    y int
}

func (p Point) equals(other Point) bool {
    return p.x == other.x && p.y == other.y
}

func (p Point) add(other Point) Point {
    return Point{p.x + other.x, p.y + other.y}
}

type EnPassant struct {
    target Point
    pieceLocation Point
}

type PieceLocations struct {
    ownPieceLocations []Point
    enemyPieceLocations []Point
}

