package board

import (
	"testing"
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

	positions := []string{
		StartingPosition, 
		"rnbqkbnr/pp1ppppp/2p5/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
		"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		"3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1",
		"8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1",
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
	}

	for idx, position := range positions {
		board := Board{}
		board.ParseFen(position)
		
		whiteMaterial := board.material[White]
		blackMaterial := board.material[Black]

		MirrorBoard(&board)

		// what was initially from whites side is now from blacks
		if whiteMaterial != board.material[Black] {
			t.Errorf(
				"(Pos: %d) Incorrect material for white after mirror: Initial: %d != Mirror %d\n",
				idx, whiteMaterial, board.material[Black],
			)
		}

		if blackMaterial != board.material[White] {
			t.Errorf(
				"(Pos: %d) Incorrect material for black after mirror: Initial: %d != Mirror %d\n",
				idx, blackMaterial, board.material[White],
			)
		}
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
