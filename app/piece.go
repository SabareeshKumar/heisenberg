package app

import (
	"errors"
	"fmt"
)

const (
	black = 1
	white = 2
)

const (
	king   = 1
	queen  = 2
	rook   = 3
	bishop = 4
	knight = 5
	pawn   = 6
)

type piece struct {
	id            int
	name          string
	color         int
	position      uint // position in powers of 2
	moveGenerator func(*piece, chan boardMove)
}

func (p *piece) move(mv boardMove) error {
	piece := game.board.pieces[mv.From]
	if piece.id != p.id {
		errMsg := fmt.Sprintf(
			"'%s' is not at position '%s'", piece.name, mv.From)
		return errors.New(errMsg)
	}
	p.position = 1 << mv.To
	return nil
}
