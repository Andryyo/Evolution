// Actions
package main

import (
	"fmt"
	"math/rand"
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
	case ACTION_ADD_SINGLE_PROPERTY:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(Property)
		card := property.ContainingCard
		result += fmt.Sprintf("Add property %#v on card %#v to creature %#v", &property, card, creature)
	case ACTION_ADD_PAIR_PROPERTY:
		creatures := a.Arguments[PARAMETER_PAIR].([]*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(Property)
		card := property.ContainingCard
		result += fmt.Sprintf("Add property %#v on card %#v to creatures %#v and %#v", &property, card, creatures[0], creatures[1])
	case ACTION_NEXT_PLAYER:
		result += "Next player"
	case ACTION_PASS:
		result += "Pass"
	case ACTION_NEW_PHASE:
		switch a.Arguments[PARAMETER_PHASE] {
			case PHASE_DEVELOPMENT:
				result += "Starting phase development"
			case PHASE_FOOD_BANK_DETERMINATION:
				result += "Starting phase food bank determination"
			case PHASE_FEEDING:
				result += "Starting phase feeding"
			case PHASE_EXTINCTION:
				result += "It's extinction time!"
		}
	case ACTION_SEQUENCE:
		result += "Unpacking action sequence"
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
	case ACTION_ADD_SINGLE_PROPERTY:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(Property)
		card := property.ContainingCard
		card.ActiveProperty = &property
		creature.Tail = append(creature.Tail, card)
		creature.Owner.RemoveCard(card)
	case ACTION_ADD_PAIR_PROPERTY:
		creatures := a.Arguments[PARAMETER_PAIR].([]*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(Property)
		card := property.ContainingCard
		players := make([]*Player, 0, 2)
		for _,creature := range creatures {
			players = append(players, creature.Owner)
			creature.Tail = append(creature.Tail, card)
		}
		card.Owners[0].RemoveCard(card)
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
	case ACTION_NEW_PHASE:
		phase := a.Arguments[PARAMETER_PHASE].(PhaseType)
		game.CurrentPhase = phase
	case ACTION_DETERMINE_FOOD_BANK:
		foodCount := 0
		switch game.PlayersCount {
			case 2: foodCount = rand.Intn(6)+1+2
			case 3: foodCount = rand.Intn(6)+rand.Intn(6)+2	
			case 4: foodCount = rand.Intn(6)+rand.Intn(6)+2+1	
		}
		game.FoodBank.Count = foodCount
		game.Actions.PushFront(NewActionNewPhase(PHASE_FEEDING))
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

func NewActionAddSingleProperty(creature Source, property Source) *Action {
	return &Action{ACTION_ADD_SINGLE_PROPERTY, map[ArgumentName]Source{PARAMETER_CREATURE: creature, PARAMETER_PROPERTY: property}}
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

func NewActionAddPairProperty(creatures Source, property Source) *Action {
	return &Action{ACTION_ADD_PAIR_PROPERTY, map[ArgumentName]Source{PARAMETER_PAIR: creatures, PARAMETER_PROPERTY: property}}
}



func (a *Action) InstantiateFilterPrototypeAction(game *Game, reason *Action) *Action {
	actionVariants := make([]map[ArgumentName]Source, 0, 1)
	actionVariants = append(actionVariants, make(map[ArgumentName]Source))
	for key,argument := range a.Arguments {
		switch t := argument.(type) {
			case FilterSourcePrototype:
				sources := InstantiateFilterSourcePrototype(game, reason, t)
				if len(sources) == 1 {
					for i := range actionVariants {
						actionVariants[i][key] = sources[0]
					}
				} else {
					tmpActionVariants := make([]map[ArgumentName]Source, 0, len(sources) * len(actionVariants))
					for _,variant := range actionVariants {
						for _,source := range sources {
							tmpVariant := make(map[ArgumentName]Source)
							for argumentName := range variant {
								tmpVariant[argumentName] = variant[argumentName]
							}
							tmpVariant[key] = source
							tmpActionVariants = append(tmpActionVariants, tmpVariant)	
						}
					}					
					actionVariants = tmpActionVariants
				}
			case Action:
				for i := range actionVariants {
					actionVariants[i][key] = t.InstantiateFilterPrototypeAction(game, reason)
				}
			case Filter:
				for i := range actionVariants {
					actionVariants[i][key] = t.InstantiateFilterPrototype(game, reason)
				}
			default:
				for i := range actionVariants {
					actionVariants[i][key] = argument
				}
		}
	}
	if len(actionVariants) == 1 {
		return &Action{a.Type, actionVariants[0]}
	} else {
		actions := make([]*Action, 0, len(actionVariants))
		for _,variant := range actionVariants {
			actions = append(actions, &Action{a.Type, variant})
		}
		return NewActionSequence(actions...)
	}
}

