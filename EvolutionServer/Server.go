// Server
package EvolutionServer

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	//"fmt"
	"github.com/Andryyo/Evolution/EvolutionEngine"
)

const channelBufSize = 100
var maxClientId int = 0
var maxGameLobbyId int = 0

type Server struct {
	gameLobbies map[int]*GameLobby
	messages  []string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	doneCh    chan bool
	errCh     chan error
	joinLobbyCh chan struct {client *Client; lobby *GameLobby; playerId *string}
	newLobbyCh chan struct{}
	updateLobbiesCh chan struct{}
}

func NewServer () *Server {
	server := &Server{}
	server.gameLobbies = make(map[int]*GameLobby)
	server.messages = []string{}
	server.clients = make(map[int]*Client)
	server.addCh = make(chan *Client)
	server.delCh = make(chan *Client)
	server.doneCh = make(chan bool)
	server.errCh = make(chan error)
	server.joinLobbyCh = make(chan struct {client *Client; lobby *GameLobby; playerId *string})
	server.newLobbyCh = make(chan struct{})
	server.updateLobbiesCh = make(chan struct{})
	return server
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	log.Printf("Client %s disconnected", c.name)
	s.delCh <- c
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) UpdateChannels() {
	s.updateLobbiesCh <- struct {}{}
}

func (s *Server) Listen() {
	log.Println("Listening...")
	http.Handle("/socket", websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}))
	for {
		select {
			case <- s.newLobbyCh:
				lobby := NewGameLobby()
				s.gameLobbies[lobby.Id] = lobby
				go s.UpdateChannels()
			case msg := <- s.joinLobbyCh:
				msg.lobby.AddClientCh <- struct{client *Client; playerId *string}{msg.client, msg.playerId}
				go s.UpdateChannels()
			case c := <-s.addCh:
				log.Println("Added new client", c.id)
				s.clients[c.id] = c
				c.lobbiesCh <- s.gameLobbies
			case c := <-s.delCh:
				log.Println("Removing client")
				if c.lobby != nil {
					c.lobby.RemoveClientCh <- c
				}
				go s.UpdateChannels()
				delete(s.clients, c.id)
			case <- s.updateLobbiesCh:
				for _,client := range s.clients {
					client.lobbiesCh <- s.gameLobbies
				}
			case err := <-s.errCh:
				log.Println("Error:", err.Error())
			case <-s.doneCh:
				return
		}
	}
}

type GameLobby struct {
	Id int
	Name string
	Game *EvolutionEngine.Game
	Players map[int]*Client
	Observers map[int]*Client
	AddClientCh chan struct {client *Client; playerId *string}
	RemoveClientCh chan *Client
	VoteCh chan struct{}
}

func NewGameLobby() *GameLobby {
	lobby := &GameLobby{}
	lobby.Id = maxGameLobbyId
	maxGameLobbyId++
	lobby.Players = make(map[int]*Client)
	lobby.Observers = make(map[int]*Client)
	lobby.AddClientCh = make(chan struct {client *Client; playerId *string})
	lobby.RemoveClientCh = make(chan *Client)
	lobby.VoteCh = make(chan struct{})
	go lobby.Listen()
	return lobby
}

func (gl *GameLobby) Listen() {
	for {
		select {
		case msg := <-gl.AddClientCh:
			gl.Add(msg.client, msg.playerId)
		case client := <-gl.RemoveClientCh:
			if client.observer {
				gl.Game.RemoveObserver(client.notifyCh)
				delete(gl.Observers, client.id)
			} else {
				if client.player != nil {
					client.player.Occupied = false
				}
				delete(gl.Players, client.id)
			}
		case <-gl.VoteCh:
			allVoteYes := true
			for _, client := range gl.Players {
				if !client.voteStart {
					allVoteYes = false
				}
			}
			if allVoteYes {
				gl.StartGame()
			}
		}
	}
}

func (gl *GameLobby) StartGame() {
	if len(gl.Players) == 0 {
		return
	}
	gl.Game = EvolutionEngine.NewGame(len(gl.Players))
	for _,client := range gl.Players {
		client.voteStart = false
		player := gl.Game.GetUnoccupiedPlayer()
		log.Println(player)
		client.SetPlayer(player)
	}
	go gl.Game.Start()
}

func (gl *GameLobby) Add(client *Client, playerId *string) {
	client.lobby = gl
	client.voteCh = gl.VoteCh
	if gl.Game == nil {
		gl.AddNewPlayer(client)
	} else if playerId == nil {
		gl.AddObserver(client)
	} else {
		gl.RestorePlayer(client, playerId)
	}
}

func (gl *GameLobby) AddObserver(client *Client) {
	client.observer = true
	gl.Observers[client.id] = client
	client.notifyCh = gl.Game.AddObserver()
	client.updateChannelsCh <- struct {}{}
}

func (gl *GameLobby) AddNewPlayer(client *Client) {
	gl.Players[client.id] = client
}

func (gl *GameLobby) RestorePlayer(client *Client, playerId *string) {
	player := gl.Game.GetPlayerById(*playerId)
	log.Println("Found player ", player)
	if player != nil {
		gl.Players[client.id] = client
		client.SetPlayer(player)
	} else {
		gl.AddObserver(client)
	}
}