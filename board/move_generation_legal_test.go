package board

import (
	"testing"
)

func TestLegalMovesWhite(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/8/4k3/8/4R3/8/8/4K3 b - - 0 1")

	// board.LegalMovesWhite()
	DrawBitboard(board.unsafeForBlack())

	// fmt.Println(len(moveList) / 4)
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
