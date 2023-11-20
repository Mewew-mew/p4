package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	var connexions []net.Conn

	for i := 0; i < 2; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		defer conn.Close()
		connexions = append(connexions, conn)
		log.Println("Le client", i+1, "s'est connecté \n")
	}

	for i, conn := range connexions {
		test := bufio.NewReader(conn)
		message, err := test.ReadString('\n')
		if err != nil {
			log.Println("read error:", err)
			return
		}
		log.Println("Message lu", i+1, " :", message)
	}

	//time.Sleep(10 * time.Second)

}
