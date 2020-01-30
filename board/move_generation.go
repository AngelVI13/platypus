package board

import (
	"math/bits"
)

// // PossibleMovesWhite Compute all possible white moves and add them to the move list
// func (board *Board) PossibleMovesWhite(moveList *MoveList) {
// 	// This represents all squares which are not white pieces (including empty squares).
// 	// Black king is added in order to avoid generating capture moves on the black king.
// 	// For example pawn takes king is not a legal move
// 	NotMyPieces = ^(board.bitboards[WP] |
// 		board.bitboards[WN] |
// 		board.bitboards[WB] |
// 		board.bitboards[WR] |
// 		board.bitboards[WQ] |
// 		board.bitboards[WK] |
// 		board.bitboards[BK])

// 	EnemyPieces = (board.bitboards[BP] |
// 		board.bitboards[BN] |
// 		board.bitboards[BB] |
// 		board.bitboards[BR] |
// 		board.bitboards[BQ])

// 	Occupied = (board.bitboards[WP] |
// 		board.bitboards[WN] |
// 		board.bitboards[WB] |
// 		board.bitboards[WR] |
// 		board.bitboards[WQ] |
// 		board.bitboards[WK] |
// 		board.bitboards[BP] |
// 		board.bitboards[BN] |
// 		board.bitboards[BB] |
// 		board.bitboards[BR] |
// 		board.bitboards[BQ] |
// 		board.bitboards[BK])

// 	Empty = ^Occupied

// 	Unsafe = board.unsafeForWhite()

// 	board.possibleWhitePawn(moveList)
// 	board.possibleKnightMoves(moveList, board.bitboards[WN], WN)
// 	board.possibleBishopMoves(moveList, board.bitboards[WB], WB)
// 	board.possibleRookMoves(moveList, board.bitboards[WR], WR)
// 	board.possibleQueenMoves(moveList, board.bitboards[WQ], WQ)
// 	board.possibleKingMoves(moveList, board.bitboards[WK], WK)
// 	board.possibleCastleWhite(moveList)
// }

// // PossibleMovesBlack Compute all possible black moves and add them to the move list
// func (board *Board) PossibleMovesBlack(moveList *MoveList) {
// 	// This represents all squares which are not white pieces (including empty squares).
// 	// Black king is added in order to avoid generating capture moves on the black king.
// 	// For example pawn takes king is not a legal move
// 	NotMyPieces = ^(board.bitboards[BP] |
// 		board.bitboards[BN] |
// 		board.bitboards[BB] |
// 		board.bitboards[BR] |
// 		board.bitboards[BQ] |
// 		board.bitboards[BK] |
// 		board.bitboards[WK])

// 	EnemyPieces = (board.bitboards[WP] |
// 		board.bitboards[WN] |
// 		board.bitboards[WB] |
// 		board.bitboards[WR] |
// 		board.bitboards[WQ])

// 	Occupied = (board.bitboards[WP] |
// 		board.bitboards[WN] |
// 		board.bitboards[WB] |
// 		board.bitboards[WR] |
// 		board.bitboards[WQ] |
// 		board.bitboards[WK] |
// 		board.bitboards[BP] |
// 		board.bitboards[BN] |
// 		board.bitboards[BB] |
// 		board.bitboards[BR] |
// 		board.bitboards[BQ] |
// 		board.bitboards[BK])

// 	Empty = ^Occupied

// 	Unsafe = board.unsafeForBlack() // move this only to make move, castling shuold not check for checks

// 	board.possibleBlackPawn(moveList)
// 	board.possibleKnightMoves(moveList, board.bitboards[BN], BN)
// 	board.possibleBishopMoves(moveList, board.bitboards[BB], BB)
// 	board.possibleRookMoves(moveList, board.bitboards[BR], BR)
// 	board.possibleQueenMoves(moveList, board.bitboards[BQ], BQ)
// 	board.possibleKingMoves(moveList, board.bitboards[BK], BK)
// 	board.possibleCastleBlack(moveList)
// }

