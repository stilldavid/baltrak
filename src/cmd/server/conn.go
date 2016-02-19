package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type tosend struct {
	Rssi  float64 `json:"rssi"`
	Count int     `json:"count"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Alt   float64 `json:"alt"`
	Spd   float64 `json:"spd"`
	Tmpi  float64 `json:"tmpint"`
	Tmpo  float64 `json:"tmpext"`
	Volts float64 `json:"volts"`
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

// -21,2,40.037207,-105.263775,1633.70,0.55,26.83,21.50,4.10
func (c *connection) writer() {
	for message := range c.send {

		split := strings.Split(string(message), ",")

		if len(split) != 9 {
			log.Println("wrong number of params.")
			return
		}

		rssi, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			log.Fatal("can't parse lat")
		}
		count, err := strconv.Atoi(split[1])
		if err != nil {
			log.Fatal("can't parse lat")
		}
		lat, err := strconv.ParseFloat(split[2], 64)
		if err != nil {
			log.Fatal("can't parse lat")
		}
		lng, err := strconv.ParseFloat(split[3], 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		alt, err := strconv.ParseFloat(split[4], 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		spd, err := strconv.ParseFloat(split[5], 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		itmp, err := strconv.ParseFloat(split[6], 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		etmp, err := strconv.ParseFloat(strings.TrimSpace(split[7]), 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}
		volts, err := strconv.ParseFloat(strings.TrimSpace(split[8]), 64)
		if err != nil {
			log.Fatal("can't parse lng")
		}

		tosend := tosend{rssi, count, lat, lng, alt, spd, itmp, etmp, volts}

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
