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
				if err := board.LoadFen(fen); err != nil {
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
			if err := board.LoadFen(fen); err != nil {
				b.Fatalf("init fen: %v", err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				board.GenerateLegalMoves()
			}
		})
	}
}

func BenchmarkLoadPGN(b *testing.B) {
	b.ReportAllocs()
	shortPGN := "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3"

	for i := 0; i < b.N; i++ {
		board := &Board{}
		if err := board.LoadPGN(shortPGN); err != nil {
			b.Fatalf("load pgn: %v", err)
		}
	}
}

func BenchmarkLoadPGN_LongGame(b *testing.B) {
	b.ReportAllocs()
	longPGN := "1. d4 Nf6 2. c4 e6 3. Nf3 d5 4. Nc3 Be7 5. Bg5 h6 6. Bh4 O-O 7. e3 b6 8. cxd5 exd5 9. Bd3 c5 10. O-O Nc6 11. Rc1 Be6 12. Qa4 Nb4 13. Bb1 a6 14. a3 b5 15. Qd1 cxd4 16. axb4 dxc3 17. bxc3 Ne4 18. Bxe7 Qxe7 19. Nd4 Rfc8 20. Bxe4 dxe4 21. Qh5 Rc4 22. f3 exf3 23. Qxf3 Re8 24. e4 Bc8 25. Rce1 Bb7"

	for i := 0; i < b.N; i++ {
		board := &Board{}
		if err := board.LoadPGN(longPGN); err != nil {
			b.Fatalf("load pgn: %v", err)
		}
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
	if err := board.LoadFen(fen); err != nil {
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
		if err := board.LoadFen(fen); err != nil {
			b.Fatalf("init fen: %v", err)
		}
		board.GenerateLegalMoves()
		b.StartTimer()
		_ = perft(board, depth)
	}
}
