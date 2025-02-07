package remote_player

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net"
	"strings"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

func handleServerConnection(
	selfName string,
	conn net.Conn,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) {

	// ask game to send worldState
	worldStateChan := make(chan action_type.InitWorldState)
	gameActionChan <- game_action.GetWorldState{Sender: worldStateChan}

	// receive worldState from game
	worldState := <-worldStateChan

	// encode worldState
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(worldState)
	if err != nil {
		fmt.Println("Erreur encoding init world state")
		return
	}

	// send worldState
	err = sendData(conn, buff.Bytes())
	if err != nil {
		fmt.Println("Erreur sending init world state")
		return
	}

	// read peer's playerInitData
	data, err := readData(conn)
	if err != nil {
		fmt.Println("Erreur reading perr's server worldState")
		return
	}

	// decode peer's playerInitData
	var peerPlayerInitData action_type.PlayerInitData
	buff = *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(&peerPlayerInitData)
	if err != nil {
		fmt.Println("Erreur decoding perr's server worldState :", err)
		return
	}

	// notify peer's player spawned
	if len(worldState.Players) != 0 {

		to := ""
		for i := range worldState.Players {
			to += worldState.Players[i].Name + " "
		}
		to = to[:len(to)-1]

		actionOutChan <- action_out.NewNotifyPlayerSpawn(to, action_out.PlayerInitData{
			Name:        peerPlayerInitData.Name,
			X:           peerPlayerInitData.X,
			Y:           peerPlayerInitData.Y,
			Orientation: peerPlayerInitData.Orientation,
			Skin: 		 peerPlayerInitData.Skin,
		})

	}

	// spawn player in game
	gameActionChan <- game_action.SpawnPlayer{
		Name:        peerPlayerInitData.Name,
		X:           peerPlayerInitData.X,
		Y:           peerPlayerInitData.Y,
		Orientation: peerPlayerInitData.Orientation,
		Skin: 		 peerPlayerInitData.Skin,
	}

	// create id
	id := rand.IntN(9223372036854775807)

	// add peer to connection manager
	packerChan := make(chan packet, 10)
	newConnChan <- ConnectionRegister{
		id: id,
		isServer:       false,
		playerNames:    []string{peerPlayerInitData.Name},
		packetChan:     packerChan,
	}

	handleConnection(selfName, conn, id, gameActionChan, actionOutChan, packerChan, connActionChan)
}

func acceptLoop(
	selfName string,
	listener net.Listener,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) {
	defer listener.Close()

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while acccepting connection :", err)
			continue
		}

		// Handle connection
		go handleServerConnection(selfName, conn, gameActionChan, actionOutChan, newConnChan, connActionChan)
	}
}

func StartServer(
	selfName string,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) {

	listener, err := net.Listen("tcp", "0.0.0.0:5060")
	if err != nil {
		fmt.Println("cannot start server")
		return
	}

	fmt.Println("server started on :")

	interfaces, _ := net.Interfaces()
	for _, intf := range interfaces {

		addrs, err := intf.Addrs()
		if err != nil {
			fmt.Println("Erreur lors de la récupération des adresse pour l'interface", intf.Name, err)
			continue
		}


		// Afficher les adresses
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				if ip.To4() != nil && !strings.HasPrefix(ip.String(), "127.") {
					fmt.Printf(" - %s:%s", ip.String(), "5060\n")
				}
			}
		}
	}

	go acceptLoop(selfName, listener, gameActionChan, actionOutChan, newConnChan, connActionChan)

}
