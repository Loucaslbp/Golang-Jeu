package quadtree

func makeFromArrayR(floorContent [][]int, topLeftX, topLeftY, width, height int) *node {

	content := floorContent[topLeftY][topLeftX]
	isLeaf := true

	// check if all case in this node has the same content
	checkLeaf:
	for x := topLeftX; x < topLeftX+width; x++ {
		for y := topLeftY; y < topLeftY+height; y++ {

			// not same content, break both loop
			if floorContent[y][x] != content {
				isLeaf = false
				break checkLeaf
			}

		}
	}

	// other properties set to 0 value
	returnNode := node{
		topLeftX: topLeftX,
		topLeftY: topLeftY,
		width:    width,
		height:   height,
		isLeaf:   isLeaf,
	}

	if isLeaf {

		// recurtion ends here
		returnNode.content = content

	} else {

		// sub divide node section
		middleWidth := width / 2
		middleHeight := height / 2

		// recurcive call to sub node
		returnNode.topLeftNode = makeFromArrayR(floorContent, topLeftX, topLeftY, middleWidth, middleHeight)
		returnNode.topRightNode = makeFromArrayR(floorContent, topLeftX+middleWidth, topLeftY, width-middleWidth, middleHeight)
		returnNode.bottomLeftNode = makeFromArrayR(floorContent, topLeftX, topLeftY+middleHeight, middleWidth, height-middleHeight)
		returnNode.bottomRightNode = makeFromArrayR(floorContent, topLeftX+middleWidth, topLeftY+middleHeight, width-middleWidth, height-middleHeight)

	}

	return &returnNode
}

// MakeFromArray construit un quadtree représentant un terrain
// étant donné un tableau représentant ce terrain.
func MakeFromArray(floorContent [][]int) Quadtree {

	// we assume floorContent is a non-empty correct 2d rectangle
	// this condition must have been verified before

	width := len(floorContent[0])
	height := len(floorContent)

	return Quadtree{
		width:  width,
		height: height,
		root:   makeFromArrayR(floorContent, 0, 0, width, height),
	}
}
