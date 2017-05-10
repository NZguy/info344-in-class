package main

import (
	"github.com/gorilla/websocket"
)

//Notifier represents a web sockets notifier
type Notifier struct {
	eventq  chan interface{}
	clientq chan *websocket.Conn
	clients map[*websocket.Conn]bool
	//TODO: add other fields you might need
	//such as another channel or a mutex
	//(either would work)
	//remember that go maps ARE NOT safe for
	//concurrent access, so you must do something
	//to protect the `clients` map
}

//NewNotifier constructs a new Notifer.
func NewNotifier() *Notifier {
	//TODO: create, initialize and return
	//a Notifier struct
	notifier := &Notifier{
		eventq:  make(chan interface{}),
		clientq: make(chan *websocket.Conn),
		clients: make(map[*websocket.Conn]bool),
	}

	return notifier
}

//Start begins a loop that checks for new events
//and broadcasts them to all web socket clients.
//This function should be called on a new goroutine
//e.g., `go mynotifer.Start()`
func (n *Notifier) Start() {
	//TODO: implement this function
	//this should check for new events written
	//to the `eventq` channel, and broadcast
	//them to all of the web socket clients
	for {
		select {
		case event := <-n.eventq:
			n.broadcast(event)
		case client := <-n.clientq:
			n.clients[client] = true
			go n.readPump(client)
		}

	}
	// TODO: Use a switch statement to consume users to add from the user channel
}

//AddClient adds a new web socket client to the Notifer
func (n *Notifier) AddClient(client *websocket.Conn) {
	//TODO: implement this
	//But remember that this will be called from
	//an HTTP handler, and each HTTP request is
	//processed on its own goroutine, so your
	//implementation here MUST be safe for concurrent use

	// TODO: Could send this user to a user channel which is consumed by the start function
	n.clientq <- client
}

//Notify will add a new event to the event queue
func (n *Notifier) Notify(event interface{}) {
	//TODO: add the `event` to the `eventq`
	n.eventq <- event
}

//readPump will read all messages (including control messages)
//send by the client and ignore them. This is necessary in order
//process the control messages. If you don't do this, the
//websocket will get stuck and start producing errors.
//see https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
func (n *Notifier) readPump(client *websocket.Conn) {
	//TODO: implement this according to the notes in the
	//Control Message section of the Gorilla Web Socket docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}

}

//broadcast sends the event to all client as a JSON-encoded object
func (n *Notifier) broadcast(event interface{}) {
	//TODO: Loop over all of the web socket clients in
	//n.clients and write the `event` parameter to the client
	//as a JSON-encoded object.
	//HINT: https://godoc.org/github.com/gorilla/websocket#Conn.WriteJSON
	//and for even better performance, try using a PreparedMessage:
	//https://godoc.org/github.com/gorilla/websocket#PreparedMessage
	//https://godoc.org/github.com/gorilla/websocket#Conn.WritePreparedMessage

	for key, _ := range n.clients {
		err := key.WriteJSON(event)
		if err != nil {
			//If you get an error while writing to a client,
			//the client has wandered off, so you should call
			//the `.Close()` method on the client, and delete
			//it from the n.clients map
			key.Close()
			delete(n.clients, key)
		}
	}

}
