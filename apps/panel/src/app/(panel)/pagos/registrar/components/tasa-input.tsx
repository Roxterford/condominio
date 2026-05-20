"use client";

import {
  Field,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { withForm } from "@/hooks/useAppForm";
import { registrarPagoDefaultValues } from "./pago-form-schema";
import { InputGroup, InputGroupAddon, InputGroupInput } from "@/components/ui/input-group";
import { Spinner } from "@/components/ui/spinner";
import { useStore } from "@tanstack/react-form-nextjs";

export interface TasaInputProps {
  isLoading?: boolean;
  fechaCoincidente?: string | null;
}

export const TasaInput = withForm({
  props: { isLoading: false, fechaCoincidente: null } as TasaInputProps,
  defaultValues: registrarPagoDefaultValues,
  render: ({ form, isLoading, fechaCoincidente }) => {
    const fechaSeleccionada = useStore(form.store, (s) => s.values.fecha);

    const fechaSeleccionadaStr = fechaSeleccionada
      ? fechaSeleccionada.toISOString().split('T')[0]
      : null;

    const hayDiferenciaFecha =
      fechaCoincidente &&
      fechaSeleccionadaStr &&
      fechaCoincidente !== fechaSeleccionadaStr;

    return (
      <Field>
        <FieldLabel htmlFor="pago_tasa">Tasa (céntimos)</FieldLabel>
        <form.AppField
          name="tasa"
          children={(field) => (
            <InputGroup>
            <InputGroupInput
              id="pago_tasa"
              type="number"
              placeholder="Monto"
              value={field.state.value}
              onChange={(v) => field.handleChange(Number(v.target.value))}
              />
              <InputGroupAddon align="inline-end">
                {isLoading && <Spinner />}
              </InputGroupAddon>
            </InputGroup>
          )}
          />
          <FieldDescription>1,00 VED = 100</FieldDescription>
          {hayDiferenciaFecha && (
          <FieldDescription>
              ⚠️ La tasa mostrada corresponde a la fecha más cercana disponible:{" "}
              <strong>{fechaCoincidente}</strong>, que difiere de la fecha seleccionada ({fechaSeleccionadaStr}). Verifique que la tasa sea válida para la fecha deseada.
          </FieldDescription>
          )}
          <FieldDescription>
            ⚠️ Verifique la tasa antes de continuar. La tasa propuesta se basa en pagos anteriores y servicios de terceros — queda bajo responsabilidad del usuario confirmar su validez.{" "}
            <a
              href="https://www.bcv.org.ve"
              target="_blank"
              rel="noopener noreferrer"
              style={{ textDecoration: "underline" }}
            >
              Más información
            </a>
          </FieldDescription>
      </Field>
    );
  },
});
