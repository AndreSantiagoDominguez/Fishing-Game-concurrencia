package game

import (
	"context"
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480

	// Centro y radio del lago
	LakeCenterX = 320
	LakeCenterY = 240
	LakeRadius  = 180
)

var (
	// color de fallback si no hay imagen
	colorLake = color.RGBA{40, 140, 200, 255}
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateFishing
	StateCaught
)

// Game implementa ebiten.Game interface
type Game struct {
	// Sincronización
	mu     sync.Mutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	// Estado del juego
	state      GameState
	score      int
	fishCaught int

	// Contadores por rareza
	commonCount    int
	rareCount      int
	epicCount      int
	legendaryCount int

	// Entidades
	player *Player
	bobber *Bobber
	fishes []*Fish

	// Canales para concurrencia (Patrón Productor-Consumidor)
	spawnChan chan *Fish
	catchChan chan FishType

	// Assets
	lakeScene *ebiten.Image

	// Control de tiempo
	frameCount int
}

// NewGame crea una nueva instancia del juego
func NewGame() (*Game, error) {
	// Inicializar random seed
	rand.Seed(time.Now().UnixNano())

	// Crear contexto para cancelación
	ctx, cancel := context.WithCancel(context.Background())

	g := &Game{
		state:     StatePlaying,
		ctx:       ctx,
		cancel:    cancel,
		fishes:    make([]*Fish, 0),
		spawnChan: make(chan *Fish, 10),
		catchChan: make(chan FishType, 10),
	}

	// Inicializar jugador (fuera del lago)
	g.player = NewPlayer(float64(LakeCenterX), float64(LakeCenterY+LakeRadius+40))
	if err := g.player.LoadSprites(); err != nil {
		return nil, fmt.Errorf("error loading player sprites: %w", err)
	}

	// Inicializar bobber
	g.bobber = NewBobber()
	if err := g.bobber.LoadSprites(); err != nil {
		return nil, fmt.Errorf("error loading bobber sprites: %w", err)
	}

	// Cargar assets
	if err := g.loadAssets(); err != nil {
		return nil, fmt.Errorf("error loading assets: %w", err)
	}

	// Iniciar goroutines del patrón Productor-Consumidor
	g.wg.Add(1)
	go g.fishSpawner() // PRODUCTOR (en spawner.go)

	g.wg.Add(1)
	go g.catchProcessor() // CONSUMIDOR (en spawner.go)

	return g, nil
}

// loadAssets carga todas las imágenes necesarias
func (g *Game) loadAssets() error {
	var err error

	// Intentar cargar escenario del lago
	g.lakeScene, _, err = ebitenutil.NewImageFromFile("assets/lake_scene.png")
	if err != nil {
		g.lakeScene = nil
		fmt.Println("Warning: failed to load lake_scene.png, using color background:", err)
	}

	// Cargar sprites de peces (globales, compartidos)
	if err := LoadFishSprites(); err != nil {
		return fmt.Errorf("failed to load fish sprites: %w", err)
	}

	return nil
}

// Update actualiza la lógica del juego (60 FPS)
func (g *Game) Update() error {
	g.frameCount++

	// Procesar nuevos peces del canal (CONSUMIDOR)
	select {
	case fish := <-g.spawnChan:
		g.mu.Lock()
		g.fishes = append(g.fishes, fish)
		g.mu.Unlock()

		// Iniciar goroutine para el movimiento del pez
		g.wg.Add(1)
		go fish.Swim(g.ctx, &g.wg)
	default:
		// No hay peces nuevos en el canal
	}

	// Manejar input del usuario
	g.handleInput()

	// Actualizar jugador (solo si está jugando, no en modo pesca)
	if g.state == StatePlaying {
		g.player.Update()
	}

	// Actualizar bobber (animación)
	g.bobber.Update()

	// Detectar colisiones con peces si el bobber está activo
	if g.bobber.active {
		g.checkFishCollisions()
	}

	// Limpiar peces que salieron del lago
	g.cleanupFishes()

	return nil
}

// Draw dibuja el juego en la pantalla
func (g *Game) Draw(screen *ebiten.Image) {
	// Dibujar escenario
	if g.lakeScene != nil {
		screen.DrawImage(g.lakeScene, nil)
	} else {
		screen.Fill(colorLake)
	}

	// Dibujar peces (con efecto de sombra bajo el agua)
	g.mu.Lock()
	for _, fish := range g.fishes {
		fish.Draw(screen)
	}
	g.mu.Unlock()

	// Dibujar bobber (antes del jugador para que quede "en el agua")
	if g.bobber.active {
		g.bobber.Draw(screen)
	}

	// Dibujar jugador
	g.player.Draw(screen)

	// Dibujar UI (puntuación, estadísticas)
	g.drawUI(screen)
}

