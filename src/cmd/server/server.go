package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tarm/goserial"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: conn}
	go c.writer()
}

func listenToSerial(s *io.ReadWriteCloser, h *hub) {
	for {
		reader := bufio.NewReader(*s)
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
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
	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
