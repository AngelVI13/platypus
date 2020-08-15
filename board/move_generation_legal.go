package board

import (
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
// Ray does not include the pinned king square
func (board *Board) getPinnedPieceRays(kingBitboard uint64, pinRays *PinRays) {
	var ray uint64
	var numPinnedPieces int

	kingIdx := bits.TrailingZeros64(kingBitboard)
	enemyRooksQueens := board.stateBoards[EnemyRooksQueens]
	enemyBishopsQueens := board.stateBoards[EnemyBishopsQueens]
	enemySide := board.Side ^ 1
	// enemy kings, knights and pawns are always pin blocking pieces
	blockingPieces := board.stateBoards[MyPieces] | board.stateBoards[EnemyKnights] | board.stateBoards[EnemyPawns] | board.bitboards[enemySide*6+WK]
	enemyBishops := board.bitboards[enemySide*6+WB]
	enemyRooks := board.bitboards[enemySide*6+WR]

	ray = 0
	newSquare := kingIdx
	numPinnedPieces = 0
	// when calculating horizontal & vertical pins - enemy bishops are blocking pieces
	blockingPieces ^= enemyBishops
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

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
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

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
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

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
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

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// when calculating diagonal & antidiagonal pins - enemy rooks are blocking pieces, bishops are not
	blockingPieces ^= enemyBishops // remove enemy bishops
	blockingPieces ^= enemyRooks   // add enemy rooks
	// generate left diagonal down
	// (newSquare%8+1 == (newSquare+9)%8) condition makes sure the diagonal does not wrap around
	for (ray&enemyBishopsQueens == 0) && (newSquare+9) < 64 && (newSquare%8+1 == (newSquare+9)%8) && numPinnedPieces <= 1 {
		newSquare += 9
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate left diagonal up
	for (ray&enemyBishopsQueens == 0) && (newSquare-9) >= 0 && (newSquare%8-1 == (newSquare-9)%8) && numPinnedPieces <= 1 {
		newSquare -= 9
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate right diagonal down
	for (ray&enemyBishopsQueens == 0) && (newSquare+7) < 64 && (newSquare%8-1 == (newSquare+7)%8) && numPinnedPieces <= 1 {
		newSquare += 7
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
		}
	}

	ray = 0
	newSquare = kingIdx
	numPinnedPieces = 0
	// generate left diagonal up
	for (ray&enemyBishopsQueens == 0) && (newSquare-7) >= 0 && (newSquare%8+1 == (newSquare-7)%8) && numPinnedPieces <= 1 {
		newSquare -= 7
		ray |= (1 << newSquare)
		if ray&enemyBishopsQueens != 0 && numPinnedPieces == 1 {
			pinRays.AddRay(ray)
			break
		}

		if ray&blockingPieces != 0 {
			numPinnedPieces++
			blockingPieces ^= ray & blockingPieces // remove identified piece from myPieces
		}
	}
}

// LegalMovesWhite Generates all legal moves for white
// todo unify legal moves white and black into 1 method
func (board *Board) LegalMovesWhite(moveList *MoveList) {
	board.UpdateBitMasks()

	board.possibleKingMoves(moveList, board.bitboards[WK])

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
		// fmt.Println("More than 1 checkers")
		// DrawBitboard(checkers)
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
	} else {
		board.possibleCastleWhite(moveList)
	}

	board.getPinnedPieceRays(board.bitboards[WK], &pinRays)

	// for i := 0; i < pinRays.Count; i++ {
	// 	DrawBitboard(pinRays.Rays[i])
	// }
	board.possibleWhitePawn(moveList, pushMask, captureMask, &pinRays)
	board.possibleKnightMoves(moveList, board.bitboards[WN], pushMask, captureMask, &pinRays)
	board.possibleBishopMoves(moveList, board.bitboards[WB], pushMask, captureMask, &pinRays)
	board.possibleRookMoves(moveList, board.bitboards[WR], pushMask, captureMask, &pinRays)
	board.possibleQueenMoves(moveList, board.bitboards[WQ], pushMask, captureMask, &pinRays)
}

// LegalMovesBlack Generates all legal moves for black
func (board *Board) LegalMovesBlack(moveList *MoveList) {
	board.UpdateBitMasks()

	board.possibleKingMoves(moveList, board.bitboards[BK])

	checkers := board.getCheckers(board.bitboards[BK])

	// captureMask & pushMask represents all squares where
	// a piece can capture on or move to respectively
	// By default those masks allow captures/moves on all squares
	var captureMask uint64 = ^uint64(0)
	var pushMask uint64 = ^uint64(0)
	var pinRays PinRays

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
				// the push mask is limited to squares between the king and the piece giving check
				pushMask = getCheckerSliderRaysToKing(board.bitboards[BK], checkers)
				break
			}
		}
		// if we are not attacked by a sliding piece (i.e a knight) then
		// the only way to escape is to capture the checking piece or move out of check
		if pushMask == ^uint64(0) {
			pushMask = 0
		}
	} else {
		board.possibleCastleBlack(moveList)
	}

	board.getPinnedPieceRays(board.bitboards[BK], &pinRays)

	board.possibleBlackPawn(moveList, pushMask, captureMask, &pinRays)
	board.possibleKnightMoves(moveList, board.bitboards[BN], pushMask, captureMask, &pinRays)
	board.possibleBishopMoves(moveList, board.bitboards[BB], pushMask, captureMask, &pinRays)
	board.possibleRookMoves(moveList, board.bitboards[BR], pushMask, captureMask, &pinRays)
	board.possibleQueenMoves(moveList, board.bitboards[BQ], pushMask, captureMask, &pinRays)
}

