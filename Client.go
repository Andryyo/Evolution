// Client
package main

import (
	"golang.org/x/net/websocket"
	"io"
	"strconv"
	"fmt"
	"log"
)

type Client struct {
	id int
	name string
	ws *websocket.Conn
	server *Server
	ch chan string
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	maxId++
	ch := make(chan string, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, "", ws, server, ch, doneCh}
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
			websocket.Message.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) Notify(s string) {
	c.ch <- s
}

func (c *Client) GetChoice() int {
	var msg string
	websocket.Message.Receive(c.ws, &msg)
	log.Println("Received client choice " + msg)
	num,_ := strconv.Atoi(msg)
	return num
}

func (c *Client) MakeChoice(actions []*Action) *Action {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return actions[0]
	}
	c.Notify("Choose one action:")
	for i, action := range actions {
		c.Notify(fmt.Sprintf("%v) %#v", i, action))
	}
	return actions[c.GetChoice()]
}
	
func (c *Client) GetName() string {
	return c.name
}