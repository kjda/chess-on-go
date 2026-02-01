package chessongo

import (
	"encoding/binary"
	"errors"
)

// MarshalBinary encodes the board state into a byte slice.
func (b *Board) MarshalBinary() ([]byte, error) {
	// Fixed size: 64 (squares) + 1 (turn) + 1 (castling) + 1 (enpassant) + 4 (half) + 4 (full) + 4 (hist count) = 79 bytes
	// Variable size: 9 * historyCount
	buf := make([]byte, 79+len(b.PositionHistory)*9)

	copy(buf[0:64], fromPieces(b.Squares))
	buf[64] = uint8(b.Turn)
	buf[65] = uint8(b.Castling)
	buf[66] = uint8(b.EnPassant)
	binary.LittleEndian.PutUint32(buf[67:71], uint32(b.HalfMoves))
	binary.LittleEndian.PutUint32(buf[71:75], uint32(b.FullMoves))
	binary.LittleEndian.PutUint32(buf[75:79], uint32(len(b.PositionHistory)))

	offset := 79
	for hash, count := range b.PositionHistory {
		binary.LittleEndian.PutUint64(buf[offset:offset+8], hash)
		buf[offset+8] = uint8(count)
		offset += 9
	}

	return buf, nil
}

// UnmarshalBinary decodes the board state from a byte slice.
func (b *Board) UnmarshalBinary(data []byte) error {
	if len(data) < 79 {
		return errors.New("insufficient data for board")
	}

	b.Reset()

	for i := 0; i < 64; i++ {
		piece := Piece(data[i])
		if piece != EMPTY {
			b.addPiece(piece, i)
		}
	}

	b.Turn = Color(data[64])
	b.Castling = int(data[65])
	b.EnPassant = Square(data[66])
	b.HalfMoves = int(binary.LittleEndian.Uint32(data[67:71]))
	b.FullMoves = int(binary.LittleEndian.Uint32(data[71:75]))

	count := int(binary.LittleEndian.Uint32(data[75:79]))
	expectedSize := 79 + count*9
	if len(data) < expectedSize {
		return errors.New("insufficient data for position history")
	}

	b.PositionHistory = make(map[uint64]int, count)
	offset := 79
	for i := 0; i < count; i++ {
		hash := binary.LittleEndian.Uint64(data[offset : offset+8])
		c := int(data[offset+8])
		b.PositionHistory[hash] = c
		offset += 9
	}

	// Recompute Zobrist hash for current position
	b.ZobristHash = b.computeZobrist()

	// Update legal moves and check status
	b.GenerateLegalMoves()
	b.IsCheck = b.ComputeIsCheck()
	b.IsCheckmate = b.IsCheck && !b.hasMoves()
	b.IsStalement = !b.IsCheckmate && !b.hasMoves()
	b.IsMaterialDraw = b.hasInsufficientMaterial()
	b.IsThreefoldRepetition = b.checkThreefoldRepetition()
	b.IsFiftyMoveRule = b.checkFiftyMoveRule()
	b.IsSeventyFiveMoveRule = b.checkSeventyFiveMoveRule()
	b.IsFinished = b.IsCheckmate || b.IsStalement || b.IsMaterialDraw || b.IsFivefoldRepetition() || b.IsSeventyFiveMoveRule

	return nil
}

func fromPieces(p [64]Piece) []byte {
	b := make([]byte, 64)
	for i, v := range p {
		b[i] = uint8(v)
	}
	return b
}
