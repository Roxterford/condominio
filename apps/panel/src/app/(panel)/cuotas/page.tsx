import { Button } from "@/components/ui/button";
import { CuotasTableData } from "@/features/administracion/components/cuotas_table/cuotas_table";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import { Plus } from "lucide-react";
import Link from "next/link";
import { CuotasPageTaps } from "./components/cuotas-page-taps";
import { TipoDeCuota } from "@/providers/graphql/graphql";

const PageQuery = graphql(/* GraphQL */ `
  query CuotasPage {
    cuotas: obtenerCuotas {
      data {
        __typename
        ... on Cuota {
          id
          monto
          mes
          anio
          registro
          recaudacion {
            unidades_aplicadas
            pagos_asociados
          }
        }

        ... on CuotaEspecial {
          detalles {
            titulo
            descripcion
          }
        }
      }
    }
  }
`);

export default async function CuotasPage() {
  return renderGraphql(await execute(PageQuery), ({ cuotas }) => {
    const cuota_table_data = cuotas.data.map<CuotasTableData>((c) => ({
    id: c.id,
    tipo:
      c.__typename === "CuotaEspecial"
        ? TipoDeCuota.Especial
        : TipoDeCuota.Regular,
    monto: c.monto,
    mes: c.mes,
    anio: c.anio,
    registro: new Date(c.registro),
    actualizacion: new Date(),
    pagos_recibidos: c.recaudacion.pagos_asociados,
    pagos_esperados: c.recaudacion.unidades_aplicadas,
    ...((c.__typename === "CuotaEspecial" &&
      ({
        detalles: {
          titulo: c.detalles.titulo,
          descripcion: c.detalles.descripcion,
        },
      } as Pick<CuotasTableData<"especial">, "detalles">)) as any),
  }));

  return (
    <>
      <header className="flex items-center justify-between">
        <div>
          <h1 className="page-title">Cuotas</h1>
          <p className="page-description">
            Gestión de cuotas mensuales y espaciales
          </p>
        </div>

        <Link href="/cuotas/registrar">
          <Button>
            <Plus /> Nueva cuota
          </Button>
        </Link>
      </header>
      <CuotasPageTaps cuotas={cuota_table_data} />
    </>
  );
  });
}
