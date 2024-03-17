package chess

import (
	"fmt"
	"strings"
)

/*
Responsible for:
- keeping track of the pieces on the board
- keeping track of availalbe moves for all pieces
*/
type Board interface {
    disablePieces(color int, disable bool)
    disableLocation(location *Point)

    getIndex(x int, y int) *Point
    addIndex(index1 *Point, index2 *Point) *Point

    getAllPiece(color int, index int) *Piece
    getPieceLocations() []Array100[*Point]
	getPiece(location *Point) *Piece
    movePiece(piece *Piece) *Piece
    setPiece(location *Point, piece *Piece)

    getVulnerable(color int) (*Point, *Point)
    setVulnerable(color int, start *Point, end *Point)

	getEnPassant(color int) (*Point, *Point)
	setEnPassant(color int, target *Point, risk *Point)
    possibleEnPassant(color int, location *Point) (*Point, *Point, *Point, *Point)

    MovesOfColor(color int) *Array1000[FastMove]
    MovesOfLocation(fromLocation *Point) *Array100[FastMove]
    LegalMovesOfColor(color int) ([]FastMove, error)
    LegalMovesOfLocation(fromLocation *Point) ([]FastMove, error)

    CalculateMoves()
    Check(color int) bool
    CheckmateAndStalemate(color int) (bool, bool, error)

	Print() string
    State() *BoardData
    Copy() (Board, error) 
    UniqueString() string
}

func newSimpleBoard(x int, y int, players int) (*SimpleBoard, error) {
    if x <= 0 || y <= 0 {
        return nil, fmt.Errorf("invalid board size")
    }

    if players <= 0 {
        return nil, fmt.Errorf("invalid number of players")
    }

    test := false

    playersDisabled := make([]bool, players)
    enPassantTargets := make([]*Point, players)
    enPassantRisks := make([]*Point, players)
    vulnerableStarts := make([]*Point, players)
    vulnerableEnds := make([]*Point, players)
    kingLocations := make([]*Point, players)
    pieceLocations := make([]Array100[*Point], players)
    moves := make([]Array1000[FastMove], players)
    allPieces := make([][]Piece, players)
    for i := 0; i < players; i++ {
        playersDisabled[i] = false
        enPassantTargets[i] = nil // TODO we don't account for this being nil
        enPassantRisks[i] = nil
        vulnerableStarts[i] = nil
        vulnerableEnds[i] = nil
        kingLocations[i] = nil
        pieceLocations[i] = Array100[*Point]{}
        moves[i] = Array1000[FastMove]{}
        allPieces[i] = []Piece{
            {i, PAWN_R},
            {i, PAWN_L},
            {i, PAWN_D},
            {i, PAWN_U},
            {i, PAWN_R_M},
            {i, PAWN_L_M},
            {i, PAWN_D_M},
            {i, PAWN_U_M},
            {i, KNIGHT},
            {i, BISHOP},
            {i, ROOK},
            {i, ROOK_M},
            {i, QUEEN},
            {i, KING_R},
            {i, KING_L},
            {i, KING_D},
            {i, KING_U},
            {i, KING_R_M},
            {i, KING_L_M},
            {i, KING_D_M},
            {i, KING_U_M},
        }
    }

    disableds := make([][]bool, y)
    indexes := make([][]Point, y)
    pieces := make([][]*Piece, y)
    toMoves := make([][]Array100[*FastMove], y)
    fromMoves := make([][]Array100[*FastMove], y)
    for i := 0; i < y; i++ {
        disableds[i] = make([]bool, x)
        indexes[i] = make([]Point, x)
        pieces[i] = make([]*Piece, x)
        toMoves[i] = make([]Array100[*FastMove], x)
        fromMoves[i] = make([]Array100[*FastMove], x)

        for j := 0; j < x; j++ {
            disableds[i][j] = false
            indexes[i][j] = Point{j, i}
            pieces[i][j] = nil
            toMoves[i][j] = Array100[*FastMove]{}
            fromMoves[i][j] = Array100[*FastMove]{}
        }
    }

	return &SimpleBoard{
        x: x,
        y: y,
        players: players,
        test: test,

        playersDisabled: playersDisabled,
        enPassantTargets: enPassantTargets,
        enPassantRisks: enPassantRisks,
        vulnerableStarts: vulnerableStarts,
        vulnerableEnds: vulnerableEnds,
        kingLocations: kingLocations,
        pieceLocations: pieceLocations,
        moves: moves,
        allPieces: allPieces,

        disableds: disableds,
        indexes: indexes,
        pieces: pieces,
        toMoves: toMoves,
        fromMoves: fromMoves,
	}, nil
}

