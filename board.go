package chessongo

import (
	"math/rand"
	"sync"
)

// Castling permissions
const (
	CASTLE_WKS = 1 //White king side castling  0001
	CASTLE_WQS = 2 //White queen side castling 0010
	CASTLE_BKS = 4 //Black king side castling  0100
	CASTLE_BQS = 8 //Black queen side castling 1000
)

// castling squares
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

type Game struct {
	Fen         string
	WhitePieces Bitboard
	BlackPieces Bitboard
	// _, pawns, knights, bishops, rooks, queens, king
	Whites                [7]Bitboard
	Blacks                [7]Bitboard
	Occupied              Bitboard
	Squares               [64]Piece
	EnPassant             Square
	Castling              int
	HalfMoves             int
	FullMoves             int
	Turn                  Color
	PseudoMoves           []Move
	LegalMoves            []Move
	PositionHistory       map[uint64]int
	ZobristHash           uint64
	IsCheck               bool
	IsCheckmate           bool
	IsStalement           bool
	IsMaterialDraw        bool
	IsThreefoldRepetition bool
	IsFiftyMoveRule       bool
	IsSeventyFiveMoveRule bool
	IsFinished            bool
	History               []GameState
}

func (g *Game) Reset() {
	g.Fen = ""
	g.WhitePieces = 0
	g.BlackPieces = 0
	for _, kind := range [6]Piece{PAWN, KNIGHT, BISHOP, ROOK, QUEEN, KING} {
		g.Whites[kind] = 0
		g.Blacks[kind] = 0
	}
	g.Occupied = 0
	g.Squares = [64]Piece{}
	g.EnPassant = 0
	g.Castling = 0
	g.HalfMoves = 0
	g.FullMoves = 0
	g.Turn = WHITE
	g.PseudoMoves = []Move{}
	g.LegalMoves = []Move{}
	g.PositionHistory = map[uint64]int{}
	g.ZobristHash = 0
	g.IsCheck = false
	g.IsCheckmate = false
	g.IsStalement = false
	g.IsMaterialDraw = false
	g.IsThreefoldRepetition = false
	g.IsFiftyMoveRule = false
	g.IsSeventyFiveMoveRule = false
	g.IsFinished = false
	g.History = []GameState{}
}

func NewGame() *Game {
	g := Game{}
	g.LoadFen(STARTING_POSITION_FEN)
	//g.LoadFen("8/PPPPPPPP/8/8/8/8/8 w - - 0 1")
	return &g
}

func CloneGame(g *Game) Game {
	clone := Game{
		Fen:             g.Fen,
		WhitePieces:     g.WhitePieces,
		BlackPieces:     g.BlackPieces,
		Whites:          [7]Bitboard{},
		Blacks:          [7]Bitboard{},
		Squares:         [64]Piece{},
		Occupied:        g.Occupied,
		EnPassant:       g.EnPassant,
		Castling:        g.Castling,
		HalfMoves:       g.HalfMoves,
		FullMoves:       g.FullMoves,
		Turn:            g.Turn,
		PseudoMoves:     []Move{},
		LegalMoves:      []Move{},
		PositionHistory: map[uint64]int{},
		ZobristHash:     g.ZobristHash,
		IsCheck:         g.IsCheck,
		IsCheckmate:     g.IsCheckmate,
		IsStalement:     g.IsStalement,
		IsMaterialDraw:  g.IsMaterialDraw,
		IsFinished:      g.IsFinished,
		History:         make([]GameState, len(g.History)),
	}
	copy(clone.Whites[:], g.Whites[:])
	copy(clone.Blacks[:], g.Blacks[:])
	copy(clone.Squares[:], g.Squares[:])
	copy(clone.History, g.History)
	for k, v := range g.PositionHistory {
		clone.PositionHistory[k] = v
	}
	return clone
}

