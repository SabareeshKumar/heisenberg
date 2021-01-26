package app

type tpnMeta struct {
	depth    int
	score    float32
	bestMove boardMove
}

func getKey(myTurn bool, lastMove *boardMove) int64 {
	board := game.board
	key := board.hash ^ turnHash[myTurn]
	myKing := game.myPieces[king][0]
	pos := myKing.position
	if canCastle(pos, pos-2) {
		key ^= myCastlingHash[0]
	}
	if canCastle(pos, pos+2) {
		key ^= myCastlingHash[1]
	}
	otherKing := game.otherPieces[king][0]
	pos = otherKing.position
	if canCastle(pos, pos-2) {
		key ^= otherCastlingHash[0]
	}
	if canCastle(pos, pos+2) {
		key ^= otherCastlingHash[1]
	}
	if lastMove == nil {
		return key
	}
	pc := board.pieces[lastMove.To]
	if pc.enpassantMove != game.moveCount {
		return key
	}
	_, file := getRankFile(lastMove.To)
	return key ^ fileHash[file]
}

func hashedResult(myTurn bool, depth int, lastMove *boardMove) *tpnMeta {
	meta, ok := game.tpnTbl[getKey(myTurn, lastMove)]
	if !ok || meta.depth < depth {
		return nil
	}
	return &meta
}

func hashResult(
	myTurn bool, depth int, score float32, bestMove boardMove,
	lastMove *boardMove) {
	game.tpnTbl[getKey(myTurn, lastMove)] = tpnMeta{
		depth:    depth,
		score:    score,
		bestMove: bestMove,
	}
}
