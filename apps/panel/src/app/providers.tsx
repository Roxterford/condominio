"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // Con SSR, queremos evitar que el cliente haga fetch
        // inmediatamente después de recibir los datos del servidor.
        staleTime: 60 * 1000,
      },
    },
  });
}

// Esta variable solo existirá en el navegador
let browserQueryClient: QueryClient | undefined = undefined;

function getQueryClient() {
  // Comprobación estándar de JS para el entorno
  if (typeof window === "undefined") {
    // Servidor: Siempre crea un cliente nuevo para cada petición
    return makeQueryClient();
  } else {
    // Cliente: Reutiliza el cliente existente para no perder la caché
    if (!browserQueryClient) browserQueryClient = makeQueryClient();
    return browserQueryClient;
  }
}

export default function Providers({ children }: { children: React.ReactNode }) {
  // Obtenemos el cliente (nuevo en servidor, persistente en cliente)
  const queryClient = getQueryClient();

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
