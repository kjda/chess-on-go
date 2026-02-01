package chessongo

import (
	"encoding/binary"
	"errors"
)

// MarshalBinary encodes the board state into a byte slice.
func (g *Game) MarshalBinary() ([]byte, error) {
	// Fixed size: 64 (squares) + 1 (turn) + 1 (castling) + 1 (enpassant) + 4 (half) + 4 (full) + 4 (hist count) = 79 bytes
	// Variable size: 9 * historyCount
	buf := make([]byte, 79+len(g.PositionHistory)*9)

	copy(buf[0:64], fromPieces(g.Squares))
	buf[64] = uint8(g.Turn)
	buf[65] = uint8(g.Castling)
	buf[66] = uint8(g.EnPassant)
	binary.LittleEndian.PutUint32(buf[67:71], uint32(g.HalfMoves))
	binary.LittleEndian.PutUint32(buf[71:75], uint32(g.FullMoves))
	binary.LittleEndian.PutUint32(buf[75:79], uint32(len(g.PositionHistory)))

	offset := 79
	for hash, count := range g.PositionHistory {
		binary.LittleEndian.PutUint64(buf[offset:offset+8], hash)
		buf[offset+8] = uint8(count)
		offset += 9
	}

	return buf, nil
}

// UnmarshalBinary decodes the board state from a byte slice.
func (g *Game) UnmarshalBinary(data []byte) error {
	if len(data) < 79 {
		return errors.New("insufficient data for board")
	}

	g.Reset()

	for i := 0; i < 64; i++ {
		piece := Piece(data[i])
		if piece != EMPTY {
			g.addPiece(piece, i)
		}
	}

	g.Turn = Color(data[64])
	g.Castling = int(data[65])
	g.EnPassant = Square(data[66])
	g.HalfMoves = int(binary.LittleEndian.Uint32(data[67:71]))
	g.FullMoves = int(binary.LittleEndian.Uint32(data[71:75]))

	count := int(binary.LittleEndian.Uint32(data[75:79]))
	expectedSize := 79 + count*9
	if len(data) < expectedSize {
		return errors.New("insufficient data for position history")
	}

	g.PositionHistory = make(map[uint64]int, count)
	offset := 79
	for i := 0; i < count; i++ {
		hash := binary.LittleEndian.Uint64(data[offset : offset+8])
		c := int(data[offset+8])
		g.PositionHistory[hash] = c
		offset += 9
	}

	// Recompute Zobrist hash for current position
	g.ZobristHash = g.computeZobrist()

	// Update legal moves and check status
	g.GenerateLegalMoves()
	g.IsCheck = g.ComputeIsCheck()
	g.IsCheckmate = g.IsCheck && !g.hasMoves()
	g.IsStalement = !g.IsCheckmate && !g.hasMoves()
	g.IsMaterialDraw = g.hasInsufficientMaterial()
	g.IsThreefoldRepetition = g.checkThreefoldRepetition()
	g.IsFiftyMoveRule = g.checkFiftyMoveRule()
	g.IsSeventyFiveMoveRule = g.checkSeventyFiveMoveRule()
	g.IsFinished = g.IsCheckmate || g.IsStalement || g.IsMaterialDraw || g.IsFivefoldRepetition() || g.IsSeventyFiveMoveRule

	return nil
}

func fromPieces(p [64]Piece) []byte {
	b := make([]byte, 64)
	for i, v := range p {
		b[i] = uint8(v)
	}
	return b
}
