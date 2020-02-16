package board

import (
	"math/bits"
)

func (board *Board) removePieceFromSq(pieceType, sq int) {
	board.bitboards[pieceType] &= (^(1 << sq))
	board.positionKey ^= PieceKeys[pieceType][sq]
}

func (board *Board) addPieceToSq(pieceType, sq int) {
	board.bitboards[pieceType] |= 1 << sq
	board.positionKey ^= PieceKeys[pieceType][sq]
}

// MakeMove makes a move
func (board *Board) MakeMove(move int) {
	fromSq := FromSq(move)
	toSq := ToSq(move)
	pieceType := board.position[fromSq]
	capturedPiece := board.position[toSq]

	// Store has value before we do any hashing in/out of pieces etc
	board.history[board.ply].positionKey = board.positionKey

	board.fiftyMove++ // increment fifty move rule

	// Remove piece from start sq in piece's bitboard
	board.removePieceFromSq(pieceType, fromSq)
	board.position[fromSq] = NoPiece
	// Add the piece to end sq in piece's bitboard
	board.addPieceToSq(pieceType, toSq)
	board.position[toSq] = pieceType

	if CastleFlag(move) == 1 {
		// PerftCastles++

		if pieceType == WK && toSq == G1 {
			board.removePieceFromSq(WR, H1)
			board.addPieceToSq(WR, F1)
			board.position[H1] = NoPiece
			board.position[F1] = WR
		} else if pieceType == WK && toSq == C1 {
			board.removePieceFromSq(WR, A1)
			board.addPieceToSq(WR, D1)
			board.position[A1] = NoPiece
			board.position[D1] = WR
		} else if pieceType == BK && toSq == G8 {
			board.removePieceFromSq(WR, H8)
			board.addPieceToSq(WR, F8)
			board.position[H8] = NoPiece
			board.position[F8] = BR
		} else if pieceType == BK && toSq == C8 {
			board.removePieceFromSq(BR, A8)
			board.addPieceToSq(BR, D8)
			board.position[A8] = NoPiece
			board.position[D8] = BR
		} else {
			panic("Incorrect castle move")
		}
	}

	// capturedPiece := Captured(move)
	if capturedPiece != NoPiece && EnPassantFlag(move) == 0 {
		board.fiftyMove = 0 // reset 50 move rule counter

		// note: it is already removed from board.position -> no need to remove it here
		board.removePieceFromSq(capturedPiece, toSq) // remove enemy piece from its board
	} else if (capturedPiece == WP || capturedPiece == BP) && EnPassantFlag(move) == 1 {
		board.fiftyMove = 0 // reset 50 move rule counter

		// PerftEnPassant++
		// Otherwise if destination piece is empty but the move is enpassant -> remove captured piece
		if board.Side == White {
			board.removePieceFromSq(BP, toSq+8)
			// note: need to remove capture enpassant pawn since it is not in fromSq or toSq
			board.position[toSq+8] = NoPiece
		} else {
			board.removePieceFromSq(WP, toSq-8)
			board.position[toSq-8] = NoPiece
		}
	}

	// if en passant was set before this move then remove it
	// en passant is only available immediately after pawn start
	if board.bitboards[EP] != 0 {
		enPassantFile := bits.TrailingZeros64(board.bitboards[EP])
		board.positionKey ^= PieceKeys[EP][enPassantFile]
		board.history[board.ply].enPassantFile = board.bitboards[EP]
		board.bitboards[EP] = 0
	}

	// hash out the castling permissions
	board.positionKey ^= CastleKeys[board.castlePermissions]

	board.history[board.ply].move = move
	board.history[board.ply].fiftyMove = board.fiftyMove
	board.history[board.ply].castlePermissions = board.castlePermissions

	// if a rook or king has moved the remove the respective castling permission from castlePerm
	board.castlePermissions &= CastlePerm[fromSq]
	board.castlePermissions &= CastlePerm[toSq]

	// hash back in the castling perm
	board.positionKey ^= CastleKeys[board.castlePermissions]

	// if a pawn start -> update EnPassant bitboard
	if PawnStartFlag(move) == 1 {
		board.bitboards[EP] = FileMasks8[toSq%8]
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
	}

	if promoted := Promoted(move); promoted > 0 {
		// PerftPromotions++
		// todo update material here
		board.addPieceToSq(promoted, toSq)
		board.position[toSq] = promoted
		// we already move the pawn to the 8th rank and since it is a promotion
		// not only we add the promoted piece but we need to remove the pawn
		// from the 8th rank
		board.removePieceFromSq(pieceType, toSq)
	}
	board.ply++     // increase halfmove counter
	board.Side ^= 1 // change side to move
	board.positionKey ^= SideKey
}

// TakeMove Reverts last move
// todo update TakeMove to handle board.position
func (board *Board) TakeMove() {
	board.ply--

	move := board.history[board.ply].move
	fromSq := FromSq(move)
	toSq := ToSq(move)
	pieceType := board.position[fromSq]

	if board.bitboards[EP] != 0 {
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
	}

	board.positionKey ^= CastleKeys[board.castlePermissions]

	board.castlePermissions = board.history[board.ply].castlePermissions
	board.fiftyMove = board.history[board.ply].fiftyMove
	board.bitboards[EP] = board.history[board.ply].enPassantFile

	if board.bitboards[EP] != 0 {
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
	}

	board.positionKey ^= CastleKeys[board.castlePermissions]

	board.Side ^= 1
	board.positionKey ^= SideKey

	// Remove piece from to sq in piece's bitboard
	board.removePieceFromSq(pieceType, toSq)
	// Add the piece to start sq in piece's bitboard
	board.addPieceToSq(pieceType, fromSq)

	if CastleFlag(move) == 1 {
		// PerftCastles++

		if pieceType == WK && toSq == G1 {
			board.removePieceFromSq(WR, F1)
			board.addPieceToSq(WR, H1)
		} else if pieceType == WK && toSq == C1 {
			board.removePieceFromSq(WR, D1)
			board.addPieceToSq(WR, A1)
		} else if pieceType == BK && toSq == G8 {
			board.removePieceFromSq(WR, F8)
			board.addPieceToSq(WR, H8)
		} else if pieceType == WK && toSq == C8 {
			board.removePieceFromSq(WR, D8)
			board.addPieceToSq(WR, A8)
		} else {
			panic("Incorrect castle move")
		}
	}

	// todo Need to keep track of what piece is captured so that I an reinsert it here -> need to keep a [64]int of all piece on the board
	// if ToSq is occupied by one of enemy's pieces -> it was a capture
	if (board.stateBoards[EnemyPieces]>>toSq)&1 == 1 {
		// todo this should be done only once for the current board state NOT on every makemove
		var startRange int
		var endRange int
		if board.Side == White {
			// if current side is white the we capture from blacks pieces
			startRange = BP
			endRange = BK
		} else {
			startRange = WP
			endRange = WK
		}

		for i := startRange; i <= endRange; i++ {
			// if destination square is on the board -> must be the correct board for the capture
			if (board.bitboards[i]>>toSq)&1 == 1 {
				// todo update material here
				board.removePieceFromSq(i, toSq) // remove enemy piece from its board
				break
			}
		}
	} else if ((board.stateBoards[Empty]>>toSq)&1 == 1) && (EnPassantFlag(move) == 1) {
		// PerftEnPassant++
		// Otherwise if destination piece is empty but the move is enpassant -> remove captured piece
		if board.Side == White {
			board.removePieceFromSq(BP, toSq+8)
		} else {
			board.removePieceFromSq(WP, toSq-8)
		}
	}
}
