package network

import (
	"bytes"
	"testing"
)

func TestMainnetMagic(t *testing.T) {
	// Safecoin mainnet, src/chainparams.cpp
	var pchMessageStart [4]byte
	pchMessageStart[0] = 0xf1
	pchMessageStart[1] = 0xed
	pchMessageStart[2] = 0xe2
	pchMessageStart[3] = 0x8f

	magicBytes := Mainnet.Marshal(nil)

	if !bytes.Equal(magicBytes, pchMessageStart[:]) {
		t.Error("encoding failed")
	}

	magic, err := Decode(pchMessageStart[:])

	if err != nil || magic != Mainnet {
		t.Error("decoding failed")
	}
}
