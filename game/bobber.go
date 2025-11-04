package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BobberState int

const (
	BobberFloating BobberState = iota
	BobberBite
	BobberCaught
)

type Bobber struct {
	X, Y     float64
	active   bool
	state    BobberState
	bobCount int
	sprite   *ebiten.Image
}

func NewBobber() *Bobber {
	return &Bobber{
		active: false,
		state:  BobberFloating,
	}
}

// LoadSprites carga el sprite del bobber
func (b *Bobber) LoadSprites() error {
	var err error
	b.sprite, _, err = ebitenutil.NewImageFromFile("assets/bobber.png")
	return err
}

// Cast lanza el anzuelo desde la posición del jugador
func (b *Bobber) Cast(playerX, playerY float64) {
	// Calcular posición en el agua (hacia el centro del lago)
	dx := float64(LakeCenterX) - playerX
	dy := float64(LakeCenterY) - playerY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Evitar división por cero
	if distance == 0 {
		distance = 1
	}

	// Normalizar y lanzar a cierta distancia
	castDistance := 70.0
	b.X = playerX + (dx/distance)*castDistance
	b.Y = playerY + (dy/distance)*castDistance

	b.active = true
	b.state = BobberFloating
	b.bobCount = 0
}

// Update actualiza el bobber (animación de flotar)
func (b *Bobber) Update() {
	if !b.active {
		return
	}
	b.bobCount++
}

// Draw dibuja el bobber
func (b *Bobber) Draw(screen *ebiten.Image) {
	if !b.active {
		return
	}

	if b.sprite == nil {
		return
	}

	// Seleccionar frame según el estado
	frameWidth := 32
	frameHeight := 32

	sx := int(b.state) * frameWidth
	sy := 0

	op := &ebiten.DrawImageOptions{}

	// Efecto de bobbing (movimiento vertical)
	bobOffset := math.Sin(float64(b.bobCount)*0.12) * 2.5

	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(b.X, b.Y+bobOffset)

	subImg := b.sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
	screen.DrawImage(subImg, op)

	// Dibujar línea desde el bobber hacia arriba (simulando la línea de pesca)
	drawFishingLine(screen, b.X, b.Y+bobOffset-20)
}

// drawFishingLine dibuja una línea simple hacia arriba
func drawFishingLine(screen *ebiten.Image, x, y float64) {
	lineImg := ebiten.NewImage(2, 30)
	lineImg.Fill(image.White.C)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(0.5, 0.5, 0.5, 0.8) // Gris semi-transparente
	op.GeoM.Translate(x-1, y-30)

	screen.DrawImage(lineImg, op)
}

func (b *Bobber) SetState(state BobberState) {
	b.state = state
}

func (b *Bobber) Reset() {
	b.active = false
	b.state = BobberFloating
	b.bobCount = 0
}
