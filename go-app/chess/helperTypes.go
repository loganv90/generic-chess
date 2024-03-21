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



type PieceLocation struct {
    piece *Piece
    location *Point
}

func (p *PieceLocation) set(other *PieceLocation) {
    p.piece = other.piece
    p.location = other.location
}



type EnPassant struct {
    target *Point
    risk *Point
}

func (e *EnPassant) set(other *EnPassant) {
    e.target = other.target
    e.risk = other.risk
}



type Vulnerable struct {
    start *Point
    end *Point
}

func (v *Vulnerable) set(other *Vulnerable) {
    v.start = other.start
    v.end = other.end
}



type Array4[T any] struct {
    array [4]T
    count int
}

func (a *Array4[T]) get() *T {
    return &a.array[a.count]
}

func (a *Array4[T]) set(value T) {
    a.array[a.count] = value
    a.count += 1
}

func (a *Array4[T]) next() {
    a.count += 1
}

func (a *Array4[T]) clear() {
    a.count = 0
}



type Array100[T any] struct {
    array [100]T
    count int
}

func (a *Array100[T]) get() *T {
    return &a.array[a.count]
}

func (a *Array100[T]) set(value T) {
    a.array[a.count] = value
    a.count += 1
}

func (a *Array100[T]) next() {
    a.count += 1
}

func (a *Array100[T]) clear() {
    a.count = 0
}



type Array1000[T any] struct {
    array [200]T
    count int
}

func (a *Array1000[T]) get() *T {
    return &a.array[a.count]
}

func (a *Array1000[T]) set(value T) {
    a.array[a.count] = value
    a.count += 1
}

func (a *Array1000[T]) next() {
    a.count += 1
}

func (a *Array1000[T]) clear() {
    a.count = 0
}

