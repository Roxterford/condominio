import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Titular } from "../schemas";
import { Button } from "@/components/ui/button";
import { Ellipsis } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";

export interface VillasTableData {
  codigo: string;
  propietario: Titular;
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
          <TableHead className="table__head">Villa</TableHead>
          <TableHead className="table__head">Propietario</TableHead>
          <TableHead className="table__head">Estado de pagos</TableHead>
          <TableHead className="table__head">Contacto</TableHead>
          <TableHead className="table__head">Acciones</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map((villa) => (
          <TableRow key={villa.codigo}>
            <TableHead>{villa.codigo}</TableHead>
            <TableHead>{villa.propietario.nombre}</TableHead>
            <TableHead>
              <Badge>{villa.estado_pagos}</Badge>
            </TableHead>

            <TableHead>
              <p>{villa.contacto.telefono}</p>
              <span className="text-gray-500">{villa.contacto.email}</span>
            </TableHead>
            <TableHead>
              <div className="flex gap-2">
                <Button variant="outline" asChild>
                  <Link href={["/villas", villa.codigo].join("/")}>
                    Ver detalles
                  </Link>
                </Button>
                <Button variant="ghost">
                  <Ellipsis />
                </Button>
              </div>
            </TableHead>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
