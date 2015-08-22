// Server
package EvolutionServer

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
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
	startGame chan bool
	errCh     chan error
	joinLobbyCh chan struct {client *Client; lobby *GameLobby; playerId *string}
	newLobbyCh chan struct{}
}

func NewServer () *Server {
	server := &Server{}
	server.gameLobbies = make(map[int]*GameLobby)
	server.messages = []string{}
	server.clients = make(map[int]*Client)
	server.addCh = make(chan *Client)
	server.delCh = make(chan *Client)
	server.doneCh = make(chan bool)
	server.startGame = make(chan bool)
	server.errCh = make(chan error)
	server.joinLobbyCh = make(chan struct {client *Client; lobby *GameLobby; playerId *string})
	server.newLobbyCh = make(chan struct{})
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

func (s *Server) Listen() {
	log.Println("Listening...")

	go func () {
		for {
			var command string
			fmt.Scanln(&command)
			log.Println("Received server command " + command)
			switch command {
				case "Start":
					s.startGame <- true
			}
		}
	}()
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
				for _,client := range s.clients {
					client.lobbiesCh <- s.gameLobbies
				}
			case msg := <- s.joinLobbyCh:
				msg.lobby.AddClientCh <- struct{client *Client; playerId *string}{msg.client, msg.playerId}
			case c := <-s.addCh:
				log.Println("Added new client", c.id)
				s.clients[c.id] = c
				c.lobbiesCh <- s.gameLobbies
			case c := <-s.delCh:
				log.Println("Removing client")
				if c.lobby != nil {
					c.lobby.RemoveClientCh <- c
				}
				delete(s.clients, c.id)
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
}

func NewGameLobby() *GameLobby {
	lobby := &GameLobby{}
	lobby.Id = maxGameLobbyId
	maxGameLobbyId++
	lobby.Players = make(map[int]*Client)
	lobby.Observers = make(map[int]*Client)
	lobby.AddClientCh = make(chan struct {client *Client; playerId *string})
	lobby.RemoveClientCh = make(chan *Client)
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
				delete(gl.Players, client.id)
			}
		}
	}
}

func (gl *GameLobby) Add(client *Client, playerId *string) {
	client.lobby = gl
	if gl.Game == nil {
		gl.AddNewPlayer(client)
	} else if playerId == nil {
		gl.AddObserver(client)
	} else {
		gl.RestorePlayer(client, playerId)
	}
	log.Println("Adding player to lobby")
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
	if player != nil {
		gl.Players[client.id] = client
		client.SetPlayer(player)
	} else {
		gl.AddObserver(client)
	}
}