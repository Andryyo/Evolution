// Actions
package main

import (
	"fmt"
	"math/rand"
	"strconv"
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
	case ACTION_END_TURN:
		result += "End turn"
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
			case PHASE_FINAL:
				result += "It's finale!"
		}
	case ACTION_SEQUENCE:
		result += "Unpacking action sequence"
	case ACTION_ADD_FILTERS:
		result += fmt.Sprintf("Add filter %#v",a.Arguments[PARAMETER_FILTERS])
	case ACTION_GET_FOOD_FROM_BANK:
		result += fmt.Sprintf("Give food from bank to creature %#v",a.Arguments[PARAMETER_CREATURE])
	case ACTION_PIRACY:
		result += fmt.Sprintf("Steal %#v for %#v from %#v", a.Arguments[PARAMETER_TRAIT], a.Arguments[PARAMETER_SOURCE_CREATURE], a.Arguments[PARAMETER_TARGET_CREATURE])
	case ACTION_DESTROY_BANK_FOOD:
		result += fmt.Sprintf("Destroy one food in bank")
	case ACTION_ATTACK:
		result += fmt.Sprintf("Attack %#v with %#v", a.Arguments[PARAMETER_TARGET_CREATURE], a.Arguments[PARAMETER_SOURCE_CREATURE])
	case ACTION_SELECT_FROM_AVAILABLE_ACTIONS:
		result += fmt.Sprint("Player selecting action")
	case ACTION_GAIN_FOOD:
		result += fmt.Sprintf("Creature %#v gain food", a.Arguments[PARAMETER_CREATURE])
	case ACTION_GAIN_ADDITIONAL_FOOD:
		result += fmt.Sprintf("Creature %#v gain additional food", a.Arguments[PARAMETER_CREATURE])
	case ACTION_EAT:
		result += fmt.Sprint("Creature %#v was eaten", a.Arguments[PARAMETER_CREATURE])
	case ACTION_BURN_FAT:
		result += fmt.Sprintf("Burn fat on creature %#v", a.Arguments[PARAMETER_CREATURE])
	case ACTION_HIBERNATE:
		result += fmt.Sprintf("Hibernate creature %#v", a.Arguments[PARAMETER_CREATURE])
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
		actions := a.Arguments[PARAMETER_ACTIONS].([]*Action)
		chooser := actions[0].Arguments[PARAMETER_PLAYER].(*Player)
		game.ExecuteAction(chooser.MakeChoice(actions))
	case ACTION_START_TURN:
		break
	case ACTION_SELECT_FROM_AVAILABLE_ACTIONS:
		game.NotifyAll("Cards in desk: " + strconv.Itoa(len(game.Deck)))
		game.NotifyAll("Food in bank: " + strconv.Itoa(game.Food))
		game.NotifyAll("Creatures:")
		game.Players.Do(func (val interface{}) {
			player := val.(*Player)
			game.NotifyAll(fmt.Sprintf("%#v(%#v):",player.Name, player.Traits))
			for i, creature := range player.Creatures {
				game.NotifyAll(fmt.Sprintf("%v) %#v", i, creature))
			}
		})
		player := game.CurrentPlayer
		actions := game.GetAlowedActions()
		action := player.MakeChoice(actions)
		if action != nil {
			game.ExecuteAction(action)
		}
	case ACTION_PASS:
		game.ExecuteAction(NewActionAddTrait(game.CurrentPlayer, TRAIT_PASS))
	case ACTION_END_TURN:
		game.ExecuteAction(
			NewActionAddFilters(&FilterAction{
				FILTER_ACTION_REPLACE,
				&ConditionActionType{ACTION_SELECT_FROM_AVAILABLE_ACTIONS},
				&ConditionActionType{ACTION_NEXT_PLAYER},
				NewActionNextPlayer(game)}))
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
		creature.Owner.RemoveCard(card)
		card.ActiveProperty = property
		creature.Tail = append(creature.Tail, card)
		card.Owners = []Source{creature}
	case ACTION_ADD_PAIR_PROPERTY:
		creatures := a.Arguments[PARAMETER_PAIR].([]*Creature)
		property := a.Arguments[PARAMETER_PROPERTY].(*Property)
		card := property.ContainingCard
		card.Owners[0].(*Player).RemoveCard(card)
		card.Owners = make([]Source, 0, 2)
		for _,creature := range creatures {
			creature.Tail = append(creature.Tail, card)
			card.Owners = append(card.Owners, creature)
		}
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
	case ACTION_BURN_FAT:
		creature := a.Arguments[PARAMETER_CREATURE]
		game.ExecuteAction(NewActionRemoveTrait(creature, TRAIT_FAT))
		game.ExecuteAction(NewActionAddTrait(creature, TRAIT_FOOD))
	case ACTION_PIRACY:
		sourceCreature := a.Arguments[PARAMETER_SOURCE_CREATURE].(*Creature)
		targetCreature := a.Arguments[PARAMETER_TARGET_CREATURE].(*Creature)
		trait := a.Arguments[PARAMETER_TRAIT].(TraitType)
		game.ExecuteAction(NewActionRemoveTrait(targetCreature, trait))
		game.ExecuteAction(NewActionAddTrait(sourceCreature, TRAIT_ADDITIONAL_FOOD))
	case ACTION_DESTROY_BANK_FOOD:
		if game.Food == 0 {
			return
		}
		game.Food--
	case ACTION_REMOVE_CREATURE:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		player := creature.Owner
		for _,card := range creature.Tail {
			game.ExecuteAction(NewActionRemoveCard(card))
		}
		player.RemoveCreature(creature)
	case ACTION_REMOVE_CARD:
		card := a.Arguments[PARAMETER_CARD].(*Card)
		game.ExecuteAction(NewActionRemoveProperty(card.ActiveProperty))
		for _, owner := range card.Owners {
			owner.(*Creature).RemoveCard(card)
		}
	case ACTION_ATTACK:
		sourceCreature := a.Arguments[PARAMETER_SOURCE_CREATURE].(*Creature)
		switch target := a.Arguments[PARAMETER_TARGET_CREATURE].(type) {
			case *Creature:
				game.ExecuteAction(NewActionEat(target))
				game.ExecuteAction(NewActionGainAdditionalFood(sourceCreature, 2))
			case *Property:
				game.ExecuteAction(NewActionRemoveCard(target.ContainingCard))
				game.ExecuteAction(NewActionGainAdditionalFood(sourceCreature, 1))
		}
	case ACTION_EXTINCT:
		game.Players.Do(func (val interface{}) {
			player := val.(*Player)
			removed := true
			for removed {
				removed = false
				for _,creature := range player.Creatures {
					if !creature.ContainsTrait(TRAIT_FED) {
						removed = true
						game.ExecuteAction(NewActionRemoveCreature(creature))
						break
					}
				}
			}
			for _,creature := range player.Creatures {
				creature.RemoveTrait(TRAIT_FED)
				for creature.ContainsTrait(TRAIT_FOOD) {
					creature.RemoveTrait(TRAIT_FOOD)
				}
				for creature.ContainsTrait(TRAIT_ADDITIONAL_FOOD) {
					creature.RemoveTrait(TRAIT_ADDITIONAL_FOOD)
				}
			}
		})
	case ACTION_TAKE_CARDS:
		cardsCounts := make(map[*Player]int)
		game.Players.Do(func (val interface{}) {
			player := val.(*Player)
			cardsCounts[player] = len(player.Creatures) + 1
		})
		for len(cardsCounts) != 0 {
			for player,cardsCount := range cardsCounts {
				if cardsCount == 0 {
					delete(cardsCounts, player)
				} else {
					cardsCounts[player] = cardsCount - 1
					game.TakeCard(player)
				}
			}
		}
	case ACTION_GAIN_FOOD:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		game.ExecuteAction(NewActionAddTrait(creature, TRAIT_FOOD))
	case ACTION_GAIN_ADDITIONAL_FOOD:
		creature := a.Arguments[PARAMETER_CREATURE].(*Creature)
		count := a.Arguments[PARAMETER_COUNT].(int)
		for i:=0; i<count;i++ {
			game.ExecuteAction(NewActionAddTrait(creature, TRAIT_ADDITIONAL_FOOD))
		}
	case ACTION_EAT:
		creature := a.Arguments[PARAMETER_CREATURE]
		game.ExecuteAction(NewActionRemoveCreature(creature))
	case ACTION_HIBERNATE:
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

