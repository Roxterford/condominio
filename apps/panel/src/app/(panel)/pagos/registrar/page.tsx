"use client";

import { keepPreviousData, useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
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

import { Button } from "@/components/ui/button";
import { CreditCard } from "lucide-react";
import { useAppForm } from "@/hooks/useAppForm";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/providers/graphql/graphql";
import {
  Moneda as MonedaGraph,
} from "@/providers/graphql/graphql";
import { Spinner } from "@/components/ui/spinner";
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

      switch (value.moneda) {
        default:
          console.error("no moneda map");
          return;
      }

      switch (value.metodo) {
        default:
          console.error("no metodo map");
          console.error(value.metodo);
          return;
      }

      console.log("el input:", { value });

    },
  });

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


  if (error) throw error;

  return (
    <>
      <h1>Registrar Pago</h1>

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
        <TasaInput form={form} />

        <MonedaSelection form={form} />

        <form.Subscribe
          selector={(s) => [s.isValid, s.isTouched]}
          children={([isValid, isTouched]) => (
            <Button type="submit" disabled={!isValid || !isTouched}>
              Registrar pago
            </Button>
          )}
        />

      </form>
    </>
  );
}
