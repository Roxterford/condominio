import * as v from "valibot";
import { TipoDeSujeto } from "@/providers/graphql/graphql";

const BaseSujeto = {
  documento_identidad: v.pipe(
    v.string(),
    v.trim(),
    v.minLength(3, "Ingrese el documento de identidad"),
  ),
  email: v.pipe(v.string(), v.trim(), v.email("Ingrese un email válido")),
  telefono: v.pipe(v.string(), v.trim(), v.minLength(7, "Ingrese el teléfono")),
};

export const RegistrarTitularSchema = v.variant("tipo", [
  v.object({
    ...BaseSujeto,
    tipo: v.literal(TipoDeSujeto.PersonaNatural),
    nombres: v.pipe(
      v.string(),
      v.trim(),
      v.minLength(2, "Ingrese los nombres"),
    ),
    apellidos: v.pipe(
      v.string(),
      v.trim(),
      v.minLength(2, "Ingrese los apellidos"),
    ),
  }),
  v.object({
    ...BaseSujeto,
    tipo: v.literal(TipoDeSujeto.EnteJuridico),
    razon_social: v.pipe(
      v.string(),
      v.trim(),
      v.minLength(3, "Ingrese la razón social"),
    ),
  }),
]);

export type RegistrarTitularForm = v.InferOutput<typeof RegistrarTitularSchema>;

export const registrarTitularDefaultValues: RegistrarTitularForm = {
  tipo: TipoDeSujeto.PersonaNatural,
  documento_identidad: "",
  nombres: "",
  apellidos: "",
  email: "",
  telefono: "",
};
