package board

import (
	"fmt"
	"strings"
	"testing"
)

func TestPossibleMovesWhite(t *testing.T) {
	board := Board{}
	// board.ParseStringArray(StartingPosition)
	board.ParseStringArray([8][8]string{
		[8]string{"r", "n", "b", "q", "k", "b", "n", "r"},
		[8]string{"p", "p", "p", "p", "p", "p", "p", "p"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"P", "P", "P", "P", "P", "P", "P", "P"},
		[8]string{"R", "N", "B", "Q", "K", "B", "N", "R"}})

	// enabl all files on enpassant bitboard
	board.bitboards[EP] = ^uint64(0)

	var moveList MoveList
	// board.PossibleMovesWhite(&moveList)
	board.PossibleMovesBlack(&moveList)
	PrintMoveList(&moveList)
	// moveList := board.PossibleMovesBlack()

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}

func TestPossiblePawnMovesWhite(t *testing.T) {
	board := Board{}
	// board.ParseStringArray(StartingPosition)
	board.ParseStringArray([8][8]string{
		[8]string{" ", "n", "b", "q", "k", "b", "n", "r"},
		[8]string{"P", "p", " ", "p", " ", "p", "P", "p"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", "p", "P", "p", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"P", "P", "P", "P", "P", "P", "P", "P"},
		[8]string{"R", "N", "B", "Q", "K", "B", "N", "R"}})

	// enable all files on enpassant bitboard
	board.bitboards[EP] = ^uint64(0)

	var moveList MoveList
	board.PossibleMovesWhite(&moveList)

	if moveList.Count+1 != 36 {
		t.Errorf("Incorrect number of moves for white pawns")
	}
}

func TestPossiblePawnMovesBlack(t *testing.T) {
	board := Board{}
	// board.ParseStringArray(StartingPosition)
	board.ParseStringArray([8][8]string{
		[8]string{"r", "n", "b", "q", "k", "b", "n", "r"},
		[8]string{"p", "p", "p", "p", "p", "p", "p", "p"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", "P", "p", "P", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"p", "P", " ", "P", " ", "P", "p", "P"},
		[8]string{" ", "N", "B", "Q", "K", "B", "N", "R"}})

	// enable all files on enpassant bitboard
	board.bitboards[EP] = ^uint64(0)

	var moveList MoveList
	board.PossibleMovesBlack(&moveList)

	if moveList.Count+1 != 36 {
		t.Errorf("Incorrect number of moves for white pawns")
	}
}

func TestGetMoveIntPawnStart(t *testing.T) {
	move := GetMoveInt(63, 0, WP, WQ, MoveFlagPawnStart)

	if FromSq(move) != 63 {
		t.Error("Wrong FROM square")
	}
	if ToSq(move) != 0 {
		t.Error("Wrong TO square")
	}
	if Promoted(move) != WQ {
		t.Error("Wrong PROMOTED piece")
	}
	if Captured(move) != WP {
		t.Error("Wrong CAPTURED square")
	}
	if EnPassantFlag(move) != 0 {
		t.Error("Wrong EN_PASSANT_FLAG square")
	}
	if PawnStartFlag(move) != 1 {
		t.Error("Wrong PAWN_START_FLAG square")
	}

	if CastleFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}
}

func TestGetMoveIntCastle(t *testing.T) {
	move := GetMoveInt(15, 27, WR, WN, MoveFlagCastle)

	if FromSq(move) != 15 {
		t.Error("Wrong FROM square")
	}
	if ToSq(move) != 27 {
		t.Error("Wrong TO square")
	}
	if Promoted(move) != WN {
		t.Error("Wrong PROMOTED piece")
	}
	if Captured(move) != WR {
		t.Error("Wrong CAPTURED square")
	}
	if EnPassantFlag(move) != 0 {
		t.Error("Wrong EN_PASSANT_FLAG square")
	}
	if PawnStartFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}

	if CastleFlag(move) != 1 {
		t.Error("Wrong PAWN_START_FLAG square")
	}
}

func TestGetMoveIntEnPassant(t *testing.T) {
	move := GetMoveInt(31, 4, WB, WR, MoveFlagEnPass)

	if FromSq(move) != 31 {
		t.Error("Wrong FROM square")
	}
	if ToSq(move) != 4 {
		t.Error("Wrong TO square")
	}
	if Promoted(move) != WR {
		t.Error("Wrong PROMOTED piece")
	}
	if Captured(move) != WB {
		t.Error("Wrong CAPTURED square")
	}
	if EnPassantFlag(move) != 1 {
		t.Error("Wrong EN_PASSANT_FLAG square")
	}
	if PawnStartFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}

	if CastleFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}
}

func BenchmarkPossibleMovesWhite(b *testing.B) {
	board := Board{}
	board.ParseStringArray(StartingPosition)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var moveList MoveList
		board.PossibleMovesWhite(&moveList)
	}
	b.StopTimer()
}

func BenchmarkUnwrappedPossibleMovesWhite(b *testing.B) {
	board := Board{}
	board.ParseStringArray(StartingPosition)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var moveList MoveList
		NotMyPieces = ^(board.bitboards[WP] |
			board.bitboards[WN] |
			board.bitboards[WB] |
			board.bitboards[WR] |
			board.bitboards[WQ] |
			board.bitboards[WK] |
			board.bitboards[BK])

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

		board.possibleWhitePawn(&moveList)
		// board.possibleKnightMoves(&moveList, board.bitboards[WN])
		// board.possibleBishopMoves(&moveList, board.bitboards[WB])
		// board.possibleRookMoves(&moveList, board.bitboards[WR])
		// board.possibleQueenMoves(&moveList, board.bitboards[WQ])
		// board.possibleKingMoves(&moveList, board.bitboards[WK])
		// board.possibleCastleWhite(
		// 	&moveList,
		// 	board.whiteCastleKingSide,
		// 	board.whiteCastleQueenSide)
	}
	b.StopTimer()
}

func BenchmarkUpdateVarsPossibleMovesWhite(b *testing.B) {
	board := Board{}
	board.ParseStringArray(StartingPosition)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var moveList strings.Builder
		NotMyPieces = ^(board.bitboards[WP] |
			board.bitboards[WN] |
			board.bitboards[WB] |
			board.bitboards[WR] |
			board.bitboards[WQ] |
			board.bitboards[WK] |
			board.bitboards[BK])

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
		for j := 0; j < 50; j++ {
			moveList.WriteString(fmt.Sprintf("hel%d", j))
		}
	}
	b.StopTimer()
}

func BenchmarkUpdateVarsPossibleMovesWhiteOptimized(b *testing.B) {
	board := Board{}
	board.ParseStringArray(StartingPosition)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var moveList MoveList
		NotMyPieces = ^(board.bitboards[WP] |
			board.bitboards[WN] |
			board.bitboards[WB] |
			board.bitboards[WR] |
			board.bitboards[WQ] |
			board.bitboards[WK] |
			board.bitboards[BK])

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
		for j := 0; j < 50; j++ {
			moveList.AddMove(GetMoveInt(63, 63, WP, WB, MoveFlagEnPass))
		}
	}
	b.StopTimer()
}
