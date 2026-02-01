package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewGame(t *testing.T) {
	g := NewGame()
	require.EqualValues(t, WHITE, g.Turn)
}

func Test_HasMoves(t *testing.T) {
	g := NewGame()
	require.False(t, g.hasMoves())
	g.GenerateLegalMoves()
	require.True(t, g.hasMoves())
}

func TestRepetitionDetection(t *testing.T) {
	g := NewGame()
	require.False(t, g.IsThreefoldRepetition)
	require.False(t, g.IsFivefoldRepetition())

	cycle := [][2]string{{"g1", "f3"}, {"g8", "f6"}, {"f3", "g1"}, {"f6", "g8"}}
	play := func(from, to string) {
		fromSq := COORDS_TO_SQUARE[from]
		toSq := COORDS_TO_SQUARE[to]
		g.MakeMove(NewMove(fromSq, toSq, g.Squares[toSq]))
	}

	for i := 0; i < 4; i++ {
		for _, mv := range cycle {
			play(mv[0], mv[1])
		}
		if i == 1 {
			require.True(t, g.IsThreefoldRepetition)
			require.False(t, g.IsFivefoldRepetition())
		}
	}

	require.True(t, g.IsThreefoldRepetition)
	require.True(t, g.IsFivefoldRepetition())
}

func TestRepetitionBrokenByEnPassantChange(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2"
	g := &Game{}
	require.NoError(t, g.LoadFen(fen))
	require.False(t, g.IsThreefoldRepetition)

	cycle := [][2]string{{"g1", "f3"}, {"g8", "f6"}, {"f3", "g1"}, {"f6", "g8"}}
	play := func(from, to string) {
		fromSq := COORDS_TO_SQUARE[from]
		toSq := COORDS_TO_SQUARE[to]
		g.MakeMove(NewMove(fromSq, toSq, g.Squares[toSq]))
	}

	for _, mv := range cycle {
		play(mv[0], mv[1])
	}

	// After a full cycle the board pieces match but en-passant is cleared, so hash differs.
	require.False(t, g.IsThreefoldRepetition)
}

func TestRepetitionBrokenByIrreversibleMove(t *testing.T) {
	g := NewGame()
	cycle := [][2]string{{"g1", "f3"}, {"g8", "f6"}, {"f3", "g1"}, {"f6", "g8"}}
	play := func(from, to string) {
		fromSq := COORDS_TO_SQUARE[from]
		toSq := COORDS_TO_SQUARE[to]
		g.MakeMove(NewMove(fromSq, toSq, g.Squares[toSq]))
	}

	// Two occurrences of the base position (start + one cycle)
	for _, mv := range cycle {
		play(mv[0], mv[1])
	}
	require.Equal(t, 2, g.PositionHistory[g.ZobristHash])
	require.False(t, g.IsThreefoldRepetition)

	// Make an irreversible pawn move pair, then repeat the cycle once.
	play("e2", "e4")
	play("e7", "e5")
	for _, mv := range cycle {
		play(mv[0], mv[1])
	}

	// Different hash due to pawn shifts; repetitions should not accumulate toward old state.
	require.False(t, g.IsThreefoldRepetition)
	require.False(t, g.IsFivefoldRepetition())
}
