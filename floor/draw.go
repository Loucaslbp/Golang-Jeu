package floor

import (
	_ "fmt"
	"image"
	"time"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"

	"github.com/hajimehoshi/ebiten/v2"
)

// Draw affiche dans une image (en général, celle qui représente l'écran),
// la partie du sol qui est visible (qui doit avoir été calculée avec Get avant).
func (f Floor) Draw(screen *ebiten.Image, xOffset, yOffset int) {

	milisec := int(time.Now().UnixMilli()) % 2000

	for y := range f.content {
		for x := range f.content[y] {
			if f.content[y][x] == -1 {
				continue
			}

			// minus 1 because f.content is 2 tiles larger than NumTile
			screenX := float64((x-1)*configuration.Global.TileSize + xOffset)
			screenY := float64((y-1)*configuration.Global.TileSize + yOffset)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(screenX, screenY)

			shiftX := f.content[y][x] * configuration.Global.TileSize
			shiftY := 0
			
			if f.content[y][x] == BlockWater {
				shiftY = (milisec / (2000/16)) * configuration.Global.TileSize
			}

			screen.DrawImage(assets.FloorImage.SubImage(
				image.Rect(shiftX, shiftY, shiftX+configuration.Global.TileSize, shiftY + configuration.Global.TileSize),
			).(*ebiten.Image), op)

			// don't draw decoration if current block is Grass
			if f.content[y][x] == BlockGrass {
				continue
			}

			// grass decoration up
			if y < len(f.content)-1 && f.content[y+1][x] == BlockGrass {
				screen.DrawImage(assets.GrassImages[0].SubImage(
					image.Rect(0, 0, shiftX+configuration.Global.TileSize, configuration.Global.TileSize),
				).(*ebiten.Image), op)
			}

			// decoration right
			if x < len(f.content[0])-1 && f.content[y][x+1] == BlockGrass {
				screen.DrawImage(assets.GrassImages[3].SubImage(
					image.Rect(0, 0, shiftX+configuration.Global.TileSize, configuration.Global.TileSize),
				).(*ebiten.Image), op)
			}

			// decoration left
			if x > 0 && f.content[y][x-1] == BlockGrass {
				screen.DrawImage(assets.GrassImages[1].SubImage(
					image.Rect(0, 0, shiftX+configuration.Global.TileSize, configuration.Global.TileSize),
				).(*ebiten.Image), op)
			}

			// decoration down
			if y > 0 && f.content[y-1][x] == BlockGrass {
				screen.DrawImage(assets.GrassImages[2].SubImage(
					image.Rect(0, 0, shiftX+configuration.Global.TileSize, configuration.Global.TileSize),
				).(*ebiten.Image), op)
			}

		}
	}

}
