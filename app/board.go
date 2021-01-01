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
		id:            0,
		name:          "White Rook",
		color:         white,
		position:      1 << 0,
		moveGenerator: rookMoves,
	}
	pieces[1] = &piece{
		id:            1,
		name:          "White Knight",
		color:         white,
		position:      1 << 1,
		moveGenerator: knightMoves,
	}
	pieces[2] = &piece{
		id:            2,
		name:          "White Bishop",
		color:         white,
		position:      1 << 2,
		moveGenerator: bishopMoves,
	}
	pieces[3] = &piece{
		id:            3,
		name:          "White Queen",
		color:         white,
		position:      1 << 3,
		moveGenerator: queenMoves,
	}
	pieces[4] = &piece{
		id:            4,
		name:          "White King",
		color:         white,
		position:      1 << 4,
		moveGenerator: kingMoves,
	}
	pieces[5] = &piece{
		id:            5,
		name:          "White Bishop",
		color:         white,
		position:      1 << 5,
		moveGenerator: bishopMoves,
	}
	pieces[6] = &piece{
		id:            6,
		name:          "White Knight",
		color:         white,
		position:      1 << 6,
		moveGenerator: knightMoves,
	}
	pieces[7] = &piece{
		id:            7,
		name:          "White Rook",
		color:         white,
		position:      1 << 7,
		moveGenerator: rookMoves,
	}
	for i := 8; i <= 15; i++ {
		pieces[i] = &piece{
			id:            i,
			name:          "White Pawn",
			color:         white,
			position:      1 << i,
			moveGenerator: pawnMoves,
		}
	}
	for i := 48; i <= 55; i++ {
		pieces[i] = &piece{
			id:            i,
			name:          "Black Pawn",
			color:         black,
			position:      1 << i,
			moveGenerator: pawnMoves,
		}
	}
	// Create Black pieces
	pieces[56] = &piece{
		id:            56,
		name:          "Black Rook",
		color:         black,
		position:      1 << 56,
		moveGenerator: rookMoves,
	}
	pieces[57] = &piece{
		id:            57,
		name:          "Black Knight",
		color:         black,
		position:      1 << 57,
		moveGenerator: knightMoves,
	}
	pieces[58] = &piece{
		id:            58,
		name:          "Black Bishop",
		color:         black,
		position:      1 << 58,
		moveGenerator: bishopMoves,
	}
	pieces[59] = &piece{
		id:            59,
		name:          "Black Queen",
		color:         black,
		position:      1 << 59,
		moveGenerator: queenMoves,
	}
	pieces[60] = &piece{
		id:            60,
		name:          "Black King",
		color:         black,
		position:      1 << 60,
		moveGenerator: kingMoves,
	}
	pieces[61] = &piece{
		id:            61,
		name:          "Black Bishop",
		color:         black,
		position:      1 << 61,
		moveGenerator: bishopMoves,
	}
	pieces[62] = &piece{
		id:            62,
		name:          "Black Knight",
		color:         black,
		position:      1 << 62,
		moveGenerator: knightMoves,
	}
	pieces[63] = &piece{
		id:            63,
		name:          "Black Rook",
		color:         black,
		position:      1 << 63,
		moveGenerator: rookMoves,
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
	brd.pieces[bm.From] = nil
	brd.pieces[bm.To] = piece
	return nil
}
