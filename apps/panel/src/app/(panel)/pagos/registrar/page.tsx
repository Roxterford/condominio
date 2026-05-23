"use client";

import { keepPreviousData, useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { VillaSelector } from "./components/villa-selector";
import { ReferenciaInput } from "./components/referencia-input";
import { TasaInput } from "./components/tasa-input";
import { MontoInput } from "./components/monto-input";
import { MetodoDePagoRadioGroup } from "./components/metodo-de-pago-radio-group";
import { FechaDelPagoDatePicker } from "./components/fecha-del-pago-date-picker";
import { MonedaSelection } from "./components/moneda-selection";
import {
  NuevoPagoFormSchema,
  registrarPagoDefaultValues,
} from "./components/pago-form-schema";

import { format } from "date-fns";
import { es } from "date-fns/locale";

import { Button } from "@/components/ui/button";
import { CreditCard } from "lucide-react";
import { useAppForm } from "@/hooks/useAppForm";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";
import {
  RegistrarPagoDto,
  Moneda as MonedaGraph,
  MetodoDePago as MetodoDePagoGraph,
  Mes as MesGraph,
} from "@/providers/graphql/graphql";
import { Spinner } from "@/components/ui/spinner";
import { useStore } from "@tanstack/react-form-nextjs";
import * as v from "valibot";

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
    tasa_hoy: obtenerTasa {
      tipo
      valor
      fecha
    }
  }
`);

const TasaQuery = graphql(/* GraphQL */ `
  query ObtenerTasaOnPagoPage($dia: Int, $mes: Mes, $anio: Int) {
    tasa: obtenerTasa(dia: $dia, mes: $mes, anio: $anio) {
      tipo
      valor
      fecha
    }
  }
`);

const RegistrarPagoMutation = graphql(/* GraphQL */ `
  mutation RegistrarPago($input: RegistrarPagoDTO!) {
    registrarPago(input: $input)
  }
