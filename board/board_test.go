package board

import (
	"strconv"
	"testing"
)

func TestParseStringArray(t *testing.T) {
	board := Board{}
	board.ParseFen(StartingPosition)

	whitePawnBinaryStr := "0000000011111111000000000000000000000000000000000000000000000000"
	whitePawnBitboard, _ := strconv.ParseUint(whitePawnBinaryStr, 2, 64)

	if board.bitboards[WP] != whitePawnBitboard {
		t.Errorf("White pawn bitboard incorrect\nExpected: %064b\nActual:   %064b", whitePawnBitboard, board.bitboards[WP])
	}

	blackPawnBinaryStr := "0000000000000000000000000000000000000000000000001111111100000000"
	blackPawnBitboard, _ := strconv.ParseUint(blackPawnBinaryStr, 2, 64)

	if board.bitboards[BP] != blackPawnBitboard {
		t.Errorf("Black pawn bitboard incorrect\nExpected: %064b\nActual:   %064b", blackPawnBitboard, board.bitboards[BP])
	}

	whiteRookBinaryStr := "1000000100000000000000000000000000000000000000000000000000000000"
	whiteRookBitboard, _ := strconv.ParseUint(whiteRookBinaryStr, 2, 64)

	if board.bitboards[WR] != whiteRookBitboard {
		t.Errorf("White rook bitboard incorrect\nExpected: %064b\nActual:   %064b", whiteRookBitboard, board.bitboards[WR])
	}

	blackRookBinaryStr := "0000000000000000000000000000000000000000000000000000000010000001"
	blackRookBitboard, _ := strconv.ParseUint(blackRookBinaryStr, 2, 64)

	if board.bitboards[BR] != blackRookBitboard {
		t.Errorf("Black rook bitboard incorrect\nExpected: %064b\nActual:   %064b", blackRookBitboard, board.bitboards[BR])
	}
}
