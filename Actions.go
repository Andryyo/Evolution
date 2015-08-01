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
		property := a.Arguments[PARAMETER_PROPERTY].(*Property)
		card := property.ContainingCard
		result += fmt.Sprintf("Add property %#v on card %#v to creature %#v", property, card, creature)
	case ACTION_ADD_PAIR_PROPERTY:
		creatures := a.Arguments[PARAMETER_PAIR].([]*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(*Property)
		card := property.ContainingCard
		result += fmt.Sprintf("Add property %#v on card %#v to creatures %#v and %#v", property, card, creatures[0], creatures[1])
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
	case ACTION_ADD_FILTERS:
		result += fmt.Sprintf("Add filter %#v",a.Arguments[PARAMETER_FILTERS])
	case ACTION_GET_FOOD_FROM_BANK:
		result += fmt.Sprintf("Give food from bank to creature %#v",a.Arguments[PARAMETER_CREATURE])
	case ACTION_PIRACY:
		result += fmt.Sprintf("Steal %#v for %#v from %#v", a.Arguments[PARAMETER_TRAIT], a.Arguments[PARAMETER_SOURCE_CREATURE], a.Arguments[PARAMETER_TARGET_CREATURE])
	default:
		result += fmt.Sprintf("%+v", a)
	}
	return result
}

