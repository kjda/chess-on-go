package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoadFen(t *testing.T) {
	g := NewGame()
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQ e5 1 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b - - 1 2",
		"rnbqkbnr/ppp2ppp/8/3pp3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 3",
		"rnb1kbnr/ppp2ppp/8/8/2qp4/5N2/PPP2PPP/RNBQK2R w KQkq - 0 6",
	}
	for _, fen := range fens {
		g.LoadFen(fen)
		require.Equal(t, fen, g.ToFen())
	}

}

func Test_ToFen(t *testing.T) {
	g := NewGame()
	require.Equal(t, STARTING_POSITION_FEN, g.ToFen())
	g.LoadFen("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	require.Equal(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", g.ToFen())
}
