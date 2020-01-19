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

// KingSpan Bitmask for selecting all king moves
const KingSpan uint64 = 460039

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
var NotMyPieces uint64
var MyPieces uint64
var EnemyPieces uint64
var Empty uint64
var Occupied uint64
var CastleRooks [4]int = [4]int{63, 56, 7, 0}

// Move Move type
type Move struct {
	Move  int // todo consider move to be uint32 instead of int
	score int
}

/* Game move - information stored in the move int from type Move
               |Ca| |--To-||-From-|
0000 0000 0000 0000 0000 0011 1111 -> From - 0x3F
0000 0000 0000 0000 1111 1100 0000 -> To - >> 6, 0x3F
0000 0000 0000 1111 0000 0000 0000 -> Captured - >> 12, 0xF
0000 0000 0001 0000 0000 0000 0000 -> En passant capt - >> 16 - 0x40000
0000 0000 0010 0000 0000 0000 0000 -> PawnStart - >> 17 - 0x80000
0000 0011 1100 0000 0000 0000 0000 -> Promotion to what piece - >> 18, 0xF
0000 0100 0000 0000 0000 0000 0000 -> Castle - >> 22 0x1000000
*/

// FromSq - macro that returns the 'from' bits from the move int
func FromSq(m int) int {
	return m & 0x3f
}

// ToSq - macro that returns the 'to' bits from the move int
func ToSq(m int) int {
	return (m >> 6) & 0x3f
}

// Captured - macro that returns the 'Captured' bits from the move int
func Captured(m int) int {
	return (m >> 12) & 0xf
}

// Promoted - macro that returns the 'Promoted' bits from the move int
func Promoted(m int) int {
	return (m >> 18) & 0xf
}

// PawnStartFlag - macro that returns the 'PawnStart' flag bits from the move int
func PawnStartFlag(m int) int {
	return (m >> 17) & 1
}

// EnPassantFlag - macro that returns the 'EnPassant' capture flag bits from the move int
func EnPassantFlag(m int) int {
	return (m >> 16) & 1
}

// CastleFlag - macro that returns the 'CastleFlag' flag bits from the move int
func CastleFlag(m int) int {
	return (m >> 22) & 1
}

// GetMoveInt creates and returns a move int from given move information
func GetMoveInt(fromSq, toSq, capturePiece, promotionPiece, flag int) int {
	return fromSq | (toSq << 6) | (capturePiece << 12) | (promotionPiece << 18) | flag
}

const (
	// MoveFlagEnPass move flag that denotes if the capture was an enpass
	MoveFlagEnPass int = 0x10000

	// MoveFlagPawnStart move flag that denotes if move was pawn start (2x)
	MoveFlagPawnStart int = 0x20000

	// MoveFlagCapture move flag that denotes if move was capture without saying what the capture was (checks capture & enpas squares)
	MoveFlagCapture int = 0xF000

	// NoFlag constant that denotes no flag is applied to move
	NoFlag int = 0

	// MoveFlagCastle move flag that denotes if move was castling
	MoveFlagCastle int = 0x400000
)

// MaxPositionMoves maximum number of possible moves for a given position
const MaxPositionMoves int = 256

// MoveList Struct to hold all generated moves for a given position
type MoveList struct {
	Moves [MaxPositionMoves]Move
	Count int // number of moves on the moves list
}

// AddMove Adds move to move list and updates count
func (moveList *MoveList) AddMove(move int) {
	moveList.Moves[moveList.Count].Move = move
	moveList.Count++
}

// PrintSquare get algebraic notation of square i.e. b2, a6 from array index
func GetSquareString(sq int) string {
	file := sq % 8
	rank := 8 - (sq / 8) - 1

	// "a"[0] -> returns the byte value of the char 'a' -> convert to int to get ascii value
	// then add the file/rank value to it and convert back to string
	// therefore this automatically translates the files from 0-7 to a-h
	fileStr := string(int("a"[0]) + file)
	rankStr := string(int("1"[0]) + rank)

	squareStr := fileStr + rankStr
	return squareStr
}

