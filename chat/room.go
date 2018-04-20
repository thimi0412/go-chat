package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thimi0412/go-chat/trace"
)

type room struct {
	// forward is channel which hold messege for transfer client
	forward chan []byte
	// join is client who about to join a chatroom
	join chan *client
	// leave is client who about to leave a chatroom
	leave chan *client
	// clients hold joining all client
	clients map[*client]bool
	// tracer receive log in chatroom
	tracer trace.Tracer
}

// newRoom creat new chatroom
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
			r.tracer.Trace("Join new client")
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Leave client")
		case msg := <-r.forward:
			r.tracer.Trace("Reception　message : ", string(msg))
			// send message for all client
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
					r.tracer.Trace("Send to client")
				default:
					// fail sending message
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("Fail send : Clean up client")
				}
			}
		}
	}
}

const (
	soketBufferSize   = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  soketBufferSize,
	WriteBufferSize: soketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("SeverHTTP", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
