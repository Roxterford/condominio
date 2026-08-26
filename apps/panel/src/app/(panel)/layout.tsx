"use client";

import { useState } from "react";

import { SidebarProvider, useSidebar } from "@/contexts/sidebar-context";
import Sidebar from "./Sidebar";
import Header from "./Header";
import { DetailPanel } from "@/components/detail-panel/detail-panel";

function DashboardContent({ children }: { children: React.ReactNode }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { currentView } = useSidebar();

  return (
    <div className="flex min-h-screen bg-[#fff]">
      {/* Sidebar (izquierda) */}
      <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />

      {/* Main */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <Header sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />

        {/* Content */}
        <main className="flex-1 p-4 md:p-6 overflow-x-hidden overflow-y-auto">
          {children}
        </main>
      </div>

      {/* DetailPanel (derecha) — se abre al llamar setView() */}
      {currentView && <DetailPanel />}
    </div>
  );
}

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <DashboardContent>{children}</DashboardContent>
    </SidebarProvider>
  );
}
