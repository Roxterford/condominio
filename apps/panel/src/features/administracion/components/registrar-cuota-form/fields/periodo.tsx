import { Field, FieldLabel } from "@/components/ui/field";
import {
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { withForm } from "@/hooks/useAppForm";
import { Mes } from "@/providers/graphql/graphql";
import { defaultValues, MESES } from "../schema";

export const PeriodoField = withForm({
  defaultValues,
  render: ({ form }) => {
    const anioActual = new Date().getFullYear();
    const anios = [anioActual - 1, anioActual, anioActual + 1];

    return (
      <section>
        <h3>Periodo</h3>
        <div className="flex gap-5">
          <form.Subscribe
            selector={(state) => state.values.anio_actual}
            children={(usarAnioActual) => (
              <Field>
                <FieldLabel>Año</FieldLabel>
                <form.AppField
                  name="anio"
                  children={(field) => (
                    <field.Select
                      value={String(field.state.value)}
                      disabled={usarAnioActual}
                      onValueChange={(value) =>
                        field.handleChange(Number(value))
                      }
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Seleccione el año" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          {anios.map((anio) => (
                            <SelectItem key={anio} value={String(anio)}>
                              {anio}
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </field.Select>
                  )}
                />
              </Field>
            )}
          />
          <Field>
            <FieldLabel>Mes</FieldLabel>
            <form.AppField
              name="mes"
              children={(field) => (
                <field.Select
                  value={field.state.value}
                  onValueChange={(value) => field.handleChange(value as Mes)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Seleccione el mes" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {MESES.map((mes) => (
                        <SelectItem key={mes.value} value={mes.value}>
                          {mes.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </field.Select>
              )}
            />
          </Field>
        </div>

        <form.AppField
          name="anio_actual"
          children={(field) => (
            <Field orientation="horizontal">
              <field.Checkbox
                id="anio_actual"
                checked={field.state.value}
                onCheckedChange={(checked) => {
                  field.handleChange(checked === true);
                  if (checked) {
                    form.setFieldValue("anio", new Date().getFullYear());
                  }
                }}
              />
              <FieldLabel htmlFor="anio_actual">Usar año actual</FieldLabel>
            </Field>
          )}
        />

        <div className="flex gap-5">
          <form.AppField
            name="fecha_emision"
            children={(field) => (
              <Field>
                <FieldLabel htmlFor="fecha_emision">
                  Fecha de emisión
                </FieldLabel>
                <field.DatePickerInput
                  placeholder="dd/mm/aaaa"
                  locale="es-VE"
                  id="fecha_emision"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e as Date)}
                />
              </Field>
            )}
          />
          <form.AppField
            name="fecha_limite"
            children={(field) => (
              <Field>
                <FieldLabel htmlFor="fecha_limite">
                  Fecha limite de pago
                </FieldLabel>
                <field.DatePickerInput
                  placeholder="dd/mm/aaaa"
                  locale="es-VE"
                  id="fecha_limite"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e as Date)}
                />
              </Field>
            )}
          />
        </div>
      </section>
    );
  },
});