import * as v from "valibot";
import { MetodoDeTransaccion, Moneda } from "@/providers/graphql/graphql";
import { NuevoProveedorSchema } from "./nuevo-proveedor.schema";

export const NuevoGastoSchema = v.object({
  concepto: v.pipe(v.string(), v.nonEmpty("El concepto es requerido")),
  metodo: v.enum(MetodoDeTransaccion),
  proveedor: v.pipe(v.string(), v.nonEmpty("El proveedor es requerido")),
  moneda: v.enum(Moneda),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  tasa: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.date(),
  comprobante: v.optional(v.string()),
});

export type NuevoGasto = v.InferOutput<typeof NuevoGastoSchema>;

export const NuevoGastoYProveedorSchema = v.object({
  concepto: v.pipe(v.string(), v.nonEmpty("El concepto es requerido")),
  metodo: v.enum(MetodoDeTransaccion),
  proveedor: NuevoProveedorSchema,
  moneda: v.enum(Moneda),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  tasa: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.pipe(v.date()),
  comprobante: v.optional(v.string()),
});

export type NuevoGastoYProveedor = v.InferOutput<
  typeof NuevoGastoYProveedorSchema
>;
