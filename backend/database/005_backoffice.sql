-- Migración 005: Back-office MultiCatálogo (Fase 0)
-- Tablas y columnas requeridas para registro con referidos, estados de pedidos, ledger y liquidaciones

-- 1. Modificaciones a la tabla usuarios
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS nombre VARCHAR(255) DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS referido_id INTEGER REFERENCES referidos(id) ON DELETE SET NULL;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Actualizar nombres y referidos del seed existente
UPDATE usuarios SET nombre = 'Administrador', activo = true WHERE email = 'admin@upse.edu.ec';
UPDATE usuarios SET nombre = 'Diana Paz (Cliente)', activo = true, referido_id = 7 WHERE email = 'cliente@upse.edu.ec';

-- 2. Modificaciones a la tabla referidos (código de invitación único)
ALTER TABLE referidos ADD COLUMN IF NOT EXISTS codigo VARCHAR(50) UNIQUE;

-- Asignar códigos a la red semilla inicial
UPDATE referidos SET codigo = 'ANA' WHERE id = 1 AND codigo IS NULL;
UPDATE referidos SET codigo = 'LUIS' WHERE id = 2 AND codigo IS NULL;
UPDATE referidos SET codigo = 'MARTA' WHERE id = 3 AND codigo IS NULL;
UPDATE referidos SET codigo = 'CARLOS' WHERE id = 4 AND codigo IS NULL;
UPDATE referidos SET codigo = 'SOFIA' WHERE id = 5 AND codigo IS NULL;
UPDATE referidos SET codigo = 'MARCO' WHERE id = 6 AND codigo IS NULL;
UPDATE referidos SET codigo = 'DIANA' WHERE id = 7 AND codigo IS NULL;

-- Sincronizar secuencias para evitar colisión de IDs en nuevos registros
SELECT setval('referidos_id_seq', (SELECT COALESCE(MAX(id), 1) FROM referidos));
SELECT setval('usuarios_id_seq', (SELECT COALESCE(MAX(id), 1) FROM usuarios));
SELECT setval('productos_id_seq', (SELECT COALESCE(MAX(id), 1) FROM productos));

-- 3. Modificaciones a la tabla productos (borrado lógico)
ALTER TABLE productos ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;

-- 4. Modificaciones a la tabla pedidos
ALTER TABLE pedidos ADD COLUMN IF NOT EXISTS estado VARCHAR(50) NOT NULL DEFAULT 'pendiente';
ALTER TABLE pedidos ADD COLUMN IF NOT EXISTS usuario_id INTEGER REFERENCES usuarios(id);

-- 5. Tabla de liquidaciones
CREATE TABLE IF NOT EXISTS liquidaciones (
    id SERIAL PRIMARY KEY,
    referido_id INTEGER NOT NULL REFERENCES referidos(id) ON DELETE CASCADE,
    total DECIMAL(12, 2) NOT NULL DEFAULT 0,
    estado VARCHAR(50) NOT NULL DEFAULT 'pagada',
    creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pagada_en TIMESTAMPTZ
);

-- 6. Modificaciones a la tabla comisiones
ALTER TABLE comisiones ADD COLUMN IF NOT EXISTS estado VARCHAR(50) NOT NULL DEFAULT 'pendiente';
ALTER TABLE comisiones ADD COLUMN IF NOT EXISTS porcentaje DECIMAL(5, 2) DEFAULT 0;
ALTER TABLE comisiones ADD COLUMN IF NOT EXISTS liquidacion_id INTEGER REFERENCES liquidaciones(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_comisiones_referido_pedido ON comisiones(referido_id, pedido_id);
CREATE INDEX IF NOT EXISTS idx_comisiones_estado ON comisiones(estado);
CREATE INDEX IF NOT EXISTS idx_referidos_codigo ON referidos(codigo);
CREATE INDEX IF NOT EXISTS idx_pedidos_usuario_id ON pedidos(usuario_id);
