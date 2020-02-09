package board

import (
	"fmt"
)

// StartingPosition 8x8 representation of normal chess starting position
const StartingPosition string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

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
	EP // en passant file bitboard
)

// IsSlider maps piece type to information if it is a sliding piece or not (i.e. rook, bishop, queen)
var IsSlider = map[int]bool{
	WP: false,
	WN: false,
	WB: true,
	WR: true,
	WQ: true,
	WK: false,
	BP: false,
	BN: false,
	BB: true,
	BR: true,
	BQ: true,
	BK: false,
}

const (
	// NotMyPieces index to bitboard with all enemy and empty squares
	NotMyPieces int = iota

	// EnemyPieces index to bitboard with all the squares of enemy pieces
	EnemyPieces

	// EnemyRooksQueens index to bitboard with all the squares of enemy rooks & queens
	EnemyRooksQueens

	// EnemyBishopsQueens index to bitboard with all the squares of enemy bishops & queens
	EnemyBishopsQueens

	// EnemyKnights index to bitboard with all the squares of enemy knights
	EnemyKnights

	// EnemyPawns index to bitboard with all the squares of enemy pawns
	EnemyPawns

	// Empty index to bitboard with all the empty squares
	Empty

	// Occupied index to bitboard with all the occupied squares
	Occupied

	// Unsafe index to bitboard with all the unsafe squares for the current side
	Unsafe
)

// Defines for colours
const (
	White int = iota
	Black
	Both
)

