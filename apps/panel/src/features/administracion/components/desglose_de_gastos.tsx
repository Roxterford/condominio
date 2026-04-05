import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { MoreHorizontal } from "lucide-react";
import { Gasto, Proveedor } from "../schemas";

export interface DesgloseDeGastosData extends Pick<
  Gasto,
  "id" | "concepto" | "moneda" | "monto" | "fecha" | "tasa"
> {
  proveedor: Pick<Proveedor, "id" | "nombre" | "rif" | "telefono" | "email">;
}
export interface DesgloseDeGastosProps {
  data: DesgloseDeGastosData[];
  onGastoPress?: (gasto: DesgloseDeGastosData) => void;
}

export function DesgloseDeGastos({
  data: gastos,
  onGastoPress: onGastoClick,
}: DesgloseDeGastosProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">ID</TableHead>
          <TableHead className="table__head">Concepto</TableHead>
          <TableHead className="table__head">Monto</TableHead>
          <TableHead className="table__head">Fecha</TableHead>
          <TableHead className="table__head">Proveedor</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {gastos.map((gasto) => (
          <TableRow key={gasto.id}>
            <TableCell className="font-medium">
              <label className="link" onClick={() => onGastoClick?.(gasto)}>
                {gasto.id.slice(-6)}
              </label>
            </TableCell>
            <TableCell>{gasto.concepto}</TableCell>
            <TableCell>$ {(gasto.monto / 100).toLocaleString("es")}</TableCell>
            <TableCell>{gasto.fecha.toLocaleDateString("es")}</TableCell>
            <TableCell>{gasto.proveedor.nombre}</TableCell>
            <TableCell className="text-right">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon" className="size-8">
                    <MoreHorizontal />
                    <span className="sr-only">Open menu</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem>Edit</DropdownMenuItem>
                  <DropdownMenuItem>Duplicate</DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem variant="destructive">
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
