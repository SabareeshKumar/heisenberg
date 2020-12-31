package app

import (
	"errors"
)

// GameState represents status of active game
type GameState struct {
	board    *boardConfig
	myPieces []*piece
}

var game GameState

// InitGame sets up a new game.
func InitGame(colorChoice int) {
	board := newBoard()
	var myPieces []*piece
	if colorChoice == black {
		myPieces = board.pieces[48:]
	} else {
		myPieces = board.pieces[:16]
	}
	game = GameState{board, myPieces}
}

// MakeMove returns the computer's move given a move made by the user.
func MakeMove(move UserMove) (UserMove, error) {
	uMove, err := move.toBoardMove()
	if !isMoveLegal(uMove) {
		return UserMove{}, errors.New("Illegal move")
	}
	board := game.board
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
