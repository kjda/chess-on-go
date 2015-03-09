package main

/*
* FROM   0-5 bits
* TO     6-10 bits
 */
type Move uint64

func NewMove(from, to Square, enPassant ...Square) Move {
	m := Move(uint64(from) | (uint64(to) << 6))
	if len(enPassant) > 0 {
		m |= Move(enPassant[0] << 12)
	}
	return m
}

func (m Move) from() Square {
	return Square(m & 0x3F)
}

func (m Move) to() Square {
	return Square((m & 0xFC0) >> 6)
}
