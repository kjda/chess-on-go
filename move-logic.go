package chessongo

import "strings"

const maxGeneratedMoves = 256

var genMovesCalls uint = 0

// Generate all peseudo moves
func (g *Game) GeneratePseudoMoves() {
	var ours [7]Bitboard
	var oursAll Bitboard
	if g.Turn == WHITE {
		ours = g.Whites
		oursAll = g.WhitePieces
	} else {
		ours = g.Blacks
		oursAll = g.BlackPieces
	}
	// Reuse underlying array capacity if available
	if cap(g.PseudoMoves) < maxGeneratedMoves {
		g.PseudoMoves = make([]Move, 0, maxGeneratedMoves)
	} else {
		g.PseudoMoves = g.PseudoMoves[:0]
	}
	g.genPawnOneStep()
	g.genPawnTwoSteps()
	g.genPawnAttacks()
	g.genFromMoves(ours[KING], oursAll, KING_ATTACKS_FROM[:])
	g.genFromMoves(ours[KNIGHT], oursAll, KNIGHT_ATTACKS_FROM[:])
	g.genRayMoves(ours[BISHOP]|ours[QUEEN], oursAll, BISHOP_DIRECTIONS[:])
	g.genRayMoves(ours[ROOK]|ours[QUEEN], oursAll, ROOK_DIRECTIONS[:])
	g.genCastling()
}

// Generate all legal moves
func (g *Game) GenerateLegalMoves() {
	// Ensure current check status is known before validating castling legality.
	g.IsCheck = g.ComputeIsCheck()
	g.GeneratePseudoMoves()
	needed := len(g.PseudoMoves)
	if cap(g.LegalMoves) < needed {
		g.LegalMoves = make([]Move, 0, needed)
	} else {
		g.LegalMoves = g.LegalMoves[:0]
	}
	for _, move := range g.PseudoMoves {
		if g.CanMove(move) {
			g.LegalMoves = append(g.LegalMoves, move)
		}
	}
}

// Generates King & Knight pseudo-legal moves
func (g *Game) genFromMoves(pieces, ours Bitboard, attackFrom []Bitboard) {
	for pieces > 0 {
		from := pieces.popLSB()
		targets := attackFrom[from] & ^ours
		for targets > 0 {
			to := targets.popLSB()
			g.PseudoMoves = append(g.PseudoMoves, NewMove(Square(from), Square(to), g.Squares[to]))
		}
	}

}

// Generate sliding-piece's pseudo-legal moves
func (g *Game) genRayMoves(pieces, ours Bitboard, directions []Direction) {
	for pieces > 0 {
		from := pieces.popLSB()
		var allTargets, targets Bitboard
		for _, direction := range directions {
			targets = RAY_MASKS[direction][from]
			blockers := targets & g.Occupied
			if blockers > 0 {
				if DIRECTION_LSB_MSP[direction] == LSB {
					targets ^= RAY_MASKS[direction][blockers.lsbIndex()]
				} else {
					targets ^= RAY_MASKS[direction][blockers.msbIndex()]
				}
			}
			allTargets |= targets & (^ours)
		}
		for allTargets > 0 {
			to := allTargets.popLSB()
			g.PseudoMoves = append(g.PseudoMoves, NewMove(Square(from), Square(to), g.Squares[to]))
		}
	}
}

// Generate castling pseudo-legal moves
func (g *Game) genCastling() {
	if g.Turn == WHITE && (g.Castling&CASTLE_WKS) > 0 && (g.Occupied&(0x3<<61)) == 0 {
		from := Square(g.Whites[KING].lsbIndex())
		to := Square(WKS_KING_TO_SQUARE)
		g.PseudoMoves = append(g.PseudoMoves, NewCastlingMove(from, to))

	}

	if g.Turn == WHITE && (g.Castling&CASTLE_WQS) > 0 && (g.Occupied&(0x7<<57)) == 0 {
		from := Square(g.Whites[KING].lsbIndex())
		to := Square(WQS_KING_TO_SQUARE)
		g.PseudoMoves = append(g.PseudoMoves, NewCastlingMove(from, to))
	}

	if g.Turn == BLACK && (g.Castling&CASTLE_BKS) > 0 && (g.Occupied&(0x3<<5)) == 0 {
		from := Square(g.Blacks[KING].lsbIndex())
		to := Square(BKS_KING_TO_SQUARE)
		g.PseudoMoves = append(g.PseudoMoves, NewCastlingMove(from, to))
	}

	if g.Turn == BLACK && (g.Castling&CASTLE_BQS) > 0 && (g.Occupied&(0x7<<1)) == 0 {
		from := Square(g.Blacks[KING].lsbIndex())
		to := Square(BQS_KING_TO_SQUARE)
		g.PseudoMoves = append(g.PseudoMoves, NewCastlingMove(from, to))
	}
	return
}

