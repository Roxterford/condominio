"use client";

import { keepPreviousData, useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { VillaSelector } from "./components/villa-selector";

import {
  Field,
  FieldDescription,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";

import { format } from "date-fns";
import { es } from "date-fns/locale";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { Input } from "@/components/ui/input";
import { CreditCard } from "lucide-react";
import { useAppForm, withForm } from "@/hooks/useAppForm";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";
import * as v from "valibot";
import {
  RegistrarPagoDto,
  Moneda as MonedaGraph,
  MetodoDePago as MetodoDePagoGraph,
} from "@/providers/graphql/graphql";
import { Spinner } from "@/components/ui/spinner";

const PageQuery = graphql(/* GraphQL */ `
  query RegistrarPagoPage($codigo_like: String!) {
    unidades: obtenerUnidades(
      filter: { codigo: { like: $codigo_like } }
      paginator: { limit: 5, page: 1 }
    ) {
      data {
        id
        codigo
      }
    }
  }
`);

const TasaQuery = graphql(/* GraphQL */ `
  query ObtenerTasaOnPagoPage($dia: Int, $mes: Mes, $anio: Int) {
    obtenerTasa(dia: $dia, mes: $mes, anio: $anio) {
      tipo
      valor
    }
  }
`);

const RegistrarPagoMutation = graphql(/* GraphQL */ `
  mutation RegistrarPago($input: RegistrarPagoDTO!) {
    registrarPago(input: $input)
  }
`);

const NuevoPagoFormSchema = v.object({
  unidad: v.pipe(v.string(), v.nonEmpty()),
  fecha: v.nullish(v.date()),
  metodo: v.nullish(v.enum(MetodoDePago)),
  referencia: v.string(),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  tasa: v.pipe(v.number(), v.integer(), v.minValue(1)),
  moneda: v.enum(Moneda),
});

type NuevoPagoForm = v.InferOutput<typeof NuevoPagoFormSchema>;

const registrarPagoDefaultValues: NuevoPagoForm = {
  unidad: "",
  fecha: null,
  metodo: null,
  referencia: "",
  monto: 0,
  tasa: 0,
  moneda: Moneda.USD,
};

export default function RegistrarPagoPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const form = useAppForm({
    defaultValues: registrarPagoDefaultValues,
    onSubmit: function sendNewPagoToAPI({ value, formApi, meta }) {
      console.log("enviando form");
      let moneda: MonedaGraph, metodo: MetodoDePagoGraph;

      // Map Moneda
      switch (value.moneda) {
        case Moneda.USD:
          moneda = MonedaGraph.Usd;
          break;
        case Moneda.VED:
          moneda = MonedaGraph.Ved;
          break;
        default:
          return;
      }

      // Map Metodo
      switch (value.metodo) {
        case MetodoDePago.Efectivo:
          metodo = MetodoDePagoGraph.Efectivo;
          break;
        case MetodoDePago.PagoMovil:
          metodo = MetodoDePagoGraph.PagoMovil;
          break;
        case MetodoDePago.Transferencia:
          metodo = MetodoDePagoGraph.Transferencia;
          break;
        default:
          return;
      }

      console.log("preparando");

      registrarPago.mutate({
        unidad: value.unidad,
        metodo,
        moneda,
        monto: value.monto,
        referencia: value.referencia,
        tasa: value.tasa,
      });
    },
  });

  const { data, error } = useQuery({
    queryKey: ["villas", searchTerm],
    queryFn: async () => {
      const result = await execute(PageQuery, {
        codigo_like: `%${searchTerm}%`,
      });

      if (result.errors) throw new Error(JSON.stringify(result.errors));
      return result.data;
    },
    placeholderData: keepPreviousData,
  });

  const registrarPago = useMutation({
    mutationFn: (input: RegistrarPagoDto) =>
      execute(RegistrarPagoMutation, { input }),
  });

  if (error) throw error;

  return (
    <>
      <h1>Registrar Pago</h1>

      <form.Subscribe
        selector={(s) => s.values}
        children={(s) => <pre>{JSON.stringify(s, null, 4)}</pre>}
      />

      <pre>{JSON.stringify(registrarPago.data, null, 4)}</pre>

      <form
        onSubmit={(e) => {
          e.preventDefault();
          form.handleSubmit();
        }}
      >
        <form.Field
          name="unidad"
          children={(f) => (
            <VillaSelector
              items={data?.unidades?.data ?? []}
              onSelect={(v) => f.handleChange(v ?? "")}
              onDebounceChange={(v) => {
                setSearchTerm(v);
              }}
            />
          )}
        />
        <MetodoDePagoRadioGroup form={form} />

        <ReferenciaInput form={form} />

        <FechaDelPagoDatePicker form={form} />

        <MontoInput form={form} />
        <TasaInput form={form} />

        <MonedaSelection form={form} />

        <Button type="submit" disabled={!form.state.isValid}>
          {registrarPago.isPending ? (
            <Spinner data-icon="inline-start" />
          ) : (
            <CreditCard />
          )}
          Registrar pago
        </Button>
      </form>
    </>
  );
}

const ReferenciaInput = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
    return (
      <Field>
        <FieldLabel htmlFor="referencia">Referencia</FieldLabel>
        <form.AppField
          name="referencia"
          children={(field) => (
            <field.Input
              value={field.state.value}
              onChange={(v) => field.handleChange(v.target.value)}
              id="referencia"
              type="text"
              placeholder="Ej. 00002333241 (BDV)"
            />
          )}
        />

        <FieldDescription>Referencia del pago</FieldDescription>
      </Field>
    );
  },
});

