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

	message3 := "Je teste depuis client1.0 \n"

	_, err = conn.Write([]byte(message3))

	test := bufio.NewWriter(conn)

	message2 := "Hey, je teste depuis client 2.0"

	_, err = test.WriteString(message2)

	err = test.Flush()

	log.Println("Je suis connecté, moi le 1")

	if err != nil {
		log.Println("write error:", err)
		return
	}

}
