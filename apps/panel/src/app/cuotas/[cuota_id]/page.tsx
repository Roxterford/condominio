import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import styles from "./page.module.css";

const PageQuery = graphql(`
  query CuotaPage($cuota_id: String!) {
    cuota: obtenerCuota(id: $cuota_id) {
      __typename
      ... on Cuota {
        id
        mes
        anio
        pagos {
          monto_total
          total
        }
      }
      ... on CuotaEspecial {
        detalles {
          titulo
        }
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

  const {
    data: { cuota },
    errors,
  } = await execute(PageQuery, { cuota_id });

  if (errors) {
    return (
      <>
        <h1>Error</h1>
        <pre>{JSON.stringify(errors, null, 2)}</pre>
      </>
    );
  }

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
      <section className={styles.infoboxes}>
        <div className={styles.infobox}>
          <h3 className={styles.infobox__title}>Presupuesto Estimado</h3>
          <p className={styles.infobox__value}>{cuota.pagos.monto_total}</p>
        </div>
        <div className={styles.infoboxes__divider}></div>
        <div className={styles.infobox}>
          <h3 className={styles.infobox__title}>Monto por Villa</h3>
          <p className={styles.infobox__value}>$100.000</p>
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

      <section>
        <h2>Resumen de Recaudación</h2>
        <p className="page-description">
          Progreso de pagos recibidos para esta cuota{" "}
          {cuota.__typename === "CuotaEspecial" ? "especial" : "regular"}
        </p>

        <div className="flex gap-10 justify-between">
          <div className={styles.infobox}>
            <h3 className={styles.infobox__title}>Pagos Recibidos</h3>
            <p className={styles.infobox__value}>260/346</p>
          </div>
          <div className={styles.infobox}>
            <h3 className={styles.infobox__title}>Monto Recaudado</h3>
            <p className={styles.infobox__value}>$2.340,00</p>
          </div>
          <div className={styles.infobox}>
            <h3 className={styles.infobox__title}>Porcentaje</h3>
            <p className={styles.infobox__value}>75%</p>
          </div>
        </div>
        <div className="text-end">
          <span className="text-muted-foreground">75%</span>
          <Progress value={75} />
        </div>

        <section>
          <h2>Detalles del Proyecto</h2>
          <section>
            <h3>Descripción</h3>
            <p>
              Este gasto corresponde a la reparación urgente del sistema de
              bombeo de agua, cuyo deterioro ha ocasionado interrupciones en el
              suministro dentro de la urbanización. La falla detectada
              compromete el acceso regular al agua potable, lo que afecta la
              calidad de vida de los residentes y genera gastos adicionales en
              alternativas temporales. La inversión incluye la compra de
              repuestos, mantenimiento técnico y pruebas de funcionamiento para
              garantizar la operatividad eficiente del sistema.
            </p>
          </section>
          <section>
            <h3>Justificación</h3>
            <p>
              El sistema de bombeo de agua es una infraestructura esencial para
              el condominio, y su reparación es una necesidad inmediata para
              evitar costos mayores en el futuro.
            </p>
          </section>

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
        </section>
      </section>

      <section>
        <h2>Estado de Pagos por Villa</h2>
        <p>Seguimiento de pagos de cada villa para esta cuota especial</p>
      </section>
    </>
  );
}
