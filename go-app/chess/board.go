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

    getPieceLocations() []Array100[Point]
	getPiece(location *Point) *Piece
    getVulnerable(color int) *Vulnerable
	getEnPassant(color int) *EnPassant
    possibleEnPassant(color int, location *Point) (*EnPassant, *EnPassant)

    MovesOfColor(color int) *Array100[FastMove]
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

func newSimpleBoard(boardSize Point, numberOfPlayers int) (*SimpleBoard, error) {
    if boardSize.x <= 0 || boardSize.y <= 0 {
        return nil, fmt.Errorf("invalid board size")
    }

    if numberOfPlayers <= 0 {
        return nil, fmt.Errorf("invalid number of players")
    }

    playersDisabled := make([]bool, numberOfPlayers)
    for p := range playersDisabled {
        playersDisabled[p] = false
    }

    pieces := make([][]Piece, boardSize.y)
    for p := range pieces {
        pieces[p] = make([]Piece, boardSize.x)
    }

    disableds := make([][]bool, boardSize.y)
    for d := range disableds {
        disableds[d] = make([]bool, boardSize.x)
    }

    kingLocations := make([]Point, numberOfPlayers)
    for k := range kingLocations {
        kingLocations[k] = Point{-1, -1}
    }

    pieceLocations := make([]Array100[Point], numberOfPlayers)
    for p := range pieceLocations {
        pieceLocations[p] = Array100[Point]{}
    }

    enPassants := make([]EnPassant, numberOfPlayers)
    for e := range enPassants {
        enPassants[e] = EnPassant{Point{-1, -1}, Point{-1, -1}}
    }

    vulnerables := make([]Vulnerable, numberOfPlayers)
    for v := range vulnerables {
        vulnerables[v] = Vulnerable{Point{-1, -1}, Point{-1, -1}}
    }

    playerMoves := make([]Array100[FastMove], numberOfPlayers)
    for p := range playerMoves {
        playerMoves[p] = Array100[FastMove]{}
    }

    fromMoves := make([][]Array100[*FastMove], boardSize.y)
    for row := range fromMoves {
        fromMoves[row] = make([]Array100[*FastMove], boardSize.x)
        for col := range fromMoves[row] {
            fromMoves[row][col] = Array100[*FastMove]{}
        }
    }

    toMoves := make([][]Array100[*FastMove], boardSize.y)
    for row := range toMoves {
        toMoves[row] = make([]Array100[*FastMove], boardSize.x)
        for col := range toMoves[row] {
            toMoves[row][col] = Array100[*FastMove]{}
        }
    }

	return &SimpleBoard{
        size: boardSize,
        players: numberOfPlayers,

        playersDisabled: playersDisabled,

		pieces: pieces,
        disableds: disableds,

        kingLocations: kingLocations,
        pieceLocations: pieceLocations,

		enPassants: enPassants,
        vulnerables: vulnerables,

        playerMoves: playerMoves,
        fromMoves: fromMoves,
        toMoves: toMoves,

        test: false,
	}, nil
}

type SimpleBoard struct {
    size Point
    players int

    playersDisabled []bool

	pieces [][]Piece
    disableds [][]bool

    kingLocations []Point
    pieceLocations []Array100[Point]

	enPassants []EnPassant
    vulnerables []Vulnerable

    playerMoves []Array100[FastMove]
    fromMoves [][]Array100[*FastMove]
    toMoves [][]Array100[*FastMove]

    test bool
}

func (b *SimpleBoard) pointOutOfBounds(p *Point) bool {
    return p.y < 0 || p.y >= b.size.y || p.x < 0 || p.x >= b.size.x || b.disableds[p.y][p.x]
}

func (b *SimpleBoard) colorOutOfBounds(color int) bool {
    return color < 0 || color >= b.players
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

func (b *SimpleBoard) getPieceLocations() []Array100[Point] {
    return b.pieceLocations
}

func (b *SimpleBoard) getPiece(location *Point) *Piece {
    if b.pointOutOfBounds(location) {
        return nil
    }

    return &b.pieces[location.y][location.x]
}

func (b *SimpleBoard) getVulnerable(color int) *Vulnerable {
    if b.colorOutOfBounds(color) {
        return nil
    }

    return &b.vulnerables[color]
}

func (b *SimpleBoard) getEnPassant(color int) *EnPassant {
    if b.colorOutOfBounds(color) {
        return nil
    }

    return &b.enPassants[color]
}

func (b *SimpleBoard) possibleEnPassant(color int, target *Point) (*EnPassant, *EnPassant) {
    if b.colorOutOfBounds(color) {
        return nil, nil
    }

    var en1 *EnPassant
    var en2 *EnPassant

	for c, e := range b.enPassants {
        if c == color {
            continue
        }

        if !e.target.equals(target) {
            continue
        }

        if en1 == nil {
            en1 = &e
        } else if en2 == nil {
            en2 = &e
        } else {
            panic("too many en passants")
        }
	}

    return en1, en2
}

func (b *SimpleBoard) MovesOfColor(color int) *Array100[FastMove] {
    if b.colorOutOfBounds(color) {
        return nil
    }

    return &b.playerMoves[color]
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

        err := move.execute()
        if err != nil {
            return nil, err
        }

        b.CalculateMoves()

        if !b.Check(color) {
            legalMoves = append(legalMoves, move)
        }

        err = move.undo()
        if err != nil {
            return nil, err
        }
    }

    b.CalculateMoves()

    return legalMoves, nil
}

