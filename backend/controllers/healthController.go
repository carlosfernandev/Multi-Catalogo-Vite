// Package controllers agrupa las funciones de controladoras HTTP.
package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthController verifica la salud de la API y sus dependencias (PostgreSQL y Redis).
type HealthController struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

// NewHealthController crea una instancia de HealthController con acceso al pool de base de datos y cliente Redis.
func NewHealthController(db *pgxpool.Pool, rdb *redis.Client) *HealthController {
	return &HealthController{db: db, rdb: rdb}
}

// Health realiza ping a PostgreSQL (503 si falla) y a Redis (estado "degraded" si cae la caché).
func (h *HealthController) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1. Verificación obligatoria de PostgreSQL
	if h.db == nil || h.db.Ping(ctx) != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":   "error",
			"database": "down",
			"message":  "Base de datos PostgreSQL no disponible",
		})
	}

	// 2. Verificación de Redis (degradación suave si cae)
	cacheStatus := "up"
	overallStatus := "ok"
	if h.rdb == nil || h.rdb.Ping(ctx).Err() != nil {
		cacheStatus = "down"
		overallStatus = "degraded"
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   overallStatus,
		"database": "up",
		"cache":    cacheStatus,
		"message":  "Servidor Go/Fiber operativo",
	})
}
