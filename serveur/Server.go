package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Server struct {
	player1 *Player
	player2 *Player
}

func main() {
	/*-- Création du serveur --*/
	s := Server{}
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	log.Println("Server started")

	/*-- Attente de deux connexions --*/
	conn1, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	s.player1 = NewPlayer(*bufio.NewReadWriter(bufio.NewReader(conn1), bufio.NewWriter(conn1)))
	log.Println("Player 1 connected")

	conn2, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	s.player2 = NewPlayer(*bufio.NewReadWriter(bufio.NewReader(conn2), bufio.NewWriter(conn2)))
	log.Println("Player 2 connected")

	/*-- Choix des couleurs --*/
	s.player1.send("1")
	s.player2.send("2")
	go s.player1.receive()
	go s.player2.receive()

	for s.player1.ready == false || s.player2.ready == false {
		select {
		case msg := <-s.player1.channel: // Si le joueur 1 a envoyé un message
			tmp := strings.Split(msg, ", ") // On récupère les données (couleur, prêt ?)
			s.player2.send(tmp[0])          // On envoie la couleur au joueur 2
			if tmp[1] == "true" {           // Si le joueur 1 est prêt
				s.player1.ready = true // On le marque comme prêt
			}

		case msg := <-s.player2.channel: // Comme le joueur 1 (inversé)
			tmp := strings.Split(msg, ", ")
			s.player1.send(tmp[0])
			if tmp[1] == "true" {
				s.player2.ready = true
			}
		default:
			// Do nothing
		}
	}
	s.broadcast("start")
	s.player2.ready = false
	s.player1.ready = false

	log.Println("Game started")

	/*-- Boucle de jeu --*/
	turn := 1
	for {
		isPlaying := true
		/*-- Jeu --*/
		for isPlaying { // Tant que la partie n'est pas finie
			if turn == 1 {
				hasPlayed := false
				for {
					select {
					case msg := <-s.player1.channel: // Si le joueur 1 a envoyé un message
						tmp := strings.Split(msg, ", ") // On récupère les données (coords, partie finie ?)
						s.player2.send(tmp[0])          // On envoie les coords au joueur 2
						if tmp[1] == "true" {           // Si la partie est finie
							isPlaying = false // On arrête la boucle
						}
						hasPlayed = true
					default:
						// Do nothing
					}
					if hasPlayed {
						break
					}
				}
			} else {
				hasPlayed := false
				for {
					select {
					case msg := <-s.player2.channel: // Comme le joueur 1 (inversé)
						tmp := strings.Split(msg, ", ")
						s.player1.send(tmp[0])
						if tmp[1] == "true" {
							isPlaying = false
						}
						hasPlayed = true
					default:
						// Do nothing
					}
					if hasPlayed {
						break
					}
				}
			}
			turn = 3 - turn // On change de joueur
			log.Println("Turn:", turn)
		}
		/*-- Fin de partie --*/
		log.Println("Game finished")

		for s.player1.ready == false || s.player2.ready == false {
			select {
			case <-s.player1.channel: // Si le joueur 1 a envoyé un message
				s.player1.ready = true // On le marque comme prêt
			case <-s.player2.channel: // Comme le joueur 1 (inversé)
				s.player2.ready = true
			default:
				// Do nothing
			}
		}
		s.player2.ready = false
		s.player1.ready = false

		/*-- Sync --*/
		s.broadcast("start")
		log.Println("Game restarted")
	}
}

func (s *Server) broadcast(msg string) {
	s.player1.send(msg)
	s.player2.send(msg)
}
