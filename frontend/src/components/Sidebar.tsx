interface Prop {
  collapsed: boolean;
  onClose: () => void;
}


import { Link } from "react-router-dom";
const Sidebar = ({ collapsed, onClose }: Prop) => {

  return (
    <aside className={`fixed inset-y-0 left-0 z-50 md:relative md:translate-x-0 ${collapsed ? '-translate-x-full md:w-20' : 'translate-x-0 w-9/12 sm:w-8/12 md:w-64'} bg-slate-900 text-white flex flex-col transition-all duration-300 overflow-hidden whitespace-nowrap`}>
      <div className={`h-16 flex items-center ${collapsed ? 'justify-center' : 'px-6 justify-between'} text-2xl font-bold border-b border-slate-700 transition-all`}>
        {collapsed ? 'MC' : 'MultiCatálogo'}

        <button className="md:hidden bg-slate-800" onClick={onClose}>
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z" /></svg>
        </button>
      </div>

      <nav className="flex-1 p-4 space-y-2">
        <Link
          to="/"
          className="flex items-center gap-4 p-3 rounded hover:bg-slate-800 transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentColor" className="shrink-0"><path d="M120-120v-320h320v320H120Zm0-400v-320h320v320H120Zm400 400v-320h320v320H520Zm0-400v-320h320v320H520Z" /></svg>
          <span className={`transition-all duration-300 ${collapsed ? 'opacity-0' : 'opacity-100'}`}>Dashboard</span>
        </Link>
        <Link
          to="/catalogo"
          className="flex items-center gap-4 p-3 rounded hover:bg-slate-800 transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentColor" className="shrink-0"><path d="M280-600v-80h560v80H280Zm0 160v-80h560v80H280Zm0 160v-80h560v80H280ZM160-600q-17 0-28.5-11.5T120-640q0-17 11.5-28.5T160-680q17 0 28.5 11.5T200-640q0 17-11.5 28.5T160-600Zm0 160q-17 0-28.5-11.5T120-480q0-17 11.5-28.5T160-520q17 0 28.5 11.5T200-480q0 17-11.5 28.5T160-440Zm0 160q-17 0-28.5-11.5T120-320q0-17 11.5-28.5T160-360q17 0 28.5 11.5T200-320q0 17-11.5 28.5T160-280Z" /></svg>
          <span className={`transition-all duration-300 ${collapsed ? 'opacity-0' : 'opacity-100'}`}>Catálogo</span>
        </Link>
        <Link
          to="/mi-red"
          className="flex items-center gap-4 p-3 rounded hover:bg-slate-800 transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentColor" className="shrink-0"><path d="M0-240v-63q0-43 44-70t116-27q13 0 25 .5t23 2.5q-14 21-21 44t-7 48v65H0Zm240 0v-65q0-32 17.5-58.5T307-410q32-20 76.5-30t96.5-10q53 0 97.5 10t76.5 30q32 20 49 46.5t17 58.5v65H240Zm540 0H600v-65q0-25-6.5-48t-20.5-44q11-2 23-2.5t24-.5q72 0 116 27t44 70v63ZM160-440q-33 0-56.5-23.5T80-520q0-34 23.5-57t56.5-23q34 0 57 23t23 57q0 33-23 56.5T160-440Zm640 0q-34 0-57.5-23.5T720-520q0-34 23.5-57t57.5-23q33 0 56.5 23t23.5 57q0 33-23.5 56.5T800-440Zm-320-40q-50 0-85-35t-35-85q0-51 35-85.5t85-34.5q51 0 85.5 34.5T600-600q0 50-34.5 85T480-480Z" /></svg>
          <span className={`transition-all duration-300 ${collapsed ? 'opacity-0' : 'opacity-100'}`}>Mi Red</span>
        </Link>
      </nav>
    </aside>
  );
};
export default Sidebar;
