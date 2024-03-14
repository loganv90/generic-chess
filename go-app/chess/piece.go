package chess

const (
    EMPTY = 0
    PAWN_R = 1
    PAWN_L = 2
    PAWN_D = 3
    PAWN_U = 4
    PAWN_R_M = 5
    PAWN_L_M = 6
    PAWN_D_M = 7
    PAWN_U_M = 8
    KNIGHT = 9
    BISHOP = 10
    ROOK = 11
    ROOK_M = 12
    QUEEN = 13
    KING_R = 14
    KING_L = 15
    KING_D = 16
    KING_U = 17
    KING_R_M = 18
    KING_L_M = 19
    KING_D_M = 20
    KING_U_M = 21
)

var piece_moved_indexes = []int{
    EMPTY,
    PAWN_R_M,
    PAWN_L_M,
    PAWN_D_M,
    PAWN_U_M,
    PAWN_R_M,
    PAWN_L_M,
    PAWN_D_M,
    PAWN_U_M,
    KNIGHT,
    BISHOP,
    ROOK,
    ROOK_M,
    QUEEN,
    KING_R_M,
    KING_L_M,
    KING_D_M,
    KING_U_M,
    KING_R_M,
    KING_L_M,
    KING_D_M,
    KING_U_M,
}

var piece_values = []int{
    0,
    100,
    100,
    100,
    100,
    100,
    100,
    100,
    100,
    300,
    300,
    500,
    500,
    900,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
}

var piece_names = []string{
    "",
    "P",
    "P",
    "P",
    "P",
    "P",
    "P",
    "P",
    "P",
    "N",
    "B",
    "R",
    "R",
    "Q",
    "K",
    "K",
    "K",
    "K",
    "K",
    "K",
    "K",
    "K",
}

