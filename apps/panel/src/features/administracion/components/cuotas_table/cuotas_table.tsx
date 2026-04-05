import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Cuota } from "../../schemas/cuota.schema";

export interface CuotasTableData extends Pick<
  Cuota,
  "id" | "monto" | "mes" | "anio" | "registro" | "actualizacion"
> {}

interface CuotasTableProps {
  data: CuotasTableData[];
}

export function CuotasTable({ data: cuotas }: CuotasTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">Descipción</TableHead>
          <TableHead className="table__head">Monto</TableHead>
          <TableHead className="table__head">Fecha límite</TableHead>
          <TableHead className="table__head">Estado</TableHead>
          <TableHead className="table__head">Pagos recibidos</TableHead>
          <TableHead className="table__head text-right">Acciones</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {cuotas.map((cuota) => (
          <TableRow key={`${cuota.mes}-${cuota.anio}`}>
            <TableCell className="font-medium">{cuota.id}</TableCell>
            <TableCell>{cuota.monto}</TableCell>
            <TableCell>
              {cuota.mes}/{cuota.anio}
            </TableCell>
            <TableCell className="text-right">
              {cuota.actualizacion?.toLocaleDateString("es-VE")}
            </TableCell>
            <TableCell className="text-right">
              <Button>Ver detalles</Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
