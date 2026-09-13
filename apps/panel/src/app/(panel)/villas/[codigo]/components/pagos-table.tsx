"use client";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  OperacionDetalle,
  acortarId,
  type Operacion,
} from "@/components/operacion-detalle/operacion-detalle";
import { useDrawer } from "@/contexts/drawer-context";
import { money } from "@/lib/money-display";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import { ChevronRight, ReceiptText } from "lucide-react";
import { VillaPageQuery } from "@/providers/graphql/graphql";
import { RegistrarPagoButton } from "./registrar-pago-button";
import {
  Empty,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
  EmptyDescription,
  EmptyContent,
} from "@/components/ui/empty";

type OperacionDetalleProps = React.ComponentProps<typeof OperacionDetalle>;

export function PagosTable({
  pagos,
  unidadTitular,
  unidad,
}: {
  pagos: VillaPageQuery["pagos"]["data"];
  unidadTitular?: OperacionDetalleProps["unidadTitular"];
  unidad: VillaPageQuery["unidad"];
}) {
  const { open } = useDrawer();

  const verDetalles = (pago: (typeof pagos)[number]) => {
    open({
      content: (
        <OperacionDetalle
          operacion={pago as Operacion}
          unidadTitular={unidadTitular}
        />
      ),
      title: "Información de la operación",
      titleBadge: (
        <Badge
          variant="secondary"
          className="font-mono text-[10px] font-medium tracking-wide"
        >
          {acortarId(pago.operacion)}
        </Badge>
      ),
      side: "right",
      size: 440,
    });
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">Concepto</TableHead>
          <TableHead className="table__head">Monto</TableHead>
          <TableHead className="table__head">Fecha</TableHead>
          <TableHead className="table__head">Estado</TableHead>
          <TableHead className="table__head text-end">Acciones</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {pagos.length ? (
          pagos.map((pago) => (
            <TableRow
              key={pago.operacion}
              onClick={() => verDetalles(pago)}
              className="cursor-pointer"
            >
              <TableCell>{pago.concepto}</TableCell>
              <TableCell>{money(pago.monto, pago.moneda)}</TableCell>
              <TableCell>
                {format(pago.fecha, "d MMM yyyy", { locale: es })}
              </TableCell>
              <TableCell>
                <Badge className="bg-green-100 text-green-700">
                  Completado
                </Badge>
              </TableCell>
              <TableCell className="text-end">
                <ChevronRight
                  className="ml-auto text-muted-foreground"
                  size={16}
                />
              </TableCell>
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell colSpan={5}>
              <EmptyState unidad={unidad} />
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}

function EmptyState({ unidad }: { unidad: VillaPageQuery["unidad"] }) {
  return (
    <Empty>
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <ReceiptText />
        </EmptyMedia>
        <EmptyTitle>No hay pagos</EmptyTitle>
        <EmptyDescription>Registra un pago para verlo aquí</EmptyDescription>
      </EmptyHeader>
      <EmptyContent className="flex-row justify-center gap-2">
        <RegistrarPagoButton unidad={unidad} />
      </EmptyContent>
    </Empty>
  );
}
