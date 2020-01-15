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

// KnightSpan Bitmask for selecting all knight moves
const KnightSpan uint64 = 43234889994

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

// RankMasks8 Array that holds bitmasks that select a given rank based on the index of
// the element i.e. index 0 selects Rank 1, 2- Rank 2 etc.
//! seems like index 0 is equal to rank8 from above
var RankMasks8 [8]uint64 = [8]uint64{
	0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_11111111,
	0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000,
	0b00000000_00000000_00000000_00000000_00000000_11111111_00000000_00000000,
	0b00000000_00000000_00000000_00000000_11111111_00000000_00000000_00000000,
	0b00000000_00000000_00000000_11111111_00000000_00000000_00000000_00000000,
	0b00000000_00000000_11111111_00000000_00000000_00000000_00000000_00000000,
	0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
	0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
}

var DiagonalMasks8 [15]uint64 = [15]uint64{
	0x1, 0x102, 0x10204, 0x1020408, 0x102040810, 0x10204081020, 0x1020408102040,
	0x102040810204080, 0x204081020408000, 0x408102040800000, 0x810204080000000,
	0x1020408000000000, 0x2040800000000000, 0x4080000000000000, 0x8000000000000000,
}

var AntiDiagonalMasks8 [15]uint64 = [15]uint64{
	0x80, 0x8040, 0x804020, 0x80402010, 0x8040201008, 0x804020100804, 0x80402010080402,
	0x8040201008040201, 0x4020100804020100, 0x2010080402010000, 0x1008040201000000,
	0x804020100000000, 0x402010000000000, 0x201000000000000, 0x100000000000000,
}

// TODO This should be in a struct
var NotWhitePieces uint64
var BlackPieces uint64
var Empty uint64
var Occupied uint64

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

	Occupied = (board.bitboards[WP] |
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

	Empty = ^Occupied


	moveList = board.PossiblePawnMovesWhite(history)
	moveList += board.PossibleWhiteKnightMoves()
	moveList += board.PossibleWhiteBishopMoves()
	moveList += board.PossibleWhiteRookMoves()
	moveList += board.PossibleWhiteQueenMoves()

	// fmt.Println(moveList)
	return moveList
}

