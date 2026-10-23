import { useState, useEffect, useCallback, useTransition } from 'react';
import {
  obtenerUsuariosAdmin,
  actualizarUsuarioAdmin,
  toggleActivoUsuario,
  type UsuarioAdmin,
} from '../../services/usuarioService';
import { ApiError } from '../../services/http';
import { useAuth } from '../../context/auth';

export default function AdminUsuarios() {
  const { user: authUser } = useAuth();
  const [usuarios, setUsuarios] = useState<UsuarioAdmin[]>([]);
  const [busqueda, setBusqueda] = useState('');
  const [cargando, setCargando] = useState(true);
  const [guardando, setGuardando] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [exitoMsg, setExitoMsg] = useState<string | null>(null);

  // Modal para edición
  const [modalAbierto, setModalAbierto] = useState(false);
  const [usuarioEnEdicion, setUsuarioEnEdicion] = useState<UsuarioAdmin | null>(null);
  const [formNombre, setFormNombre] = useState('');
  const [formRol, setFormRol] = useState<'admin' | 'cliente'>('cliente');
  const [formActivo, setFormActivo] = useState(true);

  const [, startTransition] = useTransition();

  const cargarUsuarios = useCallback(async (q?: string) => {
    try {
      setCargando(true);
      setErrorMsg(null);
      const data = await obtenerUsuariosAdmin(q);
      setUsuarios(data || []);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg('Error al cargar la lista de usuarios');
      }
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    let activo = true;
    obtenerUsuariosAdmin()
      .then((data) => {
        if (activo) {
          setUsuarios(data || []);
          setCargando(false);
        }
      })
      .catch((err) => {
        if (activo) {
          setErrorMsg(err instanceof ApiError ? err.message : 'Error al cargar usuarios');
          setCargando(false);
        }
      });
    return () => {
      activo = false;
    };
  }, []);

  const handleBuscar = (e: React.FormEvent) => {
    e.preventDefault();
    cargarUsuarios(busqueda);
  };

  const abrirModalEditar = (u: UsuarioAdmin) => {
    setUsuarioEnEdicion(u);
    setFormNombre(u.nombre || '');
    setFormRol(u.rol);
    setFormActivo(u.activo);
    setErrorMsg(null);
    setModalAbierto(true);
  };

  const cerrarModal = () => {
    setModalAbierto(false);
    setUsuarioEnEdicion(null);
  };

  const handleGuardar = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!usuarioEnEdicion) return;

    try {
      setGuardando(true);
      setErrorMsg(null);

      await actualizarUsuarioAdmin(usuarioEnEdicion.id, {
        nombre: formNombre.trim(),
        rol: formRol,
        activo: formActivo,
      });

      setExitoMsg(`Usuario "${usuarioEnEdicion.email}" actualizado exitosamente`);
      cerrarModal();

      startTransition(() => {
        cargarUsuarios(busqueda);
      });
      setTimeout(() => setExitoMsg(null), 4000);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg('Error al actualizar el usuario');
      }
    } finally {
      setGuardando(false);
    }
  };

  const handleToggleActivo = async (u: UsuarioAdmin) => {
    if (u.email === authUser?.email && u.activo) {
      setErrorMsg('No puedes desactivar tu propia cuenta de administrador');
      return;
    }

    const accion = u.activo ? 'desactivar' : 'activar';
    const confirma = window.confirm(`¿Deseas ${accion} la cuenta de "${u.email}"?`);
    if (!confirma) return;

    try {
      await toggleActivoUsuario(u.id, u.activo);
      setExitoMsg(`Cuenta de "${u.email}" ${accion === 'activar' ? 'activada' : 'desactivada'}`);
      cargarUsuarios(busqueda);
      setTimeout(() => setExitoMsg(null), 3000);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg('Error al cambiar el estado del usuario');
      }
    }
  };

  const formatearFecha = (fechaStr?: string) => {
    if (!fechaStr) return 'N/A';
    try {
      return new Date(fechaStr).toLocaleDateString('es-ES', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch {
      return fechaStr;
    }
  };

  // Métricas
  const totalUsuarios = usuarios.length;
  const totalAdmins = usuarios.filter((u) => u.rol === 'admin').length;
  const totalClientes = usuarios.filter((u) => u.rol === 'cliente').length;
  const totalActivos = usuarios.filter((u) => u.activo).length;

  return (
    <div className="p-4 md:p-8 max-w-7xl mx-auto space-y-6">
      {/* Encabezado */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200 pb-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold text-slate-800">
            Administración de Usuarios
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Gestiona roles, estados de cuenta y permisos de acceso al sistema.
          </p>
        </div>
      </div>

      {/* Tarjetas de Métricas */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Total Cuentas</p>
          <p className="text-2xl font-bold text-slate-800 mt-1">{totalUsuarios}</p>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Administradores</p>
          <p className="text-2xl font-bold text-amber-600 mt-1">{totalAdmins}</p>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Clientes / Distribuidores</p>
          <p className="text-2xl font-bold text-indigo-600 mt-1">{totalClientes}</p>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Cuentas Activas</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{totalActivos}</p>
        </div>
      </div>

      {/* Alertas */}
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

      {/* Buscador */}
      <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
        <form onSubmit={handleBuscar} className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <input
              type="text"
              value={busqueda}
              onChange={(e) => setBusqueda(e.target.value)}
              placeholder="Buscar por correo electrónico o nombre..."
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
                cargarUsuarios('');
              }}
              className="px-3 py-2 bg-slate-100 hover:bg-slate-200 text-slate-600 text-sm font-medium rounded-lg transition"
            >
              Limpiar
            </button>
          )}
        </form>
      </div>

      {/* Tabla de Usuarios */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-sm">
            <thead className="bg-slate-50 text-slate-600 font-semibold border-b border-slate-200">
              <tr>
                <th className="p-4">Usuario</th>
                <th className="p-4">Rol</th>
                <th className="p-4 text-center">Referido ID</th>
                <th className="p-4 text-center">Estado</th>
                <th className="p-4">Fecha Registro</th>
                <th className="p-4 text-right">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {cargando ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    <div className="inline-block animate-spin mr-2">⏳</div>
                    Cargando usuarios...
                  </td>
                </tr>
              ) : usuarios.length === 0 ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    No se encontraron usuarios que coincidan con la búsqueda.
                  </td>
                </tr>
              ) : (
                usuarios.map((u) => {
                  const esUsuarioActual = u.email === authUser?.email;
                  return (
                    <tr
                      key={u.id}
                      className={`hover:bg-slate-50 transition ${
                        !u.activo ? 'opacity-60 bg-slate-50/50' : ''
                      }`}
                    >
                      <td className="p-4">
                        <div className="flex items-center gap-3">
                          <div className="w-9 h-9 rounded-full bg-slate-200 text-slate-600 font-bold flex items-center justify-center flex-shrink-0 text-sm uppercase">
                            {u.nombre ? u.nombre.charAt(0) : u.email.charAt(0)}
                          </div>
                          <div>
                            <p className="font-semibold text-slate-800 flex items-center gap-1.5">
                              {u.nombre || 'Sin nombre'}
                              {esUsuarioActual && (
                                <span className="text-[10px] bg-slate-100 text-slate-600 px-1.5 py-0.5 rounded font-normal">
                                  Tú
                                </span>
                              )}
                            </p>
                            <p className="text-xs text-slate-400">{u.email}</p>
                          </div>
                        </div>
                      </td>
                      <td className="p-4">
                        <span
                          className={`inline-block px-2.5 py-1 text-xs font-semibold rounded-full uppercase ${
                            u.rol === 'admin'
                              ? 'bg-amber-100 text-amber-700'
                              : 'bg-indigo-100 text-indigo-700'
                          }`}
                        >
                          {u.rol}
                        </span>
                      </td>
                      <td className="p-4 text-center font-mono text-xs text-slate-600">
                        {u.referido_id != null ? `#${u.referido_id}` : '—'}
                      </td>
                      <td className="p-4 text-center">
                        <span
                          className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold ${
                            u.activo
                              ? 'bg-green-100 text-green-800'
                              : 'bg-red-100 text-red-800'
                          }`}
                        >
                          {u.activo ? 'Activo' : 'Inactivo'}
                        </span>
                      </td>
                      <td className="p-4 text-slate-500 text-xs">
                        {formatearFecha(u.creado_en)}
                      </td>
                      <td className="p-4 text-right">
                        <div className="inline-flex items-center gap-2">
                          <button
                            onClick={() => abrirModalEditar(u)}
                            className="px-2.5 py-1.5 text-xs font-medium bg-slate-100 hover:bg-slate-200 text-slate-700 rounded transition"
                          >
                            Editar
                          </button>
                          <button
                            onClick={() => handleToggleActivo(u)}
                            disabled={esUsuarioActual && u.activo}
                            title={
                              esUsuarioActual && u.activo
                                ? 'No puedes desactivar tu propia cuenta'
                                : undefined
                            }
                            className={`px-2.5 py-1.5 text-xs font-medium rounded transition disabled:opacity-40 disabled:cursor-not-allowed ${
                              u.activo
                                ? 'bg-red-50 hover:bg-red-100 text-red-600'
                                : 'bg-green-50 hover:bg-green-100 text-green-700'
                            }`}
                          >
                            {u.activo ? 'Desactivar' : 'Activar'}
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal Editar Usuario */}
      {modalAbierto && usuarioEnEdicion && (
        <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-md p-6 space-y-5">
            <div className="flex justify-between items-center border-b border-slate-100 pb-3">
              <div>
                <h2 className="text-xl font-bold text-slate-800">Editar Usuario</h2>
                <p className="text-xs text-slate-400 mt-0.5">{usuarioEnEdicion.email}</p>
              </div>
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
                  Nombre Completo
                </label>
                <input
                  type="text"
                  value={formNombre}
                  onChange={(e) => setFormNombre(e.target.value)}
                  placeholder="Ej: Juan Pérez"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                  Rol de Usuario
                </label>
                <select
                  value={formRol}
                  onChange={(e) => setFormRol(e.target.value as 'admin' | 'cliente')}
                  disabled={usuarioEnEdicion.email === authUser?.email}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm bg-white focus:ring-2 focus:ring-indigo-500 focus:outline-none disabled:bg-slate-100 disabled:text-slate-500"
                >
                  <option value="cliente">Cliente (Distribuidor MLM)</option>
                  <option value="admin">Administrador (Acceso Total)</option>
                </select>
                {usuarioEnEdicion.email === authUser?.email && (
                  <p className="text-[11px] text-amber-600 mt-1">
                    No puedes modificar tu propio rol de administrador.
                  </p>
                )}
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-600 uppercase mb-1">
                  Estado de Cuenta
                </label>
                <div className="flex items-center gap-3 pt-1">
                  <label className="inline-flex items-center gap-2 cursor-pointer text-sm">
                    <input
                      type="checkbox"
                      checked={formActivo}
                      onChange={(e) => setFormActivo(e.target.checked)}
                      disabled={usuarioEnEdicion.email === authUser?.email && formActivo}
                      className="w-4 h-4 text-indigo-600 rounded focus:ring-indigo-500"
                    />
                    <span className="text-slate-700">Cuenta Activa (permite inicio de sesión)</span>
                  </label>
                </div>
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
                  {guardando ? 'Guardando...' : 'Guardar Cambios'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