type SimpleBoard struct {
    x int
    y int
    players int
    test bool

    // arrays of size PLAYERS
    playersDisabled []bool
	enPassantTargets []*Point
    enPassantRisks []*Point
    vulnerableStarts []*Point
    vulnerableEnds []*Point
    kingLocations []*Point
    pieceLocations []Array100[*Point]
    moves []Array1000[FastMove]
    allPieces [][]Piece

    // arrays of size X * Y
    disableds [][]bool
    indexes [][]Point
    pieces [][]*Piece
    toMoves [][]Array100[*FastMove]
    fromMoves [][]Array100[*FastMove]
}

func (b *SimpleBoard) pointOutOfBounds(p *Point) bool {
    if p == nil {
        return true
    }

    return p.y < 0 || p.y >= b.y || p.x < 0 || p.x >= b.x || b.disableds[p.y][p.x]
}

func (b *SimpleBoard) colorOutOfBounds(c int) bool {
    return c < 0 || c >= b.players
}

func (b *SimpleBoard) disablePieces(color int, disable bool) {
    if b.colorOutOfBounds(color) {
        return
    }

    b.playersDisabled[color] = disable
}

func (b *SimpleBoard) disableLocation(location *Point) {
    if b.pointOutOfBounds(location) {
        return
    }

    b.disableds[location.y][location.x] = true
}

func (b *SimpleBoard) getIndex(x int, y int) *Point {
    if x < 0 || x >= b.x || y < 0 || y >= b.y {
        return nil
    }

    if b.disableds[y][x] {
        return nil
    }

    return &b.indexes[y][x]
}

func (b *SimpleBoard) addIndex(index1 *Point, index2 *Point) *Point {
    if index1 == nil || index2 == nil {
        return nil
    }

    return b.getIndex(index1.x + index2.x, index1.y + index2.y)
}

func (b *SimpleBoard) getAllPiece(color int, index int) *Piece {
    if b.colorOutOfBounds(color) {
        return nil
    }

    if index < 0 {
        return nil
    }

    return &b.allPieces[color][index]
}

func (b *SimpleBoard) getPieceLocations() []Array100[*Point] {
    return b.pieceLocations
}

func (b *SimpleBoard) getPiece(location *Point) *Piece {
    if b.pointOutOfBounds(location) {
        return nil
    }

    return b.pieces[location.y][location.x]
}

func (b *SimpleBoard) movePiece(piece *Piece) *Piece {
    return b.getAllPiece(piece.color, piece.movedIndex())
}

func (b *SimpleBoard) setPiece(location *Point, piece *Piece) {
    if b.pointOutOfBounds(location) {
        return
    }

    b.pieces[location.y][location.x] = piece
}

func (b *SimpleBoard) getVulnerable(color int) (*Point, *Point) {
    if b.colorOutOfBounds(color) {
        return nil, nil
    }

    return b.vulnerableStarts[color], b.vulnerableEnds[color]
}

func (b *SimpleBoard) getEnPassant(color int) (*Point, *Point) {
    if b.colorOutOfBounds(color) {
        return nil, nil
    }

    return b.enPassantTargets[color], b.enPassantRisks[color]
}

func (b *SimpleBoard) setEnPassant(color int, target *Point, risk *Point) {
    if b.colorOutOfBounds(color) {
        return
    }

    b.enPassantTargets[color] = target
    b.enPassantRisks[color] = risk
}

func (b *SimpleBoard) setVulnerable(color int, start *Point, end *Point) {
    if b.colorOutOfBounds(color) {
        return
    }

    b.vulnerableStarts[color] = start
    b.vulnerableEnds[color] = end
}

func (b *SimpleBoard) possibleEnPassant(color int, target *Point) (*Point, *Point, *Point, *Point) {
    if target == nil {
        return nil, nil, nil, nil
    }

    if b.colorOutOfBounds(color) {
        return nil, nil, nil, nil
    }

    var target1 *Point
    var target2 *Point
    var risk1 *Point
    var risk2 *Point

    for i := 0; i < b.players; i++ {
        if i == color {
            continue
        }

        t := b.enPassantTargets[i]
        r := b.enPassantRisks[i]
        if t == nil || r == nil {
            continue
        }

        if !target.equals(t) {
            continue
        }

        if target1 == nil {
            target1 = t
            risk1 = r
        } else if target2 == nil {
            target2 = t
            risk2 = r
        } else {
            panic("too many en passants")
        }
	}

    return target1, target2, risk1, risk2
}

func (b *SimpleBoard) MovesOfColor(color int) *Array1000[FastMove] {
    if b.colorOutOfBounds(color) {
        return nil
    }

    return &b.moves[color]
}

