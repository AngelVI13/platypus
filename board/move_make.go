package board

import (
	"fmt"
	"math/bits"
	"strings"
	"errors"
)

func (board *Board) removePieceFromSq(pieceType, sq int) {
	board.bitboards[pieceType] &= (^(1 << sq))
	board.positionKey ^= PieceKeys[pieceType][sq]
	board.material[board.Side] -= PieceValue[pieceType]
	// fmt.Printf("-Unhashing piece %c from sq %s\n", PieceChar[pieceType], GetSquareString(sq))
}

func (board *Board) addPieceToSq(pieceType, sq int) {
	board.bitboards[pieceType] |= 1 << sq
	board.positionKey ^= PieceKeys[pieceType][sq]
	board.material[board.Side] += PieceValue[pieceType]
	// fmt.Printf("+Hashing piece %c from sq %s\n", PieceChar[pieceType], GetSquareString(sq))
}

// MakeMove makes a move
func (board *Board) MakeMove(move int) {
	fromSq := FromSq(move)
	toSq := ToSq(move)
	pieceType := board.position[fromSq]

	// Store hash value before we do any hashing in/out of pieces etc
	board.history[board.ply].positionKey = board.positionKey

	board.fiftyMove++ // increment fifty move rule

	// Remove piece from start sq in piece's bitboard
	board.removePieceFromSq(pieceType, fromSq)
	board.position[fromSq] = NoPiece
	// Add the piece to end sq in piece's bitboard
	board.addPieceToSq(pieceType, toSq)
	board.position[toSq] = pieceType

	// handle rook moves if castling is performed
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
			board.removePieceFromSq(BR, H8)
			board.addPieceToSq(BR, F8)
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

	// We need to get the captured piece from the move int
	// instead of directly from position[toSq] because in the
	// case of en-passant captures the toSq has no piece on it
	// however, the move int will indicate that there was a capture
	// and what the captured piece actually is
	capturedPiece := Captured(move)
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

	// If en passant was set before this move then remove it!
	// En passant is only available immediately after pawn start.
	// En passant might be set a bit later in this function
	// when we check if this move was a pawn start.
	if board.bitboards[EP] != 0 {
		enPassantFile := bits.TrailingZeros64(board.bitboards[EP])
		board.positionKey ^= PieceKeys[EP][enPassantFile]
		// fmt.Printf("-Unhashing enpass file %d\n", enPassantFile)
		board.history[board.ply].enPassantFile = board.bitboards[EP]
		board.bitboards[EP] = 0
	}

	// hash out the castling permissions
	board.positionKey ^= CastleKeys[board.castlePermissions]
	// fmt.Printf("-Unhashing castle perm %d\n", board.castlePermissions)

	// store history variables
	board.history[board.ply].move = move
	board.history[board.ply].fiftyMove = board.fiftyMove
	board.history[board.ply].castlePermissions = board.castlePermissions

	// if a rook or king has moved then remove the respective castling permission from castlePerm
	board.castlePermissions &= CastlePerm[fromSq]
	board.castlePermissions &= CastlePerm[toSq]

	// hash back in the castling perm
	board.positionKey ^= CastleKeys[board.castlePermissions]
	// fmt.Printf("+Hashing castle perm %d\n", board.castlePermissions)

	// if a pawn start -> update EnPassant bitboard
	if PawnStartFlag(move) == 1 {
		board.bitboards[EP] = FileMasks8[toSq%8]
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
		// fmt.Printf("+Hashing enpass file %d\n", toSq%8)
	}

	if promoted := Promoted(move); promoted > 0 {
		// PerftPromotions++
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
	// fmt.Printf("+Hashing side key %d\n", SideKey)
}

// TakeMove Reverts last move
func (board *Board) TakeMove() {
	board.ply--

	// get history for previous move in the ply
	move := board.history[board.ply].move
	fromSq := FromSq(move)
	toSq := ToSq(move)
	// the moved piece type is the piece in the to square
	pieceType := board.position[toSq]

	if board.bitboards[EP] != 0 {
		// Unhash enpassant if it was active(hashed)
		// This might be hashed again later if the previous move also
		// had an active en passant file.
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
		// fmt.Printf("-Unhashing enpass file %d\n", bits.TrailingZeros64(board.bitboards[EP]))
	}

	// fmt.Printf("-Unhashing castle perm %d\n", board.castlePermissions)
	// unhash current castling permissions and hash in the permissions from the last move (below)
	board.positionKey ^= CastleKeys[board.castlePermissions]

	board.castlePermissions = board.history[board.ply].castlePermissions
	board.fiftyMove = board.history[board.ply].fiftyMove
	board.bitboards[EP] = board.history[board.ply].enPassantFile

	if board.bitboards[EP] != 0 {
		// if en passant was active in the previous position - hash it in
		board.positionKey ^= PieceKeys[EP][bits.TrailingZeros64(board.bitboards[EP])]
		// fmt.Printf("+Hashing enpass file %d\n", bits.TrailingZeros64(board.bitboards[EP]))
	}

	board.positionKey ^= CastleKeys[board.castlePermissions]
	// fmt.Printf("+Hashing castle perm %d\n", board.castlePermissions)

	// change side and hash in the side key
	board.Side ^= 1
	board.positionKey ^= SideKey
	// fmt.Printf("+Hashing side key %d\n", SideKey)

	// Remove piece from to sq in piece's bitboard
	board.removePieceFromSq(pieceType, toSq)
	board.position[toSq] = NoPiece
	// Add the piece to start sq in piece's bitboard
	board.addPieceToSq(pieceType, fromSq)
	board.position[fromSq] = pieceType

	if CastleFlag(move) == 1 {
		if pieceType == WK && toSq == G1 {
			board.removePieceFromSq(WR, F1)
			board.addPieceToSq(WR, H1)
			board.position[F1] = NoPiece
			board.position[H1] = WR
		} else if pieceType == WK && toSq == C1 {
			board.removePieceFromSq(WR, D1)
			board.addPieceToSq(WR, A1)
			board.position[D1] = NoPiece
			board.position[A1] = WR
		} else if pieceType == BK && toSq == G8 {
			board.removePieceFromSq(BR, F8)
			board.addPieceToSq(BR, H8)
			board.position[F8] = NoPiece
			board.position[H8] = BR
		} else if pieceType == BK && toSq == C8 {
			board.removePieceFromSq(BR, D8)
			board.addPieceToSq(BR, A8)
			board.position[D8] = NoPiece
			board.position[A8] = BR
		} else {
			panic("Incorrect castle move")
		}
	}

	capturedPiece := Captured(move)
	if capturedPiece != NoPiece && EnPassantFlag(move) == 0 {
		board.addPieceToSq(capturedPiece, toSq) // add enemy piece from its board
		board.position[toSq] = capturedPiece
	} else if (capturedPiece == WP || capturedPiece == BP) && EnPassantFlag(move) == 1 {
		// re-add enpassant captured pawn
		if board.Side == White {
			board.addPieceToSq(BP, toSq+8)
			// add captured enpassant pawn to piece bitboard
			board.position[toSq+8] = BP
		} else {
			board.addPieceToSq(WP, toSq-8)
			board.position[toSq-8] = WP
		}
	}

	if promoted := Promoted(move); promoted > 0 {
		// So far we have moved back the piece to original square
		// however, when we have promotions we have ended up
		// moving the promoted piece back. So here we need to
		// replace the promoted piece with a pawn of the corresponding
		// color
		board.removePieceFromSq(promoted, fromSq)

		pawn := BP
		if board.Side == White {
			pawn = WP
		}

		board.addPieceToSq(pawn, fromSq)
		board.position[fromSq] = pawn
	}
}

// todo how to handle with a messed up board (if first few moves are okay but the last one is not valid?)
// PerformMoves Takes a string containing a space separated list of moves 
// and applies them to the board. Example `moves`: "e2e4 d7d5" ...
func (board *Board) MakeMoves(moves string) error {
	moveSlice := strings.Split(moves, " ")
	for _, moveString := range moveSlice {
		moveList := board.GetMoves()
		move, err := GetMoveFromString(&moveList, moveString)
		if err != nil {
			return errors.New(fmt.Sprintf("Incorrect move string: %s", moveString))
		}
		board.MakeMove(move)
	}

	return nil
}
