package board

import (
	"testing"
	"fmt"
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


	moveList := board.PossibleMovesWhite("1636")
	fmt.Println(len(moveList) / 4)
	t.Errorf("Error")
}
