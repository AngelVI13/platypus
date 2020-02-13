package board

import (
	"fmt"
	"math/bits"
)

func (board *Board) getCheckers(king uint64) uint64 {
	var checkers uint64

	kingIdx := bits.TrailingZeros64(king)

	horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, board.stateBoards[Occupied])
	checkers |= horizontalMoves & board.stateBoards[EnemyRooksQueens]
	diagonalMoves := board.DiagonalAndAntiDiagonalMoves(kingIdx, board.stateBoards[Occupied])
	checkers |= diagonalMoves & board.stateBoards[EnemyBishopsQueens]

	// check if pawns are attacking the king
	var direction int
	if board.Side == White {
		// there should be black pawns above to the left and right of the king
		direction = -1
	} else {
		// there should be white pawns below to the left and right of the king
		direction = 1
	}

	var sidePawnPossibility uint64
	if kingIdx+direction*7 > 0 {
		sidePawnPossibility = (1 << (kingIdx + direction*7))
		if sidePawnPossibility&board.stateBoards[EnemyPawns] != 0 {
			checkers |= sidePawnPossibility
		}
	}
	if kingIdx+direction*9 > 0 {
		sidePawnPossibility = (1 << (kingIdx + direction*9))
		if sidePawnPossibility&board.stateBoards[EnemyPawns] != 0 {
			checkers |= sidePawnPossibility
		}
	}

	// knight moves
	possibility := KnightMoves[kingIdx]
	possibility &= board.stateBoards[EnemyKnights]
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
	for (rays&checkerBitboard == 0) && (newSquare-8) >= 0 {
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
	for (rays&checkerBitboard == 0) && (newSquare-9) >= 0 {
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
	for (rays&checkerBitboard == 0) && (newSquare-7) >= 0 {
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

// getPinnedPieceRays Get diagonal and horizontal rays that represent a pinned piece and its only available moves.
// Ray does not include king
func (board *Board) getPinnedPieceRays(kingBitboard uint64, pinRays *PinRays) {
	var ray uint64
	var numPinnedPieces int

	kingIdx := bits.TrailingZeros64(kingBitboard)
	enemyRooksQueens := board.stateBoards[EnemyRooksQueens]
	enemyBishopsQueens := board.stateBoards[EnemyBishopsQueens]
	myPieces := board.stateBoards[MyPieces]

	ray = 0
	newSquare := kingIdx
	numPinnedPieces = 0
	// generate file ray to the right
	// if there are more than 1 piece between my king and an enemy slider
	// -> not a pinned piece
	for (ray&enemyRooksQueens == 0) && (newSquare+1)%8 != 0 && numPinnedPieces <= 1 {
		newSquare++
		ray |= (1 << newSquare)
		if ray&enemyRooksQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate file ray to the left
	for (ray&enemyRooksQueens == 0) && (newSquare)%8 != 0 && numPinnedPieces <= 1 {
		newSquare--
		ray |= (1 << newSquare)
		if ray&enemyRooksQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate rank ray upwards
	for (ray&enemyRooksQueens == 0) && (newSquare-8) >= 0 && numPinnedPieces <= 1 {
		newSquare -= 8
		ray |= (1 << newSquare)
		if ray&enemyRooksQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate rank ray down
	for (ray&enemyRooksQueens == 0) && (newSquare+8) < 64 && numPinnedPieces <= 1 {
		newSquare += 8
		ray |= (1 << newSquare)
		if ray&enemyRooksQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate left diagonal down
	for (ray&enemyBishopsQueens == 0) && (newSquare+9) < 64 && numPinnedPieces <= 1 {
		newSquare += 9
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate left diagonal up
	for (ray&enemyBishopsQueens == 0) && (newSquare-9) >= 0 && numPinnedPieces <= 1 {
		newSquare -= 9
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate right diagonal down
	for (ray&enemyBishopsQueens == 0) && (newSquare+7) < 64 && numPinnedPieces <= 1 {
		newSquare += 7
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate left diagonal up
	for (ray&enemyBishopsQueens == 0) && (newSquare-7) >= 0 && numPinnedPieces <= 1 {
		newSquare -= 7
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&myPieces != 0 {
			numPinnedPieces++
			myPieces ^= ray & myPieces // remove identified piece from myPieces
		}
	}
}

// LegalMovesWhite Generates all legal moves for white
func (board *Board) LegalMovesWhite(moveList *MoveList) {
	board.UpdateBitMasks()

	board.possibleKingMoves(moveList, board.bitboards[WK], WK)

	checkers := board.getCheckers(board.bitboards[WK])

	// captureMask & pushMask represents all squares where
	// a piece can capture on or move to respectively
	// By default those masks allow captures/moves on all squares
	var captureMask uint64 = ^uint64(0)
	var pushMask uint64 = ^uint64(0)
	var pinRays PinRays

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

	fmt.Println("1 or 0 checkers")
	DrawBitboard(checkers)
	DrawBitboard(captureMask)
	DrawBitboard(pushMask)
	board.getPinnedPieceRays(board.bitboards[WK], &pinRays)

	board.possibleWhitePawn(moveList, pushMask, captureMask, &pinRays)
	board.possibleKnightMoves(moveList, board.bitboards[WN], WN, pushMask, captureMask, &pinRays)
	// board.possibleBishopMoves(moveList, board.bitboards[WB], WB)
	// board.possibleRookMoves(moveList, board.bitboards[WR], WR)
	// board.possibleQueenMoves(moveList, board.bitboards[WQ], WQ)
	// board.possibleKingMoves(moveList, board.bitboards[WK], WK)
	// board.possibleCastleWhite(moveList)
}

func (board *Board) possibleWhitePawn(moveList *MoveList, pushMask, captureMask uint64, pinRays *PinRays) {
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
	pawnMoves = (wp >> 7) & enemyPieces & (^Rank8) & (^FileA) & captureMask
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 7))
		if possibility&ray != 0 {
			// add move only if the pawn is not pinned
			moveList.AddMove(GetMoveInt(index+7, index, WP, 0, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// capture left
	pawnMoves = (wp >> 9) & enemyPieces & (^Rank8) & (^FileH) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 9))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index+9, index, WP, 0, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// move 1 square forward
	pawnMoves = (wp >> 8) & empty & (^Rank8) & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 8))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index+8, index, WP, 0, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (wp >> 16) & empty & (empty >> 8) & Rank4 & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 16))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index+16, index, WP, 0, MoveFlagPawnStart))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	// pawn promotion by capture right
	pawnMoves = (wp >> 7) & enemyPieces & Rank8 & (^FileA) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 7))
		if possibility&ray != 0 {
			// todo maybe Capture flag??
			moveList.AddMove(GetMoveInt(index+7, index, WP, WQ, NoFlag))
			moveList.AddMove(GetMoveInt(index+7, index, WP, WR, NoFlag))
			moveList.AddMove(GetMoveInt(index+7, index, WP, WB, NoFlag))
			moveList.AddMove(GetMoveInt(index+7, index, WP, WN, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by capture left
	pawnMoves = (wp >> 9) & enemyPieces & Rank8 & (^FileH) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 9))
		if possibility&ray != 0 {
			// todo maybe Capture flag??
			moveList.AddMove(GetMoveInt(index+9, index, WP, WQ, NoFlag))
			moveList.AddMove(GetMoveInt(index+9, index, WP, WR, NoFlag))
			moveList.AddMove(GetMoveInt(index+9, index, WP, WB, NoFlag))
			moveList.AddMove(GetMoveInt(index+9, index, WP, WN, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by move 1 forward
	pawnMoves = (wp >> 8) & empty & Rank8 & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index + 8))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index+8, index, WP, WQ, NoFlag))
			moveList.AddMove(GetMoveInt(index+8, index, WP, WR, NoFlag))
			moveList.AddMove(GetMoveInt(index+8, index, WP, WB, NoFlag))
			moveList.AddMove(GetMoveInt(index+8, index, WP, WN, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// En passant right
	possibility = (wp << 1) & board.bitboards[BP] & Rank5 & (^FileA) & board.bitboards[EP]
	// en passant is possible if the piece to be captured is in the capture mask
	// or the destination to where out pawn will move during the capture is in the push mask
	if possibility != 0 && ((possibility&captureMask) != 0 || (possibility>>8&pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)

		// Remove the capturing and captured pawn from the board and check
		// if the king is attacked by a rook or queen i.e this en passant capture is
		// illegal: example - 8/8/8/K2pP2q/8/8/8/3k4 w - d6 0 2
		occupied := board.stateBoards[Occupied]
		occupied ^= (1 << (index - 1)) | (1 << index)
		kingIdx := bits.TrailingZeros64(board.bitboards[WK])
		horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, occupied)
		checkers := horizontalMoves & board.stateBoards[EnemyRooksQueens]

		if checkers == 0 {
			moveList.AddMove(GetMoveInt(index-1, index-8, WP, 0, MoveFlagEnPass))
		}
	}
	// en passant left
	possibility = (wp >> 1) & board.bitboards[BP] & Rank5 & (^FileH) & board.bitboards[EP]
	if possibility != 0 && ((possibility&captureMask) != 0 || (possibility>>8&pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)

		occupied := board.stateBoards[Occupied]
		occupied ^= (1 << (index + 1)) | (1 << index)
		kingIdx := bits.TrailingZeros64(board.bitboards[WK])
		horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, occupied)
		checkers := horizontalMoves & board.stateBoards[EnemyRooksQueens]

		if checkers == 0 {
			moveList.AddMove(GetMoveInt(index+1, index-8, WP, 0, MoveFlagEnPass))
		}
	}
}

func (board *Board) possibleBlackPawn(moveList *MoveList, pushMask, captureMask uint64, pinRays *PinRays) {
	bp := board.bitboards[BP]
	enemyPieces := board.stateBoards[EnemyPieces]
	empty := board.stateBoards[Empty]
	var possibility uint64 // holds one potential capture at a time
	var index int          // index of the "possibility" capture
	var pawnMoves uint64

	// captures and moves formward forward:
	// capture right
	pawnMoves = (bp << 7) & enemyPieces & (^Rank1) & (^FileH) & captureMask
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-7, index, BP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// capture left
	pawnMoves = (bp << 9) & enemyPieces & (^Rank1) & (^FileA) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-9, index, BP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// move 1 square forward
	pawnMoves = (bp << 8) & empty & (^Rank1) & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-8, index, BP, 0, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
	// also check that resulting square is on rank4.
	// (instead of check I mean eliminate squares that do not comply with these conditions)
	pawnMoves = (bp << 16) & empty & (empty << 8) & Rank5 & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-16, index, BP, 0, MoveFlagPawnStart))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	// pawn promotion by capture right
	pawnMoves = (bp << 7) & enemyPieces & Rank1 & (^FileH) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
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
	pawnMoves = (bp << 9) & enemyPieces & Rank1 & (^FileA) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		// todo maybe Capture flag??
		moveList.AddMove(GetMoveInt(index-9, index, BP, BQ, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, BP, BR, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, BP, BB, NoFlag))
		moveList.AddMove(GetMoveInt(index-9, index, BP, BN, NoFlag))
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by move 1 forward
	pawnMoves = (bp << 8) & empty & Rank1 & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
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
	if possibility != 0 && ((possibility&captureMask) != 0 || (possibility<<8&pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index+1, index+8, BP, 0, MoveFlagEnPass))
	}
	// en passant left
	possibility = (bp << 1) & board.bitboards[WP] & Rank4 & (^FileA) & board.bitboards[EP]
	if possibility != 0 && ((possibility&captureMask) != 0 || (possibility<<8&pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)
		moveList.AddMove(GetMoveInt(index-1, index+8, BP, 0, MoveFlagEnPass))
	}
}

