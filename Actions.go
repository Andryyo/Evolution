// Actions
package main

import (
	"container/ring"
	"fmt"
)

type ActionType int

const (
	ACTION_SEQUENCE ActionType = iota
	ACTION_START_TURN
	ACTION_NEXT_PLAYER
	ACTION_PASS
	ACTION_NEW_PHASE
	ACTION_SELECT_ACTIVE_PROPERTY
	ACTION_ADD_CARD
)

type ArgumentName string

const (
	PARAMETER_PROPERTY         ArgumentName = "Property"
	PARAMETER_PHASE                         = "Phase"
	PARAMETER_PLAYER                        = "Player"
	PARAMETER_CARD                          = "Card"
	PARAMETER_ACTIONS_SEQUENCE              = "Actions sequence"
)

type Action interface {
	GetType() ActionType
	GetArguments() *map[ArgumentName]Source
	Execute(game *Game)
}

type BaseAction struct {
	arguments map[ArgumentName]Source
}

func (b *BaseAction) GetArguments() *map[ArgumentName]Source {
	return &b.arguments
}

type ActionSelectActiveProperty struct {
	BaseAction
}

func (a *ActionSelectActiveProperty) GetType() ActionType {
	return ACTION_SELECT_ACTIVE_PROPERTY
}

func (a *ActionSelectActiveProperty) Execute(game *Game) {
	a.arguments[PARAMETER_CARD].(*Card).ActiveProperty = a.arguments[PARAMETER_PROPERTY].(*Property)
}

type ActionNewPhase struct {
	BaseAction
}

func (a *ActionNewPhase) GetType() ActionType {
	return ACTION_NEW_PHASE
}

func (a *ActionNewPhase) Execute(game *Game) {
}

type ActionSequence struct {
	BaseAction
}

func (a *ActionSequence) GetType() ActionType {
	return ACTION_SEQUENCE
}

func (a *ActionSequence) Execute(game *Game) {
	for _, action := range *a.arguments[PARAMETER_ACTIONS_SEQUENCE].(*[]Action) {
		game.ExecuteAction(action)
	}
}

type ActionStartTurn struct {
	BaseAction
}

func NewActionStartTurn(game *Game, player *ring.Ring) *ActionStartTurn {
	return &ActionStartTurn{BaseAction{map[ArgumentName]Source{PARAMETER_PLAYER: player}}}
}

func (a *ActionStartTurn) GetType() ActionType {
	return ACTION_SEQUENCE
}

func (a *ActionStartTurn) Execute(game *Game) {
	game.CurrentPlayer = a.arguments[PARAMETER_PLAYER].(*ring.Ring)
	actions := game.GetAlowedActions()
	fmt.Println("Choose one action:")
	for i, action := range actions {
		fmt.Println(i, ") ", action)
	}
}

type ActionNextPlayer struct {
	BaseAction
}

func NewActionNextPlayer(game *Game) *ActionNextPlayer {
	return &ActionNextPlayer{BaseAction{map[ArgumentName]Source{PARAMETER_PLAYER: game.CurrentPlayer.Next()}}}
}

func (a *ActionNextPlayer) GetType() ActionType {
	return ACTION_NEXT_PLAYER
}

func (a *ActionNextPlayer) Execute(game *Game) {
	game.ExecuteAction(NewActionStartTurn(game, (*a.GetArguments())[PARAMETER_PLAYER].(*ring.Ring)))
}
