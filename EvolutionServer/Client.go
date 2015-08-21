// Client
package EvolutionServer

import (
	"golang.org/x/net/websocket"
	//"io"
	"fmt"
	//"log"
	"github.com/Andryyo/Evolution/EvolutionEngine"
)

type MessageType int

const (
	MESSAGE_EXECUTED_ACTION MessageType = iota
	MESSAGE_CHOICES_LIST
	MESSAGE_NEW_PLAYER
	MESSAGE_EXISTING_PLAYER
	MESSAGE_NAME
	MESSAGE_CHOICE_NUM
)

type Message struct {
	Type MessageType
	Value interface{}
}

type Choice struct {
	Actions []*EvolutionEngine.Action
	Game *EvolutionEngine.Game
}

func NewMessageChoicesList(actions []ActionDTO, state GameStateDTO) Message {
	return Message{MESSAGE_CHOICES_LIST, struct{Actions []ActionDTO; State GameStateDTO}{actions, state}}
}

func NewMessageExecutedAction(action ActionDTO, state GameStateDTO) Message {
	return Message{MESSAGE_EXECUTED_ACTION, struct{Action ActionDTO; State GameStateDTO}{action, state}}
}

type Client struct {
	id int
	name string
	ws *websocket.Conn
	server *Server
	messageToSend chan Message
	receivedMessage chan Message
	choiceAvailableCh chan Choice
	choiceCh chan int
	player *EvolutionEngine.Player
}

type ClientAdapter struct {
	player *EvolutionEngine.Player
	client *Client
	choiceAvailable chan Choice
	choiceCh chan int
	choice *Choice
}

func NewClientAdapter(client *Client) *ClientAdapter{
	adapter := &ClientAdapter{}
	adapter.client = client
	adapter.choiceAvailable = make(chan Choice)
	adapter.choiceCh = make(chan int)
	client.choiceAvailableCh = adapter.choiceAvailable
	client.choiceCh = adapter.choiceCh
	return  adapter
}

func (a *ClientAdapter) SetPlayer(player *EvolutionEngine.Player) {
	a.player = player
	a.client.player = player
}

func (a *ClientAdapter) Notify(game *EvolutionEngine.Game, action *EvolutionEngine.Action) {
	if a.client != nil {
		a.client.Notify(game, action)
	}
}
func (a *ClientAdapter) MakeChoice(game *EvolutionEngine.Game, choices []*EvolutionEngine.Action) *EvolutionEngine.Action {
	if len(choices) == 0 {
		return  nil
	}
	if len(choices) == 1 {
		return choices[0]
	}

	choice := Choice{choices, game}
	a.choice = &Choice{choices, game}
	a.choiceAvailable <- choice
	choiceNum := <- a.choiceCh
	chosenAction := choice.Actions[choiceNum]
	a.choice = nil
	return chosenAction
}

func (a *ClientAdapter) GetName() string {
	if a.client != nil {
		return  a.client.name
	} else {
		return ""
	}
}

func (a *ClientAdapter) SetClient(client *Client) {
	a.client = client
	a.client.choiceAvailableCh = a.choiceAvailable
	a.client.choiceCh = a.choiceCh
	a.client.player = a.player
	if a.choice != nil {
		a.choiceAvailable <- *a.choice
	}
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	maxId++
	client := &Client{}
	client.messageToSend = make(chan Message, channelBufSize)
	client.receivedMessage = make(chan Message, channelBufSize)
	client.id = maxId
	client.ws = ws
	client.server = server
	return client
}

func (c *Client) SetOwner(player *EvolutionEngine.Player) {
	c.player = player
}

func (c *Client) Listen() {
	go c.listenRead();
	c.listenWrite()
}

func (c *Client) listenRead() {
	for {
		var msg Message
		err := websocket.JSON.Receive(c.ws, &msg)
		if err != nil {
			switch msg.Type {
			case MESSAGE_NEW_PLAYER:
				c.server.newPlayerCh <- c
			case MESSAGE_EXISTING_PLAYER:
				c.server.existingPlayerCh <- struct {client *Client; playerId string}{c, msg.Value.(string)}
			case MESSAGE_NAME:
				c.name = msg.Value.(string)
			case MESSAGE_CHOICE_NUM:
				c.choiceCh <- msg.Value.(int)
			}
		} else {

		}
	}
}

func (c *Client) listenWrite() {
	for {
		select {
		// send message to the client
		case msg := <-c.messageToSend:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				fmt.Println(err)
			}
		// receive done request
		/*case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return*/
		case choice := <-c.choiceAvailableCh:
			c.MakeChoice(choice)
		}
	}
}

type ActionDTO struct {
	Type EvolutionEngine.ActionType
	Arguments map[EvolutionEngine.ArgumentName]string
}

