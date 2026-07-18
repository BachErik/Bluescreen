package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed bluescreens
var bluescreensFS embed.FS

type Game struct {
	img *ebiten.Image
}

func main() {
	entries, err := fs.ReadDir(bluescreensFS, "bluescreens")
	if err != nil {
		log.Fatal(err)
	}

	var paths []string
	for _, e := range entries {
		if !e.IsDir() {
			paths = append(paths, "bluescreens/"+e.Name())
		}
	}
	if len(paths) == 0 {
		log.Fatal("keine Bilder in bluescreens/ eingebettet")
	}

	data, err := bluescreensFS.ReadFile(paths[rand.IntN(len(paths))])
	if err != nil {
		log.Fatal(err)
	}

	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Definitiv kein Bluescreen")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Game{img: ebiten.NewImageFromImage(src)}); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) ||
		ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		((ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) &&
			ebiten.IsKeyPressed(ebiten.KeyC)) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	iw, ih := float64(g.img.Bounds().Dx()), float64(g.img.Bounds().Dy())

	// "contain": größter Faktor, bei dem das Bild noch komplett passt
	scale := sw / iw
	if s := sh / ih; s < scale {
		scale = s
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((sw-iw*scale)/2, (sh-ih*scale)/2) // zentrieren
	screen.DrawImage(g.img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
