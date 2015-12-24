package chessongo

import (
	"testing"
)

func Test_NewBoard(t *testing.T) {
	b := NewBoard()
	if b.Turn != WHITE {
		t.Error("Turn should be white")
	}
}

func Test_HasMoves(t *testing.T) {
	b := NewBoard()
	if b.hasMoves() == true {
		t.Error("Should not have moves")
	}
	b.GenerateLegalMoves()
	if b.hasMoves() != true {
		t.Error("Should have moves")
	}
}
