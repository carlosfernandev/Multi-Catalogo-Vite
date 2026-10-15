// Paquete controllers para la red multinivel de referidos.
package controllers

import (
	"time"

	"backend/cache"
	"backend/models"
	"backend/repository"
	"github.com/gofiber/fiber/v2"
)

// ttlRed es el tiempo de expiración del árbol de referidos en caché.
const ttlRed = 10 * time.Minute

// ObtenerRed devuelve la jerarquía completa de la red multinivel.
// Aplica Cache-Aside sobre la clave red:tree:{email_usuario}; el email se recibe
// en el query param "email" (opcional, por defecto "anonimo").
func ObtenerRed(c *fiber.Ctx) error {
	email := c.Query("email", "anonimo")
	clave := "red:tree:" + email

	// 1. Intentar leer el árbol desde Redis (Cache Hit: sin tocar PostgreSQL).
	var arbol models.Referido
	if cache.Obtener(c.Context(), clave, &arbol) {
		c.Set("X-Cache", "HIT")
		return c.JSON(arbol)
	}

	// 2. Cache Miss: obtener lista plana de referidos desde la base de datos.
	referidos, err := repository.ObtenerRed(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener red de referidos",
			Details: err.Error(),
		})
	}

	// Construir árbol jerárquico en memoria y guardarlo en Redis con TTL de 10 minutos.
	raiz := repository.ConstruirArbolRed(referidos)
	cache.Guardar(c.Context(), clave, raiz, ttlRed)
	c.Set("X-Cache", "MISS")

	return c.JSON(raiz)
}
