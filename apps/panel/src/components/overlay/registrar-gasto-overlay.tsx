import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useAppForm } from "@/hooks/useAppForm";
import {
  NuevoGastoSchema,
  NuevoGastoYProveedorSchema,
} from "@/schemas/nuevo_gasto.schema";
import { Check, DollarSign } from "lucide-react";
import * as v from "valibot";
import { Button } from "../ui/button";
import { DatePickerInput } from "../ui/date-picker-input";
import { Field, FieldDescription, FieldLabel } from "../ui/field";
import { Input } from "../ui/input";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { OverlayProps } from "./overlay";

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
  fecha: new Date(),
};

export function RegistrarGastoOverlay(props: OverlayProps) {
  const form = useAppForm({
    defaultValues,
    validators: {
      onBlur: NuevoGastoFormSchema,
    },
  });

  return (
    <Dialog {...props}>
      <DialogContent className="">
        <DialogHeader>
          <DialogTitle>Registrar Gasto</DialogTitle>
        </DialogHeader>

        <form className="grid gap-5">
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
                    }}
                  />
                )}
              />
              <FieldLabel htmlFor="proveedor-no-registrado">
                Proveedor no registrado
              </FieldLabel>
            </Field>

            <form.Subscribe
              selector={(state) => state.values.provedor_registrado}
              children={(provedor_registrado) =>
                provedor_registrado ? (
                  <Field orientation="vertical">
                    <FieldLabel>Provedor</FieldLabel>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="Seleccione" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="option-1">Opción 1</SelectItem>
                        <SelectItem value="option-2">Opción 2</SelectItem>
                        <SelectItem value="option-3">Opción 3</SelectItem>
                      </SelectContent>
                    </Select>
                  </Field>
                ) : (
                  <>
                    <Field orientation="vertical">
                      <FieldLabel htmlFor="proveedor_nombre">Nombre</FieldLabel>
                      <Input
                        type="text"
                        id="proveedor_nombre"
                        placeholder="Nombre del provedor"
                      />
                    </Field>
                    <Field orientation="vertical">
                      <FieldLabel htmlFor="proveedor_rif">CI / RIF</FieldLabel>
                      <Input
                        type="text"
                        id="proveedor_rif"
                        placeholder="J-1234567"
                      />
                    </Field>
                    <Field orientation="vertical">
                      <FieldLabel htmlFor="proveedor_telefono">
                        Teléfono
                      </FieldLabel>
                      <Input
                        type="text"
                        id="proveedor_telefono"
                        placeholder="+584123456789"
                      />
                    </Field>
                    <Field orientation="vertical">
                      <FieldLabel htmlFor="proveedor_email">
                        Correo Electrónico
                      </FieldLabel>
                      <Input
                        type="email"
                        id="proveedor_email"
                        placeholder="email@proveedor.com"
                      />
                    </Field>
                  </>
                )
              }
            />
            <Field orientation="vertical">
              <FieldLabel htmlFor="monto">Monto</FieldLabel>
              <InputGroup>
                <InputGroupAddon>
                  <DollarSign />
                </InputGroupAddon>
                {/* TODO: Agregar mascara de moneda */}
                <InputGroupInput type="number" id="monto" placeholder="0,00" />
              </InputGroup>
            </Field>
            <Field className="">
              <FieldLabel htmlFor="fecha_emision">Fecha</FieldLabel>
              <DatePickerInput
                placeholder="dd/mm/aaaa"
                locale="es-VE"
                id="fecha_emision"
              />
              <FieldDescription>
                {/* TODO: Agregar fecha relativa a la seleccionada */}
              </FieldDescription>
            </Field>
            <Field className="">
              <FieldLabel htmlFor="comprobante">Comprobante</FieldLabel>
              <Input type="file" />
              <FieldDescription>
                Archivos permitidos: PDF, JPG, PNG
              </FieldDescription>
            </Field>
          </section>

          <div className="flex gap-2 justify-end">
            <Button type="button" variant="outline">
              Cancelar
            </Button>
            <Button type="submit">
              Registrar
              <Check />
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
