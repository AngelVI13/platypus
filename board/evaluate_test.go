package board

import (
	"testing"
	"fmt"
)


func TestMaterialStartPos(t *testing.T) {
	// Create a board with starting position.
	// Expect that material for both black and white will be equal
	
	InitHashKeys()
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	
	whiteMaterial := board.material[White]
	blackMaterial := board.material[Black]

	if whiteMaterial == 0 || blackMaterial == 0 || whiteMaterial != blackMaterial {
		t.Errorf(
			"Incorrect material for startpos: White: %d != Black %d\n",
			whiteMaterial, blackMaterial,
		)
	}
}

func TestMaterialStartPosWithMoves(t *testing.T) {
	// Create a board with given position.
	// Make a capturing move. Take back the move
	// Expect that initial material for both black and white
	// will be equal to the final material after the make move and take move
	
	InitHashKeys()
	board := Board{}
	board.ParseFen("r3k2r/p1pp1pb1/bn3np1/2qPN3/4P3/2N5/PpPBBPPP/R3K2R b KQkq - 0 1")
	
	whiteMaterial := board.material[White]
	blackMaterial := board.material[Black]

	moveList := board.GetMoves()
	move := moveList.Moves[8].Move // move is b2a1q
	
	board.MakeMove(move)
	board.TakeMove()


	if whiteMaterial != board.material[White] || blackMaterial != board.material[Black] {
		t.Errorf(
			"Incorrect material for position: White: %d != Black %d\n",
			whiteMaterial, blackMaterial,
		)
	}
}

func TestMirrorBoard(t *testing.T) {
	// Create a board with starting position.
	// Expect that material for both black and white will be equal
	
	InitHashKeys()
	board := Board{}
	board.ParseFen(StartingPosition)
	
	whiteMaterial := board.material[White]
	blackMaterial := board.material[Black]
	positionKey := board.positionKey

	fmt.Println(&board)

	MirrorBoard(&board)

	fmt.Println(&board)

	if whiteMaterial != board.material[White] {
		t.Errorf(
			"Incorrect material for white startpos after mirror: Initial: %d != Mirror %d\n",
			whiteMaterial, board.material[White],
		)
	}

	if blackMaterial != board.material[Black] {
		t.Errorf(
			"Incorrect material for black startpos after mirror: Initial: %d != Mirror %d\n",
			blackMaterial, board.material[Black],
		)
	}

	fmt.Println(board.positionKey)
	if positionKey != board.positionKey {
		t.Errorf(
			"Incorrect positionKey for startpos after mirror: Initial: %d != Mirror %d\n",
			positionKey, board.positionKey,
		)
	}
}

func TestGeneratePositionKey(t *testing.T) {
	// Create a board with starting position.
	// Expect position key after ParseFen and GeneratePositionKey for the same position
	// are equal
	
	InitHashKeys()
	board := Board{}
	board.ParseFen(StartingPosition)
	
	if GeneratePositionKey(&board) != board.positionKey {
		t.Errorf(
			"PositionKey mismatch: ParseFen: %d != GeneratePositionKey %d\n",
			board.positionKey, GeneratePositionKey(&board),
		)
	}
}
