package chessongo

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const STARTING_POSITION_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const E_INVALID_FEN = "e:invalid:fen"

var RUNE_TO_PIECE = map[rune]Piece{
	'P': W_PAWN, 'N': W_KNIGHT, 'B': W_BISHOP, 'R': W_ROOK, 'Q': W_QUEEN, 'K': W_KING,
	'p': B_PAWN, 'n': B_KNIGHT, 'b': B_BISHOP, 'r': B_ROOK, 'q': B_QUEEN, 'k': B_KING,
}

var PIECE_TO_RUNE = map[Piece]rune{
	W_PAWN: 'P', W_KNIGHT: 'N', W_BISHOP: 'B', W_ROOK: 'R', W_QUEEN: 'Q', W_KING: 'K',
	B_PAWN: 'p', B_KNIGHT: 'n', B_BISHOP: 'b', B_ROOK: 'r', B_QUEEN: 'q', B_KING: 'k',
}

var STRING_TO_KIND = map[string]uint{
	"P": PAWN, "N": KNIGHT, "B": BISHOP, "R": ROOK, "Q": QUEEN, "K": KING,
	"p": PAWN, "n": KNIGHT, "b": BISHOP, "r": ROOK, "q": QUEEN, "k": KING,
}

var RUNE_TO_FILE = map[rune]int{'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7}

var RUNE_TO_RANK = map[rune]int{'1': 7, '2': 6, '3': 5, '4': 4, '5': 3, '6': 2, '7': 1, '8': 0}

var FILE_TO_STRING = map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e", 5: "f", 6: "g", 7: "h"}

var RANK_TO_STRING = map[int]string{0: "8", 1: "7", 2: "6", 3: "5", 4: "4", 5: "3", 6: "2", 7: "1"}

//initialize board from Fen string
func (b *Board) InitFromFen(fen string) error {
	b.Reset()
	b.Fen = fen
	fen = strings.Trim(fen, " ")
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return fmt.Errorf(E_INVALID_FEN)
	}
	pieces := parts[0]
	turn := parts[1]
	castling := parts[2]
	enPassant := parts[3]
	halfMoves := parts[4]
	fullMoves := parts[5]

	index := 0
	for _, r := range pieces {
		if index > 63 {
			return fmt.Errorf(E_INVALID_FEN)
		}
		switch r {
		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			b.addPiece(RUNE_TO_PIECE[r], index)
			index++
		case '1', '2', '3', '4', '5', '6', '7', '8':
			emptySquares, _ := strconv.Atoi(string(r))
			for i := 0; i < emptySquares; i++ {
				b.addPiece(EMPTY, index)
				index++
			}
		}
	}

	if turn == "w" {
		b.Turn = WHITE
	} else {
		b.Turn = BLACK
	}

	b.Castling = 0
	for _, r := range castling {
		switch r {
		case 'K':
			b.Castling |= CASTLE_WKS
		case 'Q':
			b.Castling |= CASTLE_WQS
		case 'k':
			b.Castling |= CASTLE_BKS
		case 'q':
			b.Castling |= CASTLE_BQS
		}
	}

	b.HalfMoves, _ = strconv.Atoi(string(halfMoves))
	b.FullMoves, _ = strconv.Atoi(string(fullMoves))

	if enPassant != "-" {
		enPassant = strings.ToLower(enPassant)
		file := RUNE_TO_FILE[rune(enPassant[0])]
		rank := RUNE_TO_RANK[rune(enPassant[1])]
		log.Println("SETTING EP TO")
		b.EnPassant = CoordsToSquare(rank, file)
	} else {
		b.EnPassant = 0
	}
	return nil
}

//return FEN representation of board
func (b *Board) ToFen() string {
	var pieces, turn, castling, enPassant string

	pieces = ""
	var i, emptyCount int = 0, 0
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			if b.Squares[i] != EMPTY {
				pieces += string(PIECE_TO_RUNE[b.Squares[i]])
				i++
				continue
			}
			for emptyCount = 0; file < 8 && b.Squares[i] == EMPTY; {
				emptyCount++
				i++
				file++
			}
			if emptyCount > 0 {
				pieces += strconv.Itoa(emptyCount)
				file--
			}
		}
		if rank < 7 {
			pieces += "/"
		}
	}

	if b.Turn == WHITE {
		turn = "w"
	} else {
		turn = "b"
	}

	castling = ""
	if (b.Castling & CASTLE_WKS) > 0 {
		castling += "K"
	}
	if (b.Castling & CASTLE_WQS) > 0 {
		castling += "Q"
	}
	if (b.Castling & CASTLE_BKS) > 0 {
		castling += "k"
	}
	if (b.Castling & CASTLE_BQS) > 0 {
		castling += "q"
	}
	if len(castling) == 0 {
		castling = "-"
	}

	if b.EnPassant == 0 {
		enPassant = "-"
	} else {
		rank, file := squareCoords(b.EnPassant)
		enPassant = FILE_TO_STRING[file] + RANK_TO_STRING[rank]
	}

	return fmt.Sprintf("%s %s %s %s %d %d", pieces, turn, castling, enPassant, b.HalfMoves, b.FullMoves)
}
