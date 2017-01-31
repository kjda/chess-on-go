package chessongo

/*************************************************
*	Bitboard representation
*
*	    H G F E D C B A
*                                                     first bit in rank
*	1   0 0 0 0 0 0 0 0                               <=== 56
*	2   0 0 0 0 0 0 0 0   <== WHITE                   <=== 48
*	3   0 0 0 0 0 0 0 0                               <=== 40
*	4   0 0 0 0 0 0 0 0                               <=== 32
*	5   0 0 0 0 0 0 0 0                               <=== 24
*	6   0 0 0 0 0 0 0 0                               <=== 16
*	7   0 0 0 0 0 0 0 0   <== BLACK                   <=== 8
*	8   0 0 0 0 0 0 0 0                               <=== 0
*
*	bit 0   => A8
*	bit 63  => H1
*

***************************************************/
var SQUARE_TO_COORDS = [64]string{
	"A8", "B8", "C8", "D8", "E8", "F8", "G8", "H8",
	"A7", "B7", "C7", "D7", "E7", "F7", "G7", "H7",
	"A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6",
	"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5",
	"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4",
	"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3",
	"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2",
	"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1",
}

var COORS_TO_SQUARE = map[string]Square{
	"A8": Square(0), "B8": Square(1), "C8": Square(2), "D8": Square(3), "E8": Square(4), "F8": Square(5), "G8": Square(6), "H8": Square(7),
	"A7": Square(8), "B7": Square(9), "C7": Square(10), "D7": Square(11), "E7": Square(12), "F7": Square(13), "G7": Square(14), "H7": Square(15),
	"A6": Square(16), "B6": Square(17), "C6": Square(18), "D6": Square(19), "E6": Square(20), "F6": Square(21), "G6": Square(22), "H6": Square(23),
	"A5": Square(24), "B5": Square(25), "C5": Square(26), "D5": Square(27), "E5": Square(28), "F5": Square(29), "G5": Square(30), "H5": Square(31),
	"A4": Square(32), "B4": Square(33), "C4": Square(34), "D4": Square(35), "E4": Square(36), "F4": Square(37), "G4": Square(38), "H4": Square(39),
	"A3": Square(40), "B3": Square(41), "C3": Square(42), "D3": Square(43), "E3": Square(44), "F3": Square(45), "G3": Square(46), "H3": Square(47),
	"A2": Square(48), "B2": Square(49), "C2": Square(50), "D2": Square(51), "E2": Square(52), "F2": Square(53), "G2": Square(54), "H2": Square(55),
	"A1": Square(56), "B1": Square(57), "C1": Square(58), "D1": Square(59), "E1": Square(60), "F1": Square(61), "G1": Square(62), "H1": Square(63),
}

//"Most significant bit" index lookup table
var MS1BTABLE = [256]int{}

//"Least significant bit" index lookup table
var LS1BTABLE = [64]uint{
	0, 1, 48, 2, 57, 49, 28, 3,
	61, 58, 50, 42, 38, 29, 17, 4,
	62, 55, 59, 36, 53, 51, 43, 22,
	45, 39, 33, 30, 24, 18, 12, 5,
	63, 47, 56, 27, 60, 41, 37, 16,
	54, 35, 52, 21, 44, 32, 23, 11,
	46, 26, 40, 15, 34, 20, 31, 10,
	25, 14, 19, 9, 13, 8, 7, 6,
}

//initilize bitboards
func init() {
	initMostSignificatBit()
}

//Initialze most significant bit lookup table
func initMostSignificatBit() {
	for i := 0; i < 256; i++ {
		if i > 127 {
			MS1BTABLE[i] = 7
		} else if i > 63 {
			MS1BTABLE[i] = 6
		} else if i > 31 {
			MS1BTABLE[i] = 5
		} else if i > 15 {
			MS1BTABLE[i] = 4
		} else if i > 7 {
			MS1BTABLE[i] = 3
		} else if i > 3 {
			MS1BTABLE[i] = 2
		} else if i > 2 {
			MS1BTABLE[i] = 2
		} else if i > 1 {
			MS1BTABLE[i] = 1
		} else {
			MS1BTABLE[i] = 0
		}
	}
}

//Bitboard
type Bitboard uint64

//Get least significant bit
func (bb Bitboard) lsb() Bitboard {
	return bb & (-bb)
}

//Get index of least significant(of Martin Läuter)
func (bb Bitboard) lsbIndex() uint {
	return LS1BTABLE[(bb.lsb()*0x03f79d71b4cb0a89)>>58]
}

//Get index of most significant bit(of Eugene Nalimov)
func (bb Bitboard) msbIndex() int {
	var msb int = 0
	if bb > 0xFFFFFFFF {
		bb >>= 32
		msb = 32
	}
	if bb > 0xFFFF {
		bb >>= 16
		msb += 16
	}
	if bb > 0xFF {
		bb >>= 8
		msb += 8
	}
	return msb + MS1BTABLE[bb]
}

//Pop least significant bit and return it's index
func (bb *Bitboard) popLSB() uint {
	lsb := (*bb).lsb()
	*bb -= lsb
	return lsb.lsbIndex()
}
