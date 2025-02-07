package remote_player

import (
	"strings"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

type packet struct {
	to   string
	data []byte
}

type ConnAction struct {
	id     int
	action bool
	name   string // if name == all and action == bool -> conn lost
}

type ConnectionRegister struct {
	id          int
	isServer    bool
	packetChan  chan packet
	playerNames []string
}

func inCommon(a1, a2 []string) (common []string) {
	for i1 := range a1 {
		for i2 := range a2 {
			if a1[i1] == a2[i2] {
				common = append(common, a1[i1])
			}
		}
	}
	return
}

func managePeers(
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
) {
	peers := []ConnectionRegister{}

	for {
		select {
		case conn := <-newConnChan:
			peers = append(peers, conn)

		case connAction := <-connActionChan:
			for peerIndex := range peers {
				if peers[peerIndex].id != connAction.id {
					continue
				}

				if connAction.action {

					peers[peerIndex].playerNames = append(peers[peerIndex].playerNames, connAction.name)

				} else {

					// connection closed
					if connAction.name == "all" {

						// despawn all peer's players
						for i := range peers[peerIndex].playerNames {

							// notify everyone
							actionOutChan <- action_out.NewNotifyPlayerDespawn(
								"all",
								peers[peerIndex].playerNames[i],
							)

							// despawn in game
							gameActionChan <- game_action.DespawnPlayer{
								Name: peers[peerIndex].playerNames[i],
							}
						}

						// remove peer
						peers[peerIndex] = peers[len(peers)-1]
						peers = peers[:len(peers)-1]

						break
					}

					// remove player
					// player will be despawn in game and send to peers
					for i := range peers[peerIndex].playerNames {

						if peers[peerIndex].playerNames[i] == connAction.name {

							namesLen := len(peers[peerIndex].playerNames)
							peers[peerIndex].playerNames[i] = peers[peerIndex].playerNames[namesLen-1]
							peers[peerIndex].playerNames = peers[peerIndex].playerNames[:namesLen-1]

							break
						}

					}
				}
				break
			}

		case actionOut := <-actionOutChan:

			data := actionOut.GetData()
			to := actionOut.To()

			for peerIndex := range peers {

				// from actionOut to, exclude player
				// who aren't under current peer

				// if peerTo is empty, nothing to send to current peer
				peerTo := ""
				if to == "up" {

					if peers[peerIndex].isServer {
						peerTo = "up"
					}

				} else {

					peerToNames := peers[peerIndex].playerNames
					if to != "all" {

						// send to peer and every player peer is responsible of
						// make sure peer doesn't create loop breadcast
						peerToNames = inCommon(
							strings.Split(to, " "),
							peers[peerIndex].playerNames,
						)
					}

					if len(peerToNames) == 0 {
						continue
					}

					for i := range peerToNames {
						peerTo += peerToNames[i] + " "
					}

					peerTo = peerTo[:len(peerTo)-1]

				}

				// don't send anything to peer
				if len(peerTo) == 0 {
					continue
				}

				peers[peerIndex].packetChan <- packet{
					to:   peerTo,
					data: data,
				}
			}
		}
	}
}

func OpenConnections(
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
) (chan ConnectionRegister, chan ConnAction) {

	newConnChan := make(chan ConnectionRegister, 10)
	connActionChan := make(chan ConnAction, 10)

	go managePeers(newConnChan, connActionChan, gameActionChan, actionOutChan)

	return newConnChan, connActionChan
}
