package main

import (
	"bufio"
	"log"
	"strings"
)

type Player struct {
	handler bufio.ReadWriter
	ready   bool
	channel chan string
}

func NewPlayer(handler bufio.ReadWriter) *Player {
	return &Player{handler: handler, ready: false, channel: make(chan string)}
}

func (p *Player) receive() {
	for {
		message, err := p.handler.ReadString('\n')
		if err != nil {
			return
		}
		p.channel <- strings.Replace(message, "\n", "", -1)

		log.Print("<- ", message)
	}
}

func (p *Player) send(message string) {
	p.handler.WriteString(message + "\n")
	p.handler.Flush()

	log.Print("-> ", message)
}
