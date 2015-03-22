package main

import (
//"fmt"
)

//Generates legal moves and update board's check, checkmate, stalement and castling status
func (b *Board) UpdateStatus() {
	lastMove := b.MoveHistory.getLast()

	lastMoveKind := b.Squares[lastMove.to()].kind()
	if lastMoveKind == KING {
		if b.Turn == WHITE {
			//b.Castling &= ^(CASTLE_WKS & CASTLE_WQS)
		} else {
			//b.Castling &= ^(CASTLE_WBS & CASTLE_WBS)
		}
	}
	if lastMoveKind == ROOK {
		switch lastMove.from() {
		case WKS_ROOK_SQUARE:
			//b.Castling &= ^CASTLE_WKS
		case WQS_ROOK_SQUARE:
			//b.Castling &= ^CASTLE_WQS
		case BKS_ROOK_SQUARE:
			//b.Castling &= ^CASTLE_BKS
		case BQS_ROOK_SQUARE:
			//b.Castling &= ^CASTLE_BQS
		}
	}

	b.GenerateLegalMoves()
	b.Check = b.IsCheck()
	b.Checkmate = b.Check && !b.hasMoves()
	b.Stalement = !b.Check && !b.hasMoves()

	if lastMove > 0 {
		return
	}
}

var genMovesCalls uint = 0

//Generate all legal moves
func (b *Board) GenerateLegalMoves() {
	var ours [7]Bitboard
	var oursAll Bitboard
	if b.Turn == WHITE {
		ours = b.Whites
		oursAll = b.WhitePieces
	} else {
		ours = b.Blacks
		oursAll = b.BlackPieces
	}

	pseudoMoves := NewMoveList()

	b.genFromMoves(&pseudoMoves, ours[KING], oursAll, KING_ATTACKS_FROM[:])
	b.genFromMoves(&pseudoMoves, ours[KNIGHT], oursAll, KNIGHT_ATTACKS_FROM[:])
	b.genRayMoves(&pseudoMoves, ours[BISHOP]|ours[QUEEN], oursAll, BISHOP_DIRECTIONS[:])
	b.genRayMoves(&pseudoMoves, ours[ROOK]|ours[QUEEN], oursAll, ROOK_DIRECTIONS[:])

	b.genPawnTwoSteps(&pseudoMoves)
	b.genPawnAttacks(&pseudoMoves)
	b.genCastling(&pseudoMoves)
	b.genPawnOneStep(&pseudoMoves)

	b.PossibleMoves = NewMoveList()
	for _, move := range pseudoMoves.Moves {
		if b.CanMove(move) {
			b.PossibleMoves.append(move)
		}
	}
}

//Checks whether the given move is possible or not
func (b *Board) CanMove(m Move) bool {
	b.MakeMove(m, false)
	defer b.UndoMove(m)
	if b.IsCheck() {
		return false
	}
	return true
}

//Makes a move without updating board status
func (b *Board) MakeMove(m Move, updateBoardStatus bool) {
	from := m.from()
	to := m.to()
	captured := m.captured()
	notFromBB := ^Bitboard(0x1 << from)
	toBB := Bitboard(0x1 << to)

	movingPiece := b.Squares[from]
	movingPieceKind := movingPiece.kind()
	switch movingPiece.color() {
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
	if captured != EMPTY {
		b.capturePiece(to, captured)
	}

	if updateBoardStatus {

	}
}

//Undo a move
func (b *Board) UndoMove(m Move) {
	from := m.from()
	to := m.to()
	captured := m.captured()
	fromBB := Bitboard(0x1 << from)
	notToBB := ^Bitboard(0x1 << to)
	movingPiece := b.Squares[to]
	movingPieceKind := movingPiece.kind()
	switch movingPiece.color() {
	case WHITE:
		b.Whites[movingPieceKind] |= fromBB
		b.Whites[movingPieceKind] &= notToBB
		b.WhitePieces |= fromBB
		b.WhitePieces &= notToBB
	case BLACK:
		b.Blacks[movingPieceKind] |= fromBB
		b.Blacks[movingPieceKind] &= notToBB
		b.BlackPieces |= fromBB
		b.BlackPieces &= notToBB
	}

	b.Squares[m.from()] = b.Squares[m.to()]
	b.Squares[m.to()] = captured

	b.Occupied |= fromBB
	if captured == EMPTY {
		b.Occupied &= notToBB
	} else {
		b.uncapturePiece(to, captured)
	}
}

//Remove captured piece from opponent's pieces
func (b *Board) capturePiece(sq Square, captured Piece) {
	if captured == EMPTY {
		return
	}
	sqBB := Bitboard(0x1 << sq)
	kind := captured.kind()
	switch captured.color() {
	case WHITE:
		b.Whites[kind] &= ^sqBB
		b.WhitePieces &= ^sqBB
	case BLACK:
		b.Blacks[kind] &= ^sqBB
		b.BlackPieces &= ^sqBB
	}
}

//Undo capture piece
func (b *Board) uncapturePiece(sq Square, captured Piece) {
	if captured == EMPTY {
		return
	}
	sqBB := Bitboard(0x1 << sq)
	kind := captured.kind()
	switch captured.color() {
	case WHITE:
		b.Whites[kind] |= sqBB
		b.WhitePieces |= sqBB
	case BLACK:
		b.Blacks[kind] |= sqBB
		b.BlackPieces |= sqBB
	}
}

//Generates King & Knight pseudo-legal moves
func (b *Board) genFromMoves(moves *MoveList, pieces, ours Bitboard, attackFrom []Bitboard) {
	for pieces > 0 {
		from := pieces.popLSB()
		targets := attackFrom[from] & ^ours
		for targets > 0 {
			to := targets.popLSB()
			moves.append(NewMove(Square(from), Square(to), b.Squares[to]))
		}
	}
}

//Generate sliding-piece's pseudo-legal moves
func (b *Board) genRayMoves(moves *MoveList, pieces, ours Bitboard, directions []Direction) {
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
			moves.append(NewMove(Square(from), Square(to), b.Squares[to]))
		}
	}
}

