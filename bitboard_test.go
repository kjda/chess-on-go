package chessongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLSB(t *testing.T) {
	tests := []struct {
		name   string
		value  Bitboard
		expect Bitboard
	}{
		{"singleBit", 0x8, 0x8},
		{"multipleBits", 0xa, 0x2}, // 1010 -> lowest bit at 2
		{"highBit", Bitboard(1) << 63, Bitboard(1) << 63},
	}

	for _, tt := range tests {
		require.Equalf(t, tt.expect, tt.value.lsb(), "lsb %s", tt.name)
	}
}

func TestLSBIndex(t *testing.T) {
	tests := []struct {
		name   string
		value  Bitboard
		expect uint
	}{
		{"highest", Bitboard(1) << 63, 63},
		{"middle", Bitboard(1) << 27, 27},
		{"lowMultiple", 0xa, 1},
	}

	for _, tt := range tests {
		require.Equalf(t, tt.expect, tt.value.lsbIndex(), "lsbIndex %s", tt.name)
	}
}

func TestPopLSB(t *testing.T) {
	bb := Bitboard(0b10110)
	first := bb.popLSB()
	second := bb.popLSB()

	require.Equal(t, uint(1), first)
	require.Equal(t, uint(2), second)
	require.Equal(t, Bitboard(0b10000), bb)
}

func TestMSBIndex(t *testing.T) {
	tests := []struct {
		name   string
		value  Bitboard
		expect int
	}{
		{"zero", 0, -1},
		{"low", 0b1, 0},
		{"middle", 0b1001, 3},
		{"high", Bitboard(1) << 62, 62},
	}

	for _, tt := range tests {
		require.Equalf(t, tt.expect, tt.value.msbIndex(), "msbIndex %s", tt.name)
	}
}

func TestNumberOfSetBits(t *testing.T) {
	require.Equal(t, 0, Bitboard(0).NumberOfSetBits())
	require.Equal(t, 1, Bitboard(1).NumberOfSetBits())
	require.Equal(t, 5, Bitboard(0b10110101).NumberOfSetBits())
}
