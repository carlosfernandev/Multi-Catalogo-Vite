// Paquete controllers para la red multinivel de referidos.
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"backend/models"
	"backend/repository"
)

// ObtenerRed devuelve la jerarquía completa de la red multinivel.
func ObtenerRed(c *fiber.Ctx) error {
	// Obtener lista plana de referidos desde la base de datos.
	referidos, err := repository.ObtenerRed(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener red de referidos",
			Details: err.Error(),
		})
	}

	// Construir árbol jerárquico en memoria.
	arbol := repository.ConstruirArbolRed(referidos)

	return c.JSON(arbol)
}
