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

// GetProductosAdmin devuelve todos los productos (activos e inactivos) para gestión administrativa.
func (h *ProductoController) GetProductosAdmin(c *fiber.Ctx) error {
	q := c.Query("q")
	productos, err := h.repo.ObtenerProductosAdmin(c.Context(), q)
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener lista administrativa de productos",
			Details: err.Error(),
		})
	}
	return c.JSON(productos)
}

// CrearProducto valida y da de alta un nuevo producto en el catálogo.
func (h *ProductoController) CrearProducto(c *fiber.Ctx) error {
	var req models.CrearProductoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de solicitud inválido",
			Details: err.Error(),
		})
	}

	if req.Nombre == "" {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El nombre del producto es obligatorio",
		})
	}
	if req.Precio <= 0 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El precio debe ser mayor a 0",
		})
	}
	if req.Categoria == "" {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "La categoría es obligatoria",
		})
	}
	if req.Stock < 0 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El stock no puede ser negativo",
		})
	}

	nuevo, err := h.repo.CrearProducto(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al registrar el producto",
			Details: err.Error(),
		})
	}

	// Invalidar clave de caché para que el catálogo público refleje el nuevo ítem
	h.cache.Eliminar(c.Context(), claveProductos)

	return c.Status(201).JSON(nuevo)
}

// ActualizarProducto permite la modificación de un producto existente.
func (h *ProductoController) ActualizarProducto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "ID de producto inválido",
		})
	}

	var req models.ActualizarProductoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de solicitud inválido",
			Details: err.Error(),
		})
	}

	if req.Precio != nil && *req.Precio <= 0 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El precio debe ser estrictamente mayor a 0",
		})
	}
	if req.Stock != nil && *req.Stock < 0 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El stock no puede ser negativo",
		})
	}

	actualizado, err := h.repo.ActualizarProducto(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{
				Status:  404,
				Message: "Producto no encontrado para actualizar",
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al actualizar producto",
			Details: err.Error(),
		})
	}

	// Invalidar caché tras la edición
	h.cache.Eliminar(c.Context(), claveProductos)

	return c.JSON(actualizado)
}

// EliminarProducto aplica borrado lógico desactivando el producto.
func (h *ProductoController) EliminarProducto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "ID de producto inválido",
		})
	}

	if err := h.repo.EliminarProductoLogico(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{
				Status:  404,
				Message: "Producto no encontrado",
			})
		}
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al desactivar producto",
			Details: err.Error(),
		})
	}

	// Invalidar caché tras la desactivación
	h.cache.Eliminar(c.Context(), claveProductos)

	return c.JSON(fiber.Map{
		"mensaje": "Producto desactivado exitosamente (borrado lógico)",
		"id":      id,
	})
}
