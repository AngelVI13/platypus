package board

import (
	"strconv"
	"testing"
)

func TestParseStringArray(t *testing.T) {
	board := Board{}
	board.ParseStringArray(StartingPosition)

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

func TestGenerateChess960(t *testing.T) {
	board := Board{}
	position := GenerateChess960()
	board.ParseStringArray(position) // parse position so we can print it nicely

	for _, piece := range position[1] { // rank 7
		if piece != "p" {
			t.Errorf("All pieces on the 7th rank are not black pawns: \n%s", &board)
			break
		}
	}

	for _, piece := range position[6] { // rank 2
		if piece != "P" {
			t.Errorf("All pieces on the 2nd rank are not white pawns: \n%s", &board)
			break
		}
	}

	var numBlackPieces = map[string]int{
		"r": 0,
		"b": 0,
		"n": 0,
		"q": 0,
		"k": 0,
	}

	for _, piece := range position[0] {
		if piece != " " {
			numBlackPieces[piece]++
		}
	}

	if (numBlackPieces["r"] != 2) ||
		(numBlackPieces["n"] != 2) ||
		(numBlackPieces["b"] != 2) ||
		(numBlackPieces["q"] != 1) ||
		(numBlackPieces["k"] != 1) {
		t.Errorf("Incorrect number of black pieces (major & minor):\n%s", &board)
	}

	var numWhitePieces = map[string]int{
		"R": 0,
		"B": 0,
		"N": 0,
		"Q": 0,
		"K": 0,
	}

	for _, piece := range position[7] {
		if piece != " " {
			numWhitePieces[piece]++
		}
	}

	if (numWhitePieces["R"] != 2) ||
		(numWhitePieces["N"] != 2) ||
		(numWhitePieces["B"] != 2) ||
		(numWhitePieces["Q"] != 1) ||
		(numWhitePieces["K"] != 1) {
		t.Errorf("Incorrect number of white pieces (major & minor):\n%s", &board)
	}
}
