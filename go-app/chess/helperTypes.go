package chess

type BoardState struct {
    Squares [][]*SquareData
    Turn string
	Check bool
	Mate  bool
}

type SquareData struct {
    C string
    P string
}

type PieceState struct {
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