// Generate Pawn-one-step-forward pseudo-legal moves
func (g *Game) genPawnOneStep() {
	var targets Bitboard
	var shift int = 8
	if g.Turn == WHITE {
		targets = (g.Whites[PAWN] >> 8) & ^g.Occupied
	} else {
		targets = (g.Blacks[PAWN] << 8) & ^g.Occupied
		shift = -8
	}
	for targets > 0 {
		to := Square(targets.popLSB())
		from := Square(int(to) + shift)
		if g.IsToPromotionRank(to) {
			g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], QUEEN))
			g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], ROOK))
			g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], KNIGHT))
			g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], BISHOP))
		} else {
			g.PseudoMoves = append(g.PseudoMoves, NewMove(from, to, g.Squares[to]))
		}
	}
}

// Generate Pawn-two-step-forward pseudo-legal moves
func (g *Game) genPawnTwoSteps() {
	var targets Bitboard
	var shift int
	if g.Turn == WHITE {
		rank3filtered := ((g.Whites[PAWN] & Bitboard(RANK2_MASK)) >> 8) &^ g.Occupied
		targets = ((rank3filtered & Bitboard(RANK3_MASK)) >> 8) &^ g.Occupied
		shift = 16
	} else {
		rank6filtered := ((g.Blacks[PAWN] & Bitboard(RANK7_MASK)) << 8) &^ g.Occupied
		targets = ((rank6filtered & Bitboard(RANK6_MASK)) << 8) &^ g.Occupied
		shift = -16
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		g.PseudoMoves = append(g.PseudoMoves, NewMove(Square(from), Square(to), g.Squares[to]))
	}
}

// Generate pawns left and right attacks
func (g *Game) genPawnAttacks() {
	ours, _ := g.GetPawns()
	var targets Bitboard
	enPassant := Bitboard(0)
	if g.EnPassant > 0 {
		enPassant = Bitboard(0x1 << uint(g.EnPassant))
	}
	for _, shift := range [2]int{7, 9} {
		if g.Turn == WHITE {
			if shift == 7 {
				targets = (ours & ^Bitboard(FILE_H_MASK)) >> uint(shift)
			} else {
				targets = (ours & ^Bitboard(FILE_A_MASK)) >> uint(shift)
			}
			targets &= (g.BlackPieces | enPassant)
		} else {
			if shift == 7 {
				targets = (ours & ^Bitboard(FILE_A_MASK)) << uint(shift)
			} else {
				targets = (ours & ^Bitboard(FILE_H_MASK)) << uint(shift)
			}
			targets &= (g.WhitePieces | enPassant)
		}
		for targets > 0 {
			to := Square(targets.popLSB())
			fromShift := shift
			if g.Turn == BLACK {
				fromShift *= -1
			}
			from := Square(int(to) + fromShift)
			if g.EnPassant > 0 && to == g.EnPassant {
				var capturedPiece Piece
				if g.Turn == WHITE {
					capturedPiece = g.Squares[to+8]
				} else {
					capturedPiece = g.Squares[to-8]
				}
				g.PseudoMoves = append(g.PseudoMoves, NewEnPassantMove(from, to, capturedPiece))
			} else if g.IsToPromotionRank(to) {
				g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], QUEEN))
				g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], ROOK))
				g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], KNIGHT))
				g.PseudoMoves = append(g.PseudoMoves, NewPromotionMove(from, to, g.Squares[to], BISHOP))
			} else {
				g.PseudoMoves = append(g.PseudoMoves, NewMove(from, to, g.Squares[to]))
			}
		}
	}
}

