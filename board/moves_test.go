package board

import "testing"

func TestPossibleMovesWhite(t *testing.T) {
	board := Board{}
	board.ParseStringArray(StartingPosition)

	board.PossibleMovesWhite("")
	t.Errorf("Error")
}
