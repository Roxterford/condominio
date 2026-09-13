"use client";

import { ResultadosPorPagina } from "./resultados-por-pagina";

export interface PaginacionFooterProps {
  currentPage: number;
  totalPages: number;
  limit: number;
  onLimitChange: (limit: number) => void;
}

export function PaginacionFooter({
  currentPage,
  totalPages,
  limit,
  onLimitChange,
}: PaginacionFooterProps) {
  return (
    <div className="flex items-center justify-between pt-4">
      <p className="text-muted-foreground text-sm">
        Página <span className="font-bold">{currentPage}</span> de{" "}
        <span className="font-bold">{totalPages}</span>
      </p>
      <ResultadosPorPagina limit={limit} onLimitChange={onLimitChange} />
    </div>
  );
}
