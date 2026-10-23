// Declaramos el paquete controllers para agrupar las funciones que manejan la lógica de negocio de las rutas.
package controllers

import (
	"errors"
	"os"
	"strconv"
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

	// Rechazar usuarios inactivos
	if !usuario.Activo {
		return c.Status(403).JSON(models.APIError{
			Status:  403,
			Message: "Cuenta desactivada. Contacte al administrador",
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

// Register maneja la creación de nuevas cuentas de clientes y su inclusión en la red MLM.
func (h *AuthController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de petición inválido",
		})
	}

	if req.Email == "" || req.Password == "" || req.Nombre == "" {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Nombre, Email y Password son campos obligatorios",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "La contraseña debe tener al menos 6 caracteres",
		})
	}

	// Hashear la contraseña con bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error procesando credenciales",
		})
	}
	req.Password = string(hashedBytes)

	nuevoUsuario, err := h.repo.RegistrarUsuario(c.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicado) {
			return c.Status(409).JSON(models.APIError{
				Status:  409,
				Message: "Ya existe una cuenta registrada con este correo electrónico",
			})
		}
		var errVal *repository.ErrValidacion
		if errors.As(err, &errVal) {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: errVal.Mensaje,
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al registrar usuario",
			Details: detallesSeguros(err),
		})
	}

	// Generar token JWT automático para el nuevo usuario
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "multicatalogo-super-secret-jwt-key"
	}

	claims := jwt.MapClaims{
		"sub":   nuevoUsuario.ID,
		"email": nuevoUsuario.Email,
		"rol":   nuevoUsuario.Rol,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(jwtSecret))

	return c.Status(201).JSON(fiber.Map{
		"token":       signedToken,
		"email":       nuevoUsuario.Email,
		"nombre":      nuevoUsuario.Nombre,
		"rol":         nuevoUsuario.Rol,
		"referido_id": nuevoUsuario.ReferidoID,
	})
}

// GetUsuariosAdmin devuelve la lista de todos los usuarios registrados para el panel de administración.
func (h *AuthController) GetUsuariosAdmin(c *fiber.Ctx) error {
	q := c.Query("q")
	usuarios, err := h.repo.ObtenerUsuariosAdmin(c.Context(), q)
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener usuarios para administración",
			Details: detallesSeguros(err),
		})
	}
	return c.JSON(usuarios)
}

// ActualizarUsuarioAdmin permite a un administrador cambiar el rol o activar/desactivar un usuario.
func (h *AuthController) ActualizarUsuarioAdmin(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "ID de usuario inválido",
		})
	}

	var req models.ActualizarUsuarioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de petición inválido",
			Details: err.Error(),
		})
	}

	currentUserID, _ := c.Locals("user_id").(int)
	if id == currentUserID {
		if req.Activo != nil && !*req.Activo {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: "No puedes desactivar tu propia cuenta de administrador",
			})
		}
		if req.Rol != nil && *req.Rol != "admin" {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: "No puedes remover tu propio rol de administrador",
			})
		}
	}

	actualizado, err := h.repo.ActualizarUsuarioAdmin(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{
				Status:  404,
				Message: "Usuario no encontrado",
			})
		}
		var errVal *repository.ErrValidacion
		if errors.As(err, &errVal) {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: errVal.Mensaje,
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al actualizar usuario",
			Details: detallesSeguros(err),
		})
	}

	return c.JSON(actualizado)
}
