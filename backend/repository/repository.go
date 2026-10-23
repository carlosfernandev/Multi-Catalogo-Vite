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
	"github.com/jackc/pgx/v5/pgxpool"
	"backend/models"
)

// ErrDuplicado indica conflicto de clave única (HTTP 409).
var ErrDuplicado = errors.New("el recurso ya existe")

// DBRepository es el contrato de acceso a datos que se inyecta en los controladores.
// Los controladores dependen de esta interfaz, no de PostgreSQL directamente.
type DBRepository interface {
	BuscarUsuarioPorCredenciales(ctx context.Context, email, password string) (*models.Usuario, error)
	RegistrarUsuario(ctx context.Context, req models.RegisterRequest) (*models.Usuario, error)
	ObtenerProductos(ctx context.Context) ([]models.Producto, error)
	ObtenerProductoPorID(ctx context.Context, id int) (*models.Producto, error)
	ObtenerProductosAdmin(ctx context.Context, q string) ([]models.Producto, error)
	CrearProducto(ctx context.Context, req models.CrearProductoRequest) (*models.Producto, error)
	ActualizarProducto(ctx context.Context, id int, req models.ActualizarProductoRequest) (*models.Producto, error)
	EliminarProductoLogico(ctx context.Context, id int) error
	ObtenerRed(ctx context.Context) ([]models.Referido, error)
	BuscarUsuarioPorEmail(ctx context.Context, email string) (*models.Usuario, error)
	ObtenerUsuariosAdmin(ctx context.Context, q string) ([]models.Usuario, error)
	ActualizarUsuarioAdmin(ctx context.Context, id int, req models.ActualizarUsuarioRequest) (*models.Usuario, error)
	CrearPedido(ctx context.Context, usuarioID int, req models.CheckoutRequest) (*models.Pedido, error)
	AplicarComisiones(ctx context.Context, pedidoID, usuarioID int, total float64) ([]models.Comision, error)
	ObtenerLiquidaciones(ctx context.Context, email string) (*models.ResumenLiquidaciones, error)
}

// ErrNoEncontrado indica que el registro solicitado no existe en la base de datos (mapeado a 404/401).
var ErrNoEncontrado = errors.New("recurso no encontrado")

// ErrValidacion indica un error de datos del cliente (HTTP 400), p. ej. producto inexistente o sin stock.
type ErrValidacion struct {
	Mensaje string
}

func (e *ErrValidacion) Error() string { return e.Mensaje }

// ErrCantidadInvalida indica que la cantidad solicitada para un ítem no es estrictamente positiva.
type ErrCantidadInvalida struct {
	ProductoID int
}

func (e *ErrCantidadInvalida) Error() string { return "cantidad debe ser mayor que cero" }

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
	query := `SELECT id, email, password, rol, COALESCE(nombre, ''), activo, referido_id FROM usuarios WHERE email = $1 AND password = $2`
	var usuario models.Usuario

	err := r.db.QueryRow(ctx, query, email, password).Scan(
		&usuario.ID,
		&usuario.Email,
		&usuario.Password,
		&usuario.Rol,
		&usuario.Nombre,
		&usuario.Activo,
		&usuario.ReferidoID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &usuario, nil
}

