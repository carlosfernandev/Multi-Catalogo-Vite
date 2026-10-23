import { useState, useEffect, useCallback, useTransition } from 'react';
import type { IProducto } from '../Catalogo/interfaces/iproducto';
import {
  obtenerProductosAdmin,
  crearProducto,
  actualizarProducto,
  eliminarProducto,
  toggleActivoProducto,
  type ProductoInput,
} from '../../services/productoService';
import { ApiError } from '../../services/http';

const CATEGORIAS_DEFAULT = ['Serum', 'Crema', 'Tónico', 'Mascarilla', 'Kit', 'General'];

export default function AdminProductos() {
  const [productos, setProductos] = useState<IProducto[]>([]);
  const [busqueda, setBusqueda] = useState('');
  const [cargando, setCargando] = useState(true);
  const [guardando, setGuardando] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [exitoMsg, setExitoMsg] = useState<string | null>(null);

  // Estado del modal de edición / creación
  const [modalAbierto, setModalAbierto] = useState(false);
  const [productoEnEdicion, setProductoEnEdicion] = useState<IProducto | null>(null);

  // Formulario
  const [formNombre, setFormNombre] = useState('');
  const [formDescripcion, setFormDescripcion] = useState('');
  const [formPrecio, setFormPrecio] = useState<string>('0');
  const [formCategoria, setFormCategoria] = useState('Serum');
  const [formStock, setFormStock] = useState<string>('100');
  const [formImg, setFormImg] = useState('');
  const [formGaleria, setFormGaleria] = useState('');

  const [, startTransition] = useTransition();

  const cargarProductos = useCallback(async (q?: string) => {
    try {
      setCargando(true);
      setErrorMsg(null);
      const data = await obtenerProductosAdmin(q);
      setProductos(data || []);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg('Error al cargar la lista de productos');
      }
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    let activo = true;
    obtenerProductosAdmin()
      .then((data) => {
        if (activo) {
          setProductos(data || []);
          setCargando(false);
        }
      })
      .catch((err) => {
        if (activo) {
          setErrorMsg(err instanceof ApiError ? err.message : 'Error al cargar la lista de productos');
          setCargando(false);
        }
      });
    return () => {
      activo = false;
    };
  }, []);

  const handleBuscar = (e: React.FormEvent) => {
    e.preventDefault();
    cargarProductos(busqueda);
  };

  const abrirModalCrear = () => {
    setProductoEnEdicion(null);
    setFormNombre('');
    setFormDescripcion('');
    setFormPrecio('25.00');
    setFormCategoria('Serum');
    setFormStock('50');
    setFormImg('https://picsum.photos/seed/nuevo/600');
    setFormGaleria('');
    setErrorMsg(null);
    setModalAbierto(true);
  };

  const abrirModalEditar = (prod: IProducto) => {
    setProductoEnEdicion(prod);
    setFormNombre(prod.nombre);
    setFormDescripcion(prod.descripcion || '');
    setFormPrecio(prod.precio.toString());
    setFormCategoria(prod.categoria || 'Serum');
    setFormStock((prod.stock ?? 100).toString());
    setFormImg(prod.img || '');
    setFormGaleria(prod.galeria ? prod.galeria.join('\n') : '');
    setErrorMsg(null);
    setModalAbierto(true);
  };

  const cerrarModal = () => {
    setModalAbierto(false);
    setProductoEnEdicion(null);
  };

  const handleGuardar = async (e: React.FormEvent) => {
    e.preventDefault();
    const precioNum = parseFloat(formPrecio);
    const stockNum = parseInt(formStock, 10);

    if (!formNombre.trim()) {
      setErrorMsg('El nombre es obligatorio');
      return;
    }
    if (isNaN(precioNum) || precioNum <= 0) {
      setErrorMsg('El precio debe ser un número mayor a 0');
      return;
    }
    if (isNaN(stockNum) || stockNum < 0) {
      setErrorMsg('El stock debe ser 0 o superior');
      return;
    }

    const galeriaArray = formGaleria
      .split('\n')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    const payload: ProductoInput = {
      nombre: formNombre.trim(),
      descripcion: formDescripcion.trim(),
      precio: precioNum,
      categoria: formCategoria.trim(),
      stock: stockNum,
      img: formImg.trim() || 'https://picsum.photos/seed/producto/600',
      galeria: galeriaArray,
    };

    try {
      setGuardando(true);
      setErrorMsg(null);

      if (productoEnEdicion) {
        await actualizarProducto(productoEnEdicion.id, payload);
        setExitoMsg(`Producto "${payload.nombre}" actualizado correctamente`);
      } else {
        await crearProducto(payload);
        setExitoMsg(`Producto "${payload.nombre}" creado exitosamente`);
      }

      cerrarModal();
      startTransition(() => {
        cargarProductos(busqueda);
      });
      setTimeout(() => setExitoMsg(null), 4000);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg('Error al guardar el producto');
      }
    } finally {
      setGuardando(false);
    }
  };

  const handleToggleActivo = async (prod: IProducto) => {
    const accion = prod.activo ? 'desactivar' : 'activar';
    const confirma = window.confirm(`¿Estás seguro de que deseas ${accion} el producto "${prod.nombre}"?`);
    if (!confirma) return;

    try {
      if (prod.activo) {
        await eliminarProducto(prod.id);
        setExitoMsg(`Producto "${prod.nombre}" desactivado del catálogo`);
      } else {
        await toggleActivoProducto(prod.id, false);
        setExitoMsg(`Producto "${prod.nombre}" activado nuevamente`);
      }
      cargarProductos(busqueda);
      setTimeout(() => setExitoMsg(null), 3000);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg('Error al cambiar el estado del producto');
      }
    }
  };

  return (
    <div className="p-4 md:p-8 max-w-7xl mx-auto space-y-6">
      {/* Encabezado */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200 pb-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold text-slate-800">
            Administración de Productos
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Gestiona el catálogo, inventario y visibilidad de productos del MLM.
          </p>
        </div>
        <button
          onClick={abrirModalCrear}
          className="inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg shadow-sm transition min-h-[44px]"
        >
          <span className="text-lg leading-none">+</span>
          <span>Nuevo Producto</span>
        </button>
      </div>

      {/* Mensajes de Alerta */}
      {errorMsg && (
        <div className="p-4 bg-red-50 border-l-4 border-red-500 text-red-700 rounded-r-md flex justify-between items-center">
          <p className="text-sm font-medium">{errorMsg}</p>
          <button
            onClick={() => setErrorMsg(null)}
            className="text-red-500 hover:text-red-700 font-bold ml-4"
          >
            ✕
          </button>
        </div>
      )}

      {exitoMsg && (
        <div className="p-4 bg-green-50 border-l-4 border-green-500 text-green-700 rounded-r-md flex justify-between items-center">
          <p className="text-sm font-medium">{exitoMsg}</p>
          <button
            onClick={() => setExitoMsg(null)}
            className="text-green-500 hover:text-green-700 font-bold ml-4"
          >
            ✕
          </button>
        </div>
      )}

      {/* Barra de Búsqueda y Filtros */}
      <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
        <form onSubmit={handleBuscar} className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <input
              type="text"
              value={busqueda}
              onChange={(e) => setBusqueda(e.target.value)}
              placeholder="Buscar por nombre, categoría o descripción..."
              className="w-full pl-10 pr-4 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 text-sm"
            />
            <span className="absolute left-3 top-2.5 text-slate-400">🔍</span>
          </div>
          <button
            type="submit"
            className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white text-sm font-medium rounded-lg transition"
          >
            Buscar
          </button>
          {busqueda && (
            <button
              type="button"
              onClick={() => {
                setBusqueda('');
                cargarProductos('');
              }}
              className="px-3 py-2 bg-slate-100 hover:bg-slate-200 text-slate-600 text-sm font-medium rounded-lg transition"
            >
              Limpiar
            </button>
          )}
        </form>
      </div>

      {/* Tabla de Productos */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-sm">
            <thead className="bg-slate-50 text-slate-600 font-semibold border-b border-slate-200">
              <tr>
                <th className="p-4">Producto</th>
                <th className="p-4">Categoría</th>
                <th className="p-4 text-right">Precio</th>
                <th className="p-4 text-center">Stock</th>
                <th className="p-4 text-center">Estado</th>
                <th className="p-4 text-right">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {cargando ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    <div className="inline-block animate-spin mr-2">⏳</div>
                    Cargando catálogo...
                  </td>
                </tr>
              ) : productos.length === 0 ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    No se encontraron productos registrados.
                  </td>
                </tr>
              ) : (
                productos.map((prod) => (
                  <tr
                    key={prod.id}
                    className={`hover:bg-slate-50 transition ${
                      prod.activo === false ? 'opacity-60 bg-slate-50/50' : ''
                    }`}
                  >
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <img
                          src={prod.img || 'https://picsum.photos/seed/placeholder/100'}
                          alt={prod.nombre}
                          className="w-12 h-12 rounded-lg object-cover border border-slate-200 flex-shrink-0"
                        />
                        <div>
                          <p className="font-semibold text-slate-800">{prod.nombre}</p>
                          <p className="text-xs text-slate-400 truncate max-w-xs">
                            {prod.descripcion || 'Sin descripción'}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <span className="inline-block px-2.5 py-1 text-xs font-medium rounded-full bg-slate-100 text-slate-700">
                        {prod.categoria}
                      </span>
                    </td>
                    <td className="p-4 text-right font-medium text-slate-800">
                      ${prod.precio.toFixed(2)}
                    </td>
                    <td className="p-4 text-center">
                      <span
                        className={`font-semibold ${
                          (prod.stock ?? 0) <= 5
                            ? 'text-red-600'
                            : (prod.stock ?? 0) <= 20
                            ? 'text-amber-600'
                            : 'text-slate-700'
                        }`}
                      >
                        {prod.stock ?? 100}
                      </span>
                    </td>
                    <td className="p-4 text-center">
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold ${
                          prod.activo !== false
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-800'
                        }`}
                      >
                        {prod.activo !== false ? 'Activo' : 'Inactivo'}
                      </span>
                    </td>
                    <td className="p-4 text-right">
                      <div className="inline-flex items-center gap-2">
                        <button
                          onClick={() => abrirModalEditar(prod)}
                          className="px-2.5 py-1.5 text-xs font-medium bg-slate-100 hover:bg-slate-200 text-slate-700 rounded transition"
                        >
                          Editar
                        </button>
                        <button
                          onClick={() => handleToggleActivo(prod)}
                          className={`px-2.5 py-1.5 text-xs font-medium rounded transition ${
                            prod.activo !== false
                              ? 'bg-red-50 hover:bg-red-100 text-red-600'
                              : 'bg-green-50 hover:bg-green-100 text-green-700'
                          }`}
                        >
                          {prod.activo !== false ? 'Desactivar' : 'Activar'}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal Crear / Editar */}
      {modalAbierto && (
        <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg my-8 p-6 space-y-5">
            <div className="flex justify-between items-center border-b border-slate-100 pb-3">
              <h2 className="text-xl font-bold text-slate-800">
                {productoEnEdicion ? 'Editar Producto' : 'Crear Nuevo Producto'}
              </h2>
              <button
                onClick={cerrarModal}
                className="text-slate-400 hover:text-slate-600 font-bold p-1"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleGuardar} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                  Nombre del Producto *
                </label>
                <input
                  type="text"
                  required
                  value={formNombre}
                  onChange={(e) => setFormNombre(e.target.value)}
                  placeholder="Ej: Crema Reparadora Nocturna"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                    Precio ($) *
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    min="0.01"
                    required
                    value={formPrecio}
                    onChange={(e) => setFormPrecio(e.target.value)}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                    Stock *
                  </label>
                  <input
                    type="number"
                    min="0"
                    required
                    value={formStock}
                    onChange={(e) => setFormStock(e.target.value)}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                    Categoría *
                  </label>
                  <select
                    value={formCategoria}
                    onChange={(e) => setFormCategoria(e.target.value)}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm bg-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  >
                    {CATEGORIAS_DEFAULT.map((cat) => (
                      <option key={cat} value={cat}>
                        {cat}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                  Descripción
                </label>
                <textarea
                  rows={3}
                  value={formDescripcion}
                  onChange={(e) => setFormDescripcion(e.target.value)}
                  placeholder="Detalles sobre beneficios, ingredientes o aplicación..."
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                  URL de Imagen Principal
                </label>
                <input
                  type="url"
                  value={formImg}
                  onChange={(e) => setFormImg(e.target.value)}
                  placeholder="https://picsum.photos/seed/ejemplo/600"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                  Galería de Imágenes (una URL por línea)
                </label>
                <textarea
                  rows={2}
                  value={formGaleria}
                  onChange={(e) => setFormGaleria(e.target.value)}
                  placeholder="https://picsum.photos/seed/1/600&#10;https://picsum.photos/seed/2/600"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono text-xs focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
              </div>

              <div className="flex justify-end gap-3 pt-3 border-t border-slate-100">
                <button
                  type="button"
                  onClick={cerrarModal}
                  disabled={guardando}
                  className="px-4 py-2 border border-slate-300 text-slate-600 font-medium rounded-lg hover:bg-slate-50 transition text-sm"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={guardando}
                  className="px-5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg shadow-sm transition text-sm disabled:opacity-50"
                >
                  {guardando ? 'Guardando...' : productoEnEdicion ? 'Actualizar' : 'Crear Producto'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
