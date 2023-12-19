package main

import (
	"bufio"
	"log"
	"strings"
)

type server struct {
	handler bufio.ReadWriter
	wait    bool
	ready   bool
	channel chan string
}

func (s *server) receive() {
	for {
		message, err := s.handler.ReadString('\n')
		if err != nil {
			return
		}
		s.channel <- strings.Replace(message, "\n", "", -1)
		log.Print("<- ", message)
	}
}

func (s *server) send(message string) {
	s.handler.WriteString(message + "\n")
	s.handler.Flush()
	log.Print("-> ", message)
}