// CastlePerm used to simplify hashing castle permissions
// Everytime we make a move we will take pos.castlePerm &= CastlePerm[sq]
// in this way if any of the rooks or the king moves, the castle permission will be
// disabled for that side. In any other move, the castle permissions will remain the
// same, since 15 is the max number associated with all possible castling permissions
// for both sides
var CastlePerm = [64]int{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

// MaxGameMoves Maximum number of game moves
const (
	MaxGameMoves   int = 2048
	BoardSquareNum int = 64
)

// Defines for castling rights
// The values are such that they each represent a bit from a 4 bit int value
// for example if white can castle kingside and black can castle queenside
// the 4 bit int value is going to be 1001
const (
	WhiteKingCastling  int = 1
	WhiteQueenCastling int = 2
	BlackKingCastling  int = 4
	BlackQueenCastling int = 8
)

// Undo struct
type Undo struct {
	move              int
	castlePermissions int
	enPassantFile     int
	fiftyMove         int
	posKey            uint64
}

// Board Struct to represent the chess board
type Board struct {
	bitboards         [13]uint64
	stateBoards       [9]uint64 // bitboards representing a state i.e. EnemyPieces, Empty, Occupied etc.
	Side              int
	castlePermissions int
	ply               int                // how many half moves have been made
	fiftyMove         int                // how many moves from the fifty move rule have been made
	positionKey       uint64             // position key is a unique key stored for each position (used to keep track of 3fold repetition)
	history           [MaxGameMoves]Undo // array that stores current position and variables before a move is made
}

// Reset Resets current board
func (board *Board) Reset() {
	// todo update with all variables
	for i := 0; i < 13; i++ {
		board.bitboards[i] = uint64(0)
	}
	board.Side = White
	board.castlePermissions = 0
	board.fiftyMove = 0
}

// String Return string representing the current board (from the stored bitboards)
func (board *Board) String() string {
	var position [8][8]string

	for i := 0; i < 64; i++ {
		position[i/8][i%8] = " "
	}

	// todo unicode should be used only on linux -> add a flag to disable it on linux
	for i := 0; i < 64; i++ {
		if ((board.bitboards[WP] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u265F" // P
		}
		if ((board.bitboards[WN] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u265e" // N
		}
		if ((board.bitboards[WB] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u265d" // B
		}
		if ((board.bitboards[WR] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u265c" // R
		}
		if ((board.bitboards[WQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u265b" // Q
		}
		if ((board.bitboards[WK] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u265a" // K
		}
		if ((board.bitboards[BP] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u2659" // p
		}
		if ((board.bitboards[BN] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u2658" // n
		}
		if ((board.bitboards[BB] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u2657" // b
		}
		if ((board.bitboards[BR] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u2656" // r
		}
		if ((board.bitboards[BQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u2655" // q
		}
		if ((board.bitboards[BK] >> i) & 1) == 1 {
			position[i/8][i%8] = "\u2654" // k
		}
	}

	var positionStr string
	positionStr += "\n\n    \u05c0 "
	startFileIdx := "A"[0]
	for i := startFileIdx; i < startFileIdx+8; i++ {
		positionStr += fmt.Sprintf("%s \u05c0 ", string(i))
	}
	positionStr += fmt.Sprintf("\n")

	for idx, rank := range position {
		positionStr += fmt.Sprintf("    ________________________________\n")
		// positionStr += fmt.Sprintf("\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\u2796\n")
		positionStr += fmt.Sprintf(" %d  \u05c0", (8 - idx))
		// for _, file := range rank {
		for i := (8 - 1); i >= 0; i-- {
			positionStr += fmt.Sprintf(" %s \u05c0", rank[i]) // \u05c0
		}
		positionStr += "\n"
	}
	positionStr += fmt.Sprintf("    ________________________________\n")

	// positionStr += "   \u05c0 "
	// for i := startFileIdx; i < startFileIdx+8; i++ {
	// 	positionStr += fmt.Sprintf("%s \u05c0 ", string(i))
	// }
	positionStr += fmt.Sprintf("\n")

	return positionStr
}

func (board *Board) PrintBoard() {
	var position [8][8]string

	for i := 0; i < 64; i++ {
		position[i/8][i%8] = " "
	}

	for i := 0; i < 64; i++ {
		if ((board.bitboards[WP] >> i) & 1) == 1 {
			position[i/8][i%8] = "P" // P
		}
		if ((board.bitboards[WN] >> i) & 1) == 1 {
			position[i/8][i%8] = "N" // N
		}
		if ((board.bitboards[WB] >> i) & 1) == 1 {
			position[i/8][i%8] = "B" // B
		}
		if ((board.bitboards[WR] >> i) & 1) == 1 {
			position[i/8][i%8] = "R" // R
		}
		if ((board.bitboards[WQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "Q" // Q
		}
		if ((board.bitboards[WK] >> i) & 1) == 1 {
			position[i/8][i%8] = "K" // K
		}
		if ((board.bitboards[BP] >> i) & 1) == 1 {
			position[i/8][i%8] = "p" // p
		}
		if ((board.bitboards[BN] >> i) & 1) == 1 {
			position[i/8][i%8] = "n" // n
		}
		if ((board.bitboards[BB] >> i) & 1) == 1 {
			position[i/8][i%8] = "b" // b
		}
		if ((board.bitboards[BR] >> i) & 1) == 1 {
			position[i/8][i%8] = "r" // r
		}
		if ((board.bitboards[BQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "q" // q
		}
		if ((board.bitboards[BK] >> i) & 1) == 1 {
			position[i/8][i%8] = "k" // k
		}
	}

	var positionStr string
	for _, rank := range position {
		positionStr += fmt.Sprint(rank)
		positionStr += "\n"
	}
	positionStr += fmt.Sprintf("\n")
	fmt.Println(positionStr)
}

// PrintBitboards Outputs all internal bitboards in an easy to view way
func (board *Board) PrintBitboards() {
	for idx, bitboard := range board.bitboards {
		switch idx {
		case WP:
			fmt.Printf("WhitePawn bitboard:\n")
		case BP:
			fmt.Printf("BlackPawn bitboard:\n")
		case WR:
			fmt.Printf("WhiteRook bitboard:\n")
		case BR:
			fmt.Printf("BlackRook bitboard:\n")
		case WN:
			fmt.Printf("WhiteKnight bitboard:\n")
		case BN:
			fmt.Printf("BlackKnight bitboard:\n")
		case WB:
			fmt.Printf("WhiteBishop bitboard:\n")
		case BB:
			fmt.Printf("BlackBishop bitboard:\n")
		case WQ:
			fmt.Printf("WhiteQueen bitboard:\n")
		case BQ:
			fmt.Printf("BlackQueen bitboard:\n")
		case WK:
			fmt.Printf("WhiteKing bitboard:\n")
		case BK:
			fmt.Printf("BlackKing bitboard:\n")
		case EP:
			fmt.Printf("EnPassant bitboard:\n")
		}
		DrawBitboard(bitboard)
		fmt.Println()
	}
}

// UpdateBitMasks Updates all move generation/making related bit masks
func (board *Board) UpdateBitMasks() {
	if board.Side == White {
		board.stateBoards[NotMyPieces] = ^(board.bitboards[WP] |
			board.bitboards[WN] |
			board.bitboards[WB] |
			board.bitboards[WR] |
			board.bitboards[WQ] |
			board.bitboards[WK] |
			board.bitboards[BK])

		board.stateBoards[EnemyPieces] = (board.bitboards[BP] |
			board.bitboards[BN] |
			board.bitboards[BB] |
			board.bitboards[BR] |
			board.bitboards[BQ])

		board.stateBoards[Occupied] = (board.bitboards[WP] |
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

		board.stateBoards[EnemyRooksQueens] = (board.bitboards[BQ] | board.bitboards[BR])
		board.stateBoards[EnemyBishopsQueens] = (board.bitboards[BQ] | board.bitboards[BB])
		board.stateBoards[EnemyKnights] = board.bitboards[BN]
		board.stateBoards[EnemyPawns] = board.bitboards[BP]

		board.stateBoards[Empty] = ^board.stateBoards[Occupied]
		board.stateBoards[Unsafe] = board.unsafeForWhite()
	} else {
		board.stateBoards[NotMyPieces] = ^(board.bitboards[BP] |
			board.bitboards[BN] |
			board.bitboards[BB] |
			board.bitboards[BR] |
			board.bitboards[BQ] |
			board.bitboards[BK] |
			board.bitboards[WK])

		board.stateBoards[EnemyPieces] = (board.bitboards[WP] |
			board.bitboards[WN] |
			board.bitboards[WB] |
			board.bitboards[WR] |
			board.bitboards[WQ])

		board.stateBoards[Occupied] = (board.bitboards[WP] |
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

		board.stateBoards[EnemyRooksQueens] = (board.bitboards[WQ] | board.bitboards[WR])
		board.stateBoards[EnemyBishopsQueens] = (board.bitboards[WQ] | board.bitboards[WB])
		board.stateBoards[EnemyKnights] = board.bitboards[WN]
		board.stateBoards[EnemyPawns] = board.bitboards[WP]

		board.stateBoards[Empty] = ^board.stateBoards[Occupied]
		board.stateBoards[Unsafe] = board.unsafeForBlack()
	}
}
