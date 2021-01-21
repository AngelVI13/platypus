package board

import (
	"testing"
)

func TestParseFen(t *testing.T) {
	board := Board{}
	board.ParseFen(StartingPosition)
	// board.ParseFen("rnbqkbnr/pp1p1ppp/8/2pPp3/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3")
	// fmt.Println(&board)
	// fmt.Println(board.Side)
	// DrawBitboard(uint64(board.castlePermissions))
	// DrawBitboard(board.bitboards[EP])
	// t.Errorf("Blabla\n")
}

func TestParseFenPositionKey1(t *testing.T) {
	board := Board{}
	board.ParseFen("rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3")
	
	parseFenKey := board.positionKey

	board.Reset()

	if board.positionKey == parseFenKey {
		t.Errorf("Position key for a reset board equals that of a previous valid position.\n")
	}

	// Set starting position
	board.ParseFen(StartingPosition)
	
	moveSeq := "e2e4 d7d5 e4e5 f7f5"
	err := board.MakeMoves(moveSeq)
	if (err != nil) {
		t.Errorf("Error in MakeMoves for sequence: %s\n", moveSeq)
	}

	if parseFenKey != board.positionKey {
		t.Errorf("Position key of board should equal the position Key of the parsed fen position.\n")
	}
}

func TestParseFenPositionKey2(t *testing.T) {
	board := Board{}
	board.ParseFen("r2qkbnr/ppp1p1pp/B1n1b3/3pPp2/8/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 5")
	
	parseFenKey := board.positionKey

	board.Reset()

	if board.positionKey == parseFenKey {
		t.Errorf("Position key for a reset board equals that of a previous valid position.\n")
	}

	board.ParseFen(StartingPosition)
	
	moveSeq := "e2e4 d7d5 e4e5 f7f5 g1f3 c8e6 f1a6 b8c6 e1g1"
	err := board.MakeMoves(moveSeq)
	if (err != nil) {
		t.Errorf("Error in MakeMoves for sequence: %s\n", moveSeq)
	}
	if parseFenKey != board.positionKey {
		t.Errorf("Position key of board should equal the position Key of the parsed fen position.\n")
	}
}


func TestGetMoveFromString(t *testing.T) {
	board := Board{}
	board.ParseFen(StartingPosition)
	
	// Generate moves
	moveList := board.GetMoves()

	// Check that manually choosing a move or using GetMoveFromString results in the same move
	move := moveList.Moves[12].Move // e2e4
	moveFromString, err := GetMoveFromString(&moveList, "e2e4")

	if (err != nil) || (move != moveFromString) {
		t.Errorf("Error in move fetching\n")
	}
}	

func TestPerformMoves(t *testing.T) {
	board := Board{}
	// Set starting position
	board.ParseFen(StartingPosition)
	
	moveSeq := "e2e4 d7d5 e4e5"
	err := board.MakeMoves(moveSeq)

	if (err != nil) {
		t.Errorf("Error in MakeMoves for sequence: %s\n", moveSeq)
	}
}	

func TestPerformMovesNegative(t *testing.T) {
	board := Board{}
	// Set starting position
	board.ParseFen(StartingPosition)
	
	//               {err}
	moveSeq := "e2e4 d7d4 e4e5"
	err := board.MakeMoves(moveSeq)

	if (err == nil) {
		t.Errorf("There should be an error in move sequence: %s\n", moveSeq)
	}
}
