// Filters
package main

type FilterType int

const (
	FILTER_DENY FilterType = iota
	FILTER_ALLOW
	FILTER_ACTION
	FILTER_MODIFY
)

type Filter interface {
	GetType() FilterType
	GetCondition() Condition
	CheckCondition(game *Game, action *Action) bool
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

type FilterAction struct {
	condition Condition
	action    *Action
}

func (f *FilterAction) GetType() FilterType {
	return FILTER_ACTION
}

func (f *FilterAction) GetCondition() Condition {
	return f.condition
}

func (f *FilterAction) CheckCondition(game *Game, action *Action) bool {
	return f.condition == nil || f.condition.CheckCondition(game, action)
}

func (f *FilterAction) GetAction() *Action {
	return f.action
}

type FilterModify interface {
	GetType() FilterType
	GetCondition() Condition
	CheckCondition(game *Game, action *Action) bool
	ModifyAction(action *Action) Action
}
