// Conditions
package main

type Condition interface {
	CheckCondition(game *Game, action *Action) bool
}

type ANDCondition struct {
	conditions []Condition
}

func (c *ANDCondition) CheckCondition(game *Game, action *Action) bool {
	for _, condition := range c.conditions {
		if !condition.CheckCondition(game, action) {
			return false
		}
	}
	return true
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

type ConditionActionArguments struct {
	arguments map[ArgumentName]Source
}

func (c *ConditionActionArguments) CheckCondition(game *Game, action *Action) bool {
	for key, argument := range c.arguments {
		if prototype, ok := argument.(SourcePrototype); ok {
			switch prototype {

			}
		} else {
			return argument == c.arguments[key]
		}
	}
	return true
}

type ConditionActionType struct {
	actionType ActionType
}

func (c *ConditionActionType) CheckCondition(game *Game, action *Action) bool {
	return c.actionType == action.Type
}
