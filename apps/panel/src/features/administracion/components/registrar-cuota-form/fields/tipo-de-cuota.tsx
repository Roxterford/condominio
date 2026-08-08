import { Field, FieldContent, FieldDescription, FieldLabel } from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { withForm } from "@/hooks/useAppForm";
import { TipoDeCuota } from "@/providers/graphql/graphql";
import { defaultValues } from "../schema";

export const TipoDeCuotaField = withForm({
  defaultValues,
  render: ({ form }) => (
    <section>
      <h3>Tipo de cuota</h3>
      <form.AppField
        name="tipo"
        children={(field) => (
          <RadioGroup
            value={field.state.value}
            onValueChange={(value) => field.handleChange(value as TipoDeCuota)}
            className="w-fit"
          >
            <Field orientation="horizontal">
              <RadioGroupItem value={TipoDeCuota.Regular} id="regular" />
              <FieldContent>
                <FieldLabel htmlFor="regular">Mensualidad Regular</FieldLabel>
                <FieldDescription>Cuota mensual estándar</FieldDescription>
              </FieldContent>
            </Field>
            <Field orientation="horizontal">
              <RadioGroupItem value={TipoDeCuota.Especial} id="especial" />
              <FieldContent>
                <FieldLabel htmlFor="especial">Cuota Especial</FieldLabel>
                <FieldDescription>
                  Cuota destinada a algún tipo de inprevisto
                </FieldDescription>
              </FieldContent>
            </Field>
          </RadioGroup>
        )}
      />
    </section>
  ),
});