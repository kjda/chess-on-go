package chessongo

type Ray struct {
	square    Square
	direction Direction
}

func (r *Ray) step() (bool, int) {
	rank := int(r.square)/8 + DIRECTION_SHIFT[r.direction][0]
	file := int(r.square)%8 + DIRECTION_SHIFT[r.direction][1]

	if IsCoordsOutofBoard(rank, file) {
		return false, -1

	}

	r.square = Square(rank*8 + file)
	return true, int(r.square)
}
