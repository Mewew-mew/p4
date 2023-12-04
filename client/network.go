package main

import (
	"log"
	"bufio"
)

func handleRead(reader *bufio.Reader, c chan bool) {
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("read error: ", err)
	}
	c <- true
}