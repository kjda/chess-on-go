package main

const (
	RANK0_MASK Bitboard = 0xFF00000000000000 //1111111100000000000000000000000000000000000000000000000000000000
	RANK1_MASK Bitboard = 0xFF000000000000   //0000000011111111000000000000000000000000000000000000000000000000
	RANK2_MASK Bitboard = 0xFF0000000000
	RANK3_MASK Bitboard = 0xFF00000000
	RANK4_MASK Bitboard = 0xFF000000
	RANK5_MASK Bitboard = 0xFF0000
	RANK6_MASK Bitboard = 0xFF00
	RANK7_MASK Bitboard = 0xFF

	FILE0_MASK = 0x8080808080808080 //1000000010000000100000001000000010000000100000001000000010000000
	FILE1_MASK = 0x4040404040404040
	FILE2_MASK = 0x2020202020202020
	FILE3_MASK = 0x1010101010101010
	FILE4_MASK = 0x0808080808080808
	FILE5_MASK = 0x0404040404040404
	FILE6_MASK = 0x0202020202020202
	FILE7_MASK = 0x0101010101010101
)

//Attack maps for every possible position
var KNIGHT_ATTACKS_FROM = [64]Bitboard{}
var KING_ATTACKS_FROM = [64]Bitboard{}
var BISHOP_ATTACKS_FROM = [64]Bitboard{}
var ROOK_ATTACKS_FROM = [64]Bitboard{}

var RANK_ATTACKS = [64][64]Bitboard{}

//Sliding pieces Ray masks
var RAY_MASKS = [8][64]Bitboard{}

//initialize masks
func init() {
	initAttacksFrom()
	initRayMasks()
}

//Initialize attack maps
func initAttacksFrom() {
	var kingShifts = [8][2]int{
		{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
		{1, 0}, {-1, 1}, {0, 1}, {1, 1},
	}
	initKingKnightAttacksFrom(KING_ATTACKS_FROM[:], kingShifts)

	var knightShifts = [8][2]int{
		{-2, -1}, {-2, 1}, {2, -1}, {2, 1},
		{-1, -2}, {-1, 2}, {1, -2}, {1, 2},
	}
	initKingKnightAttacksFrom(KNIGHT_ATTACKS_FROM[:], knightShifts)

	initRayAttackFrom(BISHOP_ATTACKS_FROM[:], BISHOP_DIRECTIONS[:])
	initRayAttackFrom(ROOK_ATTACKS_FROM[:], ROOK_DIRECTIONS[:])
}

/*********************************
*    Example: AttackFrom of a Knight on D5 (Index: 27)
*    0 0 0 0 0 0 0 0
*    0 0 0 0 0 0 0 0
*    0 0 0 1 0 1 0 0
*    0 0 1 0 0 0 1 0
*    0 0 0 0 * 0 0 0
*    0 0 1 0 0 0 1 0
*    0 0 0 1 0 1 0 0
*    0 0 0 0 0 0 0 0
*
*********************************/
func initKingKnightAttacksFrom(attacksFrom []Bitboard, shifts [8][2]int) {
	var rank, file, bitRank, bitFile int
	var square Square
	for squareIdx := 0; squareIdx < 64; squareIdx++ {
		square = Square(squareIdx)
		attacksFrom[squareIdx] = Bitboard(0)
		rank, file = square.rank(), square.file()
		for _, shift := range shifts {
			bitRank, bitFile = (rank + shift[0]), (file + shift[1])
			if IsCoordsOutofBoard(bitRank, bitFile) {
				continue
			}
			attacksFrom[squareIdx] |= (1 << CoordsToIndex(bitRank, bitFile))
		}
	}
}

/*********************************
*    Example Ray Attack of a Rook on D5 (Index: 27)
*    0 0 0 0 1 0 0 0
*    0 0 0 0 1 0 0 0
*    0 0 0 0 1 0 0 0
*    0 0 0 0 1 0 0 0
*    1 1 1 1 * 1 1 1
*    0 0 0 0 1 0 0 0
*    0 0 0 0 1 0 0 0
*    0 0 0 0 1 0 0 0
*
*********************************/
func initRayAttackFrom(attacksFrom []Bitboard, rayDirections []Direction) {
	for i := 0; i < 64; i++ {
		attacksFrom[i] = Bitboard(0)
		for _, direction := range rayDirections {
			ray := Ray{square: Square(i), direction: Direction(direction)}
			for {
				onBoard, nextRaySquare := ray.step()
				if !onBoard {
					break
				}
				attacksFrom[i] |= Bitboard(1) << uint(nextRaySquare)
			}
		}
	}

}

/*********************************
*    Example RayMask of square D5 (Index: 27) in North West direction
*    1 0 0 0 0 0 0 0
*    0 1 0 0 0 0 0 0
*    0 0 1 0 0 0 0 0
*    0 0 0 1 0 0 0 0
*    0 0 0 0 * 0 0 0
*    0 0 0 0 0 0 0 0
*    0 0 0 0 0 0 0 0
*    0 0 0 0 0 0 0 0
*
*********************************/
func initRayMasks() {
	for _, direction := range ALL_DIRECTIONS {
		for square := 0; square < 64; square++ {
			ray := Ray{square: Square(square), direction: Direction(direction)}
			for {
				onBoard, raySquare := ray.step()
				if !onBoard {
					break
				}
				RAY_MASKS[direction][square] |= Bitboard(1) << uint(raySquare)
			}
		}

	}
}
