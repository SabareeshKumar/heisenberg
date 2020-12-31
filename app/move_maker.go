package app

import (
	"errors"
)

// MakeMove returns the computer's move given a move made by the user.
func MakeMove(move UserMove) (UserMove, error) {
	uMove, err := move.toBoardMove()
	if !isMoveLegal(uMove) {
		return UserMove{}, errors.New("Illegal move")
	}
	alterPosition(uMove)
	myMove := myMove()
	alterPosition(myMove)
	myMoveCoord, err := myMove.toUserMove()
	if err != nil {
		return UserMove{}, err
	}
	return myMoveCoord, nil
}

func myMove() boardMove {
	// TODO: compute best move
	return boardMove{0, 63}
}

func alterPosition(_ boardMove) {
	// TODO: alter board position
}

func isMoveLegal(_ boardMove) bool {
	// TODO: check if move is legal
	return true
}
