package board

import (
	"fmt"
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

// DiagonalMasks8 Bitmask to help select all diagonals
var DiagonalMasks8 [15]uint64 = [15]uint64{
	0x1, 0x102, 0x10204, 0x1020408, 0x102040810, 0x10204081020, 0x1020408102040,
	0x102040810204080, 0x204081020408000, 0x408102040800000, 0x810204080000000,
	0x1020408000000000, 0x2040800000000000, 0x4080000000000000, 0x8000000000000000,
}

// AntiDiagonalMasks8 Bitmask to help select all anti diagonal
var AntiDiagonalMasks8 [15]uint64 = [15]uint64{
	0x80, 0x8040, 0x804020, 0x80402010, 0x8040201008, 0x804020100804, 0x80402010080402,
	0x8040201008040201, 0x4020100804020100, 0x2010080402010000, 0x1008040201000000,
	0x804020100000000, 0x402010000000000, 0x201000000000000, 0x100000000000000,
}

// CastleRooks Array containing all initial rook squares
var CastleRooks [4]int = [4]int{63, 56, 7, 0}

const (
	H1 int = 63
	A1 int = 56
	A8 int = 0
	H8 int = 7
	D1 int = 59
	F1 int = 61
	D8 int = 3
	F8 int = 5
	G1 int = 62
	C1 int = 58
	C8 int = 2
	G8 int = 6
	E1 int = 60
	E8 int = 4
)

// Move Move type
type Move struct {
	Move  int // todo consider move to be uint32 instead of int
	score int
}

/* Game move - information stored in the move int from type Move
               |Ca| |--To-||-From-|
0000 0000 0000 0000 0000 0011 1111 -> From - 0x3F
0000 0000 0000 0000 1111 1100 0000 -> To - >> 6, 0x3F
0000 0000 0000 1111 0000 0000 0000 -> Piece Type - >> 12, 0xF
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

// PieceType - macro that returns the 'PieceType' bits from the move int
func PieceType(m int) int {
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
func GetMoveInt(fromSq, toSq, pieceType, promotionPiece, flag int) int {
	return fromSq | (toSq << 6) | (pieceType << 12) | (promotionPiece << 18) | flag
}

const (
	// MoveFlagEnPass move flag that denotes if the capture was an enpass
	MoveFlagEnPass int = 0x10000

	// MoveFlagPawnStart move flag that denotes if move was pawn start (2x)
	MoveFlagPawnStart int = 0x20000

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

// PinRays Struct to hold all generated pin rays for a given position
type PinRays struct {
	Rays  [8]uint64 // an array with max possible pinned rays
	Count int       // number of generated pin rays in the struct
}

// AddRay Adds ray to pin rays and updates count
func (pinRays *PinRays) AddRay(ray uint64) {
	pinRays.Rays[pinRays.Count] = ray
	pinRays.Count++
}

// GetRay Get pin ray which corresponds to a given piece. If the given piece is not
// pinned -> return ^0 i.e. the piece can move to any square and is not limmited by a pin ray
func (pinRays *PinRays) GetRay(pieceBitboard uint64) uint64 {
	for i := 0; i < pinRays.Count; i++ {
		if pinRays.Rays[i]&pieceBitboard != 0 {
			return pinRays.Rays[i]
		}
	}
	return ^uint64(0)
}

// GetSquareString get algebraic notation of square i.e. b2, a6 from array index
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

// GetMoveString prints move in algebraic notation
func GetMoveString(move int) string {
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

		fmt.Printf("Move:%d > %s (score:%d)\n", index+1, GetMoveString(move), score)
	}
	fmt.Printf("MoveList Total %d Moves:\n\n", moveList.Count)
}

// DrawBitboard Prints a given bitboard to stdout in a human readable way
func DrawBitboard(bitboard uint64) {
	var bitboardStr [8][8]string

	for i := 0; i < 64; i++ {
		if ((bitboard >> i) & 1) == 1 {
			bitboardStr[i/8][i%8] = "X"
		} else if bitboardStr[i/8][i%8] == "" {
			// replace empty string with a dot for better readibility of bitboard
			bitboardStr[i/8][i%8] = "."
		}
	}

	var positionStr string
	positionStr += "\n"
	for idx, rank := range bitboardStr {
		positionStr += fmt.Sprintf(" %d  ", (8 - idx))
		for _, file := range rank {
			positionStr += fmt.Sprintf(" %s ", file)
		}
		positionStr += "\n"
	}

	positionStr += "\n     "
	startFileIdx := "A"[0]
	for i := startFileIdx; i < startFileIdx+8; i++ {
		positionStr += fmt.Sprintf("%s  ", string(i))
	}
	positionStr += fmt.Sprintf("\n")

	fmt.Println(positionStr)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
