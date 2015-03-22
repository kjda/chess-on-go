package main

import (
	"fmt"
)

func PrintBitboard(bb Bitboard, title string) {
	var shiftMask uint64 = 1
	//bb.printB("pp2")
	fmt.Printf("________________%s_______________________\n", title)
	for rank := 7; rank >= 0; rank-- {
		for file := 7; file >= 0; file-- {
			var squareIdx uint = uint(rank*8 + file)
			if bb&Bitboard(shiftMask<<squareIdx) > 0 {
				fmt.Printf(" X ")
			} else {
				fmt.Printf(" _ ")
			}
		}
		fmt.Println("")

	}
}

func (b *Board) PrintBoard(title string) {
	fmt.Printf("________________%s_______________________\n", title)
	for rank := 7; rank >= 0; rank-- {
		for file := 7; file >= 0; file-- {
			var squareIdx uint = uint(rank*8 + file)
			if b.Squares[squareIdx] == EMPTY {
				fmt.Printf(" -- ")
			} else {
				color := 'b'
				if b.Squares[squareIdx].color() == WHITE {
					color = 'w'
				}
				fmt.Printf(" %c%c ", color, PIECE_TO_RUNE[b.Squares[squareIdx]])
			}
		}
		fmt.Println("")

	}
}

func log(msg string) {
	fmt.Println(msg)
}
