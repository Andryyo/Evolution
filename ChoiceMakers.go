// ConsoleChoiceMaker
package main

import (
	"fmt"
)

type ChoiceMaker interface {
	MakeChoice([]*Action) *Action
}

type ConsoleChoiceMaker struct {
}

func (c ConsoleChoiceMaker) MakeChoice(actions []*Action) *Action {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return actions[0]
	}
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
