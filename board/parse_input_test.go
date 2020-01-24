package board


import (
	"fmt"
	"testing"
)

func TestParseFen(t *testing.T) {
	board := Board{}
	board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	board.ParseFen("rnbqkbnr/pp1p1ppp/8/2pPp3/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3")
	fmt.Println(&board)
	fmt.Println(board.Side)
	DrawBitboard(uint64(board.castlePermissions))
	DrawBitboard(board.bitboards[EP])
	t.Errorf("Blabla\n")
}