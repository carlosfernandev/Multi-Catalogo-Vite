// Paquete repository maneja el acceso a datos de la base de datos.
package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"backend/models"
)

// DBRepository es el contrato de acceso a datos que se inyecta en los controladores.
// Los controladores dependen de esta interfaz, no de PostgreSQL directamente.
type DBRepository interface {
	BuscarUsuarioPorCredenciales(ctx context.Context, email, password string) (*models.Usuario, error)
	ObtenerProductos(ctx context.Context) ([]models.Producto, error)
	ObtenerProductoPorID(ctx context.Context, id int) (*models.Producto, error)
	ObtenerRed(ctx context.Context) ([]models.Referido, error)
	BuscarUsuarioPorEmail(ctx context.Context, email string) (*models.Usuario, error)
	CrearPedido(ctx context.Context, usuarioID int, req models.CheckoutRequest) (*models.Pedido, error)
	AplicarComisiones(ctx context.Context, pedidoID, usuarioID int, total float64) ([]models.Comision, error)
}

// ErrValidacion indica un error de datos del cliente (HTTP 400), p. ej. producto inexistente o sin stock.
type ErrValidacion struct {
	Mensaje string
}

func (e *ErrValidacion) Error() string { return e.Mensaje }

// SimulacionFallas permite demostrar el ROLLBACK en clase sin editar el código.
type SimulacionFallas struct {
	// FallaAntesDeItems corta la transacción después de insertar el pedido y antes del detalle.
	FallaAntesDeItems bool
	// PausaAntesDeCommit espera antes del COMMIT (para detener el contenedor de PostgreSQL a mitad).
	PausaAntesDeCommit time.Duration
}

// SimulacionDesdeEntorno lee CHECKOUT_FALLA_SIMULADA=true y CHECKOUT_PAUSA_COMMIT=<segundos>.
func SimulacionDesdeEntorno() SimulacionFallas {
	sim := SimulacionFallas{FallaAntesDeItems: os.Getenv("CHECKOUT_FALLA_SIMULADA") == "true"}
	if seg, err := strconv.Atoi(os.Getenv("CHECKOUT_PAUSA_COMMIT")); err == nil && seg > 0 {
		sim.PausaAntesDeCommit = time.Duration(seg) * time.Second
	}
	if sim.FallaAntesDeItems || sim.PausaAntesDeCommit > 0 {
		log.Printf("Advertencia: simulación de fallas activa en checkout: %+v", sim)
	}
	return sim
}

// PostgresRepository implementa DBRepository usando un pool de pgx.
type PostgresRepository struct {
	db         *pgxpool.Pool
	simulacion SimulacionFallas
}

// NewPostgresRepository crea el repositorio inyectando el pool de conexiones y la simulación de fallas.
func NewPostgresRepository(db *pgxpool.Pool, simulacion SimulacionFallas) *PostgresRepository {
	return &PostgresRepository{db: db, simulacion: simulacion}
}

// BuscarUsuarioPorCredenciales busca un usuario por email y password.
func (r *PostgresRepository) BuscarUsuarioPorCredenciales(ctx context.Context, email, password string) (*models.Usuario, error) {
	query := `SELECT id, email, password, rol FROM usuarios WHERE email = $1 AND password = $2`
	var usuario models.Usuario

	err := r.db.QueryRow(ctx, query, email, password).Scan(
		&usuario.ID,
		&usuario.Email,
		&usuario.Password,
		&usuario.Rol,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Usuario no encontrado (no es error, simplemente no existe)
		}
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &usuario, nil
}

// BuscarUsuarioPorEmail busca un usuario por su email (sin validar contraseña).
func (r *PostgresRepository) BuscarUsuarioPorEmail(ctx context.Context, email string) (*models.Usuario, error) {
	query := `SELECT id, email, rol FROM usuarios WHERE email = $1`
	var usuario models.Usuario

	err := r.db.QueryRow(ctx, query, email).Scan(&usuario.ID, &usuario.Email, &usuario.Rol)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &usuario, nil
}

// ObtenerProductos devuelve la lista completa de productos.
func (r *PostgresRepository) ObtenerProductos(ctx context.Context) ([]models.Producto, error) {
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria FROM productos ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %w", err)
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		err := rows.Scan(
			&p.ID,
			&p.Nombre,
			&p.Descripcion,
			&p.Precio,
			&p.Categoria,
			&p.Img,
			&p.Galeria,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear producto: %w", err)
		}
		productos = append(productos, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar productos: %w", err)
	}

	return productos, nil
}

