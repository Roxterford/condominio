import { Field, FieldLabel } from "@/components/ui/field";
import { withForm } from "@/hooks/useAppForm";
import { defaultValues } from "../schema";

interface NuevoProveedorFieldsProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form: any;
}

export const NuevoProveedorFields = withForm({
  defaultValues,
  render: ({ form }) => (
    <>
      <Field orientation="vertical">
        <FieldLabel htmlFor="proveedor_nombre">Nombre</FieldLabel>
        <form.AppField
          name="proveedor.nombre"
          children={(field) => (
            <field.Input
              type="text"
              id="proveedor_nombre"
              placeholder="Nombre del provedor"
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>
      <Field orientation="vertical">
        <FieldLabel htmlFor="proveedor_rif">CI / RIF</FieldLabel>
        <form.AppField
          name="proveedor.rif"
          children={(field) => (
            <field.Input
              type="text"
              id="proveedor_rif"
              placeholder="CI / RIF del provedor"
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>
      <div className="flex gap-2">
        <Field orientation="vertical">
          <FieldLabel htmlFor="proveedor_telefono">Telefono</FieldLabel>
          {/* TODO: Mask */}
          <form.AppField
            name="proveedor.telefono"
            children={(field) => (
              <field.Input
                type="text"
                id="proveedor_telefono"
                placeholder="+584123456789"
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          />
        </Field>
        <Field orientation="vertical">
          <FieldLabel htmlFor="proveedor_email">Correo Electrónico</FieldLabel>
          {/* TODO: Mask */}
          <form.AppField
            name="proveedor.email"
            children={(field) => (
              <field.Input
                type="email"
                id="proveedor_email"
                placeholder="proveedor@ejemplo.com"
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          />
        </Field>
      </div>
    </>
  ),
});
