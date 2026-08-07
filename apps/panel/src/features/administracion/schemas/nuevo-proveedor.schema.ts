import * as v from "valibot";


export const NuevoProveedorSchema = v.object({
  nombre: v.pipe(v.string(), v.nonEmpty("El nombre es requerido")),
  rif: v.pipe(v.string(), v.nonEmpty("El rif es requerido")),
  telefono: v.pipe(v.string(), v.nonEmpty("El telefono es requerido")),
  email: v.optional(v.string()),
});
