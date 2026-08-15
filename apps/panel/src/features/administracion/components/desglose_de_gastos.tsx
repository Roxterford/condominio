import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
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
import { Gasto, Proveedor } from "@/providers/graphql/graphql";
import { MoreHorizontal } from "lucide-react";

export interface DesgloseDeGastoItem extends Pick<
  Gasto,
  "operacion" | "concepto" | "moneda" | "monto" | "fecha" | "tasa" | "total"
> {
  proveedor: Pick<Proveedor, "id" | "nombre" | "rif" | "telefono" | "email">;
}
export interface DesgloseDeGastosProps {
  data: DesgloseDeGastoItem[];
  showActions?: boolean;
  onGastoPress?: (gasto: DesgloseDeGastoItem) => void;
  onRemove?: (gasto: DesgloseDeGastoItem) => void;
}

export function DesgloseDeGastos({
  data: gastos,
  onGastoPress: onGastoClick,
  showActions = true,
  onRemove,
}: DesgloseDeGastosProps) {
  const remove = (gasto: DesgloseDeGastoItem) => {
    onRemove?.(gasto);
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">ID</TableHead>
          <TableHead className="table__head">Concepto</TableHead>
          <TableHead className="table__head">Monto</TableHead>
          <TableHead className="table__head">Fecha</TableHead>
          <TableHead className="table__head">Proveedor</TableHead>
          {showActions && <TableHead className="text-right">Actions</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {gastos.map((gasto) => (
          <TableRow key={gasto.operacion}>
            <TableCell className="font-medium">
              <label className="link" onClick={() => onGastoClick?.(gasto)}>
                {gasto.operacion.slice(-6)}
              </label>
            </TableCell>
            <TableCell>{gasto.concepto}</TableCell>
            <TableCell>$ {gasto.total}</TableCell>
            <TableCell>{gasto.fecha.toLocaleDateString("es")}</TableCell>
            <TableCell>{gasto.proveedor.nombre}</TableCell>
            {showActions && (
              <TableCell className="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="size-8">
                      <MoreHorizontal />
                      <span className="sr-only">Open menu</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      variant="destructive"
                      onClick={() => remove(gasto)}
                    >
                      Remover
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            )}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
