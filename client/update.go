package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strconv"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case titleState:
		if g.titleUpdate() {
			g.gameState++
			g.server.ready = false //on remet le server à non prêt 
		}
	case colorSelectState:
		if g.colorSelectUpdate() {
			g.gameState++
			g.server.ready = false //on remet le server à non prêt 
			g.server.wait = false //on remet le server en mode attente
		}
	case playState:
		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()// on met les coordonées dans la grille J1
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()// on met les coordonées dans la grille J2
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished { //verification de si c'est fini
				g.result = result
				g.gameState++
				if g.turn == p2Turn { 
					g.server.send(fmt.Sprint(lastXPositionPlayed, ", ", "true"))//on envoie la derniere coordonnées du J2 et la partie est fini
				}
			} else {
				if g.turn == p2Turn {
					g.server.send(fmt.Sprint(lastXPositionPlayed, ", ", "false"))//on envoie la derniere coordonnées du J2 et la partie n'est fini
				}
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
			g.server.ready = false
			g.server.wait = false
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration

	if !g.server.ready {
		select {
		case message := <-g.server.channel:
			if message == "1" {
				g.turn = p1Turn
			} else {
				g.turn = p2Turn
			}
			g.server.ready = true
		default:
			// Do nothing
		}
	}

	return g.server.ready && inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	changes := false

	select {
	case message := <-g.server.channel:
		if message == "start" {
			return true
		}
		g.p2Color, _ = strconv.Atoi(message)
	default:
		// Do nothing
	}

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		col = (col + 1) % globalNumColorCol
		changes = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
		changes = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		line = (line + 1) % globalNumColorLine
		changes = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
		changes = true
	}

	g.p1Color = line*globalNumColorLine + col

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && (g.p1Color != g.p2Color) {
		changes = true
		g.server.wait = true
	}

	if changes { // Si le joueur 1 a changé de couleur ou a validé son choix.
		g.server.send(fmt.Sprint(g.p1Color, ", ", g.server.wait)) // On envoie la couleur choisie par le joueur 1 au serveur.
	}

	return false
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion du moment où le prochain pion est joué par le joueur 2.
func (g *game) p2Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	select {
	case message := <-g.server.channel:
		position, _ := strconv.Atoi(message)
		if updated, yPos := g.updateGrid(p2Token, position); updated {
			g.turn = p1Turn
			lastXPositionPlayed = position
			lastYPositionPlayed = yPos
		}
	default:
		// Do nothing
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g game) resultUpdate() bool {
	select {
	case message := <-g.server.channel:
		if message == "start" {
			return true
		}
	default:
		// Do nothing
	}
	if !g.server.wait {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.server.send("ready")
			g.server.wait = true
		}
	}
	return false
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