func (board *Board) possibleWhitePawn(moveList *MoveList, pushMask, captureMask uint64, pinRays *PinRays) {
	wp := board.bitboards[WP]
	enemyPieces := board.stateBoards[EnemyPieces]
	empty := board.stateBoards[Empty]
	var possibility uint64 // holds one potential capture at a time
	var index int          // index of the "possibility" capture
	var pawnMoves uint64
	var capturedPiece int

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
		// add move only if the pawn is not pinned
		if possibility&ray != 0 {
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index+7, index, capturedPiece, NoPiece, NoFlag))
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
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index+9, index, capturedPiece, NoPiece, NoFlag))
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
			moveList.AddMove(GetMoveInt(index+8, index, NoPiece, NoPiece, NoFlag))
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
			moveList.AddMove(GetMoveInt(index+16, index, NoPiece, NoPiece, MoveFlagPawnStart))
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
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index+7, index, capturedPiece, WQ, NoFlag))
			moveList.AddMove(GetMoveInt(index+7, index, capturedPiece, WR, NoFlag))
			moveList.AddMove(GetMoveInt(index+7, index, capturedPiece, WB, NoFlag))
			moveList.AddMove(GetMoveInt(index+7, index, capturedPiece, WN, NoFlag))
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
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index+9, index, capturedPiece, WQ, NoFlag))
			moveList.AddMove(GetMoveInt(index+9, index, capturedPiece, WR, NoFlag))
			moveList.AddMove(GetMoveInt(index+9, index, capturedPiece, WB, NoFlag))
			moveList.AddMove(GetMoveInt(index+9, index, capturedPiece, WN, NoFlag))
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
			moveList.AddMove(GetMoveInt(index+8, index, NoPiece, WQ, NoFlag))
			moveList.AddMove(GetMoveInt(index+8, index, NoPiece, WR, NoFlag))
			moveList.AddMove(GetMoveInt(index+8, index, NoPiece, WB, NoFlag))
			moveList.AddMove(GetMoveInt(index+8, index, NoPiece, WN, NoFlag))
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
			moveList.AddMove(GetMoveInt(index-1, index-8, BP, NoPiece, MoveFlagEnPass))
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
			moveList.AddMove(GetMoveInt(index+1, index-8, BP, NoPiece, MoveFlagEnPass))
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
	var capturedPiece int

	// captures and moves formward forward:
	// capture right
	pawnMoves = (bp << 7) & enemyPieces & (^Rank1) & (^FileH) & captureMask
	// Find first bit which is equal to '1' i.e. first capture
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index - 7))
		if possibility&ray != 0 {
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index-7, index, capturedPiece, NoPiece, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// capture left
	pawnMoves = (bp << 9) & enemyPieces & (^Rank1) & (^FileA) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index - 9))
		if possibility&ray != 0 {
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index-9, index, capturedPiece, NoPiece, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// move 1 square forward
	pawnMoves = (bp << 8) & empty & (^Rank1) & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index - 8))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index-8, index, NoPiece, NoPiece, NoFlag))
		}
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
		ray := pinRays.GetRay(1 << (index - 16))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index-16, index, NoPiece, NoPiece, MoveFlagPawnStart))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// promotions
	// pawn promotion by capture right
	pawnMoves = (bp << 7) & enemyPieces & Rank1 & (^FileH) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index - 7))
		if possibility&ray != 0 {
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index-7, index, capturedPiece, BQ, NoFlag))
			moveList.AddMove(GetMoveInt(index-7, index, capturedPiece, BR, NoFlag))
			moveList.AddMove(GetMoveInt(index-7, index, capturedPiece, BB, NoFlag))
			moveList.AddMove(GetMoveInt(index-7, index, capturedPiece, BN, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by capture left
	pawnMoves = (bp << 9) & enemyPieces & Rank1 & (^FileA) & captureMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index - 9))
		if possibility&ray != 0 {
			capturedPiece = board.position[index]
			moveList.AddMove(GetMoveInt(index-9, index, capturedPiece, BQ, NoFlag))
			moveList.AddMove(GetMoveInt(index-9, index, capturedPiece, BR, NoFlag))
			moveList.AddMove(GetMoveInt(index-9, index, capturedPiece, BB, NoFlag))
			moveList.AddMove(GetMoveInt(index-9, index, capturedPiece, BN, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// pawn promotion by move 1 forward
	pawnMoves = (bp << 8) & empty & Rank1 & pushMask
	possibility = pawnMoves & (^(pawnMoves - 1))
	for possibility != 0 {
		index = bits.TrailingZeros64(possibility)
		ray := pinRays.GetRay(1 << (index - 8))
		if possibility&ray != 0 {
			moveList.AddMove(GetMoveInt(index-8, index, NoPiece, BQ, NoFlag))
			moveList.AddMove(GetMoveInt(index-8, index, NoPiece, BR, NoFlag))
			moveList.AddMove(GetMoveInt(index-8, index, NoPiece, BB, NoFlag))
			moveList.AddMove(GetMoveInt(index-8, index, NoPiece, BN, NoFlag))
		}
		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
	}

	// En passant right
	possibility = (bp >> 1) & board.bitboards[WP] & Rank4 & (^FileH) & board.bitboards[EP]
	if possibility != 0 && ((possibility&captureMask) != 0 || (possibility<<8&pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)

		// Remove the capturing and captured pawn from the board and check
		// if the king is attacked by a rook or queen i.e this en passant capture is
		// illegal: example - 8/8/8/K2pP2q/8/8/8/3k4 w - d6 0 2
		occupied := board.stateBoards[Occupied]
		occupied ^= (1 << (index + 1)) | (1 << index)
		kingIdx := bits.TrailingZeros64(board.bitboards[BK])
		horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, occupied)
		checkers := horizontalMoves & board.stateBoards[EnemyRooksQueens]

		if checkers == 0 {
			moveList.AddMove(GetMoveInt(index+1, index+8, WP, NoPiece, MoveFlagEnPass))
		}
	}
	// en passant left
	possibility = (bp << 1) & board.bitboards[WP] & Rank4 & (^FileA) & board.bitboards[EP]
	if possibility != 0 && ((possibility&captureMask) != 0 || (possibility<<8&pushMask) != 0) {
		index = bits.TrailingZeros64(possibility)

		occupied := board.stateBoards[Occupied]
		occupied ^= (1 << (index - 1)) | (1 << index)
		kingIdx := bits.TrailingZeros64(board.bitboards[BK])
		horizontalMoves := board.HorizontalAndVerticalMoves(kingIdx, occupied)
		checkers := horizontalMoves & board.stateBoards[EnemyRooksQueens]

		if checkers == 0 {
			moveList.AddMove(GetMoveInt(index-1, index+8, WP, NoPiece, MoveFlagEnPass))
		}
	}
}

func (board *Board) possibleKnightMoves(moveList *MoveList, knight uint64, pushMask, captureMask uint64, pinRays *PinRays) {
	// Choose bishop
	knightPossibility := knight & (^(knight - 1))
	var possibility uint64
	var capturedPiece int

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
			capturedPiece = board.position[moveIndex]
			moveList.AddMove(GetMoveInt(knightIdx, moveIndex, capturedPiece, NoPiece, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		knight &= ^knightPossibility
		knightPossibility = knight & (^(knight - 1))
	}
}

func (board *Board) possibleBishopMoves(moveList *MoveList, bishop uint64, pushMask, captureMask uint64, pinRays *PinRays) {
	// Choose bishop
	bishopPossibility := bishop & (^(bishop - 1))
	var possibility uint64
	var capturedPiece int

	for bishopPossibility != 0 {
		// Current bishop index (in bitmask)
		bishopIdx := bits.TrailingZeros64(bishopPossibility)
		pinRay := pinRays.GetRay(bishopPossibility)
		possibility = board.DiagonalAndAntiDiagonalMoves(bishopIdx, board.stateBoards[Occupied])
		possibility &= board.stateBoards[NotMyPieces] & (pushMask | captureMask) & pinRay

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			capturedPiece = board.position[moveIndex]
			moveList.AddMove(GetMoveInt(bishopIdx, moveIndex, capturedPiece, NoPiece, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		bishop &= ^bishopPossibility
		bishopPossibility = bishop & (^(bishop - 1))
	}
}

func (board *Board) possibleRookMoves(moveList *MoveList, rook uint64, pushMask, captureMask uint64, pinRays *PinRays) {
	// Choose rook
	rookPossibility := rook & (^(rook - 1))
	var possibility uint64
	var capturedPiece int

	for rookPossibility != 0 {
		// Current rook index (in bitmask)
		rookIdx := bits.TrailingZeros64(rookPossibility)
		possibility = board.HorizontalAndVerticalMoves(rookIdx, board.stateBoards[Occupied])
		pinRay := pinRays.GetRay(rookPossibility)
		possibility &= board.stateBoards[NotMyPieces] & (pushMask | captureMask) & pinRay

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			capturedPiece = board.position[moveIndex]
			moveList.AddMove(GetMoveInt(rookIdx, moveIndex, capturedPiece, NoPiece, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		rook &= ^rookPossibility
		rookPossibility = rook & (^(rook - 1))
	}
}

func (board *Board) possibleQueenMoves(moveList *MoveList, queen uint64, pushMask, captureMask uint64, pinRays *PinRays) {
	// Choose queen
	queenPossibility := queen & (^(queen - 1))
	var possibility uint64
	var capturedPiece int

	for queenPossibility != 0 {
		// Current queen index (in bitmask)
		queenIdx := bits.TrailingZeros64(queenPossibility)
		pinRay := pinRays.GetRay(queenPossibility)
		possibility = board.HorizontalAndVerticalMoves(queenIdx, board.stateBoards[Occupied])
		possibility |= board.DiagonalAndAntiDiagonalMoves(queenIdx, board.stateBoards[Occupied])
		possibility &= board.stateBoards[NotMyPieces] & (pushMask | captureMask) & pinRay

		// choose move
		movePossibility := possibility & (^(possibility - 1))
		for movePossibility != 0 {
			// possible move index (in bitmask)
			moveIndex := bits.TrailingZeros64(movePossibility)
			capturedPiece = board.position[moveIndex]
			moveList.AddMove(GetMoveInt(queenIdx, moveIndex, capturedPiece, NoPiece, NoFlag))
			possibility &= ^movePossibility                      // remove move from all possible moves
			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
		}
		queen &= ^queenPossibility
		queenPossibility = queen & (^(queen - 1))
	}
}

func (board *Board) possibleKingMoves(moveList *MoveList, king uint64) {
	var possibility uint64
	var capturedPiece int

	// Current king index (in bitmask)
	kingIdx := bits.TrailingZeros64(king)

	possibility = KingMoves[kingIdx] & board.stateBoards[NotMyPieces] & ^board.stateBoards[Unsafe]

	// choose move
	movePossibility := possibility & (^(possibility - 1))
	for movePossibility != 0 {
		// possible move index (in bitmask)
		moveIndex := bits.TrailingZeros64(movePossibility)
		capturedPiece = board.position[moveIndex]
		moveList.AddMove(GetMoveInt(kingIdx, moveIndex, capturedPiece, NoPiece, NoFlag))
		possibility &= ^movePossibility                      // remove move from all possible moves
		movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
	}
}

func (board *Board) possibleCastleWhite(moveList *MoveList) {
	// if no castling is allowed -> return early
	if (board.castlePermissions&WhiteKingCastling) == 0 && 
		(board.castlePermissions&WhiteQueenCastling) == 0 {
		return
	}

	empty := board.stateBoards[Empty]
	unsafe := board.stateBoards[Unsafe]

	kingIdx := bits.TrailingZeros64(board.bitboards[WK])
	var queenSideSqBitboard uint64 = 1<<(kingIdx-1) | 1<<(kingIdx-2)
	var kingSideSqBitboard uint64 = 1<<(kingIdx+1) | 1<<(kingIdx+2)

	if (board.castlePermissions&WhiteKingCastling) != 0 && (bits.OnesCount64((kingSideSqBitboard|board.bitboards[WK]) & ^unsafe) == 3) && (bits.OnesCount64(empty&kingSideSqBitboard) == 2) {
		moveList.AddMove(GetMoveInt(E1, G1, NoPiece, NoPiece, MoveFlagCastle))
	}
	if (board.castlePermissions&WhiteQueenCastling) != 0 && (bits.OnesCount64((queenSideSqBitboard|board.bitboards[WK]) & ^unsafe) == 3) && (bits.OnesCount64(empty&(queenSideSqBitboard|1<<(kingIdx-3))) == 3) { // on the queen side there are 3 sq that should be empty to enable castling
		moveList.AddMove(GetMoveInt(E1, C1, NoPiece, NoPiece, MoveFlagCastle))
	}
}

func (board *Board) possibleCastleBlack(moveList *MoveList) {
	// if no castling is allowed -> return early
	if (board.castlePermissions&BlackKingCastling) == 0 && 
		(board.castlePermissions&BlackQueenCastling) == 0 {
		return
	}	

	empty := board.stateBoards[Empty]
	unsafe := board.stateBoards[Unsafe]

	kingIdx := bits.TrailingZeros64(board.bitboards[BK])

	var queenSideSqBitboard uint64 = 1<<(kingIdx-1) | 1<<(kingIdx-2)
	var kingSideSqBitboard uint64 = 1<<(kingIdx+1) | 1<<(kingIdx+2)

	if (board.castlePermissions&BlackKingCastling) != 0 && (bits.OnesCount64((kingSideSqBitboard|board.bitboards[BK]) & ^unsafe) == 3) && (bits.OnesCount64(empty&kingSideSqBitboard) == 2) {
		moveList.AddMove(GetMoveInt(E8, G8, NoPiece, NoPiece, MoveFlagCastle))
	}
	if (board.castlePermissions&BlackQueenCastling) != 0 && (bits.OnesCount64((queenSideSqBitboard|board.bitboards[BK]) & ^unsafe) == 3) && (bits.OnesCount64(empty&(queenSideSqBitboard|1<<(kingIdx-3))) == 3) {
		moveList.AddMove(GetMoveInt(E8, C8, NoPiece, NoPiece, MoveFlagCastle))
	}
}
