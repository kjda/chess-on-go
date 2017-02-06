package chessongo

const (
	WHITE    = 8  //White 0000 1000
	BLACK    = 16 //BLACK 0001 0000
	NO_COLOR = 0
)

const (
	EMPTY  = iota // 000
	PAWN          // 001
	KNIGHT        // 010
	BISHOP        // 011
	ROOK          // 100
	QUEEN         // 101
	KING          // 110
)

const WHITE_MASK = 0x8
const BLACK_MASK = 0x10

const (
	W_PAWN   = WHITE | PAWN   // 01001
	W_KNIGHT = WHITE | KNIGHT // 01010
	W_BISHOP = WHITE | BISHOP // 01011
	W_ROOK   = WHITE | ROOK   // 01100
	W_QUEEN  = WHITE | QUEEN  // 01101
	W_KING   = WHITE | KING   // 01110

	B_PAWN   = BLACK | PAWN   // 10001
	B_KNIGHT = BLACK | KNIGHT // 10010
	B_BISHOP = BLACK | BISHOP // 10011
	B_ROOK   = BLACK | ROOK   // 10100
	B_QUEEN  = BLACK | QUEEN  // 10101
	B_KING   = BLACK | KING   // 10110
)

type Color uint8

type Piece uint8

func (p Piece) Kind() Piece {
	return Piece(p & (0x7))
}

func (p Piece) IsWhite() bool {
	return p&WHITE_MASK > 0
}

func (p Piece) IsBlack() bool {
	return p&BLACK_MASK > 0
}

func (p Piece) Color() Color {
	if p&WHITE_MASK > 0 {
		return WHITE
	}
	if p&BLACK_MASK > 0 {
		return BLACK
	}
	return NO_COLOR
}

func (p Piece) ToRune() rune {
	r, ok := PIECE_TO_RUNE[p]
	if !ok {
		return ' '
	}
	return r
}

func (p Piece) ToString() string {
	r, ok := PIECE_TO_RUNE[p]
	if !ok {
		return ""
	}
	if p.IsWhite() {
		return "w" + string(r)
	}
	return "b" + string(r)
}
