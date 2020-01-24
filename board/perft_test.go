package board

import (
	"fmt"
	"testing"
)

func TestPerftStartingPosition(t *testing.T) {
	board := Board{}
	// board.ParseStringArray(StartingPosition)
	board.ParseStringArray(StartingPosition)
	board.Side = White
	// Set castle permissions
	board.castlePermissions = ^(int(0))

	PerftMaxDepth = 5
	Perft(board, 0)
	fmt.Println(PerftMoveCounter)
	fmt.Println(PerftCaptures)
	fmt.Println(PerftEnPassant)
	fmt.Println(PerftPromotions)
	t.Errorf("Errror\n")
}