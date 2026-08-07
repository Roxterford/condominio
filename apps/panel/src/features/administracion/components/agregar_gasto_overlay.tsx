import { OverlayProps } from "@/components/overlay/overlay";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Badge } from "@/components/ui/badge";
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldLabel,
  FieldTitle,
} from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { SubmitEvent } from "react";

export interface AgregarGastoOverlayProps extends OverlayProps {
  onSelect?(selection: "todos" | "seleccionar" | "nuevo"): void;
}

export function AgregarGastoOverlay(props: AgregarGastoOverlayProps) {
  const handleSubmit = (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    props.onSelect?.(e.target["agregar"].value);
  };

  return (
    <Dialog {...props}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Agregar Gasto</DialogTitle>
          <DialogDescription>
            Seleccione una de las siguientes opciones
          </DialogDescription>
        </DialogHeader>
        <form className="grid gap-5" onSubmit={handleSubmit}>
          <section>
            <RadioGroup name="agregar" defaultValue="selecionar">
              <FieldLabel htmlFor="no-asociados">
                <Field orientation="horizontal">
                  <FieldContent>
                    <div className="flex gap-5">
                      <FieldTitle>Seleccionar gastos no asociados</FieldTitle>
                      <Badge className="bg-green-100 text-green-800 dark:bg-green-950 dark:text-green-300">
                        Recomendado
                      </Badge>
                    </div>
                    <FieldDescription>
                      Vincule gastos huerfanos registrados previamente
                    </FieldDescription>
                  </FieldContent>
                  <RadioGroupItem value="seleccionar" id="no-asociados" />
                </Field>
              </FieldLabel>
              <FieldLabel htmlFor="todos-no-asociados">
                <Field orientation="horizontal">
                  <FieldContent>
                    <FieldTitle>Añadir todos los gastos</FieldTitle>
                    <FieldDescription>
                      Añade todos los gastos huerfanos automaticamente
                    </FieldDescription>
                  </FieldContent>
                  <RadioGroupItem value="todos" id="todos-no-asociados" />
                </Field>
              </FieldLabel>
              <FieldLabel htmlFor="registrar-nuevo">
                <Field orientation="horizontal">
                  <FieldContent>
                    <FieldTitle>Registrar nuevo gasto</FieldTitle>
                    <FieldDescription>
                      Registre un nuevo gasto en el sistema
                    </FieldDescription>
                  </FieldContent>
                  <RadioGroupItem value="nuevo" id="registrar-nuevo" />
                </Field>
              </FieldLabel>
            </RadioGroup>
          </section>
          <div className="flex gap-2 justify-end">
            <Button type="button" variant="outline">
              Cancelar
            </Button>
            <Button type="submit">Siguiente</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
