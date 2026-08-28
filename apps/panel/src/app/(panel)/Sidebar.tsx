"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import {
  LayoutDashboard,
  Home,
  CreditCard,
  BarChart3,
  Settings,
  Wrench,
  X,
  Newspaper,
} from "lucide-react";

type SidebarProps = {
  sidebarOpen: boolean;
  setSidebarOpen: (open: boolean) => void;
};

type MenuItem = {
  title: string;
  href: string;
  icon: React.ReactNode;
};

const menuItems: MenuItem[] = [
  {
    title: "Dashboard",
    icon: <LayoutDashboard size={20} />,
    href: "/dashboard",
  },
  { title: "Villas", icon: <Home size={20} />, href: "/villas" },
  { title: "Cuotas", icon: <Newspaper size={20} />, href: "/cuotas" },
  { title: "Pagos", icon: <CreditCard size={20} />, href: "/pagos/registrar" },
  { title: "Reportes", icon: <BarChart3 size={20} />, href: "/reportes" },
  {
    title: "Configuración",
    icon: <Settings size={20} />,
    href: "/configuracion",
  },
  {
    title: "Administración",
    icon: <Wrench size={20} />,
    href: "/admin/outbox",
  },
];

export default function Sidebar({ sidebarOpen, setSidebarOpen }: SidebarProps) {
  const pathname = usePathname();

  const isActive = (href: string) => pathname === href;

  const desktopItemClass = (href: string) =>
    `w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
      isActive(href)
        ? "bg-primary text-white shadow-lg shadow-blue-100"
        : "text-gray-600 hover:bg-gray-100"
    }`;

  const mobileItemClass = (href: string) =>
    `w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
      isActive(href)
        ? "bg-blue-600 text-white shadow-lg shadow-blue-100"
        : "text-gray-600 hover:bg-gray-100"
    }`;

  return (
    <>
      {/* Desktop Sidebar (izquierda) */}
      <aside className="hidden md:flex md:w-64 bg-white border-r border-gray-200 flex-col">
        {/* Logo */}
        <div className="h-20 flex items-center px-6 border-b border-gray-100">
          <div className="w-8 h-8 rounded-2xl bg-primary mr-4" />
          <div>
            <h1 className="text-xl font-bold text-gray-800">Condominio</h1>
            <p className="text-xs text-gray-400">Dashboard</p>
          </div>
        </div>

        {/* Menu */}
        <nav className="flex-1 p-3 space-y-2">
          {menuItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={desktopItemClass(item.href)}
            >
              {item.icon}
              <span className="font-medium">{item.title}</span>
            </Link>
          ))}
        </nav>
      </aside>

      {/* Mobile Sidebar (slide-in desde izquierda) */}
      <AnimatePresence>
        {sidebarOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSidebarOpen(false)}
              className="fixed inset-0 bg-black/40 z-40 md:hidden"
            />

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
                    <p className="text-xs text-gray-400">Dashboard</p>
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
                {menuItems.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    onClick={() => setSidebarOpen(false)}
                    className={mobileItemClass(item.href)}
                  >
                    {item.icon}
                    <span className="font-medium">{item.title}</span>
                  </Link>
                ))}
              </nav>
            </motion.aside>
          </>
        )}
      </AnimatePresence>
    </>
  );
}
