// Filters
package main

type Filter interface {
	GetType() FilterType
	GetCondition() Condition
	CheckCondition(game *Game, action *Action) bool
	InstantiateFilterPrototype(reason *Action) Filter
}

type FilterDeny struct {
	condition Condition
}

func (f *FilterDeny) GetType() FilterType {
	return FILTER_DENY
}

func (f *FilterDeny) GetCondition() Condition {
	return f.condition
}

func (f *FilterDeny) CheckCondition(game *Game, action *Action) bool {
	return f.condition == nil || f.condition.CheckCondition(game, action)
}

func (f *FilterDeny) InstantiateFilterPrototype(reason *Action) Filter {
	return &FilterDeny{f.condition.InstantiateFilterPrototypeCondition(reason)}
}

type FilterAllow struct {
	condition Condition
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
	return f.condition == nil || f.condition.CheckCondition(game, action)
}

func (f *FilterAllow) InstantiateFilterPrototype(reason *Action) Filter {
	return &FilterAllow{f.condition.InstantiateFilterPrototypeCondition(reason),f.action.InstantiateFilterPrototypeAction(reason)}
}

type FilterAction struct {
	actionType FilterType
	condition Condition
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

func (f *FilterAction) InstantiateFilterPrototype(reason *Action) Filter {
	return &FilterAction{f.actionType, f.condition.InstantiateFilterPrototypeCondition(reason), f.action.InstantiateFilterPrototypeAction(reason)}
}

func InstantiateFilterSourcePrototype(reason *Action, parameter Source) Source {
	if _, ok := parameter.(FilterSourcePrototype) ; ok {
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
		}
	} else {
		return parameter
	}
	return nil
}