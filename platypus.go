package main

import (
	"fmt"
	"platypus/board"
	"time"
)

func main() {
	boardVar := board.Board{}
	boardVar.ParseStringArray(board.StartingPosition)
	fmt.Println(&boardVar)

	// todo refactor moves instead of string as an array or slice ?
	// todo right now loosing a lot of speed from this
	start := time.Now()
	for i := 0; i < 1000; i++ {
		boardVar.PossibleMovesWhite()
	}
	end := time.Since(start)
	fmt.Printf("MoveGen (1000): %s\n", end)
}