func (board *Board) HorizontalAndVerticalMoves(square int) uint64 {
	var binarySquare uint64 = 1 << square
	fileMaskIdx := square % 8
	possibilitiesHorizontal := (Occupied - 2 * binarySquare) ^ bits.Reverse64(bits.Reverse64(Occupied) - 2 * bits.Reverse64(binarySquare))
	possibilitiesVertical := ((Occupied & FileMasks8[fileMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(
		bits.Reverse64(Occupied&FileMasks8[fileMaskIdx]) - (2 * bits.Reverse64(binarySquare)))
	return (possibilitiesHorizontal & RankMasks8[square / 8]) | (possibilitiesVertical & FileMasks8[fileMaskIdx])
}

func (board *Board) DiagonalAndAntiDiagonalMoves(square int) uint64 {
	var binarySquare uint64 = 1 << square
	diagonalMaskIdx := (square / 8) + (square % 8)
	antiDiagonalMaskIdx := (square/8) + 7 - (square % 8)
	possibilitiesDiagonal := ((Occupied & DiagonalMasks8[diagonalMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(bits.Reverse64(Occupied & DiagonalMasks8[diagonalMaskIdx]) - (2 * bits.Reverse64(binarySquare)))
	possibilitiesAntiDiagonal := ((Occupied & AntiDiagonalMasks8[antiDiagonalMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(bits.Reverse64(Occupied & AntiDiagonalMasks8[antiDiagonalMaskIdx]) - (2 * bits.Reverse64(binarySquare)))
	return (possibilitiesDiagonal & DiagonalMasks8[diagonalMaskIdx]) | (possibilitiesAntiDiagonal & AntiDiagonalMasks8[antiDiagonalMaskIdx])
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

	//    original                 WP >> 7          WP & ~FileA & ~Rank8    WP & BlackPieces (currently no black pieces on the 3rd rank to capture)
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [  X X X X X X X]      [  X X X X X X X]      [               ]
	// [X X X X X X X X]      [X              ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]

	// captures and moves formward forward: x1,y1,x2,y2
	pawnMoves = (wp >> 7) & (BlackPieces) & (^Rank8) & (^FileA) // capture right
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file     start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d", (index/8 + 1), (index%8 - 1), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 9) & (BlackPieces) & (^Rank8) & (^FileH) // capture left
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file     start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d", (index/8 + 1), (index%8 + 1), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 8) & Empty & (^Rank8) // move 1 square forward
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		//                                   final_rank     final_file start_rank   start_file
		moveList += fmt.Sprintf("%d%d%d%d", (index/8 + 1), (index%8), (index / 8), (index % 8))
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
		moveList += fmt.Sprintf("%d%d%d%d", (index/8 + 2), (index%8), (index / 8), (index % 8))
		pawnMoves &= ^possibility  // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves -1 ))  // find next bit equal to '1' i.e. next capture
	}

	// promotions - format: y1,y2,Promotion Type,"P"
	pawnMoves = (wp >> 7) & BlackPieces & Rank8 & (^FileA) // pawn promotion by capture right
	possibility = pawnMoves & (^(pawnMoves -1 ))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList += fmt.Sprintf("%d%dQP %d%dRP %d%dBP %d%dNP",
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
		moveList += fmt.Sprintf("%d%dQP %d%dRP %d%dBP %d%dNP",
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
		moveList += fmt.Sprintf("%d%dQP %d%dRP %d%dBP %d%dNP",
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
				//! This move is based from normal (white perspective) rank and file (starting from 1)
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

// todo refactor possibleWhiteBishop/Rook/Queen moves since they are the same
func (board *Board) PossibleWhiteBishopMoves() string {
	moveList := ""
	wb := board.bitboards[WB]
	// Choose bishop
	bishopPossibility := wb & (^(wb - 1))
	var possibility uint64

	for bishopPossibility != 0 {
		// Current bishop index (in bitmask)
		bishopIdx := bits.TrailingZeros64(bishopPossibility)
		possibility = board.DiagonalAndAntiDiagonalMoves(bishopIdx) & NotWhitePieces

		// choose move
		movePossibility := possibility& (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList += fmt.Sprintf("%d%d%d%d", (bishopIdx / 8), (bishopIdx % 8), (moveIndex / 8), (moveIndex % 8))
			possibility &= ^movePossibility  // remove move from all possible moves
			movePossibility = possibility& (^(possibility - 1))  // calculate new possible move
		}
		wb &= ^bishopPossibility
		bishopPossibility = wb & (^(wb - 1))
	}
	return moveList
}

func (board *Board) PossibleWhiteRookMoves() string {
	moveList := ""
	wr := board.bitboards[WR]
	// Choose rook
	rookPossibility := wr & (^(wr - 1))
	var possibility uint64

	for rookPossibility != 0 {
		// Current rook index (in bitmask)
		rookIdx := bits.TrailingZeros64(rookPossibility)
		possibility = board.HorizontalAndVerticalMoves(rookIdx) & NotWhitePieces

		// choose move
		movePossibility := possibility& (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList += fmt.Sprintf("%d%d%d%d", (rookIdx / 8), (rookIdx % 8), (moveIndex / 8), (moveIndex % 8))
			possibility &= ^movePossibility  // remove move from all possible moves
			movePossibility = possibility& (^(possibility - 1))  // calculate new possible move
		}
		wr &= ^rookPossibility
		rookPossibility = wr & (^(wr - 1))
	}
	return moveList
}

func (board *Board) PossibleWhiteQueenMoves() string {
	moveList := ""
	wq := board.bitboards[WQ]
	// Choose queen
	queenPossibility := wq & (^(wq - 1))
	var possibility uint64

	for queenPossibility != 0 {
		// Current queen index (in bitmask)
		queenIdx := bits.TrailingZeros64(queenPossibility)
		possibility = (board.HorizontalAndVerticalMoves(queenIdx) | board.DiagonalAndAntiDiagonalMoves(queenIdx)) & NotWhitePieces

		// choose move
		movePossibility := possibility& (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList += fmt.Sprintf("%d%d%d%d", (queenIdx / 8), (queenIdx % 8), (moveIndex / 8), (moveIndex % 8))
			possibility &= ^movePossibility  // remove move from all possible moves
			movePossibility = possibility& (^(possibility - 1))  // calculate new possible move
		}
		wq &= ^queenPossibility
		queenPossibility = wq & (^(wq - 1))
	}
	return moveList
}

func (board *Board) PossibleWhiteKnightMoves() string {
	moveList := ""
	wn := board.bitboards[WN]
	// Choose bishop
	knightPossibility := wn & (^(wn - 1))
	var possibility uint64

	for knightPossibility != 0 {
		// Current knight index (in bitmask)
		knightIdx := bits.TrailingZeros64(knightPossibility)

		// Move knight pattern mask around depending on the current knight position idx
		// KnightSpan is a predefined knight move pattern such as
		// [               ]
		// [               ]
		// [               ]
		// [        X   X  ]
		// [      X       X]
		// [          O    ]
		// [      X       X]
		// [        X   X  ]
		if knightIdx > 18 {
			possibility = KnightSpan << (knightIdx - 18)
		} else {
			possibility = KnightSpan >> (18 - knightIdx)
		}

		// handle wrap around of knight pattern movement
		if knightIdx % 8 < 4 {
			possibility &= (^FileGH) & NotWhitePieces
		} else {
			possibility &= (^FileAB) & NotWhitePieces
		}

		// choose move
		movePossibility := possibility& (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList += fmt.Sprintf("%d%d%d%d", (knightIdx / 8), (knightIdx % 8), (moveIndex / 8), (moveIndex % 8))
			possibility &= ^movePossibility  // remove move from all possible moves
			movePossibility = possibility& (^(possibility - 1))  // calculate new possible move
		}
		wn &= ^knightPossibility
		knightPossibility = wn & (^(wn - 1))
	}
	return moveList
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
