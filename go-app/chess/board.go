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
    // these are for the move
	getPiece(location Point) (Piece, bool)
	setPiece(location Point, piece Piece) bool
    disableLocation(location Point) error
    getVulnerable(color int) (Vulnerable, error) // if these locations are attacked, the player is in check
    setVulnerable(color int, vulnerable Vulnerable) error
	getEnPassant(color int) (EnPassant, error) // if these locations are attacked, a piece is captured en passant
	setEnPassant(color int, enPassant EnPassant) error
    possibleEnPassant(color int, location Point) ([]EnPassant, error)

    // these are for the playerTransition
    disablePieces(color int, disable bool) error

    // these are for the bot
    getPieceLocations() [][]Point

    // these are for the game
    CalculateMoves() error // calcutes moves for every color
    MovesOfColor(color int) (*Array100[FastMove], error) // returns moves
    MovesOfLocation(fromLocation Point) (*Array100[FastMove], error) // returns moves
    LegalMovesOfColor(color int) ([]FastMove, error) // calculates and returns moves that do not result in check
    LegalMovesOfLocation(fromLocation Point) ([]FastMove, error) // calculates and returns moves that do not result in check
    CheckmateAndStalemate(color int) (bool, bool, error) // calculates legal moves to return checkmate and stalemate
    Check(color int) bool // returns whether the player is in check

    State() *BoardData
	Print() string
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

    disabledLocations := make([][]bool, boardSize.y)
    for d := range disabledLocations {
        disabledLocations[d] = make([]bool, boardSize.x)
    }

    kingLocations := make([]Point, numberOfPlayers)
    for k := range kingLocations {
        kingLocations[k] = Point{-1, -1}
    }

    pieceLocations := make([][]Point, numberOfPlayers)
    for p := range pieceLocations {
        pieceLocations[p] = []Point{}
    }

    enPassants := make([]EnPassant, numberOfPlayers)
    for e := range enPassants {
        enPassants[e] = EnPassant{Point{-1, -1}, Point{-1, -1}}
    }

    vulnerables := make([]Vulnerable, numberOfPlayers)
    for v := range vulnerables {
        vulnerables[v] = Vulnerable{Point{-1, -1}, Point{-1, -1}}
    }

    playerMoves := make([]*Array100[FastMove], numberOfPlayers)
    for p := range playerMoves {
        playerMoves[p] = &Array100[FastMove]{}
    }

    fromMoves := make([][]*Array100[*FastMove], boardSize.y)
    for row := range fromMoves {
        fromMoves[row] = make([]*Array100[*FastMove], boardSize.x)
        for col := range fromMoves[row] {
            fromMoves[row][col] = &Array100[*FastMove]{}
        }
    }

    toMoves := make([][]*Array100[*FastMove], boardSize.y)
    for row := range toMoves {
        toMoves[row] = make([]*Array100[*FastMove], boardSize.x)
        for col := range toMoves[row] {
            toMoves[row][col] = &Array100[*FastMove]{}
        }
    }

	return &SimpleBoard{
        size: boardSize,
        players: numberOfPlayers,

        playersDisabled: playersDisabled,
		pieces: pieces,
        disabledLocations: disabledLocations,
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
    disabledLocations [][]bool
    kingLocations []Point
    pieceLocations [][]Point

	enPassants []EnPassant
    vulnerables []Vulnerable

    playerMoves []*Array100[FastMove]
    fromMoves [][]*Array100[*FastMove]
    toMoves [][]*Array100[*FastMove]

    test bool
}

func (b *SimpleBoard) disablePieces(color int, disable bool) error {
    if b.colorOutOfBounds(color) {
        return fmt.Errorf("invalid color")
    }

    b.playersDisabled[color] = disable
    return nil
}

func (b *SimpleBoard) getPiece(location Point) (Piece, bool) {
    if b.pointOutOfBounds(location) {
        return Piece{0, 0}, false
    }

	return b.pieces[location.y][location.x], true
}

func (b *SimpleBoard) setPiece(location Point, p Piece) bool {
    if b.pointOutOfBounds(location) {
        return false
    }

    oldPiece, ok := b.getPiece(location)
    if !ok {
        return false
    }

    if oldPiece.valid() {
        if b.colorOutOfBounds(oldPiece.color) {
            return false
        }

        pieceLocations := b.pieceLocations[oldPiece.color]

        removeIndex := -1
        for i, pieceLocation := range pieceLocations {
            if pieceLocation.equals(location) {
                removeIndex = i
                break
            }
        }
        if removeIndex != -1 {
            pieceLocations[removeIndex] = pieceLocations[len(pieceLocations)-1]
            pieceLocations[len(pieceLocations)-1] = Point{}
            pieceLocations = pieceLocations[:len(pieceLocations)-1]
            b.pieceLocations[oldPiece.color] = pieceLocations
        }
    }

    if p.valid() {
        if b.colorOutOfBounds(p.color) {
            return false
        }

        pieceLocations := b.pieceLocations[p.color]

        pieceLocations = append(pieceLocations, location)
        b.pieceLocations[p.color] = pieceLocations

        if p.index > 13 {
            b.kingLocations[p.color] = location
        }
    }

	b.pieces[location.y][location.x] = p
    return true
}

