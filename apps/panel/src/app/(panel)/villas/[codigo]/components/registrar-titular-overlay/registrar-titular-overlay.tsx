"use client";

import { OverlayProps } from "@/components/overlay";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Field,
  FieldContent,
  FieldError,
  FieldLabel,
  FieldSet,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAppForm } from "@/hooks/useAppForm";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { TipoDeSujeto } from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Check } from "lucide-react";
import { type SubmitEventHandler } from "react";
import {
  registrarTitularDefaultValues,
  RegistrarTitularSchema,
} from "./schema";
import { RegistrarSujetoDto } from "@/providers/graphql/graphql";

const RegistrarTitularMutation = graphql(/* GraphQL */ `
  mutation RegistrarTitular($input: RegistrarSujetoDTO!) {
    registrarSujeto(input: $input) {
      id
      display_name
    }
  }
`);

export function RegistrarTitularOverlay(props: OverlayProps) {
  const router = useRouter();

  const registrar = useMutation({
    mutationKey: ["villas.registrar-titular"],
    mutationFn: (input: RegistrarSujetoDto) =>
      execute(RegistrarTitularMutation, { input }),
  });

  const form = useAppForm({
    defaultValues: registrarTitularDefaultValues,
    validators: {
      onChange: RegistrarTitularSchema,
      onBlur: RegistrarTitularSchema,
    },
    onSubmit: async ({ value }) => {
      const res = await registrar.mutateAsync({
        tipo: value.tipo,
        documento_identidad: value.documento_identidad,
        email: value.email,
        telefono: value.telefono,
        ...(value.tipo === TipoDeSujeto.PersonaNatural
          ? { nombres: value.nombres, apellidos: value.apellidos }
          : { razon_social: value.razon_social }),
      });

      if (res.errors?.length) {
        return toast.error(res.errors.at(0)?.message);
      }

      form.reset();
      toast.success("Titular registrado", {
        description: `${res.data?.registrarSujeto.display_name} ya forma parte de la propiedad`,
      });
      router.refresh();
      props.onDone?.();
    },
  });

  const handleSubmit: SubmitEventHandler = (event) => {
    event.preventDefault();
    event.stopPropagation();
    form.handleSubmit();
  };

  return (
    <Dialog {...props}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Registrar titular</DialogTitle>
          <DialogDescription className="text-sm text-muted-foreground">
            Añade una persona natural o ente jurídico como titular de la
            propiedad.
          </DialogDescription>
        </DialogHeader>

        <form className="grid gap-5" onSubmit={handleSubmit}>
          <form.AppField
            name="tipo"
            children={(field) => (
              <Tabs
                value={field.state.value}
                onValueChange={(value) =>
                  field.handleChange(value as TipoDeSujeto)
                }
              >
                <TabsList className="w-full">
                  <TabsTrigger value={TipoDeSujeto.PersonaNatural}>
                    Persona natural
                  </TabsTrigger>
                  <TabsTrigger value={TipoDeSujeto.EnteJuridico}>
                    Ente jurídico
                  </TabsTrigger>
                </TabsList>
              </Tabs>
            )}
          />

          <FieldSet className="gap-4">
            <form.AppField name="documento_identidad">
              {(field) => (
                <Field>
                  <FieldLabel>Documento de identidad (RIF/C.I.)</FieldLabel>
                  <FieldContent>
                    <Input
                      placeholder="Ej. V-12345678 o J-123456789"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      onBlur={field.handleBlur}
                    />
                    <FieldError errors={field.state.meta.errors} />
                  </FieldContent>
                </Field>
              )}
            </form.AppField>

            <form.Subscribe
              selector={(state) => state.values.tipo}
              children={(tipo) =>
                tipo === TipoDeSujeto.PersonaNatural ? (
                  <>
                    <form.AppField name="nombres">
                      {(field) => (
                        <Field>
                          <FieldLabel>Nombres</FieldLabel>
                          <FieldContent>
                            <Input
                              placeholder="Ej. María Fernanda"
                              value={field.state.value}
                              onChange={(e) =>
                                field.handleChange(e.target.value)
                              }
                              onBlur={field.handleBlur}
                            />
                            <FieldError errors={field.state.meta.errors} />
                          </FieldContent>
                        </Field>
                      )}
                    </form.AppField>
                    <form.AppField name="apellidos">
                      {(field) => (
                        <Field>
                          <FieldLabel>Apellidos</FieldLabel>
                          <FieldContent>
                            <Input
                              placeholder="Ej. González Pérez"
                              value={field.state.value}
                              onChange={(e) =>
                                field.handleChange(e.target.value)
                              }
                              onBlur={field.handleBlur}
                            />
                            <FieldError errors={field.state.meta.errors} />
                          </FieldContent>
                        </Field>
                      )}
                    </form.AppField>
                  </>
                ) : (
                  <form.AppField name="razon_social">
                    {(field) => (
                      <Field>
                        <FieldLabel>Razón social</FieldLabel>
                        <FieldContent>
                          <Input
                            placeholder="Ej. Inversiones La Casona, C.A."
                            value={field.state.value}
                            onChange={(e) => field.handleChange(e.target.value)}
                            onBlur={field.handleBlur}
                          />
                          <FieldError errors={field.state.meta.errors} />
                        </FieldContent>
                      </Field>
                    )}
                  </form.AppField>
                )
              }
            />

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <form.AppField name="email">
                {(field) => (
                  <Field>
                    <FieldLabel>Email</FieldLabel>
                    <FieldContent>
                      <Input
                        type="email"
                        placeholder="ej. titular@correo.com"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        onBlur={field.handleBlur}
                      />
                      <FieldError errors={field.state.meta.errors} />
                    </FieldContent>
                  </Field>
                )}
              </form.AppField>
              <form.AppField name="telefono">
                {(field) => (
                  <Field>
                    <FieldLabel>Teléfono</FieldLabel>
                    <FieldContent>
                      <Input
                        placeholder="Ej. 0414-1234567"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        onBlur={field.handleBlur}
                      />
                      <FieldError errors={field.state.meta.errors} />
                    </FieldContent>
                  </Field>
                )}
              </form.AppField>
            </div>
          </FieldSet>

          <div className="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => props.onOpenChange?.(false)}
            >
              Cancelar
            </Button>
            <form.Subscribe
              selector={(state) => state.isValid}
              children={(isValid) => (
                <Button type="submit" disabled={!isValid}>
                  Registrar titular
                  {registrar.isPending ? (
                    <Spinner />
                  ) : (
                    <Check className="ml-1 h-4 w-4" />
                  )}
                </Button>
              )}
            />
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
