package board

import (
	"fmt"
	"testing"
)

func TestPerftStartingPosition(t *testing.T) {
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	PerftMaxDepth = 6
	Perft(board, 0)
	fmt.Println(PerftMoveCounter)
	fmt.Println(PerftCaptures)
	fmt.Println(PerftEnPassant)
	fmt.Println(PerftPromotions)
	t.Errorf("Errror\n")
}

func BenchmarkPerftStartingPositionDepth5(b *testing.B) {
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	PerftMaxDepth = 5

	// This performs ~54% faster than hugo
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Perft(board, 0)
	}
	b.StopTimer()
}
