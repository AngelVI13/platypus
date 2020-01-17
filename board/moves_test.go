package board

import (
	"testing"
	"fmt"
	"time"
	"strings"
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


	start := time.Now()
	for i := 0; i < 1000; i++ {
		var moveList strings.Builder
		board.PossibleMovesWhite(&moveList)
	}
	end := time.Since(start)
	fmt.Printf("MoveGen (1000): %s\n", end)
	// moveList := board.PossibleMovesBlack()

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}

func BenchmarkPossibleMovesWhite(b *testing.B) {
	board := Board{}
	board.ParseStringArray(StartingPosition)


	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var moveList strings.Builder
		board.PossibleMovesWhite(&moveList)
	}
	b.StopTimer()
}

func BenchmarkUnwrappedPossibleMovesWhite(b *testing.B) {
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

		board.possibleWhitePawn(&moveList)
		board.possibleKnightMoves(&moveList, board.bitboards[WN])
		board.possibleBishopMoves(&moveList, board.bitboards[WB])
		board.possibleRookMoves(&moveList, board.bitboards[WR])
		board.possibleQueenMoves(&moveList, board.bitboards[WQ])
		board.possibleKingMoves(&moveList, board.bitboards[WK])
		board.possibleCastleWhite(
			&moveList,
			board.whiteCastleKingSide,
			board.whiteCastleQueenSide)
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


func BenchmarkExample1(b *testing.B) {
	board := Board{}
	board.ParseStringArray(StartingPosition)


	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.Examplefunc()
	}
	b.StopTimer()
}

func BenchmarkExample2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Examplefunc2()
	}
	b.StopTimer()
}
