import { useState, useEffect, useCallback, useTransition } from "react";
import {
  obtenerPendientesLiquidacion,
  obtenerHistorialLiquidaciones,
  procesarLiquidacion,
  type DistribuidorPendiente,
  type LiquidacionHistorial,
} from "../../services/liquidacionService";
import { ApiError } from "../../services/http";

export default function AdminLiquidaciones() {
  const [tab, setTab] = useState<"pendientes" | "historial">("pendientes");
  const [pendientes, setPendientes] = useState<DistribuidorPendiente[]>([]);
  const [historial, setHistorial] = useState<LiquidacionHistorial[]>([]);
  const [cargando, setCargando] = useState(true);
  const [liquidandoId, setLiquidandoId] = useState<number | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [exitoMsg, setExitoMsg] = useState<string | null>(null);
  const [mostrarGuia, setMostrarGuia] = useState(true);

  const [, startTransition] = useTransition();

  const cargarDatos = useCallback(async () => {
    try {
      setCargando(true);
      setErrorMsg(null);
      const [pend, hist] = await Promise.all([
        obtenerPendientesLiquidacion(),
        obtenerHistorialLiquidaciones(),
      ]);
      setPendientes(pend || []);
      setHistorial(hist || []);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg("Error al sincronizar liquidaciones");
      }
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    let activo = true;
    Promise.all([obtenerPendientesLiquidacion(), obtenerHistorialLiquidaciones()])
      .then(([pend, hist]) => {
        if (activo) {
          setPendientes(pend || []);
          setHistorial(hist || []);
          setCargando(false);
        }
      })
      .catch((err) => {
        if (activo) {
          setErrorMsg(err instanceof ApiError ? err.message : "Error al cargar datos");
          setCargando(false);
        }
      });
    return () => {
      activo = false;
    };
  }, []);

  const handleLiquidar = async (item: DistribuidorPendiente) => {
    const confirma = window.confirm(
      `¿Confirmas la liquidación de $${item.total_pendiente.toFixed(2)} a favor de "${item.referido_nombre}"?`
    );
    if (!confirma) return;

    try {
      setLiquidandoId(item.referido_id);
      setErrorMsg(null);
      const res = await procesarLiquidacion(item.referido_id);

      setExitoMsg(
        `Liquidación #LIQ-${res.id} por $${res.total.toFixed(2)} emitida con éxito a "${item.referido_nombre}".`
      );

      startTransition(() => {
        cargarDatos();
      });
      setTimeout(() => setExitoMsg(null), 5000);
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message);
      } else {
        setErrorMsg("Error al procesar la liquidación");
      }
    } finally {
      setLiquidandoId(null);
    }
  };

  const formatearFecha = (fechaStr?: string) => {
    if (!fechaStr) return "N/A";
    try {
      return new Date(fechaStr).toLocaleDateString("es-ES", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    } catch {
      return fechaStr;
    }
  };

  // Cálculos de métricas
  const deudaPendienteTotal = pendientes.reduce((acc, p) => acc + p.total_pendiente, 0);
  const totalPagadoHistorico = historial.reduce((acc, h) => acc + h.total, 0);

  return (
    <div className="p-4 md:p-8 max-w-7xl mx-auto space-y-6">
      {/* Encabezado */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200 pb-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold text-slate-800">
            Gestión de Liquidaciones MLM
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Control de comisiones devengadas, balances pendientes y pagos formalizados a distribuidores.
          </p>
        </div>
        <button
          onClick={() => setMostrarGuia((prev) => !prev)}
          className="inline-flex items-center gap-2 text-xs font-semibold px-3 py-2 bg-indigo-50 text-indigo-700 hover:bg-indigo-100 rounded-lg transition"
        >
          <span>💡</span>
          <span>{mostrarGuia ? "Ocultar Explicación" : "¿Cómo funcionan las liquidaciones?"}</span>
        </button>
      </div>

      {/* Guía Explicativa Didáctica */}
      {mostrarGuia && (
        <div className="bg-gradient-to-r from-indigo-50 via-sky-50 to-emerald-50 border border-indigo-200 p-5 rounded-2xl shadow-sm text-sm text-slate-700 space-y-3">
          <div className="flex items-center justify-between">
            <h3 className="font-bold text-indigo-900 flex items-center gap-2">
              <span>📖</span> ¿Cómo funciona el ciclo de Liquidaciones?
            </h3>
            <button
              onClick={() => setMostrarGuia(false)}
              className="text-slate-400 hover:text-slate-600 font-bold"
            >
              ✕
            </button>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-1">
            <div className="bg-white/80 p-3.5 rounded-xl border border-indigo-100">
              <p className="font-semibold text-indigo-800 text-xs uppercase mb-1">
                1. Compra del Cliente
              </p>
              <p className="text-xs text-slate-600 leading-relaxed">
                Cada vez que un cliente realiza un checkout, el backend calcula automáticamente las
                comisiones ascendentes (<strong>Nivel 1: 10%</strong>, <strong>Nivel 2: 5%</strong>,{" "}
                <strong>Nivel 3: 2%</strong>) y las guarda con estado <code>pendiente</code>.
              </p>
            </div>
            <div className="bg-white/80 p-3.5 rounded-xl border border-indigo-100">
              <p className="font-semibold text-amber-800 text-xs uppercase mb-1">
                2. Balance Pendiente
              </p>
              <p className="text-xs text-slate-600 leading-relaxed">
                En la pestaña <strong>"Pendientes de Cobro"</strong> se acumula la deuda que la empresa
                debe pagar a cada distribuidor por las ventas de su red.
              </p>
            </div>
            <div className="bg-white/80 p-3.5 rounded-xl border border-indigo-100">
              <p className="font-semibold text-emerald-800 text-xs uppercase mb-1">
                3. Liquidación y Pago
              </p>
              <p className="text-xs text-slate-600 leading-relaxed">
                Al pulsar <strong>"Liquidar"</strong>, se ejecuta una transacción ACID que suma el total,
                emite un folio en <code>liquidaciones</code> y marca las comisiones como pagadas. El saldo
                pendiente del distribuidor queda en $0.00.
              </p>
            </div>
          </div>
          <div className="bg-amber-50 border border-amber-200 text-amber-900 text-xs p-3 rounded-lg flex items-center gap-2">
            <span>✨</span>
            <span>
              <strong>¿Para ver más registros?</strong> Realiza pedidos en la tienda con una cuenta de
              cliente vinculada a un referente. Luego ven aquí y pulsa "Liquidar" para generar nuevos
              recibos de pago.
            </span>
          </div>
        </div>
      )}

      {/* Tarjetas de Métricas */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Deuda Pendiente Total</p>
          <p className="text-2xl font-bold text-amber-600 mt-1">
            ${deudaPendienteTotal.toFixed(2)}
          </p>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Total Liquidado (Pagado)</p>
          <p className="text-2xl font-bold text-emerald-600 mt-1">
            ${totalPagadoHistorico.toFixed(2)}
          </p>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Distribuidores x Cobrar</p>
          <p className="text-2xl font-bold text-indigo-600 mt-1">{pendientes.length}</p>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p className="text-xs font-semibold text-slate-400 uppercase">Comprobantes Emitidos</p>
          <p className="text-2xl font-bold text-slate-800 mt-1">{historial.length}</p>
        </div>
      </div>

      {/* Alertas */}
      {errorMsg && (
        <div className="p-4 bg-red-50 border-l-4 border-red-500 text-red-700 rounded-r-md flex justify-between items-center text-sm">
          <p className="font-medium">{errorMsg}</p>
          <button onClick={() => setErrorMsg(null)} className="text-red-500 hover:text-red-700 font-bold ml-4">
            ✕
          </button>
        </div>
      )}

      {exitoMsg && (
        <div className="p-4 bg-green-50 border-l-4 border-green-500 text-green-700 rounded-r-md flex justify-between items-center text-sm">
          <p className="font-medium">{exitoMsg}</p>
          <button onClick={() => setExitoMsg(null)} className="text-green-500 hover:text-green-700 font-bold ml-4">
            ✕
          </button>
        </div>
      )}

      {/* Selector de Pestañas */}
      <div className="flex border-b border-slate-200 gap-4">
        <button
          onClick={() => setTab("pendientes")}
          className={`pb-3 font-semibold text-sm border-b-2 transition flex items-center gap-2 ${
            tab === "pendientes"
              ? "border-indigo-600 text-indigo-600"
              : "border-transparent text-slate-500 hover:text-slate-700"
          }`}
        >
          <span>⏳ Pendientes de Cobro</span>
          {pendientes.length > 0 && (
            <span className="bg-amber-100 text-amber-800 text-xs px-2 py-0.5 rounded-full font-bold">
              {pendientes.length}
            </span>
          )}
        </button>
        <button
          onClick={() => setTab("historial")}
          className={`pb-3 font-semibold text-sm border-b-2 transition flex items-center gap-2 ${
            tab === "historial"
              ? "border-indigo-600 text-indigo-600"
              : "border-transparent text-slate-500 hover:text-slate-700"
          }`}
        >
          <span>📜 Historial de Liquidaciones</span>
          {historial.length > 0 && (
            <span className="bg-slate-100 text-slate-700 text-xs px-2 py-0.5 rounded-full font-bold">
              {historial.length}
            </span>
          )}
        </button>
      </div>

      {/* Contenido Pestaña 1: Pendientes */}
      {tab === "pendientes" && (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse text-sm">
              <thead className="bg-slate-50 text-slate-600 font-semibold border-b border-slate-200">
                <tr>
                  <th className="p-4">Distribuidor</th>
                  <th className="p-4 text-center">Nodo ID</th>
                  <th className="p-4 text-center">Ventas Acreditadas</th>
                  <th className="p-4 text-right">Saldo Pendiente</th>
                  <th className="p-4 text-right">Acción</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {cargando ? (
                  <tr>
                    <td colSpan={5} className="p-8 text-center text-slate-400">
                      <div className="inline-block animate-spin mr-2">⏳</div>
                      Cargando pendientes...
                    </td>
                  </tr>
                ) : pendientes.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="p-8 text-center text-slate-400">
                      🎉 ¡Todo al día! No hay comisiones pendientes por liquidar actualmente.
                    </td>
                  </tr>
                ) : (
                  pendientes.map((item) => (
                    <tr key={item.referido_id} className="hover:bg-slate-50 transition">
                      <td className="p-4">
                        <div className="flex items-center gap-3">
                          <div className="w-9 h-9 rounded-full bg-indigo-100 text-indigo-700 font-bold flex items-center justify-center text-sm">
                            {item.referido_nombre.charAt(0)}
                          </div>
                          <div>
                            <p className="font-semibold text-slate-800">{item.referido_nombre}</p>
                            <p className="text-xs text-slate-400">Distribuidor de la red</p>
                          </div>
                        </div>
                      </td>
                      <td className="p-4 text-center font-mono text-xs text-slate-500">
                        #{item.referido_id}
                      </td>
                      <td className="p-4 text-center">
                        <span className="inline-block px-2 py-0.5 bg-slate-100 text-slate-700 rounded text-xs font-semibold">
                          {item.cantidad_comisiones} comisiones
                        </span>
                      </td>
                      <td className="p-4 text-right">
                        <span className="font-bold text-amber-600 text-base">
                          ${item.total_pendiente.toFixed(2)}
                        </span>
                      </td>
                      <td className="p-4 text-right">
                        <button
                          onClick={() => handleLiquidar(item)}
                          disabled={liquidandoId === item.referido_id}
                          className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-medium rounded-lg shadow-sm transition text-xs disabled:opacity-50 inline-flex items-center gap-1.5"
                        >
                          {liquidandoId === item.referido_id ? (
                            <span>Procesando...</span>
                          ) : (
                            <>
                              <span>💳</span>
                              <span>Liquidar Pago</span>
                            </>
                          )}
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Contenido Pestaña 2: Historial */}
      {tab === "historial" && (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse text-sm">
              <thead className="bg-slate-50 text-slate-600 font-semibold border-b border-slate-200">
                <tr>
                  <th className="p-4">Folio</th>
                  <th className="p-4">Distribuidor</th>
                  <th className="p-4 text-right">Total Liquidado</th>
                  <th className="p-4 text-center">Estado</th>
                  <th className="p-4 text-right">Fecha de Pago</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {cargando ? (
                  <tr>
                    <td colSpan={5} className="p-8 text-center text-slate-400">
                      <div className="inline-block animate-spin mr-2">⏳</div>
                      Cargando historial...
                    </td>
                  </tr>
                ) : historial.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="p-8 text-center text-slate-400">
                      Aún no se han ejecutado liquidaciones en el sistema.
                    </td>
                  </tr>
                ) : (
                  historial.map((h) => (
                    <tr key={h.id} className="hover:bg-slate-50 transition">
                      <td className="p-4 font-mono font-bold text-indigo-600">
                        #LIQ-{h.id.toString().padStart(4, "0")}
                      </td>
                      <td className="p-4">
                        <p className="font-semibold text-slate-800">
                          {h.referido_nombre || `Referido #${h.referido_id}`}
                        </p>
                        <p className="text-xs text-slate-400">Referido ID #{h.referido_id}</p>
                      </td>
                      <td className="p-4 text-right font-bold text-emerald-600 text-base">
                        ${h.total.toFixed(2)}
                      </td>
                      <td className="p-4 text-center">
                        <span className="inline-block px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-800 uppercase">
                          {h.estado}
                        </span>
                      </td>
                      <td className="p-4 text-right text-slate-500 text-xs">
                        {formatearFecha(h.pagada_en || h.creado_en)}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}
