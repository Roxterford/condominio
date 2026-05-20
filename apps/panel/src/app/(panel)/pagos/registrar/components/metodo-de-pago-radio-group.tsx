"use client";

import {
  Field,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { withForm } from "@/hooks/useAppForm";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { registrarPagoDefaultValues } from "./pago-form-schema";

export const MetodoDePagoRadioGroup = withForm({
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