// HorizontalAndVerticalMoves Generate a bitboard of all possible horizontal and vertical moves for a given square
func (board *Board) HorizontalAndVerticalMoves(square int, occupied uint64) uint64 {
	var binarySquare uint64 = 1 << square
	fileMaskIdx := square % 8
	possibilitiesHorizontal := (occupied - 2*binarySquare) ^ bits.Reverse64(bits.Reverse64(occupied)-2*bits.Reverse64(binarySquare))
	possibilitiesVertical := ((occupied & FileMasks8[fileMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(
		bits.Reverse64(occupied&FileMasks8[fileMaskIdx])-(2*bits.Reverse64(binarySquare)))
	return (possibilitiesHorizontal & RankMasks8[square/8]) | (possibilitiesVertical & FileMasks8[fileMaskIdx])
}

// DiagonalAndAntiDiagonalMoves Generate a bitboard of all possible diagonal and anti-diagonal moves for a given square
func (board *Board) DiagonalAndAntiDiagonalMoves(square int, occupied uint64) uint64 {
	var binarySquare uint64 = 1 << square
	diagonalMaskIdx := (square / 8) + (square % 8)
	antiDiagonalMaskIdx := (square / 8) + 7 - (square % 8)
	possibilitiesDiagonal := ((occupied & DiagonalMasks8[diagonalMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(bits.Reverse64(occupied&DiagonalMasks8[diagonalMaskIdx])-(2*bits.Reverse64(binarySquare)))
	possibilitiesAntiDiagonal := ((occupied & AntiDiagonalMasks8[antiDiagonalMaskIdx]) - (2 * binarySquare)) ^ bits.Reverse64(bits.Reverse64(occupied&AntiDiagonalMasks8[antiDiagonalMaskIdx])-(2*bits.Reverse64(binarySquare)))
	return (possibilitiesDiagonal & DiagonalMasks8[diagonalMaskIdx]) | (possibilitiesAntiDiagonal & AntiDiagonalMasks8[antiDiagonalMaskIdx])
}

// func (board *Board) possibleWhitePawn(moveList *MoveList) {
// 	wp := board.bitboards[WP]
// 	var possibility uint64 // holds one potential capture at a time
// 	var index int          // index of the "possibility" capture
// 	var pawnMoves uint64

// 	// Move all bits from the white pawn bitboard to the left by 7 squares
// 	// and disable File A (we do not want any leftovers from the move to the left)
// 	// also remove any captures on the 8th rank since that will be handled by promotions
// 	// finally AND the resulting bitboard with the black pieces i.e. a capture move is any
// 	// move that can capture an enemy piece

// 	//    original                 WP >> 7          WP & ~FileA & ~Rank8    WP & BlackPieces (currently no black pieces on the 3rd rank to capture)
// 	// [               ]      [               ]      [               ]      [               ]
// 	// [               ]      [               ]      [               ]      [               ]
// 	// [               ]      [               ]      [               ]      [               ]
// 	// [               ]      [               ]      [               ]      [               ]
// 	// [               ]      [               ]      [               ]      [               ]
// 	// [               ]      [  X X X X X X X]      [  X X X X X X X]      [               ]
// 	// [X X X X X X X X]      [X              ]      [               ]      [               ]
// 	// [               ]      [               ]      [               ]      [               ]

// 	// captures and moves formward forward
// 	pawnMoves = (wp >> 7) & (EnemyPieces) & (^Rank8) & (^FileA) // capture right
// 	// Find first bit which is equal to '1' i.e. first capture
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+7, index, WP, 0, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (wp >> 9) & (EnemyPieces) & (^Rank8) & (^FileH) // capture left
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+9, index, WP, 0, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (wp >> 8) & Empty & (^Rank8) // move 1 square forward
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+8, index, WP, 0, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
// 	// also check that resulting square is on rank4.
// 	// (instead of check I mean eliminate squares that do not comply with these conditions)
// 	pawnMoves = (wp >> 16) & Empty & (Empty >> 8) & Rank4 // move 2 squares forward
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+16, index, WP, 0, MoveFlagPawnStart))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	// promotions
// 	pawnMoves = (wp >> 7) & EnemyPieces & Rank8 & (^FileA) // pawn promotion by capture right
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		// todo maybe Capture flag??
// 		moveList.AddMove(GetMoveInt(index+7, index, WP, WQ, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+7, index, WP, WR, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+7, index, WP, WB, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+7, index, WP, WN, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (wp >> 9) & EnemyPieces & Rank8 & (^FileH) // pawn promotion by capture left
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		// todo maybe Capture flag??
// 		moveList.AddMove(GetMoveInt(index+9, index, WP, WQ, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+9, index, WP, WR, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+9, index, WP, WB, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+9, index, WP, WN, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (wp >> 8) & Empty & Rank8 // pawn promotion by move 1 forward
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+8, index, WP, WQ, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+8, index, WP, WR, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+8, index, WP, WB, NoFlag))
// 		moveList.AddMove(GetMoveInt(index+8, index, WP, WN, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	// En passant right
// 	possibility = (wp << 1) & board.bitboards[BP] & Rank5 & (^FileA) & board.bitboards[EP]
// 	if possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-1, index-8, WP, 0, MoveFlagEnPass))
// 	}
// 	// en passant left
// 	possibility = (wp >> 1) & board.bitboards[BP] & Rank5 & (^FileH) & board.bitboards[EP]
// 	if possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+1, index-8, WP, 0, MoveFlagEnPass))
// 	}
// }

// func (board *Board) possibleBlackPawn(moveList *MoveList) {
// 	bp := board.bitboards[BP]
// 	var possibility uint64 // holds one potential capture at a time
// 	var index int          // index of the "possibility" capture
// 	var pawnMoves uint64

// 	// captures and moves formward forward:
// 	pawnMoves = (bp << 7) & (EnemyPieces) & (^Rank1) & (^FileH) // capture right
// 	// Find first bit which is equal to '1' i.e. first capture
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-7, index, BP, 0, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (bp << 9) & (EnemyPieces) & (^Rank1) & (^FileA) // capture left
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-9, index, BP, 0, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (bp << 8) & Empty & (^Rank1) // move 1 square forward
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-8, index, BP, 0, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	// Move all pawns 2 ranks, check that, in between and on the final square there is nothing,
// 	// also check that resulting square is on rank4.
// 	// (instead of check I mean eliminate squares that do not comply with these conditions)
// 	pawnMoves = (bp << 16) & Empty & (Empty << 8) & Rank5 // move 2 squares forward
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-16, index, BP, 0, MoveFlagPawnStart))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	// promotions
// 	pawnMoves = (bp << 7) & EnemyPieces & Rank1 & (^FileH) // pawn promotion by capture right
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		// todo maybe Capture flag??
// 		moveList.AddMove(GetMoveInt(index-7, index, BP, BQ, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-7, index, BP, BR, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-7, index, BP, BB, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-7, index, BP, BN, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (bp << 9) & EnemyPieces & Rank1 & (^FileA) // pawn promotion by capture left
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		// todo maybe Capture flag??
// 		moveList.AddMove(GetMoveInt(index-9, index, BP, BQ, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-9, index, BP, BR, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-9, index, BP, BB, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-9, index, BP, BN, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	pawnMoves = (bp << 8) & Empty & Rank1 // pawn promotion by move 1 forward
// 	possibility = pawnMoves & (^(pawnMoves - 1))
// 	for possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-8, index, BP, BQ, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-8, index, BP, BR, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-8, index, BP, BB, NoFlag))
// 		moveList.AddMove(GetMoveInt(index-8, index, BP, BN, NoFlag))
// 		pawnMoves &= ^possibility                    // remove the capture that we just analyzed
// 		possibility = pawnMoves & (^(pawnMoves - 1)) // find next bit equal to '1' i.e. next capture
// 	}

