package chess

import (
	"fmt"
	"strings"
	"unicode"
)

func NewBoard(players []string, fen string) (*board, error) {
	if len(players) < 2 {
		return nil, fmt.Errorf("not enough players")
	}

	fenSplit := strings.Split(fen, " ")

	if len(fenSplit) != 6 {
		return nil, fmt.Errorf("invalid fen")
	}

	return &board{
		players:       players,
		currentPlayer: 0,
		pieces:        createPiecesFromFen(fenSplit[0]),
	}, nil
}

func createPiecesFromFen(fenRows string) [][]piece {
	fenRowsSplit := strings.Split(fenRows, "/")
	pieceRows := [][]piece{}

	for _, row := range fenRowsSplit {
		pieces := []piece{}

		for _, char := range row {
			if unicode.IsDigit(char) {
				for i := 0; i < int(char); i++ {
					pieces = append(pieces, nil)
				}
			} else {
				pieces = append(pieces, createPieceFromChar(char))
			}
		}

		pieceRows = append(pieceRows, pieces)
	}

	return pieceRows
}

func createPieceFromChar(char rune) piece {
	switch char {
	case 'r':
		return &rook{}
	case 'n':
		return &knight{}
	case 'b':
		return &bishop{}
	case 'q':
		return &queen{}
	case 'k':
		return &king{}
	case 'p':
		return &pawn{false, 0, 1}
	case 'R':
		return &rook{}
	case 'N':
		return &knight{}
	case 'B':
		return &bishop{}
	case 'Q':
		return &queen{}
	case 'K':
		return &king{}
	case 'P':
		return &pawn{false, 0, -1}
	default:
		return nil
	}
}

type board struct {
	players       []string
	currentPlayer int
	pieces        [][]piece
}
