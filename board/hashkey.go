package board

import "math/rand"

// InitHashKeys initializes hashkeys for all pieces and possible positions, for castling rights, for side to move
func InitHashKeys() {
	for i := 0; i < 13; i++ {
		for j := 0; j < BoardSquareNum; j++ {
			PieceKeys[i][j] = rand.Uint64() // returns a random 64 bit number
		}
	}

	SideKey = rand.Uint64()
	for i := 0; i < 16; i++ {
		CastleKeys[i] = rand.Uint64()
	}
}

// PieceKeys hashkeys for each piece for each possible position for the key
var PieceKeys [13][BoardSquareNum]uint64

// SideKey the hashkey associated with the current side
var SideKey uint64

// CastleKeys haskeys associated with castling rights
var CastleKeys [16]uint64 // castling value ranges from 0-15 -> we need 16 hashkeys
