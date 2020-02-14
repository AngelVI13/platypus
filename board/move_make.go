package board

func (board *Board) removePieceFromSq(pieceType, sq int) {
	board.bitboards[pieceType] &= (^(1 << sq))
}

func (board *Board) addPieceToSq(pieceType, sq int) {
	board.bitboards[pieceType] |= 1 << sq
}

// MakeMove makes a move
func (board *Board) MakeMove(move int) {
	fromSq := FromSq(move)
	toSq := ToSq(move)
	pieceType := PieceType(move)

	// Remove piece from start sq in piece's bitboard
	board.removePieceFromSq(pieceType, fromSq)
	// Add the piece to end sq in piece's bitboard
	board.addPieceToSq(pieceType, toSq)

	if CastleFlag(move) == 1 {
		PerftCastles++

		if pieceType == WK && toSq == G1 {
			board.removePieceFromSq(WR, H1)
			board.addPieceToSq(WR, F1)
		} else if pieceType == WK && toSq == C1 {
			board.removePieceFromSq(WR, A1)
			board.addPieceToSq(WR, D1)
		} else if pieceType == BK && toSq == G8 {
			board.removePieceFromSq(WR, H8)
			board.addPieceToSq(WR, F8)
		} else if pieceType == WK && toSq == C8 {
			board.removePieceFromSq(WR, A8)
			board.addPieceToSq(WR, D8)
		} else {
			panic("Incorrect castle move")
		}
	}
	// if a rook or king has moved the remove the respective castling permission from castlePerm
	board.castlePermissions &= CastlePerm[fromSq]
	board.castlePermissions &= CastlePerm[toSq]

	// todo might have to update EnemyPieces here, otherwise we rely that somebody called GetMoves before that
	// todo Or store this bitboard in history and in current board
	// if ToSq is occupied by one of enemy's pieces -> it was a capture
	if (board.stateBoards[EnemyPieces]>>toSq)&1 == 1 {
		board.fiftyMove = 0 // reset 50 move rule counter
		PerftCaptures++

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
		PerftEnPassant++
		// Otherwise if destination piece is empty but the move is enpassant -> remove captured piece
		if board.Side == White {
			board.removePieceFromSq(BP, toSq+8)
		} else {
			board.removePieceFromSq(WP, toSq-8)
		}
	}

	// if en passant was set before this move then remove it
	// en passant is only available immediately after pawn start
	if board.bitboards[EP] != 0 {
		board.bitboards[EP] = 0
	}

	// if a pawn start -> update EnPassant bitboard
	if PawnStartFlag(move) == 1 {
		board.bitboards[EP] = FileMasks8[toSq%8]
	}

	if promoted := Promoted(move); promoted > 0 {
		PerftPromotions++
		// todo update material here
		board.addPieceToSq(promoted, toSq)
		// we already move the pawn to the 8th rank and since it is a promotion
		// not only we add the promoted piece but we need to remove the pawn
		// from the 8th rank
		board.removePieceFromSq(pieceType, toSq)
	}

	board.Side ^= 1 // change side to move
}
