import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import { Box, CircleCheckBig, TriangleAlert, User } from "lucide-react";
import { EstadoUnidadTag } from "@/components/estado-unidad-tag";
import { DeudaUnidadTag } from "@/components/deuda-unidad-tag";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { money } from "@/lib/money-display";
import { Badge } from "@/components/ui/badge";
import StatCard from "@/components/ui/StatCard";
import { EstadoDeDeuda, VillaPageQuery } from "@/providers/graphql/graphql";
import { es } from "date-fns/locale";
import { format } from "date-fns";
import { today } from "@/lib/today";
import { PagosTable } from "./components/pagos-table";
import { DeudasTable } from "./components/deudas-table";
import { RegistrarPagoButton } from "./components/registrar-pago-button";

const PageQuery = graphql(/* GraphQL */ `
  query VillaPage($codigo: String!, $estado_deuda_pendiente: String!) {
    ultimo_pago: obtenerPagos(
      filtro: { unidad: { eq: $codigo } }
      paginador: { limit: 1, page: 1 }
    ) {
      data {
        fecha
        metodo
        total
      }
    }
    unidad: obtenerUnidadPorCodigo(codigo: $codigo) {
      id
      codigo
      estado
      wallet
      deuda
      titular_primario {
        __typename
        ... on Sujeto {
          id
          cedula
          display_name
          email
          telefono
        }
      }
      titulares {
        __typename

        ... on Sujeto {
          id
          cedula
          display_name
          email
        }
      }
    }
    pagos: obtenerPagos(filtro: { unidad: { eq: $codigo } }) {
      data {
        __typename
        fecha
        operacion
        concepto
        metodo
        moneda
        monto
        registro
        tasa
        total
        unidad {
          id
          codigo
        }
      }
    }
    deudas: obtenerDeudas(filtro: { unidad: { eq: $codigo } }) {
      data {
        id
        cuota {
          __typename
          ... on Deuda__Cuota {
            id
            nombre
          }
        }
        deuda
        monto
        estado
      }
    }
    deudas_pendientes: obtenerDeudas(
      filtro: {
        unidad: { eq: $codigo }
        estado: { eq: $estado_deuda_pendiente }
      }
    ) {
      total
    }
  }
`);

export interface VillaPageProps {
  params: Promise<{ codigo: string }>;
}

export default async function VillaPage(page: VillaPageProps) {
  const { codigo } = await page.params;

  return renderGraphql(
    await execute(PageQuery, {
      codigo,
      estado_deuda_pendiente: EstadoDeDeuda.Pendiente,
    }),
    ({ unidad, pagos, deudas, deudas_pendientes, ...data }) => {
      const ultimo_pago = data.ultimo_pago.data.at(0) as
        VillaPageQuery["ultimo_pago"]["data"][0] | undefined;

      if (!unidad) {
        return <div>Unidad no encontrada</div>;
      }

      return (
        <>
          <header className="flex gap-20 items-center justify-between">
            <section className="flex gap-20 justify-between flex-1">
              <section>
                <div className="flex gap-4 items-start">
                  <div className="bg-primary/10 rounded-lg p-2">
                    <Box className="size-10 text-primary" />
                  </div>
                  <div className="grid gap-1">
                    <p className="font-medium text-nowrap">{unidad.codigo}</p>
                    <DeudaUnidadTag pending={unidad.deuda > 0} />
                  </div>
                </div>
              </section>

              <ul className="flex gap-10 justify-end">
                <li>
                  <h3 className="text-sm font-medium">Estado</h3>
                  <EstadoUnidadTag state={unidad.estado} />
                </li>
                <li>
                  <h3 className="text-sm font-medium">Cuenta</h3>
                  {unidad.wallet > 0 ? (
                    <span>{money(unidad.wallet)}</span>
                  ) : (
                    <span className="text-muted-foreground">-</span>
                  )}
                </li>
              </ul>
            </section>
            <section>
              <RegistrarPagoButton unidad={unidad} />
            </section>
          </header>

          <ul className="mt-10">
            {unidad.titular_primario && (
              <li>
                <h3 className="text-sm font-medium">Responsable</h3>
                <div className="flex gap-2 items-center">
                  <div className="bg-gray-200 w-min p-3 rounded-full">
                    <User className="size-4" />
                  </div>
                  <p className="font-medium">
                    {unidad.titular_primario?.display_name}
                  </p>
                </div>
              </li>
            )}
          </ul>

          <ul className="statcards |  mt-10">
            <li>
              <StatCard
                title={"Estado de cuenta"}
                value={deudas_pendientes.total ? "En deuda" : "Al día"}
                subtitle={""}
                icon={
                  deudas_pendientes.total ? (
                    <TriangleAlert className="text-yellow-700" />
                  ) : (
                    <CircleCheckBig className="text-green-700" />
                  )
                }
                color={
                  deudas_pendientes.total ? "bg-yellow-100" : "bg-green-100"
                }
              />
            </li>
            <li>
              <StatCard
                title={"Deuda total"}
                value={money(unidad.deuda)}
                subtitle={""}
                icon={undefined}
                color={""}
              />
            </li>
            <li>
              <StatCard
                title={"Cuotas pendientes"}
                value={deudas_pendientes.total}
                subtitle={""}
                icon={undefined}
                color={""}
              />
            </li>

            <li>
              <StatCard
                title={"Último pago"}
                value={
                  ultimo_pago?.fecha
                    ? format(
                        ultimo_pago.fecha,
                        `d MMM ${ultimo_pago.fecha.getFullYear() === today().getFullYear() ? "" : "yyyy"}`,
                        {
                          locale: es,
                        },
                      )
                    : "Ninguno"
                }
                subtitle={`${money(ultimo_pago?.total ?? 0)} · ${ultimo_pago?.metodo ?? "Nunca"}`}
                icon={undefined}
                color={""}
              />
            </li>
          </ul>

          <Tabs defaultValue="deudas">
            <TabsList variant="line">
              <TabsTrigger value="deudas">Deudas</TabsTrigger>
              <TabsTrigger value="pagos">Pagos</TabsTrigger>
              <TabsTrigger value="documentos">Documentos</TabsTrigger>
              <TabsTrigger value="notificaciones">Notificaciones</TabsTrigger>
            </TabsList>

            <TabsContent value="deudas">
              <h3>Historial de deudas</h3>
              <DeudasTable deudas={deudas.data} />
            </TabsContent>
            <TabsContent value="pagos">
              <h3>Historial de pagos</h3>
              <PagosTable
                pagos={pagos.data}
                unidadTitular={unidad.titular_primario}
              />
            </TabsContent>
            <TabsContent value="documentos">
              <h3>Información de la propiedad</h3>
              {unidad.titulares && (
                <section>
                  <h4 className="font-medium">Titulares</h4>
                  <ul className="space-y-5">
                    {unidad.titulares.map((titular) => (
                      <li key={titular.id} className="flex gap-2 items-center">
                        <div className="bg-gray-200 w-min p-3 rounded-full">
                          <User className="size-4" />
                        </div>
                        <div>
                          <div className="flex gap-3 items-center">
                            <p className="font-medium">
                              {titular.display_name}
                            </p>
                            {titular.id === unidad.titular_primario?.id && (
                              <Badge className="bg-teal-100 text-teal-700">
                                Principal
                              </Badge>
                            )}
                          </div>
                          <p className="text-sm text-gray-500">
                            {titular.email}
                          </p>
                        </div>
                      </li>
                    ))}
                  </ul>
                </section>
              )}
            </TabsContent>
          </Tabs>
        </>
      );
    },
  );
}
