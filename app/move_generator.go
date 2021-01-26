package app

import (
	"math"
)

func getRankFile(boardIndex int) (int, int) {
	boardIndex++
	rank := int(math.Ceil(float64(boardIndex) / 8.0))
	file := boardIndex % 8
	if file == 0 {
		// Remainder will be zero for last cell in a rank
		file = 8
	}
	return rank, file
}

func newBoardMove(
	from, to, castlingFrom, castlingTo, promotedPc int,
	attacks *uint) boardMove {
	*attacks |= 1 << to
	return boardMove{
		From:         from,
		To:           to,
		castlingFrom: castlingFrom,
		castlingTo:   castlingTo,
		captured:     game.board.pieces[to],
		PromotedPc:   promotedPc,
	}
}

func canCastle(kingFrom, kingTo int) bool {
	pieces := game.board.pieces
	if pieces[kingFrom] == nil || pieces[kingFrom].moveCount > 0 {
		// king already moved
		return false
	}
	if kingTo > kingFrom {
		// king side castling
		rookPos := kingFrom + 3
		if pieces[rookPos] == nil || pieces[rookPos].moveCount > 0 {
			// rook moved
			return false
		}
		if pieces[kingFrom+1] != nil || pieces[kingFrom+2] != nil {
			// pieces in between
			return false
		}
		return true
	}
	// queen side castling
	rookPos := kingFrom - 4
	if pieces[rookPos] == nil || pieces[rookPos].moveCount > 0 {
		// rook moved
		return false
	}
	if pieces[kingFrom-1] != nil || pieces[kingFrom-2] != nil ||
		pieces[kingFrom-3] != nil {
		// pieces in between
		return false
	}
	return true

}

func kingMoves(piece *piece) ([]boardMove, uint) {
	moves := make([]boardMove, 0)
	var attacks uint
	pos := piece.position
	rank, file := getRankFile(pos)
	if file >= 2 { // Left
		moves = append(moves, newBoardMove(pos, pos-1, -1, -1, -1, &attacks))
	}
	if file <= 7 { // Right
		moves = append(moves, newBoardMove(pos, pos+1, -1, -1, -1, &attacks))
	}
	if rank <= 7 { // Up
		moves = append(moves, newBoardMove(pos, pos+8, -1, -1, -1, &attacks))
	}
	if rank >= 2 { // Down
		moves = append(moves, newBoardMove(pos, pos-8, -1, -1, -1, &attacks))
	}
	if file >= 2 && rank <= 7 { // Upper diagonal left
		moves = append(moves, newBoardMove(pos, pos+7, -1, -1, -1, &attacks))
	}
	if file <= 7 && rank <= 7 { // Upper diagonal right
		moves = append(moves, newBoardMove(pos, pos+9, -1, -1, -1, &attacks))
	}
	if file >= 2 && rank >= 2 { // Lower diagonal left
		moves = append(moves, newBoardMove(pos, pos-9, -1, -1, -1, &attacks))
	}
	if file <= 7 && rank >= 2 { // Lower diagonal right
		moves = append(moves, newBoardMove(pos, pos-7, -1, -1, -1, &attacks))
	}
	if piece.moveCount > 0 {
		return moves, attacks
	}
	// Check king side castling
	if canCastle(pos, pos+2) {
		moves = append(
			moves, newBoardMove(pos, pos+2, pos+3, pos+1, -1, &attacks))
	}
	// Check queen side castling
	if canCastle(pos, pos-2) {
		moves = append(
			moves, newBoardMove(pos, pos-2, pos-4, pos-1, -1, &attacks))
	}
	return moves, attacks
}

func queenMoves(piece *piece) ([]boardMove, uint) {
	rmoves, rattacks := rookMoves(piece)
	bmoves, battacks := bishopMoves(piece)
	return append(rmoves, bmoves...), rattacks | battacks
}

func rookMoves(piece *piece) ([]boardMove, uint) {
	moves := make([]boardMove, 0)
	var attacks uint
	pos := piece.position
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	// Traverse file upwards until obstructed
	for p, r := pos+8, rank+1; r <= 8; {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p += 8
		r++
	}
	// Traverse file downwards until obstructed
	for p, r := pos-8, rank-1; r >= 1; {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p -= 8
		r--
	}
	// Traverse rank left until obstructed
	for p, f := pos-1, file-1; f >= 1; {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p--
		f--
	}
	// Traverse rank right until obstructed
	for p, f := pos+1, file+1; f <= 8; {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p++
		f++
	}
	return moves, attacks
}

