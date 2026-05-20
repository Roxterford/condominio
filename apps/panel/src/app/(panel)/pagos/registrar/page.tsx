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
import { registrarPagoDefaultValues } from "./components/pago-form-schema";

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
    onSubmit: function sendNewPagoToAPI({ value, formApi, meta }) {
      let moneda: MonedaGraph, metodo: MetodoDePagoGraph;

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

  const form_fecha = useStore(form.store, (s)=> s.values.fecha)

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
      form.setFieldValue('tasa', obtenerTasa.data?.tasa.valor ?? 0)
  }, [obtenerTasa.data])

  const registrarPago = useMutation({
    mutationFn: (input: RegistrarPagoDto) =>
      execute(RegistrarPagoMutation, { input }),
  });

  if (error) throw error;

  return (
    <>
      <h1>Registrar Pago</h1>

      <p>
        {obtenerTasa.isLoading ? "cargando tasa..." : "-"}
      </p>

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
            {format(
              toDate(obtenerTasa.data?.tasa.fecha),
              "PPPP",
              { locale: es },
            )}
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
        <TasaInput form={form} isLoading={obtenerTasa.isLoading} fechaCoincidente={obtenerTasa.data?.tasa.fecha}/>

        <MonedaSelection form={form} />


        <form.Subscribe
        selector={(s) => s.isValid}
        children={
 <Button type="submit" disabled={!form.state.isValid}>
          {registrarPago.isPending ? (
            <Spinner data-icon="inline-start" />
          ) : (
            <CreditCard />
          )}
          Registrar pago
        </Button>
        }
        />

       
      </form>
    </>
  );
}

function toDate(input: any): Date {
  if (!input) return new Date();
  if (input instanceof Date) return input;
  return new Date(input);
}
