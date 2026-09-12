import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Home, SearchX, Ellipsis } from "lucide-react";
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";
import { AvatarIniciales } from "@/components/avatar-iniciales/avatar-iniciales";
import { Titular, Unidad } from "@/providers/graphql/graphql";
import { EstadoUnidadTag } from "@/components/estado-unidad-tag";
import { money } from "@/lib/money-display";

export interface VillasTableData extends Pick<
  Unidad,
  "codigo" | "estado" | "wallet"
> {
  propietario?: Pick<Titular, "cedula" | "display_name">;
  estado_pagos: "solvente" | "pendiente";
  contacto: {
    email: string;
    telefono: string;
  };
}

export interface VillasTableProps {
  data: VillasTableData[];
  busqueda?: boolean;
  emptyTitle?: string;
  emptyDescription?: string;
}
export function VillasTable({
  data,
  busqueda,
  emptyTitle,
  emptyDescription,
}: VillasTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">Código</TableHead>
          <TableHead className="table__head">Titular primario</TableHead>
          <TableHead className="table__head">Contacto</TableHead>
          <TableHead className="table__head">Estado</TableHead>
          <TableHead className="table__head">Solvencia</TableHead>
          <TableHead className="table__head text-end">Deuda total</TableHead>
          <TableHead className="table__head"></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.length ? (
          data.map((villa) => (
            <TableRow key={villa.codigo}>
              <TableCell>
                <Link className="link" href={"/villas/" + villa.codigo}>
                  {villa.codigo}
                </Link>
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  {villa.propietario ? (
                    <>
                      <AvatarIniciales nombre={villa.propietario.display_name} />
                      <div className="grid">
                        <span>{villa.propietario.display_name}</span>
                        <span className="text-gray-500 text-xs">
                          {villa.propietario.cedula}
                        </span>
                      </div>
                    </>
                  ) : (
                    <Ignore />
                  )}
                </div>
              </TableCell>
              <TableCell>
                <p>{villa.contacto.telefono}</p>
                <span className="text-gray-500">{villa.contacto.email}</span>
              </TableCell>
              <TableCell>
                <EstadoUnidadTag state={villa.estado} />
              </TableCell>

              <TableCell>
                <Badge>{villa.estado_pagos}</Badge>
              </TableCell>

              <TableCell className="text-end">
                {villa.wallet > 0 ? money(villa.wallet) : <Ignore />}
              </TableCell>

              <TableCell>
                <div className="flex gap-2 justify-end">
                  <Button variant="outline" asChild>
                    <Link href={["/villas", villa.codigo].join("/")}>
                      Ver detalles
                    </Link>
                  </Button>
                  <Button variant="ghost">
                    <Ellipsis />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell colSpan={7}>
              {busqueda ? (
                <NotFoundState />
              ) : (
                <EmptyState
                  title={emptyTitle}
                  description={emptyDescription}
                />
              )}
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}

function Ignore() {
  return "-";
}

function EmptyState({
  title = "Lista vacía",
  description = "No hay unidades que mostrar aquí",
}: {
  title?: string;
  description?: string;
}) {
  return (
    <Empty>
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <Home />
        </EmptyMedia>
        <EmptyTitle>{title}</EmptyTitle>
        <EmptyDescription>{description}</EmptyDescription>
      </EmptyHeader>
    </Empty>
  );
}

function NotFoundState() {
  return (
    <Empty>
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <SearchX />
        </EmptyMedia>
        <EmptyTitle>Resultados no encontrados</EmptyTitle>
        <EmptyDescription>Intenta una búsqueda diferente</EmptyDescription>
      </EmptyHeader>
    </Empty>
  );
}
