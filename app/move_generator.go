package app

import (
	"math"
)

func getRankFile(boardIndex int) (int, int) {
	boardIndex += 1
	rank := int(math.Ceil(float64(boardIndex) / 8.0))
	file := boardIndex % 8
	if file == 0 {
		// Remainder will be zero for last cell in a rank
		file = 8
	}
	return rank, file
}

func kingMoves(piece *piece, moveCh chan boardMove) {
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	// TODO: handle castling
	if file >= 2 { // Left
		moveCh <- boardMove{From: pos, To: pos - 1}
	}
	if file <= 7 { // Right
		moveCh <- boardMove{From: pos, To: pos + 1}
	}
	if rank <= 7 { // Up
		moveCh <- boardMove{From: pos, To: pos + 8}
	}
	if rank >= 2 { // Down
		moveCh <- boardMove{From: pos, To: pos - 8}
	}
	if file >= 2 && rank <= 7 { // Upper diagonal left
		moveCh <- boardMove{From: pos, To: pos + 7}
	}
	if file <= 7 && rank <= 7 { // Upper diagonal right
		moveCh <- boardMove{From: pos, To: pos + 9}
	}
	if file >= 2 && rank >= 2 { // Lower diagonal left
		moveCh <- boardMove{From: pos, To: pos - 9}
	}
	if file <= 7 && rank >= 2 { // Lower diagonal right
		moveCh <- boardMove{From: pos, To: pos - 7}
	}
}

func queenMoves(piece *piece, moveCh chan boardMove) {
	go rookMoves(piece, moveCh)
	go bishopMoves(piece, moveCh)
}

func rookMoves(piece *piece, moveCh chan boardMove) {
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	// Traverse file upwards until obstructed
	for p, r := pos+8, rank+1; r <= 8; {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p += 8
		r += 1
	}
	// Traverse file downwards until obstructed
	for p, r := pos-8, rank-1; r >= 1; {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p -= 8
		r -= 1
	}
	// Traverse rank left until obstructed
	for p, f := pos-1, file-1; f >= 1; {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p -= 1
		f -= 1
	}
	// Traverse rank right until obstructed
	for p, f := pos+1, file+1; f <= 8; {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p += 1
		f += 1
	}
}

func bishopMoves(piece *piece, moveCh chan boardMove) {
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	// Traverse top right diagonal until obstructed
	p, f, r := pos+8+1, file+1, rank+1
	for f <= 8 && r <= 8 {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p += 8 + 1
		f += 1
		r += 1
	}
	// Traverse bottom right diagonal until obstructed
	p, f, r = pos-8+1, file+1, rank-1
	for f <= 8 && r >= 1 {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p = p - 8 + 1
		f += 1
		r -= 1
	}
	// Traverse bottom left diagonal until obstructed
	p, f, r = pos-8-1, file-1, rank-1
	for f >= 1 && r >= 1 {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p = p - 8 - 1
		f -= 1
		r -= 1
	}
	// Traverse top left diagonal until obstructed
	p, f, r = pos+8-1, file-1, rank+1
	for f >= 1 && r <= 8 {
		moveCh <- boardMove{From: pos, To: p}
		if pieces[p] != nil {
			break
		}
		p += 8 - 1
		f -= 1
		r += 1
	}
}

func knightMoves(piece *piece, moveCh chan boardMove) {
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	// Clock wise moves from top right
	if file <= 7 && rank <= 6 {
		moveCh <- boardMove{From: pos, To: pos + 8 + 9}
	}
	if file <= 6 && rank <= 7 {
		moveCh <- boardMove{From: pos, To: pos + 8 + 2}
	}
	if file <= 6 && rank >= 2 {
		moveCh <- boardMove{From: pos, To: pos - 6}
	}
	if file <= 7 && rank >= 3 {
		moveCh <- boardMove{From: pos, To: pos - 8 - 7}
	}
	if file >= 2 && rank >= 3 {
		moveCh <- boardMove{From: pos, To: pos - 8 - 9}
	}
	if file >= 3 && rank >= 2 {
		moveCh <- boardMove{From: pos, To: pos - 8 - 2}
	}
	if file >= 3 && rank <= 7 {
		moveCh <- boardMove{From: pos, To: pos + 6}
	}
	if file >= 2 && rank <= 6 {
		moveCh <- boardMove{From: pos, To: pos + 8 + 7}
	}
}

func pawnMoves(piece *piece, moveCh chan boardMove) {
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	// TODO: handle en-passant
	pieces := game.board.pieces
	if piece.color == white {
		if rank < 8 && pieces[pos+8] == nil { // Up
			moveCh <- boardMove{From: pos, To: pos + 8}
		}
		if file > 1 && rank < 8 && pieces[pos+7] != nil {
			// Upper diagonal left
			moveCh <- boardMove{From: pos, To: pos + 7}
		}
		if file < 8 && rank < 8 && pieces[pos+9] != nil {
			// Upper diagonal right
			moveCh <- boardMove{From: pos, To: pos + 9}
		}
		return
	}
	// Color is black so we need to move reverse in terms of board numbering
	if rank >= 2 && pieces[pos-8] == nil { // Down
		moveCh <- boardMove{From: pos, To: pos - 8}
	}
	if file >= 2 && rank >= 2 && pieces[pos-9] != nil {
		// Lower diagonal left
		moveCh <- boardMove{From: pos, To: pos - 9}
	}
	if file <= 7 && rank >= 2 && pieces[pos-7] != nil {
		// Lower diagonal right
		moveCh <- boardMove{From: pos, To: pos - 7}
	}
}
