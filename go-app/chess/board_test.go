package chess

import (
	"testing"
)

func TestNewBoardWithDefaultFen(t *testing.T) {
	board, err := NewBoard(
		[]string{"white", "black"},
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(board.pieces) != 8 || len(board.pieces[0]) != 8 {
		t.Fatalf("board should have 8 rows and 8 columns")
	}

	for row := 2; row < 6; row++ {
		for col := 0; col < 8; col++ {
			if board.pieces[row][col] != nil {
				t.Fatalf("piece at row %d and column %d should be nil", row, col)
			}
		}
	}

	for _, row := range []int{1, 6} {
		for col := 0; col < 8; col++ {
			if _, ok := board.pieces[row][col].(*pawn); !ok {
				t.Fatalf("piece at row %d and column %d should be a pawn", row, col)
			}
		}
	}

	for _, row := range []int{0, 7} {
		if _, ok := board.pieces[row][0].(*rook); !ok {
			t.Fatalf("piece at row %d and column %d should be a rook", row, 0)
		}
		if _, ok := board.pieces[row][1].(*knight); !ok {
			t.Fatalf("piece at row %d and column %d should be a knight", row, 1)
		}
		if _, ok := board.pieces[row][2].(*bishop); !ok {
			t.Fatalf("piece at row %d and column %d should be a bishop", row, 2)
		}
		if _, ok := board.pieces[row][3].(*queen); !ok {
			t.Fatalf("piece at row %d and column %d should be a queen", row, 3)
		}
		if _, ok := board.pieces[row][4].(*king); !ok {
			t.Fatalf("piece at row %d and column %d should be a queen", row, 4)
		}
		if _, ok := board.pieces[row][5].(*bishop); !ok {
			t.Fatalf("piece at row %d and column %d should be a bishop", row, 5)
		}
		if _, ok := board.pieces[row][6].(*knight); !ok {
			t.Fatalf("piece at row %d and column %d should be a knight", row, 6)
		}
		if _, ok := board.pieces[row][7].(*rook); !ok {
			t.Fatalf("piece at row %d and column %d should be a rook", row, 7)
		}
	}
}
