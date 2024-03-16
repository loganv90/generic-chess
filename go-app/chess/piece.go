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

var piece_move_functions = []func(Board, *Piece, *Point, *Array100[FastMove]) {
    func (b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {},
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

var pawn_u_directions = []*Point{
    {0, -1}, // one 
    {0, -2}, // two
    {0, -3}, // three
    {-1, -1}, // capture
    {1, -1}, // capture
}

var pawn_d_directions = []*Point{
    {0, 1},
    {0, 2},
    {0, 3},
    {-1, 1},
    {1, 1},
}

var pawn_l_directions = []*Point{
    {-1, 0},
    {-2, 0},
    {-3, 0},
    {-1, -1},
    {-1, 1},
}

var pawn_r_directions = []*Point{
    {1, 0},
    {2, 0},
    {3, 0},
    {1, -1},
    {1, 1},
}

var knight_directions = []*Point{
    {1, 2},
    {-1, 2},
    {2, 1},
    {-2, 1},
    {1, -2},
    {-1, -2},
    {2, -1},
    {-2, -1},
}

var bishop_directions = []*Point{
    {1, 1},
    {-1, 1},
    {1, -1},
    {-1, -1},
}

var rook_directions = []*Point{
    {1, 0},
    {-1, 0},
    {0, 1},
    {0, -1},
}

var queen_directions = []*Point{
    {1, 0},
    {-1, 0},
    {0, 1},
    {0, -1},
    {1, 1},
    {-1, 1},
    {1, -1},
    {-1, -1},
}

var king_ud_directions = []*Point{
    {1, 0}, // castle search 1
    {-1, 0}, // castle search 2
    {-1, 0}, // king offset 1
    {2, 0}, // king offset 2
    {-2, 0}, // rook offset 1
    {3, 0}, // rook offset 2
}

var king_lr_directions = []*Point{
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

func (p *Piece) copy(piece *Piece) {
    piece.color = p.color
    piece.index = piece_moved_indexes[p.index]
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

func (p *Piece) moves(b Board, fromLocation *Point, moves *Array100[FastMove]) {
    piece_move_functions[p.index](b, p, fromLocation, moves)
}

func addDirection(
	b Board,
    fromPiece *Piece,
    fromLocation *Point,
	moves *Array100[FastMove],
    direction *Point,
) {
    currentLocation := *fromLocation

    for {
        currentLocation.add(direction)

        currentPiece := b.getPiece(&currentLocation)
        if currentPiece == nil {
            break
        }

        if !currentPiece.valid() { // no piece
            addMoveSimple(b, fromPiece, fromLocation, currentPiece, &currentLocation, nil, moves)
        } else if currentPiece.color != fromPiece.color { // enemy piece
            addMoveSimple(b, fromPiece, fromLocation, currentPiece, &currentLocation, nil, moves)
            break
        } else { // ally piece
            addMoveAllyDefense(b, fromPiece, fromLocation, &currentLocation, moves)
            break
        }
    }
}

func addSimple(
	b Board,
    fromPiece *Piece,
    fromLocation *Point,
	moves *Array100[FastMove],
    direction *Point,
) {
    toLocation := *fromLocation
    toLocation.add(direction)

	toPiece := b.getPiece(&toLocation)
	if toPiece == nil {
		return
	}

	if !toPiece.valid() { // no piece
        addMoveSimple(b, fromPiece, fromLocation, toPiece, &toLocation, nil, moves)
	} else if toPiece.color != fromPiece.color { // enemy piece
        addMoveSimple(b, fromPiece, fromLocation, toPiece, &toLocation, nil, moves)
	} else { // ally piece
        addMoveAllyDefense(b, fromPiece, fromLocation, &toLocation, moves)
    }
}

var pawn_r_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_r_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_r_directions)
}

var pawn_l_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_l_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_l_directions)
}

var pawn_u_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_u_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_u_directions)
}

var pawn_d_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_d_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_d_directions)
}

func pawnAddForward(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove], directions []*Point) {
    to1Location := *fromLocation
    to2Location := *fromLocation
    to3Location := *fromLocation

    to1Location.add(directions[0])
    to2Location.add(directions[1])
    to3Location.add(directions[2])

    piece1 := b.getPiece(&to1Location)
    piece2 := b.getPiece(&to2Location)
    piece3 := b.getPiece(&to3Location)

    if piece1 == nil {
        return
    }

    queen := Piece{fromPiece.color, QUEEN}
    rook_m := Piece{fromPiece.color, ROOK_M}
    bishop := Piece{fromPiece.color, BISHOP}
    knight := Piece{fromPiece.color, KNIGHT}

    if !piece1.valid() {
        if piece2 == nil {
            addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &queen, moves)
            addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &rook_m, moves)
            addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &bishop, moves)
            addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &knight, moves)
        } else {
            addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, nil, moves)
        }
    }

    if piece2 == nil {
        return
    }

    if fromPiece.moved() {
        return
    }

    enPassant := EnPassant{to1Location, to2Location}

    if !piece2.valid() {
        if piece3 == nil {
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &queen, &enPassant, moves)
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &rook_m, &enPassant, moves)
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &bishop, &enPassant, moves)
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &knight, &enPassant, moves)
        } else {
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, nil, &enPassant, moves)
        }
    }
}

