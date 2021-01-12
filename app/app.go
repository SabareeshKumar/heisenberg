package app

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// InProgress denotes if game is in progress
const InProgress = 0

// Win denotes if game is won by opponent
const Win = 1

// Lost denotes if game is won by opponent
const Lost = 2

// Stalemate denotes if game is a stalemate
const Stalemate = 3

// GameState represents status of active game
type GameState struct {
	board           *boardConfig
	myColor         int
	myPieces        map[int][]*piece
	otherPieces     map[int][]*piece
	materialBalance int
	moveCount       int
}

var game GameState

// InitGame sets up a new game.
func InitGame(colorChoice int) {
	whitePieces := make(map[int][]*piece, 16)
	blackPieces := make(map[int][]*piece, 16)
	board := newBoard()
	game = GameState{}
	game.board = board
	game.moveCount = 0
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

// MakeMove verifies if the user move if valid
func MakeMove(uMove boardMove) error {
	piece := game.board.pieces[uMove.From]
	if piece.color == game.myColor {
		return errors.New("Cannot move opponent piece")
	}
	valid := false
	for _, move := range piece.moveGenerator(piece) {
		if uMove == move {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("Invalid move")
	}
	if !isMoveLegal(uMove) {
		return errors.New("Illegal move")
	}
	return game.board.alterPosition(uMove)
}

// MyMove computes a move for the engine
func MyMove() (UserMove, error) {
	myMov, _ := search(true, float32(math.MaxInt32), 1)
	err := game.board.alterPosition(myMov)
	if err != nil {
		return UserMove{}, err
	}
	myMoveCoord, err := myMov.toUserMove()
	if err != nil {
		return UserMove{}, err
	}
	return myMoveCoord, nil
}

func isMoveLegal(mv boardMove) bool {
	board := game.board
	pc := board.pieces[mv.From]
	if pc.captured {
		// Cannot move dead piece
		return false
	}
	capturedPc := board.pieces[mv.To]
	if capturedPc != nil {
		if pc.color == capturedPc.color {
			// Cannot capture own piece
			return false
		}
		if capturedPc.id == king {
			// Cannot capture king
			return false
		}
	}
	var kingPc *piece
	var otherPieces map[int][]*piece
	if pc.color == game.myColor {
		kingPc = game.myPieces[king][0]
		otherPieces = game.otherPieces
	} else {
		kingPc = game.otherPieces[king][0]
		otherPieces = game.myPieces
	}
	if mv.castlingFrom > 0 {
		// Check if castling is valid. It inherently checks if move
		// results in check.
		return isCastlingValid(mv, otherPieces)
	}
	defer board.undoMove(mv)
	board.alterPosition(mv)
	// Check if move results in king in check
	return !inCheck(kingPc, otherPieces)
}

func isCastlingValid(bm boardMove, otherPieces map[int][]*piece) bool {
	for _, pieces := range otherPieces {
		for _, piece := range pieces {
			if piece.captured {
				continue
			}
			moves := piece.moveGenerator(piece)
			for _, move := range moves {
				if move.To == bm.From ||
					move.To == bm.To ||
					move.To == bm.castlingTo {
					return false
				}
			}
		}
	}
	return true
}

func inCheck(kingPc *piece, otherPieces map[int][]*piece) bool {
	brdIndex := int(math.Log2(float64(kingPc.position)))
	for _, pieces := range otherPieces {
		for _, piece := range pieces {
			if piece.captured {
				continue
			}
			moves := piece.moveGenerator(piece)
			for _, move := range moves {
				if move.To == brdIndex {
					return true
				}
			}
		}
	}
	return false
}

// GameStatus returns status of the current game
func GameStatus(myTurn bool) int {
	if myTurn {
		if legalMoves(myTurn) > 0 {
			return InProgress
		}
		if inCheck(game.myPieces[king][0], game.otherPieces) {
			return Win
		}
		return Stalemate
	}
	if legalMoves(myTurn) > 0 {
		return InProgress
	}
	if inCheck(game.otherPieces[king][0], game.myPieces) {
		return Lost
	}
	return Stalemate
}

// IsPromotion tells whether given move corresponds to pawn promotion
func IsPromotion(mv boardMove) bool {
	piece := game.board.pieces[mv.From]
	if piece == nil || piece.id != pawn {
		return false
	}
	rank, _ := getRankFile(mv.From)
	if piece.color == white && rank == 7 {
		return true
	}
	if piece.color == black && rank == 2 {
		return true
	}
	return false
}

// PrintBoard prints entire board state. Useful for testing.
func PrintBoard() {
	names := make([]string, 0)
	lines := make([]string, 0)
	for i, pc := range game.board.pieces {
		if pc == nil {
			names = append(names, "-")
		} else {
			switch pc.id {
			case king:
				names = append(names, "K")
			case queen:
				names = append(names, "Q")
			case rook:
				names = append(names, "R")
			case bishop:
				names = append(names, "B")
			case knight:
				names = append(names, "N")
			case pawn:
				names = append(names, "P")
			}
		}
		if (i+1)%8 == 0 {
			nameStr := strings.Join(names, " ")
			lines = append(lines, "| "+nameStr+" |")
			names = make([]string, 0)
		}
	}
	fmt.Println("___________________")
	for i := len(lines) - 1; i >= 0; i-- {
		fmt.Println(lines[i])
	}
	fmt.Println("___________________")
}
