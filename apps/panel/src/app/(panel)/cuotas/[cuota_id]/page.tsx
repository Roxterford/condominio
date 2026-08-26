import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import styles from "./page.module.css";
import {
  DesgloseDeGastoItem,
  DesgloseDeGastos,
} from "@/features/administracion/components/desglose_de_gastos";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { EstadoDeudaTag } from "@/components/estado-deuda-tag";
import { EstadoDeDeuda } from "@/providers/graphql/graphql";
import Link from "next/link";

const PageQuery = graphql(/* GraphQL */ `
  query CuotaPage($cuota_id: String!) {
    cuota: obtenerCuota(id: $cuota_id) {
      __typename
      ... on Cuota {
        id
        mes
        anio
        monto
        gastos {
          __typename
          ... on Gasto {
            operacion
            concepto
            moneda
            monto
            fecha
            tasa
            total
          }
          ... on GastoAProveedor {
            proveedor {
              id
              nombre
              rif
              telefono
              email
            }
          }
        }
        recaudacion {
          unidades_aplicadas
          unidades_solventes
          monto_estimado
          monto_recaudado
        }
      }
      ... on CuotaEspecial {
        detalles {
          titulo
          descripcion
          justificacion
        }
      }
    }
    deudas: obtenerDeudas(filtro: { cuota: { eq: $cuota_id } }) {
      data {
        unidad {
          codigo
        }
        titular {
          display_name
        }
        deuda
        estado
      }
    }
  }
`);

export default async function CuotaPage({
  params,
}: {
  params: Promise<{ cuota_id: string }>;
}) {
  const { cuota_id } = await params;

  return renderGraphql(
    await execute(PageQuery, { cuota_id }),
    ({ cuota, deudas }) => {
      if (!cuota) {
        return <div>Not found</div>;
      }

      return (
        <>
          <header>
            <div>
              <h1>
                {cuota.__typename === "CuotaEspecial"
                  ? cuota.detalles.titulo
                  : `${cuota.mes} ${cuota.anio}`}
              </h1>
              <p className="page-description">
                Detalle de cuota{" "}
                {cuota.__typename === "CuotaEspecial" ? "especial" : "regular"}
              </p>
            </div>
          </header>
          <section className={[styles.infoboxes, "mt-10"].join(" ")}>
            <div className={styles.infobox}>
              <h3 className={styles.infobox__title}>Presupuesto Estimado</h3>
              <p className={styles.infobox__value}>
                ${cuota.recaudacion.monto_estimado.toLocaleString("es-VE")}
                {cuota.recaudacion.monto_estimado % 1 !== 0 ? "" : ",00"}
              </p>
            </div>
            <div className={styles.infoboxes__divider}></div>
            <div className={styles.infobox}>
              <h3 className={styles.infobox__title}>Monto por Villa</h3>
              <p className={styles.infobox__value}>
                ${cuota.monto.toLocaleString("es-VE")}
              </p>
            </div>
            <div className={styles.infoboxes__divider}></div>
            <div className={styles.infobox}>
              <h3 className={styles.infobox__title}>Fecha limite de pago</h3>
              <p className={styles.infobox__value}>2025-10-01</p>
            </div>
            <div className={styles.infoboxes__divider}></div>
            <div className={styles.infobox}>
              <h3 className={styles.infobox__title}>Estado</h3>
              <Badge className="bg-yellow-100 text-yellow-600">Pendiente</Badge>
            </div>
          </section>

          <section className="mt-10">
            <h2>Resumen de Recaudación</h2>
            <p className="page-description">
              Progreso de pagos recibidos para esta cuota{" "}
              {cuota.__typename === "CuotaEspecial" ? "especial" : "regular"}
            </p>

            <div className="flex gap-10 justify-between mt-10">
              <div className={styles.infobox}>
                <h3 className={styles.infobox__title}>Pagos Recibidos</h3>
                <p className={styles.infobox__value}>
                  {cuota.recaudacion.unidades_solventes}/
                  {cuota.recaudacion.unidades_aplicadas}
                </p>
              </div>
              <div className={styles.infobox}>
                <h3 className={styles.infobox__title}>Monto Recaudado</h3>
                <p className={styles.infobox__value}>
                  {cuota.recaudacion.monto_recaudado}
                </p>
              </div>
              <div className={styles.infobox}>
                <h3 className={styles.infobox__title}>Porcentaje</h3>
                <p className={styles.infobox__value}>
                  {(
                    (cuota.recaudacion.monto_recaudado /
                      cuota.recaudacion.monto_estimado) *
                    100
                  ).toFixed(2)}
                  %
                </p>
              </div>
            </div>
            <div className="text-end">
              <span className="text-muted-foreground">
                {(
                  (cuota.recaudacion.monto_recaudado /
                    cuota.recaudacion.monto_estimado) *
                  100
                ).toFixed(2)}
                %
              </span>
              <Progress
                value={
                  (cuota.recaudacion.monto_recaudado /
                    cuota.recaudacion.monto_estimado) *
                  100
                }
              />
            </div>

            <section className="space-y-8 mt-10">
              <h2 className="mb-2">Detalles del Proyecto</h2>
              {cuota.__typename === "CuotaEspecial" && (
                <>
                  <section>
                    <h3>Descripción</h3>
                    <p>{cuota.detalles?.descripcion}</p>
                  </section>
                  <section>
                    <h3>Justificación</h3>
                    <p>{cuota.detalles.justificacion}</p>
                  </section>
                </>
              )}

              <section>
                <h3>Documentos Adjuntos</h3>
                <ul>
                  <li></li>
                </ul>
              </section>
            </section>
          </section>

          <section>
            <h2>Información Adicional</h2>
            <section>
              <h3>Desglose de gastos</h3>
              <DesgloseDeGastos
                data={cuota.gastos.reduce<DesgloseDeGastoItem[]>(
                  (acc, it) =>
                    it.__typename !== "GastoAProveedor" ? acc : [...acc, it],
                  [],
                )}
                showActions={false}
              />
            </section>
          </section>

          <section>
            <h2>Estado de Pagos por Villa</h2>
            <p>Seguimiento de pagos de cada villa para esta cuota especial</p>

            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="table__head">Villa</TableHead>
                  <TableHead className="table__head">Propietario</TableHead>
                  <TableHead className="table__head">Estado</TableHead>
                  <TableHead className="table__head">Monto</TableHead>
                  <TableHead className="table__head text-right">
                    Acciones
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {deudas.data.map((deuda, i) => (
                  <TableRow key={i}>
                    <TableCell className="">
                      <Link
                        className="link"
                        href={`/villas/${deuda.unidad.codigo}`}
                      >
                        {deuda.unidad.codigo}
                      </Link>
                    </TableCell>
                    <TableCell className="">
                      {deuda.titular?.display_name}
                    </TableCell>
                    <TableCell className="">
                      <EstadoDeudaTag state={deuda.estado} />
                    </TableCell>
                    <TableCell className="">{deuda.deuda}</TableCell>
                    <TableCell className="text-right">
                      <Actions state={deuda.estado}></Actions>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </section>
        </>
      );
    },
  );
}

function Actions({ state }: { state: EstadoDeDeuda }) {
  if ([EstadoDeDeuda.Pendiente, EstadoDeDeuda.Saldada].includes(state))
    return (
      <>
        <Button variant="outline">Registrar pago</Button>
      </>
    );

  return (
    <>
      <Button variant="outline">Ver detalles</Button>
    </>
  );
}
