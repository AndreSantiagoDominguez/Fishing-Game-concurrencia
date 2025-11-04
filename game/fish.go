package game

import (
	"context"
	"image"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type FishType int

const (
	FishCommon FishType = iota
	FishRare
	FishEpic
	FishLegendary
)

// Sprites globales para peces (compartidos por todas las instancias)
var (
	fishSprites   map[FishType]*ebiten.Image
	fishSpritesMu sync.Mutex
)

type Fish struct {
	X, Y      float64
	vx, vy    float64 // Velocidad
	FishType  FishType
	
	// Animación
	frame      int
	frameCount int
	
	// Control de goroutine
	mu       sync.Mutex
	active   bool
}

// LoadFishSprites carga todos los sprites de peces
func LoadFishSprites() error {
	fishSprites = make(map[FishType]*ebiten.Image)
	
	var err error
	fishSprites[FishCommon], _, err = ebitenutil.NewImageFromFile("assets/fish_common.png")
	if err != nil {
		return err
	}
	
	fishSprites[FishRare], _, err = ebitenutil.NewImageFromFile("assets/fish_rare.png")
	if err != nil {
		return err
	}
	
	fishSprites[FishEpic], _, err = ebitenutil.NewImageFromFile("assets/fish_epic.png")
	if err != nil {
		return err
	}
	
	fishSprites[FishLegendary], _, err = ebitenutil.NewImageFromFile("assets/fish_legendary.png")
	if err != nil {
		return err
	}
	
	return nil
}

func NewFish(x, y float64, fishType FishType) *Fish {
	// Velocidad aleatoria
	angle := rand.Float64() * 2 * math.Pi
	speed := 0.5 + rand.Float64()*1.0
	
	return &Fish{
		X:        x,
		Y:        y,
		vx:       math.Cos(angle) * speed,
		vy:       math.Sin(angle) * speed,
		FishType: fishType,
		active:   true,
	}
}

// Swim es la goroutine que controla el movimiento del pez
// Cada pez tiene su propia goroutine
func (f *Fish) Swim(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()
	
	changeDirectionCounter := 0
	
	for {
		select {
		case <-ctx.Done():
			return
			
		case <-ticker.C:
			f.mu.Lock()
			if !f.active {
				f.mu.Unlock()
				return
			}
			
			// Actualizar posición
			f.X += f.vx
			f.Y += f.vy
			
			// Cambiar dirección aleatoriamente cada cierto tiempo
			changeDirectionCounter++
			if changeDirectionCounter > 120 { // Cada ~2 segundos
				changeDirectionCounter = 0
				if rand.Float64() < 0.3 { // 30% de probabilidad
					angle := rand.Float64() * 2 * math.Pi
					speed := 0.5 + rand.Float64()*1.0
					f.vx = math.Cos(angle) * speed
					f.vy = math.Sin(angle) * speed
				}
			}
			
			// Mantener dentro del lago (comportamiento de rebote)
			dx := f.X - LakeCenterX
			dy := f.Y - LakeCenterY
			distance := math.Sqrt(dx*dx + dy*dy)
			
			if distance > LakeRadius-20 {
				// Rebotar hacia el centro
				angle := math.Atan2(dy, dx)
				f.vx = -math.Cos(angle) * (0.5 + rand.Float64()*1.0)
				f.vy = -math.Sin(angle) * (0.5 + rand.Float64()*1.0)
			}
			
			// Actualizar frame de animación
			f.frameCount++
			if f.frameCount >= 15 {
				f.frameCount = 0
				f.frame = (f.frame + 1) % 2
			}
			
			f.mu.Unlock()
		}
	}
}

// Draw dibuja el pez con efecto de sombra (bajo el agua)
func (f *Fish) Draw(screen *ebiten.Image) {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	fishSpritesMu.Lock()
	sprite := fishSprites[f.FishType]
	fishSpritesMu.Unlock()
	
	if sprite == nil {
		return
	}
	
	// Calcular tamaño del frame según el tipo de pez
	frameWidth := 32  // Ajusta según tus sprites
	frameHeight := 24 // Ajusta según tus sprites
	
	if f.FishType == FishLegendary {
		frameWidth = 50
		frameHeight = 40
	} else if f.FishType == FishEpic {
		frameWidth = 40
		frameHeight = 32
	}
	
	sx := f.frame * frameWidth
	sy := 0
	
	op := &ebiten.DrawImageOptions{}
	
	// Efecto de sombra bajo el agua (semi-transparente y oscurecido)
	op.ColorScale.Scale(0.7, 0.7, 0.9, 0.7) // Darken y transparencia
	
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(f.X, f.Y)
	
	subImg := sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
	screen.DrawImage(subImg, op)
}

// CheckCollision verifica si el pez colisionó con un punto (anzuelo)
func (f *Fish) CheckCollision(x, y, radius float64) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	dx := f.X - x
	dy := f.Y - y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	return distance < radius
}

// Stop detiene la goroutine del pez
func (f *Fish) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.active = false
}