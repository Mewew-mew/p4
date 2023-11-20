package main

import (
	"log"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	message3 := "Je teste depuis client1 \n"

	_, err = conn.Write([]byte(message3))

	log.Println("Je suis connecté")

	if err != nil {
		log.Println("write error:", err)
		return
	}

}
