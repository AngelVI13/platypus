package main

import (
	"fmt"
	"platypus/board"
)

func main() {
	boardVar := board.Board{}
	fmt.Println(&boardVar)
}
