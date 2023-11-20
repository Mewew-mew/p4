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

	message4 := "Je teste depuis client2 \n"

	_, err = conn.Write([]byte(message4))

	log.Println("Je suis connecté")

	if err != nil {
		log.Println("write error:", err)
		return
	}

}
