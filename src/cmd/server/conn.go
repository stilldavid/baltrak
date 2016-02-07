package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type tosend struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Spd float64 `json:"spd"`
	Alt float64 `json:"alt"`
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// The hub.
	h *hub
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.h.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {

		split := strings.Split(string(message), ",")
		lat, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			log.Fatal("can't parse lat")
		}
		lng, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		alt, err := strconv.ParseFloat(strings.TrimSpace(split[2]), 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		spd, err := strconv.ParseFloat(strings.TrimSpace(split[3]), 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}

		tosend := tosend{lat, lng, spd, alt}

		err = c.ws.WriteJSON(tosend)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	h *hub
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, h: wsh.h}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()
	go c.writer()
	c.reader()
}
