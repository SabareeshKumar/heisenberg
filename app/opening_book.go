package app

import (
	"math"
	"sort"
)

type openingLine struct {
	moves  []UserMove
	weight int // Increase weight to prioritize a opening line
}

type openingBook []openingLine

func (bk openingBook) Len() int {
	return len(bk)
}

func (bk openingBook) Less(i, j int) bool {
	return bk[i].weight < bk[j].weight
}

func (bk openingBook) Swap(i, j int) {
	bk[i], bk[j] = bk[j], bk[i]
}

var book = []openingLine{
	openingLine{ // Sicilian Defense
		moves: []UserMove{
			UserMove{"e2", "e4"},
			UserMove{"c7", "c5"},
		},
		weight: 0,
	},
	openingLine{ // French Defense
		moves: []UserMove{
			UserMove{"e2", "e4"},
			UserMove{"e7", "e6"},
		},
		weight: 0,
	},
	openingLine{ // French Defense
		moves: []UserMove{
			UserMove{"e2", "e4"},
			UserMove{"e7", "e6"},
		},
		weight: 0,
	},
	openingLine{ // Ruy Lopez Opening
		moves: []UserMove{
			UserMove{"e2", "e4"},
			UserMove{"e7", "e5"},
			UserMove{"g1", "f3"},
			UserMove{"b8", "c6"},
			UserMove{"f1", "b5"},
		},
		weight: 0,
	},
}

func hashOpenings(
	myTurn bool, line openingLine, index int, lastMove *boardMove) error {
	if index == len(line.moves) {
		return nil
	}
	uMove := line.moves[index]
	bMove, err := uMove.ToBoardMove()
	if err != nil {
		return err
	}
	hashResult(myTurn, maxDepth, math.MaxInt32, bMove, lastMove)
	board := game.board
	if err := board.alterPosition(bMove); err != nil {
		return err
	}
	if err := hashOpenings(!myTurn, line, index+1, &bMove); err != nil {
		return err
	}
	board.undoMove(bMove)
	return nil
}

func loadOpeningBook(userColorChoice int) error {
	if len(book) == 0 {
		return nil
	}
	myTurn := false
	if userColorChoice == black {
		myTurn = true
	}
	sort.Sort(openingBook(book))
	for _, line := range book {
		if err := hashOpenings(myTurn, line, 0, nil); err != nil {
			return err
		}
	}
	return nil
}
