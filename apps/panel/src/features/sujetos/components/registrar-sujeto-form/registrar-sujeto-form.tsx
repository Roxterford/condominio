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
	RegistrarSujetoDocument,
	RegistrarSujetoDto,
} from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { Check, Loader2 } from "lucide-react";
import { useState, type SubmitEventHandler } from "react";
import { toast } from "sonner";
import { TIPOS_SUJETO, defaultValues, RegistrarSujetoFormSchema } from "./schema";

const RegistrarSujetoMutation = RegistrarSujetoDocument;

export function RegistrarSujetoForm() {
  const router = useRouter();
  const [submitting, setSubmitting] = useState(false);

  const registrar = useMutation({
    mutationFn: (input: RegistrarSujetoDto) =>
      execute(RegistrarSujetoMutation, { input }),
  });

  const form = useAppForm({
    defaultValues,
    validators: {
      onChange: RegistrarSujetoFormSchema,
      onBlur: RegistrarSujetoFormSchema,
    },
    onSubmit: async ({ value }) => {
      setSubmitting(true);
      const res = await registrar.mutateAsync({
        tipo: value.tipo,
        documento_identidad: value.documento_identidad,
        nombres: value.nombres || undefined,
        apellidos: value.apellidos || undefined,
        razon_social: value.razon_social || undefined,
        email: value.email,
        telefono: value.telefono,
        representante: value.representante || undefined,
      });

      setSubmitting(false);

      if (res.errors?.length) {
        return toast.error(res.errors.at(0)?.message);
      }

      toast.success("Propietario registrado con éxito");
      router.push("/villas");
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
        <FieldLabel>Tipo</FieldLabel>
        <form.AppField
          name="tipo"
          children={(field) => (
            <field.Select
              value={field.state.value}
              onValueChange={(value) =>
                field.handleChange(value as typeof field.state.value)
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Seleccione tipo" />
              </SelectTrigger>
              <SelectContent>
                {TIPOS_SUJETO.map((tipo) => (
                  <SelectItem key={tipo.value} value={tipo.value}>
                    {tipo.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </field.Select>
          )}
        />
      </Field>

      <Field orientation="vertical">
        <FieldLabel htmlFor="documento_identidad">Documento de identidad</FieldLabel>
        <form.AppField
          name="documento_identidad"
          children={(field) => (
            <field.Input
              id="documento_identidad"
              placeholder="Ej: V-12345678"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>

      <div className="flex gap-2">
        <Field orientation="vertical" className="flex-1">
          <FieldLabel htmlFor="nombres">Nombres</FieldLabel>
          <form.AppField
            name="nombres"
            children={(field) => (
              <field.Input
                id="nombres"
                placeholder="Nombres"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          />
        </Field>
        <Field orientation="vertical" className="flex-1">
          <FieldLabel htmlFor="apellidos">Apellidos</FieldLabel>
          <form.AppField
            name="apellidos"
            children={(field) => (
              <field.Input
                id="apellidos"
                placeholder="Apellidos"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          />
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel htmlFor="razon_social">Razón social</FieldLabel>
        <form.AppField
          name="razon_social"
          children={(field) => (
            <field.Input
              id="razon_social"
              placeholder="Solo para entes jurídicos"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          )}
        />
      </Field>

      <div className="flex gap-2">
        <Field orientation="vertical" className="flex-1">
          <FieldLabel htmlFor="email">Correo electrónico</FieldLabel>
          <form.AppField
            name="email"
            children={(field) => (
              <field.Input
                id="email"
                type="email"
                placeholder="correo@ejemplo.com"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          />
        </Field>
        <Field orientation="vertical" className="flex-1">
          <FieldLabel htmlFor="telefono">Teléfono</FieldLabel>
          <form.AppField
            name="telefono"
            children={(field) => (
              <field.Input
                id="telefono"
                placeholder="+584123456789"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          />
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel htmlFor="representante">Representante (ID)</FieldLabel>
        <form.AppField
          name="representante"
          children={(field) => (
            <field.Input
              id="representante"
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
