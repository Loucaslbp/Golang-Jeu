package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/remote_player"
)

func main() {

	var configFileName string
	flag.StringVar(&configFileName, "config", "config.json", "select configuration file")

	flag.Parse()

	configuration.Load(configFileName)
	PlayerSpawnX := configuration.Global.PlayerSpawnX
	PlayerSpawnY := configuration.Global.PlayerSpawnY
	playerName := configuration.Global.PlayerName
	playerSkin := configuration.Global.PlayerSkin

	if strings.Contains(playerName, " ") {
		fmt.Println("Player name cannot contains spaces")
		return 
	}

	// open connections to other player
	actionOutChan := make(chan game.ActionOut, 10)
	gameActionChan := make(chan game.GameAction, 10)
	newConnChan, connActionChan := remote_player.OpenConnections(gameActionChan, actionOutChan)

	// connect if GameType == "join"
	if configuration.Global.GameType == "join" {

		addr := configuration.Global.ServerAddr
		isConnected := remote_player.Connect(
			addr,
			action_type.PlayerInitData{
				Name: playerName,
				X:    PlayerSpawnX,
				Y:    PlayerSpawnY,
				Skin: playerSkin,
			},
			gameActionChan,
			actionOutChan,
			newConnChan,
			connActionChan,
		)

		if !isConnected {
			fmt.Println("Connot connect to server.")
			return
		}

	} else if configuration.Global.GameType == "new" {

		// create game dir
		err := os.Mkdir(configuration.Global.GameDir, 0755)
		if err != nil {
			panic("Connot create game directory")
		}

	}

	// start serveur anyway
	remote_player.StartServer(playerName, gameActionChan, actionOutChan, newConnChan, connActionChan)

	assets.Load()
	g := &game.Game{
		GameActionChan: gameActionChan,
		ActionOutChan:  actionOutChan,
		ChunkLoadedChan: make(chan game.ChunkLoaded, 10),
	}
	g.Init(configuration.Global.Seed, PlayerSpawnX, PlayerSpawnY, playerName, playerSkin)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetWindowTitle("Quadtree")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