// 	// En passant right
// 	possibility = (bp >> 1) & board.bitboards[WP] & Rank4 & (^FileH) & board.bitboards[EP]
// 	if possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index+1, index+8, BP, 0, MoveFlagEnPass))
// 	}
// 	// en passant left
// 	possibility = (bp << 1) & board.bitboards[WP] & Rank4 & (^FileA) & board.bitboards[EP]
// 	if possibility != 0 {
// 		index = bits.TrailingZeros64(possibility)
// 		moveList.AddMove(GetMoveInt(index-1, index+8, BP, 0, MoveFlagEnPass))
// 	}
// }

// func (board *Board) possibleKnightMoves(moveList *MoveList, knight uint64, pieceType int) {
// 	// Choose bishop
// 	knightPossibility := knight & (^(knight - 1))
// 	var possibility uint64

// 	for knightPossibility != 0 {
// 		// Current knight index (in bitmask)
// 		knightIdx := bits.TrailingZeros64(knightPossibility)

// 		// Move knight pattern mask around depending on the current knight position idx
// 		// KnightSpan is a predefined knight move pattern such as
// 		// [               ]
// 		// [               ]
// 		// [               ]
// 		// [        X   X  ]
// 		// [      X       X]
// 		// [          O    ]
// 		// [      X       X]
// 		// [        X   X  ]
// 		if knightIdx > 18 {
// 			possibility = KnightSpan << (knightIdx - 18)
// 		} else {
// 			possibility = KnightSpan >> (18 - knightIdx)
// 		}