func (b *SimpleBoard) disableLocation(location Point) error {
    if b.pointOutOfBounds(location) {
        return fmt.Errorf("invalid location")
    }

    b.disabledLocations[location.y][location.x] = true
    return nil
}

func (b *SimpleBoard) getVulnerable(color int) (Vulnerable, error) {
    if b.colorOutOfBounds(color) {
        return Vulnerable{}, fmt.Errorf("invalid color")
    }

    return b.vulnerables[color], nil
}

func (b *SimpleBoard) setVulnerable(color int, vulnerable Vulnerable) error {
    if b.colorOutOfBounds(color) {
        return fmt.Errorf("invalid color")
    }

    b.vulnerables[color] = vulnerable
    return nil
}

func (b *SimpleBoard) getEnPassant(color int) (EnPassant, error) {
    if b.colorOutOfBounds(color) {
        return EnPassant{}, fmt.Errorf("invalid color")
    }

	return b.enPassants[color], nil
}

func (b *SimpleBoard) setEnPassant(color int, enPassant EnPassant) error {
    if b.colorOutOfBounds(color) {
        return fmt.Errorf("invalid color")
    }

	b.enPassants[color] = enPassant
    return nil
}

func (b *SimpleBoard) possibleEnPassant(color int, target Point) ([]EnPassant, error) {
    if b.colorOutOfBounds(color) {
        return []EnPassant{}, fmt.Errorf("invalid color")
    }

    enPassants := []EnPassant{}
	for c, e := range b.enPassants {
        if c == color {
            continue
        }

        if !e.target.equals(target) {
            continue
        }

        enPassants = append(enPassants, e)
	}
    return enPassants, nil
}

func (b *SimpleBoard) MovesOfColor(color int) (*Array100[FastMove], error) {
    if b.colorOutOfBounds(color) {
        return nil, fmt.Errorf("invalid color")
    }

    return b.playerMoves[color], nil
}

func (b *SimpleBoard) MovesOfLocation(fromLocation Point) (*Array100[FastMove], error) {
    if b.pointOutOfBounds(fromLocation) {
        return nil, fmt.Errorf("invalid location")
    }

    moves := Array100[FastMove]{}
    for i := 0; i < b.fromMoves[fromLocation.y][fromLocation.x].count; i++ {
        move := b.fromMoves[fromLocation.y][fromLocation.x].array[i]
        moves.append(*move)
    }
    return &moves, nil
}

func (b *SimpleBoard) LegalMovesOfColor(color int) ([]FastMove, error) {
    movesPointer, err := b.MovesOfColor(color)
    moves := *movesPointer
    if err != nil {
        return []FastMove{}, err
    }

    legalMoves := []FastMove{}
    for i := 0; i < moves.count; i++ {
        move := moves.array[i]
        if move.allyDefense {
            continue
        }

        err = move.execute()
        if err != nil {
            return []FastMove{}, err
        }

        err = b.CalculateMoves()
        if err != nil {
            return []FastMove{}, err
        }

        if !b.Check(color) {
            legalMoves = append(legalMoves, move)
        }

        err = move.undo()
        if err != nil {
            return []FastMove{}, err
        }
    }

    err = b.CalculateMoves()
    if err != nil {
        return []FastMove{}, err
    }

    return legalMoves, nil
}

func (b *SimpleBoard) LegalMovesOfLocation(fromLocation Point) ([]FastMove, error) {
    piece, ok := b.getPiece(fromLocation)
    if !ok || !piece.valid() {
        return []FastMove{}, fmt.Errorf("piece not found")
    }

    movesPointer, err := b.MovesOfLocation(fromLocation)
    moves := *movesPointer
    if err != nil {
        return []FastMove{}, err
    }

    legalMoves := []FastMove{}
    for i := 0; i < moves.count; i++ {
        move := moves.array[i]
        if move.allyDefense {
            continue
        }

        err = move.execute()
        if err != nil {
            return []FastMove{}, err
        }

        err = b.CalculateMoves()
        if err != nil {
            return []FastMove{}, err
        }

        if !b.Check(piece.color) {
            legalMoves = append(legalMoves, move)
        }

        err = move.undo()
        if err != nil {
            return []FastMove{}, err
        }
    }

    err = b.CalculateMoves()
    if err != nil {
        return []FastMove{}, err
    }

    return legalMoves, nil
}

