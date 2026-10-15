// Paquete cache implementa las operaciones genéricas del patrón Cache-Aside sobre Redis.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"backend/config"
	"github.com/redis/go-redis/v9"
)

// Obtener intenta leer la clave desde Redis y deserializar el JSON en destino.
// Retorna true si hubo Cache Hit. Cualquier error (Redis caído, JSON inválido) se trata como Cache Miss.
func Obtener(ctx context.Context, clave string, destino interface{}) bool {
	if config.RDB == nil {
		return false
	}

	datos, err := config.RDB.Get(ctx, clave).Bytes()
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
func Guardar(ctx context.Context, clave string, valor interface{}, ttl time.Duration) {
	if config.RDB == nil {
		return
	}

	datos, err := json.Marshal(valor)
	if err != nil {
		log.Printf("Advertencia: no se pudo serializar %q: %v", clave, err)
		return
	}

	if err := config.RDB.Set(ctx, clave, datos, ttl).Err(); err != nil {
		log.Printf("Advertencia: error al guardar %q en Redis: %v", clave, err)
	}
}
