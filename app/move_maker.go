package app

import (
	"log"
)

// MakeMove returns the computer's move given a move made by the user.
func MakeMove(move Move) (Move, error) {
	fromIndex, toIndex, err := move.boardIndices()
	log.Print(fromIndex, toIndex)
	return move, err
}
