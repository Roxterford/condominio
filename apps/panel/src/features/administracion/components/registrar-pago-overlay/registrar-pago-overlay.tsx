"use client";

import { OverlayProps } from "@/components/overlay";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Field,
  FieldDescription,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemTitle,
} from "@/components/ui/item";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { useAppForm } from "@/hooks/useAppForm";
import { useDebounce } from "@/hooks/useDebounce";
import { Check } from "lucide-react";
import { useEffect, useState, type SubmitEventHandler } from "react";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";

import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/providers/graphql/graphql";

import {
  DestinoDePago,
  registrarPagoDefaultValues,
  RegistrarPagoFormSchema,
} from "./schema";

const UnidadesQuery = graphql(`
  query RegistrarPagoOverlayUnidades($codigo_like: String!) {
    unidades: obtenerUnidades(
      filter: { codigo: { like: $codigo_like } }
      paginator: { limit: 10, page: 1 }
    ) {
      data {
        id
        codigo
      }
    }
  }
`);

const DEUDAS_MOCK = [
  { id: "cuota-ago-2026", label: "Cuota Ago 2026 · $120" },
  { id: "cuota-jul-2026", label: "Cuota Jul 2026 · $120" },
  { id: "cuota-jun-2026", label: "Cuota Jun 2026 · $120" },
  { id: "proyecto-piscina", label: "Proyecto Piscina · $200" },
];

const DESTINOS = [
  {
    value: DestinoDePago.DeudaEspecifica,
    titulo: "Deuda específica",
    descripcion: "Tú eliges la cuota o proyecto a abonar",
  },
  {
    value: DestinoDePago.MasAntigua,
    titulo: "Más antigua primero (FIFO)",
    descripcion:
      "El sistema destina a la deuda de mayor antigüedad / vencimiento",
  },
  {
    value: DestinoDePago.Recargo,
    titulo: "Recargo por pago tardío primero",
    descripcion: "Aplica al interés/mora y luego al capital",
  },
  {
    value: DestinoDePago.MayorSaldo,
    titulo: "Mayor saldo primero",
    descripcion: "A la deuda pendiente más alta",
  },
  {
    value: DestinoDePago.Prorrateo,
    titulo: "Prorrateo",
    descripcion: "Reparte el monto entre todas las deudas según su saldo",
  },
  {
    value: DestinoDePago.AbonoCuenta,
    titulo: "Abono a cuenta (saldo a favor)",
    descripcion: "No se liga a una deuda; queda disponible para pagos futuros",
  },
];

export interface RegistrarPagoOverlayProps extends OverlayProps {}

