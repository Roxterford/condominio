import {
  Search,
  ChevronDown,
  Menu,
} from "lucide-react"
const  fotoPerfil = '/mocks/img/1.jpg'; 

export default function Header({
  setSidebarOpen,
}: any) {
  return (
    <header className="h-20 bg-white border-b border-gray-200 flex items-center justify-between px-4 md:px-6">

      {/* Left */}
      <div className="flex items-center gap-4 w-full">

        {/* Mobile Menu */}
        <button
          onClick={() => setSidebarOpen(true)}
          className="md:hidden"
        >
          <Menu size={24} className="text-gray-700" />
        </button>

        {/* Search */}
        <div className="relative w-full md:w-[420px]">

          <Search
            size={18}
            className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"
          />

          <input
            type="text"
            placeholder="Buscar..."
            className="w-full h-12 rounded-2xl border border-gray-200 bg-gray-50 pl-11 pr-4 text-sm outline-none focus:border-blue-500"
          />

        </div>

      </div>

      {/* Right */}
      <div className="hidden md:flex items-center gap-4 ml-6">

        <img
          src={fotoPerfil as any}
          alt=""
          className="w-11 h-11 rounded-full"
        />

        <div>

          <h3 className="text-sm font-semibold text-gray-800">
            Ed Ccs
          </h3>

          <p className="text-xs text-gray-400">
            Administrador
          </p>

        </div>

        <ChevronDown
          size={18}
          className="text-gray-400"
        />

      </div>

    </header>
  )
}