// PrintMove prints move in algebraic notation
func PrintMove(move int) string {
	// fmt.Printf("FromSq: %d, ToSq: %d, Promoted: %d\n", FromSq(move), ToSq(move), Promoted(move))

	fromSq := GetSquareString(FromSq(move))
	toSq := GetSquareString(ToSq(move))

	moveStr := fromSq + toSq

	// if this move is a promotion, add char of the piece we promote to at the end of the move string
	// i.e. if a7a8q -> we promote to Queen
	pieceChar := ""
	switch promoted := Promoted(move); promoted {
	case WN, BN:
		pieceChar = "n"
	case WB, BB:
		pieceChar = "b"
	case WR, BR:
		pieceChar = "r"
	case WQ, BQ:
		pieceChar = "q"
	}
	moveStr += pieceChar

	return moveStr
}

// PrintMoveList prints move list
func PrintMoveList(moveList *MoveList) {
	fmt.Println("MoveList:\n", moveList.Count)

	for index := 0; index < moveList.Count; index++ {

		move := moveList.Moves[index].Move
		score := moveList.Moves[index].score

		fmt.Printf("Move:%d > %s (score:%d)\n", index+1, PrintMove(move), score)
	}
	fmt.Printf("MoveList Total %d Moves:\n\n", moveList.Count)
}

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

func (board *Board) PossibleMovesWhite(moveList *MoveList) {
	// This represents all squares which are not white pieces (including empty squares).
	// Black king is added in order to avoid generating capture moves on the black king.
	// For example pawn takes king is not a legal move
	NotMyPieces = ^(board.bitboards[WP] |
		board.bitboards[WN] |
		board.bitboards[WB] |
		board.bitboards[WR] |
		board.bitboards[WQ] |
		board.bitboards[WK] |
		board.bitboards[BK])

	EnemyPieces = (board.bitboards[BP] |
		board.bitboards[BN] |
		board.bitboards[BB] |
		board.bitboards[BR] |
		board.bitboards[BQ])

	MyPieces = (board.bitboards[WP] |
		board.bitboards[WN] |
		board.bitboards[WB] |
		board.bitboards[WR] |
		board.bitboards[WQ])

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

	board.possibleWhitePawn(moveList)
	board.possibleKnightMoves(moveList, board.bitboards[WN])
	board.possibleBishopMoves(moveList, board.bitboards[WB])
	board.possibleRookMoves(moveList, board.bitboards[WR])
	board.possibleQueenMoves(moveList, board.bitboards[WQ])
	board.possibleKingMoves(moveList, board.bitboards[WK])
	board.possibleCastleWhite(moveList)
}

