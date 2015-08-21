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
var maxId int = 0

type Server struct {
	messages  []string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	doneCh    chan bool
	startGame chan bool
	errCh     chan error
	game	  *EvolutionEngine.Game
	newPlayerCh chan *Client
	existingPlayerCh chan struct {client *Client; playerId string}
}

func NewServer () *Server {
	messages := []string{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	doneCh := make(chan bool)
	startGame := make(chan bool)
	errCh := make(chan error)
	newPlayerCh := make(chan *Client)
	existingPlayerCh := make(chan struct {client *Client; playerId string})

	return &Server{
		messages,
		clients,
		addCh,
		delCh,
		doneCh,
		startGame,
		errCh,
		nil,
		newPlayerCh,
		existingPlayerCh,
	}
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
			case <-s.startGame:
				log.Println("Starting game")
				if s.game == nil {
					s.game = EvolutionEngine.NewGame(len(s.clients))
				} 
				for _,client := range s.clients {
					player := s.game.GetUnoccupiedPlayer()
					log.Println(player)
					client.SetPlayer(player)
				}
				go s.game.Start()
			case c := <-s.addCh:
				log.Println("Added new client", c.id)
				s.clients[c.id] = c
			case c := <-s.delCh:
				log.Println("Delete client")
				delete(s.clients, c.id)
			case <-s.newPlayerCh:
			case val := <-s.existingPlayerCh:
				if s.game == nil {
					continue
				}
				log.Println("Connecting existing client ", val.client, " to player", val.playerId)
				s.game.Players.Do(func (p interface {}) {
					log.Println(p)
					player := p.(*EvolutionEngine.Player)
					if (fmt.Sprintf("%p", player) == val.playerId) {
						log.Println("Set player")
						val.client.SetPlayer(player)
					}
				})
			case err := <-s.errCh:
				log.Println("Error:", err.Error())
			case <-s.doneCh:
				return
		}
	}
}