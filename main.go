package main

import (
	"fmt"
	"github.com/SabareeshKumar/heisenberg/app"
	"strings"
)

var myTurn = false

func toggleTurn() bool {
	myTurn = !myTurn
	status := app.GameStatus(myTurn)
	if status == app.InProgress {
		return true
	}
	if status == app.Win {
		fmt.Println("You won !!")
		return false
	}
	if status == app.Lost {
		fmt.Println("You lost :(")
		return false
	}
	fmt.Println("Oops. It's a stalemate")
	return false

}

func play() bool {
	app.PrintBoard()
	if myTurn {
		fmt.Print("Thinking...")
		move, err := app.MyMove()
		if err != nil {
			fmt.Println(err)
			return true
		}
		fmt.Println(move)
		return toggleTurn()
	}
	fmt.Print("\nYour move (Enter 'q' to quit game): ")
	var input string
	fmt.Scan(&input)
	if strings.ToLower(strings.Trim(input, " ")) == "q" {
		return false
	}
	move := app.UserMove{input, ""}
	fmt.Scan(&move.To)
	mv, err := move.ToBoardMove()
	if err != nil {
		fmt.Println(err)
		return true
	}
	if app.IsPromotion(mv) {
		for {
			fmt.Println("Please choose a piece type to promote to:")
			fmt.Println("1. Queen")
			fmt.Println("2. Rook")
			fmt.Println("3. Bishop")
			fmt.Println("4. Knight")
			fmt.Scan(&mv.PromotedPc)
			found := true
			switch mv.PromotedPc {
			case 1, 2, 3, 4:
				break
			default:
				found = false
			}
			if found {
				// Convert to corresponding piece ID's
				mv.PromotedPc += 1
				break
			}
		}
	}
	err = app.MakeMove(mv)
	if err != nil {
		fmt.Println(err)
		return true
	}
	return toggleTurn()
}

func main() {
	var colorChoice int
	for {
		fmt.Println("\nChoose a color:\n1. Black\n2. White")
		fmt.Scanln(&colorChoice)
		if colorChoice != 1 && colorChoice != 2 {
			fmt.Println("Invalid choice")
			continue
		}
		myTurn = (colorChoice == 1)
		app.InitGame(colorChoice)
		for play() {
		}
	}
}
