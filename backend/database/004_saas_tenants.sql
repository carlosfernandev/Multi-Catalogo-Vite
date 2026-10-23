-- Migración 004: Capa SaaS Multitenant (Unidad 3)
-- Tablas para gestión de inquilinos / distribuidores y catálogos personalizados

CREATE TABLE IF NOT EXISTS saas_tenants (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    dominio_personalizado VARCHAR(255),
    contacto_email VARCHAR(255) NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT true,
    configuracion JSONB NOT NULL DEFAULT '{}'::jsonb,
    creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS saas_catalogos (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES saas_tenants(id) ON DELETE CASCADE,
    nombre_catalogo VARCHAR(255) NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT true,
    items JSONB NOT NULL DEFAULT '[]'::jsonb,
    actualizado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índices de consulta SaaS
CREATE INDEX IF NOT EXISTS idx_saas_tenants_slug ON saas_tenants(slug);
CREATE INDEX IF NOT EXISTS idx_saas_tenants_dominio ON saas_tenants(dominio_personalizado);
CREATE INDEX IF NOT EXISTS idx_saas_catalogos_tenant_id ON saas_catalogos(tenant_id);

-- Semillas iniciales SaaS (idempotentes)
INSERT INTO saas_tenants (id, nombre, slug, dominio_personalizado, contacto_email, activo, configuracion)
VALUES 
    (1, 'Distribuidora Guayas Cosmetics', 'guayas-cosmetics', 'guayas.multicatalogo.local', 'contacto@guayascosmetics.ec', true, '{"moneda": "USD", "tema_color": "#0ea5e9", "comision_extra_pct": 2.5}'::jsonb),
    (2, 'Santa Elena Beauty Red', 'santa-elena-beauty', NULL, 'admin@upse.edu.ec', true, '{"moneda": "USD", "tema_color": "#4f46e5", "comision_extra_pct": 5.0}'::jsonb)
ON CONFLICT (id) DO NOTHING;

INSERT INTO saas_catalogos (id, tenant_id, nombre_catalogo, activo, items)
VALUES
    (1, 1, 'Catálogo Exclusivo Costa', true, '[{"orden": 1, "destacado": true, "producto_id": 1, "precio_personalizado": 42.0}, {"orden": 2, "destacado": true, "producto_id": 2, "precio_personalizado": 30.0}, {"orden": 3, "destacado": true, "producto_id": 8, "precio_personalizado": 120.0}]'::jsonb)
ON CONFLICT (id) DO NOTHING;
