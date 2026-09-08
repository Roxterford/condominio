"use client";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import {
  AlertCircle,
  ArrowLeftRight,
  CircleDollarSign,
  CreditCardMinus,
  DollarSign,
  Plus,
  Receipt,
} from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { OpPagination } from "./operacion-pagination";
import { OperacionesTable } from "./operaciones-table";
import { RegistrarPagoOverlay } from "@/features/administracion/components/registrar-pago-overlay";
import { useOverlay } from "@/hooks/useOverlay";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import StatCard from "@/components/ui/StatCard";
import { InputGroup, InputGroupInput } from "@/components/ui/input-group";

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

export function OperacionesPageContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const currentPage = Math.max(1, Number(searchParams.get("page")) || 1);

  const { data, isLoading, refetch } = useQuery({
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

  const registrarPago = useOverlay({
    closeOnDone: true,
    onDone() {
      refetch();
    },
  });

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
          <Button onClick={registrarPago.open}>
            <Plus /> Pago
          </Button>
        </div>
      </header>
      <section className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6 mt-8">
        <StatCard
          title="Ingresos"
          value="$440"
          subtitle="3 pagos · Agosto"
          color="bg-green-100"
          icon={<DollarSign className="text-green-600" />}
        />

        <StatCard
          title="Egresos"
          value="$800"
          subtitle="2 gastos"
          color="bg-yellow-100"
          icon={<CircleDollarSign className="text-yellow-600" />}
        />

        <StatCard
          title="Neto del periodo"
          value="−$360"
          subtitle="ingresos − egresos"
          color="bg-rose-100"
          icon={<AlertCircle className="text-rose-600" />}
        />

        <StatCard
          title="Operaciones"
          value="6"
          subtitle="incl. 1 compensación"
          color="bg-cyan-100"
          icon={<Receipt className="text-cyan-600" />}
        />
      </section>
      <section className="mt-5">
        <Tabs defaultValue="todas">
          <TabsList variant="line">
            <TabsTrigger value="todas">Todas</TabsTrigger>
            <TabsTrigger value="pagos">pagos</TabsTrigger>
            <TabsTrigger value="gastos">Gastos</TabsTrigger>
            <TabsTrigger value="transacciones">Transacciones</TabsTrigger>
          </TabsList>
          <TabsContent value="todas">
            {isLoading ? (
              <div className="flex justify-center py-8">
                <Spinner className="size-6" />
              </div>
            ) : (
              <>
                <div className="flex">
                  <form>
                    <InputGroup>
                      <InputGroupInput
                        placeholder="Buscar por villa o propietario"
                        className="md:min-w-64"
                      />
                    </InputGroup>
                  </form>
                  <OpPagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    onPageChange={setPage}
                    className="justify-end"
                  />
                </div>
                <OperacionesTable data={operaciones?.data ?? []} />
              </>
            )}
          </TabsContent>
        </Tabs>
      </section>
      <RegistrarPagoOverlay {...registrarPago.overlayProps} />
    </>
  );
}
