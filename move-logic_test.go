package chessongo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_GenerateLegalMoves(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		board := NewBoard()
		board.InitFromFen("r3k2r/pppbqppp/2n2n2/1B1pp3/1b1PP3/P1N1BN2/1PP1QPPP/R3K2R b KQkq - 0 8")
		for pb.Next() {
			board.GenerateLegalMoves()
		}
	})

}

func Test_GenerateLegalMoves(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1":               20,
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1":             20,
		"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2":           29,
		"rnbqkbnr/ppp2ppp/8/3pp3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 3":         28,
		"rnbqkbnr/ppp2ppp/8/3pp3/4P3/2N2N2/PPPP1PPP/R1BQKB1R b KQkq - 1 3":       38,
		"rnbqkb1r/ppp2ppp/5n2/3pp3/4P3/2N2N2/PPPP1PPP/R1BQKB1R w KQkq - 2 4":     30,
		"rnbqkb1r/ppp2ppp/5n2/3pp3/3PP3/2N2N2/PPP2PPP/R1BQKB1R b KQkq - 0 4":     36,
		"rnbqk2r/ppp2ppp/5n2/3pp3/1b1PP3/2N2N2/PPP2PPP/R1BQKB1R w KQkq - 1 5":    31,
		"r3k2r/pppnqppp/4bn2/3pp1B1/1b1PP3/2NB1N2/PPPQ1PPP/R3K2R w KQkq - 7 8":   43,
		"r3kq1r/pppb1ppp/2n2n2/1B1pp3/1b1PP3/P1NQBN2/1PP2PPP/R3K2R b KQkq - 0 9": 40,
		"r2qk2r/pppb1ppp/2n2n2/1B1pp3/1b1PP3/2N1BN2/PPP1QPPP/R3K2R b KQkq - 6 7": 39,
		"r3k2r/pppbqppp/2n2n2/1B1pp3/1b1PP3/P1NQBN2/1PP2PPP/R3K2R w - - 5 12":    38,
		"r3k2r/pppbqppp/2n2n2/1B1pp3/1b1PP3/P1N1BN2/1PP1QPPP/R3K2R b KQkq - 0 8": 41,
	}
	for fen, expectedCount := range positions {
		b := NewBoard()
		b.InitFromFen(fen)
		b.GenerateLegalMoves()
		//for _, move := range b.LegalMoves {
		//	fmt.Printf("%s \n", move.ToString())
		//}
		//fmt.Printf("\n\n")
		movesCount := len(b.LegalMoves)
		msg := fmt.Sprintf("Invalid moves count for FEN '%s', expected %d, got %d", fen, expectedCount, movesCount)
		assert.Equal(t, expectedCount, movesCount, msg)
	}
}

func Test_GenerateKing(t *testing.T) {
	var positions = map[string][]string{
		"8/8/8/8/8/8/8/7K w  - 0 0": []string{
			"h1 h2",
			"h1 g1",
			"h1 g2",
		},
		"8/8/8/8/8/8/8/6K1 w  - 0 0": []string{
			"g1 h1",
			"g1 h2",
			"g1 g2",
			"g1 f1",
			"g1 f2",
		},
		"K7/8/8/8/8/8/8/8 w  - 0 0": []string{
			"a8 b8",
			"a8 b7",
			"a8 a7",
		},
		"8/8/3K4/8/8/8/8/8 w  - 0 0": []string{
			"d6 e5",
			"d6 e6",
			"d6 e7",
			"d6 c5",
			"d6 c6",
			"d6 c7",
			"d6 d5",
			"d6 d7",
		},
		"8/8/3K4/8/8/4q3/8/8 w  - 0 0": []string{
			"d6 c6",
			"d6 c7",
			"d6 d5",
			"d6 d7",
		},
		"n7/8/3K4/n7/8/1b2q3/8/8 w  - 0 0": []string{
			"d6 d7",
		},
		"8/6n1/3K4/n7/q7/1b6/8/2r5 w  - 0 0": []string{
			"d6 e5",
			"d6 e7",
		},
		"7N/8/4k3/4n3/1B4K1/8/3R4/4R3 b  - 0 0": []string{
			"e6 f6",
		},
	}

	for fen, expectedMoves := range positions {
		b := NewBoard()
		//fmt.Println(fen)
		b.InitFromFen(fen)
		b.GenerateLegalMoves()
		//for _, move := range b.LegalMoves {
		//	fmt.Printf("%s \n", move.ToString())
		//}
		//fmt.Printf("\n\n")
		assert.Equal(t, len(expectedMoves), len(b.LegalMoves))
		for _, expectedMove := range expectedMoves {
			found := false
			for _, move := range b.LegalMoves {
				if move.ToString() == expectedMove {
					found = true
				}
			}
			msg := fmt.Sprintf("expected to find move %s for [%s]", expectedMove, fen)
			assert.Equal(t, found, true, msg)
		}
		//fmt.Printf("\n\n")

	}

}

/*
rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	h3 rnbqkbnr/pppppppp/8/8/8/7P/PPPPPPP1/RNBQKBNR b KQkq - 0 1
	h4 rnbqkbnr/pppppppp/8/8/7P/8/PPPPPPP1/RNBQKBNR b KQkq - 0 1
	g3 rnbqkbnr/pppppppp/8/8/8/6P1/PPPPPP1P/RNBQKBNR b KQkq - 0 1
	g4 rnbqkbnr/pppppppp/8/8/6P1/8/PPPPPP1P/RNBQKBNR b KQkq - 0 1
	f3 rnbqkbnr/pppppppp/8/8/8/5P2/PPPPP1PP/RNBQKBNR b KQkq - 0 1
	f4 rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq - 0 1
	e3 rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1
	e4 rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1
	d3 rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1
	d4 rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1
	c3 rnbqkbnr/pppppppp/8/8/8/2P5/PP1PPPPP/RNBQKBNR b KQkq - 0 1
	c4 rnbqkbnr/pppppppp/8/8/2P5/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1
	b3 rnbqkbnr/pppppppp/8/8/8/1P6/P1PPPPPP/RNBQKBNR b KQkq - 0 1
	b4 rnbqkbnr/pppppppp/8/8/1P6/8/P1PPPPPP/RNBQKBNR b KQkq - 0 1
	a3 rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR b KQkq - 0 1
	a4 rnbqkbnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQKBNR b KQkq - 0 1
	Na3 rnbqkbnr/pppppppp/8/8/8/N7/PPPPPPPP/R1BQKBNR b KQkq - 1 1
	Nc3 rnbqkbnr/pppppppp/8/8/8/2N5/PPPPPPPP/R1BQKBNR b KQkq - 1 1
	Nf3 rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1
	Nh3 rnbqkbnr/pppppppp/8/8/8/7N/PPPPPPPP/RNBQKB1R b KQkq - 1 1
*/
