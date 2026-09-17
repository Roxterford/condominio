import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { money } from "@/lib/money-display";
import { Moneda } from "@/providers/graphql/graphql";
import { DesgloseDeGastoItem } from "../desglose_de_gastos";
import type { EstrategiaDeDistribucion } from "./schema";

export interface ResumenDeCalculoProps {
  gastos: DesgloseDeGastoItem[];
  estrategia: EstrategiaDeDistribucion;
  unidadesActivas?: number;
  isLoadingUnidades?: boolean;
}

export function ResumenDeCalculo({
  gastos,
  estrategia,
  unidadesActivas,
  isLoadingUnidades,
}: ResumenDeCalculoProps) {
  const totalUsd = gastos.reduce((acc, g) => acc + g.total, 0);
  const montoPorVilla = unidadesActivas ? totalUsd / unidadesActivas : 0;

  return (
    <section>
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <CardTitle>
              <h3>Resumen de cálculo</h3>
            </CardTitle>
          </div>
          <CardDescription>
            Total a distribuir entre las villas según la estrategia{" "}
            <span className="font-medium capitalize text-foreground">
              {estrategia}
            </span>
          </CardDescription>
        </CardHeader>
        <CardContent className="gap-5">
          <ResumenRow label="Gastos" value={money(totalUsd, Moneda.Usd)} />
          <ResumenRow label="Fondo de reserva" value={money(0, Moneda.Usd)} />
          <ResumenRow
            label="Villas / Unidades"
            value={isLoadingUnidades ? "…" : String(unidadesActivas ?? 0)}
          />
          <Separator />
          {estrategia === "lineal" && (
            <div className="flex justify-between items-baseline gap-4">
              <span className="font-medium text-lg">Monto por villa</span>
              <span className="text-xl font-bold tabular-nums">
                {isLoadingUnidades ? "…" : money(montoPorVilla, Moneda.Usd)}
              </span>
            </div>
          )}
        </CardContent>
      </Card>
    </section>
  );
}

function ResumenRow({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div className="flex justify-between items-baseline gap-4">
      <span className="">{label}</span>
      <span className="text-sm font-semibold tabular-nums">{value}</span>
    </div>
  );
}