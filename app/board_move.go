package app

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type boardMove struct {
	From         int
	To           int
	castlingFrom int
	castlingTo   int
	captured     *piece
	PromotedPc   int
}

func (m boardMove) toUserMove() (UserMove, error) {
	fromCoord, err := toCoordinates(m.From)
	if err != nil {
		return UserMove{}, err
	}
	toCoord, err := toCoordinates(m.To)
	if err != nil {
		return UserMove{}, err
	}
	return UserMove{fromCoord, toCoord}, nil
}

func toCoordinates(index int) (string, error) {
	if index < 0 || index > 63 {
		errMsg := fmt.Sprintf("Invalid index: %d", index)
		return "", errors.New(errMsg)
	}
	index += 1
	rank := int(math.Ceil(float64(index) / 8.0))
	file := index % 8
	if file == 0 {
		// Remainder will be zero for last cell in a rank
		file = 8
	}
	// Convert to algebraic notation. 1 will be represented as 'a' and so on.
	fileStr := string(rune(97 + file - 1))
	rankStr := strconv.Itoa(rank)
	return fileStr + rankStr, nil
}
