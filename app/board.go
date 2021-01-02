package app

import (
	"errors"
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
	}
	pieces[1] = &piece{
		id:            knight,
		name:          "White Knight",
		color:         white,
		position:      1 << 1,
		moveGenerator: knightMoves,
		captured:      false,
	}
	pieces[2] = &piece{
		id:            bishop,
		name:          "White Bishop",
		color:         white,
		position:      1 << 2,
		moveGenerator: bishopMoves,
		captured:      false,
	}
	pieces[3] = &piece{
		id:            queen,
		name:          "White Queen",
		color:         white,
		position:      1 << 3,
		moveGenerator: queenMoves,
		captured:      false,
	}
	pieces[4] = &piece{
		id:            king,
		name:          "White King",
		color:         white,
		position:      1 << 4,
		moveGenerator: kingMoves,
		captured:      false,
	}
	pieces[5] = &piece{
		id:            bishop,
		name:          "White Bishop",
		color:         white,
		position:      1 << 5,
		moveGenerator: bishopMoves,
		captured:      false,
	}
	pieces[6] = &piece{
		id:            knight,
		name:          "White Knight",
		color:         white,
		position:      1 << 6,
		moveGenerator: knightMoves,
		captured:      false,
	}
	pieces[7] = &piece{
		id:            rook,
		name:          "White Rook",
		color:         white,
		position:      1 << 7,
		moveGenerator: rookMoves,
		captured:      false,
	}
	for i := 8; i <= 15; i++ {
		pieces[i] = &piece{
			id:            pawn,
			name:          "White Pawn",
			color:         white,
			position:      1 << i,
			moveGenerator: pawnMoves,
			captured:      false,
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
	}
	pieces[57] = &piece{
		id:            knight,
		name:          "Black Knight",
		color:         black,
		position:      1 << 57,
		moveGenerator: knightMoves,
		captured:      false,
	}
	pieces[58] = &piece{
		id:            bishop,
		name:          "Black Bishop",
		color:         black,
		position:      1 << 58,
		moveGenerator: bishopMoves,
		captured:      false,
	}
	pieces[59] = &piece{
		id:            queen,
		name:          "Black Queen",
		color:         black,
		position:      1 << 59,
		moveGenerator: queenMoves,
		captured:      false,
	}
	pieces[60] = &piece{
		id:            king,
		name:          "Black King",
		color:         black,
		position:      1 << 60,
		moveGenerator: kingMoves,
		captured:      false,
	}
	pieces[61] = &piece{
		id:            bishop,
		name:          "Black Bishop",
		color:         black,
		position:      1 << 61,
		moveGenerator: bishopMoves,
		captured:      false,
	}
	pieces[62] = &piece{
		id:            knight,
		name:          "Black Knight",
		color:         black,
		position:      1 << 62,
		moveGenerator: knightMoves,
		captured:      false,
	}
	pieces[63] = &piece{
		id:            rook,
		name:          "Black Rook",
		color:         black,
		position:      1 << 63,
		moveGenerator: rookMoves,
		captured:      false,
	}
	return &boardConfig{pieces}
}

func (brd *boardConfig) alterPosition(bm boardMove) error {
	piece := brd.pieces[bm.From]
	if piece == nil {
		return errors.New("Invalid move")
	}
	err := piece.move(bm)
	if err != nil {
		return err
	}
	capturedPc := brd.pieces[bm.To]
	if capturedPc != nil {
		capturedPc.captured = true
	}
	brd.pieces[bm.From] = nil
	brd.pieces[bm.To] = piece
	return nil
}
