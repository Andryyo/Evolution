// Actions
package main

import (
	"fmt"
)

type ActionType int

const (
	ACTION_SEQUENCE ActionType = iota
	ACTION_START_TURN
	ACTION_NEXT_PLAYER
	ACTION_ADD_CREATURE
	ACTION_ADD_PROPERTY
	ACTION_PASS
	ACTION_NEW_PHASE
)

type ArgumentName int

const (
	PARAMETER_PROPERTY ArgumentName = iota
	PARAMETER_PHASE
	PARAMETER_PLAYER
	PARAMETER_CARD
	PARAMETER_ACTIONS_SEQUENCE
	PARAMETER_CREATURE
)

type Action struct {
	Type      ActionType
	Arguments map[ArgumentName]Source
}

func (a *Action) GoString() string {
	result := ""
	switch a.Type {
	case ACTION_ADD_CREATURE:
		card := a.Arguments[PARAMETER_CARD].(*Card)
		result += fmt.Sprintf("Create creature using (%#v) card", card)
	case ACTION_START_TURN:
		result += fmt.Sprintf("Player %s starts turn", a.Arguments[PARAMETER_PLAYER].(*Player).Name)
	default:
		result += fmt.Sprintf("%+v", a)
	}
	return result
}

func (a *Action) Execute(game *Game) {
	switch a.Type {
	case ACTION_SEQUENCE:
		for _, action := range a.Arguments[PARAMETER_ACTIONS_SEQUENCE].([]*Action) {
			game.ExecuteAction(action)
		}
	case ACTION_START_TURN:
		player := a.Arguments[PARAMETER_PLAYER].(*Player)
		game.CurrentPlayer = player
		actions := game.GetAlowedActions()
		fmt.Println("Choose one action:")
		for i, action := range actions {
			fmt.Printf("%v) %#v\n", i, action)
		}
	case ACTION_NEXT_PLAYER:
		game.ExecuteAction(NewActionStartTurn(game.Players.Next().Value.(*Player)))
	case ACTION_ADD_CREATURE:
		card := a.Arguments[PARAMETER_CARD].(*Card)
		player := a.Arguments[PARAMETER_PLAYER].(*Player)
		creature := &Creature{card, []*Card{}, player}
		player.Creatures = append(player.Creatures, creature)
	case ACTION_ADD_PROPERTY:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(*Property)
		card := property.ContainingCard
		card.ActiveProperty = property
		creature.Tail = append(creature.Tail, card)
	}
}

func NewActionStartTurn(player Source) *Action {
	return &Action{ACTION_START_TURN, map[ArgumentName]Source{PARAMETER_PLAYER: player}}
}

func NewActionNextPlayer(game *Game) *Action {
	return &Action{ACTION_NEXT_PLAYER, map[ArgumentName]Source{}}
}

func NewActionSequence(actions ...*Action) *Action {
	return &Action{ACTION_SEQUENCE, map[ArgumentName]Source{PARAMETER_ACTIONS_SEQUENCE: actions}}
}

func NewActionNewPhase(phaseType PhaseType) *Action {
	return &Action{ACTION_NEW_PHASE, map[ArgumentName]Source{PARAMETER_PHASE: phaseType}}
}

func NewActionAddCreature(player Source, card Source) *Action {
	return &Action{ACTION_ADD_CREATURE, map[ArgumentName]Source{PARAMETER_PLAYER: player, PARAMETER_CARD: card}}
}

func NewActionAddProperty(creature Source, property Source) *Action {
	return &Action{ACTION_ADD_PROPERTY, map[ArgumentName]Source{PARAMETER_CREATURE: creature, PARAMETER_PROPERTY: property}}
}