func (b *SimpleBoard) MovesOfLocation(fromLocation *Point) *Array100[FastMove] {
    if b.pointOutOfBounds(fromLocation) {
        return nil
    }

    moves := Array100[FastMove]{}
    for i := 0; i < b.fromMoves[fromLocation.y][fromLocation.x].count; i++ {
        move := b.fromMoves[fromLocation.y][fromLocation.x].array[i]
        currentMove := moves.get()
        *currentMove = *move
        moves.next()
    }
    return &moves
}

func (b *SimpleBoard) LegalMovesOfColor(color int) ([]FastMove, error) {
    movesPointer := b.MovesOfColor(color)
    if movesPointer == nil {
        return nil, fmt.Errorf("invalid color")
    }
    moves := *movesPointer

    legalMoves := []FastMove{}

    for i := 0; i < moves.count; i++ {
        move := moves.array[i]
        if move.allyDefense {
            continue
        }

        move.execute()

        b.CalculateMoves()
        if !b.Check(color) {
            legalMoves = append(legalMoves, move)
        }

        move.undo()
    }

    b.CalculateMoves()

    return legalMoves, nil
}

func (b *SimpleBoard) LegalMovesOfLocation(fromLocation *Point) ([]FastMove, error) {
    piecePointer := b.getPiece(fromLocation)
    if piecePointer == nil {
        return nil, fmt.Errorf("invalid piece")
    }
    color := piecePointer.color

    movesPointer := b.MovesOfLocation(fromLocation)
    if movesPointer == nil {
        return nil, fmt.Errorf("invalid location")
    }
    moves := *movesPointer

    legalMoves := []FastMove{}

    for i := 0; i < moves.count; i++ {
        move := moves.array[i]
        if move.allyDefense {
            continue
        }

        move.execute()

        b.CalculateMoves()
        if !b.Check(color) {
            legalMoves = append(legalMoves, move)
        }

        move.undo()
    }

    b.CalculateMoves()

    return legalMoves, nil
}

// don't use the make function
// remove the remaining interfaces
// get indexes from board to avoid creating a lot of points
// Returning items from functions results in heap allocation
// TODO remove set piece and just return pointers from this struct
// TODO edit existing structs when doing moves
// TODO implement dynamic move calculations based on previous move
// TODO how about we don't create massive move objects with pieces and stuff
// stop excessive use of maps
// stop excessive use of pointers
// don't store the moves in the board maybe
// reduce calls to append maybe
// TODO add 3 move repetition and 50 move rule
// TODO add rule to allow checks and only lose on king capture
// TODO add rule to check for checkmate and stalemate on all players after every move
func (b *SimpleBoard) CalculateMoves() {
    for i := 0; i < b.players; i++ {
        b.moves[i].clear()
        b.pieceLocations[i].clear()
    }

    for y := 0; y < b.y; y++ {
        for x := 0; x < b.x; x++ {
            b.toMoves[y][x].clear()
            b.fromMoves[y][x].clear()
        }
    }

    for y := 0; y < b.y; y++ {
        for x := 0; x < b.x; x++ {
            if b.disableds[y][x] {
                continue
            }

            piece := b.pieces[y][x]
            if piece == nil {
                continue
            }

            if b.playersDisabled[piece.color] {
                continue
            }

            index := b.getIndex(x, y)

            p := b.pieceLocations[piece.color].get()
            *p = index
            b.pieceLocations[piece.color].next()

            if piece.isKing() {
                b.kingLocations[piece.color] = index
            }

            piece.moves(b, index, &b.moves[piece.color])
        }
    }

    for i := 0; i < b.players; i++ {
        moves := b.moves[i]

        for j := 0; j < moves.count; j++ {
            move := &moves.array[j]

            fromMove := b.fromMoves[move.fromLocation.y][move.fromLocation.x].get()
            *fromMove = move
            b.fromMoves[move.fromLocation.y][move.fromLocation.x].next()

            toMove := b.toMoves[move.toLocation.y][move.toLocation.x].get()
            *toMove = move
            b.toMoves[move.toLocation.y][move.toLocation.x].next()
        }
    }
}

func (b *SimpleBoard) Check(color int) bool {
    if b.colorOutOfBounds(color) {
        return false
    }

    kingLocation := b.kingLocations[color]
    if b.pointOutOfBounds(kingLocation) {
        return false
    }

    for i := 0; i < b.toMoves[kingLocation.y][kingLocation.x].count; i++ {
        move := b.toMoves[kingLocation.y][kingLocation.x].array[i]

        if move.color != color {
            return true
        }
    }

    vulnerableStart := b.vulnerableStarts[color]
    if b.pointOutOfBounds(b.vulnerableStarts[color]) {
        return false
    }

    vulnerableEnd := b.vulnerableEnds[color]
    if b.pointOutOfBounds(b.vulnerableEnds[color]) {
        return false
    }

    for y := vulnerableStart.y; y <= vulnerableEnd.y; y++ {
        for x := vulnerableStart.x; x <= vulnerableEnd.x; x++ {
            moves := b.toMoves[y][x]

            for i := 0; i < moves.count; i++ {
                move := moves.array[i]

                if move.color != color {
                    return true
                }
            }
        }
    }

    return false
}

