package chessongo

import (
	"fmt"
	"testing"
)

func Test_IsCoordsOutofBoard(t *testing.T) {
	for _, coords := range [][]int{{0, 0}, {4, 3}, {7, 7}} {
		if IsCoordsOutofBoard(coords[0], coords[1]) == true {
			t.Error(fmt.Sprintf("IsCoordsOutofBoard failed for %d,%d", coords[0], coords[1]))
		}
	}

	for _, coords := range [][]int{{-1, 0}, {0, -1}, {8, 7}, {1, 8}} {
		if IsCoordsOutofBoard(coords[0], coords[1]) == false {
			t.Error(fmt.Sprintf("IsCoordsOutofBoard failed for %d,%d", coords[0], coords[1]))
		}
	}
}

func Test_CoordsToIndex(t *testing.T) {
	var idx uint
	r, f := 0, 0
	idx = 0
	if idxRes := CoordsToIndex(r, f); idxRes != idx {
		t.Error(fmt.Sprintf("CoordsToIndex failed for R: %d, F: %d - Expexted: %d, Got: %d", r, f, idx, idxRes))
	}
	r, f = 7, 7
	idx = 63
	if idxRes := CoordsToIndex(r, f); idxRes != idx {
		t.Error(fmt.Sprintf("CoordsToIndex failed for R: %d, F: %d - Expexted: %d, Got: %d", r, f, idx, idxRes))
	}
}

func Test_CoordsToSquare(t *testing.T) {
	r, f := 0, 0
	expectedSq := Square(0)
	if sq := CoordsToSquare(r, f); sq != expectedSq {
		t.Error(fmt.Sprintf("CoordsToSquare failed for R: %d, F: %d - Expexted: %d, Got: %d", r, f, expectedSq, sq))
	}
	r, f = 7, 7
	expectedSq = Square(63)
	if sq := CoordsToSquare(r, f); sq != expectedSq {
		t.Error(fmt.Sprintf("CoordsToSquare failed for R: %d, F: %d - Expexted: %d, Got: %d", r, f, expectedSq, sq))
	}
}

func Test_SquareCoords(t *testing.T) {
	sq := Square(0)
	expectedRank, expectedFile := 0, 0

	if r, f := squareCoords(sq); r != expectedRank || f != expectedFile {
		t.Error(fmt.Sprintf("squareCoords failed Square: 0, Expected: (%d, %d) Got(%d, %d)", expectedRank, expectedFile, r, f))
	}
	sq = Square(63)
	expectedRank, expectedFile = 7, 7
	if r, f := squareCoords(sq); r != expectedRank || f != expectedFile {
		t.Error(fmt.Sprintf("squareCoords failed Square: 0, Expected: (%d, %d) Got(%d, %d)", expectedRank, expectedFile, r, f))
	}
}
