package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func perft(b *Board, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var nodes uint64
	// Make a copy of moves to avoid issues with shared slices if any,
	// though Clone() should handle it.
	moves := make([]Move, len(b.LegalMoves))
	copy(moves, b.LegalMoves)

	for _, m := range moves {
		nb := CloneBoard(b)
		nb.MakeMove(m)
		nodes += perft(&nb, depth-1)
	}
	return nodes
}

// *
func TestPerftInitialPosition(t *testing.T) {
	b := &Board{}
	err := b.LoadFen(STARTING_POSITION_FEN)
	require.NoError(t, err)
	b.GenerateLegalMoves()

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 20},
		{2, 400},
		{3, 8902},
		{4, 197281},
		{5, 4865609},
	}

	for _, tt := range tests {
		nodes := perft(b, tt.depth)
		require.Equalf(t, tt.expected, nodes, "Perft(initial, %d)", tt.depth)
	}
}

func TestPerftPosition2(t *testing.T) {
	// Kiwipete
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	b := &Board{}
	require.NoError(t, b.LoadFen(fen))
	b.GenerateLegalMoves()

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 48},
		{2, 2039},
		{3, 97862},
		{4, 4085603},
	}

	for _, tt := range tests {
		nodes := perft(b, tt.depth)
		require.Equalf(t, tt.expected, nodes, "Perft(pos2, %d)", tt.depth)
	}
}

func TestPerftPosition3(t *testing.T) {
	fen := "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"
	b := &Board{}
	require.NoError(t, b.LoadFen(fen))
	b.GenerateLegalMoves()

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 14},
		{2, 191},
		{3, 2812},
		{4, 43238},
		{5, 674624},
	}

	for _, tt := range tests {
		nodes := perft(b, tt.depth)
		require.Equalf(t, tt.expected, nodes, "Perft(pos3, %d)", tt.depth)
	}
}

func TestPerftPosition4(t *testing.T) {
	fen := "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	b := &Board{}
	require.NoError(t, b.LoadFen(fen))
	b.GenerateLegalMoves()

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 6},
		{2, 264},
		{3, 9467},
		{4, 422333},
	}

	for _, tt := range tests {
		nodes := perft(b, tt.depth)
		require.Equalf(t, tt.expected, nodes, "Perft(pos4, %d)", tt.depth)
	}
}

func TestPerftPosition5(t *testing.T) {
	fen := "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
	b := &Board{}
	require.NoError(t, b.LoadFen(fen))
	b.GenerateLegalMoves()

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 44},
		{2, 1486},
		{3, 62379},
		{4, 2103487},
	}

	for _, tt := range tests {
		nodes := perft(b, tt.depth)
		require.Equalf(t, tt.expected, nodes, "Perft(pos5, %d)", tt.depth)
	}
}

func TestPerftPosition6(t *testing.T) {
	fen := "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	b := &Board{}
	require.NoError(t, b.LoadFen(fen))
	b.GenerateLegalMoves()

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 46},
		{2, 2079},
		{3, 89890},
		{4, 3894594},
	}

	for _, tt := range tests {
		nodes := perft(b, tt.depth)
		require.Equalf(t, tt.expected, nodes, "Perft(pos6, %d)", tt.depth)
	}
}

//*/
