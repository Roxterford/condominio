"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Empty,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
  EmptyDescription,
} from "@/components/ui/empty";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { money } from "@/lib/money-display";
import {
  OperacionesPageQuery,
  OperacionType,
} from "@/providers/graphql/graphql";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import { MoveDownRight, MoveUpRight, SearchX } from "lucide-react";

export function OperacionesTable({
  data,
}: {
  data: OperacionesPageQuery["operaciones"]["data"];
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">Concepto</TableHead>
          <TableHead className="table__head">Referencia</TableHead>
          <TableHead className="table__head | text-center">Tipo</TableHead>
          <TableHead className="table__head">Monto</TableHead>
          <TableHead className="table__head">Método</TableHead>
          <TableHead className="table__head">Fecha</TableHead>
          <TableHead className="table__head"></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.length ? (
          data.map((operacion) => (
            <TableRow key={operacion.operacion}>
              <TableCell>
                <div className="flex gap-3">
                  <VarianteDeOperacionIcon variant={operacion.__typename} />
                  {operacion.concepto}
                </div>
              </TableCell>
              <TableCell>
                <VarianteDeOperacionReferencia data={operacion} />
              </TableCell>
              <TableCell className="text-center">
                <VarianteDeOperacionTag variant={operacion.__typename} />
              </TableCell>
              <TableCell>{money(operacion.monto, operacion.moneda)}</TableCell>
              <TableCell>{operacion.metodo}</TableCell>
              <TableCell>
                {format(operacion.fecha, "d 'de' MMMM 'de' yyyy", {
                  locale: es,
                })}
              </TableCell>
              <TableCell className="text-end">
                <Button variant="outline">Ver detalles</Button>
              </TableCell>
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell colSpan={7}>
              <NotFoundState />
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}

function VarianteDeOperacionReferencia({
  data,
}: {
  data: OperacionesPageQuery["operaciones"]["data"][0];
}) {
  let maintext = "-";

  if (data.__typename === "GastoACondominio") maintext = "Condominio";
  if (data.__typename === "GastoAProveedor") maintext = data.proveedor.nombre;
  if (data.__typename === "Pago") maintext = `Unidad ${data.unidad.codigo}`;

  return (
    <div>
      <p className="font-medium">{maintext}</p>
      <p className="text-sm text-gray-500">{data.operacion}</p>
    </div>
  );
}

function VarianteDeOperacionIcon({
  variant,
}: {
  variant: OperacionType["__typename"];
}) {
  switch (variant) {
    case "GastoACondominio":
    case "GastoAProveedor":
      return (
        <Badge className="bg-yellow-200 text-yellow-700">
          <MoveUpRight />
        </Badge>
      );
    case "Pago":
      return (
        <Badge className="bg-green-100 text-green-700">
          <MoveDownRight />
        </Badge>
      );
  }

  return <Badge>Indeterminado</Badge>;
}

function VarianteDeOperacionTag({
  variant,
}: {
  variant: OperacionType["__typename"];
}) {
  switch (variant) {
    case "GastoACondominio":
    case "GastoAProveedor":
      return <Badge className="bg-yellow-200 text-yellow-700">Gasto</Badge>;
    case "Pago":
      return <Badge className="bg-green-100 text-green-700">Pago</Badge>;
  }

  return <Badge>Indeterminado</Badge>;
}

function NotFoundState() {
  return (
    <Empty>
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <SearchX />
        </EmptyMedia>
        <EmptyTitle>Resultados no enctrados</EmptyTitle>
        <EmptyDescription>Intenta una busqueda diferente</EmptyDescription>
      </EmptyHeader>
    </Empty>
  );
}
