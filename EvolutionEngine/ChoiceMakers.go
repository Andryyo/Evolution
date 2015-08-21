// ConsoleChoiceMaker
package EvolutionEngine

import (
	"fmt"
)

type ChoiceMaker interface {
	Notify(game *Game, action *Action)
	MakeChoice([]*Action) *Action
	GetName() string
	SetPlayer(player *Player)
}

type ConsoleChoiceMaker struct {
	name string
}

func (c *ConsoleChoiceMaker) Notify(action *Action) {
	fmt.Println(action.GoString())
}

func (c *ConsoleChoiceMaker) GetChoice() int {
	var i int
	_, err := fmt.Scanln(&i)
	if err!=nil {
		fmt.Println("Error!")
	}
	return i
}

func (c *ConsoleChoiceMaker) GetName() string {
	return c.name
}

func (c *ConsoleChoiceMaker) MakeChoice(actions []*Action) *Action {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return actions[0]
	}
	//c.Notify("Choose one action:")
	for _, action := range actions {
		//c.Notify(fmt.Sprintf("%v) %#v", i, action))
		c.Notify(action)
	}
	return actions[c.GetChoice()]
}

func (c *ConsoleChoiceMaker) SetOwner(player *Player) {
}