// Seguimos dentro del paquete controllers, ya que maneja la lógica de otra sección del negocio.
package controllers

import (
	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos strconv para convertir strings a números.
	"strconv"
	// Importamos nuestro paquete de modelos para usar la estructura Producto.
	"backend/models"
)

// GetProductos es la función controladora encargada de devolver el catálogo de artículos.
func GetProductos(c *fiber.Ctx) error {
	// Declaramos e inicializamos un 'slice' (arreglo dinámico) usando nuestro modelo models.Producto.
	productos := []models.Producto{
		// Agregamos el primer producto con sus respectivos valores para ID, Nombre, Precio e Img.
		{ID: 1, Nombre: "Serum Revitalizante", Precio: 45.00, Img: "https://picsum.photos/seed/serum/150"},
		// Agregamos el segundo producto a la lista.
		{ID: 2, Nombre: "Crema Hidratante Pro", Precio: 32.50, Img: "https://picsum.photos/seed/crema/150"},
		// Agregamos el tercer producto a la lista.
		{ID: 3, Nombre: "Tónico Purificante", Precio: 28.00, Img: "https://picsum.photos/seed/tonico/150"},
		// Agregamos el cuarto producto a la lista.
		{ID: 4, Nombre: "Mascarilla Nocturna", Precio: 50.00, Img: "https://picsum.photos/seed/mascarilla/150"},
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

	// Buscamos el producto en el catálogo hardcodeado.
	productos := []models.Producto{
		{ID: 1, Nombre: "Serum Revitalizante", Precio: 45.00, Img: "https://picsum.photos/seed/serum/150"},
		{ID: 2, Nombre: "Crema Hidratante Pro", Precio: 32.50, Img: "https://picsum.photos/seed/crema/150"},
		{ID: 3, Nombre: "Tónico Purificante", Precio: 28.00, Img: "https://picsum.photos/seed/tonico/150"},
		{ID: 4, Nombre: "Mascarilla Nocturna", Precio: 50.00, Img: "https://picsum.photos/seed/mascarilla/150"},
	}

	// Iteramos sobre los productos buscando el ID solicitado.
	for _, producto := range productos {
		if producto.ID == id {
			return c.JSON(producto)
		}
	}

	// Si no encontramos el producto, devolvemos un 404.
	return c.Status(404).JSON(models.APIError{
		Status:  404,
		Message: "Producto no encontrado",
		Details: map[string]int{"id": id},
	})
}
