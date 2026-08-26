import { Field, FieldLabel } from "@/components/ui/field";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues } from "../schema";

export const ConceptoField = withForm({
  defaultValues,
  render: ({ form }) => {
    return (
      <Field>
        <FieldLabel htmlFor="concepto">Concepto</FieldLabel>
        <form.AppField
          name="concepto"
          children={(field) => (
            <field.Input
              id="concepto"
              placeholder="Ej. Limpieza"
              onChange={(e) => {
                field.handleChange(e.target.value);
              }}
            />
          )}
        />
      </Field>
    );
  },
});
