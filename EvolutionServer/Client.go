// Client
package EvolutionServer

import (
	"golang.org/x/net/websocket"
	"io"
	"fmt"
	"log"
	"strconv"
	"github.com/Andryyo/Evolution/EvolutionEngine"
)

type MessageType int

const (
	MESSAGE_EXECUTED_ACTION MessageType = iota
	MESSAGE_CHOICES_LIST
	MESSAGE_NAME
	MESSAGE_CHOICE_NUM
	MESSAGE_LOBBIES_LIST
	MESSAGE_NEW_LOBBY
	MESSAGE_JOIN_LOBBY
	MESSAGE_VOTE_START
)

type Message struct {
	Type MessageType
	Value interface{}
}

func NewMessageChoicesList(actions []ActionDTO, state GameStateDTO) Message {
	return Message{MESSAGE_CHOICES_LIST, struct{Actions []ActionDTO; State GameStateDTO}{actions, state}}
}

func NewMessageExecutedAction(action ActionDTO, state GameStateDTO) Message {
	return Message{MESSAGE_EXECUTED_ACTION, struct{Action ActionDTO; State GameStateDTO}{action, state}}
}

func NewMessageLobbiesList(lobbies []GameLobbyDTO) Message {
	return Message{MESSAGE_LOBBIES_LIST, lobbies}
}

type Client struct {
	id int
	name string
	lobby *GameLobby
	ws *websocket.Conn
	server *Server
	messageToSend chan Message
	receivedMessage chan Message
	choiceAvailableCh chan *EvolutionEngine.Choice
	choiceCh chan int
	notifyCh chan struct{Action *EvolutionEngine.Action; Game *EvolutionEngine.Game}
	doneCh   chan struct{}
	updateChannelsCh chan struct{}
	player *EvolutionEngine.Player
	lobbiesCh chan map[int]*GameLobby
	lobbies map[int]*GameLobby
	voteCh chan struct{}
	observer bool
	voteStart bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	maxClientId++
	client := &Client{}
	client.messageToSend = make(chan Message, channelBufSize)
	client.receivedMessage = make(chan Message, channelBufSize)
	client.updateChannelsCh = make(chan struct{})
	client.doneCh = make(chan struct{})
	client.id = maxClientId
	client.ws = ws
	client.server = server
	client.lobbiesCh = make(chan map[int]*GameLobby)
	return client
}

func (c *Client) SetPlayer(player *EvolutionEngine.Player) {
	if player.Occupied {
		return
	}
	log.Printf("Set player %p\n", player)
	player.Occupied = true
	c.choiceAvailableCh = player.AvailableChoiceCh
	c.choiceCh = player.ChoiceCh
	c.notifyCh = player.NotifyCh
	c.player = player
	c.updateChannelsCh <- struct{}{}
	if player.PendingChoice != nil {
		c.choiceAvailableCh <- player.PendingChoice
	}
}

func (c *Client) Listen() {
	go c.listenRead();
	c.listenWrite()
}

func (c *Client) listenRead() {
	var msg Message
	for {
		select {
			case <-c.doneCh:
				log.Println("Client disconnected")
				return
			default:
				err := websocket.JSON.Receive(c.ws, &msg)
				if err == nil {
					log.Println(msg)
					switch msg.Type {
					case MESSAGE_VOTE_START:
						c.voteStart = msg.Value.(bool)
						if c.voteCh != nil {
							c.voteCh <- struct {}{}
						}
					case MESSAGE_NEW_LOBBY:
						c.server.newLobbyCh <- struct {}{}
					case MESSAGE_JOIN_LOBBY:
						log.Println(msg.Value)
						lobbyId,_ := msg.Value.(map[string]interface{})["LobbyId"].(float64)
						log.Println("Successfully parser lobbyId ", lobbyId)
						playerId, ok := msg.Value.(map[string]interface{})["PlayerId"].(string)
						log.Println("Successfully parsed playerId ", playerId)
						if ok {
							c.server.joinLobbyCh <- struct {client *Client; lobby *GameLobby; playerId *string}{c, c.lobbies[int(lobbyId)], &playerId}
						} else {
							c.server.joinLobbyCh <- struct {client *Client; lobby *GameLobby; playerId *string}{c, c.lobbies[int(lobbyId)], nil}
						}
					case MESSAGE_NAME:
						c.name = msg.Value.(string)
					case MESSAGE_CHOICE_NUM:
						choiceNum,_ := strconv.ParseInt(msg.Value.(string), 10, 16)
						c.choiceCh <- int(choiceNum)
					}
				} else if err == io.EOF {
					log.Println("Received EOF sygnal")
					c.doneCh <- struct{}{}
				} else {
					log.Println(err)
				}
		}
	}
}

