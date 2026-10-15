// Declaramos el paquete controllers para agrupar las funciones de controladoras.
package controllers

import (
	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
)

// Health es el endpoint de health check que indica si el servidor está operativo.
func Health(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "Servidor Go/Fiber operativo",
	})
}
