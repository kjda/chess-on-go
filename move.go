package chessongo

import "fmt"

const (
	MOVE_TO_BIT         = 6
	MOVE_ENPASSANT_BIT  = 12
	MOVE_CASTLING_BIT   = 13
	MOVE_PROMOTION_BIT  = 14
	MOVE_CAPTURE_BIT    = 17
	MOVE_TO_FROM_MASK   = 0x3F
	MOVE_CAPTURED_MASK  = 0x1F
	MOVE_PROMOTION_MASK = 0x7
)

/*
* FROM             bits 0-5
* TO               bits 6-11
* EnPassant        bit  12
* Castling         bit  13
* Promotion        bits 14 - 16
* CapturedPieceKind    bits 17 - 20
 */
type Move uint32

func NewMove(from, to Square, captured Piece) Move {
	m := Move(uint32(from) | (uint32(to) << MOVE_TO_BIT) | uint32(captured)<<MOVE_CAPTURE_BIT)
	return m
}

func NewEnPassantMove(from, to Square, captured Piece) Move {
	m := NewMove(from, to, captured)
	m |= Move(uint32(0x1) << MOVE_ENPASSANT_BIT)
	return m
}

func NewPromotionMove(from, to Square, captured Piece, promotionTo Piece) Move {
	m := NewMove(from, to, captured)
	m |= Move(uint32(promotionTo) << MOVE_PROMOTION_BIT)
	return m
}

func NewCastlingMove(from, to Square) Move {
	m := NewMove(from, to, EMPTY)
	m |= Move(0x1 << MOVE_CASTLING_BIT)
	return m
}

func (m Move) From() Square {
	return Square(m & MOVE_TO_FROM_MASK)
}

func (m Move) To() Square {
	return Square((m >> MOVE_TO_BIT) & MOVE_TO_FROM_MASK)
}

func (m Move) IsCastlingMove() bool {
	return (uint32(m)>>MOVE_CASTLING_BIT)&0x1 > 0
}

func (m Move) IsEnPassant() bool {
	return (uint32(m)>>MOVE_ENPASSANT_BIT)&0x1 > 0
}

func (m Move) IsPromotionMove() bool {
	return (uint32(m)>>MOVE_PROMOTION_BIT)&MOVE_PROMOTION_MASK > 0
}

func (m Move) IsPromotingTo(promotingTo string) bool {
	promotingToPiece, ok := STRING_TO_KIND[promotingTo]
	return ok && Piece(promotingToPiece) == m.GetPromotionTo()
}

func (m Move) GetCapturedPiece() Piece {
	return Piece((m >> MOVE_CAPTURE_BIT) & MOVE_CAPTURED_MASK)
}

func (m Move) GetPromotionTo() Piece {
	return Piece((m >> MOVE_PROMOTION_BIT) & MOVE_PROMOTION_MASK)
}

func (m Move) ToString() string {
	return fmt.Sprintf("%s %s", SQUARE_TO_COORDS[m.From()], SQUARE_TO_COORDS[m.To()])
}

func (m Move) ToFromToStrings() (string, string) {
	return SQUARE_TO_COORDS[m.From()], SQUARE_TO_COORDS[m.To()]
}
