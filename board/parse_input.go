package board

import (
	"fmt"
	"strconv"
	"math/bits"
)

// PieceNotationMap maps piece notations (i.e. 'p', 'N') to piece values (i.e. 'BlackPawn', 'WhiteKnight')
var PieceNotationMap = map[string]int{
	"p": BP,
	"r": BR,
	"n": BN,
	"b": BB,
	"k": BK,
	"q": BQ,
	"P": WP,
	"R": WR,
	"N": WN,
	"B": WB,
	"K": WK,
	"Q": WQ,
}

// ParseFen parse fen position string and setup a position accordingly
func (board *Board) ParseFen(fen string) {
	piece := 0
	count := 0 // number of empty squares declared inside fen string

	board.Reset()
	char := 0

	for (count < 64) && char < len(fen) {
		switch t := string(fen[char]); t {
		case "p", "r", "n", "b", "k", "q", "P", "R", "N", "B", "K", "Q":
			// If we have a piece related char -> set the piece to corresponding value, i.e p -> BlackPawn
			piece = PieceNotationMap[t]
		case "1", "2", "3", "4", "5", "6", "7", "8":
			// otherwise it must be a count of a number of empty squares
			empty, _ := strconv.Atoi(t) // get number of empty squares and store in count
			count += empty
			char++
			continue
		case "/", " ":
			// if we have / or space then we are either at the end of the rank or at the end of the piece list
			// -> reset variables and continue the while loop
			char++
			continue
		default:
			panic("FEN error")
		}
		// compute piece color based on piece type
		color := Black
		if piece < BP {
			color = White
		}

		board.bitboards[piece] |= (1 << count)
		board.position[count] = piece
		board.material[color] += PieceValue[piece]
		board.positionKey ^= PieceKeys[piece][count]
		char++
		count++
	}

	newChar := ""
	char++ // move char from empty space to the w/b part of FEN
	// newChar should be set to the side to move part of the FEN string here
	newChar = string(fen[char])
	if newChar == "w" {
		board.Side = White
		// hash side (side key is only added for one side)
		board.positionKey ^= SideKey 
	} else if newChar == "b" {
		board.Side = Black
	} else {
		panic(fmt.Sprintf("Unknown side to move: %s", newChar))
	}

	// move char pointer 2 chars further and it should now point to the start of the castling permissions part of FEN
	char += 2

	// Iterate over the next 4 chars - they show if white is allowed to castle king or quenside and the same for black
	for i := 0; i < 4; i++ {
		newChar = string(fen[char])
		if newChar == " " {
			// when we hit a space, it means there are no more castling permissions => break
			break
		}
		switch newChar { // Depending on the char, enable the corresponding castling permission related bit
		case "K":
			board.castlePermissions |= WhiteKingCastling
		case "Q":
			board.castlePermissions |= WhiteQueenCastling
		case "k":
			board.castlePermissions |= BlackKingCastling
		case "q":
			board.castlePermissions |= BlackQueenCastling
		default:
			break
		}
		char++
	}
	// hash castle permissions
	board.positionKey ^= CastleKeys[board.castlePermissions]

	// AssertTrue(pos.castlePerm >= 0 && pos.castlePerm <= 15)
	// move to the en passant square related part of FEN
	char++
	newChar = string(fen[char])

	if newChar != "-" {
		file := newChar[0] - "a"[0]
		char++

		if file < 0 || file > 7 {
			panic(fmt.Sprintf("File out of bounds: file(%d)", file))
		}

		board.bitboards[EP] = FileMasks8[file]
		// hash en passant
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
	}
}
