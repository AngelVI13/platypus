package main

import (
	"fmt"
	"platypus/board"
)

func main() {
	boardVar := board.Board{}
	boardVar.ParseStringArray(board.StartingPosition)
	fmt.Println(&boardVar)
	boardVar.PrintBitboards()
}
