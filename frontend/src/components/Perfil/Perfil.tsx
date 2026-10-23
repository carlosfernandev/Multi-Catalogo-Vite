import { useState, useEffect } from "react";
import { obtenerPerfil, type UsuarioPerfil } from "./services/usuarioService";

export default function Perfil() {
  const [perfil, setPerfil] = useState<UsuarioPerfil | null>(null);
  const [cargando, setCargando] = useState(true);
  const [copiado, setCopiado] = useState(false);

  useEffect(() => {
    obtenerPerfil()
      .then(setPerfil)
      .catch((err) => console.error("Error al cargar perfil:", err))
      .finally(() => setCargando(false));
  }, []);

  const copiarCodigo = () => {
    if (perfil?.referido_codigo) {
      const url = `${window.location.origin}/registro?ref=${perfil.referido_codigo}`;
      navigator.clipboard.writeText(url);
      setCopiado(true);
      setTimeout(() => setCopiado(false), 2000);
    }
  };

  if (cargando) {
    return (
      <div className="p-8 text-center text-slate-500">
        <div className="inline-block animate-spin mr-2">⏳</div>
        Cargando perfil...
      </div>
    );
  }

  if (!perfil) {
    return (
      <div className="p-8 text-center text-slate-500">
        No se pudo cargar el perfil del usuario.
      </div>
    );
  }

  return (
    <div className="p-4 md:p-8 max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl md:text-3xl font-bold text-slate-800">Mi Perfil</h1>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6 space-y-4">
        <div>
          <p className="text-sm font-semibold text-slate-400 uppercase">Nombre</p>
          <p className="text-lg font-medium text-slate-800">{perfil.nombre}</p>
        </div>

        <div>
          <p className="text-sm font-semibold text-slate-400 uppercase">Email</p>
          <p className="text-lg font-medium text-slate-800">{perfil.email}</p>
        </div>

        <div>
          <p className="text-sm font-semibold text-slate-400 uppercase">Rol</p>
          <p className="text-lg font-medium text-slate-800 capitalize">{perfil.rol}</p>
        </div>

        {perfil.referido_codigo && (
          <div className="bg-indigo-50 border border-indigo-200 rounded-lg p-4 space-y-3">
            <div>
              <p className="text-sm font-semibold text-indigo-700 uppercase">Tu Código de Referido</p>
              <p className="text-2xl font-bold text-indigo-900 tracking-wider">{perfil.referido_codigo}</p>
            </div>

            <div>
              <p className="text-sm font-semibold text-indigo-700 uppercase mb-2">Enlace de Invitación</p>
              <div className="flex gap-2">
                <input
                  type="text"
                  readOnly
                  value={`${window.location.origin}/registro?ref=${perfil.referido_codigo}`}
                  className="flex-1 px-3 py-2 bg-white border border-indigo-300 rounded-lg text-sm text-slate-600"
                />
                <button
                  onClick={copiarCodigo}
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg transition text-sm"
                >
                  {copiado ? "¡Copiado!" : "Copiar"}
                </button>
              </div>
            </div>

            {perfil.referido_nivel !== undefined && perfil.referido_ventas !== undefined && (
              <div className="grid grid-cols-2 gap-4 pt-2 border-t border-indigo-200">
                <div>
                  <p className="text-xs text-indigo-600">Nivel en la red</p>
                  <p className="text-lg font-bold text-indigo-900">{perfil.referido_nivel}</p>
                </div>
                <div>
                  <p className="text-xs text-indigo-600">Ventas acumuladas</p>
                  <p className="text-lg font-bold text-indigo-900">${perfil.referido_ventas}</p>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
