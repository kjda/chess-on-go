package chessongo

import "fmt"

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
		to := targets.popLSB()
		from := int(to) + shift
		if (b.Turn == WHITE && (Bitboard(0x1<<to)&RANK8_MASK > 0)) || (b.Turn == BLACK && (Bitboard(0x1<<to)&RANK1_MASK > 0)) {
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(Square(from), Square(to), b.Squares[to], QUEEN))
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(Square(from), Square(to), b.Squares[to], ROOK))
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(Square(from), Square(to), b.Squares[to], KNIGHT))
			b.PseudoMoves = append(b.PseudoMoves, NewPromotionMove(Square(from), Square(to), b.Squares[to], BISHOP))
		} else {
			b.PseudoMoves = append(b.PseudoMoves, NewMove(Square(from), Square(to), b.Squares[to]))
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
	enPassant := Bitboard(0x1 << uint(b.EnPassant))
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
			if from.rank() == to.rank() {
				continue
			}
			if b.EnPassant > 0 && to == b.EnPassant {
				var capturedSq Square
				if b.Turn == WHITE {
					capturedSq = to + 8
				} else {
					capturedSq = to - 8
				}
				b.PseudoMoves = append(b.PseudoMoves, NewEnPassantMove(from, to, b.Squares[capturedSq]))
			} else {
				b.PseudoMoves = append(b.PseudoMoves, NewMove(from, to, b.Squares[to]))
			}
		}
	}
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
	if m.isCastlingMove() {
		var inBetweenSq Square
		if m.to() == WKS_KING_TO_SQUARE || m.to() == BQS_KING_TO_SQUARE {
			inBetweenSq = m.from() + 1
		} else if m.to() == WQS_KING_TO_SQUARE || m.to() == BQS_KING_TO_SQUARE {
			inBetweenSq = m.from() - 1
		}
		inBetweenMove := NewMove(m.from(), inBetweenSq, EMPTY)
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
	kind := b.Squares[m.to()].Kind()
	if kind == KING {
		if b.Turn == WHITE {
			b.Castling &= ^(CASTLE_WKS | CASTLE_WQS)
		} else {
			b.Castling &= ^(CASTLE_BKS | CASTLE_BQS)
		}
	}
	if kind == ROOK {
		switch m.from() {
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

	switch m.to() {
	case WKS_ROOK_ORIGINAL_SQUARE:
		fmt.Println("ROOK TO WKS_ROOK_ORIGINAL_SQUARE")
		b.Castling &= ^CASTLE_WKS
	case WQS_ROOK_ORIGINAL_SQUARE:
		fmt.Println("ROOK TO WQS_ROOK_ORIGINAL_SQUARE")
		b.Castling &= ^CASTLE_WQS
	case BKS_ROOK_ORIGINAL_SQUARE:
		fmt.Println("ROOK TO BKS_ROOK_ORIGINAL_SQUARE")
		b.Castling &= ^CASTLE_BKS
	case BQS_ROOK_ORIGINAL_SQUARE:
		fmt.Println("ROOK TO BQS_ROOK_ORIGINAL_SQUARE")
		b.Castling &= ^CASTLE_BQS
	}
	// enPassant target
	b.EnPassant = 0
	if kind == PAWN && b.Turn == WHITE {
		if m.from().rank() == 6 && m.to().rank() == 4 {
			b.EnPassant = m.from() - 8
		}
	}
	if kind == PAWN && b.Turn == BLACK {
		if m.from().rank() == 1 && m.to().rank() == 3 {
			b.EnPassant = m.from() + 8
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
}

//
func (b *Board) justMove(m Move) {
	from := m.from()
	to := m.to()
	//capturedPiece := m.captured()
	capturedPiece := b.Squares[to]
	if m.isEnPassant() && b.Turn == WHITE {
		capturedPiece = b.Squares[to+8]
	} else if m.isEnPassant() && b.Turn == BLACK {
		capturedPiece = b.Squares[to-8]
	}
	notFromBB := ^Bitboard(0x1 << from)
	toBB := Bitboard(0x1 << to)

	movingPiece := b.Squares[from]
	movingPieceKind := movingPiece.Kind()
	switch movingPiece.Color() {
	case WHITE:
		b.Whites[movingPieceKind] &= notFromBB
		b.Whites[movingPieceKind] |= toBB
		b.WhitePieces &= notFromBB
		b.WhitePieces |= toBB
	case BLACK:
		b.Blacks[movingPieceKind] &= notFromBB
		b.Blacks[movingPieceKind] |= toBB
		b.BlackPieces &= notFromBB
		b.BlackPieces |= toBB
	}

	b.Occupied &= notFromBB
	b.Occupied |= toBB

	b.Squares[m.to()] = b.Squares[m.from()]
	b.Squares[m.from()] = EMPTY
	if capturedPiece != EMPTY {
		if !m.isEnPassant() {
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