func (a *Action) Execute(game *Game) {
	switch a.Type {
	case ACTION_SEQUENCE:
		for _, action := range a.Arguments[PARAMETER_ACTIONS].([]*Action) {
			game.ExecuteAction(action)
		}
	case ACTION_SELECT:
		game.ExecuteAction(a.Arguments[PARAMETER_PLAYER].(*Player).MakeChoice(game, a.Arguments[PARAMETER_ACTIONS].([]*Action)))
	case ACTION_START_TURN:
		player := a.Arguments[PARAMETER_PLAYER].(*Player)
		actions := game.GetAlowedActions()
		action := player.MakeChoice(game, actions)
		if action != nil {
			game.ExecuteAction(action)
		}
	case ACTION_PASS:
		game.ExecuteAction(NewActionAddTrait(game.CurrentPlayer, TRAIT_PASS))
	case ACTION_NEXT_PLAYER:
		game.Players = game.Players.Next()
		game.CurrentPlayer = game.Players.Value.(*Player)
		game.ExecuteAction(NewActionStartTurn(game.CurrentPlayer))
	case ACTION_ADD_CREATURE:
		card := a.Arguments[PARAMETER_CARD].(*Card)
		player := a.Arguments[PARAMETER_PLAYER].(*Player)
		creature := &Creature{card, []*Card{}, player, []TraitType{TRAIT_REQUIRE_FOOD}}
		player.Creatures = append(player.Creatures, creature)
		player.RemoveCard(card)
	case ACTION_ADD_SINGLE_PROPERTY:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(*Property)
		card := property.ContainingCard
		card.ActiveProperty = property
		creature.Tail = append(creature.Tail, card)
		creature.Owner.RemoveCard(card)
	case ACTION_ADD_PAIR_PROPERTY:
		creatures := a.Arguments[PARAMETER_PAIR].([]*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(*Property)
		card := property.ContainingCard
		for _,creature := range creatures {
			creature.Tail = append(creature.Tail, card)
		}
		card.Owner.RemoveCard(card)
	case ACTION_ADD_TRAIT:
		trait := a.Arguments[PARAMETER_TRAIT].(TraitType)
		source := a.Arguments[PARAMETER_SOURCE].(WithTraits)
		source.AddTrait(trait)
	case ACTION_REMOVE_TRAIT:
		trait := a.Arguments[PARAMETER_TRAIT].(TraitType)
		source := a.Arguments[PARAMETER_SOURCE].(WithTraits)
		source.RemoveTrait(trait)
	case ACTION_ADD_FILTERS:
		filters := a.Arguments[PARAMETER_FILTERS].([]Filter)
		game.Filters = append(game.Filters, filters...)
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
		game.Food = foodCount
	case ACTION_GET_FOOD_FROM_BANK:
		creature := a.Arguments[PARAMETER_CREATURE]
		if game.Food == 0 {
			return
		}
		game.Food--
		game.ExecuteAction(NewActionAddTrait(creature, TRAIT_FOOD))
	case ACTION_PIRACY:
		sourceCreature := a.Arguments[PARAMETER_SOURCE_CREATURE].(*Creature)
		targetCreature := a.Arguments[PARAMETER_TARGET_CREATURE].(*Creature)
		trait := a.Arguments[PARAMETER_TRAIT].(TraitType)
		game.ExecuteAction(NewActionRemoveTrait(targetCreature, trait))
		game.ExecuteAction(NewActionAddTrait(sourceCreature, TRAIT_ADDITIONAL_FOOD))
	}
}

func NewActionStartTurn(player Source) *Action {
	return &Action{ACTION_START_TURN, map[ArgumentName]Source{PARAMETER_PLAYER: player}}
}

func NewActionNextPlayer(game *Game) *Action {
	return &Action{ACTION_NEXT_PLAYER, map[ArgumentName]Source{}}
}

func NewActionGetFoodFromBank(creature Source) *Action {
	return &Action{ACTION_GET_FOOD_FROM_BANK, map[ArgumentName]Source{PARAMETER_CREATURE: creature}}
}

func NewActionSequence(actions ...*Action) *Action {
	if len(actions) == 1 {
		return &Action{actions[0].Type, actions[0].Arguments}
	}
	return &Action{ACTION_SEQUENCE, map[ArgumentName]Source{PARAMETER_ACTIONS: actions}}
}

func NewActionSelect(actions ...*Action) *Action {
	if len(actions) == 1 {
		return &Action{actions[0].Type, actions[0].Arguments}
	}
	return &Action{ACTION_SELECT, map[ArgumentName]Source{PARAMETER_ACTIONS: actions}}
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

func NewActionAddTrait(source Source, trait Source) *Action {
	return &Action{ACTION_ADD_TRAIT, map[ArgumentName]Source{PARAMETER_SOURCE: source, PARAMETER_TRAIT: trait}}
}

func NewActionRemoveTrait(source Source, trait Source) *Action {
	return &Action{ACTION_REMOVE_TRAIT, map[ArgumentName]Source{PARAMETER_SOURCE: source, PARAMETER_TRAIT: trait}}
}

func NewActionAddFilters(filters ...Filter) *Action {
	return &Action{ACTION_ADD_FILTERS, map[ArgumentName]Source{PARAMETER_FILTERS : filters}}
}

func NewActionAddPairProperty(creatures Source, property Source) *Action {
	return &Action{ACTION_ADD_PAIR_PROPERTY, map[ArgumentName]Source{PARAMETER_PAIR: creatures, PARAMETER_PROPERTY: property}}
}

func NewActionPiracy(sourceCreature Source, targetCreature Source, trait Source) *Action {
	return &Action{ACTION_PIRACY, map[ArgumentName]Source{PARAMETER_SOURCE_CREATURE: sourceCreature, PARAMETER_TARGET_CREATURE: targetCreature, PARAMETER_TRAIT: trait}}
}

func (a *Action) InstantiateFilterPrototypeAction(game *Game, reason *Action, instantiate bool) *Action {
	instantiatedSources := make(map[ArgumentName]Source)
	for key,argument := range a.Arguments {
		instantiatedSource := InstantiateFilterSourcePrototype(game, reason, argument, instantiate)
		instantiatedSources[key] = instantiatedSource
		if !instantiate {
			continue
		}
		switch instantiatedSource.(type) {
			case OneOf:
				oneOf := instantiatedSource.(OneOf)
				actions := make([]*Action, 0, len(oneOf.Sources))
				for _, o := range oneOf.Sources {
					action := &Action{a.Type, make(map[ArgumentName]Source)}
					for k := range a.Arguments {
						action.Arguments[k] = a.Arguments[k]
					}
					action.Arguments[key] = o
					instantiatedAction := action.InstantiateFilterPrototypeAction(game, reason, instantiate)
					if instantiatedAction == nil {
						continue
					}
					if instantiatedAction.Type == ACTION_SELECT {
						actions = append(actions, instantiatedAction.Arguments[PARAMETER_ACTIONS].([]*Action)...)
					} else {
						actions = append(actions, instantiatedAction)
					}
				}
				if len(actions) == 0 {
					return nil
				} else {
					return NewActionSelect(actions...)
				}
			case AllOf:
				all := instantiatedSource.(AllOf)
				actions := make([]*Action, 0, len(all.Sources))
				for _, o := range all.Sources {
					action := &Action{a.Type, make(map[ArgumentName]Source)}
					for k := range a.Arguments {
						action.Arguments[k] = a.Arguments[k]
					}
					action.Arguments[key] = o
					actions = append(actions, action.InstantiateFilterPrototypeAction(game, reason, instantiate))
				}
				return NewActionSequence(actions...)
		}		
	}
	if !game.ActionDenied(&Action{a.Type, instantiatedSources}) {
		return &Action{a.Type, instantiatedSources}
	} else {
		return nil
	}
}

