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

// getCheckerSliderRaysToKing Get diagonal OR horizontal rays to a square. Ray does not include king or checker
func getCheckerSliderRaysToKing(kingBitboard uint64, checkerBitboard uint64) uint64 {
	var rays uint64
	kingIdx := bits.TrailingZeros64(kingBitboard)

	rays = 0
	newSquare := kingIdx
	// generate file ray to the right
	for (rays&checkerBitboard == 0) && (newSquare+1)%8 != 0 {
		newSquare++
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate file ray to the left
	for (rays&checkerBitboard == 0) && (newSquare)%8 != 0 {
		newSquare--
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate rank ray upwards
	for (rays&checkerBitboard == 0) && (newSquare-8) > 0 {
		newSquare -= 8
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate rank ray down
	for (rays&checkerBitboard == 0) && (newSquare+8) < 64 {
		newSquare += 8
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate left diagonal down
	for (rays&checkerBitboard == 0) && (newSquare+9) < 64 {
		newSquare += 9
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate left diagonal up
	for (rays&checkerBitboard == 0) && (newSquare-9) > 0 {
		newSquare -= 9
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate right diagonal down
	for (rays&checkerBitboard == 0) && (newSquare+7) < 64 {
		newSquare += 7
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	rays = 0
	newSquare = kingIdx
	// generate left diagonal up
	for (rays&checkerBitboard == 0) && (newSquare-7) > 0 {
		newSquare -= 7
		rays |= (1 << newSquare)
		if rays&checkerBitboard != 0 {
			rays ^= checkerBitboard // remove checker square from ray and return rays
			return rays
		}
	}

	// if we haven't returned by now -> couldn't generate rays
	panic("Could not generate rays.")
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
				pushMask = getCheckerSliderRaysToKing(board.bitboards[WK], checkers)
				break
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
