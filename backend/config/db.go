// Paquete config maneja la configuración de la base de datos.
package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// InitDB crea el pool de conexiones a PostgreSQL usando variables de entorno.
// El pool se devuelve (sin variables globales) para inyectarlo en el repositorio desde main.go.
func InitDB() (*pgxpool.Pool, error) {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Advertencia: no se pudo cargar .env: %v", err)
	}

	// Construir DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Crear pool de conexiones
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error al parsear configuración de DB: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error al crear pool de conexiones: %w", err)
	}

	// Verificar conexión
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error al conectar a PostgreSQL: %w", err)
	}

	log.Println("Conectado a PostgreSQL exitosamente")
	return pool, nil
}
