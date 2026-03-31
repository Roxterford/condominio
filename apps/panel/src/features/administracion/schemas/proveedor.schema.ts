import * as v from "valibot";

export const ProveedorSchema = v.object({
  id: v.string(),
  rif: v.string(),
  nombre: v.string(),
  tipo: v.string(),
  email: v.string(),
  telefono: v.string(),
  direccion: v.optional(v.string()),
  creado_en: v.string(),
  actualizado_en: v.string(),
});

export type Proveedor = v.InferOutput<typeof ProveedorSchema>;

export const NuevoProveedorSchema = v.object({
  nombre: v.pipe(v.string(), v.nonEmpty("El nombre es requerido")),
  rif: v.pipe(v.string(), v.nonEmpty("El rif es requerido")),
  telefono: v.pipe(v.string(), v.nonEmpty("El telefono es requerido")),
  email: v.optional(v.string()),
});
