package chessongo

import (
	"strings"
	"testing"
)

// Test Pawn Moves: Single push, double push, capture, en passant, promotion
func TestPawnMoves(t *testing.T) {
	tests := []struct {
		fen           string
		description   string
		expectedMoves []string
	}{
		{
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			description: "Starting position white pawns",
			expectedMoves: []string{
				"e2 e3", "e2 e4", // Single and double push
			},
		},
		{
			fen:         "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			description: "Black response to e4",
			expectedMoves: []string{
				"e7 e6", "e7 e5",
			},
		},
		{
			fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			description: "White pawn capture d5",
			expectedMoves: []string{
				"e4 d5", // Capture
				"e4 e5", // Push
			},
		},
		{
			fen:         "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
			description: "White En Passant capture on f6",
			expectedMoves: []string{
				"e5 f6", // En Passant
				"e5 e6", // Push
			},
		},
		{
			fen:         "8/4P3/8/8/8/8/8/8 w - - 0 1",
			description: "White Pawn Promotion",
			expectedMoves: []string{
				"e7 e8", // Will appear as multiple promotion moves in move list logic,
				// but ToString might just be "e7 e8".
				// The engine generates separate moves for promotions.
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGame()
			g.LoadFen(tt.fen)
			g.GenerateLegalMoves()

			var moves []string
			for _, m := range g.LegalMoves {
				moves = append(moves, m.ToString())
			}

			for _, expected := range tt.expectedMoves {
				assertContains(t, moves, expected)
			}
		})
	}
}

func TestPromotionMovesDetails(t *testing.T) {
	// White pawn ready to promote on a8
	fen := "8/P7/8/8/8/8/8/k6K w - - 0 1"
	g := NewGame()
	g.LoadFen(fen)
	g.GenerateLegalMoves()

	// expect 4 promotion moves from a7 to a8
	promotionCount := 0
	for _, m := range g.LegalMoves {
		if m.From() == CoordsToSquare(1, 0) && m.To() == CoordsToSquare(0, 0) && m.IsPromotionMove() {
			promotionCount++
		}
	}

	if promotionCount != 4 {
		t.Errorf("Expected 4 promotion moves, got %d", promotionCount)
	}
}

// Test Castling
func TestCastling(t *testing.T) {
	tests := []struct {
		fen           string
		description   string
		expectedMoves []string
		absentMoves   []string
	}{
		{
			fen:         "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			description: "White full castling rights, clear path",
			expectedMoves: []string{
				"e1 g1", // O-O
				"e1 c1", // O-O-O
			},
		},
		{
			fen:         "r3k2r/8/8/8/8/8/8/R3K2R w - - 0 1",
			description: "White no castling rights",
			absentMoves: []string{
				"e1 g1",
				"e1 c1",
			},
		},
		{
			fen:         "r3k2r/8/8/8/8/8/8/R3K2R w K - 0 1",
			description: "White only Kingside disallowed by rights",
			expectedMoves: []string{
				"e1 g1",
			},
			absentMoves: []string{
				"e1 c1",
			},
		},
		{
			fen:         "4k3/8/8/8/8/8/8/R3K2R w K - 0 1",
			description: "White Kingside, Rook Missing (implicitly handled by board logic which uses rook presence for rights usually, checking rights though)",
			// If FEN says K, logic assumes rights exist. However, move logic checks for rook presence usually?
			// Let's see move-logic.go:93: (g.Occupied&(0x3<<61)) == 0  -- checks empty squares f1, g1
			// It implies King is at e1. It doesn't explicitly check Rook presence in GeneratePseudoMoves if castling bits are set,
			// usually logic assumes if Castling bit is set, Rook is there.
			// BUT standard chess rules say if rook is captured, right is lost.
			// Let's assume FEN is source of truth for rights.
			expectedMoves: []string{
				"e1 g1",
			},
		},
		{
			fen:         "r3k2r/8/8/8/8/4r3/8/R3K2R w KQkq - 0 1",
			description: "White in check, cannot castle",
			absentMoves: []string{
				"e1 g1",
				"e1 c1",
			},
		},
		{
			fen:         "r3k2r/8/8/8/8/5r2/8/R3K2R w KQkq - 0 1",
			description: "White castling through check (f1 attacked)",
			absentMoves: []string{
				"e1 g1", // passes through f1
			},
			expectedMoves: []string{
				"e1 c1", // Queenside safe
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGame()
			g.LoadFen(tt.fen)
			g.GenerateLegalMoves()

			var moves []string
			for _, m := range g.LegalMoves {
				moves = append(moves, m.ToString())
			}

			for _, expected := range tt.expectedMoves {
				assertContains(t, moves, expected)
			}
			for _, absent := range tt.absentMoves {
				assertNotContains(t, moves, absent)
			}
		})
	}
}

