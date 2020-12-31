package app

import (
	"errors"
	"log"
)

var board *boardConfig

// CreateBoard sets up a new game board
func CreateBoard() {
	board = newBoard()
	log.Print("Created new board configuration")
}

// MakeMove returns the computer's move given a move made by the user.
func MakeMove(move UserMove) (UserMove, error) {
	uMove, err := move.toBoardMove()
	if !isMoveLegal(uMove) {
		return UserMove{}, errors.New("Illegal move")
	}
	err = board.alterPosition(uMove)
	if err != nil {
		return UserMove{}, err
	}
	myMove := myMove()
	err = board.alterPosition(myMove)
	if err != nil {
		return UserMove{}, err
	}
	myMoveCoord, err := myMove.toUserMove()
	if err != nil {
		return UserMove{}, err
	}
	return myMoveCoord, nil
}

func myMove() boardMove {
	// TODO: compute best move
	return boardMove{52, 44}
}

func isMoveLegal(_ boardMove) bool {
	// TODO: check if move is legal
	return true
}
