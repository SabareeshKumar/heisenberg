package app

import (
	"math"
)

func pawnStructure(myTurn bool) (isolated, doubled, blocked int) {
	pieces := game.board.pieces
	var pawns []*piece
	if myTurn {
		pawns = game.myPieces[pawn]
	} else {
		pawns = game.otherPieces[pawn]
	}
	columns := make([]int, 8)
	for _, piece := range pawns {
		if piece.captured {
			continue
		}
		brdIndex := int(math.Log2(float64(piece.position)))
		rank, file := getRankFile(brdIndex)
		columns[file-1] += 1
		var blockedPos int
		if myTurn && rank < 8 {
			blockedPos = 1 << (brdIndex + 8)
		} else if !myTurn && rank > 1 {
			blockedPos = 1 << (brdIndex - 8)
		} else {
			continue
		}
		blockedPiece := pieces[blockedPos]
		if blockedPiece == nil {
			continue
		}
		if blockedPiece.color == piece.color {
			doubled += 1
		} else {
			blocked += 1
		}
	}
	for _, count := range columns {
		if count == 1 {
			isolated += 1
		}
	}
	return
}

func legalMoves(myTurn bool) int {
	var pieceMap map[int][]*piece
	if myTurn {
		pieceMap = game.myPieces
	} else {
		pieceMap = game.otherPieces
	}
	moveCount := 0
	for _, pieces := range pieceMap {
		for _, piece := range pieces {
			if piece.captured {
				continue
			}
			for _, move := range piece.moveGenerator(piece) {
				if isMoveLegal(move) {
					moveCount += 1
				}
			}
		}
	}
	return moveCount
}

func eval(myTurn bool) float32 {
	moveCount := legalMoves(myTurn)
	if moveCount == 0 {
		return float32(math.MinInt32)
	}
	_moveCount := legalMoves(!myTurn)
	if _moveCount == 0 {
		return float32(math.MaxInt32)
	}
	isolated, doubled, blocked := pawnStructure(myTurn)
	_isolated, _doubled, _blocked := pawnStructure(!myTurn)
	score := float32(game.materialBalance)
	score -= 0.5 * float32(isolated - _isolated + doubled - _doubled +
		blocked - _blocked)
	score += 0.1 * float32(moveCount - _moveCount)
	return score
}
