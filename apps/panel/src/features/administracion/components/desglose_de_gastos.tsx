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
import { Gasto } from "../schemas";

export type GastoItem = Pick<
  Gasto,
  "id" | "concepto" | "moneda" | "monto" | "fecha" | "proveedor"
>;
export interface DesgloseDeGastosProps {
  gastos: GastoItem[];
  onGastoPress?: (gasto: GastoItem) => void;
}

export function DesgloseDeGastos({
  gastos,
  onGastoPress: onGastoClick,
}: DesgloseDeGastosProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="font-semibold uppercase">ID</TableHead>
          <TableHead className="font-semibold uppercase">Concepto</TableHead>
          <TableHead className="font-semibold uppercase">Monto</TableHead>
          <TableHead className="font-semibold uppercase">Fecha</TableHead>
          <TableHead className="font-semibold uppercase">Proveedor</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {gastos.map((gasto) => (
          <TableRow key={gasto.id}>
            <TableCell className="font-medium">
              <label className="link" onClick={() => onGastoClick?.(gasto)}>
                {gasto.id.split("").reverse().join("").substring(0, 6)}
              </label>
            </TableCell>
            <TableCell>{gasto.concepto}</TableCell>
            <TableCell>$ {(gasto.monto / 100).toLocaleString("es")}</TableCell>
            <TableCell>{gasto.fecha.toLocaleDateString("es")}</TableCell>
            <TableCell>{gasto.proveedor}</TableCell>
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
