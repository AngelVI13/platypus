package board

import (
	"fmt"
	"testing"
)

func TestGetMoves(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("r6r/1b2k1bq/8/8/7B/8/8/R3K2R b QK - 3 2")
	moveList := board.GetMoves()
	if moveList.Count != 9 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 8 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("8/8/8/2k5/2pP4/8/B7/4K3 b - d3 5 3")
	moveList = board.GetMoves()
	if moveList.Count != 8 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 8 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("r1bqkbnr/pppppppp/n7/8/8/P7/1PPPPPPP/RNBQKBNR w QqKk - 2 2")
	moveList = board.GetMoves()
	if moveList.Count != 19 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 19 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("r3k2r/p1pp1pb1/bn2Qnp1/2qPN3/1p2P3/2N5/PPPBBPPP/R3K2R b QqKk - 3 2")
	moveList = board.GetMoves()
	if moveList.Count != 5 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 5 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("2kr3r/p1ppqpb1/bn2Qnp1/3PN3/1p2P3/2N5/PPPBBPPP/R3K2R b QK - 3 2")
	moveList = board.GetMoves()
	if moveList.Count != 44 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 44 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("rnb2k1r/pp1Pbppp/2p5/q7/2B5/8/PPPQNnPP/RNB1K2R w QK - 3 9")
	moveList = board.GetMoves()
	if moveList.Count != 39 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 39 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("2r5/3pk3/8/2P5/8/2K5/8/8 w - - 5 4")
	moveList = board.GetMoves()
	if moveList.Count != 9 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 9 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NKPP/RNBQ3R b - - 0 8")
	moveList = board.GetMoves()
	if moveList.Count != 28 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 28 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("rnbq1k1r/pp1Pbppp/2p5/8/2B5/3n4/PPPKN1PP/RNBQ3R w - - 3 9")
	moveList = board.GetMoves()
	if moveList.Count != 37 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 37 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("rnbq1k1r/pp1P1ppp/2pb4/8/2B5/8/PPPKNnPP/RNBQ3R w - - 3 9")
	moveList = board.GetMoves()
	if moveList.Count != 36 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 36 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}

	board.ParseFen("r4rk1/1pp1qppp/p1np1n2/2b3B1/2BpP1b1/P1NP4/1PP1QPPP/R4RK1 w - - 0 11")
	moveList = board.GetMoves()
	if moveList.Count != 45 {
		PrintMoveList(&moveList)
		t.Errorf("Expected 45 possible moves, got %d instead. Board: \n%s", moveList.Count, board.String())
	}
}

func TestLegalMovesWhite(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/8/5r2/3KP2b/8/8/8/3k4 w - - 0 1")
	// board.ParseFen("8/4K3/4N3/8/4r3/8/8/3k4 w - - 0 1")
	// board.ParseFen("5kq1/2b5/8/4pP2/2K5/8/8/8 w - e6 0 1")
	// board.ParseFen(StartingPosition)

	fmt.Println(&board)
	var moveList MoveList
	board.LegalMovesWhite(&moveList)
	PrintMoveList(&moveList)

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}

func TestLegalMovesBlack(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("8/4Q3/8/8/1b6/k7/8/3K4 b - - 0 1")
	// board.ParseFen("8/4K3/4N3/8/4r3/8/8/3k4 w - - 0 1")
	// board.ParseFen("5kq1/2b5/8/4pP2/2K5/8/8/8 w - e6 0 1")
	// board.ParseFen(StartingPosition)

	fmt.Println(&board)
	var moveList MoveList
	board.LegalMovesBlack(&moveList)
	PrintMoveList(&moveList)

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}

func TestGetPinnedPieceRays(t *testing.T) {
	InitHashKeys()
	board := Board{}
	board.ParseFen("r2q2q1/8/2RRR3/q1RKR2b/2RRR3/8/q1kq2q1/8 w - - 0 1")
	board.UpdateBitMasks()

	var pinRays PinRays

	board.getPinnedPieceRays(board.bitboards[WK], &pinRays)
	for i := 0; i < pinRays.Count; i++ {
		DrawBitboard(pinRays.Rays[i])
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
		var pinRays PinRays
		board.getPinnedPieceRays(board.bitboards[WK], &pinRays)
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
