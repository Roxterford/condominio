import { Field, FieldLabel } from "@/components/ui/field";
import { InputGroup, InputGroupAddon } from "@/components/ui/input-group";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues } from "../schema";

export const TasaField = withForm({
  defaultValues,
  render: ({ form }) => (
    <div className="flex gap-2">
      <Field orientation="vertical">
        <FieldLabel htmlFor="tasa">Tasa</FieldLabel>
        <InputGroup>
          <InputGroupAddon>
            <span>Bs.</span>
          </InputGroupAddon>
          {/* TODO: mask */}
          <form.AppField
            name="tasa"
            children={(field) => (
              <field.InputGroupInput
                type="number"
                id="tasa"
                placeholder="0,00"
                onChange={(e) => field.handleChange(e.target.valueAsNumber)}
              />
            )}
          />
        </InputGroup>
      </Field>
    </div>
  ),
});