// 		// handle wrap around of knight pattern movement
// 		if knightIdx%8 < 4 {
// 			possibility &= (^FileGH) & NotMyPieces
// 		} else {
// 			possibility &= (^FileAB) & NotMyPieces
// 		}

// 		// choose move
// 		movePossibility := possibility & (^(possibility - 1))
// 		for movePossibility != 0 {
// 			// possible move index (in bitmask)
// 			moveIndex := bits.TrailingZeros64(movePossibility)
// 			moveList.AddMove(GetMoveInt(knightIdx, moveIndex, pieceType, 0, NoFlag))
// 			possibility &= ^movePossibility                      // remove move from all possible moves
// 			movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
// 		}
// 		knight &= ^knightPossibility
// 		knightPossibility = knight & (^(knight - 1))
// 	}
// }

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

// func (board *Board) possibleKingMoves(moveList *MoveList, king uint64, pieceType int) {
// 	var possibility uint64

// 	// Current king index (in bitmask)
// 	kingIdx := bits.TrailingZeros64(king)

// 	if kingIdx > 9 {
// 		possibility = KingSpan << (kingIdx - 9)
// 	} else {
// 		possibility = KingSpan >> (9 - kingIdx)
// 	}

// 	// handle wrap around of knight pattern movement
// 	if kingIdx%8 < 4 {
// 		possibility &= (^FileGH) & NotMyPieces
// 	} else {
// 		possibility &= (^FileAB) & NotMyPieces
// 	}

// 	// choose move
// 	movePossibility := possibility & (^(possibility - 1))
// 	for movePossibility != 0 {
// 		// possible move index (in bitmask)
// 		moveIndex := bits.TrailingZeros64(movePossibility)
// 		moveList.AddMove(GetMoveInt(kingIdx, moveIndex, pieceType, 0, NoFlag))
// 		possibility &= ^movePossibility                      // remove move from all possible moves
// 		movePossibility = possibility & (^(possibility - 1)) // calculate new possible move
// 	}
// }

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

// func (board *Board) unsafeForBlack() (unsafe uint64) {
// 	// todo should it update global value ??
// 	Occupied = (board.bitboards[WP] |
// 		board.bitboards[WN] |
// 		board.bitboards[WB] |
// 		board.bitboards[WR] |
// 		board.bitboards[WQ] |
// 		board.bitboards[WK] |
// 		board.bitboards[BP] |
// 		board.bitboards[BN] |
// 		board.bitboards[BB] |
// 		board.bitboards[BR] |
// 		board.bitboards[BQ] |
// 		board.bitboards[BK])

// 	// pawn
// 	unsafe = ((board.bitboards[WP] >> 7) & (^FileA))  // pawn capture right
// 	unsafe |= ((board.bitboards[WP] >> 9) & (^FileH)) // pawn capture left

// 	var possibility uint64
// 	// knight
// 	wn := board.bitboards[WN]
// 	i := wn & (^(wn - 1))
// 	for i != 0 {
// 		iLocation := bits.TrailingZeros64(i)
// 		if iLocation > 18 {
// 			possibility = KnightSpan << (iLocation - 18)
// 		} else {
// 			possibility = KnightSpan >> (18 - iLocation)
// 		}

// 		if iLocation%8 < 4 {
// 			possibility &= (^FileGH)
// 		} else {
// 			possibility &= (^FileAB)
// 		}
// 		unsafe |= possibility
// 		wn &= (^i)
// 		i = wn & (^(wn - 1))
// 	}

