package main

import (
	"log"
	"net"
	"bufio"
)


func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	conn1, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	writer1 := bufio.NewWriter(conn1)

	conn2, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	writer2 := bufio.NewWriter(conn2)

	log.Print("Tout le monde est là")

	_, err = writer1.WriteString("go\n")
	if err != nil {
		log.Println("write error:", err)
		return
	}
	err = writer1.Flush()
	if err != nil {
		log.Println("flush error:", err)
		return
	}

	_, err = writer2.WriteString("go\n")
	if err != nil {
		log.Println("write error:", err)
		return
	}
	err = writer2.Flush()
	if err != nil {
		log.Println("flush error:", err)
		return
	}

	log.Print("Messages envoyés")

	for {}
}
