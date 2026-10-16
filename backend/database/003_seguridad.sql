-- Migración 003: seguridad, hashing bcrypt e índices de rendimiento
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Índices para optimizar consultas de red, pedidos y catálogo
CREATE INDEX IF NOT EXISTS idx_referidos_parent_id ON referidos(parent_id);
CREATE INDEX IF NOT EXISTS idx_pedidos_cliente_email ON pedidos(cliente_email);
CREATE INDEX IF NOT EXISTS idx_pedido_items_pedido_id ON pedido_items(pedido_id);
CREATE INDEX IF NOT EXISTS idx_productos_categoria ON productos(categoria);

-- Actualizar contraseñas existentes a hashes bcrypt si aún están en texto plano
UPDATE usuarios
SET password = crypt(password, gen_salt('bf'))
WHERE password NOT LIKE '$2a$%' AND password NOT LIKE '$2b$%';
