package board

import (
	"fmt"
	"strconv"
)

// Indexes to access bitboards i.e. WP - white pawn, BB - black bishop
const (
	WP int = iota
	WN
	WB
	WR
	WQ
	WK
	BP
	BN
	BB
	BR
	BQ
	BK
)

// Board Struct to represent the chess board
type Board struct {
	// position  [8][8]string // todo this probably does not need to be into the struct
	bitboards [12]uint64
}

// ParseStringArray Parses a 8x8 string array that represents a
// board position in human form to 12 bitboards
func (board *Board) ParseStringArray(position [8][8]string) {
	for i := 0; i < 64; i++ {
		// MSB format.
		// Add 1 to the current square (i.e. value of i)
		// if there is a piece on that square the binaryStr is converted to binary
		// and is OR-ed with the current value of the corresponding piece bitboard
		// NOTE: string processing is started from left to right - i.e in MSB format
		binaryStr := "0000000000000000000000000000000000000000000000000000000000000000"
		binaryStr = binaryStr[i+1:] + "1" + binaryStr[:i]

		// i=0 is the A8 square on the chess board.
		// For loop iteration starts from the top left corner of the board below:
		// {"r","n","b","q","k","b","n","r"},
		// {"p","p","p","p","p","p","p","p"},
		// {" "," "," "," "," "," "," "," "},
		// {" "," "," "," "," "," "," "," "},
		// {" "," "," "," "," "," "," "," "},
		// {" "," "," "," "," "," "," "," "},
		// {"P","P","P","P","P","P","P","P"},
		// {"R","N","B","Q","K","B","N","R"},
		switch position[i/8][i%8] {
		case "P":
			board.bitboards[WP] += convertStringToBitboard(binaryStr)
		case "N":
			board.bitboards[WN] += convertStringToBitboard(binaryStr)
		case "B":
			board.bitboards[WB] += convertStringToBitboard(binaryStr)
		case "R":
			board.bitboards[WR] += convertStringToBitboard(binaryStr)
		case "Q":
			board.bitboards[WQ] += convertStringToBitboard(binaryStr)
		case "K":
			board.bitboards[WK] += convertStringToBitboard(binaryStr)
		case "p":
			board.bitboards[BP] += convertStringToBitboard(binaryStr)
		case "n":
			board.bitboards[BN] += convertStringToBitboard(binaryStr)
		case "b":
			board.bitboards[BB] += convertStringToBitboard(binaryStr)
		case "r":
			board.bitboards[BR] += convertStringToBitboard(binaryStr)
		case "q":
			board.bitboards[BQ] += convertStringToBitboard(binaryStr)
		case "k":
			board.bitboards[BK] += convertStringToBitboard(binaryStr)
		}

	}
}

// String Return string representing the current board (from the stored bitboards)
func (board *Board) String() string {
	var position [8][8]string

	for i := 0; i < 64; i++ {
		position[i/8][i%8] = " "
	}

	for i := 0; i < 64; i++ {
		if ((board.bitboards[WP] >> i) & 1) == 1 {
			position[i/8][i%8] = "P"
		}
		if ((board.bitboards[WN] >> i) & 1) == 1 {
			position[i/8][i%8] = "N"
		}
		if ((board.bitboards[WB] >> i) & 1) == 1 {
			position[i/8][i%8] = "B"
		}
		if ((board.bitboards[WR] >> i) & 1) == 1 {
			position[i/8][i%8] = "R"
		}
		if ((board.bitboards[WQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "Q"
		}
		if ((board.bitboards[WK] >> i) & 1) == 1 {
			position[i/8][i%8] = "K"
		}
		if ((board.bitboards[BP] >> i) & 1) == 1 {
			position[i/8][i%8] = "p"
		}
		if ((board.bitboards[BN] >> i) & 1) == 1 {
			position[i/8][i%8] = "n"
		}
		if ((board.bitboards[BB] >> i) & 1) == 1 {
			position[i/8][i%8] = "b"
		}
		if ((board.bitboards[BR] >> i) & 1) == 1 {
			position[i/8][i%8] = "r"
		}
		if ((board.bitboards[BQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "q"
		}
		if ((board.bitboards[BK] >> i) & 1) == 1 {
			position[i/8][i%8] = "k"
		}
	}

	var positionStr string
	for i := 0; i < 8; i++ {
		positionStr += fmt.Sprint(position[i])
		positionStr += "\n"
	}
	return positionStr
}

// convertStringToBitboard Helper for converting binary strings to bitboards
func convertStringToBitboard(binaryStr string) (bitboard uint64) {
	bitboard, err := strconv.ParseUint(binaryStr, 2, 64)
	if err != nil {
		panic(err)
	}
	return bitboard
}
