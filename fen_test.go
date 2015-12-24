package chessongo

import (
	"testing"
)

func Test_InitFromFen(t *testing.T) {
	b := NewBoard()
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQ e5 1 2",

		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b - - 1 2",
	}
	for _, fen := range fens {
		b.InitFromFen(fen)
		if b.ToFen() != fen {
			t.Error("initFromFen not working")
		}
	}

}

func Test_ToFen(t *testing.T) {
	b := NewBoard()
	if b.ToFen() != STARTING_POSITION_FEN {
		t.Error("Invalid starting position")
	}
}
