// Server
package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
)

const channelBufSize = 100
var maxId int = 0

type Server struct {
	messages  []string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan string
	doneCh    chan bool
	startGame chan bool
	errCh     chan error
}

func NewServer () *Server {
	messages := []string{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan string)
	doneCh := make(chan bool)
	startGame := make(chan bool)
	errCh := make(chan error)

	return &Server{
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		startGame,
		errCh,
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
	
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
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
	http.Handle("/client", websocket.Handler(onConnected))
	for {
		select {
			case <-s.startGame:
				log.Println("Starting game")
				players := make([]ChoiceMaker, 0, len(s.clients))
				for _, client := range s.clients {
					players = append(players, client)
				}
				NewGame(players...)
			case c := <-s.addCh:
				log.Println("Added new client")
				s.clients[c.id] = c
			case c := <-s.delCh:
				log.Println("Delete client")
				delete(s.clients, c.id)
			case err := <-s.errCh:
				log.Println("Error:", err.Error())
			case <-s.doneCh:
				return
		}
	}
}