func (board *Board) possibleKnightMoves(moveList *MoveList, knight uint64, pieceType int, pushMask, captureMask uint64, pinRays *PinRays) {
	// Choose bishop
	knightPossibility := knight & (^(knight - 1))
	var possibility uint64

	for knightPossibility != 0 {
		// Current knight index (in bitmask)
		knightIdx := bits.TrailingZeros64(knightPossibility)
		// if piece is pinned limits possibilities to move only along the pin line
		pinRay := pinRays.GetRay(knightPossibility)
		possibility = KnightMoves[knightIdx] & board.stateBoards[NotMyPieces] & (pushMask | captureMask) & pinRay
		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			moveList.AddMove(GetMoveInt(knightIdx, moveIndex, pieceType, 0, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		knight &= ^knightPossibility
		knightPossibility = knight & (^(knight - 1))
	}
}

// func (board *Board) possibleBishopMoves(moveList *MoveList, bishop uint64, pieceType int) {
// 	// Choose bishop
// 	bishopPossibility := bishop & (^(bishop - 1))
// 	var possibility uint64

// 	for bishopPossibility != 0 {
// 		// Current bishop index (in bitmask)
// 		bishopIdx := bits.TrailingZeros64(bishopPossibility)
// 		possibility = board.DiagonalAndAntiDiagonalMoves(bishopIdx) & NotMyPieces

// 		// choose move
// 		movePossibility := possibility & (^(possibility - 1))
// 		for movePossibility != 0 {
// 			// possible move index (in bitmask)
// 			moveIndex := bits.TrailingZeros64(movePossibility)
// 			moveList.AddMove(GetMoveInt(bishopIdx, moveIndex, pieceType, 0, NoFlag))
// 			possibility &= ^movePossibility                      // remove move from all possible moves
// 			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
// 		}
// 		bishop &= ^bishopPossibility
// 		bishopPossibility = bishop & (^(bishop - 1))
// 	}
// }

// func (board *Board) possibleRookMoves(moveList *MoveList, rook uint64, pieceType int) {
// 	// Choose rook
// 	rookPossibility := rook & (^(rook - 1))
// 	var possibility uint64

// 	for rookPossibility != 0 {
// 		// Current rook index (in bitmask)
// 		rookIdx := bits.TrailingZeros64(rookPossibility)
// 		possibility = board.HorizontalAndVerticalMoves(rookIdx) & NotMyPieces

// 		// choose move
// 		movePossibility := possibility & (^(possibility - 1))
// 		for movePossibility != 0 {
// 			// possible move index (in bitmask)
// 			moveIndex := bits.TrailingZeros64(movePossibility)
// 			moveList.AddMove(GetMoveInt(rookIdx, moveIndex, pieceType, 0, NoFlag))
// 			possibility &= ^movePossibility                      // remove move from all possible moves
// 			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
// 		}
// 		rook &= ^rookPossibility
// 		rookPossibility = rook & (^(rook - 1))
// 	}
// }

// func (board *Board) possibleQueenMoves(moveList *MoveList, queen uint64, pieceType int) {
// 	// Choose queen
// 	queenPossibility := queen & (^(queen - 1))
// 	var possibility uint64

// 	for queenPossibility != 0 {
// 		// Current queen index (in bitmask)
// 		queenIdx := bits.TrailingZeros64(queenPossibility)
// 		possibility = (board.HorizontalAndVerticalMoves(queenIdx) | board.DiagonalAndAntiDiagonalMoves(queenIdx)) & NotMyPieces

// 		// choose move
// 		movePossibility := possibility & (^(possibility - 1))
// 		for movePossibility != 0 {
// 			// possible move index (in bitmask)
// 			moveIndex := bits.TrailingZeros64(movePossibility)
// 			moveList.AddMove(GetMoveInt(queenIdx, moveIndex, pieceType, 0, NoFlag))
// 			possibility &= ^movePossibility                      // remove move from all possible moves
// 			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
// 		}
// 		queen &= ^queenPossibility
// 		queenPossibility = queen & (^(queen - 1))
// 	}
// }

func (board *Board) possibleKingMoves(moveList *MoveList, king uint64, pieceType int) {
	var possibility uint64

	// Current king index (in bitmask)
	kingIdx := bits.TrailingZeros64(king)

	possibility = KingMoves[kingIdx] & board.stateBoards[NotMyPieces] & ^board.stateBoards[Unsafe]

	// choose move
	movePossibility := possibility & (^(possibility - 1))
	for movePossibility != 0 {
		// possible move index (in bitmask)
		moveIndex := bits.TrailingZeros64(movePossibility)
		moveList.AddMove(GetMoveInt(kingIdx, moveIndex, pieceType, 0, NoFlag))
		possibility &= ^movePossibility                      // remove move from all possible moves
		movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
	}
}

// func (board *Board) possibleCastleWhite(moveList *MoveList) {
// 	kingIdx := bits.TrailingZeros64(board.bitboards[WK])
// 	var queenSideSqBitboard uint64 = 1<<(kingIdx-1) | 1<<(kingIdx-2)
// 	var kingSideSqBitboard uint64 = 1<<(kingIdx+1) | 1<<(kingIdx+2)

// 	// todo replace hardcoded squares with consts
// 	if (board.castlePermissions&WhiteKingCastling) != 0 && (bits.OnesCount64((kingSideSqBitboard|board.bitboards[WK]) & ^Unsafe) == 3) && (bits.OnesCount64(Empty&kingSideSqBitboard) == 2) {
// 		moveList.AddMove(GetMoveInt(60, 62, WK, 0, MoveFlagCastle))
// 	}
// 	if (board.castlePermissions&WhiteQueenCastling) != 0 && (bits.OnesCount64((queenSideSqBitboard|board.bitboards[WK]) & ^Unsafe) == 3) && (bits.OnesCount64(Empty&(queenSideSqBitboard|1<<(kingIdx-3))) == 3) { // on the queen side there are 3 sq that should be empty to enable castling
// 		moveList.AddMove(GetMoveInt(60, 58, WK, 0, MoveFlagCastle))
// 	}
// }

// func (board *Board) possibleCastleBlack(moveList *MoveList) {
// 	kingIdx := bits.TrailingZeros64(board.bitboards[BK])
// 	var queenSideSqBitboard uint64 = 1<<(kingIdx-1) | 1<<(kingIdx-2)
// 	var kingSideSqBitboard uint64 = 1<<(kingIdx+1) | 1<<(kingIdx+2)

// 	if (board.castlePermissions&BlackKingCastling) != 0 && (bits.OnesCount64((kingSideSqBitboard|board.bitboards[BK]) & ^Unsafe) == 3) && (bits.OnesCount64(Empty&kingSideSqBitboard) == 2) {
// 		moveList.AddMove(GetMoveInt(4, 6, BK, 0, MoveFlagCastle))
// 	}
// 	if (board.castlePermissions&BlackQueenCastling) != 0 && (bits.OnesCount64((queenSideSqBitboard|board.bitboards[BK]) & ^Unsafe) == 3) && (bits.OnesCount64(Empty&(queenSideSqBitboard|1<<(kingIdx-3))) == 3) {
// 		moveList.AddMove(GetMoveInt(4, 2, BK, 0, MoveFlagCastle))
// 	}
// }
