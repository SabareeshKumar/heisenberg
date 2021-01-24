package app

import (
	"math"
	// "sort"
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
	moves []boardMove, lastMove *boardMove) (boardMove, float32) {
	hResult := hashedResult(myTurn, maxDepth-depth+1, lastMove)
	if hResult != nil {
		tableHits++
		return hResult.bestMove, hResult.score
	}
	if depth > maxDepth {
		evaluationsPerSearch++
		move, score := boardMove{}, eval()
		hashResult(myTurn, 0, score, move, lastMove)
		return boardMove{}, eval()
	}
	var bestMove, bestChildMove boardMove
	board := game.board
	if depth == 1 {
		moves, _ = generateMoves(myTurn, bestMove)
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
			opponentMoves, attacks := generateMoves(!myTurn, bestChildMove)
			if inCheckSimple(myTurn, attacks) {
				board.undoMove(move)
				continue
			}
			childMv, score := search(
				!myTurn, maxScore, depth+1, opponentMoves, &move)
			board.undoMove(move)
			if score < maxScore {
				continue
			}
			bestChildMove = childMv
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
		evaluationsPerSearch++
		hashResult(myTurn, maxDepth-depth+1, minScore, bestMove, lastMove)
	}()
	for _, move := range moves {
		if !isMoveValid(move) {
			continue
		}
		board.alterPosition(move)
		opponentMoves, attacks := generateMoves(!myTurn, bestChildMove)
		if inCheckSimple(myTurn, attacks) {
			board.undoMove(move)
			continue
		}
		childMv, score := search(
			!myTurn, minScore, depth+1, opponentMoves, &move)
		board.undoMove(move)
		if score > minScore {
			continue
		}
		bestChildMove = childMv
		minScore = score
		bestMove = move
		if minScore <= bestParentScore {
			return bestMove, minScore // beta-pruning
		}
	}
	return bestMove, minScore
}

func generateMoves(myTurn bool, refutationMove boardMove) ([]boardMove, uint) {
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
	// Refutation move index
	reftMoveIndex := -1
	for i, move := range moves {
		if move == refutationMove {
			reftMoveIndex = i
			break
		}
	}
	if reftMoveIndex > 0 {
		// Place refutation move first so that it triggers a cutoff
		moves = append([]boardMove{refutationMove}, moves...)
	}
	// TODO: come up with some meaningful sorting
	// sort.Sort(sortInt(moves))
	return moves, attacks
}
