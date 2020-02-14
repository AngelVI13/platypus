package board

import (
	"fmt"
	"testing"
)

func TestPerftStartingPosition(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	PerftMaxDepth = 5
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
	board.ParseFen("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	PerftMaxDepth = 3
	Perft(board, 0)
	if PerftMoveCounter != 62379 {
		t.Errorf("Expected 62379 possible moves, got %d instead.", PerftMoveCounter)
	}
}

func TestPerftPosition2(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")
	PerftMaxDepth = 3
	Perft(board, 0)
	if PerftMoveCounter != 89890 {
		t.Errorf("Expected 89890 possible moves, got %d instead.", PerftMoveCounter)
	}
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
