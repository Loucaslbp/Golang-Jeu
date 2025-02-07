package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type NotifyPlayerMove struct {
	to string
	Name string
	X, Y, Speed int
}

func NewNotifyPlayerMove(to string, name string, x, y, speed int) NotifyPlayerMove {
	return NotifyPlayerMove{
		to: to,
		Name: name,
		X: x,
		Y: y,
		Speed: speed,
	}
}

func (a NotifyPlayerMove) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a)
    if err != nil {
        log.Fatal("connot encode player move action :", err)
    }

	prefixedData := append([]byte{action_type.PlayerMoveId}, buff.Bytes()...)
	return prefixedData

}

func (a NotifyPlayerMove) To() string {
	return a.to
}