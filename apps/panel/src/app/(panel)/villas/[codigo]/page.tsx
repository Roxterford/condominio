import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import { Box, User } from "lucide-react";
import { EstadoUnidadTag } from "@/components/estado-unidad-tag";
import { DeudaUnidadTag } from "@/components/deuda-unidad-tag";
import { Button } from "@/components/ui/button";
import { CreditCard } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { money } from "@/lib/money-display";
import { Badge } from "@/components/ui/badge";
import { EstadoDeudaTag } from "@/components/estado-deuda-tag";
import Link from "next/link";

const PageQuery = graphql(/* GraphQL */ `
  query VillaPage($codigo: String!) {
    unidad: obtenerUnidadPorCodigo(codigo: $codigo) {
      codigo
      estado
      wallet
      titular_primario {
        __typename
        ... on Sujeto {
          display_name
        }
      }
      titulares {
        __typename

        ... on Sujeto {
          cedula
          display_name
        }
      }
    }
    pagos: obtenerPagos(filtro: { unidad: { eq: $codigo } }) {
      data {
        operacion
        concepto
        monto
        moneda
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
    }),
    ({ unidad, pagos, deudas }) => {
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
                    <Box className="size-10" />
                  </div>
                  <div className="grid gap-1">
                    <p className="font-medium text-nowrap">{unidad.codigo}</p>
                    <DeudaUnidadTag pending={unidad.wallet < 0} />
                  </div>
                </div>
              </section>

              <ul className="flex gap-10 justify-end">
                <li>
                  <h3 className="text-sm font-medium">Estado</h3>
                  <EstadoUnidadTag state={unidad.estado} />
                </li>
                <li>
                  <h3 className="text-sm font-medium">Cuenta</h3>$
                  {unidad.wallet}
                </li>
              </ul>
            </section>
            <section>
              <Button>
                <CreditCard /> Registrar pago
              </Button>
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

          <Tabs defaultValue="pagos">
            <TabsList variant="line">
              <TabsTrigger value="pagos">Pagos</TabsTrigger>
              <TabsTrigger value="deudas">Deudas</TabsTrigger>
              <TabsTrigger value="documentos">Documentos</TabsTrigger>
              <TabsTrigger value="notificaciones">Notificaciones</TabsTrigger>
            </TabsList>
            <TabsContent value="pagos">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="table__head">Concepto</TableHead>
                    <TableHead className="table__head">Monto</TableHead>
                    <TableHead className="table__head">Fecha</TableHead>
                    <TableHead className="table__head">Estado</TableHead>
                    <TableHead className="table__head text-end">
                      Acciones
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {pagos.data.map((pago) => (
                    <TableRow key={pago.operacion}>
                      <TableCell>{pago.concepto}</TableCell>
                      <TableCell>{money(pago.monto, pago.moneda)}</TableCell>
                      <TableCell>Fecha</TableCell>
                      <TableCell>
                        <Badge className="bg-green-100 text-green-700">
                          Completado
                        </Badge>
                      </TableCell>
                      <TableCell className="text-end">
                        <Button variant="outline">Ver detalles</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TabsContent>
            <TabsContent value="deudas">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="table__head">Cuota</TableHead>
                    <TableHead className="table__head">Estado</TableHead>
                    <TableHead className="table__head">Monto</TableHead>
                    <TableHead className="table__head">Deuda</TableHead>
                    <TableHead className="table__head text-end">
                      Acciones
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {deudas.data.map((deuda) => (
                    <TableRow key={deuda.id}>
                      <TableCell>
                        <Link
                          className="link"
                          href={"/cuotas/" + deuda.cuota.id}
                        >
                          {deuda.cuota.nombre}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <EstadoDeudaTag state={deuda.estado} />
                      </TableCell>
                      <TableCell>{money(deuda.monto)}</TableCell>
                      <TableCell>{money(deuda.deuda)}</TableCell>
                      <TableCell className="text-end">
                        <Button variant="outline">Ver detalles</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TabsContent>
            <TabsContent value="documentos">
              {unidad.titulares && (
                <section>
                  <h3>Titulares</h3>
                  <ul>
                    {unidad.titulares.map((titular) => (
                      <li key={titular.cedula}>{titular.display_name}</li>
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
