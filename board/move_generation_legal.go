package board

import (
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

func (board *Board) LegalMovesWhite() {
	board.UpdateBitMasks()

	// todo first generate all king moves!

	checkers := board.getCheckers(board.bitboards[WK])

	// todo then if checkers bit count is > 1 -> simply return the generated king moves since these are the only moves we can make

	DrawBitboard(checkers)
}
