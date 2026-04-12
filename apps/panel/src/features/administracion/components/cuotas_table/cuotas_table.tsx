import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import Link from "next/link";
import { Cuota, Proyecto, TipoDeCuota } from "../../schemas/cuota.schema";

export type CuotasTableType = "regular" | "especial" | "default";
export interface CuotasTableData<
  T extends CuotasTableType = "default",
> extends Pick<
  Cuota,
  | "id"
  | "tipo"
  | "monto"
  | "mes"
  | "anio"
  | "registro"
  | "actualizacion"
  | "tipo"
> {
  detalles: T extends "regular"
    ? never
    : T extends "especial"
      ? Pick<Proyecto, "titulo" | "descripcion">
      : never;
}

export interface CuotasTableProps {
  type?: CuotasTableType;
  data: CuotasTableData[];
}

export function CuotasTable({
  data: cuotas,
  type = "default",
}: CuotasTableProps) {
  return (
    <Table>
      <TableHeader>
        {type === "default" ? (
          <CuotaHeaders />
        ) : type === "regular" ? (
          <CuotaRegularHeaders />
        ) : (
          <CuotaEspecialHeaders />
        )}
      </TableHeader>
      <TableBody>
        {cuotas
          .filter((cuota) => {
            if (type === "default") return true;
            if (type === "regular") return cuota.tipo === TipoDeCuota.REGULAR;
            return cuota.tipo === TipoDeCuota.ESPECIAL;
          })
          .map((cuota) =>
            type === "default" ? (
              <CuotaRow key={cuota.id} cuota={cuota} />
            ) : type === "regular" ? (
              <CuotaRegularRow key={cuota.id} cuota={cuota} />
            ) : (
              <CuotaEspecialRow key={cuota.id} cuota={cuota as any} />
            ),
          )}
      </TableBody>
    </Table>
  );
}

function CuotaHeaders() {
  return (
    <TableRow>
      <TableHead className="table__head">Cuota</TableHead>
      <TableHead className="table__head">Tipo</TableHead>
      <TableHead className="table__head">Monto</TableHead>
      <TableHead className="table__head">Fecha límite</TableHead>
      <TableHead className="table__head">Estado</TableHead>
      <TableHead className="table__head">Pagos recibidos</TableHead>
      <TableHead className="table__head text-right">Acciones</TableHead>
    </TableRow>
  );
}

function CuotaRow({ cuota }: { cuota: CuotasTableData }) {
  return (
    <TableRow key={`${cuota.mes}-${cuota.anio}`}>
      <TableCell className="font-medium">{cuota.id}</TableCell>
      <TableCell>{cuota.tipo}</TableCell>
      <TableCell>{cuota.monto}</TableCell>
      <TableCell>
        {cuota.mes}/{cuota.anio}
      </TableCell>
      <TableCell className="text-right">
        {cuota.actualizacion?.toLocaleDateString("es-VE")}
      </TableCell>
      <TableCell className="text-center">-</TableCell>
      <TableCell className="text-right">
        <Button variant="outline">Ver detalles</Button>
      </TableCell>
    </TableRow>
  );
}

function CuotaRegularHeaders() {
  return (
    <TableRow>
      <TableHead className="table__head">Periodo</TableHead>
      <TableHead className="table__head">Monto</TableHead>
      <TableHead className="table__head">Fecha límite</TableHead>
      <TableHead className="table__head">Estado</TableHead>
      <TableHead className="table__head">Pagos recibidos</TableHead>
      <TableHead className="table__head text-right">Acciones</TableHead>
    </TableRow>
  );
}

function CuotaRegularRow({ cuota }: { cuota: CuotasTableData }) {
  return (
    <TableRow key={`${cuota.mes}-${cuota.anio}`}>
      <TableCell className="font-medium">{cuota.id}</TableCell>
      <TableCell>{cuota.monto}</TableCell>
      <TableCell>
        {cuota.mes}/{cuota.anio}
      </TableCell>
      <TableCell className="text-right">
        {cuota.actualizacion?.toLocaleDateString("es-VE")}
      </TableCell>
      <TableCell className="text-center">-</TableCell>
      <TableCell className="text-right">
        <Button variant="outline">Ver detalles</Button>
      </TableCell>
    </TableRow>
  );
}

function CuotaEspecialHeaders() {
  return (
    <TableRow>
      <TableHead className="table__head">Descipción</TableHead>
      <TableHead className="table__head">Monto</TableHead>
      <TableHead className="table__head">Fecha límite</TableHead>
      <TableHead className="table__head">Estado</TableHead>
      <TableHead className="table__head">Pagos recibidos</TableHead>
      <TableHead className="table__head text-right">Acciones</TableHead>
    </TableRow>
  );
}

function CuotaEspecialRow({ cuota }: { cuota: CuotasTableData<"especial"> }) {
  return (
    <TableRow key={`${cuota.mes}-${cuota.anio}`}>
      <TableCell className="font-medium">{cuota.detalles.titulo}</TableCell>
      <TableCell>{cuota.monto}</TableCell>
      <TableCell>
        {cuota.mes}/{cuota.anio}
      </TableCell>
      <TableCell className="text-right">
        {cuota.actualizacion?.toLocaleDateString("es-VE")}
      </TableCell>
      <TableCell className="text-center">-</TableCell>
      <TableCell className="text-right">
        <Link href={`/cuotas/${cuota.id}`}>Ver detalles</Link>
      </TableCell>
    </TableRow>
  );
}
