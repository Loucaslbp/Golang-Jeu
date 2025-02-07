package floor

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Update se charge de stocker dans la structure interne (un tableau)
// de f une représentation de la partie visible du terrain à partir
// des coordonnées absolues de la case sur laquelle se situe la
// caméra.
func (f *Floor) Update(camXPos, camYPos int) []Coords {

	for y := 0; y < len(f.content); y++ {
		for x := 0; x < len(f.content[0]); x++ {
			f.content[y][x] = -1
		}
	}

	// smooth follow char
	topLeftX := camXPos - configuration.Global.ScreenCenterTileX - 1
	topLeftY := camYPos - configuration.Global.ScreenCenterTileY - 1
	bottomRightX := topLeftX + configuration.Global.NumTileX + 1
	bottomRightY := topLeftY + configuration.Global.NumTileY + 1

	chunkStartX := topLeftX / 64
	if topLeftX < 0 {
		chunkStartX -= 1
	}

	chunkEndX := bottomRightX / 64
	if bottomRightX < 0 {
		chunkEndX -= 1
	}

	chunkStartY := topLeftY / 64
	if topLeftY < 0 {
		chunkStartY -= 1
	}

	chunkEndY := bottomRightY / 64
	if bottomRightY < 0 {
		chunkEndY -= 1
	}

	chunkToLoad := []Coords{}

	for x := chunkStartX; x <= chunkEndX; x++ {
		for y := chunkStartY; y <= chunkEndY; y++ {

			chunk, ok := f.fullContent[Coords{x, y}]
			if !ok {

				// check if chunks has already been asked
				isChunkPending := false
				for k := range f.pendingChunks {
					if f.pendingChunks[k].X == x && f.pendingChunks[k].Y == y {
						isChunkPending = true
						break
					}
				}

				if !isChunkPending {	
					coord := Coords{X: x, Y: y}
					chunkToLoad = append(chunkToLoad, coord)
					f.pendingChunks = append(f.pendingChunks, coord)
				}

				continue
			}

			chunk.quadtree.GetContent(
				topLeftX-x*64,
				topLeftY-y*64,
				f.content,
			)

		}
	}

	return chunkToLoad
}
