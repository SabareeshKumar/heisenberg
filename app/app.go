package app

import (
	"errors"
)

// GameState represents status of active game
type GameState struct {
	board    *boardConfig
	myPieces map[int]*piece
}

var game GameState

// InitGame sets up a new game.
func InitGame(colorChoice int) {
	board := newBoard()
	myPieces := make(map[int]*piece, 16)
	game = GameState{board, myPieces}
	if colorChoice == white {
		for i := 48; i <= 63; i++ {
			myPieces[i] = board.pieces[i]
		}
		return
	}
	for i := 0; i <= 15; i++ {
		myPieces[i] = board.pieces[i]
	}
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
	myMov := myMove()
	err = board.alterPosition(myMov)
	if err != nil {
		return UserMove{}, err
	}
	myMoveCoord, err := myMov.toUserMove()
	if err != nil {
		return UserMove{}, err
	}
	return myMoveCoord, nil
}

func myMove() boardMove {
	// Create a channel to receive moves on the fly
	moveCh := make(chan boardMove)
	for _, piece := range game.myPieces {
		// Move calculation is concurrent
		go piece.moveGenerator(piece, moveCh)
	}
	searchResults := make(chan searchResult)
	var bestMove boardMove
	for move := range moveCh {
		if !isMoveLegal(move) {
			continue
		}
		// Search is concurrent
		go search(move, searchResults)
		// TODO: compute best move
		bestMove = move
		break
	}
	for _ = range searchResults {
		// Choose best result
		break
	}
	return bestMove
}

func isMoveLegal(_ boardMove) bool {
	// TODO: check if move is legal
	return true
}
