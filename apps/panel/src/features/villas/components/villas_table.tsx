import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Propietario } from "../schemas";

export interface VillasTableData {
  villa: number;
  propietario: Pick<
    Propietario,
    "nombres" | "apellidos" | "email" | "telefono"
  >;
  estado_pagos: "solvente" | "pendiente";
  contacto: string;
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
        {data.map((data) => (
          <TableRow key={Math.random()}>
            <TableHead>{Math.random()}</TableHead>
            <TableHead>{Math.random()}</TableHead>
            <TableHead>{Math.random()}</TableHead>
            <TableHead>{Math.random()}</TableHead>
            <TableHead>{Math.random()}</TableHead>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
