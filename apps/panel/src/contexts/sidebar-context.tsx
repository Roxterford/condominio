"use client"

import { createContext, useContext, useState, useCallback, type ReactNode } from "react"

type SetViewOptions = {
  component: ReactNode
  title?: string
}

type SidebarContextType = {
  currentView: SetViewOptions | null
  setView: (view: SetViewOptions) => void
  clearView: () => void
}

const SidebarContext = createContext<SidebarContextType | null>(null)

export function SidebarProvider({ children }: { children: ReactNode }) {
  const [currentView, setCurrentView] = useState<SetViewOptions | null>(null)

  const setView = useCallback((view: SetViewOptions) => {
    setCurrentView(view)
  }, [])

  const clearView = useCallback(() => {
    setCurrentView(null)
  }, [])

  return (
    <SidebarContext value={{ currentView, setView, clearView }}>
      {children}
    </SidebarContext>
  )
}

export function useSidebar() {
  const ctx = useContext(SidebarContext)
  if (!ctx) throw new Error("useSidebar debe usarse dentro de un SidebarProvider")
  return ctx
}
