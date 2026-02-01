package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinaryEncoding(t *testing.T) {
	g := NewGame()
	// Make some moves to populate history and complex state
	err := g.LoadPGN("1. e4 e5 2. Nf3 Nc6 3. Bb5 a6")
	require.NoError(t, err)

	data, err := g.MarshalBinary()
	require.NoError(t, err)

	b2 := &Game{}
	err = b2.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, g.Turn, b2.Turn)
	require.Equal(t, g.Castling, b2.Castling)
	require.Equal(t, g.EnPassant, b2.EnPassant)
	require.Equal(t, g.HalfMoves, b2.HalfMoves)
	require.Equal(t, g.FullMoves, b2.FullMoves)
	require.Equal(t, g.ZobristHash, b2.ZobristHash)
	require.Equal(t, g.Occupation(), b2.Occupation())
	require.Equal(t, len(g.PositionHistory), len(b2.PositionHistory))

	for k, v := range g.PositionHistory {
		require.Equal(t, v, b2.PositionHistory[k], "History count mismatch for hash %x", k)
	}

	// Verify move generation is consistent
	g.GenerateLegalMoves()
	b2.GenerateLegalMoves()
	require.Equal(t, len(g.LegalMoves), len(b2.LegalMoves))
	for i := range g.LegalMoves {
		require.Equal(t, g.LegalMoves[i], b2.LegalMoves[i])
	}
}

func (g *Game) Occupation() uint64 {
	return uint64(g.Occupied)
}
