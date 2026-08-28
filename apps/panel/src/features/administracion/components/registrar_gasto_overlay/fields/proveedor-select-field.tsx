import { Field, FieldLabel } from "@/components/ui/field";
import {
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Proveedor } from "@/providers/graphql/graphql";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues } from "../schema";

interface ProveedorSelectFieldProps {
  proveedores: Pick<Proveedor, "id" | "nombre">[];
}

export const ProveedorSelectField = withForm({
  defaultValues,
  props: {
    proveedores: [],
  } as ProveedorSelectFieldProps,
  render: ({ form, proveedores }) => (
    <Field orientation="vertical">
      <FieldLabel>Proveedor</FieldLabel>
      <form.AppField
        name="proveedor"
        children={(field) => (
          <field.Select
            value={field.state.value.toString()}
            onValueChange={field.handleChange}
          >
            <SelectTrigger>
              <SelectValue placeholder="Seleccione" />
            </SelectTrigger>
            <SelectContent>
              {proveedores.map((proveedor) => (
                <SelectItem key={proveedor.id} value={proveedor.id}>
                  {proveedor.nombre}
                </SelectItem>
              ))}
            </SelectContent>
          </field.Select>
        )}
      />
    </Field>
  ),
});
