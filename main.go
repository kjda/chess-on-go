package main

import (
	"fmt"
	"time"
)

func main() {
	b := NewBoard()
	b.InitFromFen("k2k3r/P1ppqpb1/3pP3/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b - d7 4 11")
	PrintBitboard(b.Occupied, "BOARD")
	PrintBitboard(b.WhitePieces, "WHITES")
	PrintBitboard(b.BlackPieces, "BLACKS")
	PrintBitboard(b.Whites[KING], "WHITE  KING")
	PrintBitboard(b.Blacks[KING], "BLACK KING")
	PrintBitboard(Bitboard(0x1<<b.EnPassant), "ENPASSANR")
	var elapsed time.Duration
	start := time.Now()
	for i := 0; i < 1; i++ {
		
		b.GenerateLegalMoves()
		
		}
		elapsed = time.Since(start)
		for _, m := range b.PossibleMoves.Moves {
			fr, ff := squareCoords(m.from())
			tr, tf := squareCoords(m.to())
			var enPassant, castling, promotion, captured, promottionTo rune
			if m.isEnPassant() {
				enPassant = 'Y'
			} else {
				enPassant = 'N'
			}
			if m.isPromotionMove() {
				promotion = 'Y'
			} else {
				promotion = 'N'
			}
			if m.isCastlingMove() {
				castling = 'Y'
			} else {
				castling = 'N'
			}
			promottionTo = PIECE_TO_RUNE[Piece(b.Turn)|m.getPromotionTo()]
			captured = m.captured().toRune()
			//fmt.Printf("%.3b\n", m.getPromotionTo())
			piece := b.Squares[m.from()].toRune()
			fmt.Printf("%c: (%d, %d) => (%d, %d), Cap: %c, EP: %c, C: %c, Prom: %c PromTo: %c\n", piece, fr, ff, tr, tf, captured, enPassant, castling, promotion, promottionTo)
		}
	

	fmt.Printf("\n\nTIME_ %s calls %d\n\n", elapsed,genMovesCalls)

	b.PrintBoard("BOARD FINAL")
	if b.Turn == WHITE {
		fmt.Println("WHITE")
	} else {
		fmt.Println("BLACK")
	}
	if b.IsCheck() {
		fmt.Println("CHECK")
	} else {
		fmt.Println("NO_CHECK")
	}
	movesBB := Bitboard(0)
	for _, m := range b.PossibleMoves.Moves {
		movesBB |= 0x1 << m.to()
	}
	PrintBitboard(movesBB, "MOVES")
	fmt.Println("FEN: ", b.ToFen())
}
