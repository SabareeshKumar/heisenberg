package app

import (
	"math"
)

type sortInt []boardMove

func (bm sortInt) Len() int {
	return len(bm)
}

func (bm sortInt) Less(i, j int) bool {
	score1, score2 := float32(-1), float32(-1)
	if isMoveLegal(bm[i]) {
		game.board.alterPosition(bm[i])
		score1 = eval()
		game.board.undoMove(bm[i])
	}
	if isMoveLegal(bm[j]) {
		game.board.alterPosition(bm[j])
		score1 = eval()
		game.board.undoMove(bm[j])
	}
	return score1 > score2
}

func (bm sortInt) Swap(i, j int) {
	bm[i], bm[j] = bm[j], bm[i]
}

func search(
	myTurn bool, bestParentScore float32, depth int,
	moves []boardMove, lastMove *boardMove,
	path []boardMove) (boardMove, float32, []boardMove) {
	if depth > maxDepth {
		score := quiescenceSearch(myTurn, bestParentScore, 1, moves)
		return boardMove{}, score, path
	}
	hResult := hashedResult(myTurn, maxDepth-depth+1, lastMove)
	if hResult != nil {
		tableHits++
		return hResult.bestMove, hResult.score, path
	}
	var bestMove boardMove
	var bestPath []boardMove
	reftMoves := make(map[int]bool, 0)
	board := game.board
	if depth == 1 {
		moves, _ = generateMoves(myTurn, reftMoves)
	}
	if myTurn {
		maxScore := float32(math.MinInt32)
		defer func() {
			evaluationsPerSearch++
			hashResult(myTurn, maxDepth-depth+1, maxScore, bestMove, lastMove)
		}()
		for _, move := range moves {
			if !isMoveValid(move) {
				continue
			}
			board.alterPosition(move)
			opponentMoves, attacks := generateMoves(!myTurn, reftMoves)
			if inCheckSimple(myTurn, attacks) {
				board.undoMove(move)
				continue
			}
			childMv, score, spath := search(
				!myTurn, maxScore, depth+1, opponentMoves, &move, path)
			reftMoves[childMv.hashKey()] = true
			board.undoMove(move)
			if score < maxScore {
				continue
			}
			bestPath = spath
			maxScore = score
			bestMove = move
			if maxScore >= bestParentScore {
				path = append(path, append(bestPath, bestMove)...)
				return bestMove, maxScore, path // alpha-pruning
			}
		}
		path = append(path, append(bestPath, bestMove)...)
		return bestMove, maxScore, path
	}
	minScore := float32(math.MaxInt32)
	defer func() {
		evaluationsPerSearch++
		hashResult(myTurn, maxDepth-depth+1, minScore, bestMove, lastMove)
	}()
	for _, move := range moves {
		if !isMoveValid(move) {
			continue
		}
		board.alterPosition(move)
		opponentMoves, attacks := generateMoves(!myTurn, reftMoves)
		if inCheckSimple(myTurn, attacks) {
			board.undoMove(move)
			continue
		}
		childMv, score, spath := search(
			!myTurn, minScore, depth+1, opponentMoves, &move, path)
		reftMoves[childMv.hashKey()] = true
		board.undoMove(move)
		if score > minScore {
			continue
		}
		bestPath = spath
		minScore = score
		bestMove = move
		if minScore <= bestParentScore {
			path = append(path, append(bestPath, bestMove)...)
			return bestMove, minScore, path // beta-pruning
		}
	}
	path = append(path, append(bestPath, bestMove)...)
	return bestMove, minScore, path
}

func quiescenceSearch(
	myTurn bool, bestParentScore float32, depth int,
	moves []boardMove) float32 {
	defer func() {
		evaluationsPerSearch++
	}()
	if depth > quiescenceDepth {
		return eval()
	}
	reftMoves := make(map[int]bool, 0)
	board := game.board
	if myTurn {
		maxScore := float32(math.MinInt32)
		foundMove := false
		for _, move := range moves {
			if move.captured == nil && move.PromotedPc <= 0 {
				// Discard quiet moves
				continue
			}
			if !isMoveValid(move) {
				continue
			}
			board.alterPosition(move)
			opponentMoves, attacks := generateMoves(!myTurn, reftMoves)
			if inCheckSimple(myTurn, attacks) {
				board.undoMove(move)
				continue
			}
			foundMove = true
			score := quiescenceSearch(!myTurn, maxScore, depth+1, opponentMoves)
			board.undoMove(move)
			if score < maxScore {
				continue
			}
			maxScore = score
			if maxScore >= bestParentScore {
				return maxScore // alpha-pruning
			}
		}
		if !foundMove { // All moves are quiet. Just evaluate the board.
			return eval()
		}
		return maxScore
	}
	minScore := float32(math.MaxInt32)
	foundMove := false
	for _, move := range moves {
		if move.captured == nil && move.PromotedPc <= 0 {
			// Discard quiet moves
			continue
		}
		if !isMoveValid(move) {
			continue
		}
		board.alterPosition(move)
		opponentMoves, attacks := generateMoves(!myTurn, reftMoves)
		if inCheckSimple(myTurn, attacks) {
			board.undoMove(move)
			continue
		}
		foundMove = true
		score := quiescenceSearch(!myTurn, minScore, depth+1, opponentMoves)
		board.undoMove(move)
		if score > minScore {
			continue
		}
		minScore = score
		if minScore <= bestParentScore {
			return minScore // beta-pruning
		}
	}
	if !foundMove { // All moves are quiet. Just evaluate the board.
		return eval()
	}
	return minScore
}

func generateMoves(
	myTurn bool, reftMoves map[int]bool) ([]boardMove, uint) {
	var pieceMap map[int][]*piece
	if myTurn {
		pieceMap = game.myPieces
	} else {
		pieceMap = game.otherPieces
	}
	moves := make([]boardMove, 0)
	var attacks uint
	for _, pieceId := range []int{queen, rook, knight, bishop, king, pawn} {
		for _, piece := range pieceMap[pieceId] {
			if piece.captured {
				continue
			}
			pieceMoves, pieceAttacks := piece.moveGenerator(piece)
			moves = append(moves, pieceMoves...)
			attacks |= pieceAttacks
		}
	}
	foundReftMoves := make([]boardMove, 0)
	otherMoves := make([]boardMove, 0)
	for _, move := range moves {
		if reftMoves[move.hashKey()] {
			foundReftMoves = append(foundReftMoves, move)
		} else {
			otherMoves = append(otherMoves, move)
		}
	}
	// TODO: come up with some meaningful sorting
	// Place refutation move first so that it triggers a cutoff
	return append(foundReftMoves, otherMoves...), attacks
}
