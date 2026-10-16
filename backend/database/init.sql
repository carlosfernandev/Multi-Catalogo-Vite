-- Script de inicialización de la base de datos multicatalogo
-- Se ejecuta automáticamente al levantar el contenedor PostgreSQL

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    rol VARCHAR(50) NOT NULL CHECK (rol IN ('admin', 'cliente'))
);

-- Tabla de productos
CREATE TABLE IF NOT EXISTS productos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    descripcion TEXT,
    precio DECIMAL(10, 2) NOT NULL,
    categoria VARCHAR(100),
    img VARCHAR(500),
    galeria TEXT[]
);

-- Tabla de referidos (estructura jerárquica multinivel)
CREATE TABLE IF NOT EXISTS referidos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    nivel INTEGER NOT NULL CHECK (nivel IN (1, 2, 3)),
    ventas DECIMAL(10, 2) NOT NULL DEFAULT 0,
    parent_id INTEGER REFERENCES referidos(id) ON DELETE CASCADE
);

-- Insertar usuarios (admin y cliente)
INSERT INTO usuarios (email, password, rol) VALUES
    ('admin@upse.edu.ec', '123456', 'admin'),
    ('cliente@upse.edu.ec', '123456', 'cliente')
ON CONFLICT (email) DO NOTHING;

-- Insertar productos (8 productos del mock de frontend)
INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria) VALUES
    (1, 'Serum Revitalizante', 'Serum concentrado con vitamina C y ácido hialurónico. Ilumina la piel, reduce manchas y aporta hidratación profunda desde la primera aplicación.', 45.00, 'Serum', 'https://picsum.photos/seed/serum/600', ARRAY['https://picsum.photos/seed/serum1/600', 'https://picsum.photos/seed/serum2/600', 'https://picsum.photos/seed/serum3/600']),
    (2, 'Crema Hidratante Pro', 'Crema facial de textura ligera con niacinamida y manteca de karité. Hidratación de 24 horas y barrera cutánea fortalecida para todo tipo de piel.', 32.50, 'Crema', 'https://picsum.photos/seed/crema/600', ARRAY['https://picsum.photos/seed/crema1/600', 'https://picsum.photos/seed/crema2/600', 'https://picsum.photos/seed/crema3/600']),
    (3, 'Tónico Purificante', 'Tónico sin alcohol con hamamelis y agua de rosas. Limpia poros, equilibra el pH y prepara la piel para el resto de la rutina diaria.', 28.00, 'Tónico', 'https://picsum.photos/seed/tonico/600', ARRAY['https://picsum.photos/seed/tonico1/600', 'https://picsum.photos/seed/tonico2/600', 'https://picsum.photos/seed/tonico3/600']),
    (4, 'Mascarilla Nocturna', 'Mascarilla de noche con colágeno y vitamina E. Repara la piel mientras duermes y recupera la luminosidad al despertar.', 50.00, 'Mascarilla', 'https://picsum.photos/seed/mascarilla/600', ARRAY['https://picsum.photos/seed/mascarilla1/600', 'https://picsum.photos/seed/mascarilla2/600', 'https://picsum.photos/seed/mascarilla3/600']),
    (5, 'Serum Anti-Edad Retinol', 'Serum de retinol encapsulado que suaviza líneas de expresión y mejora la firmeza. Uso nocturno recomendado.', 58.00, 'Serum', 'https://picsum.photos/seed/retinol/600', ARRAY['https://picsum.photos/seed/retinol1/600', 'https://picsum.photos/seed/retinol2/600', 'https://picsum.photos/seed/retinol3/600']),
    (6, 'Crema Contorno de Ojos', 'Contorno de ojos con cafeína y péptidos. Reduce bolsas y ojeras, hidrata la zona más delicada del rostro.', 26.00, 'Crema', 'https://picsum.photos/seed/ojos/600', ARRAY['https://picsum.photos/seed/ojos1/600', 'https://picsum.photos/seed/ojos2/600', 'https://picsum.photos/seed/ojos3/600']),
    (7, 'Tónico Exfoliante AHA-BHA', 'Exfoliación química suave con ácidos AHA y BHA. Renueva la textura de la piel y desobstruye los poros.', 34.00, 'Tónico', 'https://picsum.photos/seed/exfoliante/600', ARRAY['https://picsum.photos/seed/exfoliante1/600', 'https://picsum.photos/seed/exfoliante2/600', 'https://picsum.photos/seed/exfoliante3/600']),
    (8, 'Kit Rutina Completa', 'Kit integral con serum, crema, tónico y mascarilla. La rutina perfecta para comenzar: limpieza, tratamiento e hidratación.', 129.00, 'Kit', 'https://picsum.photos/seed/kit/600', ARRAY['https://picsum.photos/seed/kit1/600', 'https://picsum.photos/seed/kit2/600', 'https://picsum.photos/seed/kit3/600'])
