// Declaramos el paquete routes, cuya única responsabilidad es registrar los endpoints de la API.
package routes

import (
	// Importamos Fiber para poder recibir y configurar la instancia principal de la aplicación.
	"github.com/gofiber/fiber/v2"
	// Importamos nuestro paquete de controladores, donde residen las funciones que se ejecutarán en cada ruta.
	"backend/controllers"
	// Importamos las interfaces de repositorio y caché que se inyectan en los controladores.
	"backend/cache"
	"backend/repository"
)

// SetupRoutes recibe la aplicación Fiber y las dependencias (repositorio y caché),
// instancia los controladores inyectándoselas y registra los endpoints.
func SetupRoutes(app *fiber.App, repo repository.DBRepository, cacheRepo cache.CacheRepository) {
	// Instanciamos los controladores con sus dependencias (inyección por constructor).
	authCtrl := controllers.NewAuthController(repo)
	prodCtrl := controllers.NewProductoController(repo, cacheRepo)
	redCtrl := controllers.NewRedController(repo, cacheRepo)

	// Definimos una ruta HTTP GET en "/api/health" para health check.
	app.Get("/api/health", controllers.Health)

	// Definimos una ruta HTTP POST en "/api/login" y la enlazamos a la función Login del paquete controllers.
	app.Post("/api/login", authCtrl.Login)

	// Definimos una ruta HTTP GET en "/api/productos" y la enlazamos a la función GetProductos del paquete controllers.
	app.Get("/api/productos", prodCtrl.GetProductos)

	// Definimos una ruta HTTP GET en "/api/productos/:id" para obtener un producto por ID.
	app.Get("/api/productos/:id", prodCtrl.GetProductoPorID)

	// Definimos una ruta HTTP GET en "/api/red" para obtener la red multinivel de referidos.
	app.Get("/api/red", redCtrl.ObtenerRed)
}
