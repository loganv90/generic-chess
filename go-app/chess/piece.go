package chess

type CastleDirection struct {
	Point
	kingOffset *Point
	rookOffset *Point
}

type Allegiant struct {
	color string
}

func (a *Allegiant) getColor() string {
	return a.color
}

type Piece interface {
	getColor() string
	copy() Piece
    setMoved() error
	moves(Board, *Point) []Move
	print() string
}

func addDirection(
    fromLocation *Point,
	b Board,
	moves *[]Move,
    direction *Point,
	color string,
) {
    currentLocation := fromLocation.add(direction)

    for {
        p, err := b.getPiece(currentLocation)
        if err != nil {
            break
        }

		if p == nil {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, currentLocation)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
		} else if p.getColor() != color {
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, currentLocation)
			if err == nil {
				*moves = append(*moves, simpleMove)
			}
			break
		} else {
			break
		}

        currentLocation = currentLocation.add(direction)
	}
}

func addSimple(
    fromLocation *Point,
	b Board,
	moves *[]Move,
    direction *Point,
	color string,
) {
    toLocation := fromLocation.add(direction)

	p, err := b.getPiece(toLocation)
	if err != nil {
		return
	}

	if p == nil {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	} else if p.getColor() != color {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
		if err == nil {
			*moves = append(*moves, simpleMove)
		}
	}
}

func newPawn(color string, moved bool, xDir int, yDir int) *Pawn {
    var forward1 *Point
    var forward2 *Point
    var captures []*Point

	if xDir == 1 || xDir == -1 {
		forward1 = &Point{xDir, 0}
		forward2 = &Point{xDir * 2, 0}
		captures = []*Point{
			{xDir, 1},
			{xDir, -1},
		}
	} else if yDir == 1 || yDir == -1 {
		forward1 = &Point{0, yDir}
		forward2 = &Point{0, yDir * 2}
		captures = []*Point{
			{1, yDir},
			{-1, yDir},
		}
	} else {
        panic("invalid direction")
	}

	return &Pawn{
		Allegiant{color},
		moved,
		forward1,
		forward2,
		captures,
        &Point{xDir, yDir},
	}
}

type Pawn struct {
	Allegiant
	moved bool
	forward1 *Point
	forward2 *Point
	captures []*Point
    direction *Point
}

func (a *Pawn) print() string {
	return "P"
}

func (a *Pawn) copy() Piece {
	return &Pawn{
		Allegiant{a.color},
		a.moved,
		a.forward1,
		a.forward2,
		a.captures,
        a.direction,
	}
}

func (a *Pawn) setMoved() error {
    a.moved = true
    return nil
}

func (a *Pawn) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	a.addForward(b, fromLocation, moves)
	a.addCaptures(b, fromLocation, moves)
	return *moves
}

func (a *Pawn) nextLocationInvalid(b Board, toLocation *Point) bool {
    nextLocation := toLocation.add(a.direction)
    _, err := b.getPiece(nextLocation)
    return err != nil
}

func (a *Pawn) appendPromotionMoves(move Move, moves *[]Move) {
    promotionMoves, err := moveFactoryInstance.newPromotionMoves(
        move,
        []Piece{
            newQueen(a.color),
            newRook(a.color, true),
            newBishop(a.color),
            newKnight(a.color),
        },
    )
    if err != nil {
        return
    }
    for _, promotionMove := range promotionMoves {
        *moves = append(*moves, promotionMove)
    }
}

func (a *Pawn) addForward(b Board, fromLocation *Point, moves *[]Move) {
    to1Location := fromLocation.add(a.forward1)
	if piece, err := b.getPiece(to1Location); err != nil || piece != nil { // if the square is invalid or occupied
		return
	} else {
		simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, to1Location)
        if err != nil {
            return
        } else if a.nextLocationInvalid(b, to1Location) {
            a.appendPromotionMoves(simpleMove, moves)
        } else {
			*moves = append(*moves, simpleMove)
		}
	}

	if a.moved {
		return
	}

    to2Location := fromLocation.add(a.forward2)
	if piece, err := b.getPiece(to2Location); err != nil || piece != nil { // if the square is invalid or occupied
		return
	} else {
		revealEnPassantMove, err := moveFactoryInstance.newRevealEnPassantMove(b, fromLocation, to2Location, to1Location)
        if err != nil {
            return
        } else if a.nextLocationInvalid(b, to2Location) {
            a.appendPromotionMoves(revealEnPassantMove, moves)
        } else {
			*moves = append(*moves, revealEnPassantMove)
		}
	}
}

