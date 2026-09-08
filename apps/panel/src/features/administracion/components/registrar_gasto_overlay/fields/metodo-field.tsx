import { Field, FieldLabel } from "@/components/ui/field";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues } from "../schema";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { MetodoDeOperacion } from "@/providers/graphql/graphql";

const METODOS: Array<{ value: MetodoDeOperacion | null; label: string }> = [
  { value: null, label: "Seleccione un metodo" },
  { value: MetodoDeOperacion.PagoMovil, label: "Pago móvil" },
  { value: MetodoDeOperacion.TransferenciaNacional, label: "Transferencia" },
  {
    value: MetodoDeOperacion.TransferenciaInternacional,
    label: "Transferencia internacional",
  },
  { value: MetodoDeOperacion.Efectivo, label: "Efectivo" },
  { value: MetodoDeOperacion.Compensacion, label: "Compensación" },
];

export const MetodoField = withForm({
  defaultValues,
  render: ({ form }) => (
    <div className="flex gap-2">
      <Field orientation="vertical">
        <FieldLabel htmlFor="tasa">Metodo</FieldLabel>
        <form.AppField name="metodo">
          {(field) => (
            <Select<MetodoDeOperacion>
              items={METODOS}
              onValueChange={(it) => {
                field.handleChange(it!);
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Seleccione" data-diabled />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Metodos</SelectLabel>
                  {METODOS.map((it) => {
                    return (
                      <SelectItem
                        key={it.value}
                        value={it.value}
                        disabled={it.value === null}
                      >
                        {it.label}
                      </SelectItem>
                    );
                  })}
                </SelectGroup>
              </SelectContent>
            </Select>
          )}
        </form.AppField>
      </Field>
    </div>
  ),
});
