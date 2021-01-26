package app

import (
	"errors"
	"math"
)

type boardConfig struct {
	pieces []*piece
	hash   int64
}

func newBoard() *boardConfig {
	pieces := make([]*piece, 64)
	pieces[0] = newWhitePiece(rook, 1<<0)
	pieces[1] = newWhitePiece(knight, 1<<1)
	pieces[2] = newWhitePiece(bishop, 1<<2)
	pieces[3] = newWhitePiece(queen, 1<<3)
	pieces[4] = newWhitePiece(king, 1<<4)
	pieces[5] = newWhitePiece(bishop, 1<<5)
	pieces[6] = newWhitePiece(knight, 1<<6)
	pieces[7] = newWhitePiece(rook, 1<<7)
	for i := 8; i <= 15; i++ {
		pieces[i] = newWhitePiece(pawn, 1<<i)
	}
	for i := 48; i <= 55; i++ {
		pieces[i] = newBlackPiece(pawn, 1<<i)
	}
	pieces[56] = newBlackPiece(rook, 1<<56)
	pieces[57] = newBlackPiece(knight, 1<<57)
	pieces[58] = newBlackPiece(bishop, 1<<58)
	pieces[59] = newBlackPiece(queen, 1<<59)
	pieces[60] = newBlackPiece(king, 1<<60)
	pieces[61] = newBlackPiece(bishop, 1<<61)
	pieces[62] = newBlackPiece(knight, 1<<62)
	pieces[63] = newBlackPiece(rook, 1<<63)
	board := boardConfig{}
	board.pieces = pieces
	for i := 0; i <= 15; i++ {
		board.hash ^= whiteHashKeys[pieces[i].id][i]
	}
	for i := 48; i <= 63; i++ {
		board.hash ^= blackHashKeys[pieces[i].id][i]
	}
	return &board
}

func (brd *boardConfig) alterPosition(bm boardMove) error {
	pc := brd.pieces[bm.From]
	if pc == nil {
		return errors.New("Invalid move")
	}
	var keys map[int][]int64
	if pc.color == white {
		keys = whiteHashKeys
		brd.hash ^= whiteHashKeys[pc.id][bm.From]
		brd.hash ^= whiteHashKeys[pc.id][bm.To]
	} else {
		keys = blackHashKeys
		brd.hash ^= blackHashKeys[pc.id][bm.From]
		brd.hash ^= blackHashKeys[pc.id][bm.To]
	}
	capturedPc := bm.captured
	if capturedPc != nil {
		capturedPc.captured = true
		if capturedPc.color != game.myColor {
			game.materialBalance += weights[capturedPc.id]
		} else {
			game.materialBalance -= weights[capturedPc.id]
		}
		capturedPos := int(math.Log2(float64(capturedPc.position)))
		if capturedPc.color == white {
			brd.hash ^= whiteHashKeys[capturedPc.id][capturedPos]
		} else {
			brd.hash ^= blackHashKeys[capturedPc.id][capturedPos]
		}
		brd.pieces[capturedPos] = nil
	}
	brd.pieces[bm.From] = nil
	brd.pieces[bm.To] = pc
	pc.position = 1 << bm.To
	game.moveCount++
	pc.moveCount++
	if bm.castlingFrom != -1 {
		rookPc := brd.pieces[bm.castlingFrom]
		brd.pieces[bm.castlingFrom] = nil
		brd.pieces[bm.castlingTo] = rookPc
		rookPc.position = 1 << bm.castlingTo
		rookPc.moveCount++
		brd.hash ^= keys[rookPc.id][bm.castlingFrom]
		brd.hash ^= keys[rookPc.id][bm.castlingTo]
		return nil
	}
	if pc.id == pawn && int(math.Abs(float64(bm.To-bm.From))) == 16 {
		pc.enpassantMove = game.moveCount
		return nil
	}
	if bm.PromotedPc <= 0 {
		return nil
	}
	var newPc *piece
	if pc.color == black {
		newPc = newBlackPiece(bm.PromotedPc, pc.position)
	} else {
		newPc = newWhitePiece(bm.PromotedPc, pc.position)
	}
	brd.pieces[bm.To] = newPc
	newPc.promotedBy = pc
	pc.captured = true
	brd.hash ^= keys[pc.id][bm.To]
	brd.hash ^= keys[newPc.id][bm.To]
	if pc.color == game.myColor {
		game.materialBalance += weights[bm.PromotedPc]
		game.materialBalance -= weights[pc.id]
		game.myPieces[newPc.id] = append(game.myPieces[newPc.id], newPc)
		return nil
	}
	game.materialBalance -= weights[bm.PromotedPc]
	game.materialBalance += weights[pc.id]
	game.otherPieces[newPc.id] = append(game.otherPieces[newPc.id], newPc)
	return nil
}

func (brd *boardConfig) undoMove(bm boardMove) {
	game.moveCount--
	pc := brd.pieces[bm.To]
	brd.pieces[bm.From] = pc
	var keys map[int][]int64
	if pc.color == white {
		keys = whiteHashKeys
		brd.hash ^= whiteHashKeys[pc.id][bm.From]
		brd.hash ^= whiteHashKeys[pc.id][bm.To]
	} else {
		keys = blackHashKeys
		brd.hash ^= blackHashKeys[pc.id][bm.From]
		brd.hash ^= blackHashKeys[pc.id][bm.To]
	}
	capturedPc := bm.captured
	capturedPos := bm.To
	if capturedPc != nil {
		capturedPc.captured = false
		if capturedPc.color != game.myColor {
			game.materialBalance -= weights[capturedPc.id]
		} else {
			game.materialBalance += weights[capturedPc.id]
		}
		capturedPos = int(math.Log2(float64(capturedPc.position)))
		if capturedPc.color == white {
			brd.hash ^= whiteHashKeys[capturedPc.id][capturedPos]
		} else {
			brd.hash ^= blackHashKeys[capturedPc.id][capturedPos]
		}
	}
	brd.pieces[bm.To] = nil
	brd.pieces[capturedPos] = capturedPc
	pc.position = 1 << bm.From
	pc.moveCount--
	if bm.castlingFrom != -1 {
		rookPc := brd.pieces[bm.castlingTo]
		brd.pieces[bm.castlingFrom] = rookPc
		brd.pieces[bm.castlingTo] = nil
		rookPc.position = 1 << bm.castlingFrom
		rookPc.moveCount--
		brd.hash ^= keys[rookPc.id][bm.castlingFrom]
		brd.hash ^= keys[rookPc.id][bm.castlingTo]
		return
	}
	if pc.id == pawn && int(math.Abs(float64(bm.To-bm.From))) == 16 {
		pc.enpassantMove = -1
		return
	}
	pwn := pc.promotedBy
	// Check if move to be undone is a promotion move
	if pwn == nil || pc.moveCount > 0 {
		return
	}
	brd.pieces[bm.From] = pwn
	pwn.position = pc.position
	pwn.moveCount--
	pwn.captured = false
	brd.hash ^= keys[pc.id][bm.From]
	brd.hash ^= keys[pwn.id][bm.From]
	if pwn.color == game.myColor {
		game.materialBalance += weights[pwn.id]
		game.materialBalance -= weights[pc.id]
		// Remove promoted piece from list
		pieces := game.myPieces[pc.id]
		game.myPieces[pc.id] = pieces[:len(pieces)-1]
		return
	}
	game.materialBalance -= weights[pwn.id]
	game.materialBalance += weights[pc.id]
	// Remove promoted piece from list
	pieces := game.otherPieces[pc.id]
	game.otherPieces[pc.id] = pieces[:len(pieces)-1]

}
