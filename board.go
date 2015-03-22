package main

//Castling permissions
const (
	CASTLE_WKS = 1 //White king side castling 0001
	CASTLE_WQS = 2 //White queen side castling 0010
	CASTLE_BKS = 4 //Black king side castling 0100
	CASTLE_BQS = 8 //Black queen side castling 1000
)

//castling squares
const (
	WKS_KING_SQUARE = 57
	WQS_KING_SQUARE = 61
	BKS_KING_SQUARE = 6
	BQS_KING_SQUARE = 2
	WKS_ROOK_SQUARE = 56
	WQS_ROOK_SQUARE = 63
	BKS_ROOK_SQUARE = 7
	BQS_ROOK_SQUARE = 0
)

type Board struct {
	WhitePieces Bitboard
	BlackPieces Bitboard

	Whites [7]Bitboard
	Blacks [7]Bitboard

	Occupied Bitboard

	EnPassant Square
	Castling  uint

	HalfMoves int
	FullMoves int

	Turn Color

	Check     bool
	Checkmate bool
	Stalement bool

	Squares [64]Piece

	PossibleMoves MoveList

	MoveHistory MoveList
}

func NewBoard() *Board {
	b := Board{}
	b.Reset()
	return &b
}

func (b *Board) addPiece(piece Piece, index int) {
	b.Squares[index] = piece
	if piece == EMPTY {
		return
	}
	bit := Bitboard(0x1 << uint(index))
	kind := piece.kind()
	switch piece.color() {
	case WHITE:
		b.Whites[kind] |= bit
		b.WhitePieces |= bit
	case BLACK:
		b.Blacks[kind] |= bit
		b.BlackPieces |= bit
	}
	b.Occupied |= bit
}

func (b *Board) Reset() {
	for _, kind := range [6]Piece{PAWN, KNIGHT, BISHOP, ROOK, QUEEN, KING} {
		b.Whites[kind] = Bitboard(0)
		b.Blacks[kind] = Bitboard(0)
	}

	b.WhitePieces = Bitboard(0)
	b.BlackPieces = Bitboard(0)
	b.Occupied = Bitboard(0)
	b.EnPassant = 0
	b.Castling = 0
	b.HalfMoves = 0
	b.FullMoves = 0
	b.PossibleMoves = NewMoveList()
	b.Squares = [64]Piece{}
	b.Turn = WHITE
	b.Check = false
	b.Checkmate = false
	b.Stalement = false
}

//Get our pawns and opponent's
func (b *Board) GetPawns() (Bitboard, Bitboard) {
	if b.Turn == WHITE {
		return b.Whites[PAWN], b.Blacks[PAWN]
	}
	return b.Blacks[PAWN], b.Whites[PAWN]
}

//Get our color and opponent's
func (b *Board) GetColors() (Color, Color) {
	if b.Turn == WHITE {
		return WHITE, BLACK
	}
	return BLACK, WHITE
}

func (b *Board) switchTurn() {
	if b.Turn == WHITE {
		b.Turn = BLACK
		return
	}
	b.Turn = WHITE
}

func (b *Board) hasMoves() bool {
	return b.PossibleMoves.len() > 0
}
