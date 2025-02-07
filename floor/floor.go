package floor

// Floor représente les données du terrain. Pour le moment
// aucun champs n'est exporté.
//
//   - content : partie du terrain qui doit être affichée à l'écran
//   - fullContent : a map to all in-memory chunk

type Coords struct {
	X, Y int
}

type Floor struct {
	content       [][]int
	fullContent   map[Coords]*Chunk
	pendingChunks []Coords
	ScreenRatio float32
}

func (f *Floor) SetChunk(x, y int, chunk Chunk) {

	// remove chunk if it was pending
	for k := range f.pendingChunks {
		if f.pendingChunks[k].X == x && f.pendingChunks[k].Y == y {
			f.pendingChunks[k] = f.pendingChunks[len(f.pendingChunks) - 1]
			f.pendingChunks = f.pendingChunks[:len(f.pendingChunks) - 1]
			break
		}
	}

	f.fullContent[Coords{x, y}] = &chunk
}

func (f *Floor) TryGetChunk(x, y int) [][]int {

	// check if chunk is in memory
	chunk, ok := f.fullContent[Coords{x, y}]
	if !ok {
		return nil
	}

	content := make([][]int, 64)
	for y := 0; y < 64; y++ {
		content[y] = make([]int, 64)
	}

	chunk.quadtree.GetContent(0, 0, content)

	return content
}
