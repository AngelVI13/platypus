package board

// PieceValue A map used to identify a piece's value
var PieceValue = map[int]int{
   NoPiece: 0,
   WP: 100,
   WN: 325,
   WB: 325,
   WR: 550,
   WQ: 1000,
   WK: 50000,
   BP: 100,
   BN: 325,
   BB: 325,
   BR: 550,
   BQ: 1000,
   BK: 50000,
}

// PawnTable pawn table
var PawnTable = [BoardSquareNum]int{
	0,  0,  0,   0,   0,  0,  0,  0,
   10, 10,  0, -10, -10,  0, 10, 10,
	5,  0,  0,   5,   5,  0,  0,  5,
	0,  0, 10,  20,  20, 10,  0,  0,
	5,  5,  5,  10,  10,  5,  5,  5,
   10, 10, 10,  20,  20, 10, 10, 10,
   20, 20, 20,  30,  30, 20, 20, 20,
	0,  0,  0,   0,   0,  0,  0,  0,
}

// KnightTable knight table
var KnightTable = [BoardSquareNum]int{
   0, -10,  0,  0,  0, 0, -10,  0,
   0,   0,  0,  5,  5,  0,  0,  0,
   0,   0, 10, 10, 10, 10,  0,  0,
   0,   0, 10, 20, 20, 10,  0,  0,
   5,  10, 15, 20, 20, 15, 10,  5,
   5,  10, 10, 20, 20, 10, 10,  5,
   0,   0,  5, 10, 10,  5,  0,  0,
   0,   0,  0,  0,  0,  0,  0,  0,
}

// BishopTable bishop table
var BishopTable = [BoardSquareNum]int{
   0,  0, -10,  0,  0, -10,  0, 0,
   0,  0,   0, 10, 10,   0,  0, 0,
   0,  0,  10, 15, 15,  10,  0, 0,
   0, 10,  15, 20, 20,  15, 10, 0,
   0, 10,  15, 20, 20,  15, 10, 0,
   0,  0,  10, 15, 15,  10,  0, 0,
   0,  0,   0, 10, 10,   0,  0, 0,
   0,  0,   0,  0,  0,   0,  0, 0,
}

// RookTable rook table
var RookTable = [BoardSquareNum]int{
	0,  0,  5, 10, 10,  5,  0,  0,
	0,  0,  5, 10, 10,  5,  0,  0,
	0,  0,  5, 10, 10,  5,  0,  0,
	0,  0,  5, 10, 10,  5,  0,  0,
	0,  0,  5, 10, 10,  5,  0,  0,
	0,  0,  5, 10, 10,  5,  0,  0,
   25, 25, 25, 25, 25, 25, 25, 25,
	0,  0,  5, 10, 10,  5,  0,  0,
}

// KingE king endgame table
var KingE = [BoardSquareNum]int{
   -50, -10,  0,  0,  0,  0, -10, -50,
   -10,   0, 10, 10, 10, 10,   0, -10,
	 0,  10, 15, 15, 15, 15,  10,   0,
	 0,  10, 15, 20, 20, 15,  10,   0,
	 0,  10, 15, 20, 20, 15,  10,   0,
	 0,  10, 15, 15, 15, 15,  10,   0,
   -10,   0, 10, 10, 10, 10,   0, -10,
   -50, -10,  0,  0,  0,  0, -10, -50,
}

// KingO king opening/middle game table
var KingO = [BoardSquareNum]int{
	 0,  10,  10, -10, -10,   0,  20, 10,
   -30, -30, -30, -30, -30, -30, -30, -30,
   -50, -50, -50, -50, -50, -50, -50, -50,
   -70, -70, -70, -70, -70, -70, -70, -70,
   -70, -70, -70, -70, -70, -70, -70, -70,
   -70, -70, -70, -70, -70, -70, -70, -70,
   -70, -70, -70, -70, -70, -70, -70, -70,
   -70, -70, -70, -70, -70, -70, -70, -70,
}

// Mirror64 slice that is used to get a mirror version of the tables for black's evaluation
var Mirror64 = [BoardSquareNum]int{
   56, 57, 58, 59, 60, 61, 62, 63,
   48, 49, 50, 51, 52, 53, 54, 55,
   40, 41, 42, 43, 44, 45, 46, 47,
   32, 33, 34, 35, 36, 37, 38, 39,
   24, 25, 26, 27, 28, 29, 30, 31,
   16, 17, 18, 19, 20, 21, 22, 23,
	8,  9, 10, 11, 12, 13, 14, 15,
	0,  1,  2,  3,  4,  5,  6,  7,
}

// PawnPassed passed pawn bonuses depending on how far down the board it is
var PawnPassed = [8]int{0, 5, 10, 20, 35, 60, 100, 200}

