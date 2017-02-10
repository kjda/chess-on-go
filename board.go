package chessongo

//Castling permissions
const (
	CASTLE_WKS = 1 //White king side castling  0001
	CASTLE_WQS = 2 //White queen side castling 0010
	CASTLE_BKS = 4 //Black king side castling  0100
	CASTLE_BQS = 8 //Black queen side castling 1000
)

//castling squares
const (
	W_KING_INIT_SQUARE       = 60 // e1
	B_KING_INIT_SQUARE       = 4  // e8
	WKS_KING_TO_SQUARE       = 62 // g1
	WQS_KING_TO_SQUARE       = 58 // c1
	BKS_KING_TO_SQUARE       = 6  // g8
	BQS_KING_TO_SQUARE       = 2  // c8
	WKS_ROOK_ORIGINAL_SQUARE = 63 // h1
	WQS_ROOK_ORIGINAL_SQUARE = 56 // a1
	BKS_ROOK_ORIGINAL_SQUARE = 7  // h8
	BQS_ROOK_ORIGINAL_SQUARE = 0  // a8
)

type Board struct {
	Fen         string
	WhitePieces Bitboard
	BlackPieces Bitboard
	// _, pawns, knights, bishops, rooks, queens, king
	Whites         [7]Bitboard
	Blacks         [7]Bitboard
	Occupied       Bitboard
	Squares        [64]Piece
	EnPassant      Square
	Castling       int
	HalfMoves      int
	FullMoves      int
	Turn           Color
	PseudoMoves    []Move
	LegalMoves     []Move
	IsCheck        bool
	IsCheckmate    bool
	IsStalement    bool
	IsMaterialDraw bool
	IsFinished     bool
}

func (b *Board) Reset() {
	b.Fen = ""
	b.WhitePieces = Bitboard(0)
	b.BlackPieces = Bitboard(0)
	for _, kind := range [6]Piece{PAWN, KNIGHT, BISHOP, ROOK, QUEEN, KING} {
		b.Whites[kind] = Bitboard(0)
		b.Blacks[kind] = Bitboard(0)
	}
	b.Occupied = Bitboard(0)
	b.Squares = [64]Piece{}
	b.EnPassant = 0
	b.Castling = 0
	b.HalfMoves = 0
	b.FullMoves = 0
	b.Turn = WHITE
	b.PseudoMoves = []Move{}
	b.LegalMoves = []Move{}
	b.IsCheck = false
	b.IsCheckmate = false
	b.IsStalement = false
	b.IsMaterialDraw = false
	b.IsFinished = false
}

func NewBoard() *Board {
	b := Board{}
	b.InitFromFen(STARTING_POSITION_FEN)
	//b.InitFromFen("8/PPPPPPPP/8/8/8/8/8 w - - 0 1")
	return &b
}

func CloneBoard(b *Board) Board {
	clone := Board{
		Fen:            b.Fen,
		WhitePieces:    b.WhitePieces,
		BlackPieces:    b.BlackPieces,
		Whites:         [7]Bitboard{},
		Blacks:         [7]Bitboard{},
		Squares:        [64]Piece{},
		Occupied:       b.Occupied,
		EnPassant:      b.EnPassant,
		Castling:       b.Castling,
		HalfMoves:      b.HalfMoves,
		FullMoves:      b.FullMoves,
		Turn:           b.Turn,
		PseudoMoves:    []Move{},
		LegalMoves:     []Move{},
		IsCheck:        b.IsCheck,
		IsCheckmate:    b.IsCheckmate,
		IsStalement:    b.IsStalement,
		IsMaterialDraw: b.IsMaterialDraw,
		IsFinished:     b.IsFinished,
	}
	copy(clone.Whites[:], b.Whites[:])
	copy(clone.Blacks[:], b.Blacks[:])
	copy(clone.Squares[:], b.Squares[:])
	return clone
}

func (b *Board) addPiece(piece Piece, index int) {
	b.Squares[index] = piece
	if piece == EMPTY {
		return
	}
	bit := Bitboard(0x1 << uint(index))
	kind := piece.Kind()
	switch piece.Color() {
	case WHITE:
		b.Whites[kind] |= bit
		b.WhitePieces |= bit
	case BLACK:
		b.Blacks[kind] |= bit
		b.BlackPieces |= bit
	}
	b.Occupied |= bit
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

func (b *Board) hasMoves() bool {
	return len(b.LegalMoves) > 0
}

func (b *Board) ShouldIncFullMoves(m Move) bool {
	return b.Squares[m.From()].Color() == BLACK
}

func (b *Board) ShouldResetPly(m Move) bool {
	return m.GetCapturedPiece() > 0 || b.Squares[m.From()].Kind() == PAWN
}
