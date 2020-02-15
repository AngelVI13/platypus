package board

import (
	"math/bits"
)

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

func (board *Board) unsafeForBlack() (unsafe uint64) {
	// pawn
	unsafe = ((board.bitboards[WP] >> 7) & (^FileA))  // pawn capture right
	unsafe |= ((board.bitboards[WP] >> 9) & (^FileH)) // pawn capture left

	var possibility uint64
	// knight
	wn := board.bitboards[WN]
	i := wn & (^(wn - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = KnightMoves[iLocation]
		unsafe |= possibility
		wn &= (^i)
		i = wn & (^(wn - 1))
	}

	// sliding pieces
	occupiedExludingKing := board.stateBoards[Occupied] ^ board.bitboards[BK]
	// bishop/queen
	qb := board.bitboards[WQ] | board.bitboards[WB]
	i = qb & (^(qb - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.DiagonalAndAntiDiagonalMoves(iLocation, occupiedExludingKing)
		unsafe |= possibility
		qb &= (^i)
		i = qb & (^(qb - 1))
	}

	// rook/queen
	qr := board.bitboards[WQ] | board.bitboards[WR]
	i = qr & (^(qr - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.HorizontalAndVerticalMoves(iLocation, occupiedExludingKing)
		unsafe |= possibility
		qr &= (^i)
		i = qr & (^(qr - 1))
	}

	// king
	iLocation := bits.TrailingZeros64(board.bitboards[WK])
	possibility = KingMoves[iLocation]
	unsafe |= possibility
	return unsafe
}

func (board *Board) unsafeForWhite() (unsafe uint64) {
	// pawn
	unsafe = ((board.bitboards[BP] << 7) & (^FileH))  // pawn capture right
	unsafe |= ((board.bitboards[BP] << 9) & (^FileA)) // pawn capture left

	var possibility uint64
	// knight
	bn := board.bitboards[BN]
	i := bn & (^(bn - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = KnightMoves[iLocation]
		unsafe |= possibility
		bn &= (^i)
		i = bn & (^(bn - 1))
	}

	// sliding pieces
	// when calculating unsafe squares for a given colour we need to exclude the
	// current side's king because if an enemy queen is attacking our king,
	// the squares behind the king are also unsafe, however, when the king is included
	// geneation of unsafe squares will stop at the king and will not extend behind it
	occupiedExludingKing := board.stateBoards[Occupied] ^ board.bitboards[WK]
	// bishop/queen
	qb := board.bitboards[BQ] | board.bitboards[BB]
	i = qb & (^(qb - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.DiagonalAndAntiDiagonalMoves(iLocation, occupiedExludingKing)
		unsafe |= possibility
		qb &= (^i)
		i = qb & (^(qb - 1))
	}

	// rook/queen
	qr := board.bitboards[BQ] | board.bitboards[BR]
	i = qr & (^(qr - 1))
	for i != 0 {
		iLocation := bits.TrailingZeros64(i)
		possibility = board.HorizontalAndVerticalMoves(iLocation, occupiedExludingKing)
		unsafe |= possibility
		qr &= (^i)
		i = qr & (^(qr - 1))
	}

	// king
	iLocation := bits.TrailingZeros64(board.bitboards[BK])
	possibility = KingMoves[iLocation]
	unsafe |= possibility
	return unsafe
}