var piece_move_functions = []func(Board, Piece, Point, *Array100[FastMove]) {
    func (b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {},
    pawn_r_moves,
    pawn_l_moves,
    pawn_d_moves,
    pawn_u_moves,
    pawn_r_moves,
    pawn_l_moves,
    pawn_d_moves,
    pawn_u_moves,
    knight_moves,
    bishop_moves,
    rook_moves,
    rook_moves,
    queen_moves,
    king_lr_moves,
    king_lr_moves,
    king_ud_moves,
    king_ud_moves,
    king_lr_moves,
    king_lr_moves,
    king_ud_moves,
    king_ud_moves,
}

var pawn_u_directions = []Point{
    {0, -1}, // one 
    {0, -2}, // two
    {0, -3}, // three
    {-1, -1}, // capture
    {1, -1}, // capture
}

var pawn_d_directions = []Point{
    {0, 1},
    {0, 2},
    {0, 3},
    {-1, 1},
    {1, 1},
}

var pawn_l_directions = []Point{
    {-1, 0},
    {-2, 0},
    {-3, 0},
    {-1, -1},
    {-1, 1},
}

var pawn_r_directions = []Point{
    {1, 0},
    {2, 0},
    {3, 0},
    {1, -1},
    {1, 1},
}

var knight_directions = []Point{
    {1, 2},
    {-1, 2},
    {2, 1},
    {-2, 1},
    {1, -2},
    {-1, -2},
    {2, -1},
    {-2, -1},
}

var bishop_directions = []Point{
    {1, 1},
    {-1, 1},
    {1, -1},
    {-1, -1},
}

var rook_directions = []Point{
    {1, 0},
    {-1, 0},
    {0, 1},
    {0, -1},
}

var queen_directions = []Point{
    {1, 0},
    {-1, 0},
    {0, 1},
    {0, -1},
    {1, 1},
    {-1, 1},
    {1, -1},
    {-1, -1},
}

var king_ud_directions = []Point{
    {1, 0}, // castle search 1
    {-1, 0}, // castle search 2
    {-1, 0}, // king offset 1
    {2, 0}, // king offset 2
    {-2, 0}, // rook offset 1
    {3, 0}, // rook offset 2
}

var king_lr_directions = []Point{
    {0, 1},
    {0, -1},
    {0, -1},
    {0, 2},
    {0, -2},
    {0, 3},
}

type Piece struct {
    color int
    index int
}

func (p *Piece) valid() bool {
    return p.index > 0
}

func (p *Piece) value() int {
    return piece_values[p.index]
}

func (p *Piece) print() string {
    return piece_names[p.index]
}

func (p *Piece) copy() Piece {
    return Piece{p.color, piece_moved_indexes[p.index]}
}

func (p *Piece) moved() bool {
    if p.index > 4 && p.index < 9 {
        return true
    }
    if p.index == 12 {
        return true
    }
    if p.index > 17 {
        return true
    }
    return false
}

func (p *Piece) moves(b Board, fromLocation Point, moves *Array100[FastMove]) {
    piece_move_functions[p.index](b, *p, fromLocation, moves)
}

func addDirection(
	b Board,
    p Piece,
    fromLocation Point,
	moves *Array100[FastMove],
    direction Point,
) {
    currentLocation := fromLocation.add(direction)

    for {
        piece, ok := b.getPiece(currentLocation)
        if !ok {
            break
        }

        if !piece.valid() { // no piece
            moves.append(createMoveSimple(b, p, fromLocation, piece, currentLocation, Piece{0, 0}))
        } else if piece.color != p.color { // enemy piece
            moves.append(createMoveSimple(b, p, fromLocation, piece, currentLocation, Piece{0, 0}))
            break
        } else { // ally piece
            moves.append(createMoveAllyDefense(b, p, fromLocation, currentLocation))
            break
        }

        currentLocation = currentLocation.add(direction)
    }
}

func addSimple(
	b Board,
    p Piece,
    fromLocation Point,
	moves *Array100[FastMove],
    direction Point,
) {
    toLocation := fromLocation.add(direction)

	piece, ok := b.getPiece(toLocation)
	if !ok {
		return
	}

	if !piece.valid() { // no piece
        moves.append(createMoveSimple(b, p, fromLocation, piece, toLocation, Piece{0, 0}))
	} else if piece.color != p.color { // enemy piece
        moves.append(createMoveSimple(b, p, fromLocation, piece, toLocation, Piece{0, 0}))
	} else { // ally piece
        moves.append(createMoveAllyDefense(b, p, fromLocation, toLocation))
    }
}

var pawn_r_moves = func(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    pawnAddForward(b, p, fromLocation, moves, pawn_r_directions)
    pawnAddCaptures(b, p, fromLocation, moves, pawn_r_directions)
}

var pawn_l_moves = func(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    pawnAddForward(b, p, fromLocation, moves, pawn_l_directions)
    pawnAddCaptures(b, p, fromLocation, moves, pawn_l_directions)
}

var pawn_u_moves = func(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    pawnAddForward(b, p, fromLocation, moves, pawn_u_directions)
    pawnAddCaptures(b, p, fromLocation, moves, pawn_u_directions)
}

var pawn_d_moves = func(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    pawnAddForward(b, p, fromLocation, moves, pawn_d_directions)
    pawnAddCaptures(b, p, fromLocation, moves, pawn_d_directions)
}

func pawnAddForward(b Board, p Piece, fromLocation Point, moves *Array100[FastMove], directions []Point) {
    to1Location := fromLocation.add(directions[0])
    to2Location := fromLocation.add(directions[1])
    to3Location := fromLocation.add(directions[2])
    piece1, ok1 := b.getPiece(to1Location)
    piece2, ok2 := b.getPiece(to2Location)
    _, ok3 := b.getPiece(to3Location)

    if !ok1 {
        return
    }

    if !piece1.valid() {
        if !ok2 {
            moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, QUEEN}))
            moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, ROOK_M}))
            moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, BISHOP}))
            moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, KNIGHT}))
        } else {
            moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{0, 0}))
        }
    }

    if !ok2 {
        return
    }

    if p.index > 4 {
        return
    }

    if !piece2.valid() {
        if !ok3 {
            moves.append(createMoveRevealEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, QUEEN}, EnPassant{to1Location, to2Location}))
            moves.append(createMoveRevealEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, ROOK_M}, EnPassant{to1Location, to2Location}))
            moves.append(createMoveRevealEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, BISHOP}, EnPassant{to1Location, to2Location}))
            moves.append(createMoveRevealEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, KNIGHT}, EnPassant{to1Location, to2Location}))
        } else {
            moves.append(createMoveRevealEnPassant(b, p, fromLocation, piece2, to2Location, Piece{0, 0}, EnPassant{to1Location, to2Location}))
        }
    }
}

