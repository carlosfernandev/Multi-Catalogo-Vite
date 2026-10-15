// Paquete config: cliente de Redis para el almacenamiento en caché (patrón Cache-Aside).
package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RDB es el cliente global de Redis. Es nil si Redis no está disponible.
var RDB *redis.Client

// InitRedis inicializa el cliente de Redis leyendo REDIS_URL (ej. redis://localhost:6379/0)
// o, en su defecto, REDIS_HOST, REDIS_PORT y REDIS_PASSWORD.
func InitRedis() error {
	var opts *redis.Options

	if url := os.Getenv("REDIS_URL"); url != "" {
		parsed, err := redis.ParseURL(url)
		if err != nil {
			return fmt.Errorf("error al parsear REDIS_URL: %w", err)
		}
		opts = parsed
	} else {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("REDIS_PORT")
		if port == "" {
			port = "6379"
		}
		opts = &redis.Options{
			Addr:     host + ":" + port,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}
	}

	// Tiempos cortos: si Redis cae, la API responde rápido desde PostgreSQL.
	opts.DialTimeout = 500 * time.Millisecond
	opts.ReadTimeout = 300 * time.Millisecond
	opts.WriteTimeout = 300 * time.Millisecond
	opts.MaxRetries = 1

	client := redis.NewClient(opts)

	// Verificar conexión
	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return fmt.Errorf("error al conectar a Redis (%s): %w", opts.Addr, err)
	}

	RDB = client
	log.Printf("Conectado a Redis exitosamente (%s)", opts.Addr)
	return nil
}

// CloseRedis cierra el cliente de Redis.
// Debe llamarse al cerrar la aplicación (defer en main).
func CloseRedis() {
	if RDB != nil {
		RDB.Close()
		log.Println("Conexión a Redis cerrada")
	}
}
