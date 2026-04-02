import { OverlayProps } from "@/components/overlay";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Proveedor } from "@/features/administracion/schemas";
import { useAppForm } from "@/hooks/useAppForm";
import { Check, Loader2 } from "lucide-react";
import { useEffect, type SubmitEventHandler } from "react";
import { toast } from "sonner";
import { ConceptoField } from "./fields/concepto-field";
import { FechaField } from "./fields/fecha-field";
import { MontoField } from "./fields/monto-field";
import { NuevoProveedorFields } from "./fields/nuevo-proveedor-fields";
import { ProveedorSelectField } from "./fields/proveedor-select-field";
import { useRegistrarGasto } from "./hooks/use-registrar-gasto";
import { defaultValues, NuevoGastoFormSchema } from "./schema";

export interface RegistrarGastoOverlayProps extends OverlayProps {
  proveedores: Pick<Proveedor, "id" | "nombre">[];
}

export function RegistrarGastoOverlay(props: RegistrarGastoOverlayProps) {
  const { registrarGasto, registrarGastoYProveedor } = useRegistrarGasto();

  const form = useAppForm({
    defaultValues,
    validators: {
      onChange: NuevoGastoFormSchema,
      onBlur: NuevoGastoFormSchema,
    },
    onSubmit: async (data) => {
      if (data.value.provedor_registrado) {
        const response = await registrarGasto.mutateAsync(data.value);
        console.log(response);
        if (response.errors) {
          toast.error(response.errors[0].message);
        }
      } else {
        const response = await registrarGastoYProveedor.mutateAsync(data.value);
        console.log(response);
        if (response.errors) {
          toast.error(response.errors[0].message);
        }
      }
    },
  });

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

  const handleSubmit: SubmitEventHandler = (event) => {
    event.preventDefault();
    event.stopPropagation();
    form.handleSubmit();
  };

  return (
    <Dialog {...props}>
      <DialogContent className="">
        <DialogHeader>
          <DialogTitle>Registrar Gasto</DialogTitle>
        </DialogHeader>
        {/* 
        <form.Subscribe
          selector={(state) => state.values}
          children={(values) => (
            <code>
              <pre>{JSON.stringify(values, null, 2)}</pre>
            </code>
          )}
        /> */}

        <form className="grid gap-5" onSubmit={handleSubmit}>
          <section className="grid gap-5">
            <ConceptoField form={form} />

            {!!props.proveedores.length && (
              <Field orientation="horizontal">
                <form.AppField
                  name="provedor_registrado"
                  children={(field) => (
                    <field.Checkbox
                      id="proveedor-no-registrado"
                      checked={!field.state.value}
                      onCheckedChange={() => {
                        field.handleChange(!field.state.value);
                        if (!field.state.value) {
                          form.setFieldValue("proveedor", "");
                        } else {
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
                  <ProveedorSelectField
                    form={form}
                    proveedores={props.proveedores}
                  />
                ) : (
                  <NuevoProveedorFields form={form} />
                )
              }
            />

            <MontoField form={form} />
            <FechaField form={form} />
          </section>

          <div className="flex gap-2 justify-end">
            <Button type="button" variant="outline">
              Cancelar
            </Button>
            <form.Subscribe
              selector={(state) => state.isValid}
              children={(isValid) => {
                const provedor_registrado = form.getFieldValue(
                  "provedor_registrado",
                );
                const concepto = form.getFieldValue("concepto");
                const proveedor = form.getFieldValue("proveedor");
                const monto = form.getFieldValue("monto");

                const isFormValid =
                  isValid &&
                  concepto.trim() !== "" &&
                  monto > 0 &&
                  (provedor_registrado
                    ? proveedor !== ""
                    : typeof proveedor === "object" &&
                      proveedor.nombre?.trim() !== "" &&
                      proveedor.rif?.trim() !== "");

                return (
                  <Button
                    type="submit"
                    disabled={form.state.isSubmitting || !isFormValid}
                  >
                    Registrar
                    {form.state.isSubmitting ? (
                      <Loader2 className="animate-spin" />
                    ) : (
                      <Check />
                    )}
                  </Button>
                );
              }}
            />
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
