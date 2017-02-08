package chessongo

//Tells whether coords are out of board or not
func IsCoordsOutofBoard(rank, file int) bool {
	return rank < 0 || rank > 7 || file < 0 || file > 7
}

//Convert rank, file coordinates to a 64 based square
func CoordsToSquare(rank, file int) Square {
	return Square(rank*8 + file)
}

//Convert rank, file coordinates to a 64 based index
func CoordsToIndex(rank, file int) uint {
	return uint(rank*8 + file)
}

//get rank, file from Square
func squareCoords(sq Square) (int, int) {
	return int(sq) / 8, int(sq) % 8
}

var FILES = map[int]rune{
	0: 'a',
	1: 'b',
	2: 'c',
	3: 'd',
	4: 'e',
	5: 'f',
	6: 'g',
	7: 'h',
}

func file(sq Square) rune {
	return FILES[sq.File()]
}

func rank(sq Square) int {
	return 8 - sq.Rank()
}
