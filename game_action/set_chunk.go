package game_action

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/floor"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type SetChunk struct {
	X, Y int
	Content [][]int
}

func (a SetChunk) Execute(g *game.Game) {
	fmt.Println("Received chunk :", a.X, a.Y)
	g.Floor.SetChunk(a.X, a.Y, floor.NewChunk(a.Content))
}

func decodeSetChunk(data []byte) (SetChunk, error) {

	action := SetChunk{}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&action)
	if err != nil {
		fmt.Println("connot decode set chunk action:", err)
		return action, err
	}

	return action, nil

}
