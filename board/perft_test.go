package board

import (
	"fmt"
	"testing"
)

func TestPerftStartingPosition(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	PerftMaxDepth = 4
	Perft(board, 0)
	fmt.Println(PerftMoveCounter)
	fmt.Println(PerftCaptures)
	fmt.Println(PerftEnPassant)
	fmt.Println(PerftPromotions)
	t.Errorf("Errror\n")
}

func TestPerftPosition1(t *testing.T) {
	InitHashKeys()
	board := Board{}
	// rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8
	// rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPPKNnPP/RNBQ3R b - - 2 8
	board.ParseFen("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPPKNnPP/RNBQ3R b - - 2 8")
	PerftMaxDepth = 2
	Perft(board, 0)
	if PerftMoveCounter != 62379 {
		t.Errorf("Expected 62379 possible moves, got %d instead.", PerftMoveCounter)
	}
	// fmt.Println(PerftMoveCounter)
	// fmt.Println(PerftCaptures)
	// fmt.Println(PerftEnPassant)
	// fmt.Println(PerftPromotions)
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
