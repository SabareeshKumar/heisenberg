package app

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// UserMove represents move coordinates
type UserMove struct {
	From string
	To   string
}

func (m UserMove) String() string {
	return fmt.Sprintf("%s -> %s", m.From, m.To)
}

func (m UserMove) toBoardMove() (boardMove, error) {
	fromIndex, err := toIndex(m.From)
	if err != nil {
		return boardMove{}, err
	}
	toIndex, err := toIndex(m.To)
	if err != nil {
		return boardMove{}, err
	}
	pieces := game.board.pieces
	piece := pieces[fromIndex]
	if piece.id != king || int(math.Abs(float64(toIndex-fromIndex))) == 1 {
		bm := boardMove{
			From:         fromIndex,
			To:           toIndex,
			castlingFrom: -1,
			castlingTo:   -1,
			captured:     pieces[toIndex],
		}
		return bm, nil
	}
	bm := boardMove{}
	bm.From = fromIndex
	bm.To = toIndex
	bm.captured = nil
	if m.To > m.From {
		// king side castling
		bm.castlingFrom = fromIndex + 3
		bm.castlingTo = toIndex - 1
	} else {
		// queen side castling
		bm.castlingFrom = fromIndex - 4
		bm.castlingTo = toIndex + 1
	}
	return bm, nil
}

// Given a move coordinate like 'e4', this method will find the board index
// based on this conversion:
// a1-a8 will be denoted as 0-7
// b1-b8 will be denoted as 8-15 and so on.
func toIndex(m string) (int, error) {
	if len(m) != 2 {
		errMsg := fmt.Sprintf("'%s' length is not equal to 2", m)
		return 0, errors.New(errMsg)
	}
	// Split move string to file and rank
	chars := strings.Split(m, "")
	fileStr, rankStr := chars[0], chars[1]
	// Convert rank to integer
	rank, err := strconv.Atoi(rankStr)
	if err != nil {
		return 0, err
	}
	if rank < 1 || rank > 8 {
		errMsg := fmt.Sprintf("Invalid rank: '%s'", m)
		return 0, errors.New(errMsg)
	}
	fileUpper := strings.ToLower(fileStr)
	// Convert file notation to integer. 'a' will be represented as 1 and so on
	file := int([]rune(fileUpper)[0] - 97 + 1)
	if file < 1 || file > 8 {
		errMsg := fmt.Sprintf("Invalid file: '%s'", m)
		return 0, errors.New(errMsg)
	}
	// Compute board index of given coordinate
	index := 8*(rank-1) + file - 1
	if index >= 0 && index <= 63 {
		return index, nil
	}
	errMsg := fmt.Sprintf("Invalid coordinate: '%s'", m)
	return 0, errors.New(errMsg)
}
