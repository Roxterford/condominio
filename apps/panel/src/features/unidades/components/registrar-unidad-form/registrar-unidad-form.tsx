"use client";

import { Button } from "@/components/ui/button";
import { Field, FieldLabel } from "@/components/ui/field";
import {
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAppForm } from "@/hooks/useAppForm";
import { execute } from "@/providers/graphql/execute";
import {
	RegistrarUnidadDocument,
	RegistrarUnidadDto,
} from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { Check, Loader2 } from "lucide-react";
import { useState, type SubmitEventHandler } from "react";
import { toast } from "sonner";
import { ESTADOS_UNIDAD, defaultValues, RegistrarUnidadFormSchema } from "./schema";

const RegistrarUnidadMutation = RegistrarUnidadDocument;

export function RegistrarUnidadForm() {
  const router = useRouter();
  const [submitting, setSubmitting] = useState(false);

  const registrar = useMutation({
    mutationFn: (input: RegistrarUnidadDto) =>
      execute(RegistrarUnidadMutation, { input }),
  });

  const form = useAppForm({
    defaultValues,
    validators: {
      onChange: RegistrarUnidadFormSchema,
      onBlur: RegistrarUnidadFormSchema,
    },
    onSubmit: async ({ value }) => {
      setSubmitting(true);
      const res = await registrar.mutateAsync({
        codigo: value.codigo,
        estado: value.estado,
        descripcion: value.descripcion || undefined,
        titular_primario: value.titular_primario || undefined,
        contacto: value.contacto || undefined,
      });

      setSubmitting(false);

      if (res.errors?.length) {
        return toast.error(res.errors.at(0)?.message);
      }

      toast.success("Unidad registrada con éxito");
      const unidad = res.data?.registrarUnidad;
      if (unidad) {
        router.push(`/villas/${unidad.codigo}`);
      }
    },
  });

  const handleSubmit: SubmitEventHandler = (event) => {
    event.preventDefault();
    event.stopPropagation();
    form.handleSubmit();
  };

  return (
    <form className="grid gap-5" onSubmit={handleSubmit}>
      <Field orientation="vertical">
        <FieldLabel htmlFor="codigo">Código</FieldLabel>
        <form.AppField
          name="codigo"
          children={(field) => (
            <field.Input
              id="codigo"
              placeholder="Ej: A-101"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>

      <Field orientation="vertical">
        <FieldLabel>Estado</FieldLabel>
        <form.AppField
          name="estado"
          children={(field) => (
            <field.Select
              value={field.state.value}
              onValueChange={(value) =>
                field.handleChange(value as typeof field.state.value)
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Seleccione estado" />
              </SelectTrigger>
              <SelectContent>
                {ESTADOS_UNIDAD.map((estado) => (
                  <SelectItem key={estado.value} value={estado.value}>
                    {estado.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </field.Select>
          )}
        />
      </Field>

      <Field orientation="vertical">
        <FieldLabel htmlFor="descripcion">Descripción</FieldLabel>
        <form.AppField
          name="descripcion"
          children={(field) => (
            <field.Input
              id="descripcion"
              placeholder="Ej: Apartamento primer piso"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>

      <Field orientation="vertical">
        <FieldLabel htmlFor="titular_primario">Titular principal (ID)</FieldLabel>
        <form.AppField
          name="titular_primario"
          children={(field) => (
            <field.Input
              id="titular_primario"
              placeholder="ID del sujeto (opcional)"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>

      <Field orientation="vertical">
        <FieldLabel htmlFor="contacto">Contacto (ID)</FieldLabel>
        <form.AppField
          name="contacto"
          children={(field) => (
            <field.Input
              id="contacto"
              placeholder="ID del sujeto (opcional)"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>

      <div className="flex gap-2 justify-end">
        <Button type="button" variant="outline" onClick={() => form.reset()}>
          Cancelar
        </Button>
        <form.Subscribe
          selector={(state) => state.isValid}
          children={(isValid) => (
            <Button type="submit" disabled={submitting || !isValid}>
              Registrar
              {submitting ? <Loader2 className="animate-spin" /> : <Check />}
            </Button>
          )}
        />
      </div>
    </form>
  );
}
