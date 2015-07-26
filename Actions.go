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
	ACTION_ADD_TRAIT
	ACTION_REMOVE_TRAIT
	ACTION_ADD_FILTER
	ACTION_REMOVE_FILTER
	ACTION_ATTACK
)

type ArgumentName int

const (
	PARAMETER_PROPERTY ArgumentName = iota
	PARAMETER_PHASE
	PARAMETER_PLAYER
	PARAMETER_CARD
	PARAMETER_ACTIONS_SEQUENCE
	PARAMETER_CREATURE
	PARAMETER_TRAIT
	PARAMETER_SOURCE
	PARAMETER_FILTER
	PARAMETER_SOURCE_CREATURE
	PARAMETER_TARGET_CREATURE
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
	case ACTION_ADD_PROPERTY:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(Property)
		card := property.ContainingCard
		result += fmt.Sprintf("Add property %#v(%p) on card %#v(%p) to creature %#v(%p)", property, property, card, card, creature, creature)
	default:
		result += fmt.Sprintf("%+v", a)
	}
	return result
}

func (a *Action) Execute(game *Game) {
	switch a.Type {
	case ACTION_SEQUENCE:
		for _, action := range a.Arguments[PARAMETER_ACTIONS_SEQUENCE].([]*Action) {
			game.Actions.PushFront(action)
		}
	case ACTION_START_TURN:
		player := a.Arguments[PARAMETER_PLAYER].(*Player)
		actions := game.GetAlowedActions()
		action := player.MakeChoice(actions)
		if action == nil {
			game.Actions.PushFront(NewActionNextPlayer(game))
			game.Actions.PushFront(NewActionAddTrait(player, TRAIT_PASS))
			break
		}
		game.Actions.PushFront(NewActionNextPlayer(game))
		game.Actions.PushFront(action)
	case ACTION_PASS:
		game.Actions.PushFront(NewActionAddTrait(game.CurrentPlayer, TRAIT_PASS))
	case ACTION_NEXT_PLAYER:
		game.Players = game.Players.Next()
		game.CurrentPlayer = game.Players.Value.(*Player)
		game.Actions.PushFront(NewActionStartTurn(game.CurrentPlayer))
	case ACTION_ADD_CREATURE:
		card := a.Arguments[PARAMETER_CARD].(*Card)
		player := a.Arguments[PARAMETER_PLAYER].(*Player)
		creature := &Creature{card, []*Card{}, player, []TraitType{}}
		player.Creatures = append(player.Creatures, creature)
		player.RemoveCard(card)
	case ACTION_ADD_PROPERTY:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(Property)
		card := property.ContainingCard
		card.ActiveProperty = &property
		creature.Tail = append(creature.Tail, card)
		creature.Owner.RemoveCard(card)
	case ACTION_ADD_TRAIT:
		trait := a.Arguments[PARAMETER_TRAIT].(TraitType)
		source := a.Arguments[PARAMETER_SOURCE].(WithTraits)
		source.AddTrait(trait)
	case ACTION_REMOVE_TRAIT:
		trait := a.Arguments[PARAMETER_TRAIT].(TraitType)
		source := a.Arguments[PARAMETER_SOURCE].(WithTraits)
		source.RemoveTrait(trait)
	case ACTION_ADD_FILTER:
		filter := a.Arguments[PARAMETER_FILTER].(Filter)
		game.Filters = append(game.Filters, filter)
	case ACTION_REMOVE_FILTER:
		filter := a.Arguments[PARAMETER_FILTER].(Filter)
		for i,f := range game.Filters {
			if f == filter {
				game.Filters = append(game.Filters[:i],game.Filters[i+1:]...)
			} 
		}
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

func NewActionAddTrait(source Source, trait TraitType) *Action {
	return &Action{ACTION_ADD_TRAIT, map[ArgumentName]Source{PARAMETER_SOURCE: source, PARAMETER_TRAIT: trait}}
}

func NewActionRemoveTrait(source Source, trait TraitType) *Action {
	return &Action{ACTION_REMOVE_TRAIT, map[ArgumentName]Source{PARAMETER_SOURCE: source, PARAMETER_TRAIT: trait}}
}

func NewActionAddFilter(filter Filter) *Action {
	return &Action{ACTION_ADD_FILTER, map[ArgumentName]Source{PARAMETER_FILTER : filter}}
}

func NewActionRemoveFilter(filter Filter) *Action {
	return &Action{ACTION_REMOVE_FILTER, map[ArgumentName]Source{PARAMETER_FILTER : filter}}
}



func (a *Action) InstantiateFilterTemplateAction(reason *Action) *Action {
	result := &Action {a.Type, map[ArgumentName]Source {}}
	for key,argument := range a.Arguments {
		switch t := argument.(type) {
			case FilterSourceParameter:
				result.Arguments[key] = InstantiateFilterTemplateParameter(reason, t)
			case Action:
				result.Arguments[key] = t.InstantiateFilterTemplateAction(reason)
			case Filter:
				result.Arguments[key] = t.InstantiateFilterTemplate(reason)
			default:
				result.Arguments[key] = argument
		}
	}
	return result
}

