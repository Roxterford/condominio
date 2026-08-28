import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Ellipsis } from "lucide-react";
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
}
export function VillasTable({ data }: VillasTableProps) {
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
        {data.map((villa) => (
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
        ))}
      </TableBody>
    </Table>
  );
}

function Ignore() {
  return "-";
}
