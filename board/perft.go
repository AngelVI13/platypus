package board

import "fmt"

var PerftMoveCounter int
var PerftMaxDepth int
var PerftCaptures int
var PerftCastles int
var PerftEnPassant int
var PerftPromotions int

func Perft(board Board, depth int) {
	if depth < PerftMaxDepth {
		moveList := board.GetMoves()
		// PrintMoveList(&moveList)

		for i := 0; i < moveList.Count; i++ {
			move := moveList.Moves[i].Move
			moveBoard := board // copy board

			// fmt.Printf("Making move: %s\n", GetMoveString(move))
			moveBoard.MakeMove(move)
			// fmt.Println(&moveBoard)

			var currentMoveCount int
			if depth == 0 {
				currentMoveCount = PerftMoveCounter
			}

			if (depth + 1) == PerftMaxDepth {
				PerftMoveCounter++
			}
			Perft(moveBoard, depth+1)

			if depth == 0 {
				currentMoveCount = PerftMoveCounter - currentMoveCount
				fmt.Printf("%s: %d\n", GetMoveString(move), currentMoveCount)
			}
		}
	}
}
