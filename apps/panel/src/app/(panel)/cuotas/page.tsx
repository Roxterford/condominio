import { Button } from "@/components/ui/button";
import { CuotasTableData } from "@/features/administracion/components/cuotas_table/cuotas_table";
import { TipoDeCuota } from "@/features/administracion/schemas/cuota.schema";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { Plus } from "lucide-react";
import Link from "next/link";
import { CuotasPageTaps } from "./components/cuotas-page-taps";

const PageQuery = graphql(`
  query CuotasPage {
    cuotas: obtenerCuotas {
      data {
        __typename
        ... on CuotaRegular {
          id
          monto
          mes
          anio
          registro
        }

        ... on CuotaEspecial {
          id
          monto
          mes
          anio
          registro
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
  const {
    data: { cuotas },
  } = await execute(PageQuery);

  const cuota_table_data = cuotas.data.map<CuotasTableData>((c) => ({
    id: c.id,
    tipo:
      c.__typename === "CuotaEspecial"
        ? TipoDeCuota.ESPECIAL
        : TipoDeCuota.REGULAR,
    monto: c.monto,
    mes: c.mes,
    anio: c.anio,
    registro: new Date(c.registro),
    actualizacion: new Date(),
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
}