func (b *SimpleBoard) LegalMovesOfLocation(fromLocation *Point) ([]FastMove, error) {
    piece := b.getPiece(fromLocation)
    if piece == nil || !piece.valid() {
        return nil, fmt.Errorf("invalid piece")
    }

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

        err := move.execute()
        if err != nil {
            return nil, err
        }

        b.CalculateMoves()

        if !b.Check(piece.color) {
            legalMoves = append(legalMoves, move)
        }

        err = move.undo()
        if err != nil {
            return nil, err
        }
    }

    b.CalculateMoves()

    return legalMoves, nil
}

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
        b.playerMoves[i].clear()
        b.pieceLocations[i].clear()
    }

    for y := 0; y < b.size.y; y++ {
        for x := 0; x < b.size.x; x++ {
            b.toMoves[y][x].clear()
            b.fromMoves[y][x].clear()
        }
    }

    for y := 0; y < b.size.y; y++ {
        for x := 0; x < b.size.x; x++ {
            if b.disableds[y][x] {
                continue
            }

            piece := b.pieces[y][x]
            if !piece.valid() {
                continue
            }

            pieceLocations := b.pieceLocations[piece.color]
            pieceLocation := pieceLocations.get()
            pieceLocation.x = x
            pieceLocation.y = y
            pieceLocations.next()
            if piece.index > 13 {
                b.kingLocations[piece.color] = Point{x, y}
            }

            piece.moves(b, &Point{x, y}, &b.playerMoves[piece.color])
        }
    }

    for i := 0; i < b.players; i++ {
        moves := b.playerMoves[i]

        for j := 0; j < moves.count; j++ {
            move := moves.array[j]

            fromMove := b.fromMoves[move.fromLocation.y][move.fromLocation.x].get()
            *fromMove = &move
            b.fromMoves[move.fromLocation.y][move.fromLocation.x].next()

            toMove := b.toMoves[move.toLocation.y][move.toLocation.x].get()
            *toMove = &move
            b.toMoves[move.toLocation.y][move.toLocation.x].next()
        }
    }
}

func (b *SimpleBoard) Check(color int) bool {
    if b.colorOutOfBounds(color) {
        return false
    }

    if b.pointOutOfBounds(&b.kingLocations[color]) {
        return false
    }

    kingLocation := b.kingLocations[color]
    for i := 0; i < b.toMoves[kingLocation.y][kingLocation.x].count; i++ {
        move := b.toMoves[kingLocation.y][kingLocation.x].array[i]

        if move.color != color {
            return true
        }
    }

    if b.pointOutOfBounds(&b.vulnerables[color].start) {
        return false
    }

    if b.pointOutOfBounds(&b.vulnerables[color].end) {
        return false
    }

    vulnerable := b.vulnerables[color]
    for x := vulnerable.start.x; x <= vulnerable.end.x; x++ {
        for y := vulnerable.start.y; y <= vulnerable.end.y; y++ {
            for i := 0; i < b.toMoves[y][x].count; i++ {
                move := b.toMoves[y][x].array[i]

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
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*b.size.x-1)))
		for x := range row {
			builder.WriteString(fmt.Sprintf("|%s%2dx ", strings.Repeat(" ", cellWidth-4), x))
		}
		builder.WriteString("|\n")
		for x, piece := range row {
            if b.pointOutOfBounds(&Point{x, y}) {
                builder.WriteString(fmt.Sprintf("|%s", strings.Repeat("X", cellWidth)))
            } else if !piece.valid() {
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
	builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", (cellWidth+1)*b.size.x-1)))

	return builder.String()
}

func (b *SimpleBoard) State() *BoardData {
    pieces := []*PieceData{}
    for y, row := range b.pieces {
        for x, piece := range row {
            if piece.valid() {
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
        XSize: b.size.x,
        YSize: b.size.y,
        Disabled: disabled,
        Pieces: pieces,
    }
}

func (b *SimpleBoard) Copy() (Board, error) {
    simpleBoard, err := newSimpleBoard(b.size, b.players)
    if err != nil {
        return nil, err
    }

    for color, disabled := range b.playersDisabled {
        simpleBoard.playersDisabled[color] = disabled
    }

    for y, row := range b.pieces {
        for x := range row {
            simpleBoard.pieces[y][x] = b.pieces[y][x]
        }
    }

    for y, row := range b.disableds {
        for x := range row {
            simpleBoard.disableds[y][x] = b.disableds[y][x]
        }
    }

    for color, enPassant := range b.enPassants {
        simpleBoard.enPassants[color] = enPassant
    }

    for color, vulnerable := range b.vulnerables {
        simpleBoard.vulnerables[color] = vulnerable
    }

    return simpleBoard, nil
}

func (b *SimpleBoard) UniqueString() string {
    builder := strings.Builder{}

    counter := 0
    for y, row := range b.pieces {
        for x := range row {
            piece := b.pieces[y][x]

            if !piece.valid() {
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

