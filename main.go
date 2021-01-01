package main

import (
	"fmt"
	"github.com/SabareeshKumar/heisenberg/app"
	"strings"
)

func play() bool {
	fmt.Print("\nYour move (Enter 'q' to quit game): ")
	var input string
	fmt.Scan(&input)
	if strings.ToLower(strings.Trim(input, " ")) == "q" {
		return false
	}
	move := app.UserMove{input, ""}
	fmt.Scan(&move.To)
	fmt.Print("Thinking...")
	myMove, err := app.MakeMove(move)
	if err != nil {
		fmt.Println(err)
		return true
	}
	fmt.Println(myMove)
	return true
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
		app.InitGame(colorChoice)
		for play() {
		}
	}
}
