// Package middleware provee intermediarios de seguridad y autenticación para Fiber.
package middleware

import (
	"errors"
	"strings"

	"backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ClaimsJWT define la estructura de los datos contenidos en el token JWT.
type ClaimsJWT struct {
	Sub   int    `json:"sub"`
	Email string `json:"email"`
	Rol   string `json:"rol"`
	jwt.RegisteredClaims
}

// RequireAuth valida que la petición contenga un JWT válido en el header Authorization (Bearer).
// Si el token es válido, inyecta los claims en c.Locals para uso en controladores subsiguientes.
func RequireAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(models.APIError{
				Status:  401,
				Message: "Token no proporcionado",
			})
		}

		partes := strings.SplitN(authHeader, " ", 2)
		if len(partes) != 2 || strings.ToLower(partes[0]) != "bearer" {
			return c.Status(401).JSON(models.APIError{
				Status:  401,
				Message: "Formato de token inválido. Se espera: Bearer <token>",
			})
		}

		tokenString := strings.TrimSpace(partes[1])
		claims := &ClaimsJWT{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de firma no permitido")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(models.APIError{
				Status:  401,
				Message: "Token inválido o expirado",
			})
		}

		// Almacenar identidad del usuario autenticado en el contexto de Fiber
		c.Locals("user_id", claims.Sub)
		c.Locals("user_email", claims.Email)
		c.Locals("user_rol", claims.Rol)

		return c.Next()
	}
}

// RequireAdmin asegura que el usuario autenticado posea el rol de 'admin'.
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rol, ok := c.Locals("user_rol").(string)
		if !ok || rol != "admin" {
			return c.Status(403).JSON(models.APIError{
				Status:  403,
				Message: "Acceso denegado: se requiere rol de administrador",
			})
		}
		return c.Next()
	}
}
