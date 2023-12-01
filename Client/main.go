package main

import (
	"bufio"
	"log"
	"net"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Création, paramétrage et lancement du jeu.
func main() {

	g := game{}

	conn, err := net.Dial("tcp", "localhost:8080") //récupérer l adresse de la ligne de commande//
	if err != nil {
		log.Fatal(err)
	}
	g.readChan = make(chan bool, 1)
	g.reader = bufio.NewReader(conn)

	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	//Mise en place connexion
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}