//Generate castling pseudo-legal moves
func (b *Board) genCastling(moves *MoveList) {
	if b.Turn == WHITE && (b.Castling&CASTLE_WKS) > 0 && (b.Occupied&(0x2<<61)) == 0 {
		//@todo: generate move
		from := Square(b.Whites[KING].lsbIndex())
		to := Square(WKS_KING_SQUARE)
		moves.append(NewCastlingMove(from, to))

	}
	if b.Turn == WHITE && (b.Castling&CASTLE_WQS) > 0 && (b.Occupied&(0x3<<57)) == 0 {
		//@todo: generate move
		from := Square(b.Whites[KING].lsbIndex())
		to := Square(WQS_KING_SQUARE)
		moves.append(NewCastlingMove(from, to))
	}
	if b.Turn == BLACK && (b.Castling&CASTLE_BKS) > 0 && (b.Occupied&(0x3<<1)) == 0 {
		//@todo: generate move
		from := Square(b.Blacks[KING].lsbIndex())
		to := Square(BKS_KING_SQUARE)
		moves.append(NewCastlingMove(from, to))
	}
	if b.Turn == BLACK && (b.Castling&CASTLE_BQS) > 0 && (b.Occupied&(0x2<<5)) == 0 {
		//@todo: generate move
		from := Square(b.Blacks[KING].lsbIndex())
		to := Square(BQS_KING_SQUARE)
		moves.append(NewCastlingMove(from, to))
	}
}

//Generate Pawn-one-step-forward pseudo-legal moves
func (b *Board) genPawnOneStep(moves *MoveList) {
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
			moves.append(NewPromotionMove(Square(from), Square(to), b.Squares[to], QUEEN))
			moves.append(NewPromotionMove(Square(from), Square(to), b.Squares[to], ROOK))
			moves.append(NewPromotionMove(Square(from), Square(to), b.Squares[to], KNIGHT))
			moves.append(NewPromotionMove(Square(from), Square(to), b.Squares[to], BISHOP))
		} else {
			moves.append(NewMove(Square(from), Square(to), b.Squares[to]))
		}

	}
}

//Generate Pawn-two-step-forward pseudo-legal moves
func (b *Board) genPawnTwoSteps(moves *MoveList) {
	var targets Bitboard
	var shift int
	if b.Turn == WHITE {
		targets = ((b.Whites[PAWN] & RANK2_MASK) >> 16) & (^(RANK3_MASK & b.Occupied) >> 8) & ^b.Occupied
		shift = 16
	} else {
		targets = ((b.Blacks[PAWN] & RANK7_MASK) << 16) & (^(RANK6_MASK & b.Occupied) >> 8) & ^b.Occupied
		shift = -16
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		moves.append(NewMove(Square(from), Square(to), b.Squares[to]))
	}
}

//Generate pawns left and right attacks
func (b *Board) genPawnAttacks(moves *MoveList) {
	ours, theirs := b.GetPawns()
	var targets Bitboard
	enPassant := Bitboard(0x1 << uint(b.EnPassant))
	for _, shift := range [2]int{7, 9} {
		fromShift := shift
		if b.Turn == WHITE {
			targets = Bitboard(ours>>uint(shift)) & (theirs | enPassant)
		} else {
			targets = Bitboard(ours<<uint(shift)) & (theirs | enPassant)
			fromShift *= -1
		}
		for targets > 0 {
			to := Square(targets.popLSB())
			from := Square(int(to) + fromShift)
			if b.EnPassant > 0 && to == b.EnPassant {
				moves.append(NewEnPassantMove(from, to, b.Squares[to]))
			} else {
				moves.append(NewMove(from, to, b.Squares[to]))
			}

		}
	}
}

//Checks whether our king is in check or not
func (b *Board) IsCheck() bool {
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
