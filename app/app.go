package app

import (
	"errors"
	"fmt"
	"math"
)

// GameState represents status of active game
type GameState struct {
	board           *boardConfig
	myColor         int
	myPieces        map[int][]*piece
	otherPieces     map[int][]*piece
	materialBalance int
}

var game GameState

// InitGame sets up a new game.
func InitGame(colorChoice int) {
	whitePieces := make(map[int][]*piece, 16)
	blackPieces := make(map[int][]*piece, 16)
	board := newBoard()
	game = GameState{}
	game.board = board
	for i := 0; i <= 15; i++ {
		piece := board.pieces[i]
		whitePieces[piece.id] = append(whitePieces[piece.id], piece)
	}
	for i := 48; i <= 63; i++ {
		piece := board.pieces[i]
		blackPieces[piece.id] = append(blackPieces[piece.id], piece)
	}
	if colorChoice == white {
		game.myColor = black
		game.otherPieces = whitePieces
		game.myPieces = blackPieces
		return
	}
	game.myColor = white
	game.otherPieces = blackPieces
	game.myPieces = whitePieces
}

// MakeMove returns the computer's move given a move made by the user.
func MakeMove(move UserMove) (UserMove, error) {
	uMove, err := move.toBoardMove()
	piece := game.board.pieces[uMove.From]
	if piece.color == game.myColor {
		return UserMove{}, errors.New("Cannot move opponent piece")
	}
	valid := false
	for _, move := range piece.moveGenerator(piece) {
		if uMove == move {
			valid = true
			break
		}
	}
	if !valid {
		return UserMove{}, errors.New("Invalid move")
	}
	if !isMoveLegal(uMove) {
		return UserMove{}, errors.New("Illegal move")
	}
	fmt.Print("Thinking...")
	board := game.board
	err = board.alterPosition(uMove)
	if err != nil {
		return UserMove{}, err
	}
	myMov, _ := search(true, float32(math.MaxInt32), 1)
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

// func myMove() boardMove {
// 	// Create a channel to receive moves on the fly
// 	moveCh := make(chan boardMove)
// 	for _, pieces := range game.myPieces {
// 		for _, piece := range pieces {
// 			if piece.captured {
// 				continue
// 			}
// 			// Move calculation is concurrent
// 			go piece.moveGenerator(piece, moveCh)
// 		}
// 	}
// 	searchResults := make(chan searchResult)
// 	var bestMove boardMove
// 	for move := range moveCh {
// 		if !isMoveLegal(move) {
// 			continue
// 		}
// 		// Search is concurrent
// 		go search(move, searchResults)
// 		// TODO: compute best move
// 		bestMove = move
// 		break
// 	}
// 	for _ = range searchResults {
// 		// Choose best result
// 		break
// 	}
// 	return bestMove
// }

func isMoveLegal(mv boardMove) bool {
	board := game.board
	piece := board.pieces[mv.From]
	if piece.captured {
		// Cannot move dead piece
		return false
	}
	capturedPc := board.pieces[mv.To]
	if capturedPc != nil && piece.color == capturedPc.color {
		// Cannot capture own piece
		return false
	}
	// TODO: Check if move results in king in check
	// TODO: If move is castling, check its legality
	// TODO: If move is en passant, check its legality
	return true
}
