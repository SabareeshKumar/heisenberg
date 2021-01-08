package app

const (
	black = 1
	white = 2
)

const (
	king   = 1
	queen  = 2
	rook   = 3
	bishop = 4
	knight = 5
	pawn   = 6
)

var weights = map[int]int{
	king:   200,
	queen:  9,
	rook:   5,
	bishop: 3,
	knight: 3,
	pawn:   1,
}

type piece struct {
	id            int
	name          string
	color         int
	position      uint // position in powers of 2
	moveGenerator func(*piece) []boardMove
	captured      bool
	// move in which this pawn moved 2 squares so that it can be captured by
	// en-passant move by some other pawn.
	enpassantMove int
}
