// Package controllers agrupa y ensambla todos los controladores del microservicio.
package controllers

import (
	"backend/cache"
	"backend/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Controllers es el agregador que contiene las instancias de todos los controladores
// de la aplicación, facilitando la inyección de dependencias en un único punto.
type Controllers struct {
	Health   *HealthController
	Auth     *AuthController
	Producto *ProductoController
	Red      *RedController
	Checkout *CheckoutController
}

// NewControllers ensambla e inicializa cada controlador inyectándole sus dependencias necesarias.
func NewControllers(
	repo repository.DBRepository,
	cacheRepo cache.CacheRepository,
	db *pgxpool.Pool,
	rdb *redis.Client,
) *Controllers {
	return &Controllers{
		Health:   NewHealthController(db, rdb),
		Auth:     NewAuthController(repo),
		Producto: NewProductoController(repo, cacheRepo),
		Red:      NewRedController(repo, cacheRepo),
		Checkout: NewCheckoutController(repo, cacheRepo),
	}
}
