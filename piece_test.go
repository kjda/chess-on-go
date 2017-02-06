package chessongo

import (
	"testing"
)

type TestRecord struct {
	piece   Piece
	color   Color
	kind    Piece
	isWhite bool
	r       rune
}

func Test_Consts(t *testing.T) {
	if WHITE == BLACK {
		t.Error("WHITE should not equal BLACK")
	}
	if WHITE == NO_COLOR || BLACK == NO_COLOR {
		t.Error("NO_COLOR should be different from white and black")
	}
	pieces := []Piece{KING, QUEEN, ROOK, BISHOP, KNIGHT, PAWN}

	for _, p1 := range pieces {
		count := 0
		for _, p2 := range pieces {
			if p1 == p2 {
				count++
			}
		}
		if count != 1 {
			t.Error("Piece consts definition is buggy")
		}
	}
}

func Test_PieceApi(t *testing.T) {
	testData := []TestRecord{
		{W_KING, WHITE, KING, true, 'K'},
		{W_QUEEN, WHITE, QUEEN, true, 'Q'},
		{W_ROOK, WHITE, ROOK, true, 'R'},
		{W_BISHOP, WHITE, BISHOP, true, 'B'},
		{W_KNIGHT, WHITE, KNIGHT, true, 'N'},
		{W_PAWN, WHITE, PAWN, true, 'P'},

		{B_KING, BLACK, KING, false, 'k'},
		{B_QUEEN, BLACK, QUEEN, false, 'q'},
		{B_ROOK, BLACK, ROOK, false, 'r'},
		{B_BISHOP, BLACK, BISHOP, false, 'b'},
		{B_KNIGHT, BLACK, KNIGHT, false, 'n'},
		{B_PAWN, BLACK, PAWN, false, 'p'},
	}
	for _, record := range testData {
		if record.piece.Color() != record.color {
			t.Error("Piece color not working")
		}
		if record.piece.Kind() != record.kind {
			t.Error("Piece kind not working")
		}
		if record.piece.IsWhite() != record.isWhite {
			t.Error("Piece isWhite not working")
		}
		if record.piece.IsBlack() == record.isWhite {
			t.Error("Piece isBlack not working")
		}
		if record.piece.ToRune() != record.r {
			t.Error("Piece toRune not working")
		}
	}
}
