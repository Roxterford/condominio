"use client";

import { useState } from "react";

import { DrawerProvider } from "@/contexts/drawer-context";
import { DynamicDrawer } from "@/components/dynamic-drawer";
import Sidebar from "./Sidebar";
import Header from "./Header";

function DashboardContent({ children }: { children: React.ReactNode }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

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

      {/* Dynamic Drawer — se abre al llamar useDrawer().open() */}
      <DynamicDrawer />
    </div>
  );
}

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <DrawerProvider>
      <DashboardContent>{children}</DashboardContent>
    </DrawerProvider>
  );
}
