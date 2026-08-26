import { Field, FieldLabel } from "@/components/ui/field";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues } from "../schema";

export const FechaField = withForm({
  defaultValues,
  render: ({ form }) => (
    <Field className="">
      <FieldLabel htmlFor="fecha_emision">Fecha</FieldLabel>
      <form.AppField
        name="fecha"
        children={(field) => (
          <field.DatePickerInput
            placeholder="dd/mm/aaaa"
            locale="es-VE"
            id="fecha_emision"
            value={field.state.value}
            onChange={(e) => field.handleChange(e as Date)}
          />
        )}
      />
    </Field>
  ),
});