func (g *Game) IsToPromotionRank(to Square) bool {
	return (g.Turn == WHITE && (Bitboard(0x1<<uint(to))&Bitboard(RANK8_MASK) > 0)) || (g.Turn == BLACK && (Bitboard(0x1<<uint(to))&Bitboard(RANK1_MASK) > 0))
}

// Checks whether our king is in check or not
func (g *Game) ComputeIsCheck() bool {
	genMovesCalls++
	var kingBB, theirsAll, attackers Bitboard
	var theirs []Bitboard
	if g.Turn == WHITE {
		kingBB, theirs, theirsAll = g.Whites[KING], g.Blacks[:], g.BlackPieces
	} else {
		kingBB, theirs, theirsAll = g.Blacks[KING], g.Whites[:], g.WhitePieces
	}
	kingIdx := kingBB.lsbIndex()
	possibleAttackers := theirsAll & ATTACKS_TO[kingIdx]

	attackers = (theirs[ROOK] | theirs[QUEEN]) & possibleAttackers
	if attackers > 0 && g.isCheckedFromRay(kingBB, attackers, ROOK_DIRECTIONS[:]) {
		return true
	}

	attackers = (theirs[BISHOP] | theirs[QUEEN]) & possibleAttackers
	if attackers > 0 && g.isCheckedFromRay(kingBB, attackers, BISHOP_DIRECTIONS[:]) {
		return true
	}

	attackers = theirs[KNIGHT] & possibleAttackers
	for attackers > 0 {
		from := attackers.popLSB()
		if KNIGHT_ATTACKS_FROM[from]&kingBB > 0 {
			return true
		}
	}

	if g.Turn == WHITE {
		// Black pawns attack “down” the board (towards higher square indices).
		if ((g.Blacks[PAWN]&^Bitboard(FILE_A_MASK))<<7)&kingBB > 0 ||
			((g.Blacks[PAWN]&^Bitboard(FILE_H_MASK))<<9)&kingBB > 0 {
			return true
		}
	} else {
		// White pawns attack “up” the board (towards lower square indices).
		if ((g.Whites[PAWN]&^Bitboard(FILE_H_MASK))>>7)&kingBB > 0 ||
			((g.Whites[PAWN]&^Bitboard(FILE_A_MASK))>>9)&kingBB > 0 {
			return true
		}
	}

	attackers = theirs[KING] & possibleAttackers
	if attackers > 0 {
		from := attackers.popLSB()
		if KING_ATTACKS_FROM[from]&kingBB > 0 {
			return true
		}
	}
	return false
}

// checks whether target is attacked by one of the "attackers"
func (g *Game) isCheckedFromRay(target, attackers Bitboard, directions []Direction) bool {
	var targets Bitboard
	var from uint
	for attackers > 0 {
		from = attackers.popLSB()
		for _, direction := range directions {
			targets = RAY_MASKS[direction][from]
			blockers := targets & g.Occupied
			if blockers > 0 {
				if DIRECTION_LSB_MSP[direction] == LSB {
					targets ^= RAY_MASKS[direction][blockers.lsbIndex()]
				} else {
					targets ^= RAY_MASKS[direction][blockers.msbIndex()]
				}
			}
			if targets&target > 0 {
				return true
			}
		}
	}
	return false
}

// Checks whether the given move is possible or not
func (g *Game) CanMove(m Move) bool {
	if m.IsCastlingMove() {
		var inBetweenSq Square
		if m.To() == WKS_KING_TO_SQUARE || m.To() == BKS_KING_TO_SQUARE {
			inBetweenSq = m.From() + 1
		} else if m.To() == WQS_KING_TO_SQUARE || m.To() == BQS_KING_TO_SQUARE {
			inBetweenSq = m.From() - 1
		}
		inBetweenMove := NewMove(m.From(), inBetweenSq, EMPTY)
		if g.IsCheck || g.WillMoveCauseCheck(inBetweenMove) {
			return false
		}
	}
	return !g.WillMoveCauseCheck(m)
}

func (g *Game) WillMoveCauseCheck(m Move) bool {
	// Optimization: Stack-copy the board. accessing underlying arrays by value.
	// Since justMove/ComputeIsCheck don't modify maps/slices (only arrays/primitives), this is safe and allocation-free.
	clone := *g
	clone.justMove(m)
	return clone.ComputeIsCheck()
}

