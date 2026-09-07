import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { money } from "@/lib/money-display";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import {
  OperacionesPageQuery,
  OperacionType,
} from "@/providers/graphql/graphql";
import { renderGraphql } from "@/providers/graphql/render";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import {
  ArrowLeftRight,
  CreditCardMinus,
  MoveDownRight,
  MoveUpRight,
  Plus,
} from "lucide-react";

const PageQuery = graphql(/* GraphQL */ `
  query OperacionesPage {
    operaciones: obtenerOperaciones(paginador: { limit: 20, page: 1 }) {
      data {
        __typename

        ... on IOperacion {
          fecha
          operacion
          concepto
          metodo
          moneda
          moneda
          total
        }

        ... on GastoAProveedor {
          proveedor {
            nombre
          }
        }

        ... on Pago {
          unidad {
            codigo
          }
        }
      }
    }
  }
`);

export default async function OperacionesPage() {
  return renderGraphql(await execute(PageQuery), OperacionesPageContent);
}

function OperacionesPageContent(page: OperacionesPageQuery) {
  return (
    <>
      <header className="flex items-end justify-between">
        <div>
          <h1>Operaciones</h1>
          <p className="page-description">
            Finanzas · Movimientos, pagos y gastos
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <CreditCardMinus /> Gasto
          </Button>
          <Button variant="outline">
            <ArrowLeftRight /> Transacción
          </Button>
          <Button>
            <Plus /> Pago
          </Button>
        </div>
      </header>
      <section>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="table__head">Referencia</TableHead>
              <TableHead className="table__head | text-center">Tipo</TableHead>
              <TableHead className="table__head">Concepto</TableHead>
              <TableHead className="table__head">Monto</TableHead>
              <TableHead className="table__head">Método</TableHead>
              <TableHead className="table__head">Fecha</TableHead>
              <TableHead className="table__head"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {page.operaciones.data.map((operacion) => (
              <TableRow>
                <TableCell>
                  <div className="flex gap-3">
                    <VarianteDeOperacionIcon variant={operacion.__typename} />

                    <VarianteDeOperacionReferencia data={operacion} />
                  </div>
                </TableCell>
                <TableCell className="text-center">
                  <VarianteDeOperacionTag variant={operacion.__typename} />
                </TableCell>
                <TableCell>{operacion.concepto}</TableCell>
                <TableCell>{money(operacion.total)}</TableCell>
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
            ))}
          </TableBody>
        </Table>
      </section>
    </>
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
