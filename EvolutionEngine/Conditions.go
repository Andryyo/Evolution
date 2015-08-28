// Conditions
package EvolutionEngine

import (
	"fmt"
	"reflect"
	"runtime/debug"
)

type Condition interface {
	CheckCondition(game *Game, action *Action) bool
	GoString() string
	InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) Condition
}

type ANDCondition struct {
	conditions []Condition
}

func NewANDCondition(conditions ...Condition) *ANDCondition {
	return &ANDCondition{conditions}
}

func (c *ANDCondition) CheckCondition(game *Game, action *Action) bool {
	for _, condition := range c.conditions {
		if !condition.CheckCondition(game, action) {
			return false
		}
	}
	return true
}

func (c *ANDCondition) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) Condition {
	if c == nil {
		return c
	}
	var conditions []Condition
	for _,condition := range c.conditions {
		conditions = append(conditions, condition.InstantiateFilterPrototypeCondition(game, reason, instantiate))
	}
	return NewANDCondition(conditions...)
}

func (c ANDCondition) GoString() string {
	if len(c.conditions) == 0 {
		return "()"
	}
	result := "("+c.conditions[0].GoString()
	for i := 1; i < len(c.conditions) ; i++ {
		result += "&&"+c.conditions[i].GoString()
	}
	result += ")"
	return result
}

type ORCondition struct {
	conditions []Condition
}

func (c *ORCondition) CheckCondition(game *Game, action *Action) bool {
	for _, condition := range c.conditions {
		if condition.CheckCondition(game, action) {
			return true
		}
	}
	return false
}

func NewORCondition(conditions ...Condition) *ORCondition {
	return &ORCondition{conditions}
}

func (c ORCondition) GoString() string {
	if len(c.conditions) == 0 {
		return "()"
	}
	result := "("+c.conditions[0].GoString()
	for i := 1; i < len(c.conditions) ; i++ {
		result += "||"+c.conditions[i].GoString()
	}
	result += ")"
	return result
}

func (c *ORCondition) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) Condition {
	if c == nil {
		return c
	}
	var conditions []Condition
	for _,condition := range c.conditions {
		conditions = append(conditions, condition.InstantiateFilterPrototypeCondition(game, reason, instantiate))
	}
	return NewORCondition(conditions...)
}

type ConditionEqual struct {
	deepEqual bool
	sources []Source
}

func (c *ConditionEqual) CheckCondition(game *Game, action *Action) bool {
	instantiatedCondition := c.InstantiateFilterPrototypeCondition(game, action, true).(Condition)
	if condition, ok := instantiatedCondition.(*ConditionEqual) ; !ok {
		return instantiatedCondition.CheckCondition(game, action)
	} else {
		if condition.deepEqual {
			for i := 1; i<len(condition.sources) ; i++ {
				if !reflect.DeepEqual(condition.sources[i-1], condition.sources[i]) {
					return false
				}
			}	
			return true
		} else {
			for i := 1; i<len(condition.sources) ; i++ {
				if condition.sources[i-1] != condition.sources[i] {
					return false
				}
			}	
			return true
		}
	}
}

func (c ConditionEqual) GoString() string {
	return fmt.Sprintf("(Equals %+v)", c.sources)
}

func (c *ConditionEqual) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) (condition Condition) {
	if c == nil {
		return c
	}
	defer func() {
		if r := recover(); r!=nil {
			fmt.Printf("%v, %s\n", r, debug.Stack())
			condition = &ConditionFalse{}
		}
	} ()
	instantiatedSources := make([]Source, 0, len(c.sources))
	for i,source := range c.sources {
		instantiatedSource := InstantiateFilterSourcePrototype(game, reason, source, instantiate)
		instantiatedSources = append(instantiatedSources, instantiatedSource)
		if !instantiate {
			continue
		}
		switch instantiatedSource.(type) {
			case OneOf:
				oneOf := instantiatedSource.(OneOf)
				conditions := make([]Condition, 0, len(oneOf.Sources))
				for _, o := range oneOf.Sources {
					condition := &ConditionEqual{c.deepEqual, make([]Source,0,len(c.sources))}
					for _, c := range c.sources {
						condition.sources = append(condition.sources, c)
					}
					condition.sources[i] = o
					conditions = append(conditions, condition.InstantiateFilterPrototypeCondition(game, reason, instantiate))
				}
				return NewORCondition(conditions...)
			case AllOf:
				all := instantiatedSource.(AllOf)
				conditions := make([]Condition, 0, len(all.Sources))
				for _, o := range all.Sources {
					condition := &ConditionEqual{c.deepEqual, make([]Source,0,len(c.sources))}
					for _, c := range c.sources {
						condition.sources = append(condition.sources, c)
					}
					condition.sources[i] = o
					conditions = append(conditions, condition.InstantiateFilterPrototypeCondition(game, reason, instantiate))
				}
				return NewANDCondition(conditions...)
		}
	}
	return &ConditionEqual{c.deepEqual, instantiatedSources}
}

func NewConditionEqual(sources ...Source) Condition {
	return &ConditionEqual{false, sources}
}

func NewConditionDeepEqual(sources ...Source) Condition {
	return &ConditionEqual{true, sources}
}

type ConditionActionType struct {
	actionType ActionType
}

func (c *ConditionActionType) CheckCondition(game *Game, action *Action) bool {
	return c.actionType == action.Type
}

func (c ConditionActionType) GoString() string {
	return fmt.Sprintf("(Action type %#v)", c.actionType)
}

func (c *ConditionActionType) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) Condition {
	return c
}

type NOTCondition struct {
	condition Condition
}

func (c *NOTCondition) CheckCondition(game *Game, action *Action) bool {
	return !c.condition.CheckCondition(game, action)
}

func (c NOTCondition) GoString() string {
	return fmt.Sprintf("!%#v", c.condition)
}

func (c *NOTCondition) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) Condition {
	if c == nil {
		return c
	}
	return &NOTCondition{c.condition.InstantiateFilterPrototypeCondition(game, reason, instantiate)}
}

type ConditionPhase struct {
	phase PhaseType
}

func (c *ConditionPhase) CheckCondition(game *Game, action *Action) bool {
	return c.phase == game.CurrentPhase
}

func (c ConditionPhase) GoString() string {
	return fmt.Sprintf("(Game phase %#v)", c.phase)
}

func (c *ConditionPhase) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) (condition Condition) {
	return c
}

type ConditionFalse struct {
}

func (c *ConditionFalse) CheckCondition(game *Game, action *Action) bool {
	return false
}

func (c ConditionFalse) GoString() string {
	return fmt.Sprintf("(False)")
}

func (c *ConditionFalse) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) (condition Condition) {
	return c
}

type ConditionActionDenied struct {
	action *Action
}

func (c *ConditionActionDenied) CheckCondition(game *Game, action *Action) bool {
	instantiatedActions := c.action.InstantiateFilterPrototypeAction(game, action, true)
	return game.ActionDenied(instantiatedActions)
}

func (c ConditionActionDenied) GoString() string {
	return fmt.Sprintf("Action %#v denied", c.action)
}

func (c *ConditionActionDenied) InstantiateFilterPrototypeCondition(game *Game, reason *Action, instantiate bool) (condition Condition) {
	return &ConditionActionDenied{c.action.InstantiateFilterPrototypeAction(game, reason, instantiate)}
}