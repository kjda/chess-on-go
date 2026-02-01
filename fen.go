package chessongo

import (
	"fmt"
	"strconv"
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

// initialize board from FEN string
func (b *Board) LoadFen(fen string) error {
	b.Reset()
	b.Fen = fen
	i := 0
	for i < len(fen) && fen[i] == ' ' {
		i++
	}
	if i >= len(fen) {
		return fmt.Errorf(E_INVALID_FEN)
	}

	// Pieces
	idx := 0
	for ; i < len(fen) && fen[i] != ' '; i++ {
		c := fen[i]
		switch c {
		case '/':
			continue
		case '1', '2', '3', '4', '5', '6', '7', '8':
			empty := int(c - '0')
			for j := 0; j < empty; j++ {
				if idx > 63 {
					return fmt.Errorf(E_INVALID_FEN)
				}
				b.addPiece(EMPTY, idx)
				idx++
			}
		case 'p':
			b.addPiece(B_PAWN, idx)
			idx++
		case 'n':
			b.addPiece(B_KNIGHT, idx)
			idx++
		case 'b':
			b.addPiece(B_BISHOP, idx)
			idx++
		case 'r':
			b.addPiece(B_ROOK, idx)
			idx++
		case 'q':
			b.addPiece(B_QUEEN, idx)
			idx++
		case 'k':
			b.addPiece(B_KING, idx)
			idx++
		case 'P':
			b.addPiece(W_PAWN, idx)
			idx++
		case 'N':
			b.addPiece(W_KNIGHT, idx)
			idx++
		case 'B':
			b.addPiece(W_BISHOP, idx)
			idx++
		case 'R':
			b.addPiece(W_ROOK, idx)
			idx++
		case 'Q':
			b.addPiece(W_QUEEN, idx)
			idx++
		case 'K':
			b.addPiece(W_KING, idx)
			idx++
		default:
			return fmt.Errorf(E_INVALID_FEN)
		}
	}
	if idx != 64 {
		return fmt.Errorf(E_INVALID_FEN)
	}
	if i >= len(fen) || fen[i] != ' ' {
		return fmt.Errorf(E_INVALID_FEN)
	}
	i++

	// Turn
	if i >= len(fen) {
		return fmt.Errorf(E_INVALID_FEN)
	}
	if fen[i] == 'w' {
		b.Turn = WHITE
	} else if fen[i] == 'b' {
		b.Turn = BLACK
	} else {
		return fmt.Errorf(E_INVALID_FEN)
	}
	i++
	if i >= len(fen) || fen[i] != ' ' {
		return fmt.Errorf(E_INVALID_FEN)
	}
	i++

	// Castling
	b.Castling = 0
	if i < len(fen) && fen[i] == '-' {
		i++
	} else {
		for ; i < len(fen) && fen[i] != ' '; i++ {
			switch fen[i] {
			case 'K':
				b.Castling |= CASTLE_WKS
			case 'Q':
				b.Castling |= CASTLE_WQS
			case 'k':
				b.Castling |= CASTLE_BKS
			case 'q':
				b.Castling |= CASTLE_BQS
			default:
				return fmt.Errorf(E_INVALID_FEN)
			}
		}
	}
	if i >= len(fen) || fen[i] != ' ' {
		return fmt.Errorf(E_INVALID_FEN)
	}
	i++

	// En passant
	if i >= len(fen) {
		return fmt.Errorf(E_INVALID_FEN)
	}
	if fen[i] == '-' {
		b.EnPassant = 0
		i++
	} else {
		if i+1 >= len(fen) {
			return fmt.Errorf(E_INVALID_FEN)
		}
		fileChar, rankChar := fen[i], fen[i+1]
		if fileChar < 'a' || fileChar > 'h' || rankChar < '1' || rankChar > '8' {
			return fmt.Errorf(E_INVALID_FEN)
		}
		file := int(fileChar - 'a')
		rank := 8 - int(rankChar-'0')
		b.EnPassant = CoordsToSquare(rank, file)
		i += 2
	}
	if i >= len(fen) || fen[i] != ' ' {
		return fmt.Errorf(E_INVALID_FEN)
	}
	i++

	// Half-move clock
	halfMoves := 0
	for ; i < len(fen) && fen[i] != ' '; i++ {
		d := fen[i]
		if d < '0' || d > '9' {
			return fmt.Errorf(E_INVALID_FEN)
		}
		halfMoves = halfMoves*10 + int(d-'0')
	}
	if i >= len(fen) {
		return fmt.Errorf(E_INVALID_FEN)
	}
	if fen[i] != ' ' {
		return fmt.Errorf(E_INVALID_FEN)
	}
	i++
	b.HalfMoves = halfMoves

	// Full-move number
	fullMoves := 0
	for ; i < len(fen) && fen[i] != ' '; i++ {
		d := fen[i]
		if d < '0' || d > '9' {
			return fmt.Errorf(E_INVALID_FEN)
		}
		fullMoves = fullMoves*10 + int(d-'0')
	}
	if i < len(fen) {
		// trailing non-space content means invalid
		for ; i < len(fen); i++ {
			if fen[i] != ' ' {
				return fmt.Errorf(E_INVALID_FEN)
			}
		}
	}
	b.FullMoves = fullMoves

	if b.PositionHistory == nil {
		b.PositionHistory = map[uint64]int{}
	}
	b.recordPosition()

	return nil
}

// return FEN representation of board
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