func pawnAddCaptures(b Board, p Piece, fromLocation Point, moves *Array100[FastMove], directions []Point) {
    to1Location := fromLocation.add(directions[3])
    to2Location := fromLocation.add(directions[4])
    to3Location := to1Location.add(directions[0])
    to4Location := to2Location.add(directions[0])
    piece1, ok1 := b.getPiece(to1Location)
    piece2, ok2 := b.getPiece(to2Location)
    _, ok3 := b.getPiece(to3Location)
    _, ok4 := b.getPiece(to4Location)

    if ok1 {
        if ens, err := b.possibleEnPassant(p.color, to1Location); err == nil && len(ens) > 0 { // if the square is an en passant target
            if !ok3 {
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece1, to1Location, Piece{p.color, QUEEN}, ens))
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece1, to1Location, Piece{p.color, ROOK_M}, ens))
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece1, to1Location, Piece{p.color, BISHOP}, ens))
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece1, to1Location, Piece{p.color, KNIGHT}, ens))
            } else {
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece1, to1Location, Piece{0, 0}, ens))
            }
        } else if piece1.valid() && piece1.color != p.color { // if the square is occupied by an enemy piece
            if !ok3 {
                moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, QUEEN}))
                moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, ROOK_M}))
                moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, BISHOP}))
                moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{p.color, KNIGHT}))
            } else {
                moves.append(createMoveSimple(b, p, fromLocation, piece1, to1Location, Piece{0, 0}))
            }
        } else if piece1.valid() {
            moves.append(createMoveAllyDefense(b, p, fromLocation, to1Location))
        }
    }

    if ok2 {
        if ens, err := b.possibleEnPassant(p.color, to2Location); err == nil && len(ens) > 0 { // if the square is an en passant target
            if !ok4 {
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, QUEEN}, ens))
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, ROOK_M}, ens))
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, BISHOP}, ens))
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece2, to2Location, Piece{p.color, KNIGHT}, ens))
            } else {
                moves.append(createMoveCaptureEnPassant(b, p, fromLocation, piece2, to2Location, Piece{0, 0}, ens))
            }
        } else if piece2.valid() && piece2.color != p.color { // if the square is occupied by an enemy piece
            if !ok4 {
                moves.append(createMoveSimple(b, p, fromLocation, piece2, to2Location, Piece{p.color, QUEEN}))
                moves.append(createMoveSimple(b, p, fromLocation, piece2, to2Location, Piece{p.color, ROOK_M}))
                moves.append(createMoveSimple(b, p, fromLocation, piece2, to2Location, Piece{p.color, BISHOP}))
                moves.append(createMoveSimple(b, p, fromLocation, piece2, to2Location, Piece{p.color, KNIGHT}))
            } else {
                moves.append(createMoveSimple(b, p, fromLocation, piece2, to2Location, Piece{0, 0}))
            }
        } else if piece2.valid() {
            moves.append(createMoveAllyDefense(b, p, fromLocation, to2Location))
        }
    }
}

func knight_moves(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
	for _, direction := range knight_directions {
		addSimple(b, p, fromLocation, moves, direction)
	}
}

func bishop_moves(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    for _, direction := range bishop_directions {
        addDirection(b, p, fromLocation, moves, direction)
    }
}

func rook_moves(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    for _, direction := range rook_directions {
        addDirection(b, p, fromLocation, moves, direction)
    }
}

func queen_moves(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    for _, direction := range queen_directions {
        addDirection(b, p, fromLocation, moves, direction)
    }
}

var king_lr_moves = func(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    for _, direction := range queen_directions {
        addSimple(b, p, fromLocation, moves, direction)
    }

    if p.index > 17 {
        return
    }

    addCastle(b, p, fromLocation, moves, king_lr_directions[0], king_lr_directions[2], king_lr_directions[4])
    addCastle(b, p, fromLocation, moves, king_lr_directions[1], king_lr_directions[3], king_lr_directions[5])
}

