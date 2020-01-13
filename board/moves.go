package board

import (
	"fmt"
	"math/bits"
)

// The following bitmasks represent squares starting from H1-A1 -> H8-A8. Ranks are separated by "_"

// FileA Bitmask for selecting all squares that are on the A file
const FileA uint64 = 0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001

// FileH Bitmask for selecting all squares that are on the H file
const FileH uint64 = 0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000

// FileAB Bitmask for selecting all squares that are on the A & B files
const FileAB uint64 = 0b00000011_00000011_00000011_00000011_00000011_00000011_00000011_00000011

// FileGH Bitmask for selecting all squares that are on the G & H files
const FileGH uint64 = 0b11000000_11000000_11000000_11000000_11000000_11000000_11000000_11000000

// Rank8 Bitmask for selecting all squares that are on the 8th rank
const Rank8 uint64 = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_11111111

// Rank5 Bitmask for selecting all squares that are on the 5th rank
const Rank5 uint64 = 0b00000000_00000000_00000000_00000000_11111111_00000000_00000000_00000000

// Rank4 Bitmask for selecting all squares that are on the 4th rank
const Rank4 uint64 = 0b00000000_00000000_00000000_11111111_00000000_00000000_00000000_00000000

// Rank1 Bitmask for selecting all squares that are on the 1st rank
const Rank1 uint64 = 0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000

// Center Bitmask for selecting all center squares (D4, E4, D5, E5)
const Center uint64 = 0b00000000_00000000_00000000_00011000_00011000_00000000_00000000_00000000

// ExtendedCenter Bitmask for selecting all extended center squares (C3 - F3, C4 - F4, C5 - F5)
const ExtendedCenter uint64 = 0b00000000_00000000_00111100_00111100_00111100_00111100_00000000_00000000

// QueenSide Bitmask for selecting all queenside squares
const QueenSide uint64 = 0b00001111_00001111_00001111_00001111_00001111_00001111_00001111_00001111

// KingSide Bitmask for selecting all kingside squares
const KingSide uint64 = 0b11110000_11110000_11110000_11110000_11110000_11110000_11110000_11110000

// FileMasks8 Array that holds bitmasks that select a given file based on the index of
// the element i.e. index 0 selects File A, 1- FileB etc.
var FileMasks8 [8]uint64 = [8]uint64{
	0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001,
	0b00000010_00000010_00000010_00000010_00000010_00000010_00000010_00000010,
	0b00000100_00000100_00000100_00000100_00000100_00000100_00000100_00000100,
	0b00001000_00001000_00001000_00001000_00001000_00001000_00001000_00001000,
	0b00010000_00010000_00010000_00010000_00010000_00010000_00010000_00010000,
	0b00100000_00100000_00100000_00100000_00100000_00100000_00100000_00100000,
	0b01000000_01000000_01000000_01000000_01000000_01000000_01000000_01000000,
	0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000,
}

// TODO This should be in a struct
var NotWhitePieces uint64
var BlackPieces uint64
var Empty uint64

// DrawBitboard Prints a given bitboard to stdout in a human readable way
func DrawBitboard(bitboard uint64) {
	var bitboardStr [8][8]string

	for i := 0; i < 64; i++ {
		if ((bitboard >> i) & 1) == 1 {
			bitboardStr[i/8][i%8] = "X"
		}

		if bitboardStr[i/8][i%8] == "" {
			// replace empty string with a space for better readibility of bitboard
			bitboardStr[i/8][i%8] = " "
		}
	}

	for _, rank := range bitboardStr {
		fmt.Printf("%s\n", rank)
	}
	fmt.Println()
}

func (board *Board) PossibleMovesWhite(history string) (moveList string) {
	// This represents all squares which are not white pieces (including empty squares).
	// Black king is added in order to avoid generating capture moves on the black king.
	// For example pawn takes king is not a legal move
	NotWhitePieces = ^(board.bitboards[WP] |
		board.bitboards[WN] |
		board.bitboards[WB] |
		board.bitboards[WR] |
		board.bitboards[WQ] |
		board.bitboards[WK] |
		board.bitboards[BK])

	BlackPieces = (board.bitboards[BP] |
		board.bitboards[BN] |
		board.bitboards[BB] |
		board.bitboards[BR] |
		board.bitboards[BQ])

	Empty = ^(board.bitboards[WP] |
		board.bitboards[WN] |
		board.bitboards[WB] |
		board.bitboards[WR] |
		board.bitboards[WQ] |
		board.bitboards[WK] |
		board.bitboards[BP] |
		board.bitboards[BN] |
		board.bitboards[BB] |
		board.bitboards[BR] |
		board.bitboards[BQ] |
		board.bitboards[BK])

	DrawBitboard(NotWhitePieces)
	fmt.Println()
	DrawBitboard(BlackPieces)
	fmt.Println()
	DrawBitboard(Empty)
	fmt.Println()

	moveList = board.PossiblePawnMovesWhite(history)
	fmt.Println(moveList)
	return moveList
}