func (g *Game) MakeMove(m Move) {
	// Capture state for UndoMove
	capturedPiece := g.Squares[m.To()]
	if m.IsEnPassant() {
		if g.Turn == WHITE {
			capturedPiece = g.Squares[m.To()+8]
		} else {
			capturedPiece = g.Squares[m.To()-8]
		}
	}
	g.History = append(g.History, GameState{
		CapturedPiece: capturedPiece,
		Castling:      g.Castling,
		EnPassant:     g.EnPassant,
		HalfMoves:     g.HalfMoves,
		ZobristHash:   g.ZobristHash,
	})

	if g.ShouldResetPly(m) {
		g.HalfMoves = 0
	} else {
		g.HalfMoves++
	}
	if g.ShouldIncFullMoves(m) {
		g.FullMoves++
	}

	g.justMove(m)
	kind := g.Squares[m.To()].Kind()
	if kind == KING {
		if g.Turn == WHITE {
			g.Castling &= ^(CASTLE_WKS | CASTLE_WQS)
		} else {
			g.Castling &= ^(CASTLE_BKS | CASTLE_BQS)
		}
	}
	if kind == ROOK {
		switch m.From() {
		case WKS_ROOK_ORIGINAL_SQUARE:
			g.Castling &= ^CASTLE_WKS
		case WQS_ROOK_ORIGINAL_SQUARE:
			g.Castling &= ^CASTLE_WQS
		case BKS_ROOK_ORIGINAL_SQUARE:
			g.Castling &= ^CASTLE_BKS
		case BQS_ROOK_ORIGINAL_SQUARE:
			g.Castling &= ^CASTLE_BQS
		}
	}

	switch m.To() {
	case WKS_ROOK_ORIGINAL_SQUARE:
		g.Castling &= ^CASTLE_WKS
	case WQS_ROOK_ORIGINAL_SQUARE:
		g.Castling &= ^CASTLE_WQS
	case BKS_ROOK_ORIGINAL_SQUARE:
		g.Castling &= ^CASTLE_BKS
	case BQS_ROOK_ORIGINAL_SQUARE:
		g.Castling &= ^CASTLE_BQS
	}
	// enPassant target
	g.EnPassant = 0
	if kind == PAWN && g.Turn == WHITE {
		if m.From().Rank() == 6 && m.To().Rank() == 4 {
			g.EnPassant = m.From() - 8
		}
	}
	if kind == PAWN && g.Turn == BLACK {
		if m.From().Rank() == 1 && m.To().Rank() == 3 {
			g.EnPassant = m.From() + 8
		}
	}

	if g.Turn == WHITE {
		g.Turn = BLACK
	} else {
		g.Turn = WHITE
	}

	g.recordPosition()

	g.GenerateLegalMoves()

	g.IsCheck = g.ComputeIsCheck()
	g.IsCheckmate = g.IsCheck && !g.hasMoves()
	g.IsStalement = !g.IsCheckmate && !g.hasMoves()
	g.IsMaterialDraw = g.hasInsufficientMaterial()
	g.IsThreefoldRepetition = g.checkThreefoldRepetition()
	g.IsFiftyMoveRule = g.checkFiftyMoveRule()
	g.IsSeventyFiveMoveRule = g.checkSeventyFiveMoveRule()
	g.IsFinished = (g.IsCheckmate || g.IsStalement || g.IsMaterialDraw || g.IsFivefoldRepetition() || g.IsSeventyFiveMoveRule)
}

