package character


const (
	orientedDown int = iota
	orientedLeft
	orientedRight
	orientedUp
)

// Character définit les charactéristiques du personnage.
// Pour le moment seules les coordonnées absolues de l'endroit
// où il se trouve sont exportées, mais vous pourrez
// ajouter des choses au besoins lors de votre développement.
//
// Les champs non exportés sont :
//   - Orientation: l'Orientation du personnage (haut, bas, gauche, droite).
//   - animationStep : l'étape d'animation (-1 ou 1, représentant l'animation
//     d'un pas à gauche ou à droite).
//   - xInc, yInc : les incréments en X et Y à réaliser après la prochaine animation.
//   - moving : l'information de si une animation est en cours ou pas.
//   - shift : la position actuelle en pixels du personnage relativement à ses
//     coordonnées absolues.
//   - animationFrameCount : le nombre d'appels à update (ou de 1/60 de seconde) qui
//     ont eu lieu depuis la dernière étape d'animation.
type Character struct {
	Name                string
	X, Y                int
	Skin 				int
	Orientation         int
	animationStep       int
	xInc, yInc          int
	moving              bool
	shift               int
	animationFrameCount int
	speed				int
}

func (c0 *Character) IsInFrontOf(c1 *Character) bool {
	if c0.Y != c1.Y {
		return c0.Y > c1.Y
	}

	if !c0.moving && c1.moving {
		return c1.Orientation == orientedUp
	}

	if c0.moving && !c1.moving {
		return c0.Orientation == orientedDown
	}

	if c0.moving && c1.moving {
		if c0.Orientation == orientedDown && c1.Orientation == orientedDown {
			return c0.shift > c1.shift
		}

		if c0.Orientation == orientedUp && c1.Orientation == orientedUp {
			return c0.shift < c1.shift
		}
	}

	return true
}

func getSpeed(currentBlock int) int {

	// if on sand block, move slow
	if currentBlock == 1 {
		return 1
	}

	return 2
}

func (c *Character) TryMoveUp(blocking [4]bool, currentBlock int) (bool, int) {
	if c.moving {
		return false, 0
	}

	c.Orientation = orientedUp

	if blocking[0] {
		return false, 0
	}

	c.yInc = -1
	c.moving = true
	c.speed = getSpeed(currentBlock)
	return true, c.speed
}

func (c *Character) TryMoveRight(blocking [4]bool, currentBlock int) (bool, int) {
	if c.moving {
		return false, 0
	}

	c.Orientation = orientedRight

	if blocking[1] {
		return false, 0
	}

	c.xInc = 1
	c.moving = true
	c.speed = getSpeed(currentBlock)
	return true, c.speed
}

func (c *Character) TryMoveDown(blocking [4]bool, currentBlock int) (bool, int) {
	if c.moving {
		return false, 0
	}

	c.Orientation = orientedDown

	if blocking[2] {
		return false, 0
	}

	c.yInc = 1
	c.moving = true
	c.speed = getSpeed(currentBlock)
	return true, c.speed
}

func (c *Character) TryMoveLeft(blocking [4]bool, currentBlock int) (bool, int) {
	if c.moving {
		return false, 0
	}

	c.Orientation = orientedLeft

	if blocking[3] {
		return false, 0
	}

	c.xInc = -1
	c.moving = true
	c.speed = getSpeed(currentBlock)
	return true, c.speed
}

func (c *Character) RemoteMove(newX, newY, speed int) {
	switch {
	case newX > c.X:
		c.X = newX - 1
		c.Orientation = orientedRight
		c.xInc = 1
	case newX < c.X:
		c.X = newX + 1
		c.Orientation = orientedLeft
		c.xInc = -1
	case newY > c.Y:
		c.Y = newY - 1
		c.Orientation = orientedDown
		c.yInc = 1
	case newY < c.Y:
		c.Y = newY + 1
		c.Orientation = orientedUp
		c.yInc = -1
	}

	c.moving = true
	c.shift = 0
	c.speed = speed

}

func (c *Character) RemoteSetOrientation(newOrientation int) {
	c.Orientation = newOrientation
}

func New(name string, x, y, skin int) Character {
	return Character{
		Name:          name,
		animationStep: 1,
		X:             x,
		Y:             y,
		Skin: skin,
	}
}

func (c *Character) GetOffset() (int, int) {
	switch c.Orientation {
	case orientedUp:
		return 0, -c.shift
	case orientedRight:
		return c.shift, 0
	case orientedDown:
		return 0, c.shift
	default:
		return -c.shift, 0
	}
}