// TODO add 3 move repetition and 50 move rule
// TODO add rule to allow checks and only lose on king capture
// TODO add rule to check for checkmate and stalemate on all players after every move
func (b *SimpleBoard) CalculateMoves() error {
    if b.test {
        return nil
    }

    for i := 0; i < b.players; i++ {
        b.playerMoves[i].clear()
    }

    for y := 0; y < b.size.y; y++ {
        for x := 0; x < b.size.x; x++ {
            b.toMoves[y][x].clear()
            b.fromMoves[y][x].clear()
        }
    }

    for color, pieceLocations := range b.pieceLocations {
        for _, fromLocation := range pieceLocations {
            piece, _ := b.getPiece(fromLocation)
            if !piece.valid() {
                continue
            }

            if b.playersDisabled[color] {
                continue
            }

            moves := b.playerMoves[color]
            piece.moves(b, fromLocation, moves)
        }
    }

    for i := 0; i < b.players; i++ {
        moves := b.playerMoves[i]
        for j := 0; j < moves.count; j++ {
            move := moves.array[j]

            b.fromMoves[move.fromLocation.y][move.fromLocation.x].append(&move)
            b.toMoves[move.toLocation.y][move.toLocation.x].append(&move)
        }
    }

    return nil
}

// TODO edit existing structs when doing moves
// TODO implement dynamic move calculations based on previous move
// TODO how about we don't create massive move objects with pieces and stuff
// stop excessive use of maps
// stop excessive use of pointers
// don't store the moves in the board maybe
// reduce calls to append maybe
func (b *SimpleBoard) CalculateMovesPartial(move FastMove) error {
    return fmt.Errorf("not implemented")
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
            if b.pointOutOfBounds(Point{x, y}) {
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
    for y, row := range b.disabledLocations {
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

func (b *SimpleBoard) Check(color int) bool {
    if b.colorOutOfBounds(color) {
        return false
    }

    kingLocation := b.kingLocations[color]
    for i := 0; i < b.toMoves[kingLocation.y][kingLocation.x].count; i++ {
        move := b.toMoves[kingLocation.y][kingLocation.x].array[i]
        piece, ok := b.getPiece(move.fromLocation)
        if !ok || !piece.valid() {
            continue
        }

        if piece.color != color {
            return true
        }
    }

    vulnerable := b.vulnerables[color]
    for x := max(0, vulnerable.start.x); x <= vulnerable.end.x; x++ {
        for y := max(0, vulnerable.start.y); y <= vulnerable.end.y; y++ {
            for i := 0; i < b.toMoves[y][x].count; i++ {
                move := b.toMoves[y][x].array[i]
                piece, ok := b.getPiece(move.fromLocation)
                if !ok || !piece.valid() {
                    continue
                }

                if piece.color != color {
                    return true
                }
            }
        }
    }

    return false
}

func (b *SimpleBoard) getPieceLocations() [][]Point {
    return b.pieceLocations
}

func (b *SimpleBoard) copy() (*SimpleBoard, error) {
    simpleBoard, err := newSimpleBoard(b.size, b.players)
    if err != nil {
        return nil, err
    }

    for color, disabled := range b.playersDisabled {
        simpleBoard.playersDisabled[color] = disabled
    }

    for y, row := range b.pieces {
        for x := range row {
            piece := b.pieces[y][x]
            if piece.valid() {
                simpleBoard.pieces[y][x] = piece
            }
        }
    }

    for y, row := range b.disabledLocations {
        for x := range row {
            if b.disabledLocations[y][x] {
                simpleBoard.disabledLocations[y][x] = true
            }
        }
    }

    for color, kingLocation := range b.kingLocations {
        simpleBoard.kingLocations[color] = kingLocation
    }

    for color, pieceLocations := range b.pieceLocations {
        dst := make([]Point, len(pieceLocations))
        copy(dst, pieceLocations)
        simpleBoard.pieceLocations[color] = dst
    }

    for color, enPassant := range b.enPassants {
        simpleBoard.enPassants[color] = enPassant
    }

    for color, vulnerable := range b.vulnerables {
        simpleBoard.vulnerables[color] = vulnerable
    }

    return simpleBoard, nil
}

func (b *SimpleBoard) Copy() (Board, error) {
    board, err := b.copy()
    return board, err
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

func (b *SimpleBoard) pointOutOfBounds(p Point) bool {
    return p.y < 0 || p.y >= b.size.y || p.x < 0 || p.x >= b.size.x || b.disabledLocations[p.y][p.x]
}

func (b *SimpleBoard) colorOutOfBounds(color int) bool {
    return color < 0 || color >= b.players
}