// RegistrarUsuario crea una cuenta de cliente y su nodo en la red multinivel en una sola transacción.
func (r *PostgresRepository) RegistrarUsuario(ctx context.Context, req models.RegisterRequest) (*models.Usuario, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transaccion de registro: %w", err)
	}
	defer tx.Rollback(context.Background())

	// 1. Validar si el email ya existe
	var existe int
	err = tx.QueryRow(ctx, "SELECT id FROM usuarios WHERE email = $1", req.Email).Scan(&existe)
	if err == nil {
		return nil, ErrDuplicado
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error validando email: %w", err)
	}

	// 2. Determinar parent_id y nivel en la red según el código de referido (si se proporciona)
	var parentID *int
	nivel := 1
	if req.Ref != "" {
		var pID, pNivel int
		errRef := tx.QueryRow(ctx, "SELECT id, nivel FROM referidos WHERE UPPER(codigo) = UPPER($1)", req.Ref).Scan(&pID, &pNivel)
		if errors.Is(errRef, pgx.ErrNoRows) {
			return nil, &ErrValidacion{Mensaje: fmt.Sprintf("El código de referido '%s' no es válido", req.Ref)}
		} else if errRef != nil {
			return nil, fmt.Errorf("error validando codigo de referido: %w", errRef)
		}
		parentID = &pID
		if pNivel < 3 {
			nivel = pNivel + 1
		} else {
			nivel = 3
		}
	}

	// 3. Insertar el nodo correspondiente en referidos
	var nuevoReferidoID int
	// Generar un código único simple para este nuevo usuario
	codigoNuevo := fmt.Sprintf("REF%d", time.Now().Unix()%100000)
	err = tx.QueryRow(ctx,
		`INSERT INTO referidos (nombre, nivel, ventas, parent_id, codigo)
		 VALUES ($1, $2, 0, $3, $4) RETURNING id`,
		req.Nombre, nivel, parentID, codigoNuevo,
	).Scan(&nuevoReferidoID)
	if err != nil {
		return nil, fmt.Errorf("error creando nodo en red de referidos: %w", err)
	}

	// 4. Insertar el usuario vinculado
	var nuevoUsuario models.Usuario
	err = tx.QueryRow(ctx,
		`INSERT INTO usuarios (email, password, rol, nombre, activo, referido_id)
		 VALUES ($1, $2, 'cliente', $3, true, $4)
		 RETURNING id, email, rol, nombre, activo, referido_id`,
		req.Email, req.Password, req.Nombre, nuevoReferidoID,
	).Scan(
		&nuevoUsuario.ID,
		&nuevoUsuario.Email,
		&nuevoUsuario.Rol,
		&nuevoUsuario.Nombre,
		&nuevoUsuario.Activo,
		&nuevoUsuario.ReferidoID,
	)
	if err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error al confirmar registro: %w", err)
	}

	return &nuevoUsuario, nil
}

// BuscarUsuarioPorEmail busca un usuario por su email.
func (r *PostgresRepository) BuscarUsuarioPorEmail(ctx context.Context, email string) (*models.Usuario, error) {
	query := `SELECT id, email, password, rol, COALESCE(nombre, ''), activo, referido_id FROM usuarios WHERE email = $1`
	var usuario models.Usuario

	err := r.db.QueryRow(ctx, query, email).Scan(
		&usuario.ID,
		&usuario.Email,
		&usuario.Password,
		&usuario.Rol,
		&usuario.Nombre,
		&usuario.Activo,
		&usuario.ReferidoID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &usuario, nil
}

// ObtenerUsuariosAdmin devuelve la lista de todos los usuarios registrados para el panel de administración.
func (r *PostgresRepository) ObtenerUsuariosAdmin(ctx context.Context, q string) ([]models.Usuario, error) {
	query := `SELECT id, email, rol, COALESCE(nombre, ''), activo, referido_id, creado_en FROM usuarios`
	var rows pgx.Rows
	var err error
	if q != "" {
		query += ` WHERE email ILIKE $1 OR nombre ILIKE $1 ORDER BY id`
		rows, err = r.db.Query(ctx, query, "%"+q+"%")
	} else {
		query += ` ORDER BY id`
		rows, err = r.db.Query(ctx, query)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios admin: %w", err)
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Rol,
			&u.Nombre,
			&u.Activo,
			&u.ReferidoID,
			&u.CreadoEn,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		usuarios = append(usuarios, u)
	}

	return usuarios, rows.Err()
}

// ActualizarUsuarioAdmin permite modificar el rol, nombre y estado activo/inactivo de un usuario.
func (r *PostgresRepository) ActualizarUsuarioAdmin(ctx context.Context, id int, req models.ActualizarUsuarioRequest) (*models.Usuario, error) {
	var u models.Usuario
	queryGet := `SELECT id, email, rol, COALESCE(nombre, ''), activo, referido_id, creado_en FROM usuarios WHERE id = $1`
	err := r.db.QueryRow(ctx, queryGet, id).Scan(
		&u.ID, &u.Email, &u.Rol, &u.Nombre, &u.Activo, &u.ReferidoID, &u.CreadoEn,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, fmt.Errorf("error al buscar usuario para actualizar: %w", err)
	}

	if req.Nombre != nil {
		u.Nombre = *req.Nombre
	}
	if req.Rol != nil {
		if *req.Rol != "admin" && *req.Rol != "cliente" {
			return nil, &ErrValidacion{Mensaje: "El rol debe ser 'admin' o 'cliente'"}
		}
		u.Rol = *req.Rol
	}
	if req.Activo != nil {
		u.Activo = *req.Activo
	}

	queryUpdate := `
		UPDATE usuarios
		SET nombre = $1, rol = $2, activo = $3
		WHERE id = $4
		RETURNING id, email, rol, COALESCE(nombre, ''), activo, referido_id, creado_en
	`
	err = r.db.QueryRow(ctx, queryUpdate, u.Nombre, u.Rol, u.Activo, id).Scan(
		&u.ID, &u.Email, &u.Rol, &u.Nombre, &u.Activo, &u.ReferidoID, &u.CreadoEn,
	)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	return &u, nil
}

// ObtenerProductos devuelve la lista de productos activos para el catálogo público.
func (r *PostgresRepository) ObtenerProductos(ctx context.Context) ([]models.Producto, error) {
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo FROM productos WHERE activo = true ORDER BY id`
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
			&p.Stock,
			&p.Activo,
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
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo FROM productos WHERE id = $1`
	var producto models.Producto

	err := r.db.QueryRow(ctx, query, id).Scan(
		&producto.ID,
		&producto.Nombre,
		&producto.Descripcion,
		&producto.Precio,
		&producto.Categoria,
		&producto.Img,
		&producto.Galeria,
		&producto.Stock,
		&producto.Activo,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado // Producto no encontrado
		}
		return nil, fmt.Errorf("error al buscar producto: %w", err)
	}

	return &producto, nil
}

