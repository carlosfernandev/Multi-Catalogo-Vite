// Package routes registra los endpoints de la API desacoplados de la implementación concreta.
package routes

import (
	"os"

	"backend/controllers"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes recibe la aplicación Fiber y el agregador Controllers con sus dependencias inyectadas.
func SetupRoutes(app *fiber.App, ctrls *controllers.Controllers) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "multicatalogo-super-secret-jwt-key"
	}

	// Health check con diagnóstico de DB y Redis
	app.Get("/api/health", ctrls.Health.Health)

	// Autenticación pública
	app.Post("/api/login", ctrls.Auth.Login)

	// Catálogo de productos (lectura pública optimizada con Cache-Aside)
	app.Get("/api/productos", ctrls.Producto.GetProductos)
	app.Get("/api/productos/:id", ctrls.Producto.GetProductoPorID)

	// Rutas protegidas mediante token JWT (U2-T4)
	app.Get("/api/red", middleware.RequireAuth(jwtSecret), ctrls.Red.ObtenerRed)
	app.Post("/api/checkout", middleware.RequireAuth(jwtSecret), ctrls.Checkout.Checkout)
}