func NewActionBurnFat(creature Source) *Action {
	return &Action{ACTION_BURN_FAT, map[ArgumentName]Source{PARAMETER_CREATURE: creature}}
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

func NewActionGrazing(source Source) *Action {
	return &Action{ACTION_DESTROY_BANK_FOOD, map[ArgumentName]Source{PARAMETER_PROPERTY: source}}
}

func NewActionSelectFromAvailableActions() *Action {
	return &Action{ACTION_SELECT_FROM_AVAILABLE_ACTIONS, map[ArgumentName]Source{}}
}

func NewActionAttack(player Source, sourceCreature Source, targetCreature Source) *Action {
	return &Action{ACTION_ATTACK, map[ArgumentName]Source{PARAMETER_PLAYER: player, PARAMETER_SOURCE_CREATURE: sourceCreature, PARAMETER_TARGET_CREATURE: targetCreature}}
}

func NewActionRemoveCreature(creature Source) *Action {
	return &Action{ACTION_REMOVE_CREATURE, map[ArgumentName]Source{PARAMETER_CREATURE: creature}}
}

func NewActionRemoveCard(card Source) *Action {
	return &Action{ACTION_REMOVE_CARD, map[ArgumentName]Source{PARAMETER_CARD: card}}
}

func NewActionRemoveProperty(property Source) * Action {
	return &Action{ACTION_REMOVE_PROPERTY, map[ArgumentName]Source{PARAMETER_PROPERTY: property}}
}

func NewActionEndTurn() *Action {
	return &Action{ACTION_END_TURN, map[ArgumentName]Source{}}
}

func NewActionGainFood(creature Source) *Action {
	return &Action{ACTION_GAIN_FOOD, map[ArgumentName]Source{PARAMETER_CREATURE : creature}}
}

func NewActionGainAdditionalFood(creature Source, count int) *Action {
	return &Action{ACTION_GAIN_ADDITIONAL_FOOD, map[ArgumentName]Source{PARAMETER_CREATURE : creature, PARAMETER_COUNT: count}}
}

func NewActionEat(creature Source) *Action {
	return &Action{ACTION_EAT, map[ArgumentName]Source{PARAMETER_CREATURE : creature}}
}

func NewActionTakeCards() *Action {
	return &Action{ACTION_TAKE_CARDS, map[ArgumentName]Source{}}
}

func NewActionHibernate(creature Source) *Action {
	return &Action{ACTION_HIBERNATE, map[ArgumentName]Source{PARAMETER_CREATURE: creature}}
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
	return &Action{a.Type, instantiatedSources}
}

