package character

import (
	"image"
	"image/color"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw permet d'afficher le personnage dans une *ebiten.Image
// (en pratique, celle qui représente la fenêtre de jeu) en
// fonction des charactéristiques du personnage (position, orientation,
// étape d'animation, etc) et de la position de la caméra (le personnage
// est affiché relativement à la caméra).
func (c Character) Draw(screen *ebiten.Image, camX, camY, camXOffset, camYOffset int) {

	xShift := 0
	yShift := 0
	switch c.Orientation {
	case orientedDown:
		yShift = c.shift
	case orientedUp:
		yShift = -c.shift
	case orientedLeft:
		xShift = -c.shift
	case orientedRight:
		xShift = c.shift
	}

	// minus 1 for smooth follow char
	xTileForDisplay := c.X - camX + configuration.Global.ScreenCenterTileX - 1
	yTileForDisplay := c.Y - camY + configuration.Global.ScreenCenterTileY - 1
	
	xPos := (xTileForDisplay)*configuration.Global.TileSize + xShift - camXOffset
	yPos := (yTileForDisplay)*configuration.Global.TileSize - configuration.Global.TileSize/2 + 2 + yShift - camYOffset

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xPos), float64(yPos))
	op.Filter = ebiten.FilterNearest

	shiftX := configuration.Global.TileSize
	if c.moving {
		shiftX += c.animationStep * configuration.Global.TileSize
	}
	shiftY := c.Orientation * configuration.Global.TileSize

	screen.DrawImage(assets.CharacterImage.SubImage(
		image.Rect(c.Skin * 192 + shiftX, shiftY, c.Skin * 192 + shiftX+configuration.Global.TileSize, shiftY+configuration.Global.TileSize),
	).(*ebiten.Image), op)


	// draw name
	nameBound := text.BoundString(font.Minecraft, c.Name)
	nameX := (xPos + configuration.Global.TileSize /2) - nameBound.Dx() / 2
	nameY := yPos - 10

	vector.DrawFilledRect(
		screen, 
		float32(nameX - 4), 
		float32(nameY - nameBound.Dy() - 4), 
		float32(nameBound.Dx() + 8), 
		float32(nameBound.Dy() + 8), 
		color.RGBA{R: 0, G: 0, B: 0, A: 128},
		false,
	)

	text.Draw(
		screen, 
		c.Name, 
		font.Minecraft, 
		nameX, 
		nameY, 
		color.White,
	)
}