// Test Check Evasion
func TestCheckEvasion(t *testing.T) {
	// White King at e1, Black Rook at e8.
	// Valid moves: King moves away (d1, f1, etc), or something blocks.
	fen := "4r3/8/8/8/8/8/4P3/4K3 w - - 0 1"
	g := NewGame()
	g.LoadFen(fen)
	g.GenerateLegalMoves()

	moves := map[string]bool{}
	for _, m := range g.LegalMoves {
		moves[m.ToString()] = true
	}

	// e1e2 is illegal (king moving to attacked square, protected by rook) BUT wait, e8 rook attacks e file.
	// e2 is occupied by white pawn.
	// e1f1 ok. e1d1 ok.
	// e1f2 illegal (checked). e1d2 illegal (checked).

	if moves["e1 e2"] {
		t.Error("King should not capture own piece / move to check")
	}
	if !moves["e1 f1"] {
		t.Error("King should be able to move to f1")
	}
	if !moves["e1 d1"] {
		t.Error("King should be able to move to d1")
	}
	// Pawn at e2 cannot move e2-e3 because it's pinned?
	// No, vertical pin. e2-e3 would leave king in check from e8 rook?
	// Rook at e8. King at e1. Pawn at e2.
	// e2-e3 is LEGAL (moves along the pin ray, maintaining block).
	if !moves["e2 e3"] {
		t.Error("Pinned pawn e2 SHOULD be able to move to e3 (along pin ray)")
	}
	// e2-e4 is LEGAL (moves along pin ray).
	if !moves["e2 e4"] {
		t.Error("Pinned pawn e2 SHOULD be able to move to e4 (along pin ray)")
	}
}

