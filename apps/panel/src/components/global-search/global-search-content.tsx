"use client";

import { Fragment, useEffect, useRef } from "react";
import { SearchX } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { Kbd } from "@/components/ui/kbd";
import { SearchResultRow, type ResultRow } from "./search-result-row";

export type GlobalSearchState = "idle" | "loading" | "results" | "empty";

type GlobalSearchContentProps = {
  state: GlobalSearchState;
  query: string;
  rows: ResultRow[];
  activeIndex: number;
  onHoverIndex: (index: number) => void;
  onSelectRow: (row: ResultRow) => void;
  recientes: string[];
  onRecienteClick: (value: string) => void;
  totalCoincidencias: number;
  onReset: () => void;
};

export function GlobalSearchContent({
  state,
  query,
  rows,
  activeIndex,
  onHoverIndex,
  onSelectRow,
  recientes,
  onRecienteClick,
  totalCoincidencias,
  onReset,
}: GlobalSearchContentProps) {
  const resultsRef = useRef<HTMLDivElement>(null);
  const scrollRafRef = useRef<number>(0);

  useEffect(() => {
    const container = resultsRef.current;
    const active = container?.querySelector<HTMLElement>(
      '[data-search-active="true"]',
    );
    if (!container || !active) return;

    const containerRect = container.getBoundingClientRect();
    const activeRect = active.getBoundingClientRect();
    const gap = 8;
    let target = container.scrollTop;
    if (activeRect.top < containerRect.top + gap) {
      target = container.scrollTop + (activeRect.top - containerRect.top) - gap;
    } else if (activeRect.bottom > containerRect.bottom - gap) {
      target =
        container.scrollTop + (activeRect.bottom - containerRect.bottom) + gap;
    } else {
      return;
    }

    if (
      window.matchMedia("(prefers-reduced-motion: reduce)").matches ||
      Math.abs(target - container.scrollTop) < 1
    ) {
      container.scrollTop = target;
      return;
    }

    const from = container.scrollTop;
    const duration = 140;
    const start = performance.now();
    const step = (now: number) => {
      const t = Math.min(1, (now - start) / duration);
      const eased = 1 - Math.pow(1 - t, 3);
      container.scrollTop = from + (target - from) * eased;
      if (t < 1) scrollRafRef.current = requestAnimationFrame(step);
    };
    scrollRafRef.current = requestAnimationFrame(step);

    return () => cancelAnimationFrame(scrollRafRef.current);
  }, [activeIndex, state, rows]);

  if (state === "loading") {
    return (
      <div className="flex items-center gap-3 px-3 py-6 text-sm text-muted-foreground">
        <Spinner className="size-4" />
        Buscando en unidades, propietarios y operaciones…
      </div>
    );
  }

  if (state === "empty") {
    return (
      <div className="flex flex-col items-center gap-2 px-4 py-8 text-center">
        <span className="flex size-10 items-center justify-center rounded-full bg-muted text-muted-foreground">
          <SearchX className="size-5" />
        </span>
        <p className="text-sm font-medium">Sin resultados para "{query}"</p>
        <p className="max-w-70 text-xs text-muted-foreground">
          Prueba con el código de una villa (ej. villa-12), el concepto de una
          operación o el nombre de un proveedor.
        </p>
        <button
          type="button"
          onClick={onReset}
          className="mt-1 text-xs font-medium text-primary hover:underline"
        >
          Limpiar búsqueda
        </button>
      </div>
    );
  }

  return (
    <>
      {state === "idle" && recientes.length > 0 ? (
        <section className="px-2 pt-2 pb-3">
          <h4 className="px-2 py-1.5 text-[11px] font-medium tracking-wide text-muted-foreground uppercase">
            Búsquedas recientes
          </h4>
          <div className="flex flex-wrap gap-1.5 px-2 pb-1">
            {recientes.map((r) => (
              <button
                key={r}
                type="button"
                onClick={() => onRecienteClick(r)}
                className="rounded-full border border-border px-3 py-1 text-xs text-muted-foreground transition-colors hover:border-primary hover:text-primary"
              >
                {r}
              </button>
            ))}
          </div>
        </section>
      ) : null}

      <section
        ref={resultsRef}
        className="scroll-my-1 overflow-y-auto px-1 pb-1"
      >
        {rows.map((row, index) => {
          const firstOfGroup =
            index === 0 || rows[index - 1]?.groupLabel !== row.groupLabel;
          const lastOfGroup =
            index === rows.length - 1 ||
            rows[index + 1]?.groupLabel !== row.groupLabel;
          return (
            <Fragment key={row.key}>
              {firstOfGroup && (
                <h4 className="px-2 py-1.5 pt-3 text-[11px] font-medium tracking-wide text-muted-foreground uppercase">
                  {row.groupLabel}
                </h4>
              )}
              <SearchResultRow
                row={row}
                query={query}
                active={index === activeIndex}
                onHover={() => onHoverIndex(index)}
                onSelect={() => onSelectRow(row)}
              />
              {lastOfGroup && (
                <div className="mx-2 mt-1.5 mb-1 h-px bg-border" />
              )}
            </Fragment>
          );
        })}
      </section>

      <footer className="flex items-center justify-between border-t border-border px-3 py-2 text-xs text-muted-foreground">
        {state === "idle" ? (
          <span>Busca en todas las secciones del condominio</span>
        ) : (
          <span>
            {totalCoincidencias}{" "}
            {totalCoincidencias === 1 ? "coincidencia" : "coincidencias"}
          </span>
        )}
        {state !== "idle" && (
          <span className="flex items-center gap-1">
            <Kbd>↵</Kbd> para abrir
          </span>
        )}
      </footer>
    </>
  );
}