func (g *Game) justMove(m Move) {
	from := m.From()
	to := m.To()

	//capturedPiece := m.captured()
	capturedPiece := g.Squares[to]
	if m.IsEnPassant() && g.Turn == WHITE {
		capturedPiece = g.Squares[to+8]
	} else if m.IsEnPassant() && g.Turn == BLACK {
		capturedPiece = g.Squares[to-8]
	}
	fromBBNeg := ^Bitboard(0x1 << from)
	toBB := Bitboard(0x1 << to)
	movingPiece := g.Squares[from]
	movingPieceKind := movingPiece.Kind()
	switch movingPiece.Color() {
	case WHITE:
		// update bitmap of moving piece kind, unset bit of source square
		g.Whites[movingPieceKind] &= fromBBNeg
		// update bitmap of moving piece kind, set bit of source square
		g.Whites[movingPieceKind] |= toBB
		// update white pieces bitboard - unset old square
		g.WhitePieces &= fromBBNeg
		// update white pieces bitboard - set new square
		g.WhitePieces |= toBB
	case BLACK:
		g.Blacks[movingPieceKind] &= fromBBNeg
		g.Blacks[movingPieceKind] |= toBB
		g.BlackPieces &= fromBBNeg
		g.BlackPieces |= toBB
	}

	g.Occupied &= fromBBNeg
	g.Occupied |= toBB

	g.Squares[m.To()] = g.Squares[m.From()]
	g.Squares[m.From()] = EMPTY
	if capturedPiece != EMPTY {
		if !m.IsEnPassant() {
			g.capturePiece(to, capturedPiece)
		} else {
			if g.Turn == WHITE {
				capSq := to + 8
				g.capturePiece(capSq, g.Squares[capSq])
				g.Occupied &= ^Bitboard(0x1 << capSq)
				g.Squares[to+8] = EMPTY
			} else {
				capSq := to - 8
				g.capturePiece(capSq, g.Squares[capSq])
				g.Occupied &= ^Bitboard(0x1 << capSq)
				g.Squares[to-8] = EMPTY
			}
		}
	}
	if m.IsCastlingMove() {
		var rookMove Move
		if m.To() == WKS_KING_TO_SQUARE || m.To() == BKS_KING_TO_SQUARE {
			rookMove = NewMove(m.To()+1, m.To()-1, 0)
		} else if m.To() == WQS_KING_TO_SQUARE || m.To() == BQS_KING_TO_SQUARE {
			rookMove = NewMove(m.To()-2, m.To()+1, 0)
		}
		g.justMove(rookMove)
	}
	var promoteTo Piece = m.GetPromotionTo()
	if promoteTo > 0 {
		switch g.Squares[to].Color() {
		case WHITE:
			// remove advanced pawn from boards
			g.Whites[PAWN] &= ^toBB
			// add promotePiece to board
			g.Whites[promoteTo] |= toBB
			g.WhitePieces |= toBB
		case BLACK:
			// remove advanced pawn from boards
			g.Blacks[PAWN] &= ^toBB
			// add promotePiece to board
			g.Blacks[promoteTo] |= toBB
			g.BlackPieces |= toBB
		}
		g.Squares[m.To()] = Piece(uint(promoteTo) | uint(g.Turn))
	}
}

// Remove captured piece from opponent's pieces
func (g *Game) capturePiece(sq Square, captured Piece) {
	if captured == EMPTY {
		return
	}
	sqBB := Bitboard(0x1 << sq)
	kind := captured.Kind()
	switch captured.Color() {
	case WHITE:
		g.Whites[kind] &= ^sqBB
		g.WhitePieces &= ^sqBB
	case BLACK:
		g.Blacks[kind] &= ^sqBB
		g.BlackPieces &= ^sqBB
	}
}

func (g *Game) hasInsufficientMaterial() bool {
	if g.Whites[QUEEN] > 0 || g.Whites[ROOK] > 0 || g.Whites[PAWN] > 0 {
		return false
	}
	if g.Blacks[QUEEN] > 0 || g.Blacks[ROOK] > 0 || g.Blacks[PAWN] > 0 {
		return false
	}
	if g.Whites[KNIGHT] > 0 && g.Whites[BISHOP] > 0 {
		return false
	}
	if g.Blacks[KNIGHT] > 0 && g.Blacks[BISHOP] > 0 {
		return false
	}

	if g.Whites[BISHOP].NumberOfSetBits() > 1 {
		return false
	}

	if g.Blacks[BISHOP].NumberOfSetBits() > 1 {
		return false
	}

	if g.Whites[KNIGHT].NumberOfSetBits() > 1 {
		return false
	}

	if g.Blacks[KNIGHT].NumberOfSetBits() > 1 {
		return false
	}

	return true
}

