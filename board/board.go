package board

import (
	"fmt"
)

// StartingPosition 8x8 representation of normal chess starting position
const StartingPosition string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Indexes to access bitboards i.e. WP - white pawn, BB - black bishop
const (
	NoPiece int = iota
	WP
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

	// MyPieces index to bitboard with all my pieces excluding my king
	MyPieces int = iota

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
	enPassantFile     uint64
	fiftyMove         int
	positionKey       uint64
}

// Board Struct to represent the chess board
type Board struct {
	position          [BoardSquareNum]int // var to keep track of all pieces on the board
	bitboards         [14]uint64          // 0- empty, 1-12 pieces WP-BK, 13 - en passant
	stateBoards       [10]uint64          // bitboards representing a state i.e. EnemyPieces, Empty, Occupied etc.
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
	for i := 0; i < 14; i++ {
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
		if ((board.bitboards[WP] >> i) & 1) == 1 {
			position[i/8][i%8] = "P"
		} else if ((board.bitboards[WN] >> i) & 1) == 1 {
			position[i/8][i%8] = "N"
		} else if ((board.bitboards[WB] >> i) & 1) == 1 {
			position[i/8][i%8] = "B"
		} else if ((board.bitboards[WR] >> i) & 1) == 1 {
			position[i/8][i%8] = "R"
		} else if ((board.bitboards[WQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "Q"
		} else if ((board.bitboards[WK] >> i) & 1) == 1 {
			position[i/8][i%8] = "K"
		} else if ((board.bitboards[BP] >> i) & 1) == 1 {
			position[i/8][i%8] = "p"
		} else if ((board.bitboards[BN] >> i) & 1) == 1 {
			position[i/8][i%8] = "n"
		} else if ((board.bitboards[BB] >> i) & 1) == 1 {
			position[i/8][i%8] = "b"
		} else if ((board.bitboards[BR] >> i) & 1) == 1 {
			position[i/8][i%8] = "r"
		} else if ((board.bitboards[BQ] >> i) & 1) == 1 {
			position[i/8][i%8] = "q"
		} else if ((board.bitboards[BK] >> i) & 1) == 1 {
			position[i/8][i%8] = "k"
		} else {
			position[i/8][i%8] = "."
		}
	}

	var positionStr string
	positionStr += "\n"
	for idx, rank := range position {
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

	return positionStr
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

		board.stateBoards[MyPieces] = (board.bitboards[WP] |
			board.bitboards[WN] |
			board.bitboards[WB] |
			board.bitboards[WR] |
			board.bitboards[WQ])

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

		board.stateBoards[MyPieces] = (board.bitboards[BP] |
			board.bitboards[BN] |
			board.bitboards[BB] |
			board.bitboards[BR] |
			board.bitboards[BQ])

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

// GetMoves Returns a struct that holds all the possible moves for a given position
func (board *Board) GetMoves() (moveList MoveList) {
	if board.Side == White {
		board.LegalMovesWhite(&moveList)
	} else {
		board.LegalMovesBlack(&moveList)
	}
	return moveList
}
