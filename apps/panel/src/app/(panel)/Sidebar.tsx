import {
  LayoutDashboard,
  Home,
  CreditCard,
  BarChart3,
  Settings,
  X,
  Newspaper,
} from "lucide-react"

import { motion, AnimatePresence } from "framer-motion"
import Link from "next/link"

const menu = [
  {
    title: "Dashboard",
    icon: <LayoutDashboard size={20} />,
    active: true,
    href: '/dashboard'
  },
  {
    title: "Villas",
    icon: <Home size={20} />,
    href: '/villas'

  },
  {
    title: "Cuotas",
    icon: <Newspaper size={20} />,
    href: '/cuotas'

  },
  {
    title: "Pagos",
    icon: <CreditCard size={20} />,
    href: '/pagos'

  },
  {
    title: "Reportes",
    icon: <BarChart3 size={20} />,
    href: '/reportes'

  },
  {
    title: "Configuración",
    icon: <Settings size={20} />,
    href: '/configuracion'

  },
]

export default function Sidebar({
  sidebarOpen,
  setSidebarOpen,
}: any) {
  return (
    <>
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex md:w-64 bg-white border-r border-gray-200 flex-col">

        {/* Logo */}
        <div className="h-20 flex items-center px-6 border-b border-gray-100">

          <div className="w-8 h-8 rounded-2xl bg-blue-600 mr-4" />

          <div>

            <h1 className="text-xl font-bold text-gray-800">
              Condominio
            </h1>

            <p className="text-xs text-gray-400">
              Dashboard
            </p>

          </div>

        </div>

        {/* Menu */}
        <nav className="flex-1 p-3 space-y-2">

          {menu.map((item, index) => (

            <Link
              href={item.href}
              key={index}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                item.active
                  ? "bg-blue-600 text-white shadow-lg shadow-blue-100"
                  : "text-gray-600 hover:bg-gray-100"
              }`}
            >

              {item.icon}

              <span className="font-medium">
                {item.title}
              </span>

            </Link>

          ))}

        </nav>

      </aside>

      {/* Mobile Sidebar */}
      <AnimatePresence>

        {sidebarOpen && (

          <>
            {/* Overlay */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSidebarOpen(false)}
              className="fixed inset-0 bg-black/40 z-40 md:hidden"
            />

            {/* Sidebar */}
            <motion.aside
              initial={{ x: -300 }}
              animate={{ x: 0 }}
              exit={{ x: -300 }}
              transition={{ type: "spring", damping: 25 }}
              className="fixed top-0 left-0 w-72 h-full bg-white z-50 md:hidden flex flex-col"
            >

              {/* Header */}
              <div className="h-20 px-6 flex items-center justify-between border-b border-gray-100">

                <div className="flex items-center">

                  <div className="w-8 h-8 rounded-2xl bg-blue-600 mr-4" />

                  <div>

                    <h1 className="text-xl font-bold text-gray-800">
                      Condominio
                    </h1>

                    <p className="text-xs text-gray-400">
                      Dashboard
                    </p>

                  </div>

                </div>

                <button
                  onClick={() => setSidebarOpen(false)}
                  className="text-gray-500"
                >
                  <X size={22} />
                </button>

              </div>

              {/* Menu */}
              <nav className="flex-1 p-3 space-y-2">

                {menu.map((item, index) => (

                  <button
                    key={index}
                    className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                      item.active
                        ? "bg-blue-600 text-white shadow-lg shadow-blue-100"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                  >

                    {item.icon}

                    <span className="font-medium">
                      {item.title}
                    </span>

                  </button>

                ))}

              </nav>

            </motion.aside>

          </>

        )}

      </AnimatePresence>
    </>
  )
}