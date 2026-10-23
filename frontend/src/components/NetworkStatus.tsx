import { useState, useEffect } from "react";

export const NetworkStatus = () => {
  const [isOnline, setIsOnline] = useState<boolean>(() =>
    typeof navigator !== "undefined" ? navigator.onLine : true
  );
  const [showRestored, setShowRestored] = useState<boolean>(false);

  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      setShowRestored(true);
      const timer = setTimeout(() => {
        setShowRestored(false);
      }, 4000);
      return () => clearTimeout(timer);
    };

    const handleOffline = () => {
      setIsOnline(false);
      setShowRestored(false);
    };

    window.addEventListener("online", handleOnline);
    window.addEventListener("offline", handleOffline);

    return () => {
      window.removeEventListener("online", handleOnline);
      window.removeEventListener("offline", handleOffline);
    };
  }, []);

  if (isOnline && !showRestored) {
    return null;
  }

  return (
    <div
      role="status"
      aria-live="polite"
      className="fixed bottom-4 right-4 z-50 max-w-sm transition-all duration-300 ease-in-out"
    >
      {!isOnline ? (
        <div className="flex items-center gap-3 bg-amber-500 text-slate-950 px-4 py-3 rounded-lg shadow-xl font-medium text-sm border border-amber-600">
          <svg
            className="w-5 h-5 flex-shrink-0 text-slate-900 animate-pulse"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 4.243a9 9 0 01-2.828-6.364m0 0l2.828 2.828M3 3l18 18"
            />
          </svg>
          <div>
            <p className="font-semibold">Modo sin conexión</p>
            <p className="text-xs text-slate-900/90">
              Navegando con datos cacheados de la PWA.
            </p>
          </div>
        </div>
      ) : (
        <div className="flex items-center gap-2 bg-emerald-600 text-white px-4 py-2.5 rounded-lg shadow-lg font-medium text-sm">
          <svg
            className="w-5 h-5 flex-shrink-0 text-white"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M5 13l4 4L19 7"
            />
          </svg>
          <span>Conexión a Internet restablecida</span>
        </div>
      )}
    </div>
  );
};

export default NetworkStatus;
