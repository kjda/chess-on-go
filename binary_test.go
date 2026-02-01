package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinaryEncoding(t *testing.T) {
	b := NewBoard()
	// Make some moves to populate history and complex state
	err := b.LoadPGN("1. e4 e5 2. Nf3 Nc6 3. Bb5 a6")
	require.NoError(t, err)

	data, err := b.MarshalBinary()
	require.NoError(t, err)

	b2 := &Board{}
	err = b2.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, b.Turn, b2.Turn)
	require.Equal(t, b.Castling, b2.Castling)
	require.Equal(t, b.EnPassant, b2.EnPassant)
	require.Equal(t, b.HalfMoves, b2.HalfMoves)
	require.Equal(t, b.FullMoves, b2.FullMoves)
	require.Equal(t, b.ZobristHash, b2.ZobristHash)
	require.Equal(t, b.Occupation(), b2.Occupation())
	require.Equal(t, len(b.PositionHistory), len(b2.PositionHistory))

	for k, v := range b.PositionHistory {
		require.Equal(t, v, b2.PositionHistory[k], "History count mismatch for hash %x", k)
	}

	// Verify move generation is consistent
	b.GenerateLegalMoves()
	b2.GenerateLegalMoves()
	require.Equal(t, len(b.LegalMoves), len(b2.LegalMoves))
	for i := range b.LegalMoves {
		require.Equal(t, b.LegalMoves[i], b2.LegalMoves[i])
	}
}

func (b *Board) Occupation() uint64 {
	return uint64(b.Occupied)
}
