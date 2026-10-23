// Paquete controllers para liquidaciones y comisiones históricas.
package controllers

import (
	"errors"

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

// ObtenerPendientes lista los distribuidores que tienen comisiones pendientes de pago.
func (h *LiquidacionController) ObtenerPendientes(c *fiber.Ctx) error {
	pendientes, err := h.repo.ObtenerPendientesLiquidacion(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener saldos pendientes de comisiones",
			Details: err.Error(),
		})
	}
	return c.JSON(pendientes)
}

// ObtenerHistorial devuelve la lista histórica de liquidaciones ejecutadas.
func (h *LiquidacionController) ObtenerHistorial(c *fiber.Ctx) error {
	historial, err := h.repo.ObtenerHistorialLiquidaciones(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener historial de liquidaciones",
			Details: err.Error(),
		})
	}
	return c.JSON(historial)
}

// Liquidar procesa el pago y cierre de comisiones pendientes para un distribuidor.
func (h *LiquidacionController) Liquidar(c *fiber.Ctx) error {
	var req models.LiquidarRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de petición inválido",
			Details: err.Error(),
		})
	}

	if req.ReferidoID <= 0 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "referido_id inválido o requerido",
		})
	}

	liq, err := h.repo.LiquidarComisiones(c.Context(), req.ReferidoID)
	if err != nil {
		var errVal *repository.ErrValidacion
		if errors.As(err, &errVal) {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: errVal.Mensaje,
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al procesar la liquidación",
			Details: err.Error(),
		})
	}

	return c.Status(201).JSON(liq)
}
