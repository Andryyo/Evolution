// Filters
package main

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

func (f *FilterDeny) GetType() FilterType {
	return FILTER_DENY
}

func (f *FilterDeny) GetCondition() Condition {
	return f.condition
}

func (f *FilterDeny) CheckCondition(game *Game, action *Action) bool {
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

type FilterAllow struct {
	condition Condition
	removeCondition Condition
	action *Action
}

func (f *FilterAllow) GetType() FilterType {
	return FILTER_ALLOW
}

func (f *FilterAllow) GetCondition() Condition {
	return f.condition
}

func (f *FilterAllow) GetAction() *Action {
	return f.action
}

func (f *FilterAllow) CheckCondition(game *Game, action *Action) bool {
	return f.condition == nil || f.condition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAllow) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAllow) InstantiateFilterPrototype(game *Game, reason *Action) Filter {
	return &FilterAllow{f.condition.InstantiateFilterPrototypeCondition(game, reason), f.removeCondition.InstantiateFilterPrototypeCondition(game, reason), f.action.InstantiateFilterPrototypeAction(game, reason)}
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
	return f.condition == nil || f.condition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAction) CheckRemoveCondition(game *Game, action *Action) bool {
	return f.removeCondition != nil && f.removeCondition.InstantiateFilterPrototypeCondition(game, action).CheckCondition(game, action)
}

func (f *FilterAction) InstantiateFilterPrototype(game *Game, reason *Action) Filter {
	return &FilterAction{
				f.actionType, 
				f.condition.InstantiateFilterPrototypeCondition(game, reason), 
				f.removeCondition.InstantiateFilterPrototypeCondition(game, reason),
				f.action.InstantiateFilterPrototypeAction(game, reason),
			}
}

func InstantiateFilterSourcePrototype(game *Game, reason *Action, parameter Source) []Source {
	if _, ok := parameter.(FilterSourcePrototype) ; ok {
		switch parameter {
			case FILTER_SOURCE_PARAMETER_PLAYER:
				return []Source{reason.Arguments[PARAMETER_PLAYER]}
			case FILTER_SOURCE_PARAMETER_PROPERTY:
				return []Source{reason.Arguments[PARAMETER_PROPERTY]}
			case FILTER_SOURCE_PARAMETER_CREATURE:
				return []Source{reason.Arguments[PARAMETER_CREATURE]}
			case FILTER_SOURCE_PARAMETER_SOURCE_CREATURE:
				return []Source{reason.Arguments[PARAMETER_SOURCE_CREATURE]}
			case FILTER_SOURCE_PARAMETER_TARGET_CREATURE:
				return []Source{reason.Arguments[PARAMETER_TARGET_CREATURE]}
			case FILTER_SOURCE_PARAMETER_TRAIT:
				return []Source{reason.Arguments[PARAMETER_TRAIT]}
			case FILTER_SOURCE_PARAMETER_ALL_PLAYERS:
				result := make([]Source, 0, game.PlayersCount)
				game.Players.Do(func 	(player interface{}) {
					result = append(result, player.(*Player))
				})
				return result
				
		}
	} else {
		return []Source{parameter}
	}
	panic("Invalid instantiation parameter")
}