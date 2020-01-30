package board

import (
	"fmt"
	"math/bits"
)

func (board *Board) getCheckers(king uint64) uint64 {
	var checkers uint64

	kingIdx := bits.TrailingZeros64(king)

	horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, board.stateBoards[Occupied])
	checkers |= horizontalMoves & board.stateBoards[EnemyPieces]
	diagonalMoves := board.DiagonalAndAntiDiagonalMoves(kingIdx, board.stateBoards[Occupied])
	checkers |= diagonalMoves & board.stateBoards[EnemyPieces]

	// knight moves
	possibility := KnightMoves[kingIdx]
	possibility &= board.stateBoards[EnemyPieces]
	checkers |= possibility

	return checkers
}

func (board *Board) LegalMovesWhite(moveList *MoveList) {
	board.UpdateBitMasks()

	board.possibleKingMoves(moveList, board.bitboards[WK], WK)

	checkers := board.getCheckers(board.bitboards[WK])

	if bits.OnesCount64(checkers) > 1 {
		// if there are more than 1 checking piece -> only king moves are possible
		return
	}

	fmt.Println("One or 0 checkers")
	DrawBitboard(checkers)
}