/*
 1. moving piece letter(exclude pawn)
    1.1 originating file letter of the moving piece
    1.2 OR: the originating rank digit of the moving piece
    1.3 OR: originating square
 2. if capturing pawn -> include originating file
 3. x for caputures
 4. destination square
 5. PawnPromotion -> "="  followed by promoted piece rune in uppercase
*/
func (g *Game) GetMoveSan(m Move) string {
	pgn := g.GetMoveSanWithoutSuffix(m)

	g.MakeMove(m)
	if g.IsCheckmate {
		pgn += "#"
	} else if g.IsCheck {
		pgn += "+"
	}
	g.UndoMove(m)

	return pgn
}

// GetMoveSanWithoutSuffix returns the SAN notation for a move without checking for check (+) or checkmate (#).
// This prevents expensive board cloning and is sufficient for PGN parsing matching.
func (g *Game) GetMoveSanWithoutSuffix(m Move) string {
	var sb strings.Builder
	from := m.From()
	to := m.To()

	if m.IsCastlingMove() {
		if m.To() == WKS_KING_TO_SQUARE || m.To() == BKS_KING_TO_SQUARE {
			sb.WriteString("O-O")
		}
		if m.To() == WQS_KING_TO_SQUARE || m.To() == BQS_KING_TO_SQUARE {
			sb.WriteString("O-O-O")
		}
	} else {
		movingKind := g.Squares[from].Kind()
		if movingKind != PAWN { // -------> 1.
			sb.WriteString(strings.ToUpper(string(g.Squares[from].ToRune())))
		}

		// Disambiguation:
		othersOfSameKind, onSameFileCount, onSameRankCount := g.GetOthersOfSameKindMovingToSameTargetCounts(m)
		if othersOfSameKind > 0 {
			if onSameFileCount == 0 {
				sb.WriteString(m.From().FileLetter()) // -------> 1.1
			} else if onSameRankCount == 0 {
				sb.WriteString(m.From().FileLetter()) // -------> 1.2
			} else {
				sb.WriteString(m.From().Coords()) // -------> 1.2
			}
		}

		isCapturing := m.GetCapturedPiece() != EMPTY
		if isCapturing && movingKind == PAWN {
			sb.WriteString(from.FileLetter()) // -------> 2.
		}
		if isCapturing {
			sb.WriteString("x") // -------> 3.
		}

		sb.WriteString(to.Coords()) // -------> 4.

		if m.IsPromotionMove() {
			sb.WriteString("=")
			sb.WriteString(strings.ToUpper(string(m.GetPromotionTo()))) // -------> 5.
		}
	}
	return sb.String()
}

func (g *Game) GetOthersOfSameKindMovingToSameTargetCounts(themove Move) (otherOfSameKind int, onSameFileCount int, onSameRankCount int) {
	movingPiece := g.Squares[themove.From()]
	to := themove.To()
	for _, m := range g.LegalMoves {
		if m == themove || m.To() != to || g.Squares[m.From()].Kind() != movingPiece.Kind() {
			continue
		}
		otherOfSameKind += 1
		if m.From().File() == to.File() {
			onSameFileCount += 1
		}
		if m.From().Rank() == to.Rank() {
			onSameRankCount += 1
		}
	}
	return
}

func (g *Game) UndoMove(m Move) {
	if len(g.History) == 0 {
		return
	}
	// Decrement history count for current position
	if g.PositionHistory != nil {
		g.PositionHistory[g.ZobristHash]--
		if g.PositionHistory[g.ZobristHash] <= 0 {
			delete(g.PositionHistory, g.ZobristHash)
		}
	}

	// Pop state
	state := g.History[len(g.History)-1]
	g.History = g.History[:len(g.History)-1]

	// Restore simple fields
	g.Castling = state.Castling
	g.EnPassant = state.EnPassant
	g.HalfMoves = state.HalfMoves
	g.ZobristHash = state.ZobristHash

	// Flip Turn back
	if g.Turn == WHITE {
		g.Turn = BLACK
		g.FullMoves--
	} else {
		g.Turn = WHITE
	}

	g.unmakeMove(m, state.CapturedPiece)

	// Re-calculate derived state
	g.GenerateLegalMoves()
	g.IsCheck = g.ComputeIsCheck()
	g.IsCheckmate = g.IsCheck && !g.hasMoves()
	g.IsStalement = !g.IsCheckmate && !g.hasMoves()
	g.IsMaterialDraw = g.hasInsufficientMaterial()
	g.IsThreefoldRepetition = g.checkThreefoldRepetition()
	g.IsFiftyMoveRule = g.checkFiftyMoveRule()
	g.IsSeventyFiveMoveRule = g.checkSeventyFiveMoveRule()
	g.IsFinished = (g.IsCheckmate || g.IsStalement || g.IsMaterialDraw || g.IsFivefoldRepetition() || g.IsSeventyFiveMoveRule)
}

