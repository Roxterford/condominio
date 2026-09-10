"use client";

import { AvatarIniciales } from "@/components/avatar-iniciales/avatar-iniciales";
import { Badge } from "@/components/ui/badge";
import { money } from "@/lib/money-display";
import { OperacionesPageQuery } from "@/providers/graphql/graphql";
import { format, formatDistanceToNow } from "date-fns";
import { es } from "date-fns/locale";

type Operacion = OperacionesPageQuery["operaciones"]["data"][number];
type OperacionConActualizacion = Operacion & { actualizado_en?: Date };

export function OperacionDetalle({ operacion }: { operacion: Operacion }) {
  return (
    <div className="space-y-6 p-6">
      <TagSection operacion={operacion} />
      <ResumenBoxes operacion={operacion} />
      <CajasDeDetalle operacion={operacion} />
      {operacion.__typename === "GastoAProveedor" && (
        <ProveedorSection proveedor={operacion.proveedor} />
      )}
    </div>
  );
}

function TagSection({ operacion }: { operacion: Operacion }) {
  if (operacion.__typename === "Pago") {
    return (
      <div className="flex flex-wrap gap-2">
        <Badge className="bg-green-100 text-green-700">Pago</Badge>
        <Badge variant="secondary">Unidad {operacion.unidad.codigo}</Badge>
      </div>
    );
  }

  if (operacion.__typename === "GastoAProveedor") {
    return (
      <div className="flex flex-wrap gap-2">
        <Badge className="bg-yellow-200 text-yellow-700">Gasto</Badge>
        <Badge variant="secondary">{operacion.proveedor.nombre}</Badge>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap gap-2">
      <Badge className="bg-yellow-200 text-yellow-700">Gasto</Badge>
      <Badge variant="secondary">Condominio</Badge>
    </div>
  );
}

function ResumenBoxes({ operacion }: { operacion: Operacion }) {
  const positivo = operacion.__typename === "Pago";

  return (
    <div className="grid grid-cols-3 divide-x divide-border">
      <Box titulo="Total">
        <p
          className={[
            "text-sm font-semibold",
            positivo ? "text-green-700" : "text-yellow-700",
          ].join(" ")}
        >
          {money(operacion.total, operacion.moneda)}
        </p>
      </Box>
      <Box titulo="Tasa">
        <p className="text-sm font-semibold">Bs. {formatearTasa(operacion.tasa)}</p>
      </Box>
      <Box titulo="Fecha">
        <p className="text-sm font-semibold">
          {format(operacion.fecha, "d MMM", { locale: es })}
        </p>
        <p className="text-xs text-muted-foreground">
          {formatDistanceToNow(operacion.fecha, {
            addSuffix: true,
            locale: es,
          })}
        </p>
      </Box>
    </div>
  );
}

function CajasDeDetalle({ operacion }: { operacion: Operacion }) {
  const actualizadoEn = (operacion as OperacionConActualizacion).actualizado_en;

  return (
<div className="space-y-4">
        <Box titulo="Concepto">
          <p className="text-sm font-medium text-foreground">{operacion.concepto}</p>
        </Box>

        <div className="grid grid-cols-2 divide-x divide-border">
          <Box titulo="Monto">
            <p
              className={[
                "text-sm font-semibold",
                operacion.__typename === "Pago"
                  ? "text-green-700"
                  : "text-yellow-700",
              ].join(" ")}
            >
              {money(operacion.monto, operacion.moneda)}
            </p>
          </Box>
          <Box titulo="Registro">
            <p className="text-sm font-medium">
              {format(operacion.registro, "d MMM yyyy '·' HH:mm", {
                locale: es,
              })}
            </p>
          </Box>
        </div>

        {actualizadoEn && (
          <Box titulo="Actualización">
            <p className="text-sm font-medium">
              {format(actualizadoEn, "d MMM yyyy '·' HH:mm", { locale: es })}
            </p>
          </Box>
        )}

        <Box titulo="ID">
          <p className="font-mono text-xs break-all">{operacion.operacion}</p>
        </Box>
      </div>
  );
}

function Box({
  titulo,
  children,
}: {
  titulo: string;
  children: React.ReactNode;
}) {
  return (
    <div className="min-w-0 px-4 py-3">
      <p className="text-xs text-muted-foreground">{titulo}</p>
      <div className="mt-0.5">{children}</div>
    </div>
  );
}

function ProveedorSection({
  proveedor,
}: {
  proveedor: Extract<Operacion, { __typename: "GastoAProveedor" }>["proveedor"];
}) {
  return (
    <section>
      <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
        Proveedor
      </h3>
      <div className="mt-3">
        <AvatarIniciales nombre={proveedor.nombre} shape="rounded" />
      </div>
      <dl className="mt-3">
        <Fila label="Nombre" value={proveedor.nombre} strong />
        <Fila label="RIF" value={proveedor.rif} mono />
        <Fila label="Teléfono" value={proveedor.telefono ?? "—"} />
        <Fila label="Correo" value={proveedor.email ?? "—"} />
      </dl>
    </section>
  );
}

function Fila({
  label,
  value,
  strong = false,
  mono = false,
}: {
  label: string;
  value: string;
  strong?: boolean;
  mono?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-gray-100 py-3 last:border-0">
      <dt className="text-sm text-muted-foreground shrink-0">{label}</dt>
      <dd
        className={[
          "text-sm font-medium text-right min-w-0 break-all",
          strong && "font-semibold text-foreground",
          mono && "font-mono text-xs",
        ]
          .filter(Boolean)
          .join(" ")}
      >
        {value}
      </dd>
    </div>
  );
}

function formatearTasa(tasa: number): string {
  return new Intl.NumberFormat("es-VE", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(tasa);
}