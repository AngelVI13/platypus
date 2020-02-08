package board

import (
	"testing"
)

func TestLegalMovesWhite(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("5kq1/2b5/8/4pP2/2K5/8/8/8 w - e6 0 1")

	var moveList MoveList
	board.LegalMovesWhite(&moveList)
	PrintMoveList(&moveList)

	// fmt.Println(len(moveList) / 4)
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.UpdateBitMasks()
		board.getCheckers(board.bitboards[WK])
	}
	b.StopTimer()
}
