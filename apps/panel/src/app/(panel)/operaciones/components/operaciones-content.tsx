"use client";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { ArrowLeftRight, CreditCardMinus, Plus } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { OpPagination } from "./operacion-pagination";
import { OperacionesTable } from "./operaciones-table";

const PageQuery = graphql(/* GraphQL */ `
  query OperacionesPage($page: Int!) {
    operaciones: obtenerOperaciones(paginador: { limit: 20, page: $page }) {
      data {
        __typename

        ... on IOperacion {
          fecha
          operacion
          concepto
          metodo
          monto
          moneda
          total
        }

        ... on GastoAProveedor {
          proveedor {
            nombre
          }
        }

        ... on Pago {
          unidad {
            codigo
          }
        }
      }

      limit
      page
      pages
      total
    }
  }
`);

export function OperacionesContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const currentPage = Math.max(1, Number(searchParams.get("page")) || 1);

  const { data, isLoading } = useQuery({
    queryKey: ["operaciones", currentPage],
    queryFn: async () => {
      const result = await execute(PageQuery, { page: currentPage });
      return result.data;
    },
    placeholderData: keepPreviousData,
  });

  const operaciones = data?.operaciones;
  const totalPages = operaciones?.pages ?? 1;

  const setPage = (page: number) => {
    router.push(`/operaciones?page=${page}`);
  };

  return (
    <>
      <header className="flex items-end justify-between">
        <div>
          <h1>Operaciones</h1>
          <p className="page-description">
            Finanzas · Movimientos, pagos y gastos
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <CreditCardMinus /> Gasto
          </Button>
          <Button variant="outline">
            <ArrowLeftRight /> Transacción
          </Button>
          <Button>
            <Plus /> Pago
          </Button>
        </div>
      </header>
      <section>
        {isLoading ? (
          <div className="flex justify-center py-8">
            <Spinner className="size-6" />
          </div>
        ) : (
          <>
            <OpPagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setPage}
              className="justify-end"
            />
            <OperacionesTable data={operaciones?.data ?? []} />
          </>
        )}
      </section>
    </>
  );
}