func (g *Game) unmakeMove(m Move, captured Piece) {
	from := m.From()
	to := m.To()

	movingPieceKind := g.Squares[to].Kind()
	movingColor := g.Turn

	if m.IsPromotionMove() {
		// The piece at `to` is the promoted piece.
		promotedKind := m.GetPromotionTo()

		toBB := Bitboard(0x1 << to)
		fromBB := Bitboard(0x1 << from)

		if movingColor == WHITE {
			g.Whites[promotedKind] &= ^toBB
			g.WhitePieces &= ^toBB
			// Restore Pawn at `from`
			g.Whites[PAWN] |= fromBB
			g.WhitePieces |= fromBB
		} else {
			g.Blacks[promotedKind] &= ^toBB
			g.BlackPieces &= ^toBB
			g.Blacks[PAWN] |= fromBB
			g.BlackPieces |= fromBB
		}
		g.Squares[to] = EMPTY
		g.Squares[from] = Piece(uint(PAWN) | uint(movingColor))
		g.Occupied &= ^toBB
		g.Occupied |= fromBB

	} else {
		toBB := Bitboard(0x1 << to)
		fromBB := Bitboard(0x1 << from)

		if movingColor == WHITE {
			g.Whites[movingPieceKind] &= ^toBB
			g.Whites[movingPieceKind] |= fromBB
			g.WhitePieces &= ^toBB
			g.WhitePieces |= fromBB
		} else {
			g.Blacks[movingPieceKind] &= ^toBB
			g.Blacks[movingPieceKind] |= fromBB
			g.BlackPieces &= ^toBB
			g.BlackPieces |= fromBB
		}

		g.Occupied &= ^toBB
		g.Occupied |= fromBB

		g.Squares[from] = g.Squares[to]
		g.Squares[to] = EMPTY
	}

	if captured != EMPTY {
		if m.IsEnPassant() {
			var capSq Square
			if movingColor == WHITE {
				capSq = to + 8
			} else {
				capSq = to - 8
			}
			g.addPiece(captured, int(capSq))
		} else {
			g.addPiece(captured, int(to))
		}
	}

	if m.IsCastlingMove() {
		var rookFrom, rookTo Square

		if m.To() == WKS_KING_TO_SQUARE { // g1
			rookFrom = 61 // f1
			rookTo = 63   // h1
		} else if m.To() == WQS_KING_TO_SQUARE { // c1
			rookFrom = 59 // d1
			rookTo = 56   // a1
		} else if m.To() == BKS_KING_TO_SQUARE { // g8
			rookFrom = 5 // f8
			rookTo = 7   // h8
		} else if m.To() == BQS_KING_TO_SQUARE { // c8
			rookFrom = 3 // d8
			rookTo = 0   // a8
		}

		// Move rook back
		rFromBB := Bitboard(0x1 << rookFrom)
		rToBB := Bitboard(0x1 << rookTo)

		if movingColor == WHITE {
			g.Whites[ROOK] &= ^rFromBB
			g.Whites[ROOK] |= rToBB
			g.WhitePieces &= ^rFromBB
			g.WhitePieces |= rToBB
		} else {
			g.Blacks[ROOK] &= ^rFromBB
			g.Blacks[ROOK] |= rToBB
			g.BlackPieces &= ^rFromBB
			g.BlackPieces |= rToBB
		}
		g.Occupied &= ^rFromBB
		g.Occupied |= rToBB

		g.Squares[rookTo] = g.Squares[rookFrom]
		g.Squares[rookFrom] = EMPTY
	}
}
