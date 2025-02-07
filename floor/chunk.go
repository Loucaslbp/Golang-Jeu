package floor

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"os"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/quadtree"

	"github.com/aquilax/go-perlin"
)

type Chunk struct {
	quadtree quadtree.Quadtree
}

func NewChunk(content [][]int) Chunk {
	return Chunk{
		quadtree: quadtree.MakeFromArray(content),
	} 
}

func decodeContent(data []byte) ([][]int, error) {

	content := make([][]int, 64)
	for y := 0; y < 64; y++ {
		content[y] = make([]int, 64)
	}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)

	err := dec.Decode(&content)
	if err != nil {
		return content, errors.New("error decoding chunk content")
	}

	return content, nil
}

func encodeContent(content [][]int) ([]byte, error) {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(content)
	if err != nil {
		return buff.Bytes(), errors.New("error decoding chunk content")
	}

	return buff.Bytes(), nil
}

const noiseZoom = 20.0
const ilandSize = 64
var noise *perlin.Perlin

func genChunk(chunkX, chunkY int) [][]int {

	if noise == nil {
		panic("noise has not been init")
	}

	content := make([][]int, 64)
	for y := 0; y < 64; y++ {
		content[y] = make([]int, 64)
	}

	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {

			tileX := float64(chunkX*64 + x)
			tileY := float64(chunkY*64 + y)

			noiseValue := noise.Noise2D(tileX/noiseZoom, tileY/noiseZoom)

			distance := math.Sqrt(tileX*tileX + tileY*tileY)
			noiseValue -= (distance / ilandSize) * (distance / ilandSize)

			if noiseValue < -0.3 {

				content[y][x] = 4

			} else if noiseValue < -0.1 {

				content[y][x] = 1

			} else if noiseValue < 0.3 {

				content[y][x] = 0

			} else {
				content[y][x] = 2
			}
		}
	}

	return content

}

func LoadChunkContent(x, y int, worldDir string) [][]int {

	var content [][]int

	filename := fmt.Sprintf("%s/%d-%d.chunk", worldDir, x, y)
	file, err := os.Open(filename)
	if err != nil {

		// if file doesn't exist, gen the chunk and save it
		content = genChunk(x, y)

		// encode new generated content
		data, err := encodeContent(content)
		if err != nil {
			panic(err)
		}

		// write the content
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Println(err)
			panic("Error writing chunk file")
		}

		return content

	}

	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		panic("Error reading chunk file")
	}

	// Read the file into a byte slice
	data := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(data)
	if err != nil {
		panic("Error reading chunk file")
	}

	// decode file content
	content, err = decodeContent(data)
	if err != nil {
		panic(err)
	}

	return content
}