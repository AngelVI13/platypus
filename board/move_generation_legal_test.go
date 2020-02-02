package board

import (
	"testing"
)

func TestLegalMovesWhite(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/4k3/8/8/4r3/6b1/8/4K3 w - - 0 1")

	var moveList MoveList
	board.LegalMovesWhite(&moveList)
	PrintMoveList(&moveList)

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}

func TestGetCheckerSliderRaysToKing(t *testing.T) {
	var kingBitboard uint64 = (1 << 27)
	var checkerBitboard uint64 = (1 << 59)

	DrawBitboard(kingBitboard)
	DrawBitboard(checkerBitboard)
	DrawBitboard(kingBitboard | checkerBitboard)
	DrawBitboard(getCheckerSliderRaysToKing(kingBitboard, checkerBitboard))
	t.Errorf("Error")
}

func BenchmarkGetCheckers(b *testing.B) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/4r3/1b6/5nb1/2nP4/4K3/8/8 w - - 0 1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.UpdateBitMasks()
		board.getCheckers(board.bitboards[WK])
	}
	b.StopTimer()
}
