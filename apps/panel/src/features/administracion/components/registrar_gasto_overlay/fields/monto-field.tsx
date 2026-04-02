import { Field, FieldLabel } from "@/components/ui/field";
import { InputGroup, InputGroupAddon } from "@/components/ui/input-group";
import {
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";
import { withForm } from "@/hooks/useAppForm";
import { DollarSign } from "lucide-react";
import { defaultValues } from "../schema";

export const MontoField = withForm({
  defaultValues,
  render: ({ form }) => (
    <div className="flex gap-2">
      <Field orientation="vertical">
        <FieldLabel htmlFor="moneda">Moneda</FieldLabel>
        <form.AppField
          name="moneda"
          children={(field) => (
            <field.Select
              value={field.state.value}
              onValueChange={(value) => {
                field.handleChange(value as Moneda);
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Seleccione" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={Moneda.USD}>Dólar (USD)</SelectItem>
                <SelectItem value={Moneda.VED}>Bolívar (VED)</SelectItem>
              </SelectContent>
            </field.Select>
          )}
        />
      </Field>
      <Field orientation="vertical">
        <FieldLabel htmlFor="monto">Monto</FieldLabel>
        <InputGroup>
          <InputGroupAddon>
            <form.Subscribe
              selector={(state: { values: { moneda: string } }) =>
                state.values.moneda
              }
            >
              {(moneda: string) =>
                moneda === Moneda.USD ? <DollarSign /> : <span>Bs.</span>
              }
            </form.Subscribe>
          </InputGroupAddon>
          {/* TODO: mask */}
          <form.AppField
            name="monto"
            children={(field) => (
              <field.InputGroupInput
                type="number"
                id="monto"
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
