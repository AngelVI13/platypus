package board


import ("fmt")

var PerftMoveCounter int
var PerftMaxDepth int
var PerftCaptures int
var PerftCastles int
var PerftEnPassant int
var PerftPromotions int


func Perft(board Board, depth int) {
	if depth < PerftMaxDepth {
		// fmt.Printf("Perft: Current depth: %d\n", depth)
		// fmt.Printf("Perft: Current side: %d\n", board.Side)
		var moveList MoveList
		if board.Side == White {
			board.PossibleMovesWhite(&moveList)
		} else {
			board.PossibleMovesBlack(&moveList)
		}

		for i := 0; i < moveList.Count; i++ {
			move := moveList.Moves[i].Move
			moveBoard := board  // copy board

			if moveBoard.MakeMove(move) != true {
				// fmt.Printf("Illegal move: %s\nBoard is:\n%s\n\n", GetMoveString(move), &moveBoard)
				// DrawBitboard(Unsafe)
				continue
			}

			fmt.Printf("Made move %s\nBoard is:\n%s\n\n", GetMoveString(move), &moveBoard)
			if (depth + 1) == PerftMaxDepth {
				PerftMoveCounter++
			}
			Perft(moveBoard, depth+1)
		}
	}
}