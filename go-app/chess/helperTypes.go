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
    m FastMove
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

type PieceLocation struct {
    piece Piece
    location Point
}

type EnPassant struct {
    target Point
    risk Point
}

type Vulnerable struct {
    start Point
    end Point
}

type Array4[T any] struct {
    array [4]T
    count int
}

func (a *Array4[T]) append(item T) {
    if a.count >= len(a.array) {
        a.clear()
    }

    a.array[a.count] = item
    a.count += 1
}

func (a *Array4[T]) clear() {
    a.count = 0
}

type Array100[T any] struct {
    array [100]T
    count int
}

func (a *Array100[T]) append(item T) {
    if a.count >= len(a.array) {
        a.clear()
    }

    a.array[a.count] = item
    a.count += 1
}

func (a *Array100[T]) clear() {
    a.count = 0
}

