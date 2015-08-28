// Filters
package EvolutionEngine

import "fmt"

type Filter interface {
	GetType() FilterType
	GetCondition() Condition
	CheckCondition(game *Game, action *Action) bool
	CheckRemoveCondition(game *Game, action *Action) bool
	InstantiateFilterPrototype(game *Game, reason *Action, instantiate bool) Filter
}

type FilterDeny struct {
	condition Condition
	removeCondition Condition
}

type InstantiationOff struct {
	source Source
}

type InstantiationOn struct {
	source Source
}

type TraitsCount struct {
	source Source
	traits Source
}

type Accessor struct {
	source Source
	mode AccessorMode
}

func (t TraitsCount) GoString() string {
	return fmt.Sprintf("%#v count of %#v", t.traits, t.source)
}

type OneOf struct {
	Sources []Source
}

type AllOf struct {
	Sources []Source
}

type TypeOf struct {
	source Source
}

func (f *FilterDeny) GetType() FilterType {
	return FILTER_DENY
}

func (f *FilterDeny) GetCondition() Condition {
	return f.condition
}

func (f *FilterDeny) CheckCondition(game *Game, action *Action) (result bool) {
	if action.Type == ACTION_ADD_TRAIT {
		//fmt.Printf("Deny %#v because %#v:%#v\n", action, f.condition.InstantiateFilterPrototypeCondition(game, action, true), f.condition.CheckCondition(game, action))
	}
	return f.condition == nil || f.condition.CheckCondition(game, action)
}

func (f *FilterDeny) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.CheckCondition(game, action)
}

func (f *FilterDeny) InstantiateFilterPrototype(game *Game, reason *Action, instantiate bool) Filter {
	var condition Condition
	var removeCondition Condition
	if f.condition != nil {
		condition = f.condition.InstantiateFilterPrototypeCondition(game, reason, instantiate)
	}
	if f.removeCondition != nil {
		removeCondition = f.removeCondition.InstantiateFilterPrototypeCondition(game, reason, instantiate)
	}
	return &FilterDeny{condition, removeCondition}
}

func (f FilterDeny) GoString() string {
	return fmt.Sprintf("Deny %#v", f.condition)
}

type FilterAllow struct {
	condition Condition
	removeCondition Condition
	actions []*Action
}

func (f *FilterAllow) GetType() FilterType {
	return FILTER_ALLOW
}

func (f *FilterAllow) GetCondition() Condition {
	return f.condition
}

func (f *FilterAllow) GetActions(game *Game) []*Action {
	var instantiatedActions []*Action
	for _,action := range f.actions {
		source := action.InstantiateFilterPrototypeAction(game, nil, true)
		if source != nil {
			instantiatedActions = append(instantiatedActions, source)
		}
	}
	return instantiatedActions
}

func (f *FilterAllow) CheckCondition(game *Game, action *Action) bool {
	return f.condition == nil || f.condition.CheckCondition(game, action)
}

func (f *FilterAllow) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.CheckCondition(game, action)
}

func (f *FilterAllow) InstantiateFilterPrototype(game *Game, reason *Action, instantiate bool) Filter {
	var condition Condition
	var removeCondition Condition
	if f.condition != nil {
		condition = f.condition.InstantiateFilterPrototypeCondition(game, reason, instantiate)
	}
	if f.removeCondition != nil {
		removeCondition = f.removeCondition.InstantiateFilterPrototypeCondition(game, reason, instantiate)
	}
	var instantiatedActions []*Action
	for _,action := range f.actions {
		source := action.InstantiateFilterPrototypeAction(game, reason, instantiate)
		instantiatedActions = append(instantiatedActions, source)
	}
	return &FilterAllow{condition, removeCondition, instantiatedActions}
}

func NewFilterAllow(condition Condition, removeCondition Condition, actions ...*Action) Filter {
	return &FilterAllow{condition, removeCondition, actions}
}

type FilterAction struct {
	actionType FilterType
	condition Condition
	removeCondition Condition
	action    *Action
}

func (f *FilterAction) GetType() FilterType {
	return f.actionType
}

func (f *FilterAction) GetAction() *Action {
	return f.action
}

func (f *FilterAction) GetCondition() Condition {
	return f.condition
}

func (f *FilterAction) CheckCondition(game *Game, action *Action) bool {
	return f.condition == nil || f.condition.CheckCondition(game, action)
}

func (f *FilterAction) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.CheckCondition(game, action)
}

func (f *FilterAction) InstantiateFilterPrototype(game *Game, reason *Action, instantiate bool) Filter {
	var condition Condition
	var removeCondition Condition
	if f.condition != nil {
		condition = f.condition.InstantiateFilterPrototypeCondition(game, reason, instantiate)		
	}
	if f.removeCondition != nil {
		removeCondition = f.removeCondition.InstantiateFilterPrototypeCondition(game, reason, instantiate)		
	}
	action := f.action.InstantiateFilterPrototypeAction(game, reason, instantiate)
	return &FilterAction{f.actionType, condition, removeCondition, action}
}

