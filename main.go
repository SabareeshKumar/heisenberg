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
	err := app.MakeMove(move)
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
