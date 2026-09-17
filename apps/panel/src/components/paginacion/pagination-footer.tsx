"use client";

import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ResultadosPorPagina } from "./resultados-por-pagina";

export interface PaginacionFooterProps {
  currentPage: number;
  totalPages: number;
  limit: number;
  onLimitChange: (limit: number) => void;
  onPageChange: (page: number) => void;
}

export function PaginacionFooter({
  currentPage,
  totalPages,
  limit,
  onLimitChange,
  onPageChange,
}: PaginacionFooterProps) {
  return (
    <div className="flex items-center justify-between pt-4">
      <p className="text-muted-foreground text-sm">
        Página <span className="font-bold">{currentPage}</span> de{" "}
        <span className="font-bold">{totalPages}</span>
      </p>
      <div className="flex items-center gap-3">
        <ResultadosPorPagina limit={limit} onLimitChange={onLimitChange} />
        <div className="flex items-center gap-1">
          <Button
            variant="outline"
            size="icon"
            aria-label="Página anterior"
            disabled={currentPage <= 1}
            onClick={() => onPageChange(currentPage - 1)}
          >
            <ChevronLeftIcon />
          </Button>
          <Button
            variant="outline"
            size="icon"
            aria-label="Página siguiente"
            disabled={currentPage >= totalPages}
            onClick={() => onPageChange(currentPage + 1)}
          >
            <ChevronRightIcon />
          </Button>
        </div>
      </div>
    </div>
  );
}
