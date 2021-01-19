package app

import (
	"math"
	"sort"
)

type sortInt []boardMove

func (bm sortInt) Len() int {
	return len(bm)
}

func (bm sortInt) Less(i, j int) bool {
	if bm[i].From != bm[j].From {
		return bm[i].From < bm[j].From
	}
	return bm[i].To < bm[j].To
}

func (bm sortInt) Swap(i, j int) {
	bm[i], bm[j] = bm[j], bm[i]
}

func search(
	myTurn bool, bestParentScore float32, depth int,
	moves []boardMove, lastMove *boardMove) (boardMove, float32) {
	hResult := hashedResult(myTurn, maxDepth-depth+1, lastMove)
	if hResult != nil {
		tableHits += 1
		return hResult.bestMove, hResult.score
	}
	if depth > maxDepth {
		evaluationsPerSearch += 1
		move, score := boardMove{}, eval()
		hashResult(myTurn, 0, score, move, lastMove)
		return boardMove{}, eval()
	}
	var bestMove boardMove
	board := game.board
	if depth == 1 {
		moves, _ = generateMoves(myTurn)
	}
	if myTurn {
		maxScore := float32(math.MinInt32)
		defer func() {
			evaluationsPerSearch += 1
			hashResult(myTurn, maxDepth-depth+1, maxScore, bestMove, lastMove)
		}()
		for _, move := range moves {
			if !isMoveValid(move) {
				continue
			}
			board.alterPosition(move)
			opponentMoves, attacks := generateMoves(!myTurn)
			if inCheckSimple(myTurn, attacks) {
				board.undoMove(move)
				continue
			}
			_, score := search(
				!myTurn, maxScore, depth+1, opponentMoves, &move)
			board.undoMove(move)
			if score < maxScore {
				continue
			}
			maxScore = score
			bestMove = move
			if maxScore >= bestParentScore {
				return bestMove, maxScore // alpha-pruning
			}
		}
		return bestMove, maxScore
	}
	minScore := float32(math.MaxInt32)
	defer func() {
		evaluationsPerSearch += 1
		hashResult(myTurn, maxDepth-depth+1, minScore, bestMove, lastMove)
	}()
	for _, move := range moves {
		if !isMoveValid(move) {
			continue
		}
		board.alterPosition(move)
		opponentMoves, attacks := generateMoves(!myTurn)
		if inCheckSimple(myTurn, attacks) {
			board.undoMove(move)
			continue
		}
		_, score := search(!myTurn, minScore, depth+1, opponentMoves, &move)
		board.undoMove(move)
		if score > minScore {
			continue
		}
		minScore = score
		bestMove = move
		if minScore <= bestParentScore {
			return bestMove, minScore // beta-pruning
		}
	}
	return bestMove, minScore
}

func generateMoves(myTurn bool) ([]boardMove, uint) {
	var pieceMap map[int][]*piece
	if myTurn {
		pieceMap = game.myPieces
	} else {
		pieceMap = game.otherPieces
	}
	moves := make([]boardMove, 0)
	var attacks uint
	for _, pieces := range pieceMap {
		for _, piece := range pieces {
			if piece.captured {
				continue
			}
			pieceMoves, pieceAttacks := piece.moveGenerator(piece)
			moves = append(moves, pieceMoves...)
			attacks |= pieceAttacks
		}
	}
	// TODO: come up with some meaningful sorting
	sort.Sort(sortInt(moves))
	return moves, attacks
}
