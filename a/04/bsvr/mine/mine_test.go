package mine

import (
	"encoding/hex"
	"testing"

	"github.com/Univ-Wyo-Education/S19-4010/a/04/bsvr/block"
	"github.com/Univ-Wyo-Education/S19-4010/a/04/bsvr/hash"
)

func TestMineBlock(t *testing.T) {

	tests := []struct {
		bk               block.BlockType
		expectedSealHash string
		expectedNonce    uint64
	}{
		{
			bk: block.BlockType{
				Index:         0,
				Desc:          block.GenesisDesc,
				ThisBlockHash: []byte{},
				PrevBlockHash: []byte{},
				Nonce:         0,
				Seal:          []byte{},
			},
			expectedSealHash: "00008ff893670a368c3a3a4b7e1e9b06d2eb71cb9577990a8240edbe2837466b",
			expectedNonce:    79763,
		},
		{
			bk: block.BlockType{
				Index:         1,
				Desc:          "My First Block",
				ThisBlockHash: []byte{},
				PrevBlockHash: []byte{},
				Nonce:         0,
				Seal:          []byte{},
			},
			expectedSealHash: "0000b3228d0b2a6c862571b0e306e0d4ecadf092ec4b1f3cfd752777cebd09f8",
			expectedNonce:    169048,
		},
		{
			bk: block.BlockType{
				Index:         1,
				Desc:          "My First Block",
				ThisBlockHash: []byte{},
				PrevBlockHash: MustDecodeString("136c53391115ab7ff717bd24e62dd0df2c270500d7194290169a83488022548e"),
				Nonce:         0,
				Seal:          []byte{},
			},
			expectedSealHash: "00003ca4e7e2c4d08d97786f2259a4ffc8b65f099cf6d069105f50f33d6d6443",
			expectedNonce:    43233,
		},
	}

	for ii, test := range tests {
		test.bk.ThisBlockHash = hash.HashOf(block.SerializeBlock(&test.bk))
		tests[ii] = test
	}

	for ii, test := range tests {
		bk := &test.bk
		MineBlock(bk, "0000")
		SealString := hex.EncodeToString(bk.Seal)
		if SealString != test.expectedSealHash {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.expectedSealHash, SealString)
		}
		if bk.Nonce != test.expectedNonce {
			t.Errorf("Test %d, expected %d got %d\n", ii, test.expectedNonce, bk.Nonce)
		}
	}

}

func MustDecodeString(s string) []byte {
	rv, err := hex.DecodeString("136c53391115ab7ff717bd24e62dd0df2c270500d7194290169a83488022548e")
	if err != nil {
		panic(err)
	}
	return rv
}
