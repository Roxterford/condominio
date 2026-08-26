import * as v from "valibot";
import { EstadoDeUnidad } from "@/providers/graphql/graphql";

export const ESTADOS_UNIDAD: { value: EstadoDeUnidad; label: string }[] = [
  { value: EstadoDeUnidad.Activa, label: "Activa" },
  { value: EstadoDeUnidad.Inhabitada, label: "Inhabilitada" },
  { value: EstadoDeUnidad.Exenta, label: "Exenta" },
  { value: EstadoDeUnidad.EnLitigio, label: "En litigio" },
  { value: EstadoDeUnidad.Suspendida, label: "Suspendida" },
  { value: EstadoDeUnidad.Preventa, label: "Preventa" },
];

export const RegistrarUnidadFormSchema = v.object({
  codigo: v.pipe(v.string(), v.minLength(1, "El código es requerido")),
  estado: v.enum(EstadoDeUnidad),
  descripcion: v.string(),
  titular_primario: v.string(),
  contacto: v.string(),
});

export type RegistrarUnidadForm = v.InferOutput<typeof RegistrarUnidadFormSchema>;

export const defaultValues: RegistrarUnidadForm = {
  codigo: "",
  estado: EstadoDeUnidad.Activa,
  descripcion: "",
  titular_primario: "",
  contacto: "",
};
