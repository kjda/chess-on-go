package main

import (
	"fmt"
	"kjda/chess/chessongo"
	"time"
)

func main() {
	b := chessongo.NewBoard()
	b.InitFromFen("k2k3r/P1ppqpb1/3pP3/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b - d7 4 11")
	chessongo.PrintBitboard(b.Occupied, "BOARD")
	chessongo.PrintBitboard(b.WhitePieces, "WHITES")
	chessongo.PrintBitboard(b.BlackPieces, "BLACKS")
	chessongo.PrintBitboard(b.Whites[chessongo.KING], "WHITE  KING")
	chessongo.PrintBitboard(b.Blacks[chessongo.KING], "BLACK KING")
	chessongo.PrintBitboard(chessongo.Bitboard(0x1<<b.EnPassant), "ENPASSANR")
	var elapsed time.Duration
	start := time.Now()
	b.GenerateLegalMoves()
	elapsed = time.Since(start)
	fmt.Printf("\nLegalMoves Generated in: %s \n", elapsed)

	fmt.Println("FEN: ", b.ToFen())
}
