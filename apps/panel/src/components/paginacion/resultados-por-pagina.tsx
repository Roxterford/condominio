"use client";

import { cn } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export const RESULTADOS_POR_PAGINA = [5, 10, 20, 50] as const;

export interface ResultadosPorPaginaProps {
  limit: number;
  onLimitChange: (limit: number) => void;
  className?: string;
}
export function ResultadosPorPagina({
  limit,
  onLimitChange,
  className,
}: ResultadosPorPaginaProps) {
  return (
    <div className={cn("flex items-center gap-2", className)}>
      <span className="text-muted-foreground text-sm whitespace-nowrap">
        Resultados por página
      </span>
      <Select
        value={String(limit)}
        onValueChange={(value) => onLimitChange(Number(value))}
      >
        <SelectTrigger className="w-20">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {RESULTADOS_POR_PAGINA.map((n) => (
            <SelectItem key={n} value={String(n)}>
              {n}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
