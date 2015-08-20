// Client
package EvolutionServer

import (
	"golang.org/x/net/websocket"
	"io"
	"strconv"
	"fmt"
	"log"
	"github.com/Andryyo/Evolution/EvolutionEngine"
)

type MessageType int

const (
	MESSAGE_EXECUTED_ACTION MessageType = iota
	MESSAGE_CHOICES_LIST
)

type Message struct {
	Type MessageType
	Value interface{}
}

type MessageExecutedActionValue struct {
	Action ActionDTO
	State GameStateDTO
}

func NewMessageChoicesList(actions []ActionDTO) Message {
	return Message{MESSAGE_CHOICES_LIST, actions}
}

func NewMessageExecutedAction(action ActionDTO, state GameStateDTO) Message {
	return Message{MESSAGE_EXECUTED_ACTION, MessageExecutedActionValue{action, state}}
}

type Client struct {
	id int
	name string
	ws *websocket.Conn
	server *Server
	ch chan Message
	doneCh chan bool
	player *EvolutionEngine.Player
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	maxId++
	ch := make(chan Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, "", ws, server, ch, doneCh, nil}
}

func (c *Client) SetOwner(player *EvolutionEngine.Player) {
	c.player = player
}

func (c *Client) Listen() {
	c.requestName()
	c.listenWrite()
}

func (c *Client) requestName() {
	websocket.Message.Send(c.ws, "Please enter your name")
	var msg string
	err := websocket.Message.Receive(c.ws, &msg)
	log.Println("Received name " + msg)
	if err == io.EOF {
		c.doneCh <- true
	} else if err != nil {
		c.server.Err(err)
	} else {
		c.name = msg
	}
	log.Println("Client choosed name " + c.name)
}

func (c *Client) listenWrite() {
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				fmt.Println(err)
			}
		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
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

func (c *Client) NewGameStateDTO(game *EvolutionEngine.Game) GameStateDTO {
	state := GameStateDTO{}
	state.Phase = game.CurrentPhase
	state.FoodBank = game.Food
	state.CardsInDesk = len(game.Deck)
	if (c.player != nil) {
		state.PlayerCards = make([]CardDTO, 0, len(c.player.Cards))
		state.PlayerId = fmt.Sprintf("%p", c.player)
		for _,card := range c.player.Cards {
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
	c.ch <- NewMessageExecutedAction(NewActionDTO(action), c.NewGameStateDTO(game))
}

func (c *Client) GetChoice() int {
	var msg string
	err := websocket.Message.Receive(c.ws, &msg)
	if err != nil {
		c.doneCh <- true
	}
	log.Println("Received client choice " + msg)
	num,_ := strconv.Atoi(msg)
	return num
}

func (c *Client) MakeChoice(actions []*EvolutionEngine.Action) *EvolutionEngine.Action {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return actions[0]
	}
	actionsDTOs := make([]ActionDTO, 0, len(actions))
	for _, action := range actions {
		actionsDTOs = append(actionsDTOs, NewActionDTO(action))
	}
	c.ch <- Message{MESSAGE_CHOICES_LIST, actionsDTOs}
	return actions[c.GetChoice()]
}
	
func (c *Client) GetName() string {
	return c.name
}