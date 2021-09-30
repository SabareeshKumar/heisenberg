package app

import (
	"math"
	"sort"
)

type moveSorter struct {
	moves  []boardMove
	myTurn bool
}

func (s *moveSorter) Len() int {
	return len(s.moves)
}

func (s *moveSorter) Less(i, j int) bool {
	iReft := reftTbl[s.myTurn][s.moves[i].hashKey()]
	jReft := reftTbl[s.myTurn][s.moves[j].hashKey()]
	if iReft != jReft {
		// Place refutation move first so that it triggers a cutoff
		return iReft > jReft
	}
	return s.moves[i].captured != nil
}

func (s *moveSorter) Swap(i, j int) {
	s.moves[i], s.moves[j] = s.moves[j], s.moves[i]
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
	bestMove.From = -1
	bestMove.To = -1
	var bestPath []boardMove
	board := game.board
	if depth == 1 {
		moves, _ = generateMoves(myTurn)
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
			opponentMoves, attacks := generateMoves(!myTurn)
			if inCheckSimple(myTurn, attacks) {
				board.undoMove(move)
				continue
			}
			childMv, score, spath := search(
				!myTurn, maxScore, depth+1, opponentMoves, &move, path)
			reftTbl[myTurn][childMv.hashKey()] += (maxDepth - depth + 1)
			board.undoMove(move)
			if bestMove.From != -1 && score <= maxScore {
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
		opponentMoves, attacks := generateMoves(!myTurn)
		if inCheckSimple(myTurn, attacks) {
			board.undoMove(move)
			continue
		}
		childMv, score, spath := search(
			!myTurn, minScore, depth+1, opponentMoves, &move, path)
		reftTbl[myTurn][childMv.hashKey()] += (maxDepth - depth + 1)
		board.undoMove(move)
		if bestMove.From != -1 && score >= minScore {
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
			opponentMoves, attacks := generateMoves(!myTurn)
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
		opponentMoves, attacks := generateMoves(!myTurn)
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

func generateMoves(myTurn bool) ([]boardMove, uint) {
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
	sort.Sort(&moveSorter{moves, myTurn})
	return moves, attacks
}