func NewActionDTO(action *EvolutionEngine.Action) ActionDTO{
	dto := ActionDTO{action.Type, map[EvolutionEngine.ArgumentName]string{}}
	for key,value := range action.Arguments {
		switch v := value.(type) {
			case *EvolutionEngine.Player: dto.Arguments[key] = fmt.Sprintf("%p",v)
			case *EvolutionEngine.Creature: dto.Arguments[key] = fmt.Sprintf("%p", v)
			case *EvolutionEngine.Card: dto.Arguments[key] = fmt.Sprintf("%p", v)
			case *EvolutionEngine.Property: dto.Arguments[key] = fmt.Sprintf("%p", v)
			default : dto.Arguments[key] = fmt.Sprintf("%v", v)
		}
	}
	return dto
}

type GameStateDTO struct {
	Phase EvolutionEngine.PhaseType
	FoodBank int
	CardsInDesk int
	CurrentPlayerId string
	PlayerId string
	PlayerCards []CardDTO
	Players []PlayerDTO
}

func NewGameStateDTO(player *EvolutionEngine.Player, game *EvolutionEngine.Game) GameStateDTO {
	state := GameStateDTO{}
	state.Phase = game.CurrentPhase
	state.FoodBank = game.Food
	state.CardsInDesk = len(game.Deck)
	if (player != nil) {
		state.PlayerCards = make([]CardDTO, 0, len(player.Cards))
		state.PlayerId = fmt.Sprintf("%p", player)
		for _,card := range player.Cards {
			state.PlayerCards = append(state.PlayerCards, NewCardDTO(card))
		}
	} else {
		state.PlayerCards = make([]CardDTO, 0)
	}
	state.CurrentPlayerId = fmt.Sprintf("%p", game.CurrentPlayer)
	state.Players = make([]PlayerDTO, 0, game.PlayersCount)
	game.Players.Do(func (val interface{}) {
		state.Players = append(state.Players, NewPlayerDTO(val.(*EvolutionEngine.Player)))
	})
	return state
}

type CardDTO struct {
	Id string
	ActiveProperty PropertyDTO
	Properties []PropertyDTO
}

func NewCardDTO(card *EvolutionEngine.Card) CardDTO {
	cardDTO := CardDTO{}
	cardDTO.Id = fmt.Sprintf("%p", card)
	cardDTO.ActiveProperty = NewPropertyDTO(card.ActiveProperty)
	cardDTO.Properties = make([]PropertyDTO, 0, len(card.Properties))
	for _,property := range card.Properties {
		cardDTO.Properties = append(cardDTO.Properties, NewPropertyDTO(property))
	}
	return cardDTO
}

type PropertyDTO struct {
	Id string
	Traits []EvolutionEngine.TraitType
}

func NewPropertyDTO(property *EvolutionEngine.Property) PropertyDTO {
	return PropertyDTO{fmt.Sprintf("%p",property), property.Traits}
}

type PlayerDTO struct {
	Id string
	Creatures []CreatureDTO
}

func NewPlayerDTO(player *EvolutionEngine.Player) PlayerDTO {
	playerDTO := PlayerDTO{}
	playerDTO.Id = fmt.Sprintf("%p", player)
	playerDTO.Creatures = make([]CreatureDTO, 0, len(player.Creatures))
	for _,creature := range player.Creatures {
		playerDTO.Creatures = append(playerDTO.Creatures, NewCreatureDTO(creature))
	}
	return playerDTO
}

type CreatureDTO struct {
	Id string
	Traits []EvolutionEngine.TraitType
	Cards []CardDTO
}

func NewCreatureDTO(creature *EvolutionEngine.Creature) CreatureDTO {
	creatureDTO := CreatureDTO{}
	creatureDTO.Id = fmt.Sprintf("%p", creature)
	creatureDTO.Cards = make([]CardDTO, 0, len(creature.Tail))
	for _,card := range creature.Tail {
		creatureDTO.Cards = append(creatureDTO.Cards, NewCardDTO(card))
	}
	creatureDTO.Traits = creature.Traits
	return creatureDTO
}

func (c *Client) Notify(game *EvolutionEngine.Game, action *EvolutionEngine.Action) {
	c.messageToSend <- NewMessageExecutedAction(NewActionDTO(action), NewGameStateDTO(c.player, game))
}

func (c *Client) MakeChoice(choice Choice) {
	actionsDTOs := make([]ActionDTO, 0, len(choice.Actions))
	for _, action := range choice.Actions {
		actionsDTOs = append(actionsDTOs, NewActionDTO(action))
	}
	c.messageToSend <- NewMessageChoicesList(actionsDTOs, NewGameStateDTO(c.player, choice.Game))
}
	
func (c *Client) GetName() string {
	return c.name
}