func InstantiateFilterSourcePrototype(game *Game, reason *Action, parameter Source, instantiate bool) Source {
	switch t := parameter.(type) {
		case []Filter:
			filters := make([]Filter, 0, len(t))
			for _,element := range t {
				filters = append(filters, element.InstantiateFilterPrototype(game, reason, instantiate))
			}
			return filters
		case []*Action:
			actions := make([]*Action, 0, len(t))
			for _,element := range t {
				instantiatedAction := element.InstantiateFilterPrototypeAction(game, reason, instantiate)
				if instantiatedAction != nil {
					actions = append(actions, instantiatedAction)
				}
			}
			return actions
		case InstantiationOff:
			if instantiate {
				return InstantiateFilterSourcePrototype(game, reason, t.source, false)
			} else {
				return InstantiationOff{InstantiateFilterSourcePrototype(game, reason, t.source, instantiate)}
			}
		case InstantiationOn:
			if !instantiate {
				return InstantiateFilterSourcePrototype(game, reason, t.source, true)
			} else {
				return InstantiationOn{InstantiateFilterSourcePrototype(game, reason, t.source, instantiate)}
			}
		case TypeOf:
			if !instantiate {
				return TypeOf{t.source}
			}
			instantiatedSource := InstantiateFilterSourcePrototype(game, reason, t.source, instantiate)
			if all,ok := instantiatedSource.(AllOf) ; ok {
				results := make([]Source, 0, len(all.Sources))
				for _,source := range all.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TypeOf{source}, instantiate))
				}
				return AllOf{results}
			}
			if oneOf,ok := instantiatedSource.(OneOf) ; ok {
				results := make([]Source, 0, len(oneOf.Sources))
				for _,source := range oneOf.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TypeOf{source}, instantiate))
				}
				return OneOf{results}
			}
			switch instantiatedSource.(type) {
				case *Creature:
					return TYPE_CREATURE
				case *Property:
					return TYPE_PROPERTY
				default:
					return nil
			}
		case Accessor:
			if !instantiate {
				return Accessor{InstantiateFilterSourcePrototype(game, reason, t.source, instantiate), t.mode}
			}
			instantiatedSource := InstantiateFilterSourcePrototype(game, reason, t.source, instantiate)
			if all,ok := instantiatedSource.(AllOf) ; ok {
				results := make([]Source, 0, len(all.Sources))
				for _,source := range all.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, Accessor{source, t.mode}, instantiate))
				}
				return AllOf{results}
			}
			if oneOf,ok := instantiatedSource.(OneOf) ; ok {
				results := make([]Source, 0, len(oneOf.Sources))
				for _,source := range oneOf.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, Accessor{source, t.mode}, instantiate))
				}
				return OneOf{results}
			}
			switch t.mode {
				case ACCESSOR_MODE_ONE_OF_CREATURE_PROPERTIES:
					creature := instantiatedSource.(*Creature)
					result := make([]Source, 0, len(creature.Tail))
					for _, card := range creature.Tail {
						result = append(result, card.ActiveProperty)
					}
					return OneOf{result}
				case ACCESSOR_MODE_CREATURE_OWNER:
					creature := instantiatedSource.(*Creature)
					return creature.Owner
				case ACCESSOR_MODE_PROPERTY_OWNER:
					property := instantiatedSource.(*Property)
					return property.ContainingCard.Owners[0]
				case ACCESSOR_MODE_CREATURES:
					player := instantiatedSource.(*Player)
					result := make([]Source, 0, len(player.Creatures))
					for _,creature := range player.Creatures {
						result = append(result, creature)
					}
					return OneOf{result}
				case ACCESSOR_MODE_TRAITS:
					container := instantiatedSource.(WithTraits)
					return container.GetTraits()
				default:
					return nil
			}
		case TraitsCount:
			instantiatedSource := InstantiateFilterSourcePrototype(game, reason, t.source, instantiate)
			instantiatedTrait := InstantiateFilterSourcePrototype(game, reason, t.traits, instantiate)
			if !instantiate {
				return TraitsCount{instantiatedSource, instantiatedTrait}
			}
			if all,ok := instantiatedSource.(AllOf) ; ok {
				results := make([]Source, 0, len(all.Sources))
				for _,source := range all.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TraitsCount{source, instantiatedTrait}, instantiate))
				}
				return AllOf{results}
			}
			if oneOf,ok := instantiatedSource.(OneOf) ; ok {
				results := make([]Source, 0, len(oneOf.Sources))
				for _,source := range oneOf.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TraitsCount{source, instantiatedTrait}, instantiate))
				}
				return OneOf{results}
			}
			if oneOfTrait, ok := instantiatedTrait.(OneOf) ; ok {
				results := make([]Source, 0, len(oneOfTrait.Sources))
				for _,trait := range oneOfTrait.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TraitsCount{instantiatedSource, trait}, instantiate))
				}
				return OneOf{results}
			}
			if allOfTrait, ok := instantiatedTrait.(AllOf) ; ok {
				count := 0
				for _,sourceTrait := range instantiatedSource.(WithTraits).GetTraits() {
					for _,trait := range allOfTrait.Sources {
						if sourceTrait == trait.(TraitType) {
							count ++
						}
					}
				
				}
				return count
			}
			count := 0
			for _,sourceTrait := range instantiatedSource.(WithTraits).GetTraits() {
				if sourceTrait == instantiatedTrait.(TraitType) {
					count ++
				}
			}
			return count
		case *Action:
			return t.InstantiateFilterPrototypeAction(game, reason, instantiate)
		case Filter:
			return t.InstantiateFilterPrototype(game, reason, instantiate)
		case Condition:
			return t.InstantiateFilterPrototypeCondition(game, reason, instantiate)
		case FilterSourcePrototype:
			if !instantiate {
				return t
			}
			switch parameter {
				case FILTER_SOURCE_PARAMETER_PLAYER:
					return reason.Arguments[PARAMETER_PLAYER]
				case FILTER_SOURCE_PARAMETER_PROPERTY:
					return reason.Arguments[PARAMETER_PROPERTY]
				case FILTER_SOURCE_PARAMETER_CREATURE:
					return reason.Arguments[PARAMETER_CREATURE]
				case FILTER_SOURCE_PARAMETER_SOURCE_CREATURE:
					return reason.Arguments[PARAMETER_SOURCE_CREATURE]
				case FILTER_SOURCE_PARAMETER_TARGET_CREATURE:
					return reason.Arguments[PARAMETER_TARGET_CREATURE]
				case FILTER_SOURCE_PARAMETER_TRAIT:
					return reason.Arguments[PARAMETER_TRAIT]
				case FILTER_SOURCE_PARAMETER_ALL_PLAYERS:
					result := make([]Source, 0, game.PlayersCount)
					game.Players.Do(func (player interface{}) {
						result = append(result, player.(*Player))
					})
					return AllOf{result}
				case FILTER_SOURCE_PARAMETER_LEFT_CREATURE:
					return reason.Arguments[PARAMETER_PAIR].([]*Creature)[0]
				case FILTER_SOURCE_PARAMETER_RIGHT_CREATURE:
					return reason.Arguments[PARAMETER_PAIR].([]*Creature)[1]
				case FILTER_SOURCE_PARAMETER_PAIR:
					return reason.Arguments[PARAMETER_PAIR]		
				case FILTER_SOURCE_PARAMETER_ANY_FOOD:
					return OneOf{[]Source{TRAIT_FOOD, TRAIT_ADDITIONAL_FOOD}}
				case FILTER_SOURCE_PARAMETER_ALL_FOOD:
					return AllOf{[]Source{TRAIT_FOOD, TRAIT_ADDITIONAL_FOOD}}
				case FILTER_SOURCE_PARAMETER_ALL_FOOD_AND_FAT:
					return AllOf{[]Source{TRAIT_FOOD, TRAIT_ADDITIONAL_FOOD, TRAIT_FAT}}
				case FILTER_SOURCE_PARAMETER_FOOD_AND_FAT_LIMIT:
					return AllOf{[]Source{TRAIT_REQUIRE_FOOD, TRAIT_FAT_TISSUE}}
				case FILTER_SOURCE_PARAMETER_PHASE:
					return reason.Arguments[PARAMETER_PHASE]
				case FILTER_SOURCE_PARAMETER_SOURCE:
					return reason.Arguments[PARAMETER_SOURCE]
				case FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT:
					return game.Food
				case FILTER_SOURCE_PARAMETER_CURRENT_PLAYER:
					return game.CurrentPlayer
				case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS:
					result := make([]Source, 0, len(game.CurrentPlayer.Cards))
					for _, card := range game.CurrentPlayer.Cards {
						result = append(result, card)
					}
					return OneOf{result}
				case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES:
					result := make([]Source, 0, len(game.CurrentPlayer.Cards)*2)
					for _, card := range game.CurrentPlayer.Cards {
						for _, property := range card.Properties {
							result = append(result, property)
						}
					}
					return OneOf{result}
				case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES_PAIR:
					result := make([]Source, 0, len(game.CurrentPlayer.Creatures))
					for _, first := range game.CurrentPlayer.Creatures {
						for _, second := range game.CurrentPlayer.Creatures {
							if first != second {
								result = append(result, []*Creature{first, second})
							}
						}
					}
					return OneOf{result}
				case FILTER_SOURCE_PARAMETER_ONE_OF_CREATURES:
					result := make([]Source, 0, 10)
					game.Players.Do(
						func (val interface{}) {
							for _, creature := range val.(*Player).Creatures {
								result = append(result, creature)
							}
					})
					return OneOf{result}
				case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES:
					result := make([]Source, 0, len(game.CurrentPlayer.Cards)*2)
					for _, creature := range game.CurrentPlayer.Creatures {
						result = append(result, creature)
					}
					return OneOf{result}
				case FILTER_SOURCE_PARAMETER_BANK_CARDS_COUNT:
					return len(game.Deck)
				default:
					panic("Unknown filter source parameter")
			}
		default:
			return parameter
	}
}