// 	// bishop/queen
// 	qb := board.bitboards[WQ] | board.bitboards[WB]
// 	i = qb & (^(qb - 1))
// 	for i != 0 {
// 		iLocation := bits.TrailingZeros64(i)
// 		possibility = board.DiagonalAndAntiDiagonalMoves(iLocation)
// 		unsafe |= possibility
// 		qb &= (^i)
// 		i = qb & (^(qb - 1))
// 	}

// 	// rook/queen
// 	qr := board.bitboards[WQ] | board.bitboards[WR]
// 	i = qr & (^(qr - 1))
// 	for i != 0 {
// 		iLocation := bits.TrailingZeros64(i)
// 		possibility = board.HorizontalAndVerticalMoves(iLocation)
// 		unsafe |= possibility
// 		qr &= (^i)
// 		i = qr & (^(qr - 1))
// 	}

// 	// king
// 	iLocation := bits.TrailingZeros64(board.bitboards[WK])
// 	if iLocation > 9 {
// 		possibility = KingSpan << (iLocation - 9)
// 	} else {
// 		possibility = KingSpan >> (9 - iLocation)
// 	}

// 	if iLocation%8 < 4 {
// 		possibility &= (^FileGH)
// 	} else {
// 		possibility &= (^FileAB)
// 	}
// 	unsafe |= possibility
// 	return unsafe
// }

// func (board *Board) unsafeForWhite() (unsafe uint64) {
// 	// todo should it update global value ??
// 	Occupied = (board.bitboards[WP] |
// 		board.bitboards[WN] |
// 		board.bitboards[WB] |
// 		board.bitboards[WR] |
// 		board.bitboards[WQ] |
// 		board.bitboards[WK] |
// 		board.bitboards[BP] |
// 		board.bitboards[BN] |
// 		board.bitboards[BB] |
// 		board.bitboards[BR] |
// 		board.bitboards[BQ] |
// 		board.bitboards[BK])

// 	// pawn
// 	unsafe = ((board.bitboards[BP] << 7) & (^FileH))  // pawn capture right
// 	unsafe |= ((board.bitboards[BP] << 9) & (^FileA)) // pawn capture left

// 	var possibility uint64
// 	// knight
// 	bn := board.bitboards[BN]
// 	i := bn & (^(bn - 1))
// 	for i != 0 {
// 		iLocation := bits.TrailingZeros64(i)
// 		if iLocation > 18 {
// 			possibility = KnightSpan << (iLocation - 18)
// 		} else {
// 			possibility = KnightSpan >> (18 - iLocation)
// 		}

// 		if iLocation%8 < 4 {
// 			possibility &= (^FileGH)
// 		} else {
// 			possibility &= (^FileAB)
// 		}
// 		unsafe |= possibility
// 		bn &= (^i)
// 		i = bn & (^(bn - 1))
// 	}

// 	// bishop/queen
// 	qb := board.bitboards[BQ] | board.bitboards[BB]
// 	i = qb & (^(qb - 1))
// 	for i != 0 {
// 		iLocation := bits.TrailingZeros64(i)
// 		possibility = board.DiagonalAndAntiDiagonalMoves(iLocation)
// 		unsafe |= possibility
// 		qb &= (^i)
// 		i = qb & (^(qb - 1))
// 	}

// 	// rook/queen
// 	qr := board.bitboards[BQ] | board.bitboards[BR]
// 	i = qr & (^(qr - 1))
// 	for i != 0 {
// 		iLocation := bits.TrailingZeros64(i)
// 		possibility = board.HorizontalAndVerticalMoves(iLocation)
// 		unsafe |= possibility
// 		qr &= (^i)
// 		i = qr & (^(qr - 1))
// 	}

// 	// king
// 	iLocation := bits.TrailingZeros64(board.bitboards[BK])
// 	if iLocation > 9 {
// 		possibility = KingSpan << (iLocation - 9)
// 	} else {
// 		possibility = KingSpan >> (9 - iLocation)
// 	}

// 	if iLocation%8 < 4 {
// 		possibility &= (^FileGH)
// 	} else {
// 		possibility &= (^FileAB)
// 	}
// 	unsafe |= possibility
// 	return unsafe
// }
