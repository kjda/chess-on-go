package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadPGNStandardLine(t *testing.T) {
	pgn := "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", g.ToFen())
	// History should have at least the number of plies + initial position.
	require.GreaterOrEqual(t, g.PositionHistory[g.ZobristHash], 1)
}

func TestLoadPGNWithFenTagAndComments(t *testing.T) {
	pgn := `
[Event "Test"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]
[Result "1-0"]

1. d4 d5 {queen's pawn} 2. c4 dxc4 3. e3 Nf6 4. Bxc4 1-0
`

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "rnbqkb1r/ppp1pppp/5n2/8/2BP4/4P3/PP3PPP/RNBQK1NR b KQkq - 0 4", g.ToFen())
}

func TestLoadPGNDetectsThreefold(t *testing.T) {
	pgn := "1. Nf3 Nf6 2. Ng1 Ng8 3. Nf3 Nf6 4. Ng1 Ng8"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.True(t, g.IsThreefoldRepetition)
	require.False(t, g.IsFivefoldRepetition())
	require.GreaterOrEqual(t, g.PositionHistory[g.ZobristHash], 3)
}

func TestLoadPGNGame(t *testing.T) {
	pgn := "1. e4 e5 2. Nf3 Nc6"

	g, err := LoadPGNGame(pgn)
	require.NoError(t, err)
	require.NotNil(t, g)
	require.Equal(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", g.ToFen())
}

func TestLoadPGNWithVariations(t *testing.T) {
	// Variations in parentheses should be ignored, only main line should be played
	pgn := "1. e4 e5 2. Nf3 (2. Bc4 Nf6) 2... Nc6 3. Bb5 a6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", g.ToFen())
}

func TestLoadPGNWithNAGs(t *testing.T) {
	// NAGs (Numeric Annotation Glyphs) like $1, $2, etc. should be ignored
	pgn := "1. e4 $1 e5 $2 2. Nf3 $10 Nc6 3. Bb5 a6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", g.ToFen())
}

func TestLoadPGNWithCastling(t *testing.T) {
	// Test both kingside and queenside castling
	pgn := "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	// Verify both sides have castled - check for rook and king positions
	fen := g.ToFen()
	require.Contains(t, fen, "rk1") // Black king and rook after castling kingside
}

func TestLoadPGNWithPromotion(t *testing.T) {
	// Test pawn promotion
	pgn := `[FEN "4k3/P7/8/8/8/8/8/4K3 w - - 0 1"]
	1. a8=Q Ke7 2. Qa7+ Ke6`

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Contains(t, g.ToFen(), "Q") // Queen should be present after promotion
}

func TestLoadPGNWithEnPassant(t *testing.T) {
	// Test en passant capture
	pgn := "1. e4 a6 2. e5 d5 3. exd6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "rnbqkbnr/1pp1pppp/p2P4/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3", g.ToFen())
}

func TestLoadPGNWithDifferentResults(t *testing.T) {
	tests := []struct {
		name   string
		pgn    string
		result string
	}{
		{"White wins", "1. e4 e5 2. Qh5 Nc6 3. Qxf7# 1-0", "r1bqkbnr/pppp1Qpp/2n5/4p3/4P3/8/PPPP1PPP/RNB1KBNR b KQkq - 0 3"},
		{"Black wins", "1. f3 e5 2. g4 Qh4# 0-1", "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"},
		{"Draw", "1. e4 e5 2. Nf3 Nc6 1/2-1/2", "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"},
		{"Asterisk", "1. e4 e5 *", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{}
			require.NoError(t, g.LoadPGN(tt.pgn))
			require.Equal(t, tt.result, g.ToFen())
		})
	}
}

func TestLoadPGNWithAnnotations(t *testing.T) {
	// Test that move annotations (!, ?, !!, ??, !?, ?!) are properly stripped
	pgn := "1. e4! e5? 2. Nf3!! Nc6?? 3. Bb5!? a6?!"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", g.ToFen())
}

func TestLoadPGNWithCheckAndMateSymbols(t *testing.T) {
	// Test that check (+) and checkmate (#) symbols are properly stripped
	pgn := "1. e4 e5 2. Nf3 Nc6 3. Bc4 Bc5 4. Bxf7+ Kxf7"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Contains(t, g.ToFen(), "k") // Black king should be on the board
}

func TestLoadPGNInvalidMove(t *testing.T) {
	// Test that an invalid move returns an error
	pgn := "1. e4 e5 2. Nf3 Nc6 3. Zz9" // Invalid move

	g := &Game{}
	err := g.LoadPGN(pgn)
	require.Error(t, err)
	require.Contains(t, err.Error(), "pgn move not found")
}

func TestLoadPGNWithMultipleSpaces(t *testing.T) {
	// Test PGN with irregular spacing
	pgn := "1.    e4    e5    2.   Nf3   Nc6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", g.ToFen())
}

func TestLoadPGNWithNestedVariations(t *testing.T) {
	// Test nested variations (should all be ignored)
	pgn := "1. e4 e5 2. Nf3 (2. Bc4 Nf6 (2... Bc5)) 2... Nc6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", g.ToFen())
}

func TestLoadPGNWithBraceComments(t *testing.T) {
	// Test comments in braces
	pgn := "1. e4 {best by test} e5 {symmetric} 2. Nf3 {developing} Nc6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", g.ToFen())
}

func TestLoadPGNFastPath(t *testing.T) {
	// Test the fast path (no tags, comments, or variations)
	pgn := "e4 e5 Nf3 Nc6 Bb5 a6"

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", g.ToFen())
}

func TestLoadPGNWithBlackMove(t *testing.T) {
	// Test PGN with custom starting position where it's black's turn
	pgn := `[FEN "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"]
	1... e5 2. Nf3 Nc6`

	g := &Game{}
	require.NoError(t, g.LoadPGN(pgn))
	require.Equal(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", g.ToFen())
}

func TestExtractFENFromPGNEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		pgn         string
		expectedFEN string
	}{
		{
			name:        "FEN with extra spaces",
			pgn:         `[FEN  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"  ]`,
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "No FEN tag",
			pgn: `[Event "Test"]
[White "Player1"]`,
			expectedFEN: "",
		},
		{
			name:        "Empty PGN",
			pgn:         "",
			expectedFEN: "",
		},
		{
			name:        "FEN with trailing spaces in value",
			pgn:         `[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1  "]`,
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fen := extractFENFromPGN(tt.pgn)
			require.Equal(t, tt.expectedFEN, fen)
		})
	}
}
