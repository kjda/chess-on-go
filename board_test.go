package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewBoard(t *testing.T) {
	b := NewBoard()
	require.EqualValues(t, WHITE, b.Turn)
}

func Test_HasMoves(t *testing.T) {
	b := NewBoard()
	require.False(t, b.hasMoves())
	b.GenerateLegalMoves()
	require.True(t, b.hasMoves())
}

func TestRepetitionDetection(t *testing.T) {
	b := NewBoard()
	require.False(t, b.IsThreefoldRepetition)
	require.False(t, b.IsFivefoldRepetition())

	cycle := [][2]string{{"g1", "f3"}, {"g8", "f6"}, {"f3", "g1"}, {"f6", "g8"}}
	play := func(from, to string) {
		fromSq := COORDS_TO_SQUARE[from]
		toSq := COORDS_TO_SQUARE[to]
		b.MakeMove(NewMove(fromSq, toSq, b.Squares[toSq]))
	}

	for i := 0; i < 4; i++ {
		for _, mv := range cycle {
			play(mv[0], mv[1])
		}
		if i == 1 {
			require.True(t, b.IsThreefoldRepetition)
			require.False(t, b.IsFivefoldRepetition())
		}
	}

	require.True(t, b.IsThreefoldRepetition)
	require.True(t, b.IsFivefoldRepetition())
}

func TestRepetitionBrokenByEnPassantChange(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2"
	b := &Board{}
	require.NoError(t, b.LoadFen(fen))
	require.False(t, b.IsThreefoldRepetition)

	cycle := [][2]string{{"g1", "f3"}, {"g8", "f6"}, {"f3", "g1"}, {"f6", "g8"}}
	play := func(from, to string) {
		fromSq := COORDS_TO_SQUARE[from]
		toSq := COORDS_TO_SQUARE[to]
		b.MakeMove(NewMove(fromSq, toSq, b.Squares[toSq]))
	}

	for _, mv := range cycle {
		play(mv[0], mv[1])
	}

	// After a full cycle the board pieces match but en-passant is cleared, so hash differs.
	require.False(t, b.IsThreefoldRepetition)
}

func TestRepetitionBrokenByIrreversibleMove(t *testing.T) {
	b := NewBoard()
	cycle := [][2]string{{"g1", "f3"}, {"g8", "f6"}, {"f3", "g1"}, {"f6", "g8"}}
	play := func(from, to string) {
		fromSq := COORDS_TO_SQUARE[from]
		toSq := COORDS_TO_SQUARE[to]
		b.MakeMove(NewMove(fromSq, toSq, b.Squares[toSq]))
	}

	// Two occurrences of the base position (start + one cycle)
	for _, mv := range cycle {
		play(mv[0], mv[1])
	}
	require.Equal(t, 2, b.PositionHistory[b.ZobristHash])
	require.False(t, b.IsThreefoldRepetition)

	// Make an irreversible pawn move pair, then repeat the cycle once.
	play("e2", "e4")
	play("e7", "e5")
	for _, mv := range cycle {
		play(mv[0], mv[1])
	}

	// Different hash due to pawn shifts; repetitions should not accumulate toward old state.
	require.False(t, b.IsThreefoldRepetition)
	require.False(t, b.IsFivefoldRepetition())
}
