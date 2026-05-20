"use client";

import {
  Field,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import { withForm } from "@/hooks/useAppForm";
import { registrarPagoDefaultValues } from "./pago-form-schema";

export const ReferenciaInput = withForm({
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
