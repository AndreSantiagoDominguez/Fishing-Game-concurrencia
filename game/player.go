package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Direction int

const (
	DirectionDown Direction = iota
	DirectionUp
	DirectionLeft
	DirectionRight
)

const (
	PlayerSpeed = 2.5
	FrameDelay  = 8 // frames entre cambios (ajusta si quieres más/menos rapidez)
)

type Player struct {
	X, Y          float64
	direction     Direction
	moving        bool
	frame         int
	frameCount    int
	walkUp        *ebiten.Image
	walkDown      *ebiten.Image
	walkLeft      *ebiten.Image
	walkRight     *ebiten.Image
	fishingSheet  *ebiten.Image
	fishingFrames []*ebiten.Image
	isFishing     bool
	fishFrame     int
	fishFCount    int
}

// NewPlayer crea un jugador en x,y
func NewPlayer(x, y float64) *Player {
	return &Player{
		X:         x,
		Y:         y,
		direction: DirectionDown,
	}
}

// LoadSprites carga todos los sprites del jugador
func (p *Player) LoadSprites() error {
	var err error
	p.walkUp, _, err = ebitenutil.NewImageFromFile("assets/fisherman_walk_up.png")
	if err != nil {
		return err
	}
	p.walkDown, _, err = ebitenutil.NewImageFromFile("assets/fisherman_walk_down.png")
	if err != nil {
		return err
	}
	p.walkLeft, _, err = ebitenutil.NewImageFromFile("assets/fisherman_walk_left.png")
	if err != nil {
		return err
	}
	p.walkRight, _, err = ebitenutil.NewImageFromFile("assets/fisherman_walk_right.png")
	if err != nil {
		return err
	}

	// Cargar y dividir sprite sheet de pesca (3 frames)
	p.fishingSheet, _, err = ebitenutil.NewImageFromFile("assets/fisherman_fishing.png")
	if err != nil {
		return err
	}
	p.loadFishingFrames()
	return nil
}

// loadFishingFrames divide la hoja de pesca en 3 frames
func (p *Player) loadFishingFrames() {
	if p.fishingSheet == nil {
		return
	}
	totalW := p.fishingSheet.Bounds().Dx()
	totalH := p.fishingSheet.Bounds().Dy()
	frameW := totalW / 3
	frameH := totalH

	for i := 0; i < 3; i++ {
		frame := p.fishingSheet.SubImage(image.Rect(i*frameW, 0, (i+1)*frameW, frameH)).(*ebiten.Image)
		p.fishingFrames = append(p.fishingFrames, frame)
	}
}

// Update procesa movimiento y animaciones
func (p *Player) Update() {
	oldX, oldY := p.X, p.Y
	p.moving = false

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= PlayerSpeed
		p.direction = DirectionUp
		p.moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += PlayerSpeed
		p.direction = DirectionDown
		p.moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= PlayerSpeed
		p.direction = DirectionLeft
		p.moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += PlayerSpeed
		p.direction = DirectionRight
		p.moving = true
	}

	if p.isInWater() {
		p.X = oldX
		p.Y = oldY
		p.moving = false
	}

	if p.X < 20 {
		p.X = 20
	}
	if p.X > ScreenWidth-20 {
		p.X = ScreenWidth - 20
	}
	if p.Y < 20 {
		p.Y = 20
	}
	if p.Y > ScreenHeight-20 {
		p.Y = ScreenHeight - 20
	}

	if p.moving {
		p.frameCount++
		if p.frameCount >= FrameDelay {
			p.frameCount = 0
			p.frame = (p.frame + 1) % 2
		}
	} else {
		p.frame = 0
	}

	if p.isFishing {
		p.fishFCount++
		if p.fishFCount >= FrameDelay {
			p.fishFCount = 0
			p.fishFrame = (p.fishFrame + 1) % len(p.fishingFrames)
		}
	}
}

// Draw dibuja al jugador
func (p *Player) Draw(screen *ebiten.Image) {
	if p.isFishing && len(p.fishingFrames) > 0 {
		frame := p.fishingFrames[p.fishFrame]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(frame.Bounds().Dx())/2, -float64(frame.Bounds().Dy())/2)
		op.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(frame, op)
		return
	}

	var sprite *ebiten.Image
	switch p.direction {
	case DirectionUp:
		sprite = p.walkUp
	case DirectionDown:
		sprite = p.walkDown
	case DirectionLeft:
		sprite = p.walkLeft
	case DirectionRight:
		sprite = p.walkRight
	}
	if sprite == nil {
		return
	}
	totalW := sprite.Bounds().Dx()
	totalH := sprite.Bounds().Dy()
	frameW := totalW / 2
	sx := p.frame * frameW
	sub := sprite.SubImage(image.Rect(sx, 0, sx+frameW, totalH)).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameW)/2, -float64(totalH)/2)
	op.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(sub, op)
}

func (p *Player) isInWater() bool {
	dx := p.X - LakeCenterX
	dy := p.Y - LakeCenterY
	return math.Sqrt(dx*dx+dy*dy) < LakeRadius-10
}

func (p *Player) IsNearWater() bool {
	dx := p.X - LakeCenterX
	dy := p.Y - LakeCenterY
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist >= LakeRadius-30 && dist <= LakeRadius+60
}

func (p *Player) Cast() {
	p.isFishing = true
	p.fishFrame = 0
	p.fishFCount = 0
}

func (p *Player) StopFishing() {
	p.isFishing = false
	p.fishFrame = 0
	p.fishFCount = 0
}