ON CONFLICT (id) DO NOTHING;

-- Insertar referidos (estructura jerárquica de 3 niveles, 7 referidos)
-- Raíz del árbol (id 0 es el usuario autenticado, no se guarda en DB)
-- Nivel 1 (directos)
INSERT INTO referidos (id, nombre, nivel, ventas, parent_id) VALUES
    (1, 'Ana García', 1, 1200.00, NULL),
    (2, 'Luis Poveda', 1, 850.00, NULL),
    (3, 'Marta Sánchez', 1, 430.00, NULL)
ON CONFLICT (id) DO NOTHING;

-- Nivel 2 (indirectos de Ana)
INSERT INTO referidos (id, nombre, nivel, ventas, parent_id) VALUES
    (4, 'Carlos Ruiz', 2, 500.00, 1),
    (5, 'Sofía León', 2, 430.00, 1)
ON CONFLICT (id) DO NOTHING;

-- Nivel 2 (indirectos de Luis)
INSERT INTO referidos (id, nombre, nivel, ventas, parent_id) VALUES
    (6, 'Marco Díaz', 2, 380.00, 2)
ON CONFLICT (id) DO NOTHING;

-- Nivel 3 (indirectos de Carlos)
INSERT INTO referidos (id, nombre, nivel, ventas, parent_id) VALUES
    (7, 'Diana Paz', 3, 300.00, 4)
ON CONFLICT (id) DO NOTHING;

-- =====================================================================
-- Fase 4: checkout transaccional (pedidos + detalle + stock simulado)
-- Sentencias idempotentes: se pueden re-ejecutar sobre una DB existente.
-- =====================================================================

-- Stock simulado de cada producto (se descuenta dentro de la transacción del checkout)
ALTER TABLE productos ADD COLUMN IF NOT EXISTS stock INTEGER NOT NULL DEFAULT 100;
ALTER TABLE productos DROP CONSTRAINT IF EXISTS productos_stock_no_negativo;
ALTER TABLE productos ADD CONSTRAINT productos_stock_no_negativo CHECK (stock >= 0);

-- Registro principal del pedido
CREATE TABLE IF NOT EXISTS pedidos (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
    cliente_nombre VARCHAR(255) NOT NULL,
    cliente_email VARCHAR(255) NOT NULL,
    direccion VARCHAR(500) NOT NULL,
    ciudad VARCHAR(100) NOT NULL,
    telefono VARCHAR(50) NOT NULL,
    metodo_pago VARCHAR(50) NOT NULL,
    total DECIMAL(12, 2) NOT NULL DEFAULT 0,
    creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Detalle del pedido (la FK a productos aborta la transacción si el producto no existe)
CREATE TABLE IF NOT EXISTS pedido_items (
    id SERIAL PRIMARY KEY,
    pedido_id INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    producto_id INTEGER NOT NULL REFERENCES productos(id),
    cantidad INTEGER NOT NULL CHECK (cantidad > 0),
    precio_unitario DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(12, 2) GENERATED ALWAYS AS (cantidad * precio_unitario) STORED
);
