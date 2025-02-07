package floor

import (
	"github.com/aquilax/go-perlin"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Init initialise les structures de données internes de f.
func (f *Floor) Init(seed int64) {

	// init perlin noise
	noise = perlin.NewPerlin(2, 2, 3, seed)

	// init fullContent map
	f.fullContent = make(map[Coords]*Chunk)

	// init content
	f.UpdateContentSize()
}

func (f *Floor) UpdateContentSize() {
	// add 2 to be able to smooth flolow camera
	f.content = make([][]int, configuration.Global.NumTileY+2)
	for y := 0; y < len(f.content); y++ {
		f.content[y] = make([]int, configuration.Global.NumTileX+2)
	}
}
