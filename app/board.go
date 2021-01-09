package app

import (
	"errors"
	"math"
)

type boardConfig struct {
	pieces []*piece
}

func newBoard() *boardConfig {
	pieces := make([]*piece, 64)
	pieces[0] = &piece{
		id:            rook,
		name:          "White Rook",
		color:         white,
		position:      1 << 0,
		moveGenerator: rookMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[1] = &piece{
		id:            knight,
		name:          "White Knight",
		color:         white,
		position:      1 << 1,
		moveGenerator: knightMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[2] = &piece{
		id:            bishop,
		name:          "White Bishop",
		color:         white,
		position:      1 << 2,
		moveGenerator: bishopMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[3] = &piece{
		id:            queen,
		name:          "White Queen",
		color:         white,
		position:      1 << 3,
		moveGenerator: queenMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[4] = &piece{
		id:            king,
		name:          "White King",
		color:         white,
		position:      1 << 4,
		moveGenerator: kingMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[5] = &piece{
		id:            bishop,
		name:          "White Bishop",
		color:         white,
		position:      1 << 5,
		moveGenerator: bishopMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[6] = &piece{
		id:            knight,
		name:          "White Knight",
		color:         white,
		position:      1 << 6,
		moveGenerator: knightMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[7] = &piece{
		id:            rook,
		name:          "White Rook",
		color:         white,
		position:      1 << 7,
		moveGenerator: rookMoves,
		captured:      false,
		enpassantMove: -1,
	}
	for i := 8; i <= 15; i++ {
		pieces[i] = &piece{
			id:            pawn,
			name:          "White Pawn",
			color:         white,
			position:      1 << i,
			moveGenerator: pawnMoves,
			captured:      false,
			enpassantMove: -1,
		}
	}
	for i := 48; i <= 55; i++ {
		pieces[i] = &piece{
			id:            pawn,
			name:          "Black Pawn",
			color:         black,
			position:      1 << i,
			moveGenerator: pawnMoves,
			captured:      false,
			enpassantMove: -1,
		}
	}
	// Create Black pieces
	pieces[56] = &piece{
		id:            rook,
		name:          "Black Rook",
		color:         black,
		position:      1 << 56,
		moveGenerator: rookMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[57] = &piece{
		id:            knight,
		name:          "Black Knight",
		color:         black,
		position:      1 << 57,
		moveGenerator: knightMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[58] = &piece{
		id:            bishop,
		name:          "Black Bishop",
		color:         black,
		position:      1 << 58,
		moveGenerator: bishopMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[59] = &piece{
		id:            queen,
		name:          "Black Queen",
		color:         black,
		position:      1 << 59,
		moveGenerator: queenMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[60] = &piece{
		id:            king,
		name:          "Black King",
		color:         black,
		position:      1 << 60,
		moveGenerator: kingMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[61] = &piece{
		id:            bishop,
		name:          "Black Bishop",
		color:         black,
		position:      1 << 61,
		moveGenerator: bishopMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[62] = &piece{
		id:            knight,
		name:          "Black Knight",
		color:         black,
		position:      1 << 62,
		moveGenerator: knightMoves,
		captured:      false,
		enpassantMove: -1,
	}
	pieces[63] = &piece{
		id:            rook,
		name:          "Black Rook",
		color:         black,
		position:      1 << 63,
		moveGenerator: rookMoves,
		captured:      false,
		enpassantMove: -1,
	}
	return &boardConfig{pieces}
}

func (brd *boardConfig) alterPosition(bm boardMove) error {
	piece := brd.pieces[bm.From]
	if piece == nil {
		return errors.New("Invalid move")
	}
	capturedPc := bm.captured
	if capturedPc != nil {
		capturedPc.captured = true
		if piece.color == game.myColor {
			game.materialBalance += weights[capturedPc.id]
		} else {
			game.materialBalance -= weights[capturedPc.id]
		}
	}
	brd.pieces[bm.From] = nil
	brd.pieces[bm.To] = piece
	piece.position = 1 << bm.To
	game.moveCount += 1
	piece.moveCount += 1
	if piece.id == pawn && int(math.Abs(float64(bm.To-bm.From))) == 16 {
		piece.enpassantMove = game.moveCount
	}
	if bm.castlingFrom == -1 {
		return nil
	}
	rookPc := brd.pieces[bm.castlingFrom]
	brd.pieces[bm.castlingFrom] = nil
	brd.pieces[bm.castlingTo] = rookPc
	rookPc.position = 1 << bm.castlingTo
	rookPc.moveCount += 1
	return nil
}

func (brd *boardConfig) undoMove(bm boardMove) {
	game.moveCount -= 1
	piece := brd.pieces[bm.To]
	brd.pieces[bm.From] = piece
	capturedPc := bm.captured
	brd.pieces[bm.To] = capturedPc
	piece.position = 1 << bm.From
	piece.moveCount -= 1
	if piece.id == pawn && int(math.Abs(float64(bm.To-bm.From))) == 16 {
		piece.enpassantMove = -1
	}
	if bm.castlingFrom != -1 {
		rookPc := brd.pieces[bm.castlingTo]
		brd.pieces[bm.castlingFrom] = rookPc
		brd.pieces[bm.castlingTo] = nil
		rookPc.position = 1 << bm.castlingFrom
		rookPc.moveCount -= 1
	}
	if capturedPc == nil {
		return
	}
	capturedPc.captured = false
	if piece.color == game.myColor {
		game.materialBalance -= weights[capturedPc.id]
	} else {
		game.materialBalance += weights[capturedPc.id]
	}
}
