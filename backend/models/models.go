// Declaramos el paquete models para agrupar las estructuras de datos de la aplicación.
package models

// LoginRequest define la estructura esperada para el cuerpo de la petición (JSON) al iniciar sesión.
type LoginRequest struct {
	// Email representa el correo del usuario; la etiqueta `json:"email"` indica cómo se mapea el JSON a esta variable.
	Email    string `json:"email"`
	// Password representa la contraseña del usuario; se extrae del campo "password" del JSON entrante.
	Password string `json:"password"`
}

// Usuario representa un usuario del sistema en la base de datos.
type Usuario struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // No se exporta en JSON
	Rol      string `json:"rol"`
}

// Producto define la estructura de los datos de un artículo en nuestro catálogo.
type Producto struct {
	// ID es el identificador único numérico del producto.
	ID         int      `json:"id"`
	// Nombre es la descripción en texto del producto.
	Nombre     string   `json:"nombre"`
	// Descripcion es el texto detallado del producto.
	Descripcion string  `json:"descripcion"`
	// Precio es el costo del producto, almacenado como un número con decimales (float64).
	Precio     float64  `json:"precio"`
	// Categoria es la clasificación del producto.
	Categoria  string   `json:"categoria"`
	// Img es la URL o ruta que apunta a la fotografía o imagen del producto.
	Img        string   `json:"img"`
	// Galeria es un arreglo de URLs de imágenes adicionales.
	Galeria    []string `json:"galeria"`
}

// Referido representa un nodo en la red multinivel de referidos.
type Referido struct {
	ID       int        `json:"id"`
	Nombre   string     `json:"nombre"`
	Nivel    int        `json:"nivel"`
	Ventas   float64    `json:"ventas"`
	ParentID *int       `json:"parent_id,omitempty"`
	Hijos    []Referido `json:"hijos,omitempty"` // Se construye en memoria desde parent_id
}

// APIError estandariza las respuestas de error HTTP en toda la API.
type APIError struct {
	// Status es el código de estado HTTP.
	Status int `json:"status"`
	// Message es el mensaje de error descriptivo.
	Message string `json:"message"`
	// Details es información adicional opcional sobre el error.
	Details interface{} `json:"details,omitempty"`
}