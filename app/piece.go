package app

type piece struct {
	id            int
	name          string
	color         int
	position      uint // position in powers of 2
	moveGenerator func(*piece) []boardMove
	captured      bool
	moveCount     int // Number of times piece has moved
	enpassantMove int // First move of pawn
	promotedBy    *piece
}

type pieceMeta struct {
	name          string
	moveGenerator func(*piece) []boardMove
}

var weights = map[int]int{
	king:   200,
	queen:  9,
	rook:   5,
	bishop: 3,
	knight: 3,
	pawn:   1,
}

var blackMeta = map[int]pieceMeta{
	king:   pieceMeta{"Black King", kingMoves},
	queen:  pieceMeta{"Black Queen", queenMoves},
	rook:   pieceMeta{"Black Rook", rookMoves},
	bishop: pieceMeta{"Black Bishop", bishopMoves},
	knight: pieceMeta{"Black Knight", knightMoves},
	pawn:   pieceMeta{"Black Pawn", blackPawnMoves},
}

var whiteMeta = map[int]pieceMeta{
	king:   pieceMeta{"White King", kingMoves},
	queen:  pieceMeta{"White Queen", queenMoves},
	rook:   pieceMeta{"White Rook", rookMoves},
	bishop: pieceMeta{"White Bishop", bishopMoves},
	knight: pieceMeta{"White Knight", knightMoves},
	pawn:   pieceMeta{"White Pawn", whitePawnMoves},
}

func newBlackPiece(pieceType int, position uint) *piece {
	return &piece{
		id:            pieceType,
		name:          blackMeta[pieceType].name,
		color:         black,
		position:      position,
		moveGenerator: blackMeta[pieceType].moveGenerator,
		captured:      false,
		moveCount:     0,
		enpassantMove: -1,
	}
}

func newWhitePiece(pieceType int, position uint) *piece {
	return &piece{
		id:            pieceType,
		name:          whiteMeta[pieceType].name,
		color:         white,
		position:      position,
		moveGenerator: whiteMeta[pieceType].moveGenerator,
		captured:      false,
		moveCount:     0,
		enpassantMove: -1,
		promotedBy:    nil,
	}
}
