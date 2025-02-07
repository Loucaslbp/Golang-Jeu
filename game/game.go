package game

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/camera"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/character"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/floor"
)

type ActionOut interface {
	GetData() []byte
	To() string
}

type GameAction interface {
	Execute(g *Game)
}

// Game est le type permettant de représenter les données du jeu.
// Aucun champs n'est exporté pour le moment.
//
// Les champs non exportés sont :
//   - camera : la représentation de la caméra
//   - floor : la représentation du terrain
//   - character : la représentation du personnage

type Game struct {
	Camera    camera.Camera
	Floor     floor.Floor
	Character character.Character
	OtherPlayers []character.Character
	ActionOutChan chan ActionOut
	GameActionChan chan GameAction
	ChunkLoadedChan chan ChunkLoaded
}
