package board

import (
	"fmt"
	"testing"
)

func TestPerftStartingPosition(t *testing.T) {
	board := Board{}
	// board.ParseStringArray(StartingPosition)
	board.ParseStringArray([8][8]string{
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{"k", " ", " ", "P", "p", " ", " ", "Q"},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
		[8]string{" ", " ", " ", " ", "K", " ", " ", " "}})
	board.Side = Black
	// Set castle permissions
	board.castlePermissions = (int(0))
	board.bitboards[EP] |= ^(uint64(0))

	DrawBitboard(board.unsafeForBlack())
	DrawBitboard(board.bitboards[EP])

	PerftMaxDepth = 1
	Perft(board, 0)
	fmt.Println(PerftMoveCounter)
	fmt.Println(PerftCaptures)
	fmt.Println(PerftEnPassant)
	fmt.Println(PerftPromotions)
	t.Errorf("Errror\n")
}