// Conditions
package main

import "fmt"

type Condition interface {
	CheckCondition(game *Game, action *Action) bool
	GoString() string
	InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition
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

func (c *ANDCondition) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
	var conditions []Condition
	for _,condition := range c.conditions {
		conditions = append(conditions, condition.InstantiateFilterPrototypeCondition(game, reason))
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

func (c *ORCondition) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
	var conditions []Condition
	for _,condition := range c.conditions {
		conditions = append(conditions, condition.InstantiateFilterPrototypeCondition(game, reason))
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

func (c *ConditionActionArguments) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
	sources := InstantiateFilterSourcePrototype(game, reason, c.source)
	if len(sources) == 1 {
		return &ConditionActionArguments{c.argumentName, sources[0]}
	} else {
		conditions := make([]Condition, 0, len(sources))
		for _,source := range sources {
			conditions = append(conditions, &ConditionActionArguments{c.argumentName, source})
		}
		return NewANDCondition(conditions...)
	}
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

func (c *ConditionActionType) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
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

func (c *NOTCondition) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
	return &NOTCondition{c.condition.InstantiateFilterPrototypeCondition(game, reason)}
}

type ConditionTraitCountEqual struct {
	source Source
	trait  Source
	value int
}

func (c *ConditionTraitCountEqual) CheckCondition(game *Game, action *Action) bool {
	var source WithTraits
	if argumentName, ok := c.source.(ArgumentName) ; ok {		
		source = action.Arguments[argumentName].(WithTraits)
	} else {
		source = c.source.(WithTraits)
	}
	traitsCount := 0
	for _, t := range source.GetTraits() {
		if t == c.trait {
			traitsCount++
		}
	}
	return traitsCount == c.value
}

func (c *ConditionTraitCountEqual) GoString() string {
	return fmt.Sprintf("(%#v have %v traits %#v)", c.source, c.value, c.trait)
}

func (c *ConditionTraitCountEqual) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
	sources := InstantiateFilterSourcePrototype(game, reason, c.source)
	traits := InstantiateFilterSourcePrototype(game, reason, c.trait)
	if len(sources) == 1 && len(traits) == 1 {
		return &ConditionTraitCountEqual{sources[0],traits[0],c.value}
	} else {
		conditions := make([]Condition, 0, len(sources)*len(traits))
		for _,source := range sources {
			for _,trait := range traits {
				conditions = append(conditions, &ConditionTraitCountEqual{source, trait, c.value})
			}
		}
		return NewANDCondition(conditions...)
	}
}

type ConditionPhase struct {
	phase PhaseType
}

func (c *ConditionPhase) CheckCondition(game *Game, action *Action) bool {
	return c.phase == game.CurrentPhase
}

func (c *ConditionPhase) GoString() string {
	return fmt.Sprintf("(Game phase %#v)", c.phase)
}

func (c *ConditionPhase) InstantiateFilterPrototypeCondition(game *Game, reason *Action) Condition {
	return c
}