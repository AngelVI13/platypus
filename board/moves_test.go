package board

import (
	"testing"
)

func TestGetMoveIntPawnStart(t *testing.T) {
	move := GetMoveInt(63, 0, WP, WQ, MoveFlagPawnStart)

	if FromSq(move) != 63 {
		t.Error("Wrong FROM square")
	}
	if ToSq(move) != 0 {
		t.Error("Wrong TO square")
	}
	if Promoted(move) != WQ {
		t.Error("Wrong PROMOTED piece")
	}
	if PieceType(move) != WP {
		t.Error("Wrong PIECE_TYPE square")
	}
	if EnPassantFlag(move) != 0 {
		t.Error("Wrong EN_PASSANT_FLAG square")
	}
	if PawnStartFlag(move) != 1 {
		t.Error("Wrong PAWN_START_FLAG square")
	}

	if CastleFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}
}

func TestGetMoveIntCastle(t *testing.T) {
	move := GetMoveInt(15, 27, WR, WN, MoveFlagCastle)

	if FromSq(move) != 15 {
		t.Error("Wrong FROM square")
	}
	if ToSq(move) != 27 {
		t.Error("Wrong TO square")
	}
	if Promoted(move) != WN {
		t.Error("Wrong PROMOTED piece")
	}
	if PieceType(move) != WR {
		t.Error("Wrong PIECE_TYPE square")
	}
	if EnPassantFlag(move) != 0 {
		t.Error("Wrong EN_PASSANT_FLAG square")
	}
	if PawnStartFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}

	if CastleFlag(move) != 1 {
		t.Error("Wrong PAWN_START_FLAG square")
	}
}

func TestGetMoveIntEnPassant(t *testing.T) {
	move := GetMoveInt(31, 4, WB, WR, MoveFlagEnPass)

	if FromSq(move) != 31 {
		t.Error("Wrong FROM square")
	}
	if ToSq(move) != 4 {
		t.Error("Wrong TO square")
	}
	if Promoted(move) != WR {
		t.Error("Wrong PROMOTED piece")
	}
	if PieceType(move) != WB {
		t.Error("Wrong PIECE_TYPE square")
	}
	if EnPassantFlag(move) != 1 {
		t.Error("Wrong EN_PASSANT_FLAG square")
	}
	if PawnStartFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}

	if CastleFlag(move) != 0 {
		t.Error("Wrong PAWN_START_FLAG square")
	}
}
