// Filters
package main

type FilterType int

const (
	FILTER_DENY FilterType = iota
	FILTER_ALLOW
	FILTER_MODIFY
	FILTER_ACTION_REPLACE
	FILTER_ACTION_EXECUTE_BEFORE
	FILTER_ACTION_EXECUTE_AFTER
)

type FilterSourceParameter int

const (
	FILTER_SOURCE_PARAMETER_PLAYER FilterSourceParameter = iota
	FILTER_SOURCE_PARAMETER_PROPERTY 
	FILTER_SOURCE_PARAMETER_SOURCE_CREATURE
	FILTER_SOURCE_PARAMETER_TARGET_CREATURE
	FILTER_SOURCE_PARAMETER_CREATURE
)

type Filter interface {
	GetType() FilterType
	GetCondition() Condition
	CheckCondition(game *Game, action *Action) bool
	InstantiateFilterTemplate(reason *Action) Filter
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

func (f *FilterDeny) InstantiateFilterTemplate(reason *Action) Filter {
	return &FilterDeny{f.condition.InstantiateFilterTemplateCondition(reason)}
}

type FilterAllow struct {
	action *Action
}

func (f *FilterAllow) GetType() FilterType {
	return FILTER_ALLOW
}

func (f *FilterAllow) GetCondition() Condition {
	return nil
}

func (f *FilterAllow) GetAction() *Action {
	return f.action
}

func (f *FilterAllow) CheckCondition(game *Game, action *Action) bool {
	return true
}

func (f *FilterAllow) InstantiateFilterTemplate(reason *Action) Filter {
	return &FilterAllow{f.action.InstantiateFilterTemplateAction(reason)}
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

func (f *FilterAction) InstantiateFilterTemplate(reason *Action) Filter {
	return &FilterAction{f.actionType, f.condition.InstantiateFilterTemplateCondition(reason), f.action.InstantiateFilterTemplateAction(reason)}
}

func InstantiateFilterTemplateParameter(reason *Action, parameter Source) Source {
	if _, ok := parameter.(FilterSourceParameter) ; ok {
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
		}
	} else {
		return parameter
	}
	return nil
}