// Declaramos el paquete models para agrupar las estructuras de datos de la aplicación.
package models

import "time"

// LoginRequest define la estructura esperada para el cuerpo de la petición (JSON) al iniciar sesión.
type LoginRequest struct {
	// Email representa el correo del usuario; la etiqueta `json:"email"` indica cómo se mapea el JSON a esta variable.
	Email    string `json:"email"`
	// Password representa la contraseña del usuario; se extrae del campo "password" del JSON entrante.
	Password string `json:"password"`
}

// RegisterRequest define la estructura para crear una nueva cuenta de cliente.
type RegisterRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Ref      string `json:"ref,omitempty"` // Código de referido opcional
}

// Usuario representa un usuario del sistema en la base de datos.
type Usuario struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Nombre     string    `json:"nombre,omitempty"`
	Password   string    `json:"-"` // No se exporta en JSON
	Rol        string    `json:"rol"`
	Activo     bool      `json:"activo"`
	ReferidoID *int      `json:"referido_id,omitempty"`
	CreadoEn   time.Time `json:"creado_en,omitempty"`
}

// ActualizarUsuarioRequest define los campos modificables por un administrador.
type ActualizarUsuarioRequest struct {
	Nombre *string `json:"nombre,omitempty"`
	Rol    *string `json:"rol,omitempty"`
	Activo *bool   `json:"activo,omitempty"`
}

// Producto define la estructura de los datos de un artículo en nuestro catálogo.
type Producto struct {
	// ID es el identificador único numérico del producto.
	ID          int      `json:"id"`
	// Nombre es la descripción en texto del producto.
	Nombre      string   `json:"nombre"`
	// Descripcion es el texto detallado del producto.
	Descripcion string   `json:"descripcion"`
	// Precio es el costo del producto, almacenado como un número con decimales (float64).
	Precio      float64  `json:"precio"`
	// Categoria es la clasificación del producto.
	Categoria   string   `json:"categoria"`
	// Img es la URL o ruta que apunta a la fotografía o imagen del producto.
	Img         string   `json:"img"`
	// Galeria es un arreglo de URLs de imágenes adicionales.
	Galeria     []string `json:"galeria"`
	// Stock disponible en inventario.
	Stock       int      `json:"stock"`
	// Activo determina si el producto está visible en el catálogo público (borrado lógico).
	Activo      bool     `json:"activo"`
}

// CrearProductoRequest define el payload para dar de alta un nuevo producto.
type CrearProductoRequest struct {
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Precio      float64  `json:"precio"`
	Categoria   string   `json:"categoria"`
	Img         string   `json:"img"`
	Galeria     []string `json:"galeria"`
	Stock       int      `json:"stock"`
}

// ActualizarProductoRequest define el payload para modificar un producto existente.
type ActualizarProductoRequest struct {
	Nombre      *string   `json:"nombre,omitempty"`
	Descripcion *string   `json:"descripcion,omitempty"`
	Precio      *float64  `json:"precio,omitempty"`
	Categoria   *string   `json:"categoria,omitempty"`
	Img         *string   `json:"img,omitempty"`
	Galeria     *[]string `json:"galeria,omitempty"`
	Stock       *int      `json:"stock,omitempty"`
	Activo      *bool     `json:"activo,omitempty"`
}

// Referido representa un nodo en la red multinivel de referidos.
type Referido struct {
	ID       int        `json:"id"`
	Nombre   string     `json:"nombre"`
	Codigo   string     `json:"codigo,omitempty"`
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

// ClienteCheckout son los datos de envío que llena el cliente en Checkout.tsx.
type ClienteCheckout struct {
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	Direccion string `json:"direccion"`
	Ciudad    string `json:"ciudad"`
	Telefono  string `json:"telefono"`
}

// CheckoutItem es un producto del carrito. El precio NO se recibe del cliente: se lee de PostgreSQL.
type CheckoutItem struct {
	ProductoID int `json:"producto_id"`
	Cantidad   int `json:"cantidad"`
}

// CheckoutRequest es el cuerpo JSON de POST /api/checkout.
type CheckoutRequest struct {
	// UsuarioEmail es el correo de la cuenta autenticada que realiza la compra.
	UsuarioEmail string          `json:"usuario_email"`
	Cliente      ClienteCheckout `json:"cliente"`
	MetodoPago   string          `json:"metodo_pago"`
	Items        []CheckoutItem  `json:"items"`
}

// PedidoItem es una línea del detalle de un pedido confirmado.
type PedidoItem struct {
	ProductoID     int     `json:"producto_id"`
	Nombre         string  `json:"nombre"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

// Pedido es la respuesta de un checkout confirmado (tras el COMMIT).
type Pedido struct {
	ID        int          `json:"pedido_id"`
	Numero    string       `json:"numero"`
	UsuarioID int          `json:"-"`
	Total     float64      `json:"total"`
	Estado    string       `json:"estado"`
	Fecha     time.Time    `json:"fecha"`
	Items     []PedidoItem `json:"items"`
}

// Comision es el monto acreditado a un referido ascendente por un pedido.
type Comision struct {
	ReferidoID int     `json:"referido_id"`
	Nombre     string  `json:"nombre"`
	Nivel      int     `json:"nivel"`
	Monto      float64 `json:"monto"`
}

// LiquidacionRegistro representa una comisión histórica registrada en la BD.
type LiquidacionRegistro struct {
	ID             int       `json:"id"`
	PedidoID       int       `json:"pedido_id"`
	PedidoNumero   string    `json:"pedido_numero"`
	ReferidoID     int       `json:"referido_id"`
	ReferidoNombre string    `json:"referido_nombre"`
	Nivel          int       `json:"nivel"`
	Monto          float64   `json:"monto"`
	Fecha          time.Time `json:"fecha"`
}

// ResumenLiquidaciones agrupa el total liquidado y los registros detallados.
type ResumenLiquidaciones struct {
	TotalComisiones float64               `json:"total_comisiones"`
	Cantidad        int                   `json:"cantidad"`
	Liquidaciones   []LiquidacionRegistro `json:"liquidaciones"`
}

// LiquidarRequest payload para liquidar las comisiones pendientes de un distribuidor.
type LiquidarRequest struct {
	ReferidoID int `json:"referido_id"`
}

// Liquidacion representa un pago emitido a un distribuidor en la tabla liquidaciones.
type Liquidacion struct {
	ID             int        `json:"id"`
	ReferidoID     int        `json:"referido_id"`
	ReferidoNombre string     `json:"referido_nombre,omitempty"`
	Total          float64    `json:"total"`
	Estado         string     `json:"estado"`
	CreadoEn       time.Time  `json:"creado_en"`
	PagadaEn       *time.Time `json:"pagada_en,omitempty"`
}

// DistribuidorPendiente representa el saldo acumulado pendiente de cobro de un distribuidor.
type DistribuidorPendiente struct {
	ReferidoID         int     `json:"referido_id"`
	ReferidoNombre     string  `json:"referido_nombre"`
	TotalPendiente     float64 `json:"total_pendiente"`
	CantidadComisiones int     `json:"cantidad_comisiones"`
}