// ObtenerProductoPorID devuelve un producto específico por su ID.
func (r *PostgresRepository) ObtenerProductoPorID(ctx context.Context, id int) (*models.Producto, error) {
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria FROM productos WHERE id = $1`
	var producto models.Producto

	err := r.db.QueryRow(ctx, query, id).Scan(
		&producto.ID,
		&producto.Nombre,
		&producto.Descripcion,
		&producto.Precio,
		&producto.Categoria,
		&producto.Img,
		&producto.Galeria,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Producto no encontrado
		}
		return nil, fmt.Errorf("error al buscar producto: %w", err)
	}

	return &producto, nil
}

// ObtenerRed devuelve la lista plana de referidos desde la base de datos.
// La construcción del árbol jerárquico se hace en memoria después.
func (r *PostgresRepository) ObtenerRed(ctx context.Context) ([]models.Referido, error) {
	query := `SELECT id, nombre, nivel, ventas, parent_id FROM referidos ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener referidos: %w", err)
	}
	defer rows.Close()

	var referidos []models.Referido
	for rows.Next() {
		var ref models.Referido
		var parentID *int
		err := rows.Scan(
			&ref.ID,
			&ref.Nombre,
			&ref.Nivel,
			&ref.Ventas,
			&parentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear referido: %w", err)
		}
		ref.ParentID = parentID
		referidos = append(referidos, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar referidos: %w", err)
	}

	return referidos, nil
}

// ConstruirArbolRed construye la jerarquía multinivel a partir de una lista plana de referidos.
// Esta función es pura (no accede a DB) y se usa después de ObtenerRed.
func ConstruirArbolRed(referidos []models.Referido) *models.Referido {
	// Agrupar por padre: los nodos sin parent_id son directos (nivel 1) del usuario.
	hijosDe := make(map[int][]models.Referido)
	var raices []models.Referido
	for _, ref := range referidos {
		if ref.ParentID == nil {
			raices = append(raices, ref)
		} else {
			hijosDe[*ref.ParentID] = append(hijosDe[*ref.ParentID], ref)
		}
	}

	// Armar cada subárbol de forma recursiva para no perder niveles profundos (nivel 3).
	var armar func(nodo models.Referido) models.Referido
	armar = func(nodo models.Referido) models.Referido {
		for _, hijo := range hijosDe[nodo.ID] {
			nodo.Hijos = append(nodo.Hijos, armar(hijo))
		}
		return nodo
	}

	// Retornar el nodo raíz ficticio "Tú" con todos los referidos como hijos
	raiz := &models.Referido{
		ID:     0,
		Nombre: "Tú",
		Nivel:  0,
		Ventas: 2400, // Ventas propias del usuario (hardcodeado por ahora)
		Hijos:  make([]models.Referido, 0, len(raices)),
	}

	for _, r := range raices {
		raiz.Hijos = append(raiz.Hijos, armar(r))
	}

	return raiz
}

