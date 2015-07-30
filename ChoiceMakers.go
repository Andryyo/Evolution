// ConsoleChoiceMaker
package main

import (
	"fmt"
)

type ChoiceMaker interface {
	MakeChoice(*Game, []*Action) *Action
}

type ConsoleChoiceMaker struct {
}

func (c ConsoleChoiceMaker) MakeChoice(game *Game, actions []*Action) *Action {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return actions[0]
	}
	fmt.Println("Creatures:")
	game.Players.Do(func (val interface{}) {
		player := val.(*Player)
		fmt.Printf("%#v:\n",player.Name)
		for i, creature := range player.Creatures {
			fmt.Printf("%v) %#v\n", i, creature)
		}
	})
	fmt.Println("Choose one action:")
	for i, action := range actions {
		fmt.Printf("%v) %#v\n", i, action)
	}
	var choiceNum int
	//_, e := fmt.Scanln(&choiceNum)
	_, e := fmt.Scanln(&choiceNum)
	if e != nil {
		fmt.Println(e)
	}
	return actions[choiceNum]
	//return actions[1]
}
