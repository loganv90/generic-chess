package chess

const (
    PAWN_R = 0
    PAWN_L = 1
    PAWN_D = 2
    PAWN_U = 3
    PAWN_R_M = 4
    PAWN_L_M = 5
    PAWN_D_M = 6
    PAWN_U_M = 7
    KNIGHT = 8
    BISHOP = 9
    ROOK = 10
    ROOK_M = 11
    QUEEN = 12
    KING_R = 13
    KING_L = 14
    KING_D = 15
    KING_U = 16
    KING_R_M = 17
    KING_L_M = 18
    KING_D_M = 19
    KING_U_M = 20
)

var piece_moved_indexes = []int{
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
    ROOK_M,
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

var piece_move_functions = []func(Board, *Piece, *Point, *Array1000[FastMove]) {
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

func (p *Piece) value() int {
    return piece_values[p.index]
}

func (p *Piece) print() string {
    return piece_names[p.index]
}

func (p *Piece) movedIndex() int {
    return piece_moved_indexes[p.index]
}

func (p *Piece) isKing() bool {
    return p.index > QUEEN
}

func (p *Piece) isPawn() bool {
    return p.index < KNIGHT
}

func (p *Piece) moved() bool {
    if p.index > PAWN_U && p.index < KNIGHT {
        return true
    }
    if p.index == ROOK_M {
        return true
    }
    if p.index > KING_U {
        return true
    }
    return false
}

func (p *Piece) moves(b Board, fromLocation *Point, moves *Array1000[FastMove]) {
    piece_move_functions[p.index](b, p, fromLocation, moves)
}

func addDirection(
	b Board,
    fromPiece *Piece,
    fromLocation *Point,
	moves *Array1000[FastMove],
    direction *Point,
) {
    currentLocation := fromLocation
    currentPiece := fromPiece

    for {
        currentLocation = b.addIndex(currentLocation, direction)
        if currentLocation == nil {
            break
        }

        currentPiece = b.getPiece(currentLocation)
        if currentPiece == nil { // no piece
            addMoveSimple(b, fromPiece, fromLocation, currentPiece, currentLocation, nil, moves)
        } else if currentPiece.color != fromPiece.color { // enemy piece
            addMoveSimple(b, fromPiece, fromLocation, currentPiece, currentLocation, nil, moves)
            break
        } else { // ally piece
            addMoveAllyDefense(b, fromPiece, fromLocation, currentLocation, moves)
            break
        }
    }
}

func addSimple(
	b Board,
    fromPiece *Piece,
    fromLocation *Point,
	moves *Array1000[FastMove],
    direction *Point,
) {
    toLocation := b.addIndex(fromLocation, direction)
    if toLocation == nil {
        return
    }

	toPiece := b.getPiece(toLocation)
	if toPiece == nil { // no piece
        addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, nil, moves)
	} else if toPiece.color != fromPiece.color { // enemy piece
        addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, nil, moves)
	} else { // ally piece
        addMoveAllyDefense(b, fromPiece, fromLocation, toLocation, moves)
    }
}

var pawn_r_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_r_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_r_directions)
}

var pawn_l_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_l_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_l_directions)
}

var pawn_u_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_u_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_u_directions)
}

var pawn_d_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    pawnAddForward(b, fromPiece, fromLocation, moves, pawn_d_directions)
    pawnAddCaptures(b, fromPiece, fromLocation, moves, pawn_d_directions)
}

