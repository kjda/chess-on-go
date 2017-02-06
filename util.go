package chessongo

import (
	"fmt"
)

func PrintBitboard(bb Bitboard, title string) {
	var shiftMask uint64 = 1
	//bb.printB("pp2")
	fmt.Printf("________________%s_______________________\n", title)
	for r := 7; r >= 0; r-- {
		for f := 7; f >= 0; f-- {
			var squareIdx uint = uint(r*8 + f)
			if bb&Bitboard(shiftMask<<squareIdx) > 0 {
				fmt.Printf("%c%d: X ", file(Square(squareIdx)), rank(Square(squareIdx)))
			} else {
				fmt.Printf("%c%d: _ ", file(Square(squareIdx)), rank(Square(squareIdx)))
			}
		}
		fmt.Println("")

	}
}

func (b *Board) PrintBoard(title string) {
	fmt.Printf("________________%s_______________________\n", title)
	for i, v := range b.Squares {
		if i%8 == 0 {
			fmt.Println("")
		}
		if v.ToRune() == ' ' {
			fmt.Printf("   %c%d:-   ", file(Square(i)), rank(Square(i)))
		} else {
			fmt.Printf("   %c%d:%c   ", file(Square(i)), rank(Square(i)), v.ToRune())
		}
	}
	fmt.Printf("\n_______________________________________\n")
}
