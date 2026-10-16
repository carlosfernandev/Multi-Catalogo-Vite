// Paquete repository maneja el acceso a datos de la base de datos.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
}

// PostgresRepository implementa DBRepository usando un pool de pgx.
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository crea el repositorio inyectando el pool de conexiones.
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
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
	// Crear mapa por ID para acceso rápido
	refMap := make(map[int]*models.Referido)
	for i := range referidos {
		refMap[referidos[i].ID] = &referidos[i]
	}

	// Construir árbol jerárquico
	var raices []*models.Referido
	for i := range referidos {
		ref := &referidos[i]
		if ref.ParentID == nil {
			// Es un nodo raíz (nivel 1, directo del usuario)
			raices = append(raices, ref)
		} else {
			// Es hijo de otro nodo
			if padre, exists := refMap[*ref.ParentID]; exists {
				padre.Hijos = append(padre.Hijos, *ref)
			}
		}
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
		raiz.Hijos = append(raiz.Hijos, *r)
	}

	return raiz
}
