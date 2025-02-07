package game

import "gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"

// Layout détermine la taille de l'image sur laquelle Ebitengine
// affiche les images du jeu en fonction de la taille de la fenêtre
// dans laquelle l'affichage a lieu. Pour le moment, cette taille
// ne dépend pas de celle de la fenêtre.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	// ajust game NumTileX and NumTileY from screen 
	newAspectRatio := float32(outsideWidth) /  float32(outsideHeight)
	if g.Floor.ScreenRatio != newAspectRatio {
		g.Floor.ScreenRatio = newAspectRatio
		configuration.Global.NumTileY = int(float32(configuration.Global.NumTileX) / g.Floor.ScreenRatio)
		configuration.SetComputedFields()
		g.Floor.UpdateContentSize()
	}

	screenWidth = configuration.Global.ScreenWidth
	screenHeight = configuration.Global.ScreenHeight
	if configuration.Global.DebugMode {
		screenWidth += configuration.Global.NumTileForDebug * configuration.Global.TileSize
		screenHeight += configuration.Global.TileSize
	}
	return
}