var king_ud_moves = func(b Board, p Piece, fromLocation Point, moves *Array100[FastMove]) {
    for _, direction := range queen_directions {
        addSimple(b, p, fromLocation, moves, direction)
    }

    if p.index > 17 {
        return
    }

    addCastle(b, p, fromLocation, moves, king_ud_directions[0], king_ud_directions[2], king_ud_directions[4])
    addCastle(b, p, fromLocation, moves, king_ud_directions[1], king_ud_directions[3], king_ud_directions[5])
}

func findRookForCastle(b Board, p Piece, fromLocation Point, direction Point) (Point, Piece, bool) {
    currentLocation := fromLocation.add(direction)

    for {
        piece, ok := b.getPiece(currentLocation)
        if !ok {
            return Point{}, Piece{}, false
        }

        if !piece.valid() {
            currentLocation = currentLocation.add(direction)
            continue
        }

        if piece.index == ROOK && piece.color == p.color {
            return currentLocation, piece, true
        }

        return Point{}, Piece{}, false
    }
}

func findEdgeForCastle(b Board, p Piece, fromLocation Point, direction Point) Point {
    previousLocation := fromLocation
    currentLocation := previousLocation.add(direction)

    for {
        _, ok := b.getPiece(currentLocation)
        if !ok {
            return previousLocation
        }

        previousLocation = currentLocation
        currentLocation = previousLocation.add(direction)
    }
}

func addCastle(b Board, p Piece, fromLocation Point, moves *Array100[FastMove], direction Point, kingOffset Point, rookOffset Point) {
    fromRookLocation, rook, ok := findRookForCastle(b, p, fromLocation, direction)
    if !ok {
        return
    }

    edgeLocation := findEdgeForCastle(b, p, fromRookLocation, direction)
    toLocation := edgeLocation.add(kingOffset)
    toRookLocation := edgeLocation.add(rookOffset)

    xCheckedMin := min(fromLocation.x, fromRookLocation.x)
    xCheckedMax := max(fromLocation.x, fromRookLocation.x)
    yCheckedMin := min(fromLocation.y, fromRookLocation.y)
    yCheckedMax := max(fromLocation.y, fromRookLocation.y)

    xToMin := min(toLocation.x, toRookLocation.x)
    xToMax := max(toLocation.x, toRookLocation.x)
    yToMin := min(toLocation.y, toRookLocation.y)
    yToMax := max(toLocation.y, toRookLocation.y)

    clr := true
    for x := xCheckedMin - 1; x >= xToMin && clr; x-- {
        if piece, ok := b.getPiece(Point{x, fromLocation.y}); !ok || piece.valid() {
            clr = false
            break
        }
    }
    for y := yCheckedMin - 1; y >= yToMin && clr; y-- {
        if piece, ok := b.getPiece(Point{fromLocation.x, y}); !ok || piece.valid() {
            clr = false
            break
        }
    }
    for x := xCheckedMax + 1; x <= xToMax && clr; x++ {
        if piece, ok := b.getPiece(Point{x, fromLocation.y}); !ok || piece.valid() {
            clr = false
            break
        }
    }
    for y := yCheckedMax + 1; y <= yToMax && clr; y++ {
        if piece, ok := b.getPiece(Point{fromLocation.y, y}); !ok || piece.valid() {
            clr = false
            break
        }
    }
    if !clr {
        return
    }

    var minx int
    var maxx int
    var miny int
    var maxy int
    if toLocation.x > fromLocation.x {
        minx = fromLocation.x + 1
        maxx = toLocation.x - 1
        miny = fromLocation.y
        maxy = fromLocation.y
    } else if toLocation.x < fromLocation.x {
        minx = toLocation.x + 1
        maxx = fromLocation.x - 1
        miny = fromLocation.y
        maxy = fromLocation.y
    } else if toLocation.y > fromLocation.y {
        minx = fromLocation.x
        maxx = fromLocation.x
        miny = fromLocation.y + 1
        maxy = toLocation.y - 1
    } else if toLocation.y < fromLocation.y {
        minx = fromLocation.x
        maxx = fromLocation.x
        miny = toLocation.y + 1
        maxy = fromLocation.y - 1
    }

    moves.append(createMoveCastle(b, p, fromLocation, rook, fromRookLocation, toLocation, toRookLocation, Vulnerable{Point{minx, miny}, Point{maxx, maxy}}))
}

