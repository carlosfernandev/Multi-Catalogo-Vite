// Declaramos el paquete controllers para agrupar las funciones que manejan la lógica de negocio de las rutas.
package controllers

import (
	// Importamos el framework Fiber para tener acceso al contexto (c *fiber.Ctx) de la petición HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos nuestro paquete de modelos para poder usar la estructura LoginRequest.
	"backend/models"
	// Importamos el repository para acceder a la base de datos.
	"backend/repository"
)

// AuthController agrupa los handlers de autenticación con su repositorio inyectado.
type AuthController struct {
	repo repository.DBRepository
}

// NewAuthController crea el controlador inyectando el repositorio.
func NewAuthController(repo repository.DBRepository) *AuthController {
	return &AuthController{repo: repo}
}

// Login es la función controladora que se ejecutará cuando el cliente envíe sus credenciales.
func (h *AuthController) Login(c *fiber.Ctx) error {
	// Creamos una variable 'req' del tipo LoginRequest (ubicada en nuestro paquete models) para almacenar los datos.
	var req models.LoginRequest

	// Intentamos parsear (transformar) el cuerpo JSON entrante y guardarlo en la variable 'req'.
	if err := c.BodyParser(&req); err != nil {
		// Si ocurre un error al parsear (ej. JSON mal formado), retornamos un estado HTTP 400 (Bad Request).
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de petición inválido",
		})
	}

	// Validamos que los campos Email y Password no estén vacíos.
	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Email y Password son obligatorios",
		})
	}

	// Buscar usuario en la base de datos usando el repository.
	usuario, err := h.repo.BuscarUsuarioPorCredenciales(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al buscar usuario",
			Details: err.Error(),
		})
	}

	// Si no se encontró el usuario, retornar 401.
	if usuario == nil {
		return c.Status(401).JSON(models.APIError{
			Status:  401,
			Message: "Credenciales incorrectas",
		})
	}

	// Generar token ficticio (en producción usar JWT real).
	token := "fake-jwt-token-" + usuario.Rol

	// Retornar respuesta exitosa con token, email y rol.
	return c.JSON(fiber.Map{
		"token": token,
		"email": usuario.Email,
		"rol":   usuario.Rol,
	})
}