func (g *Game) addPiece(piece Piece, index int) {
	g.Squares[index] = piece
	if piece == EMPTY {
		return
	}
	bit := Bitboard(0x1 << uint(index))
	kind := piece.Kind()
	switch piece.Color() {
	case WHITE:
		g.Whites[kind] |= bit
		g.WhitePieces |= bit
	case BLACK:
		g.Blacks[kind] |= bit
		g.BlackPieces |= bit
	}
	g.Occupied |= bit
}

// Get our pawns and opponent's
func (g *Game) GetPawns() (Bitboard, Bitboard) {
	if g.Turn == WHITE {
		return g.Whites[PAWN], g.Blacks[PAWN]
	}
	return g.Blacks[PAWN], g.Whites[PAWN]
}

// Get our color and opponent's
func (g *Game) GetColors() (Color, Color) {
	if g.Turn == WHITE {
		return WHITE, BLACK
	}
	return BLACK, WHITE
}

func (g *Game) hasMoves() bool {
	return len(g.LegalMoves) > 0
}

func (g *Game) ShouldIncFullMoves(m Move) bool {
	return g.Squares[m.From()].Color() == BLACK
}

func (g *Game) ShouldResetPly(m Move) bool {
	return m.GetCapturedPiece() > 0 || g.Squares[m.From()].Kind() == PAWN
}

var (
	zobristOnce       sync.Once
	zobristPiece      [12][64]uint64
	zobristCastling   [16]uint64
	zobristEnPassant  [8]uint64
	zobristTurnToMove uint64
)

func initZobrist() {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 12; i++ {
		for j := 0; j < 64; j++ {
			zobristPiece[i][j] = rng.Uint64()
		}
	}
	for i := 0; i < 16; i++ {
		zobristCastling[i] = rng.Uint64()
	}
	for i := 0; i < 8; i++ {
		zobristEnPassant[i] = rng.Uint64()
	}
	zobristTurnToMove = rng.Uint64()
}

func ensureZobrist() {
	zobristOnce.Do(initZobrist)
}

func zobristPieceIndex(p Piece) int {
	switch p {
	case W_PAWN:
		return 0
	case W_KNIGHT:
		return 1
	case W_BISHOP:
		return 2
	case W_ROOK:
		return 3
	case W_QUEEN:
		return 4
	case W_KING:
		return 5
	case B_PAWN:
		return 6
	case B_KNIGHT:
		return 7
	case B_BISHOP:
		return 8
	case B_ROOK:
		return 9
	case B_QUEEN:
		return 10
	case B_KING:
		return 11
	default:
		return -1
	}
}

func (g *Game) computeZobrist() uint64 {
	ensureZobrist()
	h := uint64(0)
	for sq, piece := range g.Squares {
		idx := zobristPieceIndex(piece)
		if idx >= 0 {
			h ^= zobristPiece[idx][sq]
		}
	}

	h ^= zobristCastling[g.Castling&0xF]

	if g.EnPassant != 0 {
		file := g.EnPassant.File()
		h ^= zobristEnPassant[file]
	}

	if g.Turn == BLACK {
		h ^= zobristTurnToMove
	}

	return h
}

func (g *Game) recordPosition() {
	if g.PositionHistory == nil {
		g.PositionHistory = map[uint64]int{}
	}
	g.ZobristHash = g.computeZobrist()
	g.PositionHistory[g.ZobristHash] = g.PositionHistory[g.ZobristHash] + 1
}

func (g *Game) checkThreefoldRepetition() bool {
	return g.PositionHistory != nil && g.PositionHistory[g.ZobristHash] >= 3
}

func (g *Game) IsFivefoldRepetition() bool {
	return g.PositionHistory != nil && g.PositionHistory[g.ZobristHash] >= 5
}

func (g *Game) checkFiftyMoveRule() bool {
	return g.HalfMoves >= 100
}

func (g *Game) checkSeventyFiveMoveRule() bool {
	return g.HalfMoves >= 150
}

type GameState struct {
	CapturedPiece Piece
	Castling      int
	EnPassant     Square
	HalfMoves     int
	ZobristHash   uint64
}