// drawUI dibuja la interfaz de usuario
func (g *Game) drawUI(screen *ebiten.Image) {
	// Fondo semi-transparente
	uiRect := ebiten.NewImage(240, 160)
	uiRect.Fill(color.RGBA{0, 0, 0, 160})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	screen.DrawImage(uiRect, op)

	// Obtener datos con mutex
	g.mu.Lock()
	scoreText := fmt.Sprintf("Puntos: %d", g.score)
	totalText := fmt.Sprintf("Total Capturados: %d", g.fishCaught)
	commonText := fmt.Sprintf("Comunes: %d", g.commonCount)
	rareText := fmt.Sprintf("Raros: %d", g.rareCount)
	epicText := fmt.Sprintf("Épicos: %d", g.epicCount)
	legendText := fmt.Sprintf("Legendarios: %d", g.legendaryCount)
	
	// Contar peces en el lago por tipo
	commonInLake := g.countFishType(FishCommon)
	rareInLake := g.countFishType(FishRare)
	epicInLake := g.countFishType(FishEpic)
	legendInLake := g.countFishType(FishLegendary)
	g.mu.Unlock()

	// Mostrar estadísticas
	ebitenutil.DebugPrintAt(screen, scoreText, 20, 20)
	ebitenutil.DebugPrintAt(screen, totalText, 20, 36)
	ebitenutil.DebugPrintAt(screen, commonText, 20, 56)
	ebitenutil.DebugPrintAt(screen, rareText, 20, 76)
	ebitenutil.DebugPrintAt(screen, epicText, 20, 96)
	ebitenutil.DebugPrintAt(screen, legendText, 20, 116)
	
	// Mostrar peces en el lago
	lakeInfo := fmt.Sprintf("En el Lago: %d C, %d R, %d E, %d L", 
		commonInLake, rareInLake, epicInLake, legendInLake)
	ebitenutil.DebugPrintAt(screen, lakeInfo, 20, 136)

	// Controles
	ebitenutil.DebugPrintAt(screen, "WASD: Mover | ESPACIO: Lanzar | R: Recoger", 10, ScreenHeight-20)
}

// handleInput maneja la entrada del usuario
func (g *Game) handleInput() {
	// Lanzar anzuelo con ESPACIO
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.state == StatePlaying {
		if g.player.IsNearWater() {
			g.state = StateFishing
			g.player.Cast()
			g.bobber.Cast(g.player.X, g.player.Y)
		}
	}

	// Recoger anzuelo con R
	if ebiten.IsKeyPressed(ebiten.KeyR) && g.state == StateFishing {
		g.state = StatePlaying
		g.bobber.Reset()
		g.player.StopFishing()
	}
}

// checkFishCollisions verifica si el anzuelo tocó algún pez
func (g *Game) checkFishCollisions() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i, fish := range g.fishes {
		if fish.CheckCollision(g.bobber.X, g.bobber.Y, 15) {
			// ¡Pez capturado!
			
			// IMPORTANTE: Desactivar bobber INMEDIATAMENTE para evitar múltiples capturas
			g.bobber.active = false
			g.bobber.SetState(BobberCaught)
			
			// Enviar al canal para que catchProcessor lo procese
			select {
			case g.catchChan <- fish.FishType:
			default:
			}

			// Remover pez de la lista
			g.fishes = append(g.fishes[:i], g.fishes[i+1:]...)
			fish.Stop()

			// Iniciar goroutine para resetear después de captura
			g.wg.Add(1)
			go g.resetAfterCatch()

			break // Solo capturar UN pez
		}
	}
}

// resetAfterCatch vuelve al modo de juego normal después de capturar un pez
func (g *Game) resetAfterCatch() {
	defer g.wg.Done()
	
	time.Sleep(1 * time.Second)
	
	g.mu.Lock()
	g.state = StatePlaying
	g.bobber.Reset()
	g.player.StopFishing()
	g.mu.Unlock()
}

// cleanupFishes elimina peces que están muy lejos del lago
func (g *Game) cleanupFishes() {
	g.mu.Lock()
	defer g.mu.Unlock()

	validFishes := make([]*Fish, 0, len(g.fishes))
	for _, fish := range g.fishes {
		// Mantener solo peces dentro de un área razonable
		dx := fish.X - LakeCenterX
		dy := fish.Y - LakeCenterY
		distance := dx*dx + dy*dy

		if distance < (LakeRadius+100)*(LakeRadius+100) {
			validFishes = append(validFishes, fish)
		} else {
			fish.Stop()
		}
	}
	g.fishes = validFishes
}

// countFishType cuenta cuántos peces de un tipo hay en el lago
// IMPORTANTE: Esta función NO usa mutex, debe ser llamada dentro de un lock
func (g *Game) countFishType(fishType FishType) int {
	count := 0
	for _, fish := range g.fishes {
		if fish.FishType == fishType {
			count++
		}
	}
	return count
}

// Layout define el tamaño de la pantalla
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Cleanup limpia recursos al cerrar
func (g *Game) Cleanup() {
	g.cancel()           // Cancelar todas las goroutines
	close(g.spawnChan)   // Cerrar canales
	close(g.catchChan)
	g.wg.Wait()          // Esperar a que todas las goroutines terminen
}

// getPointsForFish retorna los puntos según el tipo de pez
// Esta función es usada por catchProcessor en spawner.go
func (g *Game) getPointsForFish(fishType FishType) int {
	switch fishType {
	case FishCommon:
		return 10
	case FishRare:
		return 25
	case FishEpic:
		return 50
	case FishLegendary:
		return 100
	default:
		return 0
	}
}