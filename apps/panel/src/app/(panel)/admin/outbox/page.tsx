"use client";

import { Settings } from "lucide-react";

import { OutboxAdminView } from "./components/outbox-admin-view";

export default function AdminOutboxPage() {
  return (
    <>
      <header className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold text-gray-800">Administración</h1>
          <p className="text-lg text-gray-500 mt-2">
            Outbox de eventos y cola de mensajes no entregados (DLQ)
          </p>
        </div>
        <Settings className="text-gray-400" size={24} />
      </header>

      <OutboxAdminView />
    </>
  );
}