func (a *Pawn) addCaptures(b Board, fromLocation *Point, moves *[]Move) {
	for _, capture := range a.captures {
        toLocation := fromLocation.add(capture)

		if piece, err := b.getPiece(toLocation); err != nil { // if the square is invalid
			continue
        } else if ens, err := b.possibleEnPassant(a.color, toLocation); err == nil && len(ens) > 0 { // if the square is an en passant target
			captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(b, fromLocation, toLocation)
            if err != nil {
                continue
            } else if a.nextLocationInvalid(b, toLocation) {
                a.appendPromotionMoves(captureEnPassantMove, moves)
            } else {
                *moves = append(*moves, captureEnPassantMove)
            }
		} else if piece != nil && piece.getColor() != a.color { // if the square is occupied by an enemy piece
			simpleMove, err := moveFactoryInstance.newSimpleMove(b, fromLocation, toLocation)
            if err != nil {
                continue
            } else if a.nextLocationInvalid(b, toLocation) {
                a.appendPromotionMoves(simpleMove, moves)
            } else {
                *moves = append(*moves, simpleMove)
            }
		}
	}
}

var knightSimples = []*Point{
	{1, 2},
	{-1, 2},
	{2, 1},
	{-2, 1},
	{1, -2},
	{-1, -2},
	{2, -1},
	{-2, -1},
}

func newKnight(color string) *Knight {
	return &Knight{
		Allegiant{color},
	}
}

type Knight struct {
	Allegiant
}

func (n *Knight) print() string {
	return "N"
}

func (n *Knight) copy() Piece {
	return &Knight{
		Allegiant{n.color},
	}
}

func (n *Knight) setMoved() error {
    return nil
}

func (n *Knight) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	n.addSimples(b, fromLocation, moves)
	return *moves
}

func (n *Knight) addSimples(b Board, fromLocation *Point, moves *[]Move) {
	for _, simple := range knightSimples {
		addSimple(fromLocation, b, moves, simple, n.color)
	}
}

var bishopDirections = []*Point{
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newBishop(color string) *Bishop {
	return &Bishop{
		Allegiant{color},
	}
}

type Bishop struct {
	Allegiant
}

func (s *Bishop) print() string {
	return "B"
}

func (s *Bishop) copy() Piece {
	return &Bishop{
		Allegiant{s.color},
	}
}

func (s *Bishop) setMoved() error {
    return nil
}

func (s *Bishop) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	s.addDirections(b, fromLocation, moves)
	return *moves
}

func (s *Bishop) addDirections(b Board, fromLocation *Point, moves *[]Move) {
	for _, direction := range bishopDirections {
		addDirection(fromLocation, b, moves, direction, s.color)
	}
}

var rookDirections = []*Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func newRook(color string, moved bool) *Rook {
	return &Rook{
		Allegiant{color},
		moved,
	}
}

type Rook struct {
	Allegiant
	moved bool
}

func (r *Rook) print() string {
	return "R"
}

func (r *Rook) copy() Piece {
	return &Rook{
		Allegiant{r.color},
		true,
	}
}

func (r *Rook) setMoved() error {
    r.moved = true
    return nil
}

func (r *Rook) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	r.addDirections(b, fromLocation, moves)
	return *moves
}

func (r *Rook) addDirections(b Board, fromLocation *Point, moves *[]Move) {
	for _, direction := range rookDirections {
		addDirection(fromLocation, b, moves, direction, r.color)
	}
}

var queenDirections = []*Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newQueen(color string) *Queen {
	return &Queen{
		Allegiant{color},
	}
}

type Queen struct {
	Allegiant
}

func (q *Queen) print() string {
	return "Q"
}

func (q *Queen) copy() Piece {
	return &Queen{
		Allegiant{q.color},
	}
}

func (q *Queen) setMoved() error {
    return nil
}

func (q *Queen) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	q.addDirections(b, fromLocation, moves)
	return *moves
}

func (q *Queen) addDirections(b Board, fromLocation *Point, moves *[]Move) {
	for _, direction := range queenDirections {
		addDirection(fromLocation, b, moves, direction, q.color)
	}
}

var kingSimples = []*Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, 1},
	{1, -1},
	{-1, -1},
}

func newKing(color string, moved bool, xDir int, yDir int) *King {
	var castles []*CastleDirection

	if xDir == 1 || xDir == -1 {
		castles = []*CastleDirection{
			{Point{0, 1}, &Point{0, -2}, &Point{0, -3}},
			{Point{0, -1}, &Point{0, 1}, &Point{0, 2}},
		}
	} else if yDir == 1 || yDir == -1 {
		castles = []*CastleDirection{
			{Point{1, 0}, &Point{-1, 0}, &Point{-2, 0}},
			{Point{-1, 0}, &Point{2, 0}, &Point{3, 0}},
		}
	} else {
        panic("invalid direction")
	}

	return &King{
		Allegiant{color},
		moved,
		castles,
	}
}

type King struct {
	Allegiant
	moved   bool
	castles []*CastleDirection
}

func (k *King) print() string {
	return "K"
}