func (c *Client) listenWrite() {
	for {
		select {
		case msg := <-c.messageToSend:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				log.Println(msg, err)
			}
		case msg := <-c.notifyCh:
			message := NewMessageExecutedAction(NewActionDTO(msg.Action), NewGameStateDTO(c.player, msg.Game))
			go func() {c.messageToSend <- message}()
		case <-c.doneCh:
			log.Println("Sending client remove sygnal to server")
			if c.player != nil {
				c.player.Occupied = false
			}
			c.server.Del(c)
			c.doneCh <- struct{}{}
			return
		case choice := <-c.choiceAvailableCh:
			c.MakeChoice(choice)
		case lobbies := <-c.lobbiesCh:
			c.lobbies = lobbies
			message := NewMessageLobbiesList(NewGameLobbiesListDTO(lobbies))
			go func() {c.messageToSend <- message}()
		case <-c.updateChannelsCh:
		}
	}
}

type ActionDTO struct {
	Type EvolutionEngine.ActionType
	Arguments map[EvolutionEngine.ArgumentName]interface{}
}

func NewActionDTO(action *EvolutionEngine.Action) ActionDTO{
	dto := ActionDTO{action.Type, map[EvolutionEngine.ArgumentName]interface{}{}}
	for key,value := range action.Arguments {
		dto.Arguments[key] = EncodeActionArgument(value)
	}
	return dto
}

func EncodeActionArgument(argument interface {}) interface {} {
	switch v := argument.(type) {
		case *EvolutionEngine.Player: return fmt.Sprintf("%p",v)
		case *EvolutionEngine.Creature: return fmt.Sprintf("%p", v)
		case *EvolutionEngine.Card: return fmt.Sprintf("%p", v)
		case *EvolutionEngine.Property: return fmt.Sprintf("%p", v)
		case []*EvolutionEngine.Creature:
			result := make([]string, 0, len(v))
			for _, creature := range v {
				result = append(result, fmt.Sprintf("%p", creature))
			}
			return result
		default :  return v
	}
}

type GameStateDTO struct {
	Phase EvolutionEngine.PhaseType
	FoodBank int
	CardsInDesk int
	CurrentPlayerId string
	PlayerId string
	PlayerCards map[string]CardDTO
	Players map[string]PlayerDTO
}

func NewGameStateDTO(player *EvolutionEngine.Player, game *EvolutionEngine.Game) GameStateDTO {
	state := GameStateDTO{}
	state.Phase = game.CurrentPhase
	state.FoodBank = game.Food
	state.CardsInDesk = len(game.Deck)
	if (player != nil) {
		state.PlayerCards = make(map[string]CardDTO)
		state.PlayerId = fmt.Sprintf("%p", player)
		for _,card := range player.Cards {
			cardDTO := NewCardDTO(card)
			state.PlayerCards[cardDTO.Id] = cardDTO
		}
	} else {
		state.PlayerCards = make(map[string]CardDTO, 0)
	}
	state.CurrentPlayerId = fmt.Sprintf("%p", game.CurrentPlayer)
	state.Players = make(map[string]PlayerDTO)
	game.Players.Do(func (val interface{}) {
		player := NewPlayerDTO(val.(*EvolutionEngine.Player))
		state.Players[player.Id] = player
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
	Creatures map[string]CreatureDTO
}

func NewPlayerDTO(player *EvolutionEngine.Player) PlayerDTO {
	playerDTO := PlayerDTO{}
	playerDTO.Id = fmt.Sprintf("%p", player)
	playerDTO.Creatures = make(map[string]CreatureDTO)
	for _,creature := range player.Creatures {
		creature := NewCreatureDTO(creature)
		playerDTO.Creatures[creature.Id] = creature
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

type GameLobbyDTO struct{
	Id int
	PlayersCount int
	ObserversCount int
}

func NewGameLobbiesListDTO(lobbies map[int]*GameLobby) []GameLobbyDTO {
	result := make([]GameLobbyDTO, 0, len(lobbies))
	for _,lobby := range lobbies {
		result = append(result, NewGameLobbyDTO(lobby))
	}
	return result
}

func NewGameLobbyDTO(lobby *GameLobby) GameLobbyDTO {
	lobbyDTO := GameLobbyDTO{}
	lobbyDTO.Id = lobby.Id
	lobbyDTO.PlayersCount = len(lobby.Players)
	lobbyDTO.ObserversCount = len(lobby.Observers)
	return lobbyDTO
}

func (c *Client) MakeChoice(choice *EvolutionEngine.Choice) {
	actionsDTOs := make([]ActionDTO, 0, len(choice.Actions))
	for _, action := range choice.Actions {
		actionsDTOs = append(actionsDTOs, NewActionDTO(action))
	}
	c.messageToSend <- NewMessageChoicesList(actionsDTOs, NewGameStateDTO(c.player, choice.Game))
}
	
func (c *Client) GetName() string {
	return c.name
}