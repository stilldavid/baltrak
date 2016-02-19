package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/tarm/goserial"
)

func listenToSerial(s *io.ReadWriteCloser, h *hub) {
	for {
		reader := bufio.NewReader(*s)

		// read until newline
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			// could do this more gracefully...
			panic(err)
		}

		h.broadcast <- reply
	}
}

func main() {
	port := flag.String("port", "/dev/cu.usbmodem1411", "Serial port device (defaults to /dev/cu.usbmodem1411)")
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
	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Fatal("server error:", err)
	}
}
