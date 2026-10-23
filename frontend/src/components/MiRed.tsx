// src/components/MiRed.tsx
// Red multinivel visual (Tema 5):
// - Resumen: referidos activos, ventas de la red, comisiones del mes y nivel alcanzado
// - Árbol jerárquico recursivo de referidos con comisión calculada por nivel
// - Historial de liquidaciones reales registradas por pedidos en backend
import { useState, useEffect } from 'react';
import { contarRed, sumarVentasRed, sumarComisiones, nivelAlcanzado, comisionDeReferido, TASA_COMISION, type Referido } from './MiRed/data/red';
import { getRed } from './MiRed/services/redService';
import { getLiquidaciones, type LiquidacionRegistro } from '../services/liquidacionService';
import { useAuth } from '../context/auth';

// Colores y etiquetas por nivel para el árbol
const NIVEL_ESTILO: Record<number, { badge: string; border: string; tasa: string }> = {
  1: { badge: "bg-indigo-100 text-indigo-700", border: "border-indigo-200", tasa: "10 %" },
  2: { badge: "bg-purple-100 text-purple-700", border: "border-purple-200", tasa: "5 %" },
  3: { badge: "bg-pink-100 text-pink-700", border: "border-pink-200", tasa: "2 %" },
};

// Componente recursivo: renderiza un referido y, debajo, a sus hijos
// Nota: la prop se llama "nodo" porque "ref" es una prop reservada en React 19
const NodoReferido = ({ nodo }: { nodo: Referido }) => {
  const estilo = NIVEL_ESTILO[nodo.nivel] ?? NIVEL_ESTILO[1];
  const comision = comisionDeReferido(nodo);

  return (
    <li>
      <div className={`bg-white rounded-lg border ${estilo.border} shadow-sm p-4 flex flex-wrap items-center justify-between gap-3`}>
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-600">
            {nodo.nombre.charAt(0)}
          </div>
          <div>
            <p className="font-semibold text-slate-800">{nodo.nombre}</p>
            <span className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${estilo.badge}`}>
              Nivel {nodo.nivel} · {estilo.tasa} comisión
            </span>
          </div>
        </div>
        <div className="text-right">
          <p className="text-sm text-slate-500">Ventas</p>
          <p className="font-bold text-slate-800">${nodo.ventas.toLocaleString()}</p>
          <p className="text-xs text-indigo-600 font-semibold mt-1">
            +${comision.toFixed(2)} para ti
          </p>
        </div>
      </div>

      {/* Hijos del referido (más profundidad en el árbol) */}
      {nodo.hijos && nodo.hijos.length > 0 && (
        <ul className="mt-3 ml-2 sm:ml-6 space-y-3 border-l-2 border-slate-200 pl-2 sm:pl-4">
          {nodo.hijos.map((hijo) => (
            <NodoReferido key={hijo.id} nodo={hijo} />
          ))}
        </ul>
      )}
    </li>
  );
};

const MiRed = () => {
  const { user } = useAuth();
  const email = user?.email;
  const [red, setRed] = useState<Referido | null>(null);
  const [liquidaciones, setLiquidaciones] = useState<LiquidacionRegistro[]>([]);
  const [totalLiquidado, setTotalLiquidado] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([getRed(email), getLiquidaciones().catch(() => ({ total_comisiones: 0, cantidad: 0, liquidaciones: [] }))])
      .then(([dataRed, dataLiq]) => {
        setRed(dataRed);
        setLiquidaciones(dataLiq.liquidaciones || []);
        setTotalLiquidado(dataLiq.total_comisiones || 0);
        setLoading(false);
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : "Error al cargar la red");
        setLoading(false);
      });
  }, [email]);

  if (loading) {
    return (
      <div className="flex items-center justify-center p-12">
        <p className="text-slate-500 font-medium">Cargando red de referidos...</p>
      </div>
    );
  }

  if (error || !red) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 p-6 rounded-lg">
        <p className="font-semibold">Error al cargar la red</p>
        <p className="text-sm mt-1">{error ?? "No se encontraron datos de la red."}</p>
      </div>
    );
  }

  // KPIs derivados de los datos de la red (funciones puras en data/red.ts)
  const referidosActivos = contarRed(red);
  const ventasRed = sumarVentasRed(red);
  const comisionesMes = sumarComisiones(red);
  const nivel = nivelAlcanzado(red.hijos?.length ?? 0);

  const resumen = [
    { etiqueta: "Referidos Activos", valor: referidosActivos.toString(), color: "text-indigo-600" },
    { etiqueta: "Ventas de la Red", valor: `$${ventasRed.toLocaleString()}`, color: "text-indigo-600" },
    { etiqueta: "Comisiones del Mes", valor: `$${comisionesMes.toFixed(2)}`, color: "text-green-600" },
    { etiqueta: "Nivel Alcanzado", valor: nivel, color: "text-amber-600" },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-2">Mi Red de Referidos</h1>
      <p className="text-slate-500 mb-6">
        Estructura multinivel: gana el {TASA_COMISION[1] * 100} % de nivel 1, {TASA_COMISION[2] * 100} % de nivel 2 y {TASA_COMISION[3] * 100} % de nivel 3.
      </p>

      {/* ===== TARJETAS DE RESUMEN ===== */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {resumen.map((item) => (
          <div key={item.etiqueta} className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
            <p className="text-sm text-slate-500 uppercase font-semibold">{item.etiqueta}</p>
            <p className={`text-3xl font-bold mt-2 ${item.color}`}>{item.valor}</p>
          </div>
        ))}
      </div>

      {/* ===== ÁRBOL DE REFERIDOS ===== */}
      <div className="bg-white rounded-lg border border-slate-200 shadow-sm p-6">
        <h2 className="text-lg font-bold text-slate-800 mb-4">Jerarquía de tu Red</h2>

        {/* Raíz: el usuario */}
        <div className="bg-slate-900 text-white rounded-lg p-4 mb-4 flex flex-wrap items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-indigo-500 flex items-center justify-center font-bold">
              {red.nombre.charAt(0)}
            </div>
            <div>
              <p className="font-semibold">{red.nombre}</p>
              <p className="text-xs text-slate-300">Nivel alcanzado: {nivel}</p>
            </div>
          </div>
          <p className="font-bold text-green-400">
            Comisiones del mes: ${comisionesMes.toFixed(2)}
          </p>
        </div>

        <ul className="space-y-3">
          {red.hijos?.map((hijo) => (
            <NodoReferido key={hijo.id} nodo={hijo} />
          ))}
        </ul>
      </div>

      {/* ===== HISTORIAL DE LIQUIDACIONES / COMISIONES REALES ===== */}
      <div className="bg-white rounded-lg border border-slate-200 shadow-sm p-6 mt-8">
        <div className="flex flex-wrap items-center justify-between gap-4 mb-4">
          <div>
            <h2 className="text-lg font-bold text-slate-800">Historial de Liquidaciones de Comisiones</h2>
            <p className="text-sm text-slate-500">
              Comisiones acumuladas e ingresadas por pedidos reales procesados por la API.
            </p>
          </div>
          <div className="text-right">
            <span className="text-xs uppercase font-semibold text-slate-400">Total liquidado registrado</span>
            <p className="text-xl font-bold text-green-600">${totalLiquidado.toFixed(2)}</p>
          </div>
        </div>

        {liquidaciones.length === 0 ? (
          <div className="p-8 text-center bg-slate-50 rounded-lg border border-dashed border-slate-200">
            <p className="text-slate-500 font-medium">Aún no hay liquidaciones registradas en la base de datos.</p>
            <p className="text-xs text-slate-400 mt-1">
              Las comisiones se generan y registran automáticamente al completar pedidos en el Checkout.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-slate-700 text-xs uppercase font-semibold border-b border-slate-200">
                <tr>
                  <th className="p-3">Pedido</th>
                  <th className="p-3">Beneficiario</th>
                  <th className="p-3">Nivel</th>
                  <th className="p-3">Monto Comisión</th>
                  <th className="p-3">Fecha</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {liquidaciones.map((liq) => (
                  <tr key={liq.id} className="hover:bg-slate-50 transition">
                    <td className="p-3 font-semibold text-indigo-600">{liq.pedido_numero}</td>
                    <td className="p-3 text-slate-800 font-medium">{liq.referido_nombre}</td>
                    <td className="p-3">
                      <span className={`inline-block px-2 py-0.5 text-xs font-semibold rounded-full ${NIVEL_ESTILO[liq.nivel]?.badge ?? "bg-slate-100"}`}>
                        Nivel {liq.nivel}
                      </span>
                    </td>
                    <td className="p-3 font-bold text-green-600">+${liq.monto.toFixed(2)}</td>
                    <td className="p-3 text-slate-400 text-xs">
                      {new Date(liq.fecha).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default MiRed;