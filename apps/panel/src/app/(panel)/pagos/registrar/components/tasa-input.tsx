"use client";

import {
  Field,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { withForm } from "@/hooks/useAppForm";
import { registrarPagoDefaultValues } from "./pago-form-schema";
import { InputGroup, InputGroupInput } from "@/components/ui/input-group";

export const TasaInput = withForm({
  defaultValues: registrarPagoDefaultValues,
  render: ({ form }) => {
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
            </InputGroup>
          )}
        />
        <FieldDescription>1,00 VED = 100</FieldDescription>
        <FieldDescription>
          ⚠️ Verifique la tasa antes de continuar. La tasa se ingresa
          manualmente y queda bajo responsabilidad del usuario confirmar su
          validez.{" "}
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