func bishopMoves(piece *piece) ([]boardMove, uint) {
	moves := make([]boardMove, 0)
	var attacks uint
	pos := piece.position
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	// Traverse top right diagonal until obstructed
	p, f, r := pos+8+1, file+1, rank+1
	for f <= 8 && r <= 8 {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p += 8 + 1
		f++
		r++
	}
	// Traverse bottom right diagonal until obstructed
	p, f, r = pos-8+1, file+1, rank-1
	for f <= 8 && r >= 1 {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p = p - 8 + 1
		f++
		r--
	}
	// Traverse bottom left diagonal until obstructed
	p, f, r = pos-8-1, file-1, rank-1
	for f >= 1 && r >= 1 {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p = p - 8 - 1
		f--
		r--
	}
	// Traverse top left diagonal until obstructed
	p, f, r = pos+8-1, file-1, rank+1
	for f >= 1 && r <= 8 {
		moves = append(moves, newBoardMove(pos, p, -1, -1, -1, &attacks))
		if pieces[p] != nil {
			break
		}
		p += 8 - 1
		f--
		r++
	}
	return moves, attacks
}

func knightMoves(piece *piece) ([]boardMove, uint) {
	moves := make([]boardMove, 0)
	var attacks uint
	pos := piece.position
	rank, file := getRankFile(pos)
	// Clock wise moves from top right
	if file <= 7 && rank <= 6 {
		moves = append(moves, newBoardMove(pos, pos+8+9, -1, -1, -1, &attacks))
	}
	if file <= 6 && rank <= 7 {
		moves = append(moves, newBoardMove(pos, pos+8+2, -1, -1, -1, &attacks))
	}
	if file <= 6 && rank >= 2 {
		moves = append(moves, newBoardMove(pos, pos-6, -1, -1, -1, &attacks))
	}
	if file <= 7 && rank >= 3 {
		moves = append(moves, newBoardMove(pos, pos-8-7, -1, -1, -1, &attacks))
	}
	if file >= 2 && rank >= 3 {
		moves = append(moves, newBoardMove(pos, pos-8-9, -1, -1, -1, &attacks))
	}
	if file >= 3 && rank >= 2 {
		moves = append(moves, newBoardMove(pos, pos-8-2, -1, -1, -1, &attacks))
	}
	if file >= 3 && rank <= 7 {
		moves = append(moves, newBoardMove(pos, pos+6, -1, -1, -1, &attacks))
	}
	if file >= 2 && rank <= 6 {
		moves = append(moves, newBoardMove(pos, pos+8+7, -1, -1, -1, &attacks))
	}
	return moves, attacks
}

func whitePromotionMoves(piece *piece, file, pos int) ([]boardMove, uint) {
	moves := make([]boardMove, 0)
	var attacks uint
	pieces := game.board.pieces
	if pieces[pos+8] == nil { // Up
		for _, promPiece := range promotablePieces {
			moves = append(
				moves, newBoardMove(pos, pos+8, -1, -1, promPiece, &attacks))
		}
	}
	if file >= 2 && pieces[pos+7] != nil { // Upper diagonal left
		for _, promPiece := range promotablePieces {
			moves = append(
				moves, newBoardMove(pos, pos+7, -1, -1, promPiece, &attacks))
		}
	}
	if file <= 7 && pieces[pos+9] != nil { // Upper diagonal right
		for _, promPiece := range promotablePieces {
			moves = append(
				moves, newBoardMove(pos, pos+9, -1, -1, promPiece, &attacks))
		}
	}
	return moves, attacks
}

