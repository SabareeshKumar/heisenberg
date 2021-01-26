package app

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

var evaluationsPerSearch = 0
var tableHits = 0

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
	tpnTbl          map[int64]tpnMeta // Transposition table
	lastMove        *boardMove
}

var game GameState

// Hash keys for each white piece and square
var blackHashKeys map[int][]int64

// Hash keys for each black piece and square
var whiteHashKeys map[int][]int64

// Hash keys for each turn
var turnHash map[bool]int64

// Hash keys for each file. Will be used to generate hashes for enpassant
// squares
var fileHash map[int]int64

// Hash keys for king & queen side castling rights of computer
var myCastlingHash []int64

// Hash keys for king & queen side castling rights of opponent
var otherCastlingHash []int64

// InitGame sets up a new game.
func InitGame(userColorChoice int) error {
	whitePieces := make(map[int][]*piece, 16)
	blackPieces := make(map[int][]*piece, 16)
	board := newBoard()
	game = GameState{}
	game.board = board
	game.moveCount = 0
	game.tpnTbl = make(map[int64]tpnMeta, tpnTblSize)
	for i := 0; i <= 15; i++ {
		piece := board.pieces[i]
		whitePieces[piece.id] = append(whitePieces[piece.id], piece)
	}
	for i := 48; i <= 63; i++ {
		piece := board.pieces[i]
		blackPieces[piece.id] = append(blackPieces[piece.id], piece)
	}
	if userColorChoice == white {
		game.myColor = black
		game.otherPieces = whitePieces
		game.myPieces = blackPieces
	} else {
		game.myColor = white
		game.otherPieces = blackPieces
		game.myPieces = whitePieces
	}
	return loadOpeningBook(userColorChoice)
}

// MakeMove verifies if the user move if valid
func MakeMove(uMove boardMove) error {
	board := game.board
	piece := board.pieces[uMove.From]
	if piece.color == game.myColor {
		return errors.New("Cannot move opponent piece")
	}
	valid := false
	moves, _ := piece.moveGenerator(piece)
	for _, move := range moves {
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
	game.lastMove = &uMove
	return board.alterPosition(uMove)
}

// MyMove computes a move for the engine
func MyMove() (UserMove, error) {
	evaluationsPerSearch = 0
	tableHits = 0
	myMov, _ := search(
		true, float32(math.MaxInt32), 1, []boardMove{}, game.lastMove)
	if debugMode {
		fmt.Printf("Evaluated %d board states:\n", evaluationsPerSearch)
		fmt.Println("Number of table hits:", tableHits)
		fmt.Println("Hash table length:", len(game.tpnTbl))
		fmt.Println("Move count:", game.moveCount)
	}
	err := game.board.alterPosition(myMov)
	if err != nil {
		return UserMove{}, err
	}
	myMoveCoord, err := myMov.toUserMove()
	if err != nil {
		return UserMove{}, err
	}
	game.lastMove = &myMov
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

func isMoveValid(mv boardMove) bool {
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
	var otherPieces map[int][]*piece
	if pc.color == game.myColor {
		otherPieces = game.otherPieces
	} else {
		otherPieces = game.myPieces
	}
	if mv.castlingFrom > 0 {
		// Check if castling is valid. It inherently checks if move
		// results in check.
		return isCastlingValid(mv, otherPieces)
	}
	return true
}

func isCastlingValid(bm boardMove, otherPieces map[int][]*piece) bool {
	for _, pieces := range otherPieces {
		for _, piece := range pieces {
			if piece.captured {
				continue
			}
			moves, _ := piece.moveGenerator(piece)
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

func inCheckSimple(myTurn bool, attacks uint) bool {
	var kingPc *piece
	if myTurn {
		kingPc = game.myPieces[king][0]
	} else {
		kingPc = game.otherPieces[king][0]
	}
	return (1 << kingPc.position & attacks) != 0
}

func inCheck(kingPc *piece, otherPieces map[int][]*piece) bool {
	brdIndex := kingPc.position
	for _, pieces := range otherPieces {
		for _, piece := range pieces {
			if piece.captured {
				continue
			}
			moves, _ := piece.moveGenerator(piece)
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
	if legalMoves(myTurn) > 0 {
		return InProgress
	}
	if myTurn {
		if inCheck(game.myPieces[king][0], game.otherPieces) {
			return Win
		}
		return Stalemate
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

// CreateHashKeys creates all hash keys needed throughout the game.
func CreateHashKeys() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	encounteredKeys := make(map[int64]bool)
	// Create keys for white pieces
	whiteHashKeys = make(map[int][]int64)
	for pieceId := range whiteMeta {
		squares := make([]int64, 64)
		for i := 0; i < 64; {
			newKey := r.Int63()
			if encounteredKeys[newKey] {
				continue
			}
			encounteredKeys[newKey] = true
			squares[i] = newKey
			i++
		}
		whiteHashKeys[pieceId] = squares
	}
	// Create keys for black pieces
	blackHashKeys = make(map[int][]int64)
	for pieceId := range blackMeta {
		squares := make([]int64, 64)
		for i := 0; i < 64; {
			newKey := r.Int63()
			if encounteredKeys[newKey] {
				continue
			}
			encounteredKeys[newKey] = true
			squares[i] = newKey
			i++
		}
		blackHashKeys[pieceId] = squares
	}
	turnHash = make(map[bool]int64)
	for {
		myTurnKey := r.Int63()
		if encounteredKeys[myTurnKey] {
			continue
		}
		encounteredKeys[myTurnKey] = true
		turnHash[true] = myTurnKey
		break
	}
	for {
		opponentTurnKey := r.Int63()
		if encounteredKeys[opponentTurnKey] {
			continue
		}
		encounteredKeys[opponentTurnKey] = true
		turnHash[false] = opponentTurnKey
		break
	}
	fileHash = make(map[int]int64)
	for i := 1; i <= 8; {
		newKey := r.Int63()
		if encounteredKeys[newKey] {
			continue
		}
		encounteredKeys[newKey] = true
		fileHash[i] = newKey
		i++
	}
	myCastlingHash = make([]int64, 2)
	for i := 0; i < 2; {
		newKey := r.Int63()
		if encounteredKeys[newKey] {
			continue
		}
		encounteredKeys[newKey] = true
		myCastlingHash[i] = newKey
		i++
	}
	otherCastlingHash = make([]int64, 2)
	for i := 0; i < 2; {
		newKey := r.Int63()
		if encounteredKeys[newKey] {
			continue
		}
		encounteredKeys[newKey] = true
		otherCastlingHash[i] = newKey
		i++
	}
}
