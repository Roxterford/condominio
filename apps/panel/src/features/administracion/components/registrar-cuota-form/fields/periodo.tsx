import { Field, FieldLabel } from "@/components/ui/field";
import {
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { withForm } from "@/hooks/useAppForm";
import { Loader2 } from "lucide-react";
import { useEffect, useLayoutEffect } from "react";
import { Mes } from "@/providers/graphql/graphql";
import { defaultValues, MESES } from "../schema";

const useIsomorphicLayoutEffect =
  typeof window !== "undefined" ? useLayoutEffect : useEffect;

export interface PeriodoFieldProps {
  periodos?: Map<number, Set<Mes>>;
  isLoading?: boolean;
}

export const PeriodoField = withForm({
  props: {} as PeriodoFieldProps,
  defaultValues,

  render: ({
    form,
    isLoading = false,
    periodos = new Map<number, Set<Mes>>(),
  }) => {
    const anios = Array.from(periodos.keys());

    useIsomorphicLayoutEffect(() => {
      if (isLoading) return;
      const anio = form.getFieldValue("anio");
      const disponibles = periodos.get(anio);
      if (!disponibles || disponibles.size === 0) return;
      const mesActual = form.getFieldValue("mes");
      if (!disponibles.has(mesActual)) {
        form.setFieldValue("mes", Array.from(disponibles)[0]);
      }
    }, [isLoading, periodos, form]);

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
                      disabled={usarAnioActual || isLoading}
                      onValueChange={(value) => {
                        const anio = Number(value);
                        field.handleChange(anio);
                        const meses = periodos.get(anio);
                        form.setFieldValue(
                          "mes",
                          meses && meses.size > 0
                            ? Array.from(meses)[0]
                            : form.getFieldValue("mes"),
                        );
                      }}
                    >
                      <SelectTrigger
                        className={
                          isLoading
                            ? "relative [&>svg:last-child]:hidden"
                            : undefined
                        }
                      >
                        <SelectValue placeholder="Seleccione el año" />
                        {isLoading && (
                          <Loader2 className="absolute top-1/2 right-2 size-4 -translate-y-1/2 animate-spin text-muted-foreground" />
                        )}
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
            <form.Subscribe
              selector={(state) => state.values.anio}
              children={(anio) => {
                const mesesDisponibles = periodos.get(anio) ?? new Set<Mes>();

                return (
                  <form.AppField
                    name="mes"
                    children={(field) => (
                      <field.Select
                        value={field.state.value}
                        disabled={isLoading}
                        onValueChange={(value) =>
                          field.handleChange(value as Mes)
                        }
                      >
                        <SelectTrigger
                          className={
                            isLoading
                              ? "relative [&>svg:last-child]:hidden"
                              : undefined
                          }
                        >
                          <SelectValue placeholder="Seleccione el mes" />
                          {isLoading && (
                            <Loader2 className="absolute top-1/2 right-2 size-4 -translate-y-1/2 animate-spin text-muted-foreground" />
                          )}
                        </SelectTrigger>
                        <SelectContent>
                          <SelectGroup>
                            {MESES.map((mes) => (
                              <SelectItem
                                key={mes.value}
                                value={mes.value}
                                disabled={!mesesDisponibles.has(mes.value)}
                              >
                                {mes.label}
                              </SelectItem>
                            ))}
                          </SelectGroup>
                        </SelectContent>
                      </field.Select>
                    )}
                  />
                );
              }}
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
                  onChange={(e) => e && field.handleChange(e)}
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
                  onChange={(e) => e && field.handleChange(e)}
                />
              </Field>
            )}
          />
        </div>
      </section>
    );
  },
});
