package remote_player

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

type relayActionOut struct {
	data []byte
	to string
}

func (a relayActionOut) GetData() []byte { return a.data }
func (a relayActionOut) To() string { return a.to }

func sendData(conn net.Conn, data []byte) error {

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(len(data)))

	_, err := conn.Write(buf)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func readData(conn net.Conn) ([]byte, error) {
	lengthBuf := make([]byte, 8)

	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}

	packetLen := binary.BigEndian.Uint64(lengthBuf)

	dataBuf := make([]byte, packetLen)
	_, err = io.ReadFull(conn, dataBuf)
	if err != nil {
		return nil, err
	}

	return dataBuf, nil
}

func handleConnection(
	selfName string,
	conn net.Conn,
	peerId int,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	packetChan chan packet,
	connActionChan chan ConnAction,
) {
	defer conn.Close()

	// send data to peer
	go func() {
		for {
			currentPacket := <-packetChan

			// incorporate to in the packet
			toInByte := []byte(currentPacket.to)
			buf := make([]byte, 8)
			binary.BigEndian.PutUint64(buf, uint64(len(toInByte)))
			buf = append(buf, toInByte...)
			buf = append(buf, currentPacket.data...)

			// send packet
			err := sendData(conn, buf)
			if err != nil {
				fmt.Println("Error while sending data")

				// stop sending data to this peer
				// this will send to all that peer's players despawned
				connActionChan <- ConnAction{
					id: peerId,
					action: false,
					name:   "all",
				}

				return
			}
		}
	}()

	// read peer data
	for {
		data, err := readData(conn)
		if err != nil || len(data) < 8 {

			// stop sending data to this peer
			// this will send to all that peer's players despawed
			connActionChan <- ConnAction{
				id: peerId,
				action: false,
				name:   "all",
			}

			return
		}

		// decode packer destination
		toLen := binary.BigEndian.Uint64(data[:8])
		to := string(data[8 : 8+toLen])
		packet := data[8+toLen:]

		// relay package
		if to != configuration.Global.PlayerName && to != "up" {
			fmt.Println("relay packet to", to)
			actionOutChan <- relayActionOut{
				data: packet,
				to: to,
			}
		}

		// don't execute action if it is not for us
		if to == "up" || strings.Contains(to, selfName) {

			remoteAction, err := game_action.GetRemoteActionFrom(packet)
			if err != nil {
				fmt.Println("Error while parsing peer packet :", err)
				continue
			}

			if playerSpawnAction, ok := remoteAction.(game_action.SpawnPlayer); ok {
				connActionChan <- ConnAction{
					id: peerId,
					action: true,
					name: playerSpawnAction.Name,
				}
			}

			if playerDespawnAction, ok := remoteAction.(game_action.DespawnPlayer); ok {
				connActionChan <- ConnAction{
					id: peerId,
					action: false,
					name: playerDespawnAction.Name,
				}
			}
			
			gameActionChan <- remoteAction
		}
	}
}
