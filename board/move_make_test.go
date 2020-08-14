package board

import (
	"testing"
)

func TestMakeMoveStartPos(t *testing.T) {
	// Create a board with starting position.
	// Make a move, take the move back.
	// Expect that original position key will match
	// The key after the move is taken back

	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	moveList := board.GetMoves()
	move := moveList.Moves[0].Move
	originalKey := board.positionKey

	board.MakeMove(move)
	board.TakeMove()
	
	if board.positionKey != originalKey {
		t.Errorf("PosKey mismatch: %X != %X\n", board.positionKey, originalKey)
	}
}

func TestMakeMoveTwoStartPos(t *testing.T) {
	// Create a board with starting position.
	// Make two random moves, take them back.
	// Expect that original position key will match
	// the key after the move is taken back

	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// first move
	moveList := board.GetMoves()
	move := moveList.Moves[0].Move
	originalKey := board.positionKey

	board.MakeMove(move)
	
	// second move
	moveList = board.GetMoves()
	move = moveList.Moves[0].Move
	
	board.MakeMove(move)
	
	// take back moves
	board.TakeMove()
	board.TakeMove()
	
	if board.positionKey != originalKey {
		t.Errorf("PosKey mismatch: %X != %X\n", board.positionKey, originalKey)
	}
}

func TestMakeMoveEnPassant(t *testing.T) {
	// Create a board with starting position.
	// Make an enpassant move, take the move back.
	// Expect that original position key will match
	// the key after the move is taken back

	InitHashKeys()
	board := Board{}
	board.ParseFen("8/8/8/2k5/2pP4/8/B7/4K3 b - d3 5 3")

	moveList := board.GetMoves()
	move := moveList.Moves[7].Move
	originalKey := board.positionKey

	board.MakeMove(move)
	board.TakeMove()
	
	if board.positionKey != originalKey {
		t.Errorf("PosKey mismatch: %X != %X\n", board.positionKey, originalKey)
	}
}

func TestMakeMoveCapture(t *testing.T) {
	// Create a board with starting position.
	// Make a capture move, take the move back.
	// Expect that original position key will match
	// the key after the move is taken back

	InitHashKeys()
	board := Board{}
	board.ParseFen("8/8/8/2k5/2pP4/8/B7/4K3 b - d3 5 3")

	moveList := board.GetMoves()
	move := moveList.Moves[6].Move
	originalKey := board.positionKey

	board.MakeMove(move)
	board.TakeMove()
	
	if board.positionKey != originalKey {
		t.Errorf("PosKey mismatch: %X != %X\n", board.positionKey, originalKey)
	}
}

func TestMakeMoveCastle(t *testing.T) {
	// Create a board with starting position.
	// Make a castling move, take the move back.
	// Expect that original position key will match
	// the key after the move is taken back

	InitHashKeys()
	board := Board{}
	board.ParseFen("r3k2r/p1pp1pb1/bn3np1/2qPN3/1p2P3/2N5/PPPBBPPP/R3K2R b QqKk - 3 2")

	moveList := board.GetMoves()
	move := moveList.Moves[3].Move
	originalKey := board.positionKey

	board.MakeMove(move)
	board.TakeMove()
	
	if board.positionKey != originalKey {
		t.Errorf("PosKey mismatch: %X != %X\n", board.positionKey, originalKey)
	}
}

func TestMakeMovePromotion(t *testing.T) {
	// Create a board with starting position.
	// Make a promotion move, take the move back.
	// Expect that original position key will match
	// the key after the move is taken back

	InitHashKeys()
	board := Board{}
	board.ParseFen("r3k2r/p1pp1pb1/bn3np1/2qPN3/4P3/2N5/PpPBBPPP/R3K2R b KQkq - 0 1")

	moveList := board.GetMoves()
	move := moveList.Moves[8].Move
	originalKey := board.positionKey

	board.MakeMove(move)
	board.TakeMove()

	if board.positionKey != originalKey {
		t.Errorf("PosKey mismatch: %X != %X\n", board.positionKey, originalKey)
	}
}
