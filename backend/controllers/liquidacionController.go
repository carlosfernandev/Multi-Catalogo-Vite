// Paquete controllers para liquidaciones y comisiones históricas.
package controllers

import (
	"backend/models"
	"backend/repository"
	"github.com/gofiber/fiber/v2"
)

// LiquidacionController gestiona la consulta de comisiones liquidadas.
type LiquidacionController struct {
	repo repository.DBRepository
}

// NewLiquidacionController inicializa el controlador inyectando el repositorio de datos.
func NewLiquidacionController(repo repository.DBRepository) *LiquidacionController {
	return &LiquidacionController{repo: repo}
}

// ObtenerLiquidaciones devuelve el historial de comisiones de la red.
// Si el usuario es cliente, se filtra por su email; si es admin, ve el global.
func (h *LiquidacionController) ObtenerLiquidaciones(c *fiber.Ctx) error {
	userEmail, _ := c.Locals("usuario_email").(string)

	resumen, err := h.repo.ObtenerLiquidaciones(c.Context(), userEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Status:  fiber.StatusInternalServerError,
			Message: "Error al consultar liquidaciones",
			Details: err.Error(),
		})
	}

	return c.JSON(resumen)
}