func pawnAddForward(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove], directions []*Point) {
    to1Location := b.addIndex(fromLocation, directions[0])
    if to1Location == nil { // location doesn't exist
        return
    }

    to2Location := b.addIndex(fromLocation, directions[1])
    to3Location := b.addIndex(fromLocation, directions[2])
    piece1 := b.getPiece(to1Location)
    piece2 := b.getPiece(to2Location)
    queen := b.getAllPiece(fromPiece.color, QUEEN)
    rook_m := b.getAllPiece(fromPiece.color, ROOK_M)
    bishop := b.getAllPiece(fromPiece.color, BISHOP)
    knight := b.getAllPiece(fromPiece.color, KNIGHT)

    if piece1 == nil { // no piece on location
        if to2Location == nil { // location doesn't exist
            addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, queen, moves)
            addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, rook_m, moves)
            addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, bishop, moves)
            addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, knight, moves)
            return
        } else {
            addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, nil, moves)
        }
    } else {
        return
    }

    if fromPiece.moved() {
        return
    }

    if piece2 == nil { // no piece on location
        if to3Location == nil { // location doesn't exist
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, to2Location, queen, to1Location, to2Location, moves)
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, to2Location, rook_m, to1Location, to2Location, moves)
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, to2Location, bishop, to1Location, to2Location, moves)
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, to2Location, knight, to1Location, to2Location, moves)
        } else {
            addMoveRevealEnPassant(b, fromPiece, fromLocation, piece2, to2Location, nil, to1Location, to2Location, moves)
        }
    }
}

func pawnAddCaptures(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove], directions []*Point) {
    to1Location := b.addIndex(fromLocation, directions[3])
    to2Location := b.addIndex(fromLocation, directions[4])
    to3Location := b.addIndex(to1Location, directions[0])
    to4Location := b.addIndex(to2Location, directions[0])
    piece1 := b.getPiece(to1Location)
    piece2 := b.getPiece(to2Location)
    queen := b.getAllPiece(fromPiece.color, QUEEN)
    rook_m := b.getAllPiece(fromPiece.color, ROOK_M)
    bishop := b.getAllPiece(fromPiece.color, BISHOP)
    knight := b.getAllPiece(fromPiece.color, KNIGHT)

    if to1Location != nil {
        if t1, t2, r1, r2 := b.possibleEnPassant(fromPiece.color, to1Location); t1 != nil { // if the square is an en passant target
            if to3Location == nil {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, to1Location, queen, t1, t2, r1, r2, moves)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, to1Location, rook_m, t1, t2, r1, r2, moves)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, to1Location, bishop, t1, t2, r1, r2, moves)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, to1Location, knight, t1, t2, r1, r2, moves)
            } else {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece1, to1Location, nil, t1, t2, r1, r2, moves)
            }
        } else if piece1 != nil && piece1.color != fromPiece.color { // if the square is occupied by an enemy piece
            if to3Location == nil {
                addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, queen, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, rook_m, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, bishop, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, knight, moves)
            } else {
                addMoveSimple(b, fromPiece, fromLocation, piece1, to1Location, nil, moves)
            }
        } else if piece1 != nil { // if the square is occupied by an ally piece
            addMoveAllyDefense(b, fromPiece, fromLocation, to1Location, moves)
        }
    }

    if to2Location != nil {
        if t1, t2, r1, r2 := b.possibleEnPassant(fromPiece.color, to2Location); t1 != nil { // if the square is an en passant target
            if to4Location == nil {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, to2Location, queen, t1, t2, r1, r2, moves)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, to2Location, rook_m, t1, t2, r1, r2, moves)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, to2Location, bishop, t1, t2, r1, r2, moves)
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, to2Location, knight, t1, t2, r1, r2, moves)
            } else {
                addMoveCaptureEnPassant(b, fromPiece, fromLocation, piece2, to2Location, nil, t1, t2, r1, r2, moves)
            }
        } else if piece2 != nil && piece2.color != fromPiece.color { // if the square is occupied by an enemy piece
            if to4Location == nil {
                addMoveSimple(b, fromPiece, fromLocation, piece2, to2Location, queen, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece2, to2Location, rook_m, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece2, to2Location, bishop, moves)
                addMoveSimple(b, fromPiece, fromLocation, piece2, to2Location, knight, moves)
            } else {
                addMoveSimple(b, fromPiece, fromLocation, piece2, to2Location, nil, moves)
            }
        } else if piece2 != nil { // if the square is occupied by an ally piece
            addMoveAllyDefense(b, fromPiece, fromLocation, to2Location, moves)
        }
    }
}

func knight_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
	for _, direction := range knight_directions {
		addSimple(b, fromPiece, fromLocation, moves, direction)
	}
}

func bishop_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    for _, direction := range bishop_directions {
        addDirection(b, fromPiece, fromLocation, moves, direction)
    }
}

