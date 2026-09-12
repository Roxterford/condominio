"use client";

import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  CuotaEspecial,
  CuotaRegular,
  Proyecto,
  TipoDeCuota,
} from "@/providers/graphql/graphql";
import { useDrawer } from "@/contexts/drawer-context";
import { execute } from "@/providers/graphql/execute";
import { useRouter } from "next/navigation";
import { useSingleDoubleClick } from "@/hooks/useSingleDoubleClick";
import {
  CuotaDetalle,
  CuotaDetalleQuery,
} from "@/components/cuota-detalle/cuota-detalle";
import { TipoCuotaTag } from "@/components/tipo-cuota-tag";
import { ChevronRight } from "lucide-react";
import type { ReactNode } from "react";

export type CuotasTableType = "regular" | "especial" | "default";

export interface CuotasTableData<
  T extends CuotasTableType = "default",
> extends Pick<
  CuotaEspecial | CuotaRegular,
  "id" | "__typename" | "monto" | "mes" | "anio" | "registro" | "actualizacion"
> {
  detalles: T extends "regular"
    ? never
    : T extends "especial"
      ? Pick<Proyecto, "titulo" | "descripcion">
      : never;
  pagos_recibidos: number;
  pagos_esperados: number;
}

export interface CuotasTableProps {
  type?: CuotasTableType;
  data: CuotasTableData[];
}

type ColumnConfig<T extends CuotasTableType = "default"> = {
  label: string;
  getValue: (cuota: CuotasTableData<T>) => ReactNode;
  className?: string;
};

const createColumns = <T extends CuotasTableType>(
  firstColLabel: string,
  getFirstValue: (cuota: CuotasTableData<T>) => ReactNode,
  firstClassName?: string,
): ColumnConfig<T>[] => [
  { label: firstColLabel, getValue: getFirstValue, className: firstClassName },
  { label: "Monto", getValue: (c) => "$" + c.monto.toLocaleString("es-VE") },
  { label: "Fecha límite", getValue: (c) => `${c.mes}/${c.anio}` },
  {
    label: "Estado",
    getValue: (c) => c.actualizacion?.toLocaleDateString("es-VE") ?? "-",
    className: "text-right",
  },
  {
    label: "Pagos recibidos",
    getValue: (c) => (
      <>
        {`${c.pagos_recibidos}/${c.pagos_esperados}`}

        <Progress value={(c.pagos_recibidos / c.pagos_esperados) * 100} />
      </>
    ),
    className: "text-center",
  },
];

const COLUMNS = createColumns("Cuota", (c) => c.id, "font-medium");

const SPECIAL_COLUMNS = createColumns<"especial">(
  "Descripción",
  (c) => c.detalles?.titulo ?? "-",
  "font-medium",
);

function tipoDeCuota(cuota: CuotasTableData): TipoDeCuota {
  return cuota.__typename === "CuotaEspecial"
    ? TipoDeCuota.Especial
    : TipoDeCuota.Regular;
}

export function CuotasTable({
  data: cuotas,
  type = "default",
}: CuotasTableProps) {
  const { open } = useDrawer();
  const router = useRouter();
  const columns = type === "especial" ? SPECIAL_COLUMNS : COLUMNS;

  const filteredCuotas = cuotas.filter((cuota) => {
    if (type === "default") return true;
    if (type === "regular") return cuota.__typename === "CuotaRegular";
    return cuota.__typename === "CuotaEspecial";
  });

  const verDetalles = (cuota: CuotasTableData) => {
    open({
      title: "Información de la cuota",
      titleBadge: <TipoCuotaTag type={tipoDeCuota(cuota)} />,
      side: "right",
      size: 460,
      loader: async () => {
        const result = await execute(CuotaDetalleQuery, {
          cuota_id: cuota.id,
        });
        const data = result.data;
        if (!data?.cuota) {
          return (
            <p className="p-6 text-sm text-muted-foreground">No encontrada</p>
          );
        }
        return (
          <CuotaDetalle
            cuota={data.cuota}
            deudasPendientes={data.deudas.data}
          />
        );
      },
    });
  };

  const onClickFila = useSingleDoubleClick({
    onSingle: verDetalles,
    onDouble: (cuota) => router.push(`/cuotas/${cuota.id}`),
  });

  return (
    <Table>
      <TableHeader>
        <TableRow>
          {columns.map((col) => (
            <TableHead
              key={col.label}
              className={`table__head ${col.className ?? ""}`}
            >
              {col.label}
            </TableHead>
          ))}
          <TableHead className="table__head text-right">Acciones</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {filteredCuotas.map((cuota) => (
          <TableRow
            key={`${cuota.mes}-${cuota.anio}`}
            onClick={(e) => onClickFila(cuota, e)}
            className="cursor-pointer"
          >
            {columns.map((col) => (
              <TableCell key={col.label} className={col.className}>
                {col.getValue(cuota as CuotasTableData<typeof type>)}
              </TableCell>
            ))}
            <TableCell className="text-right">
              <Button
                variant="outline"
                size="sm"
                onClick={(e) => {
                  e.stopPropagation();
                  verDetalles(cuota);
                }}
              >
                Ver detalles
                <ChevronRight className="ml-1 size-4" />
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
