package camera

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Update met à jour la position de la caméra à chaque pas
// de temps, c'est-à-dire tous les 1/60 secondes.
func (c *Camera) Update(characterPosX, characterPosY, characterXOffset, characterYOffset int) {

	switch configuration.Global.CameraMode {
	case Static:
		c.updateStatic()
	case FollowCharacter:
		c.updateFollowCharacter(characterPosX, characterPosY, characterXOffset, characterYOffset)
	}
}

// updateStatic est la mise-à-jour d'une caméra qui reste
// toujours à la position (0,0). Cette fonction ne fait donc
// rien.
func (c *Camera) updateStatic() {}

// updateFollowCharacter est la mise-à-jour d'une caméra qui
// suit toujours le personnage. Elle prend en paramètres deux
// entiers qui indiquent les coordonnées du personnage et place
// la caméra au même endroit.
func (c *Camera) updateFollowCharacter(characterPosX, characterPosY, characterXOffset, characterYOffset int) {

	minCameraX := c.X - configuration.Global.NumTileX / 5
	minCameraY := c.Y - configuration.Global.NumTileY / 5

	maxCameraX := c.X + configuration.Global.NumTileX / 5
	maxCameraY := c.Y + configuration.Global.NumTileY / 5

	if characterPosX <= minCameraX {
		if characterXOffset <= 0 {
			c.XOffset = characterXOffset
		}

		if characterPosX != minCameraX {
			c.X = characterPosX + configuration.Global.NumTileX / 5
		}
	}

	if characterPosX >= maxCameraX {
		if characterXOffset >= 0 {
			c.XOffset = characterXOffset
		}

		if characterPosX != maxCameraX {
			c.X = characterPosX - configuration.Global.NumTileX / 5
		}
	}


	if characterPosY <= minCameraY {
		if characterYOffset <= 0 {
			c.YOffset = characterYOffset
		}

		if characterPosY != minCameraY {
			c.Y = characterPosY + configuration.Global.NumTileY / 5
		}
	}

	if characterPosY >= maxCameraY {
		if characterYOffset >= 0 {
			c.YOffset = characterYOffset
		}

		if characterPosY != maxCameraY {
			c.Y = characterPosY - configuration.Global.NumTileY / 5
		}
	}

	
}
