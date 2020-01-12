package board

import (
	"fmt"
	"math/rand"
	"strconv"
)

// StartingPosition 8x8 representation of normal chess starting position
var StartingPosition [8][8]string = [8][8]string{
	[8]string{"r", "n", "b", "q", "k", "b", "n", "r"},
	[8]string{"p", "p", "p", "p", "p", "p", "p", "p"},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{"P", "P", "P", "P", "P", "P", "P", "P"},
	[8]string{"R", "N", "B", "Q", "K", "B", "N", "R"}}

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

// GenerateChess960 Returns a 8x8 string array with a generated chess 960 position
func GenerateChess960() (position [8][8]string) {
	position = [8][8]string{
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"p", "p", "p", "p", "p", "p", "p", "p"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"P", "P", "P", "P", "P", "P", "P", "P"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "}}

	// -- Bishops --
	random1 := rand.Intn(8)
	position[0][random1] = "b"
	position[7][random1] = "B"

	// If first bishop is on an "even" square on rank1 then the second bishop
	// must be on an "odd" square in order to ensure 1 dark bishop and 1 light bishop
	random2 := rand.Intn(8)
	for random2%2 == random1%2 {
		random2 = rand.Intn(8)
	}
	position[0][random2] = "b"
	position[7][random2] = "B"

	// -- Queen --
	random3 := rand.Intn(8)
	// Find a place for the queen that is not already taken by the bishops
	for (random3 == random1) || (random3 == random2) {
		random3 = rand.Intn(8)
	}
	position[0][random3] = "q"
	position[7][random3] = "Q"

	// -- Knights --
	// Since we have placed already 3 pieces (2 bishops and a queen)
	// we are left with 5 possible squares. We take a random number "n"
	// between [1; 5] and find the "n"-th empty square and put the first knight there
	random4a := rand.Intn(5) + 1 // +1 makes the range [1; 5] instead of [0; 5)
	emptySquareCounter := 0
	var firstKnightIndex int // 8-based index to determine where the knight should be placed
	for idx, piece := range position[0] {
		if piece == " " {
			emptySquareCounter++
		}
		if emptySquareCounter == random4a {
			firstKnightIndex = idx
			break
		}
	}
	position[0][firstKnightIndex] = "n"
	position[7][firstKnightIndex] = "N"

	// The same process is applied for the second knight, however, there are
	// only 4 remaining empty squares
	random4b := rand.Intn(4) + 1 // +1 makes the range [1; 4] instead of [0; 4)
	emptySquareCounter = 0
	var secondKnightIndex int // 8-based index to determine where the knight should be placed
	for idx, piece := range position[0] {
		if piece == " " {
			emptySquareCounter++
		}
		if emptySquareCounter == random4b {
			secondKnightIndex = idx
			break
		}
	}
	position[0][secondKnightIndex] = "n"
	position[7][secondKnightIndex] = "N"

	// -- Rooks and King --
	// There are only 3 remaining empty squares.
	// Place the king in the middle one and the two rooks on the remaining squares
	for idx, piece := range position[0] {
		if piece == " " {
			position[0][idx] = "r"
			position[7][idx] = "R"
			break
		}
	}

	for idx, piece := range position[0] {
		if piece == " " {
			position[0][idx] = "k"
			position[7][idx] = "K"
			break
		}
	}

	for idx, piece := range position[0] {
		if piece == " " {
			position[0][idx] = "r"
			position[7][idx] = "R"
			break
		}
	}

	return position
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