func (board *Board) PossibleMovesBlack(moveList *MoveList) {
	// This represents all squares which are not white pieces (including empty squares).
	// Black king is added in order to avoid generating capture moves on the black king.
	// For example pawn takes king is not a legal move
	NotMyPieces = ^(board.bitboards[BP] |
		board.bitboards[BN] |
		board.bitboards[BB] |
		board.bitboards[BR] |
		board.bitboards[BQ] |
		board.bitboards[BK] |
		board.bitboards[WK])

	EnemyPieces = (board.bitboards[WP] |
		board.bitboards[WN] |
		board.bitboards[WB] |
		board.bitboards[WR] |
		board.bitboards[WQ])

	MyPieces = (board.bitboards[BP] |
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

	// todo the following could be split into goroutines
	board.possibleBlackPawn(moveList)
	board.possibleKnightMoves(moveList, board.bitboards[BN])
	board.possibleBishopMoves(moveList, board.bitboards[BB])
	board.possibleRookMoves(moveList, board.bitboards[BR])
	board.possibleQueenMoves(moveList, board.bitboards[BQ])
	board.possibleKingMoves(moveList, board.bitboards[BK])
	board.possibleCastleBlack(moveList)
}

func (board *Board) HorizontalAndVerticalMoves(square int) uint64 {
	//! Requires Occupied to be up to date
	var binarySquare uint64 = 1 << square
	fileMaskIdx := square % 8
	possibilitiesHorizontal := (Occupied - 2*binarySquare) ^ bits.Reverse64(bits.Reverse64(Occupied)-2*bits.Reverse64(binarySquare))
	possibilitiesVertical := ((Occupied & FileMasks8[fileMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(
		bits.Reverse64(Occupied&FileMasks8[fileMaskIdx])-(2*bits.Reverse64(binarySquare)))
	return (possibilitiesHorizontal & RankMasks8[square/8]) | (possibilitiesVertical & FileMasks8[fileMaskIdx])
}

func (board *Board) DiagonalAndAntiDiagonalMoves(square int) uint64 {
	//! Requires Occupied to be up to date
	var binarySquare uint64 = 1 << square
	diagonalMaskIdx := (square / 8) + (square % 8)
	antiDiagonalMaskIdx := (square / 8) + 7 - (square % 8)
	possibilitiesDiagonal := ((Occupied & DiagonalMasks8[diagonalMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(bits.Reverse64(Occupied&DiagonalMasks8[diagonalMaskIdx])-(2*bits.Reverse64(binarySquare)))
	possibilitiesAntiDiagonal := ((Occupied & AntiDiagonalMasks8[antiDiagonalMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(bits.Reverse64(Occupied&AntiDiagonalMasks8[antiDiagonalMaskIdx])-(2*bits.Reverse64(binarySquare)))
	return (possibilitiesDiagonal & DiagonalMasks8[diagonalMaskIdx]) | (possibilitiesAntiDiagonal & AntiDiagonalMasks8[antiDiagonalMaskIdx])
}

func (board *Board) possibleWhitePawn(moveList *MoveList) {
	// todo movelist as string looks pretty stupid ???

	wp := board.bitboards[WP]
	var possibility uint64 // holds one potential capture at a time
	var index int          // index of the "possibility" capture
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

	// captures and moves formward forward
	pawnMoves = (wp >> 7) & (EnemyPieces) & (^Rank8) & (^FileA) // capture right
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+7, index, 0, 0, MoveFlagCapture))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 9) & (EnemyPieces) & (^Rank8) & (^FileH) // capture left
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+9, index, 0, 0, MoveFlagCapture))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 8) & Empty & (^Rank8) // move 1 square forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+8, index, 0, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (wp >> 16) & Empty & (Empty >> 8) & Rank4 // move 2 squares forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+16, index, 0, 0, MoveFlagPawnStart))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	pawnMoves = (wp >> 7) & EnemyPieces & Rank8 & (^FileA) // pawn promotion by capture right
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index+7, index, 0, WQ, NoFlag))
		moveList.AddMove(GetMoveInt(index+7, index, 0, WR, NoFlag))
		moveList.AddMove(GetMoveInt(index+7, index, 0, WB, NoFlag))
		moveList.AddMove(GetMoveInt(index+7, index, 0, WN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 9) & EnemyPieces & Rank8 & (^FileH) // pawn promotion by capture left
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index+9, index, 0, WQ, NoFlag))
		moveList.AddMove(GetMoveInt(index+9, index, 0, WR, NoFlag))
		moveList.AddMove(GetMoveInt(index+9, index, 0, WB, NoFlag))
		moveList.AddMove(GetMoveInt(index+9, index, 0, WN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 8) & Empty & Rank8 // pawn promotion by move 1 forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+8, index, 0, WQ, NoFlag))
		moveList.AddMove(GetMoveInt(index+8, index, 0, WR, NoFlag))
		moveList.AddMove(GetMoveInt(index+8, index, 0, WB, NoFlag))
		moveList.AddMove(GetMoveInt(index+8, index, 0, WN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// En passant right
	possibility = (wp << 1) & board.bitboards[BP] & Rank5 & (^FileA) & board.bitboards[EP]
	if possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-1, index-8, 0, 0, MoveFlagEnPass))
	}
	// en passant left
	possibility = (wp >> 1) & board.bitboards[BP] & Rank5 & (^FileH) & board.bitboards[EP]
	if possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+1, index-8, 0, 0, MoveFlagEnPass))
	}
}

func (board *Board) possibleBlackPawn(moveList *MoveList) {
	bp := board.bitboards[BP]
	var possibility uint64 // holds one potential capture at a time
	var index int          // index of the "possibility" capture
	var pawnMoves uint64

	// captures and moves formward forward:
	pawnMoves = (bp << 7) & (EnemyPieces) & (^Rank1) & (^FileH) // capture right
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-7, index, 0, 0, MoveFlagCapture))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (bp << 9) & (EnemyPieces) & (^Rank1) & (^FileA) // capture left
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-9, index, 0, 0, MoveFlagCapture))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (bp << 8) & Empty & (^Rank1) // move 1 square forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-8, index, 0, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (bp << 16) & Empty & (Empty << 8) & Rank5 // move 2 squares forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-16, index, 0, 0, MoveFlagPawnStart))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	pawnMoves = (bp << 7) & EnemyPieces & Rank1 & (^FileH) // pawn promotion by capture right
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index-7, index, 0, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-7, index, 0, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-7, index, 0, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-7, index, 0, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (bp << 9) & EnemyPieces & Rank1 & (^FileA) // pawn promotion by capture left
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index-9, index, 0, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, 0, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, 0, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, 0, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (bp << 8) & Empty & Rank1 // pawn promotion by move 1 forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-8, index, 0, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-8, index, 0, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-8, index, 0, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-8, index, 0, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// En passant right
	possibility = (bp >> 1) & board.bitboards[WP] & Rank4 & (^FileH) & board.bitboards[EP]
	if possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+1, index+8, 0, 0, MoveFlagEnPass))
	}
	// en passant left
	possibility = (bp << 1) & board.bitboards[WP] & Rank4 & (^FileA) & board.bitboards[EP]
	if possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-1, index+8, 0, 0, MoveFlagEnPass))
	}
}

