package game

import (
	"math"
	"math/rand"
	"time"
)

// Límites de peces por tipo en el lago
const (
	MaxCommonFish    = 10
	MaxRareFish      = 6
	MaxEpicFish      = 4
	MaxLegendaryFish = 1
)

// ============================================================================
// PRODUCTOR: fishSpawner
// ============================================================================
// Esta goroutine genera nuevos peces periódicamente y los envía al canal
// Respeta los límites por tipo para evitar saturación
func (g *Game) fishSpawner() {
	defer g.wg.Done()

	ticker := time.NewTicker(3 * time.Second) // Intentar generar pez cada 3 segundos
	defer ticker.Stop()

	for {
		select {
		case <-g.ctx.Done():
			// El juego se está cerrando
			return

		case <-ticker.C:
			// Decidir qué tipo de pez crear (probabilidades normales)
			fishType := g.randomFishType()

			// Verificar si hay espacio para este tipo de pez
			if g.canSpawnFish(fishType) {
				// Crear el pez
				fish := g.spawnFishOfType(fishType)

				// Enviar al canal (no bloqueante)
				select {
				case g.spawnChan <- fish:
					// Pez enviado exitosamente al canal
				default:
					// Canal lleno, descartar este pez
				}
			}
			// Si no hay espacio, simplemente no se crea nada este tick
		}
	}
}

// canSpawnFish verifica si se puede crear un pez de este tipo
// Cuenta cuántos hay actualmente y compara con el límite
func (g *Game) canSpawnFish(fishType FishType) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Contar peces de este tipo en el lago
	count := g.countFishType(fishType)

	// Verificar contra el límite correspondiente
	switch fishType {
	case FishCommon:
		return count < MaxCommonFish
	case FishRare:
		return count < MaxRareFish
	case FishEpic:
		return count < MaxEpicFish
	case FishLegendary:
		return count < MaxLegendaryFish
	default:
		return false
	}
}

// spawnFishOfType crea un pez del tipo especificado en una posición aleatoria
func (g *Game) spawnFishOfType(fishType FishType) *Fish {
	// Generar posición aleatoria dentro del lago (círculo)
	angle := rand.Float64() * 2 * math.Pi
	radius := rand.Float64() * (LakeRadius - 20) // Un poco dentro del borde

	x := LakeCenterX + radius*math.Cos(angle)
	y := LakeCenterY + radius*math.Sin(angle)

	return NewFish(x, y, fishType)
}

// randomFishType determina el tipo de pez basado en probabilidades
// IMPORTANTE: Las probabilidades NUNCA cambian, son siempre las mismas
func (g *Game) randomFishType() FishType {
	roll := rand.Float64()

	// Probabilidades fijas:
	// Común:      60%
	// Raro:       25%
	// Épico:      12%
	// Legendario:  3%

	switch {
	case roll < 0.60:
		return FishCommon
	case roll < 0.85:
		return FishRare
	case roll < 0.97:
		return FishEpic
	default:
		return FishLegendary
	}
}

// ============================================================================
// CONSUMIDOR: catchProcessor
// ============================================================================
// Esta goroutine lee del canal de capturas y actualiza la puntuación
func (g *Game) catchProcessor() {
	defer g.wg.Done()

	for {
		select {
		case <-g.ctx.Done():
			return

		case fishType := <-g.catchChan:
			// Calcular puntos por el pez capturado
			points := g.getPointsForFish(fishType)

			// Actualizar estadísticas (con mutex para thread-safety)
			g.mu.Lock()
			g.score += points
			g.fishCaught++

			// Actualizar contador específico del tipo de pez
			switch fishType {
			case FishCommon:
				g.commonCount++
			case FishRare:
				g.rareCount++
			case FishEpic:
				g.epicCount++
			case FishLegendary:
				g.legendaryCount++
			}

			g.mu.Unlock()
		}
	}
}