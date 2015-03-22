package main

import (
	"testing"
)

func Test_LSB(t *testing.T) {
	bb := Bitboard(0xa) //1010
	if bb.lsb() != 0x2 {
		t.Error("lsb failed")
	}
}

func Test_LSBIndex(t *testing.T) {
	bb := Bitboard(0x1 << 63)
	if bb.lsbIndex() != 63 {
		t.Error("lsbIndex failed:1")
	}
	bb = Bitboard(0xa)
	if bb.lsbIndex() != 1 {
		t.Error("lsbIndex failed:2")
	}
}

func Test_PopLsp(t *testing.T) {
	bb := Bitboard(0xa) //1010
	lsbInex := bb.popLSB()
	if lsbInex != 1 || bb != 0x8 {
		t.Error("popLSB failed")
	}
}

func Test_MSBIndex(t *testing.T) {
	bb := Bitboard(0xa) //1010
	if bb.msbIndex() != 3 {
		t.Error("msbIndex failed")
	}
}
