package chess

type BoardData struct {
    XSize int
    YSize int
    Pieces []*PieceData
    Turn string
    Check bool
    Checkmate bool
    Stalemate bool
}

type PieceData struct {
    T string // Type
    C string // Color
    X int // X position
    Y int // Y position
}

type PieceState struct {
    X int
    Y int
    Moves []*MoveData
    Turn bool
}

type MoveData struct {
    X int
    Y int
}

type Point struct {
    x int
    y int
}

func (p *Point) equals(other *Point) bool {
    return p.x == other.x && p.y == other.y
}

func (p *Point) add(other *Point) *Point {
    return &Point{p.x + other.x, p.y + other.y}
}

type EnPassant struct {
    target *Point
    pieceLocation *Point
}

type PieceLocations struct {
    ownPieceLocations []*Point
    enemyPieceLocations []*Point
}

