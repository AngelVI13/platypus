package board

import (
	"fmt"
	"testing"
)

func TestLegalMovesWhite(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/4K3/8/6N1/4r3/8/8/3k4 w - - 0 1")
	// board.ParseFen("5kq1/2b5/8/4pP2/2K5/8/8/8 w - e6 0 1")
	// board.ParseFen(StartingPosition)

	fmt.Println(&board)
	var moveList MoveList
	board.LegalMovesWhite(&moveList)
	PrintMoveList(&moveList)

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}

func TestGetPinnedPieceRays(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("r2q2q1/8/2RRR3/q1RKR2b/2RRR3/8/q1kq2q1/8 w - - 0 1")
	board.UpdateBitMasks()

	rays := board.getPinnedPieceRays(board.bitboards[WK])
	for _, ray := range rays {
		DrawBitboard(ray)
	}
	t.Errorf("Error")
}

func TestGetCheckerSliderRaysToKing(t *testing.T) {
	var kingBitboard uint64 = (1 << 3)
	var checkerBitboard uint64 = (1 << 24)

	DrawBitboard(kingBitboard)
	DrawBitboard(checkerBitboard)
	DrawBitboard(kingBitboard | checkerBitboard)
	DrawBitboard(getCheckerSliderRaysToKing(kingBitboard, checkerBitboard))
	t.Errorf("Error")
}

func BenchmarkGetCheckers(b *testing.B) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/2k1r3/1b6/5nb1/2nP4/4K3/8/8 w - - 0 1")
	board.UpdateBitMasks()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.getCheckers(board.bitboards[WK])
	}
	b.StopTimer()
}

func BenchmarkGetPinnedRays(b *testing.B) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("q2q2q1/8/2RRR3/q1RKR2q/2RRR3/8/q1kq2q1/8 w - - 0 1")
	board.UpdateBitMasks()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.getPinnedPieceRays(board.bitboards[WK])
	}
	b.StopTimer()
}

func BenchmarkUpdateBitMasks(b *testing.B) {
	InitHashKeys()
	board := Board{}
	board.ParseFen(StartingPosition)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.UpdateBitMasks()
	}
	b.StopTimer()
}
