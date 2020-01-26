package board

var PerftMoveCounter int
var PerftMaxDepth int
var PerftCaptures int
var PerftCastles int
var PerftEnPassant int
var PerftPromotions int

func Perft(board Board, depth int) {
	if depth < PerftMaxDepth {
		var moveList MoveList
		if board.Side == White {
			board.PossibleMovesWhite(&moveList)
		} else {
			board.PossibleMovesBlack(&moveList)
		}

		for i := 0; i < moveList.Count; i++ {
			move := moveList.Moves[i].Move
			moveBoard := board // copy board

			if moveBoard.MakeMove(move) != true {
				// fmt.Printf("Illegal move: %s\nBoard is:\n%s\n\n", GetMoveString(move), &moveBoard)
				// DrawBitboard(Unsafe)
				continue
			}

			// var currentMoveCount int
			// if depth == 0 {
			// 	currentMoveCount = PerftMoveCounter
			// }

			if (depth + 1) == PerftMaxDepth {
				PerftMoveCounter++
			}
			Perft(moveBoard, depth+1)

			// if depth == 0 {
			// 	currentMoveCount = PerftMoveCounter - currentMoveCount
			// 	fmt.Printf("%s: %d\n", GetMoveString(move), currentMoveCount)
			// }
		}
	}
}
