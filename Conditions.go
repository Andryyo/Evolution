// Conditions
package main

import "fmt"

type ConditionParameter int

const (
	CONDITION_PARAMETER_PLAYER ConditionParameter = iota
)

type Condition interface {
	CheckCondition(game *Game, action *Action) bool
	GoString() string
	InstantiateFilterTemplateCondition(reason *Action) Condition
}

type ANDCondition struct {
	conditions []Condition
}

func NewANDCondition(condifions ...Condition) *ANDCondition {
	return &ANDCondition{condifions}
}

func (c *ANDCondition) CheckCondition(game *Game, action *Action) bool {
	for _, condition := range c.conditions {
		if !condition.CheckCondition(game, action) {
			return false
		}
	}
	return true
}

func (c *ANDCondition) InstantiateFilterTemplateCondition(reason *Action) Condition {
	var conditions []Condition
	for _,condition := range c.conditions {
		conditions = append(conditions, condition.InstantiateFilterTemplateCondition(reason))
	}
	return NewANDCondition(conditions...)
}

func (c *ANDCondition) GoString() string {
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

func NewORCondition(condifions ...Condition) *ORCondition {
	return &ORCondition{condifions}
}

func (c *ORCondition) GoString() string {
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

func (c *ORCondition) InstantiateFilterTemplateCondition(reason *Action) Condition {
	var conditions []Condition
	for _,condition := range c.conditions {
		conditions = append(conditions, condition.InstantiateFilterTemplateCondition(reason))
	}
	return NewORCondition(conditions...)
}

type ConditionActionArguments struct {
	argumentName ArgumentName
	source Source
}

func (c *ConditionActionArguments) CheckCondition(game *Game, action *Action) bool {
	if prototype, ok := c.source.(SourcePrototype); ok {
		switch prototype {

		}
	} else {
		return c.source == action.Arguments[c.argumentName]
	}
	return true
}

func (c *ConditionActionArguments) GoString() string {
	return fmt.Sprintf("(%+v == %+v)", c.source, c.argumentName)
}

func (c *ConditionActionArguments) InstantiateFilterTemplateCondition(reason *Action) Condition {
	result := &ConditionActionArguments{c.argumentName, InstantiateFilterTemplateParameter(reason, c.source)}
	return result
}

type ConditionActionType struct {
	actionType ActionType
}

func (c *ConditionActionType) CheckCondition(game *Game, action *Action) bool {
	return c.actionType == action.Type
}

func (c *ConditionActionType) GoString() string {
	return fmt.Sprintf("(Action type %#v)", c.actionType)
}

func (c *ConditionActionType) InstantiateFilterTemplateCondition(reason *Action) Condition {
	return c
}

type NOTCondition struct {
	condition Condition
}

func (c *NOTCondition) CheckCondition(game *Game, action *Action) bool {
	return !c.condition.CheckCondition(game, action)
}

func (c *NOTCondition) GoString() string {
	return fmt.Sprintf("!%#v", c.condition)
}

func (c *NOTCondition) InstantiateFilterTemplateCondition(reason *Action) Condition {
	return c.condition.InstantiateFilterTemplateCondition(reason)
}

type ConditionTrait struct {
	source Source
	trait  TraitType
}

func (c *ConditionTrait) CheckCondition(game *Game, action *Action) bool {
	var source WithTraits
	if argumentName, ok := c.source.(ArgumentName) ; ok {		
		source = action.Arguments[argumentName].(WithTraits)
	} else {
		source = c.source.(WithTraits)
	}
	for _, t := range source.GetTraits() {
		if t == c.trait {
			return true
		}
	}
	return false
}

func (c *ConditionTrait) GoString() string {
	return fmt.Sprintf("(%#v have trait %#v)", c.source, c.trait)
}

func (c *ConditionTrait) InstantiateFilterTemplateCondition(reason *Action) Condition {
	return &ConditionTrait{InstantiateFilterTemplateParameter(reason, c.source), c.trait}
}

type ConditionPropertyName struct {
	name Source
}

func (c *ConditionPropertyName) CheckCondition(game *Game, action *Action) bool {
	property := action.Arguments[PARAMETER_PROPERTY].(Property)
	return property.Name == c.name.(string)
}

func (c *ConditionPropertyName) GoString() string {
	return fmt.Sprintf("(Property have name %#v)", c.name)
}

func (c *ConditionPropertyName) InstantiateFilterTemplateCondition(reason *Action) Condition {
	return &ConditionPropertyName{InstantiateFilterTemplateParameter(reason, c.name)}
}