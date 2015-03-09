package main

//Generate all Pesudo-lega moves
func (b *Board) GenMoves() {
	b.GenPawnMoves()
	b.GenCastling()
	if b.Turn == WHITE {
		b.GenFromMoves(b.WhiteKing, b.WhitePieces, KING_ATTACKS_FROM[:])
		b.GenFromMoves(b.WhiteKnights, b.WhitePieces, KNIGHT_ATTACKS_FROM[:])
		b.GenRayMoves(b.WhiteBishops|b.WhiteQueens, b.WhitePieces, BISHOP_DIRECTIONS[:])
		b.GenRayMoves(b.WhiteRooks|b.WhiteQueens, b.WhitePieces, ROOK_DIRECTIONS[:])
	} else {
		b.GenFromMoves(b.BlackKing, b.BlackPieces, KING_ATTACKS_FROM[:])
		b.GenFromMoves(b.BlackKnights, b.BlackPieces, KNIGHT_ATTACKS_FROM[:])
		b.GenRayMoves(b.BlackBishops|b.BlackQueens, b.BlackPieces, BISHOP_DIRECTIONS[:])
		b.GenRayMoves(b.BlackRooks|b.BlackQueens, b.BlackPieces, ROOK_DIRECTIONS[:])
	}
	moveBB := Bitboard(0)

	for _, m := range b.Moves {
		bit := uint(m.to())
		moveBB |= (1 << bit)
	}
	PrintBitboard(moveBB, "MOVES")
}

func (b *Board) addMoves(from uint, movesBB Bitboard) {
	for movesBB > 0 {
		to := movesBB.popLSB()
		b.Moves = append(b.Moves, NewMove(Square(from), Square(to)))
	}
}

//Generates King & Knight pseudo-legal moves
func (b *Board) GenFromMoves(pieces, ours Bitboard, attackFrom []Bitboard) {
	for pieces > 0 {
		from := pieces.popLSB()
		targets := attackFrom[from] & ^ours
		b.addMoves(from, targets)
	}
}

//Generate sliding-piece's pseudo-legal moves
func (b *Board) GenRayMoves(pieces, ours Bitboard, directions []Direction) {
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
		b.addMoves(from, allTargets)
	}
}

//Generate castling pseudo-legal moves
func (b *Board) GenCastling() {
	if b.Turn == WHITE && (b.Castling&CASTLE_WKS) > 0 && (b.Occupied&(0x2<<61)) == 0 {
		//@todo: generate move
	}
	if b.Turn == WHITE && (b.Castling&CASTLE_WQS) > 0 && (b.Occupied&(0x3<<57)) == 0 {
		//@todo: generate move
	}
	if b.Turn == BLACK && (b.Castling&CASTLE_BKS) > 0 && (b.Occupied&(0x3<<1)) == 0 {
		//@todo: generate move
	}
	if b.Turn == BLACK && (b.Castling&CASTLE_BQS) > 0 && (b.Occupied&(0x2<<5)) == 0 {
		//@todo: generate move
	}
}

//Generate all pawn moves
func (b *Board) GenPawnMoves() {
	b.GenPawnOneStep()
	b.GenPawnTwoSteps()
	b.GenPawnAttacks()
}

//Generate Pawn-one-step-forward pseudo-legal moves
func (b *Board) GenPawnOneStep() {
	var targets Bitboard
	var shift int = 8
	if b.Turn == WHITE {
		targets = (b.WhitePawns >> 8) & ^b.Occupied
	} else {
		targets = (b.BlackPawns << 8) & ^b.Occupied
		shift = -8
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		b.Moves = append(b.Moves, NewMove(Square(from), Square(to)))
	}
}

//Generate Pawn-two-step-forward pseudo-legal moves
func (b *Board) GenPawnTwoSteps() {
	var targets Bitboard
	var shift int
	if b.Turn == WHITE {
		targets = ((b.WhitePawns & RANK1_MASK) >> 16) & (^(RANK2_MASK & b.Occupied) >> 8) & ^b.Occupied
		shift = 16
	} else {
		targets = ((b.BlackPawns & RANK6_MASK) << 16) & (^(RANK5_MASK & b.Occupied) >> 8) & ^b.Occupied
		shift = -16
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		b.Moves = append(b.Moves, NewMove(Square(from), Square(to)))
	}
}

//Generate pawns left and right attacks
func (b *Board) GenPawnAttacks() {
	ours, theirs := b.GetPawns()
	var targets Bitboard
	enPassant := Bitboard(0x1 << uint(b.EnPassant))
	for _, shift := range [2]int{7, 9} {
		if b.Turn == WHITE {
			targets = Bitboard(ours>>uint(shift)) & (theirs | enPassant)
		} else {
			targets = Bitboard(ours<<uint(shift)) & (theirs | enPassant)
			shift *= -1
		}
		for targets > 0 {
			to := Square(targets.popLSB())
			from := Square(int(to) + shift)
			var enPassant Square = 0
			if b.EnPassant > 0 && to == b.EnPassant {
				enPassant = Square(b.EnPassant)
			}
			b.Moves = append(b.Moves, NewMove(from, to, enPassant))
		}
	}
}