// ObtenerProductosAdmin devuelve todos los productos (activos e inactivos) con soporte opcional de búsqueda.
func (r *PostgresRepository) ObtenerProductosAdmin(ctx context.Context, q string) ([]models.Producto, error) {
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo FROM productos`
	var rows pgx.Rows
	var err error
	if q != "" {
		query += ` WHERE nombre ILIKE $1 OR descripcion ILIKE $1 OR categoria ILIKE $1 ORDER BY id`
		rows, err = r.db.Query(ctx, query, "%"+q+"%")
	} else {
		query += ` ORDER BY id`
		rows, err = r.db.Query(ctx, query)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos admin: %w", err)
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
			&p.Stock,
			&p.Activo,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear producto: %w", err)
		}
		productos = append(productos, p)
	}

	return productos, rows.Err()
}

// CrearProducto inserta un nuevo artículo en la tabla productos.
func (r *PostgresRepository) CrearProducto(ctx context.Context, req models.CrearProductoRequest) (*models.Producto, error) {
	if req.Galeria == nil {
		req.Galeria = []string{}
	}
	query := `
		INSERT INTO productos (nombre, descripcion, precio, categoria, img, galeria, stock, activo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		RETURNING id, nombre, descripcion, precio, categoria, img, galeria, stock, activo
	`
	var p models.Producto
	err := r.db.QueryRow(ctx, query, req.Nombre, req.Descripcion, req.Precio, req.Categoria, req.Img, req.Galeria, req.Stock).Scan(
		&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria, &p.Stock, &p.Activo,
	)
	if err != nil {
		return nil, fmt.Errorf("error al insertar producto: %w", err)
	}
	return &p, nil
}

// ActualizarProducto modifica campos seleccionados de un producto existente.
func (r *PostgresRepository) ActualizarProducto(ctx context.Context, id int, req models.ActualizarProductoRequest) (*models.Producto, error) {
	prod, err := r.ObtenerProductoPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Nombre != nil {
		prod.Nombre = *req.Nombre
	}
	if req.Descripcion != nil {
		prod.Descripcion = *req.Descripcion
	}
	if req.Precio != nil {
		prod.Precio = *req.Precio
	}
	if req.Categoria != nil {
		prod.Categoria = *req.Categoria
	}
	if req.Img != nil {
		prod.Img = *req.Img
	}
	if req.Galeria != nil {
		prod.Galeria = *req.Galeria
	}
	if req.Stock != nil {
		prod.Stock = *req.Stock
	}
	if req.Activo != nil {
		prod.Activo = *req.Activo
	}

	query := `
		UPDATE productos 
		SET nombre = $1, descripcion = $2, precio = $3, categoria = $4, img = $5, galeria = $6, stock = $7, activo = $8
		WHERE id = $9
		RETURNING id, nombre, descripcion, precio, categoria, img, galeria, stock, activo
	`
	var p models.Producto
	err = r.db.QueryRow(ctx, query, prod.Nombre, prod.Descripcion, prod.Precio, prod.Categoria, prod.Img, prod.Galeria, prod.Stock, prod.Activo, id).Scan(
		&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria, &p.Stock, &p.Activo,
	)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar producto: %w", err)
	}
	return &p, nil
}

// EliminarProductoLogico marca un producto con activo = false para no romper pedidos históricos.
func (r *PostgresRepository) EliminarProductoLogico(ctx context.Context, id int) error {
	res, err := r.db.Exec(ctx, `UPDATE productos SET activo = false WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error al desactivar producto: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrNoEncontrado
	}
	return nil
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

	// 2. Detalle del pedido: bloqueo pesimista con FOR UPDATE, validación de stock y descuento.
	pedido.Items = make([]models.PedidoItem, 0, len(req.Items))
	for _, item := range req.Items {
		if item.Cantidad <= 0 {
			return nil, &ErrCantidadInvalida{ProductoID: item.ProductoID}
		}

		linea := models.PedidoItem{ProductoID: item.ProductoID, Cantidad: item.Cantidad}

		// Bloqueo pesimista de fila (FOR UPDATE) para prevenir condiciones de carrera concurrentes.
		var stockActual int
		err = tx.QueryRow(ctx,
			`SELECT nombre, precio, stock FROM productos WHERE id = $1 FOR UPDATE`,
			item.ProductoID,
		).Scan(&linea.Nombre, &linea.PrecioUnitario, &stockActual)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, &ErrValidacion{Mensaje: fmt.Sprintf("El producto %d no existe", item.ProductoID)}
			}
			return nil, fmt.Errorf("error al consultar el producto con bloqueo: %w", err)
		}

		if stockActual < item.Cantidad {
			return nil, &ErrValidacion{Mensaje: fmt.Sprintf("Stock insuficiente para el producto %d", item.ProductoID)}
		}

		// Descontar stock del producto bloqueado.
		_, err = tx.Exec(ctx,
			`UPDATE productos SET stock = stock - $1 WHERE id = $2`,
			item.Cantidad, item.ProductoID,
		)
		if err != nil {
			return nil, fmt.Errorf("error al actualizar el stock: %w", err)
		}

		// Insertar en pedido_items.
		err = tx.QueryRow(ctx,
			`INSERT INTO pedido_items (pedido_id, producto_id, cantidad, precio_unitario)
			 VALUES ($1, $2, $3, $4) RETURNING subtotal`,
			pedido.ID, item.ProductoID, item.Cantidad, linea.PrecioUnitario,
		).Scan(&linea.Subtotal)
		if err != nil {
			return nil, fmt.Errorf("error al insertar el detalle del pedido: %w", err)
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

// ObtenerLiquidaciones devuelve el historial de comisiones registradas en la BD.
// Si se pasa un email de usuario específico y no es admin, filtra por el usuario;
// si es admin o no se filtra, devuelve el histórico general de liquidaciones.
func (r *PostgresRepository) ObtenerLiquidaciones(ctx context.Context, email string) (*models.ResumenLiquidaciones, error) {
	query := `
		SELECT 
			c.id, c.pedido_id, ('PED-' || c.pedido_id) as pedido_numero,
			c.referido_id, COALESCE(r.nombre, 'Referido ' || c.referido_id) as referido_nombre,
			c.nivel, c.monto, c.creado_en
		FROM comisiones c
		LEFT JOIN pedidos p ON c.pedido_id = p.id
		LEFT JOIN referidos r ON c.referido_id = r.id
		LEFT JOIN usuarios u ON r.usuario_id = u.id
		WHERE ($1 = '' OR $1 = 'admin@upse.edu.ec' OR u.email = $1)
		ORDER BY c.creado_en DESC, c.id DESC
		LIMIT 100`

	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("error al consultar liquidaciones: %w", err)
	}
	defer rows.Close()

	var liquidaciones []models.LiquidacionRegistro
	var total float64

	for rows.Next() {
		var reg models.LiquidacionRegistro
		if err := rows.Scan(
			&reg.ID, &reg.PedidoID, &reg.PedidoNumero,
			&reg.ReferidoID, &reg.ReferidoNombre,
			&reg.Nivel, &reg.Monto, &reg.Fecha,
		); err != nil {
			return nil, fmt.Errorf("error al escanear registro de liquidacion: %w", err)
		}
		total += reg.Monto
		liquidaciones = append(liquidaciones, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando liquidaciones: %w", err)
	}

	if liquidaciones == nil {
		liquidaciones = []models.LiquidacionRegistro{}
	}

	return &models.ResumenLiquidaciones{
		TotalComisiones: total,
		Cantidad:        len(liquidaciones),
		Liquidaciones:   liquidaciones,
	}, nil
}
