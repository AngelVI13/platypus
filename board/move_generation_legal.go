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

	// captureMask & pushMask represents all squares where
	// a piece can capture on or move to respectively
	// By default those masks allow captures/moves on all squares
	var captureMask uint64 = ^uint64(0)
	var pushMask uint64 = ^uint64(0)

	checkersNum := bits.OnesCount64(checkers)
	if checkersNum > 1 {
		// if there are more than 1 checking piece -> only king moves are possible
		return
	} else if checkersNum == 1 {
		// if only 1 checker, we can evade check by capturing the checking piece
		captureMask = checkers

		// iterate over bitboards to find out what piece type is the checking piece
		for pieceType, bitboard := range board.bitboards {
			if bitboard&checkers != 0 && IsSlider[pieceType] {
				fmt.Println(pieceType)
				// the push mask is limited to squares between the king and the piece giving check
				pushMask = 0 // todo calculate rays
			}
		}
		// if we are not attacked by a sliding piece (i.e a knight) then
		// the only way to escape is to capture the checking piece or move out of check
		if pushMask == ^uint64(0) {
			pushMask = 0
		}
	}

	fmt.Println("One or 0 checkers")
	DrawBitboard(checkers)
	DrawBitboard(captureMask)
	DrawBitboard(pushMask)
}
