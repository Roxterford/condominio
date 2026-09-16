"use client";

import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useCallback, useEffect, useMemo, useState } from "react";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import type { GlobalSearchQuery as GlobalSearchResult } from "@/providers/graphql/graphql";
import { useDebounce } from "@/hooks/useDebounce";
import { buscarMockSujetos, type MockSujeto } from "./mock-sujetos";

export const LIMITE_GLOBAL_SEARCH = 5;

const GLOBAL_SEARCH_QUERY = graphql(/* GraphQL */ `
  query GlobalSearch($busqueda: String!, $limit: Int!) {
    unidades: obtenerUnidades(
      filter: { codigo: { like: $busqueda } }
      paginator: { limit: $limit, page: 1 }
    ) {
      data {
        codigo
        estado
        deuda
        titular_primario {
          __typename
          ... on Sujeto {
            id
            cedula
            display_name
          }
          ... on Persona {
            nombres
            apellidos
          }
          ... on Ente {
            razon_social
          }
        }
      }
      total
    }

    operaciones: obtenerOperaciones(
      paginador: { limit: $limit, page: 1 }
      filtro: {
        or: [
          { concepto: { like: $busqueda } }
          { unidad: { like: $busqueda } }
          { proveedor_nombre: { like: $busqueda } }
        ]
      }
    ) {
      data {
        __typename
        ... on IOperacion {
          fecha
          operacion
          concepto
          monto
          moneda
        }
        ... on Pago {
          unidad {
            id
            codigo
          }
        }
        ... on GastoAProveedor {
          proveedor {
            nombre
          }
        }
      }
      total
    }

    proveedores: obtenerProveedores {
      id
      nombre
      rif
    }
  }
`);

type UnidadDeBusqueda = NonNullable<
  NonNullable<GlobalSearchResult["unidades"]>["data"]
>[number];
type OperacionDeBusqueda = NonNullable<
  NonNullable<GlobalSearchResult["operaciones"]>["data"]
>[number];
type ProveedorDeBusqueda = NonNullable<
  GlobalSearchResult["proveedores"]
>[number];

const STORAGE_KEY = "condominio:busquedas-recientes";

export function useRecentSearches(max = 5) {
  const [recientes, setRecientes] = useState<string[]>([]);

  useEffect(() => {
    try {
      const raw = window.localStorage.getItem(STORAGE_KEY);
      if (raw) setRecientes(JSON.parse(raw) as string[]);
    } catch {
      // almacenamiento no disponible
    }
  }, []);

  const agregarBusqueda = useCallback(
    (value: string) => {
      const v = value.trim();
      if (!v) return;
      setRecientes((prev) => {
        const next = [
          v,
          ...prev.filter(
            (x) => x.toLocaleLowerCase() !== v.toLocaleLowerCase(),
          ),
        ].slice(0, max);
        try {
          window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
        } catch {
          // almacenamiento no disponible
        }
        return next;
      });
    },
    [max],
  );

  return { recientes, agregarBusqueda };
}

export type GlobalSearchData = {
  debouncedTerm: string;
  isLoading: boolean;
  unidades: UnidadDeBusqueda[];
  unidadesTotal: number;
  operaciones: OperacionDeBusqueda[];
  operacionesTotal: number;
  proveedores: ProveedorDeBusqueda[];
  sujetos: MockSujeto[];
  totalCoincidencias: number;
};

export function useGlobalSearch(term: string): GlobalSearchData {
  const [debouncedTerm, setDebouncedTerm] = useState("");
  const onDebounce = useDebounce(setDebouncedTerm);

  useEffect(() => {
    onDebounce(term);
  }, [term, onDebounce]);

  const { data, isLoading } = useQuery({
    queryKey: ["global-search", debouncedTerm],
    enabled: debouncedTerm.trim().length > 0,
    queryFn: async () => {
      const result = await execute(GLOBAL_SEARCH_QUERY, {
        busqueda: `%${debouncedTerm.trim()}%`,
        limit: LIMITE_GLOBAL_SEARCH,
      });
      return result.data;
    },
    placeholderData: keepPreviousData,
  });

  const proveedores = useMemo(() => {
    const todos = data?.proveedores ?? [];
    const q = debouncedTerm.trim().toLocaleLowerCase("es");
    if (!q) return [];
    return todos
      .filter((p) => p.nombre.toLocaleLowerCase("es").includes(q))
      .slice(0, LIMITE_GLOBAL_SEARCH);
  }, [data?.proveedores, debouncedTerm]);

  const sujetos = useMemo(
    () => buscarMockSujetos(debouncedTerm, LIMITE_GLOBAL_SEARCH),
    [debouncedTerm],
  );

  const unidadesTotal = data?.unidades?.total ?? 0;
  const operacionesTotal = data?.operaciones?.total ?? 0;
  const totalCoincidencias =
    unidadesTotal + operacionesTotal + proveedores.length + sujetos.length;

  return {
    debouncedTerm: debouncedTerm.trim(),
    isLoading,
    unidades: data?.unidades?.data ?? [],
    unidadesTotal,
    operaciones: data?.operaciones?.data ?? [],
    operacionesTotal,
    proveedores,
    sujetos,
    totalCoincidencias,
  };
}
