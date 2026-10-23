import { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuth } from "../context/auth";
import { checkHealth } from "../services/healthService";

// NUEVO: Agregamos la función a la interfaz de props
interface SidebarProps {
  collapsed: boolean;
  onClose: () => void;
}

// Definimos las opciones de navegación según el rol (Tema 5)
interface NavItem {
  to: string;
  label: string;
  title: string;
  icon: string; // ruta del icono SVG (stroke)
  soloAdmin?: boolean;
}

const NAV_ITEMS: NavItem[] = [
  {
    to: "/",
    label: "Dashboard",
    title: "Dashboard",
    soloAdmin: true,
    icon: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
  },
  {
    to: "/tienda",
    label: "Tienda",
    title: "Tienda",
    icon: "M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z",
  },
  {
    to: "/catalogo",
    label: "Catálogo",
    title: "Catálogo",
    icon: "M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z",
  },
  {
    to: "/mi-red",
    label: "Mi Red",
    title: "Mi Red",
    soloAdmin: true,
    icon: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z",
  },
  {
    to: "/admin/productos",
    label: "Admin Productos",
    title: "Administración de Productos",
    soloAdmin: true,
    icon: "M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4",
  },
];

const Sidebar = ({ collapsed, onClose }: SidebarProps) => {
  const { user } = useAuth();
  const { pathname } = useLocation();
  const [backendOnline, setBackendOnline] = useState<boolean | null>(null);

  useEffect(() => {
    checkHealth()
      .then((res) => {
        setBackendOnline(res.status === "ok");
      })
      .catch(() => {
        setBackendOnline(false);
      });
  }, []);

  // Filtramos las opciones según el rol del usuario
  const items = NAV_ITEMS.filter((item) => !item.soloAdmin || user?.rol === "admin");

  return (
    <aside
      className={`
        fixed inset-y-0 left-0 z-50 transform bg-red-900 text-white flex flex-col transition-all duration-300 ease-in-out overflow-hidden whitespace-nowrap
        md:relative md:translate-x-0
        ${collapsed ? "-translate-x-full md:w-20" : "translate-x-0 w-64 md:w-64"}
      `}
    >
      <div className={`p-4 md:p-6 text-xl font-bold border-b border-slate-700 flex items-center ${collapsed ? 'md:justify-center' : 'justify-start'} whitespace-nowrap`}>
        <span className="md:hidden">C's Multicatalogo</span>
        <span className="hidden md:inline">{collapsed ? "CMC" : "Carlos' MultiCatálogo"}</span>
      </div>

      <nav className="flex-1 p-3 space-y-2 overflow-y-auto">
        {items.map((item) => {
          // Resaltamos la opción activa según la ruta actual
          const esActivo =
            item.to === "/"
              ? pathname === "/"
              : pathname.startsWith(item.to);

          return (
            <Link
              key={item.to}
              to={item.to}
              onClick={onClose}
              className={`flex items-center gap-3 p-3 rounded transition ${
                esActivo
                  ? "bg-indigo-600 text-white"
                  : "hover:bg-slate-800"
              } ${collapsed ? 'md:justify-center' : ''}`}
              title={item.title}
            >
              <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={item.icon} />
              </svg>
              <span className={`whitespace-nowrap ${collapsed ? 'md:hidden' : ''}`}>
                {item.label}
              </span>
            </Link>
          );
        })}
      </nav>

      <div className="p-3 border-t border-slate-700 text-xs text-slate-300">
        {collapsed ? (
          <div className="flex flex-col items-center gap-1.5">
            <p className="text-center uppercase">{user?.rol}</p>
            <span
              className={`w-2.5 h-2.5 rounded-full ${
                backendOnline === true
                  ? "bg-green-400"
                  : backendOnline === false
                  ? "bg-red-400"
                  : "bg-amber-400"
              }`}
              title={
                backendOnline === true
                  ? "API Go: Conectada"
                  : backendOnline === false
                  ? "API Go: Desconectada"
                  : "API Go: Verificando..."
              }
            />
          </div>
        ) : (
          <div>
            <p>
              Conectado como <span className="font-semibold uppercase">{user?.rol}</span>
            </p>
            <div className="flex items-center gap-2 mt-1.5 text-[11px] text-slate-300">
              <span
                className={`w-2 h-2 rounded-full flex-shrink-0 ${
                  backendOnline === true
                    ? "bg-green-400"
                    : backendOnline === false
                    ? "bg-red-400"
                    : "bg-amber-400"
                }`}
              />
              <span>
                {backendOnline === true
                  ? "API Go: Conectada"
                  : backendOnline === false
                  ? "API Go: Desconectada"
                  : "API Go: Verificando..."}
              </span>
            </div>
          </div>
        )}
      </div>
    </aside>
  );
};

export default Sidebar;