const TasaInput = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
    return (
      <Field>
        <FieldLabel htmlFor="pago_tasa">Tasa (centimos)</FieldLabel>
        <form.AppField
          name="tasa"
          children={(field) => (
            <Input
              id="pago_tasa"
              type="number"
              placeholder="Monto"
              value={field.state.value}
              onChange={(v) => field.handleChange(Number(v.target.value))}
            />
          )}
        />
        <FieldDescription>1,00 VED = 100</FieldDescription>
      </Field>
    );
  },
});

const MontoInput = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
    return (
      <Field>
        <FieldLabel htmlFor="pago_monto">Monto (centimos)</FieldLabel>
        <form.AppField
          name="monto"
          children={(field) => (
            <Input
              id="pago_monto"
              type="number"
              placeholder="Monto"
              value={field.state.value}
              onChange={(v) => field.handleChange(Number(v.target.value))}
            />
          )}
        />
        <FieldDescription>1,00 USD/VED = 100</FieldDescription>
      </Field>
    );
  },
});

const MetodoDePagoRadioGroup = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
    return (
      <FieldSet className="w-full max-w-xs">
        <FieldLegend variant="label">Metodo de pago</FieldLegend>
        <form.AppField
          name="metodo"
          children={(field) => (
            <RadioGroup
              value={field.state.value ?? undefined}
              onValueChange={(v) => field.handleChange(v as MetodoDePago)}
            >
              <Field orientation="horizontal">
                <RadioGroupItem value="PAGOMOVIL" id="tp-pagomovil" />
                <FieldLabel htmlFor="tp-pagomovil" className="font-normal">
                  Pagomovil
                </FieldLabel>
              </Field>
              <Field orientation="horizontal">
                <RadioGroupItem value="TRANSFERENCIA" id="tp-transferencia" />
                <FieldLabel htmlFor="tp-transferencia" className="font-normal">
                  Transferencia
                </FieldLabel>
              </Field>
              <Field orientation="horizontal">
                <RadioGroupItem value="EFECTIVO" id="tp-efectivo" />
                <FieldLabel htmlFor="tp-efectivo" className="font-normal">
                  Efectivo
                </FieldLabel>
              </Field>
            </RadioGroup>
          )}
        />
      </FieldSet>
    );
  },
});

const FechaDelPagoDatePicker = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
    return (
      <Field className="">
        <FieldLabel htmlFor="pago_fecha">Fecha del pago</FieldLabel>
        <form.AppField
          name="fecha"
          children={(field) => (
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  id="pago_fecha"
                  className="justify-start font-normal"
                >
                  {field.state.value ? (
                    format(field.state.value, "PPP", { locale: es })
                  ) : (
                    <span>Seleccione</span>
                  )}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                  mode="single"
                  selected={field.state.value ?? undefined}
                  onSelect={(v) => field.handleChange(v ?? null)}
                  defaultMonth={field.state.value ?? undefined}
                />
              </PopoverContent>
            </Popover>
          )}
        />
      </Field>
    );
  },
});

const MonedaSelection = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
    return (
      <Field>
        <FieldLabel>Moneda</FieldLabel>
        <form.AppField
          name="moneda"
          children={(field) => (
            <Select
              value={field.state.value}
              onValueChange={(v) => field.handleChange(v as Moneda)}
            >
              <SelectTrigger className="w-full max-w-48">
                <SelectValue placeholder="Seleccione moneda" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Moneda</SelectLabel>
                  <SelectItem value="USD">
                    (USD) Dólar estadounidense
                  </SelectItem>
                  <SelectItem value="VED">(VED) Bolívar Digital</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          )}
        />
      </Field>
    );
  },
});
