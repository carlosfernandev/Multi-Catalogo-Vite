// Seguimos dentro del paquete controllers, ya que maneja la lógica de otra sección del negocio.
package controllers

import (
	"errors"
	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos strconv para convertir strings a números.
	"strconv"
	// Importamos time para definir el TTL de la caché.
	"time"
	// Importamos el paquete cache para aplicar el patrón Cache-Aside con Redis.
	"backend/cache"
	// Importamos nuestro paquete de modelos para usar la estructura Producto.
	"backend/models"
	// Importamos el repository para acceder a la base de datos.
	"backend/repository"
)

// claveProductos es la clave de Redis donde se almacena el catálogo completo.
const claveProductos = "productos:all"

// ttlProductos es el tiempo de expiración del catálogo en caché.
const ttlProductos = 5 * time.Minute

// ProductoController agrupa los handlers del catálogo con sus dependencias inyectadas.
type ProductoController struct {
	repo  repository.DBRepository
	cache cache.CacheRepository
}

// NewProductoController crea el controlador inyectando repositorio y caché.
func NewProductoController(repo repository.DBRepository, cache cache.CacheRepository) *ProductoController {
	return &ProductoController{repo: repo, cache: cache}
}

// GetProductos es la función controladora encargada de devolver el catálogo de artículos.
func (h *ProductoController) GetProductos(c *fiber.Ctx) error {
	// 1. Intentar leer el catálogo desde Redis.
	var productos []models.Producto
	if h.cache.Obtener(c.Context(), claveProductos, &productos) {
		// 2. Cache Hit: responder sin tocar PostgreSQL.
		c.Set("X-Cache", "HIT")
		return c.JSON(productos)
	}

	// 3. Cache Miss: obtener productos desde la base de datos usando el repository.
	productos, err := h.repo.ObtenerProductos(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener productos",
			Details: err.Error(),
		})
	}

	// Guardar el resultado en Redis con TTL de 5 minutos.
	h.cache.Guardar(c.Context(), claveProductos, productos, ttlProductos)
	c.Set("X-Cache", "MISS")

	// Fiber convierte automáticamente el slice de estructuras a formato JSON y lo envía como respuesta al cliente.
	return c.JSON(productos)
}

// GetProductoPorID es la función controladora encargada de devolver un producto específico por su ID.
func (h *ProductoController) GetProductoPorID(c *fiber.Ctx) error {
	// Obtenemos el parámetro "id" de la URL como string.
	idParam := c.Params("id")

	// Convertimos el string a entero. Si falla, devolvemos un error 400.
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "ID inválido: debe ser un número",
			Details: map[string]string{"param": idParam},
		})
	}

	// Buscar producto en la base de datos usando el repository.
	producto, err := h.repo.ObtenerProductoPorID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{
				Status:  404,
				Message: "Producto no encontrado",
				Details: map[string]int{"id": id},
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al buscar producto",
			Details: err.Error(),
		})
	}

	// Si no se encontró el producto, devolvemos un 404.
	if producto == nil {
		return c.Status(404).JSON(models.APIError{
			Status:  404,
			Message: "Producto no encontrado",
			Details: map[string]int{"id": id},
		})
	}

	return c.JSON(producto)
}
