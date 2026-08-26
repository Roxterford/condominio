"use client"

import { motion, AnimatePresence } from "framer-motion"
import { X } from "lucide-react"

import { useSidebar } from "@/contexts/sidebar-context"

export function DetailPanel() {
  const { currentView, clearView } = useSidebar()

  return (
    <AnimatePresence>
      {currentView && (
        <>
          {/* Overlay */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={clearView}
            className="fixed inset-0 bg-black/20 z-30 hidden md:block"
          />

          {/* Panel */}
          <motion.aside
            initial={{ x: 420 }}
            animate={{ x: 0 }}
            exit={{ x: 420 }}
            transition={{ type: "spring", damping: 25, stiffness: 250 }}
            className="fixed top-0 right-0 h-full w-[400px] bg-gray-50 shadow-2xl z-40 flex flex-col border-l border-gray-200"
          >
            {/* Header */}
            <div className="h-20 px-6 flex items-center justify-between bg-white border-b border-gray-200 shrink-0">
              {currentView.title ? (
                <span className="inline-flex items-center px-3 py-1 rounded-lg bg-blue-50 text-blue-700 text-sm font-mono font-medium">
                  {currentView.title}
                </span>
              ) : (
                <h2 className="text-lg font-semibold text-gray-800">Detalle</h2>
              )}
              <button
                onClick={clearView}
                className="text-gray-400 hover:text-gray-600 transition-colors"
              >
                <X size={22} />
              </button>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto p-6">
              {currentView.component}
            </div>
          </motion.aside>
        </>
      )}
    </AnimatePresence>
  )
}
