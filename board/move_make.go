package board

import (
	"fmt"
)

// MakeMove makes a move
func (board *Board) MakeMove(move int) {
	fmt.Println("make move")
	fromSq := FromSq(move)
	toSq := ToSq(move)
	pieceType := PieceType(move)

	// Remove piece from start sq in piece's bitboard
	board.bitboards[pieceType] &= (^(1 << fromSq))
	// Add the piece to end sq in piece's bitboard
	board.bitboards[pieceType] |= 1 << toSq

	//! add handling for castling and remove castling permissions afterwards

	// todo might have to update EnemyPieces here, otherwise we rely that somebody called GetMoves before that
	// todo Or store this bitboard in history and in current board

	// if ToSq is occupied by one of enemy's pieces -> it was a capture
	if (EnemyPieces >> toSq) & 1 == 1 {
		board.fiftyMove = 0  // reset 50 move rule counter

		for i := WP; i <= BK; i++ {
			if (board.bitboards[i] >> toSq) & 1 == 1 {
				// todo update material here
				board.bitboards[i] &= (^(1 << toSq))  // remove enemy piece from its board
				break
			}
		}
	} else if ((Empty >> toSq) & 1 == 1) && (EnPassantFlag(move) == 1) {
		// Otherwise if destination piece is empty but the move is enpassant -> remove captured piece
		if board.Side == White {
			board.bitboards[BP] &= (^(1 << (toSq + 8)))
		} else {
			board.bitboards[WP] &= (^(1 << (toSq - 8)))
		}
	}

	// if en passant was set before this move then remove it
	// en passant is only available immediately after pawn start
	if board.bitboards[EP] != 0 {
		board.bitboards[EP] = 0
	}

	// if a pawn start -> update EnPassant bitboard
	if PawnStartFlag(move) == 1 {
		board.bitboards[EP] = FileMasks8[toSq % 8]
	}

	if promoted := Promoted(move); promoted > 0 {
		// todo update material here
		board.bitboards[promoted] |= (1 << toSq)
	}

	board.Side ^= 1  // change side to move
}