// ConsoleChoiceMaker
package main

import (
	"fmt"
)

type ChoiceMaker interface {
	Notify(string)
	GetChoice() int
	MakeChoice([]*Action) *Action
	GetName() string
}

type ConsoleChoiceMaker struct {
	name string
}

func (c *ConsoleChoiceMaker) Notify(s string) {
	fmt.Printf(s)
}

func (c *ConsoleChoiceMaker) GetChoice() int {
	var i int
	fmt.Scanln(&i)
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
	c.Notify("Choose one action:")
	for i, action := range actions {
		c.Notify(fmt.Sprintf("%v) %#v", i, action))
	}
	return actions[c.GetChoice()]
}
