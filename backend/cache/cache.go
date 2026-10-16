// Paquete cache implementa las operaciones del patrón Cache-Aside sobre Redis.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheRepository es el contrato que usan los controladores para la caché.
// Permite inyectar Redis en producción o una implementación falsa en pruebas.
type CacheRepository interface {
	// Obtener deserializa la clave en destino; retorna true si hubo Cache Hit.
	Obtener(ctx context.Context, clave string, destino interface{}) bool
	// Guardar serializa valor a JSON y lo almacena con el TTL indicado.
	Guardar(ctx context.Context, clave string, valor interface{}, ttl time.Duration)
	// Eliminar borra las claves indicadas (invalidación de caché).
	Eliminar(ctx context.Context, claves ...string)
	// EliminarPorPatron borra todas las claves que coinciden con el patrón (ej. "red:tree:*").
	EliminarPorPatron(ctx context.Context, patron string)
}

// RedisCache implementa CacheRepository con go-redis.
// Si el cliente es nil (Redis no disponible), todas las operaciones son no-op.
type RedisCache struct {
	cliente *redis.Client
}

// NewRedisCache crea la caché inyectando el cliente de Redis.
func NewRedisCache(cliente *redis.Client) *RedisCache {
	return &RedisCache{cliente: cliente}
}

// Obtener intenta leer la clave desde Redis y deserializar el JSON en destino.
// Cualquier error (Redis caído, JSON inválido) se trata como Cache Miss.
func (r *RedisCache) Obtener(ctx context.Context, clave string, destino interface{}) bool {
	if r.cliente == nil {
		return false
	}

	datos, err := r.cliente.Get(ctx, clave).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("Advertencia: error al leer %q de Redis: %v", clave, err)
		}
		return false
	}

	if err := json.Unmarshal(datos, destino); err != nil {
		log.Printf("Advertencia: JSON inválido en caché %q: %v", clave, err)
		return false
	}
	return true
}

// Guardar serializa valor a JSON y lo almacena en Redis con el TTL indicado.
// Un fallo al guardar no interrumpe la respuesta al cliente.
func (r *RedisCache) Guardar(ctx context.Context, clave string, valor interface{}, ttl time.Duration) {
	if r.cliente == nil {
		return
	}

	datos, err := json.Marshal(valor)
	if err != nil {
		log.Printf("Advertencia: no se pudo serializar %q: %v", clave, err)
		return
	}

	if err := r.cliente.Set(ctx, clave, datos, ttl).Err(); err != nil {
		log.Printf("Advertencia: error al guardar %q en Redis: %v", clave, err)
	}
}

// Eliminar borra las claves de Redis; un fallo solo se registra en el log.
func (r *RedisCache) Eliminar(ctx context.Context, claves ...string) {
	if r.cliente == nil || len(claves) == 0 {
		return
	}
	if err := r.cliente.Del(ctx, claves...).Err(); err != nil {
		log.Printf("Advertencia: error al eliminar %v de Redis: %v", claves, err)
	}
}

// EliminarPorPatron recorre las claves con SCAN (no bloquea Redis como KEYS) y las borra.
func (r *RedisCache) EliminarPorPatron(ctx context.Context, patron string) {
	if r.cliente == nil {
		return
	}
	iter := r.cliente.Scan(ctx, 0, patron, 100).Iterator()
	var claves []string
	for iter.Next(ctx) {
		claves = append(claves, iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Printf("Advertencia: error al buscar %q en Redis: %v", patron, err)
		return
	}
	r.Eliminar(ctx, claves...)
}
