package chessongo

import "strings"

var genMovesCalls uint = 0

//Generate all peseudo moves
func (b *Board) GeneratePseudoMoves() {
	var ours [7]Bitboard
	var oursAll Bitboard
	if b.Turn == WHITE {
		ours = b.Whites
		oursAll = b.WhitePieces
	} else {
		ours = b.Blacks
		oursAll = b.BlackPieces
	}
	b.PseudoMoves = []Move{}
	b.genPawnOneStep()
	b.genPawnTwoSteps()
	b.genPawnAttacks()
	b.genFromMoves(ours[KING], oursAll, KING_ATTACKS_FROM[:])
	b.genFromMoves(ours[KNIGHT], oursAll, KNIGHT_ATTACKS_FROM[:])
	b.genRayMoves(ours[BISHOP]|ours[QUEEN], oursAll, BISHOP_DIRECTIONS[:])
	b.genRayMoves(ours[ROOK]|ours[QUEEN], oursAll, ROOK_DIRECTIONS[:])
	b.genCastling()
}

//Generate all legal moves
func (b *Board) GenerateLegalMoves() {
	b.GeneratePseudoMoves()
	b.LegalMoves = []Move{}
	for _, move := range b.PseudoMoves {
		if b.CanMove(move) {
			b.LegalMoves = append(b.LegalMoves, move)
		}
	}
}

//Generates King & Knight pseudo-legal moves
func (b *Board) genFromMoves(pieces, ours Bitboard, attackFrom []Bitboard) {
	for pieces > 0 {
		from := pieces.popLSB()
		targets := attackFrom[from] & ^ours
		for targets > 0 {
			to := targets.popLSB()
			b.PseudoMoves = append(b.PseudoMoves, NewMove(Square(from), Square(to), b.Squares[to]))
		}
	}

}

//Generate sliding-piece's pseudo-legal moves
func (b *Board) genRayMoves(pieces, ours Bitboard, directions []Direction) {
	for pieces > 0 {
		from := pieces.popLSB()
		var allTargets, targets Bitboard
		for _, direction := range directions {
			targets = RAY_MASKS[direction][from]
			blockers := targets & b.Occupied
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
			b.PseudoMoves = append(b.PseudoMoves, NewMove(Square(from), Square(to), b.Squares[to]))
		}
	}
}

//Generate castling pseudo-legal moves
func (b *Board) genCastling() {
	if b.Turn == WHITE && (b.Castling&CASTLE_WKS) > 0 && (b.Occupied&(0x3<<61)) == 0 {
		from := Square(b.Whites[KING].lsbIndex())
		to := Square(WKS_KING_TO_SQUARE)
		b.PseudoMoves = append(b.PseudoMoves, NewCastlingMove(from, to))

	}

	if b.Turn == WHITE && (b.Castling&CASTLE_WQS) > 0 && (b.Occupied&(0x7<<57)) == 0 {
		from := Square(b.Whites[KING].lsbIndex())
		to := Square(WQS_KING_TO_SQUARE)
		b.PseudoMoves = append(b.PseudoMoves, NewCastlingMove(from, to))
	}

	if b.Turn == BLACK && (b.Castling&CASTLE_BKS) > 0 && (b.Occupied&(0x3<<5)) == 0 {
		from := Square(b.Blacks[KING].lsbIndex())
		to := Square(BKS_KING_TO_SQUARE)
		b.PseudoMoves = append(b.PseudoMoves, NewCastlingMove(from, to))
	}

	if b.Turn == BLACK && (b.Castling&CASTLE_BQS) > 0 && (b.Occupied&(0x7<<1)) == 0 {
		from := Square(b.Blacks[KING].lsbIndex())
		to := Square(BQS_KING_TO_SQUARE)
		b.PseudoMoves = append(b.PseudoMoves, NewCastlingMove(from, to))
	}
	return
}

//Generate Pawn-one-step-forward pseudo-legal moves
func (b *Board) genPawnOneStep() {
	var targets Bitboard
	var shift int = 8
	if b.Turn == WHITE {
		targets = (b.Whites[PAWN] >> 8) & ^b.Occupied
	} else {
		targets = (b.Blacks[PAWN] << 8) & ^b.Occupied
		shift = -8
	}
	for targets > 0 {
		to := Square(targets.popLSB())
		from := Square(int(to) + shift)
		if b.IsToPromotionRank(to) {
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], QUEEN))
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], ROOK))
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], KNIGHT))
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], BISHOP))
		} else {
			b.PseudoMoves = append(b.PseudoMoves, NewMove(from, to, b.Squares[to]))
		}
	}
}

