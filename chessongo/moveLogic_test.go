package chessongo

import (
	"testing"
)

//go test -test.bench=".*" -test.cpu="8"
func Benchmark_GenerateLegalMoves(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		board := NewBoard()
		board.InitFromFen("k2k3r/P1ppqpb1/3pP3/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b - d7 4 11")
		for pb.Next() {

			board.GenerateLegalMoves()
		}
	})

}