func (board *Board) PossiblePawnMovesWhite(history string) (moveList string) {
	// todo movelist as string looks pretty stupid ???

	wp := board.bitboards[WP]
	var possibility uint64  // holds one potential capture at a time
	var index int  // index of the "possibility" capture
	var pawnMoves uint64

	// Move all bits from the white pawn bitboard to the left by 7 squares
	// and disable File A (we do not want any leftovers from the move to the left)
	// also remove any captures on the 8th rank since that will be handled by promotions
	// finally AND the resulting bitboard with the black pieces i.e. a capture move is any
	// move that can capture an enemy piece
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [X X X X X X X X]
	// [               ]
	//    WP >> 7
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [  X X X X X X X]
	// [X              ]
	// [               ]
	//    WP & ~FileA & ~Rank8
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [  X X X X X X X]
	// [               ]
	// [               ]
	//    WP & BlackPieces (currently no black pieces on the 3rd rank to capture)
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// [               ]
	// captures and moves formward forward: x1,y1,x2,y2
	pawnMoves = (wp >> 7) & (BlackPieces) & (^Rank8) & (^FileA) // capture right
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file     start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d ", (index/8 + 1), (index%8 - 1), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 9) & (BlackPieces) & (^Rank8) & (^FileH) // capture left
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file     start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d ", (index/8 + 1), (index%8 + 1), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 8) & Empty & (^Rank8) // move 1 square forward
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d ", (index/8 + 1), (index%8), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (wp >> 16) & Empty & (Empty >> 8) & Rank4 // move 2 squares forward
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d ", (index/8 + 2), (index%8), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	// promotions - format: y1,y2,Promotion Type,"P"
	pawnMoves = (wp >> 7) & BlackPieces & Rank8 & (^FileA) // pawn promotion by capture right
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList += fmt.Sprintf("%d%dQP %d%dRP %d%dBP %d%dNP ",
								(index/8 - 1), (index % 8),
								(index/8 - 1), (index % 8),
								(index/8 - 1), (index % 8),
								(index/8 - 1), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 9) & BlackPieces & Rank8 & (^FileH) // pawn promotion by capture left
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList += fmt.Sprintf("%d%dQP %d%dRP %d%dBP %d%dNP ",
								(index/8 + 1), (index % 8),
								(index/8 + 1), (index % 8),
								(index/8 + 1), (index % 8),
								(index/8 + 1), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 8) & Empty & Rank8 // pawn promotion by move 1 forward
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList += fmt.Sprintf("%d%dQP %d%dRP %d%dBP %d%dNP ",
								(index/8), (index % 8),
								(index/8), (index % 8),
								(index/8), (index % 8),
								(index/8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	if histLen := len(history); histLen >= 4 {
		// if history move, moved 2 squares forward and inside the same file
		if (history[histLen - 1] == history[histLen - 3]) &&
			(abs(int(history[histLen-2] - history[histLen-4])) == 2) {
			enPassFile := history[histLen -1] - "0"[0]  // find enpass file (type byte)
			// en passant right
			// shows piece to remove, not the destination
			possibility = (wp << 1) & board.bitboards[BP] & Rank5 & (^FileA) & FileMasks8[enPassFile]
			if possibility != 0 {
				index = bits.TrailingZeros64(possibility)
				moveList += fmt.Sprintf("%d%d E", (index%8 -1), (index%8))
			}

			// en passant left
			// shows piece to remove, not the destination
			possibility = (wp >> 1) & board.bitboards[BP] & Rank5 & (^FileH) & FileMasks8[enPassFile]
			if possibility != 0 {
				index = bits.TrailingZeros64(possibility)
				moveList += fmt.Sprintf("%d%d E", (index%8 +1), (index%8))
			}
		}
	}

	return moveList
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}