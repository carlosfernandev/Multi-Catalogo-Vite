// Declaramos el paquete controllers para agrupar las funciones que manejan la lógica de negocio de las rutas.
package controllers

import (
	"errors"
	"os"
	"time"

	// Importamos el framework Fiber para tener acceso al contexto (c *fiber.Ctx) de la petición HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos JWT y bcrypt para seguridad robusta.
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

// detallesSeguros sanitiza el mensaje de error: solo se expone en modo desarrollo.
func detallesSeguros(err error) string {
	if os.Getenv("APP_ENV") == "dev" && err != nil {
		return err.Error()
	}
	return ""
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

	// Buscar usuario en la base de datos por email.
	usuario, err := h.repo.BuscarUsuarioPorEmail(c.Context(), req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(401).JSON(models.APIError{
				Status:  401,
				Message: "Credenciales incorrectas",
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al buscar usuario",
			Details: detallesSeguros(err),
		})
	}

	// Si no se encontró el usuario, retornar 401.
	if usuario == nil {
		return c.Status(401).JSON(models.APIError{
			Status:  401,
			Message: "Credenciales incorrectas",
		})
	}

	// Validación de contraseña mediante bcrypt con fallback a texto plano para compatibilidad de migración.
	errBcrypt := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(req.Password))
	if errBcrypt != nil && usuario.Password != req.Password {
		return c.Status(401).JSON(models.APIError{
			Status:  401,
			Message: "Credenciales incorrectas",
		})
	}

	// Obtener la clave secreta para la firma del token desde variables de entorno.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "multicatalogo-super-secret-jwt-key"
	}

	// Generar JWT real con claims estándar (sub, email, rol, iat, exp 2h) firmado con HS256.
	claims := jwt.MapClaims{
		"sub":   usuario.ID,
		"email": usuario.Email,
		"rol":   usuario.Rol,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al generar token de autenticación",
			Details: detallesSeguros(err),
		})
	}

	// Retornar respuesta exitosa con token JWT real, email y rol.
	return c.JSON(fiber.Map{
		"token": signedToken,
		"email": usuario.Email,
		"rol":   usuario.Rol,
	})
}