// CrearPedido guarda el pedido y su detalle dentro de una única transacción ACID:
// BEGIN → INSERT pedidos → INSERT pedido_items + descuento de stock → total → COMMIT.
// Si cualquier paso falla, el defer ejecuta ROLLBACK y no queda ningún "pedido fantasma".
func (r *PostgresRepository) CrearPedido(ctx context.Context, usuarioID int, req models.CheckoutRequest) (*models.Pedido, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar la transacción: %w", err)
	}
	// ROLLBACK si se sale de la función sin COMMIT (tras un COMMIT exitoso no hace nada).
	defer func() {
		if err := tx.Rollback(context.Background()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("Advertencia: error en ROLLBACK: %v", err)
		}
	}()

	// 1. Registro principal del pedido.
	pedido := &models.Pedido{UsuarioID: usuarioID, Estado: "confirmado"}
	err = tx.QueryRow(ctx,
		`INSERT INTO pedidos (usuario_id, cliente_nombre, cliente_email, direccion, ciudad, telefono, metodo_pago)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		usuarioID, req.Cliente.Nombre, req.Cliente.Email, req.Cliente.Direccion,
		req.Cliente.Ciudad, req.Cliente.Telefono, req.MetodoPago,
	).Scan(&pedido.ID)
	if err != nil {
		return nil, fmt.Errorf("error al insertar el pedido: %w", err)
	}

	if r.simulacion.FallaAntesDeItems {
		return nil, errors.New("falla simulada: caída del sistema antes de procesar los items")
	}

	// 2. Detalle del pedido: el precio se toma de PostgreSQL y se descuenta el stock simulado.
	pedido.Items = make([]models.PedidoItem, 0, len(req.Items))
	for _, item := range req.Items {
		linea := models.PedidoItem{ProductoID: item.ProductoID, Cantidad: item.Cantidad}

		// La FK pedido_items.producto_id aborta la transacción si el producto no existe.
		err = tx.QueryRow(ctx,
			`INSERT INTO pedido_items (pedido_id, producto_id, cantidad, precio_unitario)
			 VALUES ($1, $2, $3, COALESCE((SELECT precio FROM productos WHERE id = $2), 0))
			 RETURNING precio_unitario, subtotal`,
			pedido.ID, item.ProductoID, item.Cantidad,
		).Scan(&linea.PrecioUnitario, &linea.Subtotal)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return nil, &ErrValidacion{Mensaje: fmt.Sprintf("El producto %d no existe", item.ProductoID)}
			}
			return nil, fmt.Errorf("error al insertar el detalle del pedido: %w", err)
		}

		// Descuento atómico del stock: solo si alcanza la cantidad solicitada.
		err = tx.QueryRow(ctx,
			`UPDATE productos SET stock = stock - $1 WHERE id = $2 AND stock >= $1 RETURNING nombre`,
			item.Cantidad, item.ProductoID,
		).Scan(&linea.Nombre)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, &ErrValidacion{Mensaje: fmt.Sprintf("Stock insuficiente para el producto %d", item.ProductoID)}
			}
			return nil, fmt.Errorf("error al actualizar el stock: %w", err)
		}

		pedido.Items = append(pedido.Items, linea)
	}

	// 3. Total del pedido calculado en la base de datos a partir del detalle.
	err = tx.QueryRow(ctx,
		`UPDATE pedidos SET total = (SELECT COALESCE(SUM(subtotal), 0) FROM pedido_items WHERE pedido_id = $1)
		 WHERE id = $1 RETURNING total, creado_en`,
		pedido.ID,
	).Scan(&pedido.Total, &pedido.Fecha)
	if err != nil {
		return nil, fmt.Errorf("error al calcular el total del pedido: %w", err)
	}

	if r.simulacion.PausaAntesDeCommit > 0 {
		log.Printf("Simulación: pausa de %s antes del COMMIT del pedido %d", r.simulacion.PausaAntesDeCommit, pedido.ID)
		time.Sleep(r.simulacion.PausaAntesDeCommit)
	}

	// 4. COMMIT: a partir de aquí el pedido es definitivo.
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	pedido.Numero = fmt.Sprintf("MC-%06d", pedido.ID)
	return pedido, nil
}

// AplicarComisiones suma la venta al nodo del comprador en la red y acredita comisiones
// a su cadena ascendente (niveles 3, 2 y 1) en una sola transacción.
// Las tasas por nivel son las mismas que usa el frontend (TASA_COMISION): 10 %, 5 % y 2 %.
// Si el usuario no pertenece a la red, no hace nada.
func (r *PostgresRepository) AplicarComisiones(ctx context.Context, pedidoID, usuarioID int, total float64) ([]models.Comision, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción de comisiones: %w", err)
	}
	defer tx.Rollback(context.Background())

	// 1. La compra suma a las ventas del nodo del comprador (alimenta los KPIs del Dashboard).
	var nodoID int
	err = tx.QueryRow(ctx,
		`UPDATE referidos SET ventas = ventas + $1 WHERE usuario_id = $2 RETURNING id`,
		total, usuarioID,
	).Scan(&nodoID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al actualizar ventas del referido: %w", err)
	}

	// 2. Recorrer la cadena ascendente, registrar cada comisión y acumular el saldo.
	query := `
		WITH RECURSIVE ascendentes AS (
			SELECT id, nivel, parent_id FROM referidos WHERE id = $1
			UNION ALL
			SELECT r.id, r.nivel, r.parent_id
			FROM referidos r JOIN ascendentes a ON r.id = a.parent_id
		),
		tasas (nivel, tasa) AS (VALUES (1, 0.10), (2, 0.05), (3, 0.02)),
		nuevas AS (
			INSERT INTO comisiones (pedido_id, referido_id, nivel, monto)
			SELECT $2, a.id, a.nivel, ROUND($3::numeric * t.tasa, 2)
			FROM ascendentes a JOIN tasas t USING (nivel)
			RETURNING referido_id, nivel, monto
		)
		UPDATE referidos ref
		SET saldo_comision = ref.saldo_comision + nuevas.monto
		FROM nuevas
		WHERE ref.id = nuevas.referido_id
		RETURNING ref.id, ref.nombre, nuevas.nivel, nuevas.monto`
	rows, err := tx.Query(ctx, query, nodoID, pedidoID, total)
	if err != nil {
		return nil, fmt.Errorf("error al calcular comisiones: %w", err)
	}
	var comisiones []models.Comision
	for rows.Next() {
		var com models.Comision
		if err := rows.Scan(&com.ReferidoID, &com.Nombre, &com.Nivel, &com.Monto); err != nil {
			rows.Close()
			return nil, fmt.Errorf("error al leer comisión: %w", err)
		}
		comisiones = append(comisiones, com)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al calcular comisiones: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error al confirmar comisiones: %w", err)
	}
	return comisiones, nil
}
