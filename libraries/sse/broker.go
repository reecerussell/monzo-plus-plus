package sse

import (
	"log"
	"net/http"
)

type Broker struct {
	Notifier chan []byte

	newClients    chan chan []byte
	closedClients chan chan []byte
	clients       map[chan []byte]bool
}

func NewBroker() *Broker {
	b := &Broker{
		Notifier:      make(chan []byte, 1),
		newClients:    make(chan chan []byte),
		closedClients: make(chan chan []byte),
		clients:       make(map[chan []byte]bool),
	}

	go b.listen()

	return b
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
			log.Printf("Client added. %d registered clients", len(b.clients))
		case s := <-b.closedClients:
			delete(b.clients, s)
			log.Printf("Removed client. %d registered clients", len(b.clients))
		case event := <-b.Notifier:
			for clientMessageChan := range b.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	mc := make(chan []byte)
	b.newClients <- mc

	defer func() {
		b.closedClients <- mc
	}()

	notify := r.Context().Done()
	go func() {
		<-notify
		b.closedClients <- mc
	}()

	for {
		log.Printf("[BROKER]: %s\n", <-mc)
		flusher.Flush()
	}
}