const (
	// PawnIsolated isolated pawn bonus
	PawnIsolated = -10
	// PawnDoubled doubled pawn bonus
	PawnDoubled = -10
	// RookOpenFile rook on open file bonus
	RookOpenFile = 10
	// RookSemiOpenfile rook on semi-open file bonus
	RookSemiOpenfile = 5
	// QueenOpenFile queen on open file bonus
	QueenOpenFile = 5
	// QueenSemiOpenFile queen on semi-open file bonus
	QueenSemiOpenFile = 3
	// BishopPair bonus
	BishopPair = 30
	// KingNearOpenFile king on or near open file bonus
	KingNearOpenFile = -10
)

// EndGameMaterial defines the boundary limit for the endgame
var EndGameMaterial = 1*PieceValue[WR] + 2*PieceValue[WN] + 2*PieceValue[WP] + PieceValue[WK]


// EvalPosition evaluate position and return value
func (board *Board) EvalPosition() int {
   score := board.material[White] - board.material[Black] // get current score

   // !!! FIX THIS !!!!
   // if MaterialDraw(pos) == true && pos.pieceNum[WhitePawn] == 0 && pos.pieceNum[BlackPawn] == 0 {
   //    return 0
   // }

   // for sq := 0; sq < BoardSquareNum; sq++ {
   //    switch pos.Pieces[sq] {
   //    case OffBoard, Empty:
   //       continue
   //    case WhitePawn:
   //       evalWhitePawn(pos, sq, &score)
   //    case BlackPawn:
   //       evalBlackPawn(pos, sq, &score)
   //    case WhiteKnight:
   //       score += KnightTable[Sq120ToSq64[sq]]
   //    case BlackKnight:
   //       score -= KnightTable[Mirror64[Sq120ToSq64[sq]]]
   //    case WhiteBishop:
   //       score += BishopTable[Sq120ToSq64[sq]]
   //    case BlackBishop:
   //       score -= BishopTable[Mirror64[Sq120ToSq64[sq]]]
   //    case WhiteRook:
   //       evalWhiteRook(pos, sq, &score)
   //    case BlackRook:
   //       evalBlackRook(pos, sq, &score)
   //    case WhiteQueen:
   //       evalWhiteQueen(pos, sq, &score)
   //    case BlackQueen:
   //       evalBlackQueen(pos, sq, &score)
   //    case WhiteKing:
   //       evalWhiteKing(pos, sq, &score)
   //    case BlackKing:
   //       evalBlackKing(pos, sq, &score)
   //    default:
   //       panic("Unknown square type")
   //    }
   // }

   // if pos.pieceNum[WhiteBishop] >= 2 {
   //    score += BishopPair
   // }
   // if pos.pieceNum[BlackBishop] >= 2 {
   //    score -= BishopPair
   // }

   // if pos.side == White {
   //    return score
   // }
   // return -score

   // todo remove this
   return score
}

// MirrorBoard takes in a position and modifies it to be the mirrored version of it
// Only used for testing.
func MirrorBoard(board *Board) {
   var swapPiece = [14]int{NoPiece, BP, BN, BB, BR, BQ, BK, WP, WN, WB, WR, WQ, WK, EP }
   var tempMaterial = [2]int{0, 0}
   var tempPosition [64]int

   // Mirror side to move
   board.Side ^= 1

   // todo this (for loop) does not mirror bitboards: 
   //      example e2e4 c7c6 should mirror as c2c3 e7e5 but this mirrors it as e2e4 c7c5 (i.e. just reverses the colors)
   // Mirror bit boards
   for idx, _ := range board.bitboards {
      // swap piece bitboards i.e. WP becomes BP now
      board.bitboards[idx], board.bitboards[swapPiece[idx]] = board.bitboards[swapPiece[idx]], board.bitboards[idx]
   }

   // Mirror pieces in position variable & update material score
   // Invert pieces' positions
   for sq := range board.position {
      tempPosition[sq] = board.position[Mirror64[sq]]
   }

   // reset board
   for sq := range board.position {
      board.position[sq] = NoPiece
   }

   // Swap white & black pieces
   for sq := range board.position {
      mirroredPiece := swapPiece[tempPosition[sq]]
      board.position[sq] = mirroredPiece

      if mirroredPiece < BP {
         tempMaterial[White] += PieceValue[mirroredPiece]
      } else {
         tempMaterial[Black] += PieceValue[mirroredPiece]
      }
   }
   board.material = tempMaterial
   
   // Mirror stateBoards
   board.UpdateBitMasks()

   // Mirror castling permissions
   var tempCastlePerm = 0
   if board.castlePermissions&WhiteKingCastling != 0 {
      tempCastlePerm |= BlackKingCastling
   }

   if board.castlePermissions&WhiteQueenCastling != 0 {
      tempCastlePerm |= BlackQueenCastling
   }

   if board.castlePermissions&BlackKingCastling != 0 {
      tempCastlePerm |= WhiteKingCastling
   }

   if board.castlePermissions&BlackQueenCastling != 0 {
      tempCastlePerm |= WhiteQueenCastling
   }
   board.castlePermissions = tempCastlePerm

   board.positionKey = GeneratePositionKey(board)
}
