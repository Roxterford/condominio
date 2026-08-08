import { Field, FieldContent, FieldDescription, FieldLabel } from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues, type EstrategiaDeDistribucion } from "../schema";

export const EstrategiaDeDistribucionField = withForm({
  defaultValues,
  render: ({ form }) => (
    <section>
      <h3>Estrategia de distribución</h3>
      <form.AppField
        name="estrategia"
        children={(field) => (
          <RadioGroup
            value={field.state.value}
            onValueChange={(value) =>
              field.handleChange(value as EstrategiaDeDistribucion)
            }
            className="w-fit"
          >
            <Field orientation="horizontal">
              <RadioGroupItem value="lineal" id="lineal" />
              <FieldContent>
                <FieldLabel htmlFor="lineal">Lineal</FieldLabel>
                <FieldDescription>
                  Los gastos se distribuyen de forma equitativa segun el numero
                  de villas activas
                </FieldDescription>
              </FieldContent>
            </Field>
            <Field orientation="horizontal" data-disabled>
              <RadioGroupItem disabled value="individual" id="individual" />
              <FieldContent>
                <FieldLabel htmlFor="individual">Individual</FieldLabel>
                <FieldDescription>
                  Los gastos se atribuyen a una villa espesifica
                </FieldDescription>
              </FieldContent>
            </Field>
            <Field orientation="horizontal" data-disabled>
              <RadioGroupItem disabled value="alicuota" id="alicuota" />
              <FieldContent>
                <FieldLabel htmlFor="alicuota">Alicuota</FieldLabel>
                <FieldDescription>
                  Los gastos se distribuyen entre las villas activas segun su
                  porcentaje de copropiedad
                </FieldDescription>
              </FieldContent>
            </Field>
          </RadioGroup>
        )}
      />
    </section>
  ),
});