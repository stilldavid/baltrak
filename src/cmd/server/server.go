package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/tarm/goserial"
)

type history struct {
	History []sentence `json:"history"`
}

func listenToSerial(s *io.ReadWriteCloser, h *hub) {
	for {
		reader := bufio.NewReader(*s)

		// read until newline
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			// could do this more gracefully...
			log.Fatal("\n\n\nSERIAL DISCONNECTED\n\n\n")
			return
		}

		// send to channel
		h.broadcast <- reply

		writeToFile(reply)
	}
}

func histHandler(w http.ResponseWriter, r *http.Request) {
	sentences, err := readFromFile()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	hist := history{sentences}

	js, err := json.Marshal(hist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	port := flag.String("port", "/dev/cu.usbmodem1421", "Serial port device (defaults to /dev/cu.usbmodem1411)")
	flag.Parse()

	c := &serial.Config{Name: *port, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	h := newHub()
	go h.run()

	go func() {
		listenToSerial(&s, h)
	}()

	http.Handle("/", http.FileServer(http.Dir("./frontend/public/")))
	http.Handle("/ws", wsHandler{h: h})
	http.HandleFunc("/tiles/", tileHandler)
	http.HandleFunc("/history.json", histHandler)
	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Fatal("server error:", err)
	}
}
