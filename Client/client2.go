package main

import (
	"bufio"
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

	message4 := "Je teste depuis client 3.0 \n"

	_, err = conn.Write([]byte(message4))

	test := bufio.NewWriter(conn)

	message5 := "Hey, je teste depuis client 4.0"

	_, err = test.WriteString(message5)

	err = test.Flush()

	log.Println("Je suis connecté")

	if err != nil {
		log.Println("write error:", err)
		return
	}

}
