package main

const (
	CASTLE_WKS = 1 //White king side castling 0001
	CASTLE_WQS = 2 //White queen side castling 0010
	CASTLE_BKS = 4 //Black king side castling 0100
	CASTLE_BQS = 8 //Black queen side castling 1000
)

type Board struct {
	WhiteKing    Bitboard
	WhiteQueens  Bitboard
	WhiteRooks   Bitboard
	WhiteBishops Bitboard
	WhiteKnights Bitboard
	WhitePawns   Bitboard

	BlackKing    Bitboard
	BlackQueens  Bitboard
	BlackRooks   Bitboard
	BlackBishops Bitboard
	BlackKnights Bitboard
	BlackPawns   Bitboard

	WhitePieces Bitboard
	BlackPieces Bitboard
	Occupied    Bitboard

	EnPassant Square
	Castling  uint

	HalfMoves int
	FullMoves int

	Turn Color

	Squares [64]Piece

	Moves []Move
}

func NewBoard() *Board {
	b := Board{}
	b.Reset()
	return &b
}

func (b *Board) addPiece(pieceType Piece, index int) {
	bit := Bitboard(0x1 << uint(index))
	switch pieceType {
	case W_PAWN:
		b.WhitePawns = b.WhitePawns | bit
	case W_ROOK:
		b.WhiteRooks = b.WhiteRooks | bit
	case W_QUEEN:
		b.WhiteQueens = b.WhiteQueens | bit
	case W_KNIGHT:
		b.WhiteKnights = b.WhiteKnights | bit
	case W_BISHOP:
		b.WhiteBishops = b.WhiteBishops | bit
	case W_KING:
		b.WhiteKing = b.WhiteKing | bit
	case B_PAWN:
		b.BlackPawns = b.BlackPawns | bit
	case B_ROOK:
		b.BlackRooks = b.BlackRooks | bit
	case B_QUEEN:
		b.BlackQueens = b.BlackQueens | bit
	case B_KNIGHT:
		b.BlackKnights = b.BlackKnights | bit
	case B_BISHOP:
		b.BlackBishops = b.BlackBishops | bit
	case B_KING:
		b.BlackKing = b.BlackKing | bit
	}
	if pieceType != EMPTY {
		b.WhitePieces = b.WhitePawns | b.WhiteKnights | b.WhiteBishops | b.WhiteRooks | b.WhiteQueens | b.WhiteKing
		b.BlackPieces = b.BlackPawns | b.BlackKnights | b.BlackBishops | b.BlackRooks | b.BlackQueens | b.BlackKing
		b.Occupied = b.WhitePieces | b.BlackPieces
	}
	b.Squares[index] = pieceType
}

func (b *Board) Reset() {
	b.WhitePawns = Bitboard(0)
	b.WhiteKnights = Bitboard(0)
	b.WhiteBishops = Bitboard(0)
	b.WhiteRooks = Bitboard(0)
	b.WhiteQueens = Bitboard(0)
	b.WhiteKing = Bitboard(0)
	b.BlackPawns = Bitboard(0)
	b.BlackKnights = Bitboard(0)
	b.BlackBishops = Bitboard(0)
	b.BlackRooks = Bitboard(0)
	b.BlackQueens = Bitboard(0)
	b.BlackKing = Bitboard(0)
	b.WhitePieces = Bitboard(0)
	b.BlackPieces = Bitboard(0)
	b.Occupied = Bitboard(0)
	b.EnPassant = 0
	b.Castling = 0
	b.HalfMoves = 0
	b.FullMoves = 0
	b.Moves = make([]Move, 0, 32)
	b.Squares = [64]Piece{}
	b.Turn = WHITE
}

//Get our pawns and opponent's
func (b *Board) GetPawns() (Bitboard, Bitboard) {
	if b.Turn == WHITE {
		return b.WhitePawns, b.BlackPawns
	}
	return b.BlackPawns, b.WhitePawns
}

//Get our color and opponent's
func (b *Board) GetColors() (Color, Color) {
	if b.Turn == WHITE {
		return WHITE, BLACK
	}
	return BLACK, WHITE
}
