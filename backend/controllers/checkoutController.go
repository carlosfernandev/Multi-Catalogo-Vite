// Paquete controllers: procesamiento de la compra enviada desde Checkout.tsx.
package controllers

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"backend/cache"
	"backend/models"
	"backend/repository"
	"github.com/gofiber/fiber/v2"
)

// maxCantidadPorItem limita cantidades absurdas en una sola línea del carrito.
const maxCantidadPorItem = 100

// timeoutComisiones limita cuánto puede tardar la Goroutine de comisiones.
const timeoutComisiones = 30 * time.Second

// CheckoutController procesa compras con sus dependencias inyectadas.
type CheckoutController struct {
	repo  repository.DBRepository
	cache cache.CacheRepository
}

// NewCheckoutController crea el controlador inyectando repositorio y caché.
func NewCheckoutController(repo repository.DBRepository, cache cache.CacheRepository) *CheckoutController {
	return &CheckoutController{repo: repo, cache: cache}
}

// Checkout recibe el carrito, lo guarda en una transacción ACID y responde 201 con el pedido.
func (h *CheckoutController) Checkout(c *fiber.Ctx) error {
	var req models.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Cuerpo de petición inválido"})
	}

	if msg := validarCheckout(&req); msg != "" {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: msg})
	}

	// El comprador debe ser una cuenta existente.
	usuario, err := h.repo.BuscarUsuarioPorEmail(c.Context(), req.UsuarioEmail)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al buscar usuario", Details: err.Error()})
	}
	if usuario == nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Usuario no registrado", Details: req.UsuarioEmail})
	}

	pedido, err := h.repo.CrearPedido(c.Context(), usuario.ID, req)
	if err != nil {
		var errValidacion *repository.ErrValidacion
		if errors.As(err, &errValidacion) {
			return c.Status(400).JSON(models.APIError{Status: 400, Message: errValidacion.Mensaje})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "No se pudo procesar el pedido", Details: err.Error()})
	}

	// COMMIT exitoso: las comisiones se calculan en segundo plano para no hacer esperar al cliente.
	go h.calcularComisiones(pedido.ID, pedido.Total, usuario.ID, usuario.Email)

	log.Printf("Checkout: pedido %s confirmado (total %.2f), respondiendo 201", pedido.Numero, pedido.Total)
	return c.Status(201).JSON(pedido)
}

// calcularComisiones corre en una Goroutine: no usa el contexto de Fiber (se recicla al
// terminar la petición) y recupera cualquier panic para no tumbar el servidor.
func (h *CheckoutController) calcularComisiones(pedidoID int, total float64, usuarioID int, email string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error: panic en Goroutine de comisiones del pedido %d: %v", pedidoID, r)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeoutComisiones)
	defer cancel()

	inicio := time.Now()
	comisiones, err := h.repo.AplicarComisiones(ctx, pedidoID, usuarioID, total)
	if err != nil {
		log.Printf("Error: comisiones del pedido %d no aplicadas: %v", pedidoID, err)
		return
	}
	for _, com := range comisiones {
		log.Printf("Goroutine: pedido %d -> %s (nivel %d) +%.2f", pedidoID, com.Nombre, com.Nivel, com.Monto)
	}
	log.Printf("Goroutine: %d comisiones del pedido %d aplicadas en %v", len(comisiones), pedidoID, time.Since(inicio))

	if len(comisiones) == 0 {
		return
	}
	// Paso 5: invalidar la caché del árbol para que el Dashboard muestre los KPIs nuevos.
	// La red es la misma para todos los usuarios, así que además de red:tree:{email}
	// se borran las demás copias cacheadas (red:tree:*), incluida la del admin.
	h.cache.Eliminar(ctx, "red:tree:"+email)
	h.cache.EliminarPorPatron(ctx, "red:tree:*")
	log.Printf("Goroutine: caché red:tree:* invalidada tras el pedido %d", pedidoID)
}

// validarCheckout normaliza el cuerpo y retorna un mensaje de error (vacío si es válido).
func validarCheckout(req *models.CheckoutRequest) string {
	req.UsuarioEmail = strings.TrimSpace(req.UsuarioEmail)
	req.MetodoPago = strings.TrimSpace(req.MetodoPago)
	cli := &req.Cliente
	cli.Nombre = strings.TrimSpace(cli.Nombre)
	cli.Email = strings.TrimSpace(cli.Email)
	cli.Direccion = strings.TrimSpace(cli.Direccion)
	cli.Ciudad = strings.TrimSpace(cli.Ciudad)
	cli.Telefono = strings.TrimSpace(cli.Telefono)

	switch {
	case req.UsuarioEmail == "":
		return "usuario_email es obligatorio"
	case cli.Nombre == "" || cli.Email == "" || cli.Direccion == "" || cli.Ciudad == "" || cli.Telefono == "":
		return "Todos los datos del cliente son obligatorios"
	case req.MetodoPago == "":
		return "metodo_pago es obligatorio"
	case len(req.Items) == 0:
		return "El carrito está vacío"
	}

	for _, item := range req.Items {
		if item.ProductoID <= 0 || item.Cantidad <= 0 || item.Cantidad > maxCantidadPorItem {
			return "Cada item debe tener producto_id válido y cantidad entre 1 y 100"
		}
	}
	return ""
}