func rook_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    for _, direction := range rook_directions {
        addDirection(b, fromPiece, fromLocation, moves, direction)
    }
}

func queen_moves(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    for _, direction := range queen_directions {
        addDirection(b, fromPiece, fromLocation, moves, direction)
    }
}

var king_lr_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    for _, direction := range queen_directions {
        addSimple(b, fromPiece, fromLocation, moves, direction)
    }

    if fromPiece.moved() {
        return
    }

    addCastle(b, fromPiece, fromLocation, moves, king_lr_directions[0], king_lr_directions[2], king_lr_directions[4])
    addCastle(b, fromPiece, fromLocation, moves, king_lr_directions[1], king_lr_directions[3], king_lr_directions[5])
}

var king_ud_moves = func(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove]) {
    for _, direction := range queen_directions {
        addSimple(b, fromPiece, fromLocation, moves, direction)
    }

    if fromPiece.moved() {
        return
    }

    addCastle(b, fromPiece, fromLocation, moves, king_ud_directions[0], king_ud_directions[2], king_ud_directions[4])
    addCastle(b, fromPiece, fromLocation, moves, king_ud_directions[1], king_ud_directions[3], king_ud_directions[5])
}

func addCastle(b Board, fromPiece *Piece, fromLocation *Point, moves *Array1000[FastMove], direction *Point, kingOffset *Point, rookOffset *Point) {
    // find rook for castle
    fromRookLocation := fromLocation
    var rook *Piece

    for {
        fromRookLocation = b.addIndex(fromRookLocation, direction)
        if fromRookLocation == nil { // exceeded board
            return
        }

        rook = b.getPiece(fromRookLocation)
        if rook == nil { // no piece at location
            continue
        }

        if rook.index == ROOK && rook.color == fromPiece.color { // found ally unmoved rook
            break
        }

        return
    }

    // find edge for castle
    edgeLocation := fromRookLocation
    currentLocation := edgeLocation

    for {
        currentLocation = b.addIndex(currentLocation, direction)
        if currentLocation == nil { // exceeded board
            break
        }

        edgeLocation = b.addIndex(edgeLocation, direction)
    }

    // everything else
    toLocation := b.addIndex(edgeLocation, kingOffset)
    toRookLocation := b.addIndex(edgeLocation, rookOffset)

    xCheckedMin := min(fromLocation.x, fromRookLocation.x)
    xCheckedMax := max(fromLocation.x, fromRookLocation.x)
    yCheckedMin := min(fromLocation.y, fromRookLocation.y)
    yCheckedMax := max(fromLocation.y, fromRookLocation.y)

    xToMin := min(toLocation.x, toRookLocation.x)
    xToMax := max(toLocation.x, toRookLocation.x)
    yToMin := min(toLocation.y, toRookLocation.y)
    yToMax := max(toLocation.y, toRookLocation.y)

    clr := true
    var index *Point 
    var piece *Piece

    for x := xCheckedMin - 1; x >= xToMin && clr; x-- {
        index = b.getIndex(x, fromLocation.y)
        if piece = b.getPiece(index); piece != nil {
            clr = false
            break
        }
    }
    for y := yCheckedMin - 1; y >= yToMin && clr; y-- {
        index = b.getIndex(fromLocation.x, y)
        if piece = b.getPiece(index); piece != nil {
            clr = false
            break
        }
    }
    for x := xCheckedMax + 1; x <= xToMax && clr; x++ {
        index = b.getIndex(x, fromLocation.y)
        if piece = b.getPiece(index); piece != nil {
            clr = false
            break
        }
    }
    for y := yCheckedMax + 1; y <= yToMax && clr; y++ {
        index = b.getIndex(fromLocation.x, y)
        if piece = b.getPiece(index); piece != nil {
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
    vulnerableStart := b.getIndex(minx, miny)
    vulnerableEnd := b.getIndex(maxx, maxy)

    addMoveCastle(b, fromPiece, fromLocation, toLocation, rook, fromRookLocation, toRookLocation, vulnerableStart, vulnerableEnd, moves)
}

