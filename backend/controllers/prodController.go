// Seguimos dentro del paquete controllers, ya que maneja la lógica de otra sección del negocio.
package controllers

import (
	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos strconv para convertir strings a números.
	"strconv"
	// Importamos nuestro paquete de modelos para usar la estructura Producto.
	"backend/models"
	// Importamos el repository para acceder a la base de datos.
	"backend/repository"
)

// GetProductos es la función controladora encargada de devolver el catálogo de artículos.
func GetProductos(c *fiber.Ctx) error {
	// Obtener productos desde la base de datos usando el repository.
	productos, err := repository.ObtenerProductos(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener productos",
			Details: err.Error(),
		})
	}

	// Fiber convierte automáticamente el slice de estructuras a formato JSON y lo envía como respuesta al cliente.
	return c.JSON(productos)
}

// GetProductoPorID es la función controladora encargada de devolver un producto específico por su ID.
func GetProductoPorID(c *fiber.Ctx) error {
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
	producto, err := repository.ObtenerProductoPorID(c.Context(), id)
	if err != nil {
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
