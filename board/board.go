package board

import (
	"fmt"
	"math/bits"
	"strconv"
)

// StartingPosition 8x8 representation of normal chess starting position
const StartingPosition string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const SideChar string = "wb-"
const PieceChar string = ".PNBRQKpnbrqk"

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

// MaxGameMoves Maximum number of game moves
const (
	MaxGameMoves   int = 2048
	BoardSquareNum int = 64
)

// CastlePerm used to simplify hashing castle permissions
// Everytime we make a move we will take pos.castlePerm &= CastlePerm[sq]
// in this way if any of the rooks or the king moves, the castle permission will be
// disabled for that side. In any other move, the castle permissions will remain the
// same, since 15 is the max number associated with all possible castling permissions
// for both sides
var CastlePerm = [BoardSquareNum]int{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

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
	material          [2]int             // material scores for black and white
	ply               int                // how many half moves have been made
	fiftyMove         int                // how many moves from the fifty move rule have been made
	positionKey       uint64             // position key is a unique key stored for each position (used to keep track of 3fold repetition)
	history           [MaxGameMoves]Undo // array that stores current position and variables before a move is made
}

// Reset Resets current board
func (board *Board) Reset() {
	for i := 0; i < 14; i++ {
		board.bitboards[i] = uint64(0)
	}
	for i := 0; i < BoardSquareNum; i++ {
		board.position[i] = 0
	}
	for i := 0; i < 10; i++ {
		board.stateBoards[i] = uint64(0)
	}

	board.Side = White
	board.castlePermissions = 0
	board.material[White] = 0
	board.material[Black] = 0
	board.ply = 0
	board.fiftyMove = 0
	board.positionKey = 0
}

// String Return string representing the current board (from the stored position)
func (board *Board) String() string {
	var position [8][8]string

	for i := 0; i < 64; i++ {
		if board.position[i] == WP {
			position[i/8][i%8] = "P"
		} else if board.position[i] == WN {
			position[i/8][i%8] = "N"
		} else if board.position[i] == WB {
			position[i/8][i%8] = "B"
		} else if board.position[i] == WR {
			position[i/8][i%8] = "R"
		} else if board.position[i] == WQ {
			position[i/8][i%8] = "Q"
		} else if board.position[i] == WK {
			position[i/8][i%8] = "K"
		} else if board.position[i] == BP {
			position[i/8][i%8] = "p"
		} else if board.position[i] == BN {
			position[i/8][i%8] = "n"
		} else if board.position[i] == BB {
			position[i/8][i%8] = "b"
		} else if board.position[i] == BR {
			position[i/8][i%8] = "r"
		} else if board.position[i] == BQ {
			position[i/8][i%8] = "q"
		} else if board.position[i] == BK {
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


	// ---
	positionStr += fmt.Sprintf("side:%c\n", SideChar[board.Side])

	enPassantFile := "-"
	if board.bitboards[EP] != 0 {
		enPassantFile = strconv.Itoa(bits.TrailingZeros64(board.bitboards[EP]))
	}
	positionStr += fmt.Sprintf("enPasFile:%s\n", enPassantFile)
	
	// Compute castling permissions
	wKCA := "-"
	if board.castlePermissions&WhiteKingCastling != 0 {
		wKCA = "K"
	}

	wQCA := "-"
	if board.castlePermissions&WhiteQueenCastling != 0 {
		wQCA = "Q"
	}

	bKCA := "-"
	if board.castlePermissions&BlackKingCastling != 0 {
		bKCA = "k"
	}

	bQCA := "-"
	if board.castlePermissions&BlackQueenCastling != 0 {
		bQCA = "q"
	}

	positionStr += fmt.Sprintf("castle:%s%s%s%s\n", wKCA, wQCA, bKCA, bQCA)
	positionStr += fmt.Sprintf("PosKey:%d\n", board.positionKey)

	// ---


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
