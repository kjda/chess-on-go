package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadPGNStandardLine(t *testing.T) {
	pgn := "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6"

	b := &Board{}
	require.NoError(t, b.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", b.ToFen())
	// History should have at least the number of plies + initial position.
	require.GreaterOrEqual(t, b.PositionHistory[b.ZobristHash], 1)
}

func TestLoadPGNWithFenTagAndComments(t *testing.T) {
	pgn := `
[Event "Test"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]
[Result "1-0"]

1. d4 d5 {queen's pawn} 2. c4 dxc4 3. e3 Nf6 4. Bxc4 1-0
`

	b := &Board{}
	require.NoError(t, b.LoadPGN(pgn))
	require.Equal(t, "rnbqkb1r/ppp1pppp/5n2/8/2BP4/4P3/PP3PPP/RNBQK1NR b KQkq - 0 4", b.ToFen())
}

func TestLoadPGNDetectsThreefold(t *testing.T) {
	pgn := "1. Nf3 Nf6 2. Ng1 Ng8 3. Nf3 Nf6 4. Ng1 Ng8"

	b := &Board{}
	require.NoError(t, b.LoadPGN(pgn))
	require.True(t, b.IsThreefoldRepetition)
	require.False(t, b.IsFivefoldRepetition())
	require.GreaterOrEqual(t, b.PositionHistory[b.ZobristHash], 3)
}
