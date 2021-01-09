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

func newBoardMove(from, to, castlingFrom, castlingTo int) boardMove {
	return boardMove{from, to, castlingFrom, castlingTo, game.board.pieces[to]}
}

func canCastle(kingFrom, kingTo int) bool {
	pieces := game.board.pieces
	if pieces[kingFrom] == nil || pieces[kingFrom].moveCount > 1 {
		// king already moved
		return false
	}
	if kingTo > kingFrom {
		// king side castling
		rookPos := kingFrom + 3
		if pieces[rookPos] == nil || pieces[rookPos].moveCount > 1 {
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
	if pieces[rookPos] == nil || pieces[rookPos].moveCount > 1 {
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

func kingMoves(piece *piece) []boardMove {
	moves := make([]boardMove, 0)
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	if file >= 2 { // Left
		moves = append(moves, newBoardMove(pos, pos-1, -1, -1))
	}
	if file <= 7 { // Right
		moves = append(moves, newBoardMove(pos, pos+1, -1, -1))
	}
	if rank <= 7 { // Up
		moves = append(moves, newBoardMove(pos, pos+8, -1, -1))
	}
	if rank >= 2 { // Down
		moves = append(moves, newBoardMove(pos, pos-8, -1, -1))
	}
	if file >= 2 && rank <= 7 { // Upper diagonal left
		moves = append(moves, newBoardMove(pos, pos+7, -1, -1))
	}
	if file <= 7 && rank <= 7 { // Upper diagonal right
		moves = append(moves, newBoardMove(pos, pos+9, -1, -1))
	}
	if file >= 2 && rank >= 2 { // Lower diagonal left
		moves = append(moves, newBoardMove(pos, pos-9, -1, -1))
	}
	if file <= 7 && rank >= 2 { // Lower diagonal right
		moves = append(moves, newBoardMove(pos, pos-7, -1, -1))
	}
	if piece.moveCount > 0 {
		return moves
	}
	// Check king side castling
	if canCastle(pos, pos+2) {
		moves = append(moves, newBoardMove(pos, pos+2, pos+3, pos+1))
	}
	// Check queen side castling
	if canCastle(pos, pos-2) {
		moves = append(moves, newBoardMove(pos, pos-2, pos-4, pos-1))
	}
	return moves
}

func queenMoves(piece *piece) []boardMove {
	return append(rookMoves(piece), bishopMoves(piece)...)
}

func rookMoves(piece *piece) []boardMove {
	moves := make([]boardMove, 0)
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	// Traverse file upwards until obstructed
	for p, r := pos+8, rank+1; r <= 8; {
		moves = append(moves, newBoardMove(pos, p, -1, -1))
		if pieces[p] != nil {
			break
		}
		p += 8
		r += 1
	}
	// Traverse file downwards until obstructed
	for p, r := pos-8, rank-1; r >= 1; {
		moves = append(moves, newBoardMove(pos, p, -1, -1))
		if pieces[p] != nil {
			break
		}
		p -= 8
		r -= 1
	}
	// Traverse rank left until obstructed
	for p, f := pos-1, file-1; f >= 1; {
		moves = append(moves, newBoardMove(pos, p, -1, -1))
		if pieces[p] != nil {
			break
		}
		p -= 1
		f -= 1
	}
	// Traverse rank right until obstructed
	for p, f := pos+1, file+1; f <= 8; {
		moves = append(moves, newBoardMove(pos, p, -1, -1))
		if pieces[p] != nil {
			break
		}
		p += 1
		f += 1
	}
	return moves
}

func bishopMoves(piece *piece) []boardMove {
	moves := make([]boardMove, 0)
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	// Traverse top right diagonal until obstructed
	p, f, r := pos+8+1, file+1, rank+1
	for f <= 8 && r <= 8 {
		moves = append(moves, newBoardMove(pos, p, -1, -1))
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
		moves = append(moves, newBoardMove(pos, p, -1, -1))
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
		moves = append(moves, newBoardMove(pos, p, -1, -1))
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
		moves = append(moves, newBoardMove(pos, p, -1, -1))
		if pieces[p] != nil {
			break
		}
		p += 8 - 1
		f -= 1
		r += 1
	}
	return moves
}

func knightMoves(piece *piece) []boardMove {
	moves := make([]boardMove, 0)
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	// Clock wise moves from top right
	if file <= 7 && rank <= 6 {
		moves = append(moves, newBoardMove(pos, pos+8+9, -1, -1))
	}
	if file <= 6 && rank <= 7 {
		moves = append(moves, newBoardMove(pos, pos+8+2, -1, -1))
	}
	if file <= 6 && rank >= 2 {
		moves = append(moves, newBoardMove(pos, pos-6, -1, -1))
	}
	if file <= 7 && rank >= 3 {
		moves = append(moves, newBoardMove(pos, pos-8-7, -1, -1))
	}
	if file >= 2 && rank >= 3 {
		moves = append(moves, newBoardMove(pos, pos-8-9, -1, -1))
	}
	if file >= 3 && rank >= 2 {
		moves = append(moves, newBoardMove(pos, pos-8-2, -1, -1))
	}
	if file >= 3 && rank <= 7 {
		moves = append(moves, newBoardMove(pos, pos+6, -1, -1))
	}
	if file >= 2 && rank <= 6 {
		moves = append(moves, newBoardMove(pos, pos+8+7, -1, -1))
	}
	return moves
}

func pawnMoves(piece *piece) []boardMove {
	moves := make([]boardMove, 0)
	pos := int(math.Log2(float64(piece.position)))
	rank, file := getRankFile(pos)
	pieces := game.board.pieces
	if piece.color == white {
		if rank <= 7 && pieces[pos+8] == nil { // Up
			moves = append(moves, newBoardMove(pos, pos+8, -1, -1))
			if rank == 2 && pieces[pos+16] == nil { // Up twice
				moves = append(moves, newBoardMove(pos, pos+16, -1, -1))
			}
		}
		if file >= 2 && rank < 8 {
			// Upper diagonal left
			if pieces[pos+7] != nil {
				moves = append(moves, newBoardMove(pos, pos+7, -1, -1))
			} else {
				sidePc := pieces[pos-1]
				if sidePc != nil && sidePc.enpassantMove == game.moveCount {
					moves = append(moves, newBoardMove(pos, pos+7, -1, -1))
				}
			}
		}
		if file <= 7 && rank < 8 {
			// Upper diagonal right
			if pieces[pos+9] != nil {
				moves = append(moves, newBoardMove(pos, pos+9, -1, -1))
			} else {
				sidePc := pieces[pos+1]
				if sidePc != nil && sidePc.enpassantMove == game.moveCount {
					moves = append(moves, newBoardMove(pos, pos+9, -1, -1))
				}
			}
		}
		return moves
	}
	// Color is black so we need to move reverse in terms of board numbering
	if rank >= 2 && pieces[pos-8] == nil { // Down
		moves = append(moves, newBoardMove(pos, pos-8, -1, -1))
		if rank == 7 && pieces[pos-16] == nil { // Down twice
			moves = append(moves, newBoardMove(pos, pos-16, -1, -1))
		}
	}
	if file >= 2 && rank >= 2 {
		// Lower diagonal left
		if pieces[pos-9] != nil {
			moves = append(moves, newBoardMove(pos, pos-9, -1, -1))
		} else {
			sidePc := pieces[pos-1]
			if sidePc != nil && sidePc.enpassantMove == game.moveCount {
				moves = append(moves, newBoardMove(pos, pos-9, -1, -1))
			}
		}
	}
	if file <= 7 && rank >= 2 {
		// Lower diagonal right
		if pieces[pos-7] != nil {
			moves = append(moves, newBoardMove(pos, pos-7, -1, -1))
		} else {
			sidePc := pieces[pos+1]
			if sidePc != nil && sidePc.enpassantMove == game.moveCount {
				moves = append(moves, newBoardMove(pos, pos-7, -1, -1))
			}
		}
	}
	return moves
}
