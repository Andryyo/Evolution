// Evolution project main.go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Let's start game!")
	game := NewGame("Andrii", "Opponent")
	game.ShuffleDeck()
}