// Test SAN Generation
func TestSAN(t *testing.T) {
	tests := []struct {
		fen         string
		moveFrom    string // e.g. "e2"
		moveTo      string // e.g. "e4"
		promote     string // "" or "Q"
		expectedSan string
	}{
		{
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moveFrom:    "e2",
			moveTo:      "e4",
			expectedSan: "e4",
		},
		{
			fen:         "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
			moveFrom:    "f1",
			moveTo:      "b5",
			expectedSan: "Bb5",
		},
		{
			fen:         "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3",
			moveFrom:    "g8",
			moveTo:      "f6",
			expectedSan: "Nf6",
		},
		// Disambiguation
		{
			fen:         "7k/8/8/8/8/2N1N3/8/K7 w - - 0 1",
			moveFrom:    "c3",
			moveTo:      "d5",
			expectedSan: "Ncd5", // N on c3 and N on e3 can both reach d5. File differs.
		},
		{
			fen:         "7k/8/8/8/8/2N1N3/8/K7 w - - 0 1",
			moveFrom:    "e3",
			moveTo:      "d5",
			expectedSan: "Ned5",
		},
		{
			fen:         "1r5k/8/8/8/8/1r6/8/K7 b - - 0 1",
			moveFrom:    "b8",
			moveTo:      "b4",
			expectedSan: "R8b4", // Rooks on same file, rank differs.
		},
		{
			fen:         "1r5k/8/8/8/8/1r6/8/K7 b - - 0 1",
			moveFrom:    "b3",
			moveTo:      "b4",
			expectedSan: "R3b4",
		},
		// Disambiguation: two queens on different files (file disambiguation)
		{
			fen:         "7k/8/8/8/8/7K/8/3Q1Q2 w - - 0 1",
			moveFrom:    "d1",
			moveTo:      "e2",
			expectedSan: "Qde2", // Qd1 and Qf1 can both reach e2. Different files.
		},
		{
			fen:         "7k/8/8/8/8/7K/8/3Q1Q2 w - - 0 1",
			moveFrom:    "f1",
			moveTo:      "e2",
			expectedSan: "Qfe2",
		},
		// Disambiguation: two queens on same file (rank disambiguation)
		{
			fen:         "7k/8/8/3Q4/8/8/8/3Q3K w - - 0 1",
			moveFrom:    "d1",
			moveTo:      "d3",
			expectedSan: "Q1d3", // Qd1 and Qd5 on same file. Rank differs.
		},
		{
			fen:         "7k/8/8/3Q4/8/8/8/3Q3K w - - 0 1",
			moveFrom:    "d5",
			moveTo:      "d3",
			expectedSan: "Q5d3",
		},
		// Disambiguation: three queens, need full coordinates for one
		{
			// Queens on d1, d5, f1. All can reach d3.
			// d1→d3: d5 shares file, f1 shares rank → need full coords "Qd1d3"
			fen:         "7k/8/8/3Q4/8/7K/8/3Q1Q2 w - - 0 1",
			moveFrom:    "d1",
			moveTo:      "d3",
			expectedSan: "Qd1d3",
		},
		{
			// d5→d3: d1 shares file, f1 does not share file or rank → rank suffices "Q5d3"
			fen:         "7k/8/8/3Q4/8/7K/8/3Q1Q2 w - - 0 1",
			moveFrom:    "d5",
			moveTo:      "d3",
			expectedSan: "Q5d3",
		},
		{
			// f1→d3: neither d1 nor d5 shares file f → file suffices "Qfd3"
			fen:         "7k/8/8/3Q4/8/7K/8/3Q1Q2 w - - 0 1",
			moveFrom:    "f1",
			moveTo:      "d3",
			expectedSan: "Qfd3",
		},
		// Disambiguation: two rooks on same rank (file disambiguation)
		{
			fen:         "7k/8/8/8/8/8/8/R3R2K w - - 0 1",
			moveFrom:    "a1",
			moveTo:      "c1",
			expectedSan: "Rac1", // Ra1 and Re1 on same rank. File differs.
		},
		{
			fen:         "7k/8/8/8/8/8/8/R3R2K w - - 0 1",
			moveFrom:    "e1",
			moveTo:      "c1",
			expectedSan: "Rec1",
		},
		// Captures
		{
			fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			moveFrom:    "e4",
			moveTo:      "d5",
			expectedSan: "exd5",
		},
		// Check
		{
			fen:         "rnbqkbnr/ppp2ppp/8/3pp3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 3",
			moveFrom:    "f1",
			moveTo:      "b5",
			expectedSan: "Bb5+", // Check
		},
		{
			fen:         "rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq - 0 2", // After 1. f3 e5 2. g4 ...
			moveFrom:    "d8",                                                            // Queen
			moveTo:      "h4",
			expectedSan: "Qh4#",
		},
		// Pawn Promotion Capture with correct SAN (User reported bug)
		{
			fen:         "1n6/2P5/8/8/8/8/8/k6K w - - 0 1",
			moveFrom:    "c7",
			moveTo:      "b8",
			promote:     "Q",
			expectedSan: "cxb8=Q", // Capture and promote
		},
		// Pawn Promotion No Capture
		{
			fen:         "8/2P5/8/8/8/8/8/k6K w - - 0 1",
			moveFrom:    "c7",
			moveTo:      "c8",
			promote:     "R",
			expectedSan: "c8=R",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedSan, func(t *testing.T) {
			// Special handling for checkmate scenario setup if needed
			// But for simplicity use the fen provided.
			// For Fool's Mate checkmate, let's fix the FEN to be "before the move"
			fen := tt.fen
			if tt.expectedSan == "Qh4#" {
				fen = "rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq - 0 2"
			}
			g := NewGame()
			g.LoadFen(fen)
			g.GenerateLegalMoves()

			found := false
			for _, m := range g.LegalMoves {
				sFrom, sTo := m.ToFromToStrings()
				if sFrom == tt.moveFrom && sTo == tt.moveTo {
					if tt.promote != "" {
						if !m.IsPromotionMove() {
							continue
						}
						// Check promotion piece
						// Using move method if available or logic
						// We know how to check promotion from move bits if needed
						// But GetMoveSan should handle it.
						// We filter by promotion kind
						// Q=5, R=4, B=3, N=2
						p := m.GetPromotionTo()
						var pStr string
						switch p {
						case QUEEN:
							pStr = "Q"
						case ROOK:
							pStr = "R"
						case BISHOP:
							pStr = "B"
						case KNIGHT:
							pStr = "N"
						}
						if pStr != tt.promote {
							continue
						}
					} else {
						// Ensure we don't pick a promotion move if we expect a non-promotion (rare for pawns to last rank without, impossible actually)
						// But for safety
						if m.IsPromotionMove() {
							continue
						}
					}

					san := g.GetMoveSan(m)
					if san != tt.expectedSan {
						t.Errorf("Expected SAN %s, got %s", tt.expectedSan, san)
					}
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Move %s->%s not found in legal moves for FEN %s", tt.moveFrom, tt.moveTo, fen)
			}
		})
	}
}

func TestUndoMoveLogic(t *testing.T) {
	// Simple game sequence: 1. e4 e5 2. Nf3 Nc6
	fenStart := STARTING_POSITION_FEN
	g := NewGame()
	g.LoadFen(fenStart)
	g.GenerateLegalMoves()

	movesToMakeStr := []string{"e2 e4", "e7 e5", "g1 f3", "b8 c6"}
	var madeMoves []Move

	// Record hashes/state at each step
	hashes := []uint64{g.ZobristHash}
	fens := []string{g.ToFen()}

	for _, moveStr := range movesToMakeStr {
		parts := strings.Split(moveStr, " ")
		fromStr, toStr := parts[0], parts[1]

		var move Move
		found := false
		for _, m := range g.LegalMoves {
			sFrom, sTo := m.ToFromToStrings()
			if sFrom == fromStr && sTo == toStr {
				move = m
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Move %s not found", moveStr)
		}

		g.MakeMove(move)
		madeMoves = append(madeMoves, move)
		hashes = append(hashes, g.ZobristHash)
		fens = append(fens, g.ToFen())
	}

	// Now Undo backwards
	// steps: 0(start), 1(e4), 2(e5), 3(Nf3), 4(Nc6)
	// hashes has 0..4.
	// madeMoves has 0..3 (e4, e5, Nf3, Nc6).

	for i := len(madeMoves) - 1; i >= 0; i-- {
		move := madeMoves[i]
		g.UndoMove(move)

		expectedHash := hashes[i]
		expectedFen := fens[i]

		if g.ZobristHash != expectedHash {
			t.Errorf("Hash mismatch after undoing move %d (%s). Expected %x, got %x", i, movesToMakeStr[i], expectedHash, g.ZobristHash)
		}
		if g.ToFen() != expectedFen {
			t.Errorf("FEN mismatch after undoing move %d (%s). Expected %s, got %s", i, movesToMakeStr[i], expectedFen, g.ToFen())
		}
	}
}

func TestUndoMoveRecursive(t *testing.T) {
	// Better test: Do a random walk of depth N, then undo all, assert state equals start.
	g := NewGame()
	g.LoadFen(STARTING_POSITION_FEN)
	startHash := g.ZobristHash
	startFen := g.ToFen()

	depth := 4
	var perform func(d int)
	perform = func(d int) {
		if d == 0 {
			return
		}
		moves := g.LegalMoves
		if len(moves) == 0 {
			return
		}
		// Pick first move (deterministic) or iterate all?
		// Iterating all is too big for unit test usually.
		// Let's pick 2 distinct moves to branch.
		limit := 2
		if len(moves) < limit {
			limit = len(moves)
		}

		for i := 0; i < limit; i++ {
			m := moves[i]
			// clone backup ? No need, logic is reversible.
			// But to check if Undo works we rely on it working :P

			// Hash before
			hashBefore := g.ZobristHash
			fenBefore := g.ToFen()

			g.MakeMove(m)

			perform(d - 1)

			g.UndoMove(m)

			// Verify restoration
			if g.ZobristHash != hashBefore {
				t.Errorf("Hash mismatch after undo at depth %d. Move: %s. Expected %x, got %x", d, m.ToString(), hashBefore, g.ZobristHash)
			}
			if g.ToFen() != fenBefore {
				t.Errorf("FEN mismatch after undo at depth %d. Move: %s. Expected %s, got %s", d, m.ToString(), fenBefore, g.ToFen())
			}
		}
	}

	perform(depth)

	if g.ZobristHash != startHash {
		t.Error("Final hash mismatch after traversal")
	}
	if g.ToFen() != startFen {
		t.Error("Final FEN mismatch after traversal")
	}
}

// Helper to assert that a string slice contains a specific string
func assertContains(t *testing.T, list []string, item string) {
	t.Helper()
	found := false
	for _, s := range list {
		if s == item {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected list to contain %q, but it didn't. List: %v", item, list)
	}
}

// Helper to assert that a string slice does NOT contain a specific string
func assertNotContains(t *testing.T, list []string, item string) {
	t.Helper()
	found := false
	for _, s := range list {
		if s == item {
			found = true
			break
		}
	}
	if found {
		t.Errorf("Expected list to NOT contain %q, but it did.", item)
	}
}
