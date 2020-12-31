package main

import (
	"bufio"
	"fmt"
	"github.com/SabareeshKumar/heisenberg/app"
	"log"
	"os"
)

func listen() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var move app.UserMove
	fmt.Fscanln(reader, &move.From, &move.To)
	move, err := app.MakeMove(move)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Fprintln(writer, move)
}

func main() {
	for {
		listen()
	}
}
