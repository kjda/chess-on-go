package main

const (
	WHITE    = 8  //White 0000 1000
	BLACK    = 16 //BLACK 0001 0000
	NO_COLOR = 0
)

const (
	EMPTY = iota
	PAWN
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

const WHITE_MASK = 0x8
const BLACK_MASK = 0x10

const (
	W_PAWN   = WHITE | PAWN   //1001
	W_KNIGHT = WHITE | KNIGHT //1010
	W_BISHOP = WHITE | BISHOP
	W_ROOK   = WHITE | ROOK
	W_QUEEN  = WHITE | QUEEN
	W_KING   = WHITE | KING

	B_PAWN   = BLACK | PAWN   //10001
	B_KNIGHT = BLACK | KNIGHT //10010
	B_BISHOP = BLACK | BISHOP
	B_ROOK   = BLACK | ROOK
	B_QUEEN  = BLACK | QUEEN
	B_KING   = BLACK | KING
)

type Color uint8

type Piece uint8

func (p Piece) kind() Piece {
	return Piece(p & (0x7))
}

func (p Piece) isWhite() bool {
	return p&WHITE_MASK > 0
}

func (p Piece) isBlack() bool {
	return p&BLACK_MASK > 0
}

func (p Piece) color() Color {
	if p&WHITE_MASK > 0 {
		return WHITE
	}
	if p&BLACK_MASK > 0 {
		return BLACK
	}
	return NO_COLOR
}

func (p Piece) toRune() rune {
	return PIECE_TO_RUNE[p]
}
