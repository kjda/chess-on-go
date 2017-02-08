package chessongo

const (
	RANK1_MASK Bitboard = 0xFF00000000000000 //1111111100000000000000000000000000000000000000000000000000000000
	RANK2_MASK Bitboard = 0xFF000000000000   //0000000011111111000000000000000000000000000000000000000000000000
	RANK3_MASK Bitboard = 0xFF0000000000
	RANK4_MASK Bitboard = 0xFF00000000
	RANK5_MASK Bitboard = 0xFF000000
	RANK6_MASK Bitboard = 0xFF0000
	RANK7_MASK Bitboard = 0xFF00
	RANK8_MASK Bitboard = 0xFF //0000000000000000000000000000000000000000000000000000000011111111

	FILE_A_MASK = 0x0101010101010101 //00000000000100000001000000010000000100000001000000010000000100000001
	FILE_B_MASK = 0x0202020202020202
	FILE_C_MASK = 0x0404040404040404
	FILE_D_MASK = 0x0808080808080808
	FILE_E_MASK = 0x1010101010101010
	FILE_F_MASK = 0x2020202020202020
	FILE_G_MASK = 0x4040404040404040
	FILE_H_MASK = 0x8080808080808080 //1000000010000000100000001000000010000000100000001000000010000000
)

//Attack maps for every possible position
var KNIGHT_ATTACKS_FROM = [64]Bitboard{}
var KING_ATTACKS_FROM = [64]Bitboard{}
var BISHOP_ATTACKS_FROM = [64]Bitboard{}
var ROOK_ATTACKS_FROM = [64]Bitboard{}

var ATTACKS_TO = [64]Bitboard{}

var RANK_ATTACKS = [64][64]Bitboard{}

//Sliding pieces Ray masks
var RAY_MASKS = [8][64]Bitboard{}

//initialize masks
func init() {
	initAttacksFrom()
	initRayMasks()
	initAttacksTo()
}

//Initialize attack maps
func initAttacksFrom() {
	initKingKnightAttacksFrom(KING_ATTACKS_FROM[:], KING_RANK_FILE_SHIFTS)
	initKingKnightAttacksFrom(KNIGHT_ATTACKS_FROM[:], KNIGHT_RANK_FILE_SHIFTS)
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
		rank, file = square.Rank(), square.File()
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

func initAttacksTo() {
	for i := 0; i < 64; i++ {
		ATTACKS_TO[i] = Bitboard(0)
		for _, direction := range ALL_DIRECTIONS {
			ray := Ray{square: Square(i), direction: Direction(direction)}
			for {
				onBoard, nextRaySquare := ray.step()
				if !onBoard {
					break
				}
				ATTACKS_TO[i] |= Bitboard(1) << uint(nextRaySquare)
			}
		}
	}
	var rank, file, bitRank, bitFile int
	var square Square
	for i := 0; i < 64; i++ {
		square = Square(i)
		rank, file = square.Rank(), square.File()
		for _, shift := range KING_RANK_FILE_SHIFTS {
			bitRank, bitFile = (rank + shift[0]), (file + shift[1])
			if IsCoordsOutofBoard(bitRank, bitFile) {
				continue
			}
			ATTACKS_TO[i] |= (1 << CoordsToIndex(bitRank, bitFile))
		}
	}
	for i := 0; i < 64; i++ {
		square = Square(i)
		rank, file = square.Rank(), square.File()
		for _, shift := range KNIGHT_RANK_FILE_SHIFTS {
			bitRank, bitFile = (rank + shift[0]), (file + shift[1])
			if IsCoordsOutofBoard(bitRank, bitFile) {
				continue
			}
			ATTACKS_TO[i] |= (1 << CoordsToIndex(bitRank, bitFile))
		}
	}
}
