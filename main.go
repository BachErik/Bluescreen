package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	imagePaths []string
)

func main() {
	// Initialisiere den Zufallsgenerator
	rand.Seed(time.Now().UnixNano())

	// Sammle alle Bildpfade aus dem Unterordner "bluescreens"
	if err := filepath.Walk("bluescreens", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			imagePaths = append(imagePaths, path)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Wähle ein zufälliges Bild aus dem Array
	randomIndex := rand.Intn(len(imagePaths))
	selectedImagePath := imagePaths[randomIndex]

	// Lade das ausgewählte Bild
	img, _, err := ebitenutil.NewImageFromFile(selectedImagePath, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	// Ermittle die Bildschirmauflösung
	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()

	// Passe die Größe des Bildes an die Bildschirmauflösung an
	imgWidth, imgHeight := img.Size()
	scaleX := float64(screenWidth) / float64(imgWidth)
	scaleY := float64(screenHeight) / float64(imgHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	adjustedWidth := int(float64(imgWidth) * scale)
	adjustedHeight := int(float64(imgHeight) * scale)

	// Skaliere das Bild auf die angepasste Größe
	scaledImg, _ := ebiten.NewImage(adjustedWidth, adjustedHeight, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	scaledImg.DrawImage(img, op)

	// Erstelle das Fenster im Vollbildmodus mit angepasster Größe
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Definitiv kein Bluescreen")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	// Definiere das Spiel
	game := &Game{
		img: scaledImg,
	}

	// Starte das Spiel
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	img *ebiten.Image
}

func (g *Game) Update(screen *ebiten.Image) error {
	// Beende das Spiel, wenn "q" oder "Q" gedrückt wird oder STRG + C betätigt wird
	if ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyEscape) || (ebiten.IsKeyPressed(ebiten.KeyC) && ebiten.IsKeyPressed(ebiten.KeyControl)) {
		os.Exit(0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Zeichne das Bild auf den Bildschirm
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebiten.ScreenSizeInFullscreen()
}
