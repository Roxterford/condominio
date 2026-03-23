import * as v from "valibot";

const NuevoProveedorSchema = v.object({
  nombre: v.string(),
  rif: v.string(),
  telefono: v.string(),
  email: v.optional(v.string())
});

export const NuevoGastoSchema = v.object({
  concepto: v.string(),
  proveedor: v.string(),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.date(),
  comprobante: v.optional(v.string())
});

export const NuevoGastoYProveedorSchema = v.object({
  concepto: v.string(),
  proveedor: NuevoProveedorSchema,
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.date(),
  comprobante: v.optional(v.string())
});

