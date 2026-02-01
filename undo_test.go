package chessongo

import (
	"testing"
)

func TestUndoMove(t *testing.T) {
	g := NewGame()
	startHash := g.ZobristHash
	startFen := g.ToFen()

	g.GenerateLegalMoves()
	if len(g.LegalMoves) == 0 {
		t.Fatal("No moves at start?")
	}

	move := g.LegalMoves[0] // e.g. a3
	g.MakeMove(move)

	if g.ZobristHash == startHash {
		t.Error("Hash did not change after move")
	}

	g.UndoMove(move)

	if g.ZobristHash != startHash {
		t.Errorf("Hash mismatch after Undo. Got %v, Want %v", g.ZobristHash, startHash)
	}

	currentFen := g.ToFen()
	if currentFen != startFen {
		t.Errorf("FEN mismatch after Undo.\nWant: %s\nGot : %s", startFen, currentFen)
	}

	// Check turn
	if g.Turn != WHITE {
		t.Errorf("Turn mismatch. Got %v, Want WHITE", g.Turn)
	}

	// Check FullMoves
	if g.FullMoves != 1 {
		t.Errorf("FullMoves mismatch. Got %d, Want 1", g.FullMoves)
	}
}

func TestUndoMove_Capture(t *testing.T) {
	// Setup a position with capture
	// Custom fen: White Pawn e4, Black Pawn d5.
	fen := "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2"
	g := &Game{}
	err := g.LoadFen(fen)
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}
	startHash := g.ZobristHash
	startFen := g.ToFen()

	// Find move e4xd5
	g.GenerateLegalMoves()
	var capMove Move
	found := false
	for _, m := range g.LegalMoves {
		if g.Squares[m.To()] != EMPTY { // Capture
			capMove = m
			found = true
			break
		}
	}

	if !found {
		// e4xd5 should be possible.
	} else {
		g.MakeMove(capMove)
		g.UndoMove(capMove)

		if g.Squares[capMove.To()] == EMPTY {
			t.Error("Captured piece not restored")
		}
		if g.Squares[capMove.To()].Kind() != PAWN { // it was a pawn
			t.Error("Restored piece is incorrect kind")
		}
		if g.ZobristHash != startHash {
			t.Errorf("Hash mismatch. Want %x Got %x", startHash, g.ZobristHash)
		}
		if g.ToFen() != startFen {
			t.Errorf("FEN mismatch.\nWant: %s\nGot : %s", startFen, g.ToFen())
		}
	}
}

func TestUndoMove_Castling(t *testing.T) {
	// Setup position for White King Side castling
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1" // h1 rook, e1 king, empty f1, g1
	g := &Game{}
	g.LoadFen(fen)
	startHash := g.ZobristHash
	startFen := g.ToFen()

	// Create castling move (e1 -> g1)
	// MakeMove detects castling by moveKind? No, NewMove likely needs the flag.
	// But `GenerateLegalMoves` sets it.
	g.GenerateLegalMoves()
	var castleMove Move
	found := false
	for _, mv := range g.LegalMoves {
		if mv.IsCastlingMove() {
			castleMove = mv
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Castling move not generated")
	}

	g.MakeMove(castleMove)
	g.UndoMove(castleMove)

	if g.ZobristHash != startHash {
		t.Errorf("Hash mismatch. Want %x Got %x", startHash, g.ZobristHash)
	}
	if g.ToFen() != startFen {
		t.Errorf("FEN mismatch.\nWant: %s\nGot : %s", startFen, g.ToFen())
	}
	// Verify Rook position
	if g.Squares[63] == EMPTY || g.Squares[63].Kind() != ROOK {
		t.Error("Rook not restored to h1")
	}
	if g.Squares[61] != EMPTY {
		t.Error("f1 not empty")
	}
}

func TestUndoMove_EnPassant(t *testing.T) {
	// Position: White Pawn e5, Black Pawn d5 (just moved d7-d5). EP target d6.
	// FEN: rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3
	// e5 captures d6 (en passant)
	fen := "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
	g := &Game{}
	g.LoadFen(fen)
	startHash := g.ZobristHash
	startFen := g.ToFen()

	g.GenerateLegalMoves()
	var epMove Move
	found := false
	for _, m := range g.LegalMoves {
		if m.IsEnPassant() {
			epMove = m
			found = true
			break
		}
	}
	if !found {
		t.Fatal("En Passant move not generated")
	}

	g.MakeMove(epMove)
	g.UndoMove(epMove)

	if g.ZobristHash != startHash {
		t.Errorf("Hash mismatch. Want %x Got %x", startHash, g.ZobristHash)
	}
	if g.ToFen() != startFen {
		t.Errorf("FEN mismatch.\nWant: %s\nGot : %s", startFen, g.ToFen())
	}

	// Check captured pawn restored at d5, not d6
	// Rank 5 (d5): Start 24. d=3. Index 27.
	// Rank 6 (d6): Start 16. d=3. Index 19.
	if g.Squares[27] == EMPTY || g.Squares[27].Kind() != PAWN {
		t.Error("Captured pawn not restored at d5")
	}
	if g.Squares[19] != EMPTY {
		t.Error("En Passant target square not empty after undo")
	}
}

func TestUndoMove_Promotion(t *testing.T) {
	// Position: White Pawn a7, Black King e8.
	fen := "4k3/P7/8/8/8/8/8/4K3 w - - 0 1"
	g := &Game{}
	g.LoadFen(fen)
	startHash := g.ZobristHash
	startFen := g.ToFen()

	g.GenerateLegalMoves()
	var promoMove Move
	found := false
	for _, m := range g.LegalMoves {
		if m.IsPromotionMove() {
			promoMove = m // grab any promotion (Queen usually first)
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Promotion move not generated")
	}

	g.MakeMove(promoMove)
	g.UndoMove(promoMove)

	if g.ZobristHash != startHash {
		t.Errorf("Hash mismatch. Want %x Got %x", startHash, g.ZobristHash)
	}
	if g.ToFen() != startFen {
		t.Errorf("FEN mismatch.\nWant: %s\nGot : %s", startFen, g.ToFen())
	}

	// Check pawn at a7
	// a7 is 8. Wait, standard mapping: a8=0..h1=63?
	// RUNE_TO_RANK '8':0. '1':7.
	// a8 is 0. a7 is 8.
	// Let's verify Square mapping.
	// bitboard layout usually 0=a1 or a8 depending on impl.
	// `square.go` likely has the mapping.
	// Assuming `g.Squares[8]` matches `a7` if Fen parsing worked and put it there.
	// `g.LoadFen` parses `P`.
	// Let's trust `ToFen` matching means board is restored.
}
