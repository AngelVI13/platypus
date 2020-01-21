package board

import (
	"fmt"
)

// MakeMove makes a move
func (board *Board) MakeMove(move int) {
	// todo add checks for side to move
	fromSq := FromSq(move)
	toSq := ToSq(move)
	pieceType := PieceType(move)

	DrawBitboard(board.bitboards[pieceType])
	// Remove piece from start sq in piece's bitboard
	board.bitboards[pieceType] &= (^(1 << fromSq))
	// Add the piece to end sq in piece's bitboard
	board.bitboards[pieceType] |= 1 << toSq

	DrawBitboard(board.bitboards[pieceType])
	// todo might have to update EnemyPieces here, otherwise we rely that somebody called GetMoves before that
	// todo Or store this bitboard in history and in current board

	// if ToSq is occupied by one of enemy's pieces -> it was a capture
	if (EnemyPieces >> toSq) & 1 == 1 {
		board.fiftyMove = 0  // reset 50 move rule counter
		fmt.Printf("Move capture on sq: %s\n", GetSquareString(toSq))
		for i := BP; i <= BK; i++ { // todo currently only looks at white captures of black pieces
			if (board.bitboards[i] >> toSq) & 1 == 1 {
				// todo update material here
				DrawBitboard(board.bitboards[i])
				board.bitboards[i] &= (^(1 << toSq))  // remove enemy piece from its board
				DrawBitboard(board.bitboards[i])
				break
			}
		}
	}

}