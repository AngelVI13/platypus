package board

import (
	"testing"
	"fmt"
	"time"
)

func TestPossibleMovesWhite(t *testing.T) {
	board := Board{}
	// board.ParseStringArray(StartingPosition)
	board.ParseStringArray([8][8]string{
		[8]string{"r", "n", "b", "q", "k", "b", "n", "r"},
		[8]string{"p", "p", "p", "p", "p", "p", "p", "p"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"P", "P", "P", "P", "P", "P", "P", "P"},
		[8]string{"R", "N", "B", "Q", "K", "B", "N", "R"}})


	start := time.Now()
	for i := 0; i < 1000; i++ {
		board.PossibleMovesWhite()
	}
	end := time.Since(start)
	fmt.Printf("MoveGen (1000): %s\n", end)
	// moveList := board.PossibleMovesBlack()

	// fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}
