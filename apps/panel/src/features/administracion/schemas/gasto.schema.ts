import * as v from "valibot";
import { MonedaSchema } from "./moneda.schema";
import { NuevoProveedorSchema } from "./proveedor.schema";

export const NuevoGastoSchema = v.object({
  concepto: v.pipe(v.string(), v.nonEmpty("El concepto es requerido")),
  proveedor: v.pipe(v.string(), v.nonEmpty("El proveedor es requerido")),
  moneda: MonedaSchema,
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.date(),
  comprobante: v.optional(v.string()),
});

export type NuevoGasto = v.InferOutput<typeof NuevoGastoSchema>;

export const NuevoGastoYProveedorSchema = v.object({
  concepto: v.pipe(v.string(), v.nonEmpty("El concepto es requerido")),
  proveedor: NuevoProveedorSchema,
  moneda: MonedaSchema,
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.pipe(v.date()),
  comprobante: v.optional(v.string()),
});

export type NuevoGastoYProveedor = v.InferOutput<
  typeof NuevoGastoYProveedorSchema
>;
