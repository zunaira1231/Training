package code

import (
	"github.com/gorilla/websocket"
	"time"
)

// client represents a single chatting user.
type Client struct {

	// socket is the web socket for this client.
	Socket *websocket.Conn

	// send is a channel on which messages are sent.
	Send chan *Message

	// room is the room this client is chatting in.
	Room_ *Room

	// userData holds information about the user
	userData map[string]interface{}
}

func (c *Client) Read() {

	defer c.Socket.Close()

	for {
		var msg *Message
		err := c.Socket.ReadJSON(&msg)

		if err != nil {
			return
		}

	msg.When = time.Now()
	msg.Name = c.userData["name"].(string)
	c.Room_.Forward <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()

	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)

		if err != nil {
			break
		}
	}}


/*
func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		c.Room_.Forward <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}*/