func (board *Board) possibleKnightMoves(moveList *MoveList, knight uint64) {
	// Choose bishop
	knightPossibility := knight & (^(knight - 1))
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
		if knightIdx%8 < 4 {
			possibility &= (^FileGH) & NotMyPieces
		} else {
			possibility &= (^FileAB) & NotMyPieces
		}

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList.AddMove(GetMoveInt(knightIdx, moveIndex, 0, 0, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		knight &= ^knightPossibility
		knightPossibility = knight & (^(knight - 1))
	}
}

func (board *Board) possibleBishopMoves(moveList *MoveList, bishop uint64) {
	// Choose bishop
	bishopPossibility := bishop & (^(bishop - 1))
	var possibility uint64

	for bishopPossibility != 0 {
		// Current bishop index (in bitmask)
		bishopIdx := bits.TrailingZeros64(bishopPossibility)
		possibility = board.DiagonalAndAntiDiagonalMoves(bishopIdx) & NotMyPieces

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList.AddMove(GetMoveInt(bishopIdx, moveIndex, 0, 0, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		bishop &= ^bishopPossibility
		bishopPossibility = bishop & (^(bishop - 1))
	}
}

func (board *Board) possibleRookMoves(moveList *MoveList, rook uint64) {
	// Choose rook
	rookPossibility := rook & (^(rook - 1))
	var possibility uint64

	for rookPossibility != 0 {
		// Current rook index (in bitmask)
		rookIdx := bits.TrailingZeros64(rookPossibility)
		possibility = board.HorizontalAndVerticalMoves(rookIdx) & NotMyPieces

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList.AddMove(GetMoveInt(rookIdx, moveIndex, 0, 0, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		rook &= ^rookPossibility
		rookPossibility = rook & (^(rook - 1))
	}
}

func (board *Board) possibleQueenMoves(moveList *MoveList, queen uint64) {
	// Choose queen
	queenPossibility := queen & (^(queen - 1))
	var possibility uint64

	for queenPossibility != 0 {
		// Current queen index (in bitmask)
		queenIdx := bits.TrailingZeros64(queenPossibility)
		possibility = (board.HorizontalAndVerticalMoves(queenIdx) | board.DiagonalAndAntiDiagonalMoves(queenIdx)) & NotMyPieces

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList.AddMove(GetMoveInt(queenIdx, moveIndex, 0, 0, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		queen &= ^queenPossibility
		queenPossibility = queen & (^(queen - 1))
	}
}

func (board *Board) possibleKingMoves(moveList *MoveList, king uint64) {
	var possibility uint64

	// Current king index (in bitmask)
	kingIdx := bits.TrailingZeros64(king)

	if kingIdx > 9 {
		possibility = KingSpan << (kingIdx - 9)
	} else {
		possibility = KingSpan >> (9 - kingIdx)
	}

	// handle wrap around of knight pattern movement
	if kingIdx%8 < 4 {
		possibility &= (^FileGH) & NotMyPieces
	} else {
		possibility &= (^FileAB) & NotMyPieces
	}

	// choose move
	movePossibility := possibility & (^(possibility - 1))
	for movePossibility != 0 {
		// possible move index (in bitmask)
		moveIndex := bits.TrailingZeros64(movePossibility)
		moveList.AddMove(GetMoveInt(kingIdx, moveIndex, 0, 0, NoFlag))
		possibility &= ^movePossibility                      // remove move from all possible moves
		movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
	}
}

func (board *Board) possibleCastleWhite(moveList *MoveList) {
	if (board.castlePermissions&WhiteKingCastling) != 0 && (((1 << CastleRooks[0]) & board.bitboards[WR]) != 0) {
		moveList.AddMove(GetMoveInt(60, 62, 0, 0, MoveFlagCastle))
	}
	if (board.castlePermissions&WhiteQueenCastling) != 0 && (((1 << CastleRooks[1]) & board.bitboards[WR]) != 0) {
		moveList.AddMove(GetMoveInt(60, 58, 0, 0, MoveFlagCastle))
	}
}

func (board *Board) possibleCastleBlack(moveList *MoveList) {
	if (board.castlePermissions&BlackKingCastling) != 0 && (((1 << CastleRooks[2]) & board.bitboards[BR]) != 0) {
		moveList.AddMove(GetMoveInt(4, 6, 0, 0, MoveFlagCastle))
	}
	if (board.castlePermissions&BlackQueenCastling) != 0 && (((1 << CastleRooks[3]) & board.bitboards[BR]) != 0) {
		moveList.AddMove(GetMoveInt(4, 2, 0, 0, MoveFlagCastle))
	}
}

func (board *Board) unsafeForBlack() (unsafe uint64) {
	// todo should it update global value ??
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

	// pawn
	unsafe = ((board.bitboards[WP] >> 7) & (^FileA))  // pawn capture right
	unsafe |= ((board.bitboards[WP] >> 9) & (^FileH)) // pawn capture left

	var possibility uint64
	// knight
	wn := board.bitboards[WN]
	i := wn & (^(wn - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		if iLocation > 18 {
			possibility = KnightSpan << (iLocation - 18)
		} else {
			possibility = KnightSpan >> (18 - iLocation)
		}

		if iLocation%8 < 4 {
			possibility &= (^FileGH)
		} else {
			possibility &= (^FileAB)
		}
		unsafe |= possibility
		wn &= (^i)
		i = wn & (^(wn - 1))
	}

	// bishop/queen
	qb := board.bitboards[WQ] | board.bitboards[WB]
	i = qb & (^(qb - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.DiagonalAndAntiDiagonalMoves(iLocation)
		unsafe |= possibility
		qb &= (^i)
		i = qb & (^(qb - 1))
	}

	// rook/queen
	qr := board.bitboards[WQ] | board.bitboards[WR]
	i = qr & (^(qr - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.HorizontalAndVerticalMoves(iLocation)
		unsafe |= possibility
		qr &= (^i)
		i = qr & (^(qr - 1))
	}

	// king
	iLocation := bits.TrailingZeros64(board.bitboards[WK])
	if iLocation > 9 {
		possibility = KingSpan << (iLocation - 9)
	} else {
		possibility = KingSpan >> (9 - iLocation)
	}

	if iLocation%8 < 4 {
		possibility &= (^FileGH)
	} else {
		possibility &= (^FileAB)
	}
	unsafe |= possibility
	return unsafe
}

func (board *Board) unsafeForWhite() (unsafe uint64) {
	// todo should it update global value ??
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

	// pawn
	unsafe = ((board.bitboards[BP] << 7) & (^FileH))  // pawn capture right
	unsafe |= ((board.bitboards[BP] << 9) & (^FileA)) // pawn capture left

	var possibility uint64
	// knight
	bn := board.bitboards[BN]
	i := bn & (^(bn - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		if iLocation > 18 {
			possibility = KnightSpan << (iLocation - 18)
		} else {
			possibility = KnightSpan >> (18 - iLocation)
		}

		if iLocation%8 < 4 {
			possibility &= (^FileGH)
		} else {
			possibility &= (^FileAB)
		}
		unsafe |= possibility
		bn &= (^i)
		i = bn & (^(bn - 1))
	}

	// bishop/queen
	qb := board.bitboards[BQ] | board.bitboards[BB]
	i = qb & (^(qb - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.DiagonalAndAntiDiagonalMoves(iLocation)
		unsafe |= possibility
		qb &= (^i)
		i = qb & (^(qb - 1))
	}

	// rook/queen
	qr := board.bitboards[BQ] | board.bitboards[BR]
	i = qr & (^(qr - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.HorizontalAndVerticalMoves(iLocation)
		unsafe |= possibility
		qr &= (^i)
		i = qr & (^(qr - 1))
	}

	// king
	iLocation := bits.TrailingZeros64(board.bitboards[BK])
	if iLocation > 9 {
		possibility = KingSpan << (iLocation - 9)
	} else {
		possibility = KingSpan >> (9 - iLocation)
	}

	if iLocation%8 < 4 {
		possibility &= (^FileGH)
	} else {
		possibility &= (^FileAB)
	}
	unsafe |= possibility
	return unsafe
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
