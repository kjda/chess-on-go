package main

import (
	"fmt"
)

func main() {
	b := NewBoard()
	b.InitFromFen("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq a3 4 11")
	PrintBitboard(b.Occupied, "BOARD")
	PrintBitboard(b.WhitePieces, "WHITES")
	PrintBitboard(b.BlackPieces, "BLACKS")
	PrintBitboard(b.WhiteKing, "WHITE  KING")
	PrintBitboard(b.BlackKing, "BLACK KING")
	b.GenMoves()
	fmt.Println("FEN: ", b.ToFen())
}
