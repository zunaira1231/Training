package code

import (
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"trace"

	//"trace"

	"github.com/gorilla/websocket"

)

type Room struct {

	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	Forward chan *Message

	// join is a channel for clients wishing to join the room.
	Join chan *Client

	// leave is a channel for clients wishing to leave the room.
	Leave chan *Client

	// clients holds all current clients in this room.
	Clients map[*Client]bool

	// tracer will receive trace information of activity
	// in the room.
	Tracer_ trace.Tracer
}

// newRoom makes a new room that is ready to
// go.
func NewRoom() *Room {
	return &Room{
		Forward: make(chan *Message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		Tracer_:  trace.Off(),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			// joining
			r.Clients[client] = true
			r.Tracer_.Trace("New client joined")
		case client := <-r.Leave:
			// leaving
			delete(r.Clients, client)
			close(client.Send)
			r.Tracer_.Trace("Client left")
		case msg := <-r.Forward:
			r.Tracer_.Trace("Message received: ", msg.Message)
			// forward message to all clients
			for client := range r.Clients {
				client.Send <- msg
				r.Tracer_.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &Client{
		Socket: socket,
		Send:   make(chan *Message, messageBufferSize),
		Room_:   r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.Join <- client
	defer func() { r.Leave <- client }()
	go client.Write()
	client.Read()
}
