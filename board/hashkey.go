package board

import (
	"math/bits"
	"math/rand"
)

// PieceKeys hashkeys for each piece for each possible position for the key
var PieceKeys [14][BoardSquareNum]uint64

// SideKey the hashkey associated with the current side
var SideKey uint64

// CastleKeys haskeys associated with castling rights
var CastleKeys [16]uint64 // castling value ranges from 0-15 -> we need 16 hashkeys

// KnightMoves an array of bitboards indicating every square the knight can go to from a given board index
var KnightMoves [BoardSquareNum]uint64

// KingMoves an array of bitboards indicating every square the king can go to from a given board index
var KingMoves [BoardSquareNum]uint64

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

	// Pregeneration of possible knight moves
	for i := 0; i < BoardSquareNum; i++ {
		var possibility uint64
		if i > 18 {
			possibility = KnightSpan << (i - 18)
		} else {
			possibility = KnightSpan >> (18 - i)
		}
		if i%8 < 4 {
			possibility &= (^FileGH)
		} else {
			possibility &= (^FileAB)
		}
		KnightMoves[i] = possibility
	}

	// Pregeneration of possible king moves
	for i := 0; i < BoardSquareNum; i++ {
		var possibility uint64

		if i > 9 {
			possibility = KingSpan << (i - 9)
		} else {
			possibility = KingSpan >> (9 - i)
		}

		if i%8 < 4 {
			possibility &= (^FileGH)
		} else {
			possibility &= (^FileAB)
		}

		KingMoves[i] = possibility
	}
}

// GeneratePositionKey takes a position and calculates a unique hashkey for it
func GeneratePositionKey(board *Board) (hashKey uint64) {
	for square, piece := range board.position {
		if piece != NoPiece {
			hashKey ^= PieceKeys[piece][square]
		}	
	}
	
	if board.bitboards[EP] != 0 {
		// take bitboard of en passant file, select only the first rank, find which file it is
		enPassantFile := bits.TrailingZeros64(board.bitboards[EP])
		hashKey ^= PieceKeys[EP][enPassantFile]
	}

	if board.Side == White {
		hashKey ^= SideKey
	}

	hashKey ^= CastleKeys[board.castlePermissions]

	return hashKey
}
