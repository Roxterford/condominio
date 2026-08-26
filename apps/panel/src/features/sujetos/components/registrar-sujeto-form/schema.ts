import * as v from "valibot";
import { TipoDeSujeto } from "@/providers/graphql/graphql";

export const TIPOS_SUJETO: { value: TipoDeSujeto; label: string }[] = [
  { value: TipoDeSujeto.PersonaNatural, label: "Persona Natural" },
  { value: TipoDeSujeto.EnteJuridico, label: "Ente Jurídico" },
];

export const RegistrarSujetoFormSchema = v.object({
  tipo: v.enum(TipoDeSujeto),
  documento_identidad: v.pipe(
    v.string(),
    v.minLength(1, "El documento es requerido"),
  ),
  nombres: v.string(),
  apellidos: v.string(),
  razon_social: v.string(),
  email: v.pipe(v.string(), v.email("Correo inválido")),
  telefono: v.pipe(v.string(), v.minLength(1, "El teléfono es requerido")),
  representante: v.string(),
});

export type RegistrarSujetoForm = v.InferOutput<typeof RegistrarSujetoFormSchema>;

export const defaultValues: RegistrarSujetoForm = {
  tipo: TipoDeSujeto.PersonaNatural,
  documento_identidad: "",
  nombres: "",
  apellidos: "",
  razon_social: "",
  email: "",
  telefono: "",
  representante: "",
};
