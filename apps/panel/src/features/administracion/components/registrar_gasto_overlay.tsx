import { OverlayProps } from "@/components/overlay";
import { Button } from "@/components/ui/button";
import { DatePickerInput } from "@/components/ui/date-picker-input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldDescription, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { InputGroup, InputGroupAddon } from "@/components/ui/input-group";
import {
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAppForm, withForm } from "@/hooks/useAppForm";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { Moneda as MonedaGraphql } from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";
import { Check, DollarSign } from "lucide-react";
import { useEffect, type SubmitEventHandler } from "react";
import * as v from "valibot";
import { Proveedor } from "../schemas";
import {
  NuevoGastoSchema,
  NuevoGastoYProveedorSchema,
  type NuevoGasto,
} from "../schemas/gasto.schema";
import { Moneda } from "../schemas/moneda.schema";

const NuevoGastoFormSchema = v.variant("provedor_registrado", [
  v.object({
    provedor_registrado: v.literal(true),
    ...NuevoGastoSchema.entries,
  }),
  v.object({
    provedor_registrado: v.literal(false),
    ...NuevoGastoYProveedorSchema.entries,
  }),
]);

type NuevoGastoForm = v.InferOutput<typeof NuevoGastoFormSchema>;

const defaultValues: NuevoGastoForm = {
  provedor_registrado: true,
  concepto: "",
  proveedor: "",
  monto: 0,
  moneda: Moneda.USD,
  fecha: new Date(),
};

export interface RegistrarGastoOverlayProps extends OverlayProps {
  proveedores: Pick<Proveedor, "id" | "nombre">[];
}

const RegistrarGastoMutation = graphql(`
  mutation RegistrarGasto($input: RegistrarGastoDTO!) {
    registrarGasto(input: $input) {
      id
    }
  }
`);

