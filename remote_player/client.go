package remote_player

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

func Connect(
	addr string,
	self action_type.PlayerInitData,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) bool {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to", addr, ":", err)
		return false
	}

	// encode player init data
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err = enc.Encode(self)
	if err != nil {
		fmt.Println("Erreur encoding init player init data")
		return false
	}

	// send player init data
	err = sendData(conn, buff.Bytes())
	if err != nil {
		fmt.Println("Erreur sending init world state")
		return false
	}

	// read server worldState
	data, err := readData(conn)
	if err != nil {
		fmt.Println("Erreur reading server worldState")
		return false
	}

	// decode server worldState
	var worldState action_type.InitWorldState
	buff = *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(&worldState)
	if err != nil {
		fmt.Println("Erreur decoding server worldState :", err)
		return false
	}

	playerNames := make([]string, len(worldState.Players))

	// spawn player in game
	for i := range worldState.Players {
		fmt.Println(worldState.Players[i])
		playerNames[i] = worldState.Players[i].Name
		gameActionChan <- game_action.SpawnPlayer{
			Name:        worldState.Players[i].Name,
			X:           worldState.Players[i].X,
			Y:           worldState.Players[i].Y,
			Orientation: worldState.Players[i].Orientation,
			Skin: worldState.Players[i].Skin,
		}
	}

	// create id
	id := rand.IntN(9223372036854775807)

	// add peer to connection manager
	packerChan := make(chan packet, 10)
	newConnChan <- ConnectionRegister{
		id: id,
		isServer:       true,
		playerNames:    playerNames,
		packetChan:     packerChan,
	}

	go handleConnection(self.Name, conn, id, gameActionChan, actionOutChan, packerChan, connActionChan)

	return true
}