//Generate Pawn-two-step-forward pseudo-legal moves
func (b *Board) genPawnTwoSteps() {
	var targets Bitboard
	var shift int
	if b.Turn == WHITE {
		rank3filtered := ((b.Whites[PAWN] & RANK2_MASK) >> 8) &^ b.Occupied
		targets = ((rank3filtered & RANK3_MASK) >> 8) &^ b.Occupied
		shift = 16
	} else {
		rank6filtered := ((b.Blacks[PAWN] & RANK7_MASK) << 8) &^ b.Occupied
		targets = ((rank6filtered & RANK6_MASK) << 8) &^ b.Occupied
		shift = -16
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		b.PseudoMoves = append(b.PseudoMoves, NewMove(Square(from), Square(to), b.Squares[to]))
	}
}

//Generate pawns left and right attacks
func (b *Board) genPawnAttacks() {
	ours, _ := b.GetPawns()
	var targets Bitboard
	enPassant := Bitboard(0)
	if b.EnPassant > 0 {
		enPassant = Bitboard(0x1 << uint(b.EnPassant))
	}
	for _, shift := range [2]int{7, 9} {
		fromShift := shift
		if b.Turn == WHITE {
			targets = Bitboard(ours>>uint(shift)) & (b.BlackPieces | enPassant)
		} else {
			targets = Bitboard(ours<<uint(shift)) & (b.WhitePieces | enPassant)
			fromShift *= -1
		}
		for targets > 0 {
			to := Square(targets.popLSB())
			from := Square(int(to) + fromShift)
			if b.EnPassant > 0 && to == b.EnPassant {
				var capturedSq Square
				if b.Turn == WHITE {
					capturedSq = to + 8
				} else {
					capturedSq = to - 8
				}
				b.PseudoMoves = append(b.PseudoMoves, NewEnPassantMove(from, to, b.Squares[capturedSq]))
			} else if b.IsToPromotionRank(to) {
				b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], QUEEN))
				b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], ROOK))
				b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], KNIGHT))
				b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(from, to, b.Squares[to], BISHOP))
			} else {
				b.PseudoMoves = append(b.PseudoMoves, NewMove(from, to, b.Squares[to]))
			}
		}
	}
}

func (b *Board) IsToPromotionRank(to Square) bool {
	return (b.Turn == WHITE && (Bitboard(0x1<<to)&RANK8_MASK > 0)) || (b.Turn == BLACK && (Bitboard(0x1<<to)&RANK1_MASK > 0))
}