export function RegistrarPagoOverlay(props: RegistrarPagoOverlayProps) {
  const form = useAppForm({
    defaultValues: registrarPagoDefaultValues,
    validators: {
      onChange: RegistrarPagoFormSchema,
      onBlur: RegistrarPagoFormSchema,
    },
    onSubmit: ({ value }) => {
      toast.success("Pago registrado (prototipo)", {
        description: `Unidad ${value.unidad} · ${value.monto} ${value.moneda}`,
      });
      console.log("RegistrarPagoOverlay payload:", value);
    },
  });

  const [vistaDestino, setVistaDestino] = useState<"lista" | "detalle">(
    "lista",
  );

  const [busquedaUnidad, setBusquedaUnidad] = useState("");

  const unidades = useQuery({
    queryKey: ["unidades-pago-overlay", busquedaUnidad],
    enabled: busquedaUnidad.trim().length > 0,
    queryFn: () =>
      execute(UnidadesQuery, { codigo_like: `%${busquedaUnidad}%` }),
  });

  useEffect(() => {
    if (props.open) setVistaDestino("lista");
  }, [props.open]);

  const handleSubmit: SubmitEventHandler = (event) => {
    event.preventDefault();
    event.stopPropagation();
    form.handleSubmit();
  };

  return (
    <Dialog {...props} modal={false}>
      <DialogContent className="md:min-w-lg max-h-[90vh] max-w-2xl">
        <DialogHeader>
          <DialogTitle>Registrar pago</DialogTitle>
          <DialogDescription>
            Registra un abono de una unidad y su destino de aplicación.
          </DialogDescription>
        </DialogHeader>

        <form
          className="grid gap-4 overflow-y-auto max-h-[50vh] -mx-6 px-6"
          onSubmit={handleSubmit}
        >
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field className="sm:col-span-2">
              <FieldLabel>Unidad</FieldLabel>
              <form.AppField
                name="unidad"
                children={(field) => (
                  <UnidadCombobox
                    items={unidades.data?.data?.unidades?.data ?? []}
                    onBusquedaChange={setBusquedaUnidad}
                    onSeleccionar={(id) => {
                      field.handleChange(id);
                    }}
                  />
                )}
              />
            </Field>

            <Field>
              <FieldLabel>Monto (USD)</FieldLabel>
              <form.AppField
                name="monto"
                children={(field) => (
                  <Input
                    type="number"
                    placeholder="120"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(Number(e.target.value))}
                  />
                )}
              />
            </Field>

            <Field>
              <FieldLabel>Método</FieldLabel>
              <form.AppField
                name="metodo"
                children={(field) => (
                  <Select
                    value={field.state.value ?? undefined}
                    onValueChange={(v) => field.handleChange(v as MetodoDePago)}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Método" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="PAGOMOVIL">Pago móvil</SelectItem>
                      <SelectItem value="TRANSFERENCIA">
                        Transferencia
                      </SelectItem>
                      <SelectItem value="EFECTIVO">Efectivo</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </Field>

            <Field>
              <FieldLabel>Referencia</FieldLabel>
              <form.AppField
                name="referencia"
                children={(field) => (
                  <Input
                    placeholder="Ej. 000123456789"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                )}
              />
            </Field>

            <Field>
              <FieldLabel>Moneda</FieldLabel>
              <form.AppField
                name="moneda"
                children={(field) => (
                  <Select
                    value={field.state.value}
                    onValueChange={(v) => field.handleChange(v as Moneda)}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="USD">USD</SelectItem>
                      <SelectItem value="VED">VED</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </Field>

            <form.Subscribe
              selector={(state) => state.values.moneda}
              children={(moneda) => (
                <Field data-disabled={moneda === "VED"}>
                  <FieldLabel>Tasa (VED/USD)</FieldLabel>
                  <form.AppField
                    name="tasa"
                    children={(field) => (
                      <Input
                        type="number"
                        disabled={moneda === "VED"}
                        value={field.state.value}
                        onChange={(e) =>
                          field.handleChange(Number(e.target.value))
                        }
                      />
                    )}
                  />
                  <FieldDescription>
                    {moneda === "VED"
                      ? "No aplica conversión al pagar en VED."
                      : "Tasa manual; verifique antes de continuar."}
                  </FieldDescription>
                </Field>
              )}
            />

            <Field className="sm:col-span-2">
              <FieldLabel>Concepto</FieldLabel>
              <form.AppField
                name="concepto"
                children={(field) => (
                  <Input
                    placeholder="Pago cuota Agosto"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                )}
              />
            </Field>
          </div>

          <Separator />

          <FieldSet className="gap-3">
            <FieldLegend variant="label">Aplicar a deuda</FieldLegend>
            <form.Subscribe selector={(state) => state.values.destino}>
              {(destinoSeleccionado) => (
                <div className="overflow-hidden rounded-md">
                  <div
                    className="flex transition-transform duration-200 ease-out"
                    style={{
                      transform:
                        vistaDestino === "detalle"
                          ? "translateX(-100%)"
                          : "translateX(0)",
                    }}
                  >
                    <div className="w-full shrink-0 space-y-2">
                      {DESTINOS.map((opcion) => (
                        <button
                          type="button"
                          key={opcion.value}
                          onClick={() => {
                            form.setFieldValue(
                              "destino",
                              opcion.value as DestinoDePago,
                            );
                            setVistaDestino("detalle");
                          }}
                          className="flex w-full cursor-pointer items-start gap-2 rounded-md border p-3 text-left transition-colors hover:bg-muted/50 data-[active=true]:border-primary"
                          data-active={destinoSeleccionado === opcion.value}
                        >
                          <span className="mt-1.5 size-2 shrink-0 rounded-full bg-primary/60" />
                          <div>
                            <div className="text-sm font-medium">
                              {opcion.titulo}
                            </div>
                            <div className="text-xs text-muted-foreground">
                              {opcion.descripcion}
                            </div>
                          </div>
                        </button>
                      ))}
                    </div>

                    <div className="w-full shrink-0 space-y-2">
                      <button
                        type="button"
                        onClick={() => setVistaDestino("lista")}
                        className="mb-1 text-xs font-medium text-primary hover:underline"
                      >
                        ← Atrás
                      </button>
                      <DestinoDetalle form={form} />
                    </div>
                  </div>
                </div>
              )}
            </form.Subscribe>
          </FieldSet>

          <form.Subscribe
            selector={(state) => ({
              monto: state.values.monto,
              tasa: state.values.tasa,
              moneda: state.values.moneda,
            })}
            children={({ monto, tasa, moneda }) =>
              moneda === "USD" && monto > 0 && tasa > 0 ? (
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Badge variant="secondary">Vista previa</Badge>≈{" "}
                  {(monto * tasa).toLocaleString("es-VE")} VED
                </div>
              ) : null
            }
          />

          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => props.onOpenChange?.(false)}
            >
              Cancelar
            </Button>
            <form.Subscribe
              selector={(state) => state.isValid}
              children={(isValid) => (
                <Button type="submit" disabled={!isValid}>
                  Registrar pago
                  <Check />
                </Button>
              )}
            />
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function DestinoDetalle(props: { form: any }) {
  return (
    <props.form.Subscribe
      selector={(state: any) => state.values.destino}
      children={(destino: any) => {
        if (destino === DestinoDePago.DeudaEspecifica) {
          return (
            <div className="space-y-2">
              <div className="text-sm font-medium">Selecciona la deuda</div>
              <props.form.AppField
                name="deuda_especifica"
                children={(field: any) => (
                  <Select
                    value={field.state.value || undefined}
                    onValueChange={field.handleChange}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Deuda a abonar" />
                    </SelectTrigger>
                    <SelectContent>
                      {DEUDAS_MOCK.map((deuda) => (
                        <SelectItem key={deuda.id} value={deuda.id}>
                          {deuda.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          );
        }

        if (destino === DestinoDePago.AbonoCuenta) {
          return (
            <div className="text-sm text-muted-foreground">
              El pago queda registrado como <b>saldo a favor</b> de la unidad;
              no requiere destino.
            </div>
          );
        }

        const esProrrateo = destino === DestinoDePago.Prorrateo;
        return (
          <div className="space-y-2">
            <div className="text-sm font-medium">Deuda más antigua</div>
            <div className="rounded-md border p-3">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm font-medium">Cuota Jun 2026</div>
                  <div className="text-xs text-muted-foreground">
                    Vencida hace 57 días
                    <span className="text-destructive">· +$3,60 (3%)</span>
                  </div>
                </div>
                <div className="font-semibold">$120</div>
              </div>
            </div>
            {esProrrateo && (
              <>
                <div className="rounded-md border p-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="text-sm font-medium">Cuota Jul 2026</div>
                      <div className="text-xs text-muted-foreground">
                        Vencida hace 27 días
                      </div>
                    </div>
                    <div className="font-semibold">≈ $40</div>
                  </div>
                </div>
                <div className="rounded-md border p-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="text-sm font-medium">Cuota Ago 2026</div>
                      <div className="text-xs text-muted-foreground">
                        Vence en 5 días
                      </div>
                    </div>
                    <div className="font-semibold">≈ $40</div>
                  </div>
                </div>
              </>
            )}
            <div className="text-xs text-muted-foreground">
              {destino === DestinoDePago.Recargo
                ? "Se aplica primero al recargo por pago tardío y luego al capital."
                : destino === DestinoDePago.MayorSaldo
                  ? "El pago se aplicará a la deuda de mayor saldo."
                  : destino === DestinoDePago.Prorrateo
                    ? "Reparto proporcional según el saldo de cada deuda."
                    : "El pago se aplicará a la deuda de mayor antigüedad."}
            </div>
          </div>
        );
      }}
    />
  );
}

function UnidadCombobox(props: {
  items: Array<Record<"id" | "codigo", string>>;
  onBusquedaChange: (value: string) => void;
  onSeleccionar: (id: string) => void;
}) {
  const debouncedChange = useDebounce((value: string) => {
    props.onBusquedaChange(value);
  }, 300);

  return (
    <Combobox
      items={props.items}
      itemToStringLabel={(villa: Record<"id" | "codigo", string>) =>
        villa.codigo
      }
      onValueChange={(v) =>
        props.onSeleccionar(
          (v as Record<"id" | "codigo", string> | null)?.id ?? "",
        )
      }
    >
      <ComboboxInput
        placeholder="Buscar unidad (código)…"
        onChange={(e) => debouncedChange(e.target.value)}
      />
      <ComboboxContent>
        <ComboboxEmpty>Sin resultados.</ComboboxEmpty>
        <ComboboxList>
          {(villa: Record<"id" | "codigo", string>) => (
            <ComboboxItem
              key={villa.codigo}
              value={villa as Record<"id" | "codigo", string>}
            >
              <Item size="xs" className="p-0">
                <ItemContent>
                  <ItemTitle className="whitespace-nowrap">
                    Villa {villa.codigo}
                  </ItemTitle>
                  <ItemDescription>{villa.id}</ItemDescription>
                </ItemContent>
              </Item>
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
