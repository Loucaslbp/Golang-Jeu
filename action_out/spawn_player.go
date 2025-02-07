package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type PlayerInitData struct {
	Name string
	X, Y, Orientation int
	Skin int
}

type NotifyPlayerSpawn struct {
	to string
	player PlayerInitData
}

func NewNotifyPlayerSpawn(to string, player PlayerInitData) NotifyPlayerSpawn {
	return NotifyPlayerSpawn{
		to: to,
		player: player,
	}
}

func (a NotifyPlayerSpawn) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a.player)
    if err != nil {
        log.Fatal("connot encode player spawn action :", err)
    }

	prefixedData := append([]byte{action_type.PlayerSpawnId}, buff.Bytes()...)
	return prefixedData

}

func (a NotifyPlayerSpawn) To() string {
	return a.to
}