func (k *King) copy() Piece {
	return &King{
		Allegiant{k.color},
		k.moved,
		k.castles,
	}
}

func (k *King) setMoved() error {
    k.moved = true
    return nil
}

func (k *King) moves(b Board, fromLocation *Point) []Move {
	moves := &[]Move{}
	k.addSimples(b, fromLocation, moves)
	k.addCastles(b, fromLocation, moves)
	return *moves
}

func (k *King) addSimples(b Board, fromLocation *Point, moves *[]Move) {
	for _, simple := range kingSimples {
		addSimple(fromLocation, b, moves, simple, k.color)
	}
}

func (k *King) addCastles(b Board, fromLocation *Point, moves *[]Move) {
	if k.moved {
		return
	}

	for _, castle := range k.castles {
        vulnerableLocations := []*Point{}
        currentLocation := fromLocation.add(&castle.Point)

		rookFound := false
		for {
			piece, err := b.getPiece(currentLocation)
			if err != nil {
				break
			}

			if piece == nil {
                currentLocation = currentLocation.add(&castle.Point)
				continue
			}

			rook, ok := piece.(*Rook)
			if !ok || rook.moved || rook.color != k.color {
				break
			}

			rookFound = true
			break
		}
		if !rookFound {
			continue
		}

		xToKing := fromLocation.x
		xToRook := currentLocation.x
		yToKing := fromLocation.y
		yToRook := currentLocation.y
		if castle.kingOffset.x > 0 {
			xToKing = castle.kingOffset.x
		} else if castle.kingOffset.x < 0 {
			xToKing = b.Size().x - 1 + castle.kingOffset.x
		}
		if castle.kingOffset.y > 0 {
			yToKing = castle.kingOffset.y
		} else if castle.kingOffset.y < 0 {
			yToKing = b.Size().y - 1 + castle.kingOffset.y
		}
		if castle.rookOffset.x > 0 {
			xToRook = castle.rookOffset.x
		} else if castle.rookOffset.x < 0 {
			xToRook = b.Size().x - 1 + castle.rookOffset.x
		}
		if castle.rookOffset.y > 0 {
			yToRook = castle.rookOffset.y
		} else if castle.rookOffset.y < 0 {
			yToRook = b.Size().y - 1 + castle.rookOffset.y
		}
        if xToKing < 0 || xToKing >= b.Size().x || yToKing < 0 || yToKing >= b.Size().y {
            continue
        }
        if xToRook < 0 || xToRook >= b.Size().x || yToRook < 0 || yToRook >= b.Size().y {
            continue
        }

        xCheckedMin := min(fromLocation.x, currentLocation.x)
        xCheckedMax := max(fromLocation.x, currentLocation.x)
        yCheckedMin := min(fromLocation.y, currentLocation.y)
        yCheckedMax := max(fromLocation.y, currentLocation.y)

        xToMin := min(xToKing, xToRook)
        xToMax := max(xToKing, xToRook)
        yToMin := min(yToKing, yToRook)
        yToMax := max(yToKing, yToRook)

        clear := true
        for x := xCheckedMin; x > xToMin && clear; x-- {
            if piece, err := b.getPiece(&Point{x, fromLocation.y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        for y := yCheckedMin; y > yToMin && clear; y-- {
            if piece, err := b.getPiece(&Point{fromLocation.x, y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        for x := xCheckedMax; x < xToMax && clear; x++ {
            if piece, err := b.getPiece(&Point{x, fromLocation.y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        for y := yCheckedMax; y < yToMax && clear; y++ {
            if piece, err := b.getPiece(&Point{fromLocation.y, y}); err != nil || piece != nil {
                clear = false
                break
            }
        }
        if !clear {
            continue
        }

        if xToKing > fromLocation.x {
            for x := fromLocation.x + 1; x < xToKing; x++ {
                vulnerableLocations = append(vulnerableLocations, &Point{x, fromLocation.y})
            }
        } else if xToKing < fromLocation.x {
            for x := fromLocation.x - 1; x > xToKing; x-- {
                vulnerableLocations = append(vulnerableLocations, &Point{x, fromLocation.y})
            }
        } else if yToKing > fromLocation.y {
            for y := fromLocation.y + 1; y < yToKing; y++ {
                vulnerableLocations = append(vulnerableLocations, &Point{fromLocation.x, y})
            }
        } else if yToKing < fromLocation.y {
            for y := fromLocation.y - 1; y > yToKing; y-- {
                vulnerableLocations = append(vulnerableLocations, &Point{fromLocation.x, y})
            }
        }

		castleMove, err := moveFactoryInstance.newCastleMove(
            b,
            fromLocation,
            currentLocation,
            &Point{xToKing, yToKing},
            &Point{xToRook, yToRook},
            vulnerableLocations,
        )
		if err == nil {
			*moves = append(*moves, castleMove)
		}
	}
}