func whitePawnMoves(piece *piece) ([]boardMove, uint) {
	pos := piece.position
	rank, file := getRankFile(pos)
	if rank == 7 {
		return whitePromotionMoves(piece, file, pos)
	}
	moves := make([]boardMove, 0)
	var attacks uint
	pieces := game.board.pieces
	if rank <= 7 && pieces[pos+8] == nil { // Up
		moves = append(moves, newBoardMove(pos, pos+8, -1, -1, -1, &attacks))
		if rank == 2 && pieces[pos+16] == nil { // Up twice
			moves = append(
				moves, newBoardMove(pos, pos+16, -1, -1, -1, &attacks))
		}
	}
	if file >= 2 && rank < 8 {
		// Upper diagonal left
		if pieces[pos+7] != nil {
			moves = append(
				moves, newBoardMove(pos, pos+7, -1, -1, -1, &attacks))
		} else {
			sidePc := pieces[pos-1]
			if sidePc != nil && sidePc.enpassantMove == game.moveCount {
				bm := newBoardMove(pos, pos+7, -1, -1, -1, &attacks)
				bm.captured = sidePc
				moves = append(moves, bm)
			}
		}
	}
	if file <= 7 && rank < 8 {
		// Upper diagonal right
		if pieces[pos+9] != nil {
			moves = append(
				moves, newBoardMove(pos, pos+9, -1, -1, -1, &attacks))
		} else {
			sidePc := pieces[pos+1]
			if sidePc != nil && sidePc.enpassantMove == game.moveCount {
				bm := newBoardMove(pos, pos+9, -1, -1, -1, &attacks)
				bm.captured = sidePc
				moves = append(moves, bm)
			}
		}
	}
	return moves, attacks
}

func blackPromotionMoves(piece *piece, file, pos int) ([]boardMove, uint) {
	moves := make([]boardMove, 0)
	var attacks uint
	pieces := game.board.pieces
	// Color is black so we need to move reverse in terms of board numbering
	if pieces[pos-8] == nil { // Down
		for _, promPiece := range promotablePieces {
			moves = append(
				moves, newBoardMove(pos, pos-8, -1, -1, promPiece, &attacks))
		}
	}
	if file >= 2 && pieces[pos-9] != nil { // Lower diagonal left
		for _, promPiece := range promotablePieces {
			moves = append(
				moves, newBoardMove(pos, pos-9, -1, -1, promPiece, &attacks))
		}
	}
	if file <= 7 && pieces[pos-7] != nil { // Lower diagonal right
		for _, promPiece := range promotablePieces {
			moves = append(
				moves, newBoardMove(pos, pos-7, -1, -1, promPiece, &attacks))
		}
	}
	return moves, attacks
}

func blackPawnMoves(piece *piece) ([]boardMove, uint) {
	pos := piece.position
	rank, file := getRankFile(pos)
	if rank == 2 {
		return blackPromotionMoves(piece, file, pos)
	}
	moves := make([]boardMove, 0)
	var attacks uint
	pieces := game.board.pieces
	// Color is black so we need to move reverse in terms of board numbering
	if rank >= 2 && pieces[pos-8] == nil { // Down
		moves = append(moves, newBoardMove(pos, pos-8, -1, -1, -1, &attacks))
		if rank == 7 && pieces[pos-16] == nil { // Down twice
			moves = append(
				moves, newBoardMove(pos, pos-16, -1, -1, -1, &attacks))
		}
	}
	if file >= 2 && rank >= 2 {
		// Lower diagonal left
		if pieces[pos-9] != nil {
			moves = append(
				moves, newBoardMove(pos, pos-9, -1, -1, -1, &attacks))
		} else {
			sidePc := pieces[pos-1]
			if sidePc != nil && sidePc.enpassantMove == game.moveCount {
				bm := newBoardMove(pos, pos-9, -1, -1, -1, &attacks)
				bm.captured = sidePc
				moves = append(moves, bm)
			}
		}
	}
	if file <= 7 && rank >= 2 {
		// Lower diagonal right
		if pieces[pos-7] != nil {
			moves = append(
				moves, newBoardMove(pos, pos-7, -1, -1, -1, &attacks))
		} else {
			sidePc := pieces[pos+1]
			if sidePc != nil && sidePc.enpassantMove == game.moveCount {
				bm := newBoardMove(pos, pos-7, -1, -1, -1, &attacks)
				bm.captured = sidePc
				moves = append(moves, bm)
			}
		}
	}
	return moves, attacks
}
