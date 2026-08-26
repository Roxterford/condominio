"use client";

import {
  Field,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { withForm } from "@/hooks/useAppForm";
import { registrarPagoDefaultValues } from "./pago-form-schema";

export const MontoInput = withForm({
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
