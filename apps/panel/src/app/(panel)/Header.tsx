import { ChevronDown, Menu } from "lucide-react";
import { GlobalSearch } from "@/components/global-search/global-search";
const fotoPerfil = "/mocks/img/1.jpg";

export default function Header({ setSidebarOpen }: any) {
  return (
    <header className="h-20 bg-white border-b border-gray-200 flex items-center justify-between px-4 md:px-6">
      {/* Left */}
      <div className="flex items-center gap-4 w-full">
        {/* Mobile Menu */}
        <button onClick={() => setSidebarOpen(true)} className="md:hidden">
          <Menu size={24} className="text-gray-700" />
        </button>

        {/* Search */}
        <GlobalSearch />
      </div>

      {/* Right */}
      <div className="hidden md:flex items-center gap-4 ml-6">
        <img
          src={fotoPerfil as any}
          alt=""
          className="w-11 h-11 rounded-full"
        />

        <div>
          <h3 className="text-sm font-semibold text-gray-800">Ed Ccs</h3>

          <p className="text-xs text-gray-400">Administrador</p>
        </div>

        <ChevronDown size={18} className="text-gray-400" />
      </div>
    </header>
  );
}
