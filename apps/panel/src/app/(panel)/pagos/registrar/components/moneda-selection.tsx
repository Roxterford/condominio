"use client";

import {
  Field,
  FieldLabel,
} from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { withForm } from "@/hooks/useAppForm";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";
import { registrarPagoDefaultValues } from "./pago-form-schema";

export const MonedaSelection = withForm({
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
