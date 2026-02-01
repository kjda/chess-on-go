package chessongo

import "testing"

func BenchmarkInitializeBoardFromFEN(b *testing.B) {
	b.ReportAllocs()
	fens := []string{
		STARTING_POSITION_FEN,
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/8/8/8/8/8/8/8 w - - 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	}

	for _, fen := range fens {
		fen := fen
		b.Run(fen, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				board := &Board{}
				if err := board.InitFromFen(fen); err != nil {
					b.Fatalf("init fen: %v", err)
				}
			}
		})
	}
}

func BenchmarkGenerateLegalMoves(b *testing.B) {
	fens := []string{
		STARTING_POSITION_FEN,
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/8/8/8/8/8/8/8 w - - 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	}

	for _, fen := range fens {
		fen := fen
		b.Run(fen, func(b *testing.B) {
			board := &Board{}
			if err := board.InitFromFen(fen); err != nil {
				b.Fatalf("init fen: %v", err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				board.GenerateLegalMoves()
			}
		})
	}
}

// Benchmarks to get a quick feel for move generation speed via perft.
func BenchmarkPerft(b *testing.B) {
	positions := []struct {
		name  string
		fen   string
		depth int
	}{
		{
			name:  "InitialDepth3",
			fen:   STARTING_POSITION_FEN,
			depth: 3,
		},
		{
			name:  "Position2Depth3",
			fen:   "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth: 3,
		},
	}

	for _, p := range positions {
		p := p
		b.Run(p.name, func(b *testing.B) {
			benchmarkPerft(b, p.fen, p.depth)
		})
	}
}

func benchmarkPerft(b *testing.B, fen string, depth int) {
	// Warm up once to set bytes and validate the position.
	board := &Board{}
	if err := board.InitFromFen(fen); err != nil {
		b.Fatalf("init fen: %v", err)
	}
	board.GenerateLegalMoves()
	baseNodes := perft(board, depth)
	b.SetBytes(int64(baseNodes))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		board = &Board{}
		if err := board.InitFromFen(fen); err != nil {
			b.Fatalf("init fen: %v", err)
		}
		board.GenerateLegalMoves()
		b.StartTimer()
		_ = perft(board, depth)
	}
}