//Checks whether our king is in check or not
func (b *Board) ComputeIsCheck() bool {
	genMovesCalls++
	var kingBB, theirsAll, attackers, targets Bitboard
	var theirs []Bitboard
	if b.Turn == WHITE {
		kingBB, theirs, theirsAll = b.Whites[KING], b.Blacks[:], b.BlackPieces
	} else {
		kingBB, theirs, theirsAll = b.Blacks[KING], b.Whites[:], b.WhitePieces
	}
	kingIdx := kingBB.lsbIndex()
	possibleAttackers := theirsAll & ATTACKS_TO[kingIdx]

	attackers = (theirs[ROOK] | theirs[QUEEN]) & possibleAttackers
	if attackers > 0 && b.isCheckedFromRay(kingBB, attackers, ROOK_DIRECTIONS[:]) {
		return true
	}

	attackers = (theirs[BISHOP] | theirs[QUEEN]) & possibleAttackers
	if attackers > 0 && b.isCheckedFromRay(kingBB, attackers, BISHOP_DIRECTIONS[:]) {
		return true
	}

	attackers = theirs[KNIGHT] & possibleAttackers
	for attackers > 0 {
		from := attackers.popLSB()
		if KNIGHT_ATTACKS_FROM[from]&kingBB > 0 {
			return true
		}
	}

	enPassant := Bitboard(0x1 << uint(b.EnPassant))
	if b.Turn == WHITE {
		targets = Bitboard(b.Blacks[PAWN]<<uint(7)) & (b.WhitePieces | enPassant)
		targets |= Bitboard(b.Blacks[PAWN]<<uint(9)) & (b.WhitePieces | enPassant)
	} else {
		targets = Bitboard(b.Whites[PAWN]>>uint(7)) & (b.BlackPieces | enPassant)
		targets |= Bitboard(b.Whites[PAWN]>>uint(9)) & (b.BlackPieces | enPassant)
	}
	if targets&kingBB > 0 {
		return true
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

//checks whether target is attacked by one of the "attackers"
func (b *Board) isCheckedFromRay(target, attackers Bitboard, directions []Direction) bool {
	var targets Bitboard
	var from uint
	for attackers > 0 {
		from = attackers.popLSB()
		for _, direction := range directions {
			targets = RAY_MASKS[direction][from]
			blockers := targets & b.Occupied
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

//Checks whether the given move is possible or not
func (b *Board) CanMove(m Move) bool {
	if m.IsCastlingMove() {
		var inBetweenSq Square
		if m.To() == WKS_KING_TO_SQUARE || m.To() == BQS_KING_TO_SQUARE {
			inBetweenSq = m.From() + 1
		} else if m.To() == WQS_KING_TO_SQUARE || m.To() == BQS_KING_TO_SQUARE {
			inBetweenSq = m.From() - 1
		}
		inBetweenMove := NewMove(m.From(), inBetweenSq, EMPTY)
		if b.IsCheck || b.WillMoveCauseCheck(inBetweenMove) {
			return false
		}
	}
	return !b.WillMoveCauseCheck(m)
}

func (b *Board) WillMoveCauseCheck(m Move) bool {
	clone := CloneBoard(b)
	clone.justMove(m)
	return clone.ComputeIsCheck() == true
}

//
func (b *Board) MakeMove(m Move) {
	b.justMove(m)
	kind := b.Squares[m.To()].Kind()
	if kind == KING {
		if b.Turn == WHITE {
			b.Castling &= ^(CASTLE_WKS | CASTLE_WQS)
		} else {
			b.Castling &= ^(CASTLE_BKS | CASTLE_BQS)
		}
	}
	if kind == ROOK {
		switch m.From() {
		case WKS_ROOK_ORIGINAL_SQUARE:
			b.Castling &= ^CASTLE_WKS
		case WQS_ROOK_ORIGINAL_SQUARE:
			b.Castling &= ^CASTLE_WQS
		case BKS_ROOK_ORIGINAL_SQUARE:
			b.Castling &= ^CASTLE_BKS
		case BQS_ROOK_ORIGINAL_SQUARE:
			b.Castling &= ^CASTLE_BQS
		}
	}

	switch m.To() {
	case WKS_ROOK_ORIGINAL_SQUARE:
		b.Castling &= ^CASTLE_WKS
	case WQS_ROOK_ORIGINAL_SQUARE:
		b.Castling &= ^CASTLE_WQS
	case BKS_ROOK_ORIGINAL_SQUARE:
		b.Castling &= ^CASTLE_BKS
	case BQS_ROOK_ORIGINAL_SQUARE:
		b.Castling &= ^CASTLE_BQS
	}
	// enPassant target
	b.EnPassant = 0
	if kind == PAWN && b.Turn == WHITE {
		if m.From().Rank() == 6 && m.To().Rank() == 4 {
			b.EnPassant = m.From() - 8
		}
	}
	if kind == PAWN && b.Turn == BLACK {
		if m.From().Rank() == 1 && m.To().Rank() == 3 {
			b.EnPassant = m.From() + 8
		}
	}

	if b.Turn == WHITE {
		b.Turn = BLACK
	} else {
		b.Turn = WHITE
	}

	b.GenerateLegalMoves()

	b.IsCheck = b.ComputeIsCheck()
	b.IsCheckmate = b.IsCheck && !b.hasMoves()
	b.IsStalement = !b.IsCheckmate && !b.hasMoves()
	b.IsMaterialDraw = b.hasInsufficientMaterial()
	b.IsFinished = (b.IsCheckmate || b.IsStalement || b.IsMaterialDraw)
}

//
func (b *Board) justMove(m Move) {
	from := m.From()
	to := m.To()

	//capturedPiece := m.captured()
	capturedPiece := b.Squares[to]
	if m.IsEnPassant() && b.Turn == WHITE {
		capturedPiece = b.Squares[to+8]
	} else if m.IsEnPassant() && b.Turn == BLACK {
		capturedPiece = b.Squares[to-8]
	}
	fromBBNeg := ^Bitboard(0x1 << from)
	toBB := Bitboard(0x1 << to)
	movingPiece := b.Squares[from]
	movingPieceKind := movingPiece.Kind()
	switch movingPiece.Color() {
	case WHITE:
		// update bitmap of moving piece kind, unset bit of source square
		b.Whites[movingPieceKind] &= fromBBNeg
		// update bitmap of moving piece kind, set bit of source square
		b.Whites[movingPieceKind] |= toBB
		// update white pieces bitboard - unset old square
		b.WhitePieces &= fromBBNeg
		// update white pieces bitboard - set new square
		b.WhitePieces |= toBB
	case BLACK:
		b.Blacks[movingPieceKind] &= fromBBNeg
		b.Blacks[movingPieceKind] |= toBB
		b.BlackPieces &= fromBBNeg
		b.BlackPieces |= toBB
	}

	b.Occupied &= fromBBNeg
	b.Occupied |= toBB

	b.Squares[m.To()] = b.Squares[m.From()]
	b.Squares[m.From()] = EMPTY
	if capturedPiece != EMPTY {
		if !m.IsEnPassant() {
			b.capturePiece(to, capturedPiece)
		} else {
			if b.Turn == WHITE {
				b.capturePiece(to+8, b.Squares[to+8])
				b.Squares[to+8] = EMPTY
			} else {
				b.capturePiece(to-8, b.Squares[to-8])
				b.Squares[to-8] = EMPTY
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
		b.justMove(rookMove)
	}
	var promoteTo Piece = m.GetPromotionTo()
	if promoteTo > 0 {
		switch b.Squares[to].Color() {
		case WHITE:
			// remove advanced pawn from boards
			b.Whites[PAWN] &= ^toBB
			// add promotePiece to board
			b.Whites[promoteTo] |= toBB
			b.WhitePieces |= toBB
		case BLACK:
			// remove advanced pawn from boards
			b.Blacks[PAWN] &= ^toBB
			// add promotePiece to board
			b.Blacks[promoteTo] |= toBB
			b.BlackPieces |= toBB
		}
		b.Squares[m.To()] = Piece(uint(promoteTo) | uint(b.Turn))
	}
}

//Remove captured piece from opponent's pieces
func (b *Board) capturePiece(sq Square, captured Piece) {
	if captured == EMPTY {
		return
	}
	sqBB := Bitboard(0x1 << sq)
	kind := captured.Kind()
	switch captured.Color() {
	case WHITE:
		b.Whites[kind] &= ^sqBB
		b.WhitePieces &= ^sqBB
	case BLACK:
		b.Blacks[kind] &= ^sqBB
		b.BlackPieces &= ^sqBB
	}
}

func (b *Board) hasInsufficientMaterial() bool {
	if b.Whites[QUEEN] > 0 || b.Whites[ROOK] > 0 || b.Whites[PAWN] > 0 {
		return false
	}
	if b.Blacks[QUEEN] > 0 || b.Blacks[ROOK] > 0 || b.Blacks[PAWN] > 0 {
		return false
	}
	if b.Whites[KNIGHT] > 0 && b.Whites[BISHOP] > 0 {
		return false
	}
	if b.Blacks[KNIGHT] > 0 && b.Blacks[BISHOP] > 0 {
		return false
	}

	if b.Whites[BISHOP].NumberOfSetBits() > 1 {
		return false
	}

	if b.Blacks[BISHOP].NumberOfSetBits() > 1 {
		return false
	}

	if b.Whites[KNIGHT].NumberOfSetBits() > 1 {
		return false
	}

	if b.Blacks[KNIGHT].NumberOfSetBits() > 1 {
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
func (b *Board) GetMoveSan(m Move) string {
	pgn := ""
	from := m.From()
	to := m.To()

	if m.IsCastlingMove() {
		if m.To() == WKS_KING_TO_SQUARE || m.To() == BKS_KING_TO_SQUARE {
			pgn += "O-O"
		}
		if m.To() == WQS_KING_TO_SQUARE || m.To() == BQS_KING_TO_SQUARE {
			pgn += "O-O-O"
		}
	} else {
		movingKind := b.Squares[from].Kind()
		if movingKind != PAWN { // -------> 1.
			pgn += strings.ToUpper(string(b.Squares[from].ToRune()))
		}

		// Disambiguation:
		othersOfSameKind, onSameFileCount, onSameRankCount := b.GetOthersOfSameKindMovingToSameTargetCounts(m)
		if othersOfSameKind > 0 {
			if onSameFileCount == 0 {
				pgn += m.From().FileLetter() // -------> 1.1
			} else if onSameRankCount == 0 {
				pgn += m.From().FileLetter() // -------> 1.2
			} else {
				pgn += m.From().Coords() // -------> 1.2
			}
		}

		isCapturing := m.GetCapturedPiece() != EMPTY
		if isCapturing && movingKind == PAWN {
			pgn += from.FileLetter() // -------> 2.
		}
		if isCapturing {
			pgn += "x" // -------> 3.
		}

		pgn += to.Coords() // -------> 4.

		if m.IsPromotionMove() {
			pgn += "=" + strings.ToUpper(string(m.GetPromotionTo())) // -------> 5.
		}
	}

	clone := CloneBoard(b)
	clone.MakeMove(m)
	if clone.IsCheckmate {
		pgn += "#"
	} else if clone.IsCheck {
		pgn += "+"
	}

	return pgn
}

func (b *Board) GetOthersOfSameKindMovingToSameTargetCounts(themove Move) (otherOfSameKind int, onSameFileCount int, onSameRankCount int) {
	movingPiece := b.Squares[themove.From()]
	to := themove.To()
	for _, m := range b.LegalMoves {
		if m == themove || m.To() != to || b.Squares[m.From()].Kind() != movingPiece.Kind() {
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