export function RegistrarGastoOverlay(props: RegistrarGastoOverlayProps) {
  useEffect(() => {
    if (!props.proveedores.length) {
      form.setFieldValue("provedor_registrado", false);
      form.setFieldValue("proveedor", {
        nombre: "",
        rif: "",
        telefono: "",
        email: "",
      });
    }
  }, [props.proveedores.length]);

  const form = useAppForm({
    defaultValues,
    validators: {
      onBlur: NuevoGastoFormSchema,
    },
  });
  const registrarGasto = useMutation({
    mutationFn: (data: NuevoGasto) => {
      const moneda =
        data.moneda === Moneda.USD ? MonedaGraphql.Usd : MonedaGraphql.Ved;

      return execute(RegistrarGastoMutation, {
        input: {
          concepo: data.concepto,
          proveedor: data.proveedor,
          monto: data.monto,
          fecha: data.fecha,
          moneda: moneda,
        },
      });
    },
  });

  const handleSubmit: SubmitEventHandler = (event) => {
    event.preventDefault();
    // registrarGasto.mutateAsync();
  };

  return (
    <Dialog {...props}>
      <DialogContent className="">
        <DialogHeader>
          <DialogTitle>Registrar Gasto</DialogTitle>
        </DialogHeader>

        <form className="grid gap-5" onSubmit={handleSubmit}>
          <section className="grid gap-5">
            <Field>
              <FieldLabel htmlFor="concepto">Consepto</FieldLabel>
              <form.AppField
                name="concepto"
                children={(field) => (
                  <field.Input
                    id="concepto"
                    placeholder="Ej. Limpieza"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                  />
                )}
              />
            </Field>
            {!!props.proveedores.length && (
              <Field orientation="horizontal">
                <form.AppField
                  name="provedor_registrado"
                  children={(field) => (
                    <field.Checkbox
                      id="proveedor-no-registrado"
                      name="provedor_registrado"
                      checked={!field.state.value}
                      onCheckedChange={() => {
                        field.handleChange(!field.state.value);
                        // Resetear el campo de proveedor cuando se marca como no registrado
                        if (!field.state.value) {
                          form.setFieldValue("proveedor", "");
                        } else {
                          // Limpiar el campo de proveedor cuando se marca como registrado
                          form.setFieldValue("proveedor", {
                            nombre: "",
                            rif: "",
                            telefono: "",
                            email: "",
                          });
                        }
                      }}
                    />
                  )}
                />
                <FieldLabel htmlFor="proveedor-no-registrado">
                  Proveedor no registrado
                </FieldLabel>
              </Field>
            )}

            <form.Subscribe
              selector={(state) => state.values.provedor_registrado}
              children={(provedor_registrado) =>
                provedor_registrado ? (
                  <Field orientation="vertical">
                    <FieldLabel>Provedor</FieldLabel>
                    <form.AppField
                      name="proveedor"
                      children={(field) => (
                        <field.Select
                          value={field.state.value.toString()}
                          onValueChange={(value) => {
                            field.handleChange(value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccione" />
                          </SelectTrigger>
                          <SelectContent>
                            {props.proveedores.map((proveedor) => (
                              <SelectItem
                                key={proveedor.id}
                                value={proveedor.id}
                              >
                                {proveedor.nombre}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </field.Select>
                      )}
                    />
                  </Field>
                ) : (
                  <NuevoProveedorFields form={form} />
                )
              }
            />
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
                        <SelectItem value={Moneda.VED}>
                          Bolívar (VED)
                        </SelectItem>
                      </SelectContent>
                    </field.Select>
                  )}
                />
              </Field>
              <Field orientation="vertical">
                <FieldLabel htmlFor="monto">Monto</FieldLabel>
                <InputGroup>
                  <InputGroupAddon>
                    <form.Subscribe selector={(state) => state.values.moneda}>
                      {(moneda) =>
                        moneda === Moneda.USD ? (
                          <DollarSign />
                        ) : (
                          <span>Bs.</span>
                        )
                      }
                    </form.Subscribe>
                  </InputGroupAddon>
                  {/* TODO: Agregar mascara de moneda */}
                  <form.AppField
                    name="monto"
                    children={(field) => (
                      <field.InputGroupInput
                        type="number"
                        id="monto"
                        value={field.state.value}
                        onChange={(e) =>
                          field.handleChange(Number(e.target.value) ?? "")
                        }
                        placeholder="0,00"
                      />
                    )}
                  />
                </InputGroup>
              </Field>
            </div>
            <Field className="">
              <FieldLabel htmlFor="fecha_emision">Fecha</FieldLabel>
              <form.AppField
                name="fecha"
                children={(field) => (
                  <DatePickerInput
                    placeholder="dd/mm/aaaa"
                    locale="es-VE"
                    id="fecha_emision"
                    value={field.state.value}
                    onChange={(date) => field.handleChange(date || new Date())}
                  />
                )}
              />
              <FieldDescription>
                {/* TODO: Agregar fecha relativa a la seleccionada */}
              </FieldDescription>
            </Field>
            {/* TODO: Comprobante */}
            {/* <Field className="">
              <FieldLabel htmlFor="comprobante">Comprobante</FieldLabel>
              <Input type="file" />
              <FieldDescription>
                Archivos permitidos: PDF, JPG, PNG
              </FieldDescription>
            </Field> */}
          </section>

          <div className="flex gap-2 justify-end">
            <Button type="button" variant="outline">
              Cancelar
            </Button>
            <form.Subscribe
              selector={(state) => state.isValid}
              children={(isValid) => (
                <Button
                  type="submit"
                  disabled={form.state.isSubmitting || !isValid}
                >
                  Registrar
                  <Check />
                </Button>
              )}
            />
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

const NuevoProveedorFields = withForm({
  defaultValues: defaultValues as NuevoGastoForm,
  render: ({ form }) => (
    <>
      <Field orientation="vertical">
        <FieldLabel htmlFor="proveedor_nombre">Nombre</FieldLabel>
        <form.AppField
          name="proveedor.nombre"
          children={(field) => (
            <Input
              type="text"
              id="proveedor_nombre"
              placeholder="Nombre del provedor"
              value={field.state.value}
              onChange={(e) => {
                field.handleChange(e.target.value);
              }}
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
              placeholder="J-1234567"
              value={field.state.value}
              onChange={(e) => {
                field.handleChange(e.target.value);
              }}
            />
          )}
        />
      </Field>
      <Field orientation="vertical">
        <FieldLabel htmlFor="proveedor_telefono">Teléfono</FieldLabel>
        <form.AppField
          name="proveedor.telefono"
          children={(field) => (
            <field.Input
              id="proveedor_telefono"
              type="text"
              placeholder="+584123456789"
              value={field.state.value}
              onChange={(e) => {
                field.handleChange(e.target.value);
              }}
            />
          )}
        />
      </Field>
      <Field orientation="vertical">
        <FieldLabel htmlFor="proveedor_email">Correo Electrónico</FieldLabel>
        <form.AppField
          name="proveedor.email"
          children={(field) => (
            <field.Input
              type="email"
              id="proveedor_email"
              placeholder="email@proveedor.com"
              value={field.state.value}
              onChange={(e) => {
                field.handleChange(e.target.value);
              }}
            />
          )}
        />
      </Field>
    </>
  ),
});
