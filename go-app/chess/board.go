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
    getVulnerables(color int) ([]Point, error) // if these locations are attacked, the player is in check
    setVulnerables(color int, locations []Point) error
	getEnPassant(color int) (EnPassant, error) // if these locations are attacked, a piece is captured en passant
	setEnPassant(color int, enPassant EnPassant) error
    possibleEnPassant(color int, location Point) ([]EnPassant, error)

    // these are for the playerTransition
    disablePieces(color int, disable bool) error

    // these are for the bot
    getPieceLocations() [][]Point

    // these are for the game
    CalculateMoves() error // calcutes moves for every color
    MovesOfColor(color int) ([]Move, error) // returns moves
    MovesOfLocation(fromLocation Point) ([]Move, error) // returns moves
    LegalMovesOfColor(color int) ([]Move, error) // calculates and returns moves that do not result in check
    LegalMovesOfLocation(fromLocation Point) ([]Move, error) // calculates and returns moves that do not result in check
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

    vulnerableLocations := make([][]Point, numberOfPlayers)
    for v := range vulnerableLocations {
        vulnerableLocations[v] = []Point{}
    }

    fromMoves := make([][]Array100[Move], boardSize.y)
    for row := range fromMoves {
        fromMoves[row] = make([]Array100[Move], boardSize.x)
        for col := range fromMoves[row] {
            fromMoves[row][col] = Array100[Move]{}
        }
    }

    toMoves := make([][]Array100[Move], boardSize.y)
    for row := range toMoves {
        toMoves[row] = make([]Array100[Move], boardSize.x)
        for col := range toMoves[row] {
            toMoves[row][col] = Array100[Move]{}
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
        vulnerableLocations: vulnerableLocations,
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
    vulnerableLocations [][]Point
    fromMoves [][]Array100[Move]
    toMoves [][]Array100[Move]
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
        return nil, false
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

    if oldPiece != nil {
        if b.colorOutOfBounds(oldPiece.getColor()) {
            return false
        }

        pieceLocations := b.pieceLocations[oldPiece.getColor()]

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
            b.pieceLocations[oldPiece.getColor()] = pieceLocations
        }
    }

    if p != nil {
        if b.colorOutOfBounds(p.getColor()) {
            return false
        }

        pieceLocations := b.pieceLocations[p.getColor()]

        pieceLocations = append(pieceLocations, location)
        b.pieceLocations[p.getColor()] = pieceLocations

        if _, ok := p.(*King); ok {
            b.kingLocations[p.getColor()] = location
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

func (b *SimpleBoard) getVulnerables(color int) ([]Point, error) {
    if b.colorOutOfBounds(color) {
        return []Point{}, fmt.Errorf("invalid color")
    }

    return b.vulnerableLocations[color], nil
}

func (b *SimpleBoard) setVulnerables(color int, locations []Point) error {
    if b.colorOutOfBounds(color) {
        return fmt.Errorf("invalid color")
    }

    b.vulnerableLocations[color] = locations
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

func (b *SimpleBoard) MovesOfColor(color int) ([]Move, error) {
    if b.colorOutOfBounds(color) {
        return []Move{}, fmt.Errorf("invalid color")
    }

    moves := []Move{}
    for _, pieceLocation := range b.pieceLocations[color] {
        for i := 0; i < b.fromMoves[pieceLocation.y][pieceLocation.x].count; i++ {
            move := b.fromMoves[pieceLocation.y][pieceLocation.x].array[i]
            moves = append(moves, move)
        }
    }
    return moves, nil
}

func (b *SimpleBoard) MovesOfLocation(fromLocation Point) ([]Move, error) {
    if b.pointOutOfBounds(fromLocation) {
        return []Move{}, fmt.Errorf("invalid location")
    }

    moves := []Move{}
    for i := 0; i < b.fromMoves[fromLocation.y][fromLocation.x].count; i++ {
        move := b.fromMoves[fromLocation.y][fromLocation.x].array[i]
        moves = append(moves, move)
    }
    return moves, nil
}

func (b *SimpleBoard) LegalMovesOfColor(color int) ([]Move, error) {
    moves, err := b.MovesOfColor(color)
    if err != nil {
        return []Move{}, err
    }

    legalMoves := []Move{}
    for _, move := range moves {
        if _, ok := move.(*AllyDefenseMove); ok {
            continue
        }

        err = move.execute()
        if err != nil {
            return []Move{}, err
        }

        err = b.CalculateMoves()
        if err != nil {
            return []Move{}, err
        }

        if !b.Check(color) {
            legalMoves = append(legalMoves, move)
        }

        err = move.undo()
        if err != nil {
            return []Move{}, err
        }
    }

    err = b.CalculateMoves()
    if err != nil {
        return []Move{}, err
    }

    return legalMoves, nil
}

func (b *SimpleBoard) LegalMovesOfLocation(fromLocation Point) ([]Move, error) {
    piece, ok := b.getPiece(fromLocation)
    if piece == nil || !ok {
        return []Move{}, fmt.Errorf("piece not found")
    }
    color := piece.getColor()

    moves, err := b.MovesOfLocation(fromLocation)
    if err != nil {
        return []Move{}, err
    }

    legalMoves := []Move{}
    for _, move := range moves {
        if _, ok := move.(*AllyDefenseMove); ok {
            continue
        }

        err = move.execute()
        if err != nil {
            return []Move{}, err
        }

        err = b.CalculateMoves()
        if err != nil {
            return []Move{}, err
        }

        if !b.Check(color) {
            legalMoves = append(legalMoves, move)
        }

        err = move.undo()
        if err != nil {
            return []Move{}, err
        }
    }

    err = b.CalculateMoves()
    if err != nil {
        return []Move{}, err
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

    for row := range b.toMoves {
        for col := range b.toMoves[row] {
            b.toMoves[row][col].clear()
            b.fromMoves[row][col].clear()
        }
    }

    for _, pieceLocations := range b.pieceLocations {
        for _, fromLocation := range pieceLocations {
            piece, ok := b.getPiece(fromLocation)
            if piece == nil || !ok {
                continue
            }

            if b.playersDisabled[piece.getColor()] {
                continue
            }

            moves := piece.moves(b, fromLocation)
            for _, move := range moves {
                action := move.getAction()
                b.fromMoves[action.fromLocation.y][action.fromLocation.x].append(move)
                b.toMoves[action.toLocation.y][action.toLocation.x].append(move)
            }
        }
    }

    return nil
}

// TODO implement dynamic move calculations based on previous move
// TODO how about we don't create massive move objects with pieces and stuff
// stop excessive use of maps
// stop excessive use of pointers
// don't store the moves in the board maybe
// reduce calls to append maybe
func (b *SimpleBoard) CalculateMovesPartial(move Move) error {
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
            } else if piece == nil {
				builder.WriteString(fmt.Sprintf("|%s", strings.Repeat(" ", cellWidth)))
			} else {
				p := piece.print()
				if len(p) > 1 {
					p = p[:1]
				}

                pColor := fmt.Sprintf("%d", piece.getColor())
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
            if piece != nil {
                disabled := b.playersDisabled[piece.getColor()]
                pieces = append(pieces, &PieceData{
                    T: piece.print(),
                    C: piece.getColor(),
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
        action := move.getAction()
        piece, ok := b.getPiece(action.fromLocation)
        if piece == nil || !ok {
            continue
        }

        if piece.getColor() != color {
            return true
        }
    }

    for _, vulnerableLocation := range b.vulnerableLocations[color] {
        for i := 0; i < b.toMoves[vulnerableLocation.y][vulnerableLocation.x].count; i++ {
            move := b.toMoves[vulnerableLocation.y][vulnerableLocation.x].array[i]
            action := move.getAction()
            piece, ok := b.getPiece(action.fromLocation)
            if piece == nil || !ok {
                continue
            }

            if piece.getColor() != color {
                return true
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
            if piece != nil {
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

    for color, vulnerableLocations := range b.vulnerableLocations {
        dst := make([]Point, len(vulnerableLocations))
        copy(dst, vulnerableLocations)
        simpleBoard.vulnerableLocations[color] = dst
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

            if piece == nil {
                counter += 1
                continue
            }

            if counter > 0 {
                builder.WriteString(fmt.Sprintf("%d", counter))
                counter = 0
            }

            if b.playersDisabled[piece.getColor()] {
                builder.WriteString("d")
                continue
            }

            builder.WriteString(piece.print())
            builder.WriteString(fmt.Sprintf("%d", piece.getColor()))
            if !piece.getMoved() {
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

