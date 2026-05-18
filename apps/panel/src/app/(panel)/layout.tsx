"use client"
import { useState } from "react"

import Sidebar from "./Sidebar"
import Header from "./Header"

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {

  const [sidebarOpen, setSidebarOpen] = useState(false)

  return (
    <div className="flex min-h-screen bg-[#F5F6FA]">

      {/* Sidebar Desktop */}
      <Sidebar
  sidebarOpen={sidebarOpen}
  setSidebarOpen={setSidebarOpen}
/>

      {/* Main */}
      <div className="flex-1 flex flex-col">

        {/* Header */}
        <Header
          sidebarOpen={sidebarOpen}
          setSidebarOpen={setSidebarOpen}
        />

        {/* Content */}
        <main className="flex-1 p-4 md:p-6 overflow-x-hidden overflow-y-auto">
          {children}
        </main>

      </div>

    </div>
  )
}