package board

import (
	"io/ioutil"
	"encoding/json"
	"testing"
)

// todo add automatic testing of all positions from test_positions.json

func TestPerftStartingPosition(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	PerftMoveCounter = 0
	PerftMaxDepth = 5
	Perft(&board, 0)
	// fmt.Println(PerftMoveCounter)
	// fmt.Println(PerftCaptures)
	// fmt.Println(PerftEnPassant)
	// fmt.Println(PerftPromotions)
	if PerftMoveCounter != 4865609 {
		t.Errorf("Expected 4865609 moves at depth 5 from starting position, got %d\n", PerftMoveCounter)
	}
}

func TestPerftPosition1(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")

	PerftMoveCounter = 0
	PerftMaxDepth = 3
	Perft(&board, 0)
	if PerftMoveCounter != 62379 {
		t.Errorf("Expected 62379 possible moves, got %d instead.", PerftMoveCounter)
	}
}

func TestPerftPosition2(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")

	PerftMoveCounter = 0
	PerftMaxDepth = 3
	Perft(&board, 0)
	if PerftMoveCounter != 89890 {
		t.Errorf("Expected 89890 possible moves, got %d instead.", PerftMoveCounter)
	}
}

type PerftPosition struct {
	Depth int `json:"depth"`
	Nodes int `json:"nodes"`
	Fen string `json:"fen"`
}

func TestPerftPositions(t *testing.T) {
	testFile := "../test_positions.json"

	InitHashKeys()
	var positions []PerftPosition
	
	dat, err := ioutil.ReadFile(testFile)
	if err != nil {
		panic(err)
	}
	
	err = json.Unmarshal(dat, &positions)
	if err != nil {
		panic(err)
	}

	for i, position := range positions {
		board := Board{}
		board.ParseFen(position.Fen)
		
		PerftMoveCounter = 0
		PerftMaxDepth = position.Depth
		Perft(&board, 0)

		if PerftMoveCounter != position.Nodes {
			t.Errorf("I=%d: Expected %d possible moves, got %d instead. \nFEN: %s\n", 
				i, position.Nodes, PerftMoveCounter, position.Fen)
		}
	}
}

func BenchmarkPerftStartingPositionDepth3(b *testing.B) {
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	PerftMoveCounter = 0
	PerftMaxDepth = 3

	// This performs ~54% faster than hugo
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Perft(&board, 0)
	}
	b.StopTimer()
}
