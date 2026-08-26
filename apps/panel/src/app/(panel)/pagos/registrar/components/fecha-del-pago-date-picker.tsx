"use client";

import { format } from "date-fns";
import { es } from "date-fns/locale";
import {
  Field,
  FieldLabel,
} from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { withForm } from "@/hooks/useAppForm";
import { registrarPagoDefaultValues } from "./pago-form-schema";

export const FechaDelPagoDatePicker = withForm({
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
