// Filters
package main

import "fmt"

type Filter interface {
	GetType() FilterType
	GetCondition() Condition
	CheckCondition(game *Game, action *Action) bool
	CheckRemoveCondition(game *Game, action *Action) bool
	InstantiateFilterPrototype(game *Game, reason *Action) Filter
}

type FilterDeny struct {
	condition Condition
	removeCondition Condition
}

type SourceWrapper struct {
	source Source
}

type TraitsCount struct {
	source Source
	traits Source
}

type OneOf struct {
	Sources []Source
}

type AllOf struct {
	Sources []Source
}

func (f *FilterDeny) GetType() FilterType {
	return FILTER_DENY
}

func (f *FilterDeny) GetCondition() Condition {
	return f.condition
}

func (f *FilterDeny) CheckCondition(game *Game, action *Action) (result bool) {
	return f.condition == nil || f.condition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterDeny) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterDeny) InstantiateFilterPrototype(game *Game, reason *Action) Filter {
	var condition Condition
	var removeCondition Condition
	if f.condition != nil {
		condition = f.condition.InstantiateFilterPrototypeCondition(game, reason)
	}
	if f.removeCondition != nil {
		removeCondition = f.removeCondition.InstantiateFilterPrototypeCondition(game, reason)
	}
	return &FilterDeny{condition, removeCondition}
}

func (f *FilterDeny) GoString() string {
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

func (f *FilterAllow) GetActions() []*Action {
	return f.actions
}

func (f *FilterAllow) CheckCondition(game *Game, action *Action) bool {
	return f.condition == nil || f.condition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAllow) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAllow) InstantiateFilterPrototype(game *Game, reason *Action) Filter {
	var condition Condition
	var removeCondition Condition
	if f.condition != nil {
		condition = f.condition.InstantiateFilterPrototypeCondition(game, reason)
	}
	if f.removeCondition != nil {
		removeCondition = f.removeCondition.InstantiateFilterPrototypeCondition(game, reason)
	}
	var instantiatedActions []*Action
	for _,action := range f.actions {
		source := action.InstantiateFilterPrototypeAction(game, reason)
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
	/*if f.action.Type == ACTION_ADD_TRAIT && f.action.Arguments[PARAMETER_TRAIT].(TraitType) == TRAIT_FED {
		fmt.Printf("Checking action %#v\n for action next player", action)
		condition := f.condition.InstantiateFilterPrototypeCondition(game, action)
		fmt.Printf("%#v,%#v\n", condition, condition.CheckCondition(game, action))
	}*/
	return f.condition == nil || f.condition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAction) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAction) InstantiateFilterPrototype(game *Game, reason *Action) Filter {
	var condition Condition
	var removeCondition Condition
	if f.condition != nil {
		condition = f.condition.InstantiateFilterPrototypeCondition(game, reason)		
	}
	if f.removeCondition != nil {
		removeCondition = f.removeCondition.InstantiateFilterPrototypeCondition(game, reason)		
	}
	action := f.action.InstantiateFilterPrototypeAction(game, reason)
	return &FilterAction{f.actionType, condition, removeCondition, action}
}

func InstantiateFilterSourcePrototype(game *Game, reason *Action, parameter Source) Source {
	switch t := parameter.(type) {
		/*case []Source:
			sources := make([]Source, 0, len(t))
			for _,element := range t {
					sources = append(sources, InstantiateFilterSourcePrototype(game, reason, element))
			}
			return sources*/
		case []Filter:
			filters := make([]Filter, 0, len(t))
			for _,element := range t {
				filters = append(filters, element.InstantiateFilterPrototype(game, reason))
			}
			return filters
		case SourceWrapper:
			return t.source
		case TraitsCount:
			instantiatedSource := InstantiateFilterSourcePrototype(game, reason, t.source)
			instantiatedTrait := InstantiateFilterSourcePrototype(game, reason, t.traits)
			if all,ok := instantiatedSource.(AllOf) ; ok {
				results := make([]Source, 0, len(all.Sources))
				for _,source := range all.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TraitsCount{source, instantiatedTrait}))
				}
				return AllOf{results}
			}
			if oneOf,ok := instantiatedSource.(OneOf) ; ok {
				results := make([]Source, 0, len(oneOf.Sources))
				for _,source := range oneOf.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TraitsCount{source, instantiatedTrait}))
				}
				return OneOf{results}
			}
			if oneOfTrait, ok := instantiatedTrait.(OneOf) ; ok {
				results := make([]Source, 0, len(oneOfTrait.Sources))
				for _,trait := range oneOfTrait.Sources {
					results = append(results, InstantiateFilterSourcePrototype(game, reason, TraitsCount{instantiatedSource, trait}))
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
		case Action:
			return t.InstantiateFilterPrototypeAction(game, reason)
		case Filter:
			return t.InstantiateFilterPrototype(game, reason)
		case FilterSourcePrototype:
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
				case FILTER_SOURCE_PARAMETER_CREATURE_PROPERTIES:
					creature := reason.Arguments[PARAMETER_CREATURE].(*Creature)
					properties := make([]Property, 0, len(creature.Tail))
					for _,card := range creature.Tail {
						properties = append(properties, *card.ActiveProperty)
					}
					return properties
				case FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT:
					return game.Food
				default:
					panic("Unknown filter source parameter")
			}
		default:
			return parameter
	}
}