func (b *SimpleBoard) CheckmateAndStalemate(color int) (bool, bool, error) {
    legalMoves, err := b.LegalMovesOfColor(color)
    if err != nil {
        return false, false, err
    }

    if len(legalMoves) > 0 {
        return false, false, nil
    }

    if b.Check(color) {
        return true, false, nil
    }

    return false, true, nil
}

func (b *SimpleBoard) Print() string {
	var builder strings.Builder
	var cellWidth int = 12

	for y, row := range b.pieces {
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*b.x-1)))
		for x := range row {
			builder.WriteString(fmt.Sprintf("|%s%2dx ", strings.Repeat(" ", cellWidth-4), x))
		}
		builder.WriteString("|\n")
		for x, piece := range row {
            if b.pointOutOfBounds(b.getIndex(x, y)) {
                builder.WriteString(fmt.Sprintf("|%s", strings.Repeat("X", cellWidth)))
            } else if piece == nil {
				builder.WriteString(fmt.Sprintf("|%s", strings.Repeat(" ", cellWidth)))
			} else {
				p := piece.print()
				if len(p) > 1 {
					p = p[:1]
				}

                pColor := fmt.Sprintf("%d", piece.color)
				if len(pColor) > 8 {
					pColor = pColor[:8]
				}

				builder.WriteString(fmt.Sprintf("| %-1s %-8s ", p, pColor))
			}
		}
		builder.WriteString("|\n")
		for range row {
			builder.WriteString(fmt.Sprintf("|%s%2dy ", strings.Repeat(" ", cellWidth-4), y))
		}
		builder.WriteString("|\n")
	}
	builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*b.x-1)))

	return builder.String()
}

func (b *SimpleBoard) State() *BoardData {
    pieces := []*PieceData{}
    for y, row := range b.pieces {
        for x, piece := range row {
            if piece != nil {
                disabled := b.playersDisabled[piece.color]
                pieces = append(pieces, &PieceData{
                    T: piece.print(),
                    C: piece.color,
                    X: x,
                    Y: y,
                    D: disabled,
                })
            }
        }
    }

    disabled := []*DisabledData{}
    for y, row := range b.disableds {
        for x, d := range row {
            if d {
                disabled = append(disabled, &DisabledData{
                    X: x,
                    Y: y,
                })
            }
        }
    }

    return &BoardData{
        XSize: b.x,
        YSize: b.y,
        Disabled: disabled,
        Pieces: pieces,
    }
}

func (b *SimpleBoard) Copy() (Board, error) {
    simpleBoard, err := newSimpleBoard(b.x, b.y, b.players)
    if err != nil {
        return nil, err
    }

    for i := 0; i < b.players; i++ {
        simpleBoard.playersDisabled[i] = b.playersDisabled[i]
        simpleBoard.enPassantTargets[i] = simpleBoard.getIndex(b.enPassantTargets[i].x, b.enPassantTargets[i].y)
        simpleBoard.enPassantRisks[i] = simpleBoard.getIndex(b.enPassantRisks[i].x, b.enPassantRisks[i].y)
        simpleBoard.vulnerableStarts[i] = simpleBoard.getIndex(b.vulnerableStarts[i].x, b.vulnerableStarts[i].y)
        simpleBoard.vulnerableEnds[i] = simpleBoard.getIndex(b.vulnerableEnds[i].x, b.vulnerableEnds[i].y)
    }

    for i := 0; i < b.y; i++ {
        for j := 0; j < b.x; j++ {
            simpleBoard.disableds[i][j] = b.disableds[i][j]
            simpleBoard.pieces[i][j] = b.pieces[i][j]
        }
    }

    return simpleBoard, nil
}

func (b *SimpleBoard) UniqueString() string {
    builder := strings.Builder{}

    counter := 0
    for y, row := range b.pieces {
        for x := range row {
            piece := b.pieces[y][x]

            if piece == nil {
                counter += 1
                continue
            }

            if counter > 0 {
                builder.WriteString(fmt.Sprintf("%d", counter))
                counter = 0
            }

            if b.playersDisabled[piece.color] {
                builder.WriteString("d")
                continue
            }

            builder.WriteString(piece.print())
            builder.WriteString(fmt.Sprintf("%d", piece.color))
            if piece.moved() {
                builder.WriteString("m")
            }
        }
    }

    return builder.String()
}

