package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
)

// Update met à jour les données du jeu à chaque 1/60 de seconde.
// Il faut bien faire attention à l'ordre des mises-à-jour car elles
// dépendent les unes des autres (par exemple, pour le moment, la
// mise-à-jour de la caméra dépend de celle du personnage et la définition
// du terrain dépend de celle de la caméra).
func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		configuration.Global.DebugMode = !configuration.Global.DebugMode
	}

	// move current player
	blocking := g.Floor.Blocking(g.Character.X, g.Character.Y, g.Camera.X, g.Camera.Y)
	charBlock := g.Floor.GetBlock(g.Character.X, g.Character.Y, g.Camera.X, g.Camera.Y)

	// check if other player block
	if !configuration.Global.CanStackPlayer {
		for i := range g.OtherPlayers {
			player := &g.OtherPlayers[i]

			if player.X == g.Character.X && player.Y == g.Character.Y-1 {
				blocking[0] = true
			}

			if player.X == g.Character.X+1 && player.Y == g.Character.Y {
				blocking[1] = true
			}

			if player.X == g.Character.X && player.Y == g.Character.Y+1 {
				blocking[2] = true
			}

			if player.X == g.Character.X-1 && player.Y == g.Character.Y {
				blocking[3] = true
			}
		}
	}

	prevOrientation := g.Character.Orientation
	if ebiten.IsKeyPressed(ebiten.KeyUp) {

		if hasMoved, speed := g.Character.TryMoveUp(blocking, charBlock); hasMoved {
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X,
				g.Character.Y-1,
				speed,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {

		if hasMoved, speed := g.Character.TryMoveRight(blocking, charBlock); hasMoved {
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X+1,
				g.Character.Y,
				speed,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {

		if hasMoved, speed := g.Character.TryMoveDown(blocking, charBlock); hasMoved {
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X,
				g.Character.Y+1,
				speed,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {

		if hasMoved, speed := g.Character.TryMoveLeft(blocking, charBlock); hasMoved {
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X-1,
				g.Character.Y,
				speed,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	}

	// update player
	g.Character.Update()
	for index := range g.OtherPlayers {
		g.OtherPlayers[index].Update()
	}

	charXOffset, charYOffset := g.Character.GetOffset()
	g.Camera.Update(g.Character.X, g.Character.Y, charXOffset, charYOffset)

	chunkToLoad := g.Floor.Update(g.Camera.X, g.Camera.Y)

	// ask chunks if GameType = "join"
	if configuration.Global.GameType == "join" {

		for i := range chunkToLoad {
			fmt.Println("asking chunk :", chunkToLoad[i].X, chunkToLoad[i].Y)
			g.ActionOutChan <- action_out.NewAskChunk(
				"up",
				chunkToLoad[i].X,
				chunkToLoad[i].Y,
				configuration.Global.PlayerName,
			)
		}

	} else {

		// load chunk on another goroutine to avoid bloking
		for i := range chunkToLoad {
			go LoadChunk(chunkToLoad[i].X, chunkToLoad[i].Y, "", g.ActionOutChan, g.ChunkLoadedChan)
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) && configuration.Global.NumTileX >= 5 {

		n := configuration.Global.NumTileX / 5
		if n == 0 {
			n = 1
		}

		configuration.Global.NumTileX -= n
		configuration.Global.NumTileY = int(float32(configuration.Global.NumTileX) / g.Floor.ScreenRatio)
		configuration.SetComputedFields()
		g.Floor.UpdateContentSize()
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) && configuration.Global.NumTileX <= 200 {

		n := configuration.Global.NumTileX / 5
		if n == 0 {
			n = 1
		}

		configuration.Global.NumTileX += n
		configuration.Global.NumTileY = int(float32(configuration.Global.NumTileX) / g.Floor.ScreenRatio)
		configuration.SetComputedFields()
		g.Floor.UpdateContentSize()
	}

	// execute all actions
	for {
		select {
		case action := <-g.GameActionChan:
			action.Execute(g)
		case chunk := <-g.ChunkLoadedChan:
			g.Floor.SetChunk(chunk.coord.X, chunk.coord.Y, chunk.chunk)
		default:
			return nil
		}
	}
}