func pawnAddCaptures(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove], directions []*Point) {
    to1Location := *fromLocation
    to2Location := *fromLocation
    to3Location := to1Location
    to4Location := to2Location

    to1Location.add(directions[3])
    to2Location.add(directions[4])
    to3Location.add(directions[0])
    to4Location.add(directions[0])

    piece1 := b.getPiece(&to1Location)
    piece2 := b.getPiece(&to2Location)
    piece3 := b.getPiece(&to3Location)
    piece4 := b.getPiece(&to4Location)

    queen := Piece{fromPiece.color, QUEEN}
    rook_m := Piece{fromPiece.color, ROOK_M}
    bishop := Piece{fromPiece.color, BISHOP}
    knight := Piece{fromPiece.color, KNIGHT}

    if piece1 != nil {
        if en1, en2 := b.possibleEnPassant(fromPiece.color, &to1Location); en1 != nil { // if the square is an en passant target
            if piece3 == nil {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, &to1Location, &queen, moves, en1, en2)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, &to1Location, &rook_m, moves, en1, en2)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, &to1Location, &bishop, moves, en1, en2)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, &to1Location, &knight, moves, en1, en2)
            } else {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, &to1Location, nil, moves, en1, en2)
            }
        } else if piece1.valid() && piece1.color != fromPiece.color { // if the square is occupied by an enemy piece
            if piece3 == nil {
                addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &queen, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &rook_m, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &bishop, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, &knight, moves)
            } else {
                addMoveSimple(b, fromPiece, fromLocation, piece1, &to1Location, nil, moves)
            }
        } else if piece1.valid() {
            addMoveAllyDefense(b, fromPiece, fromLocation, &to1Location, moves)
        }
    }

    if piece2 != nil {
        if en1, en2 := b.possibleEnPassant(fromPiece.color, &to2Location); en1 != nil { // if the square is an en passant target
            if piece4 == nil {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &queen, moves, en1, en2)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &rook_m, moves, en1, en2)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &bishop, moves, en1, en2)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, &knight, moves, en1, en2)
            } else {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, &to2Location, nil, moves, en1, en2)
            }
        } else if piece2.valid() && piece2.color != fromPiece.color { // if the square is occupied by an enemy piece
            if piece4 == nil {
                addMoveSimple(b, fromPiece, fromLocation, piece2, &to2Location, &queen, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece2, &to2Location, &rook_m, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece2, &to2Location, &bishop, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece2, &to2Location, &knight, moves)
            } else {
                addMoveSimple(b, fromPiece, fromLocation, piece2, &to2Location, nil, moves)
            }
        } else if piece2.valid() {
            addMoveAllyDefense(b, fromPiece, fromLocation, &to2Location, moves)
        }
    }
}

func knight_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
	for _, direction := range knight_directions {
		addSimple(b, fromPiece, fromLocation, moves, direction)
	}
}

func bishop_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    for _, direction := range bishop_directions {
        addDirection(b, fromPiece, fromLocation, moves, direction)
    }
}

func rook_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    for _, direction := range rook_directions {
        addDirection(b, fromPiece, fromLocation, moves, direction)
    }
}

func queen_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    for _, direction := range queen_directions {
        addDirection(b, fromPiece, fromLocation, moves, direction)
    }
}

var king_lr_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    for _, direction := range queen_directions {
        addSimple(b, fromPiece, fromLocation, moves, direction)
    }

    if fromPiece.moved() {
        return
    }

    addCastle(b, fromPiece, fromLocation, moves, king_lr_directions[0], king_lr_directions[2], king_lr_directions[4])
    addCastle(b, fromPiece, fromLocation, moves, king_lr_directions[1], king_lr_directions[3], king_lr_directions[5])
}

var king_ud_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove]) {
    for _, direction := range queen_directions {
        addSimple(b, fromPiece, fromLocation, moves, direction)
    }

    if fromPiece.moved() {
        return
    }

    addCastle(b, fromPiece, fromLocation, moves, king_ud_directions[0], king_ud_directions[2], king_ud_directions[4])
    addCastle(b, fromPiece, fromLocation, moves, king_ud_directions[1], king_ud_directions[3], king_ud_directions[5])
}

func addCastle(b Board, fromPiece *Piece, fromLocation *Point, moves *Array100[FastMove], direction *Point, kingOffset *Point, rookOffset *Point) {
    // find rook for castle
    fromRookLocation := *fromLocation
    var rook *Piece

    for {
        fromRookLocation.add(direction)

        rook = b.getPiece(&fromRookLocation)
        if rook == nil {
            return
        }

        if !rook.valid() {
            continue
        }

        if rook.index == ROOK && rook.color == fromPiece.color {
            break
        }

        return
    }

    // find edge for castle
    edgeLocation := fromRookLocation
    currentLocation := edgeLocation

    for {
        currentLocation.add(direction)

        piece := b.getPiece(&currentLocation)
        if piece == nil {
            break
        }

        edgeLocation.add(direction)
    }

    // everything else
    toLocation := edgeLocation
    toRookLocation := edgeLocation
    toLocation.add(kingOffset)
    toRookLocation.add(rookOffset)

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
        if piece := b.getPiece(&Point{x, fromLocation.y}); piece == nil || piece.valid() {
            clr = false
            break
        }
    }
    for y := yCheckedMin - 1; y >= yToMin && clr; y-- {
        if piece := b.getPiece(&Point{fromLocation.x, y}); piece == nil || piece.valid() {
            clr = false
            break
        }
    }
    for x := xCheckedMax + 1; x <= xToMax && clr; x++ {
        if piece := b.getPiece(&Point{x, fromLocation.y}); piece == nil || piece.valid() {
            clr = false
            break
        }
    }
    for y := yCheckedMax + 1; y <= yToMax && clr; y++ {
        if piece := b.getPiece(&Point{fromLocation.y, y}); piece == nil || piece.valid() {
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
    vulnerable := Vulnerable{Point{minx, miny}, Point{maxx, maxy}}

    addMoveCastle(b, fromPiece, fromLocation, &toLocation, rook, &fromRookLocation, &toRookLocation, &vulnerable, moves)
}

