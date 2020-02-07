package board

import (
	"fmt"
	"math/bits"
)

func (board *Board) getCheckers(king uint64) uint64 {
	var checkers uint64

	kingIdx := bits.TrailingZeros64(king)

	horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, board.stateBoards[Occupied])
	fmt.Println("Horizontal moves")
	DrawBitboard(horizontalMoves)
	fmt.Println("Enemepy pieces")
	DrawBitboard(board.stateBoards[EnemyPieces])
	DrawBitboard(horizontalMoves & board.stateBoards[EnemyPieces])
	checkers |= horizontalMoves & board.stateBoards[EnemyPieces]
	diagonalMoves := board.DiagonalAndAntiDiagonalMoves(kingIdx, board.stateBoards[Occupied])
	fmt.Println("diagonalMoves moves")
	DrawBitboard(diagonalMoves)
	fmt.Println("Enemepy pieces")
	DrawBitboard(board.stateBoards[EnemyPieces])
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
		fmt.Println("More than 1 checkers")
		DrawBitboard(checkers)
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

	board.possibleWhitePawn(moveList, pushMask, captureMask)
	// board.possibleKnightMoves(moveList, board.bitboards[WN], WN)
	// board.possibleBishopMoves(moveList, board.bitboards[WB], WB)
	// board.possibleRookMoves(moveList, board.bitboards[WR], WR)
	// board.possibleQueenMoves(moveList, board.bitboards[WQ], WQ)
	// board.possibleKingMoves(moveList, board.bitboards[WK], WK)
	// board.possibleCastleWhite(moveList)
}

func (board *Board) possibleWhitePawn(moveList *MoveList, pushMask, captureMask uint64) {
	wp := board.bitboards[WP]
	enemyPieces := board.stateBoards[EnemyPieces]
	empty := board.stateBoards[Empty]
	var possibility uint64 // holds one potential capture at a time
	var index int          // index of the "possibility" capture
	var pawnMoves uint64

	// Move all bits from the white pawn bitboard to the left by 7 squares
	// and disable File A (we do not want any leftovers from the move to the left)
	// also remove any captures on the 8th rank since that will be handled by promotions
	// finally AND the resulting bitboard with the black pieces i.e. a capture move is any
	// move that can capture an enemy piece

	//    original                 WP >> 7          WP & ~FileA & ~Rank8    WP & BlackPieces (currently no black pieces on the 3rd rank to capture)
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]
	// [               ]      [  X X X X X X X]      [  X X X X X X X]      [               ]
	// [X X X X X X X X]      [X              ]      [               ]      [               ]
	// [               ]      [               ]      [               ]      [               ]

	// captures and moves formward forward

	// capture right
	pawnMoves = (wp >> 7) & enemyPieces & (^Rank8) & (^FileA)
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+7, index, WP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// capture left
	pawnMoves = (wp >> 9) & enemyPieces & (^Rank8) & (^FileH)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+9, index, WP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// move 1 square forward
	pawnMoves = (wp >> 8) & empty & (^Rank8)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & pushMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+8, index, WP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (wp >> 16) & empty & (empty >> 8) & Rank4
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & pushMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+16, index, WP, 0, MoveFlagPawnStart))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	// pawn promotion by capture right
	pawnMoves = (wp >> 7) & enemyPieces & Rank8 & (^FileA)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index+7, index, WP, WQ, NoFlag))
		moveList.AddMove(GetMoveInt(index+7, index, WP, WR, NoFlag))
		moveList.AddMove(GetMoveInt(index+7, index, WP, WB, NoFlag))
		moveList.AddMove(GetMoveInt(index+7, index, WP, WN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by capture left
	pawnMoves = (wp >> 9) & enemyPieces & Rank8 & (^FileH)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index+9, index, WP, WQ, NoFlag))
		moveList.AddMove(GetMoveInt(index+9, index, WP, WR, NoFlag))
		moveList.AddMove(GetMoveInt(index+9, index, WP, WB, NoFlag))
		moveList.AddMove(GetMoveInt(index+9, index, WP, WN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (wp >> 8) & empty & Rank8 // pawn promotion by move 1 forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & pushMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+8, index, WP, WQ, NoFlag))
		moveList.AddMove(GetMoveInt(index+8, index, WP, WR, NoFlag))
		moveList.AddMove(GetMoveInt(index+8, index, WP, WB, NoFlag))
		moveList.AddMove(GetMoveInt(index+8, index, WP, WN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// En passant right
	possibility = (wp << 1) & board.bitboards[BP] & Rank5 & (^FileA) & board.bitboards[EP]
	// en passant is possible if the piece to be captured is in the capture mask
	// or the destination to where out pawn will move during the capture is in the push mask
	if possibility != 0 && ((possibility & captureMask) != 0 || (possibility >> 8 & pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-1, index-8, WP, 0, MoveFlagEnPass))
	}
	// en passant left
	possibility = (wp >> 1) & board.bitboards[BP] & Rank5 & (^FileH) & board.bitboards[EP]
	if possibility != 0 && ((possibility & captureMask) != 0 || (possibility >> 8 & pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)
		fmt.Println(index)
		moveList.AddMove(GetMoveInt(index+1, index-8, WP, 0, MoveFlagEnPass))
	}
}

func (board *Board) possibleBlackPawn(moveList *MoveList, pushMask, captureMask uint64) {
	bp := board.bitboards[BP]
	enemyPieces := board.stateBoards[EnemyPieces]
	empty := board.stateBoards[Empty]
	var possibility uint64 // holds one potential capture at a time
	var index int          // index of the "possibility" capture
	var pawnMoves uint64

	// captures and moves formward forward:
	// capture right
	pawnMoves = (bp << 7) & enemyPieces & (^Rank1) & (^FileH)
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-7, index, BP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// capture left
	pawnMoves = (bp << 9) & enemyPieces & (^Rank1) & (^FileA)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-9, index, BP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// move 1 square forward
	pawnMoves = (bp << 8) & empty & (^Rank1)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & pushMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-8, index, BP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (bp << 16) & empty & (empty << 8) & Rank5
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & pushMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-16, index, BP, 0, MoveFlagPawnStart))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	// pawn promotion by capture right
	pawnMoves = (bp << 7) & enemyPieces & Rank1 & (^FileH)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index-7, index, BP, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-7, index, BP, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-7, index, BP, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-7, index, BP, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by capture left
	pawnMoves = (bp << 9) & enemyPieces & Rank1 & (^FileA)
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & captureMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index-9, index, BP, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, BP, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, BP, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, BP, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	pawnMoves = (bp << 8) & empty & Rank1 // pawn promotion by move 1 forward
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 && (possibility & pushMask) != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-8, index, BP, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-8, index, BP, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-8, index, BP, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-8, index, BP, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// En passant right
	possibility = (bp >> 1) & board.bitboards[WP] & Rank4 & (^FileH) & board.bitboards[EP]
	if possibility != 0 && ((possibility & captureMask) != 0 || (possibility << 8 & pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+1, index+8, BP, 0, MoveFlagEnPass))
	}
	// en passant left
	possibility = (bp << 1) & board.bitboards[WP] & Rank4 & (^FileA) & board.bitboards[EP]
	if possibility != 0 && ((possibility & captureMask) != 0 || (possibility << 8 & pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-1, index+8, BP, 0, MoveFlagEnPass))
	}
}