`);

export default function RegistrarPagoPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const form = useAppForm({
    defaultValues: registrarPagoDefaultValues,
    validators: {
      onChange: NuevoPagoFormSchema,
    },
    onSubmit: function sendRegistrarPago({ value, formApi, meta }) {
      const data = v.parse(NuevoPagoFormSchema, value);

      // data.output
      console.log("enviando formulario...");
      let moneda: MonedaGraph, metodo: MetodoDePagoGraph;

      switch (value.moneda) {
        case Moneda.USD:
          moneda = MonedaGraph.Usd;
          break;
        case Moneda.VED:
          moneda = MonedaGraph.Ved;
          break;
        default:
          console.error("no moneda map");
          return;
      }

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
          console.error("no metodo map");
          console.error(value.metodo);
          return;
      }

      console.log("el input:", { value });

      registrarPago.mutate({
        fecha: data.fecha,
        unidad: data.unidad,
        metodo,
        moneda,
        monto: data.monto,
        referencia: data.referencia,
        tasa: data.tasa,
      });
    },
  });

  const form_fecha = useStore(form.store, (s) => s.values.fecha);

  const { data: page, error } = useQuery({
    queryKey: ["villas", searchTerm],
    queryFn: async () => {
      const result = await execute(PageQuery, {
        codigo_like: `%${searchTerm}%`,
      });

      // if (result.errors) throw new Error(JSON.stringify(result.errors));
      return result.data;
    },
    placeholderData: keepPreviousData,
  });

  const obtenerTasa = useQuery({
    enabled: !!form_fecha,
    queryKey: ["tasa", form_fecha?.toISOString()],
    queryFn: async () => {
      const fecha = form.state.values.fecha!;

      const dia = fecha.getDate();
      const anio = fecha.getFullYear();
      const mes = [
        MesGraph.Enero,
        MesGraph.Febrero,
        MesGraph.Marzo,
        MesGraph.Abril,
        MesGraph.Mayo,
        MesGraph.Junio,
        MesGraph.Julio,
        MesGraph.Agosto,
        MesGraph.Septiembre,
        MesGraph.Octubre,
        MesGraph.Noviembre,
        MesGraph.Diciembre,
      ][fecha.getMonth()];

      const result = await execute(TasaQuery, { dia, mes, anio });
      return result.data;
    },
  });

  useEffect(() => {
    if (obtenerTasa.data)
      form.setFieldValue("tasa", obtenerTasa.data?.tasa.valor ?? 0);
  }, [obtenerTasa.data]);

  const registrarPago = useMutation({
    mutationFn: async (input: RegistrarPagoDto) => {
      try {
        const response = await execute(RegistrarPagoMutation, { input });
        console.log({ response });
        return response;
      } catch (e) {
        console.error(e);
      }
    },
  });

  if (error) throw error;

  return (
    <>
      <h1>Registrar Pago</h1>

      <p>{obtenerTasa.isLoading ? "cargando tasa..." : "-"}</p>

      <pre>
        <form.Subscribe selector={(s) => s}>
          {(s) => (
            <>
              {JSON.stringify(s.values, null, 4)}
              <br />
              {s.isValid ? "valido" : "invalido"}
            </>
          )}
        </form.Subscribe>
      </pre>

      {page?.tasa_hoy && (
        <>
          <p>
            Tasa hoy: {page.tasa_hoy.valor} ({page.tasa_hoy.tipo})
          </p>
          <p>Fecha: {page.tasa_hoy.fecha}</p>
          <br />
          <br />
          <p>
            Tasa segun dia: {obtenerTasa.data?.tasa.valor} (
            {obtenerTasa.data?.tasa.tipo})
          </p>
          <p>
            Fecha:{" "}
            {format(toDate(obtenerTasa.data?.tasa.fecha), "PPPP", {
              locale: es,
            })}
          </p>
          <p>
            Busqueda:{" "}
            {format(toDate(form.state.values.fecha), "PPPP", { locale: es })}
          </p>
        </>
      )}

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
              items={page?.unidades?.data ?? []}
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
        <TasaInput
          form={form}
          isLoading={obtenerTasa.isLoading}
          fechaCoincidente={obtenerTasa.data?.tasa.fecha}
        />

        <MonedaSelection form={form} />

        <form.Subscribe
          selector={(s) => [s.isValid, s.isTouched]}
          children={([isValid, isTouched]) => (
            <Button type="submit" disabled={!isValid || !isTouched}>
              {registrarPago.isPending ? (
                <Spinner data-icon="inline-start" />
              ) : (
                <CreditCard />
              )}
              Registrar pago
            </Button>
          )}
        />

        {registrarPago.isSuccess && (
          <>
            {registrarPago.data?.data?.registrarPago ? (
              <div className="flex items-center gap-2 rounded-md border border-green-500 bg-green-50 p-3 text-green-700">
                <span className="font-medium">
                  ✓ Pago registrado exitosamente.
                </span>
              </div>
            ) : (
              <div className="flex flex-col gap-1 rounded-md border border-red-500 bg-red-50 p-3 text-red-700">
                <span className="font-medium">
                  ✗ Error al registrar el pago
                </span>
                {registrarPago.data?.errors && (
                  <span className="text-sm">
                    {registrarPago.data.errors
                      .map((e: any) => e.message)
                      .join(", ")}
                  </span>
                )}
                <pre className="mt-2 overflow-auto rounded bg-red-100 p-2 text-xs">
                  {JSON.stringify(registrarPago.data, null, 2)}
                </pre>
              </div>
            )}
          </>
        )}
        {registrarPago.isError && (
          <div className="flex flex-col gap-1 rounded-md border border-red-500 bg-red-50 p-3 text-red-700">
            <span className="font-medium">✗ Error al registrar el pago</span>
            <span className="text-sm">{registrarPago.error?.message}</span>
            <pre>{JSON.stringify(registrarPago.error, null, 4)}</pre>
          </div>
        )}
      </form>
    </>
  );
}

function toDate(input: any): Date {
  if (!input) return new Date();
  if (input instanceof Date) return input;
  return new Date(input);
}
