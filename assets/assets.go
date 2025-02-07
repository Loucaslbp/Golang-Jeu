package assets

import (
	"bytes"
	"image"
	"log"
	"math"

	_ "embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed floor.png
var floorBytes []byte

// FloorImage contient une version compatible avec Ebitengine de l'image
// qui contient les différents éléments qui peuvent s'afficher au sol
// (herbe, sable, etc).
// Dans la version du projet qui vous est fournie, ces éléments sont des
// carrés de 16 pixels de côté. Vous pourrez changer cela si vous le voulez.
var FloorImage *ebiten.Image

//go:embed character.png
var characterBytes []byte

// CharacterImage contient une version compatible avec Ebitengine de
// l'image qui contient les différentes étapes de l'animation du
// personnage.
// Dans la version du projet qui vous est fournie, ce personnage tient
// dans un carré de 16 pixels de côté. Vous pourrez changer cela si vous
// le voulez.
var CharacterImage *ebiten.Image

//go:embed grass.png
var grassBytes []byte
var GrassImages [4]*ebiten.Image


// Load est la fonction en charge de transformer, à l'exécution du programme,
// les images du jeu en structures de données compatibles avec Ebitengine.
// Ces structures de données sont stockées dans les variables définies ci-dessus.
func Load() {
	decoded, _, err := image.Decode(bytes.NewReader(floorBytes))
	if err != nil {
		log.Fatal(err)
	}
	FloorImage = ebiten.NewImageFromImage(decoded)

	decoded, _, err = image.Decode(bytes.NewReader(characterBytes))
	if err != nil {
		log.Fatal(err)
	}
	CharacterImage = ebiten.NewImageFromImage(decoded)


	decoded, _, err = image.Decode(bytes.NewReader(grassBytes))
	if err != nil {
		log.Fatal(err)
	}
	grassBaseImage := ebiten.NewImageFromImage(decoded)

	for i := 0; i < 4; i++ {
		
		// Create new grass image 
		w, h := grassBaseImage.Bounds().Dx(), grassBaseImage.Bounds().Dy()
		img := ebiten.NewImage(w, h)
		op := &ebiten.DrawImageOptions{}

		// Déplacer l'origine au centre de l'image
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

		// Appliquer la rotation
		op.GeoM.Rotate(float64(i) * math.Pi / 2)

		// Replacer l'image au coin supérieur gauche
		op.GeoM.Translate(float64(w)/2, float64(h)/2)

		// Dessiner l'image d'origine avec la rotation
		img.DrawImage(grassBaseImage, op)

		// Sauvegarder l'image
		GrassImages[i] = img
	}
}
