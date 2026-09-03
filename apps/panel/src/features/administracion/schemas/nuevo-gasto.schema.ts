import * as v from "valibot";
import { MetodoDeOperacion, Moneda } from "@/providers/graphql/graphql";
import { NuevoProveedorSchema } from "./nuevo-proveedor.schema";

export const NuevoGastoSchema = v.object({
  concepto: v.pipe(v.string(), v.nonEmpty("El concepto es requerido")),
  metodo: v.enum(MetodoDeOperacion),
  proveedor: v.pipe(v.string(), v.nonEmpty("El proveedor es requerido")),
  moneda: v.enum(Moneda),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  tasa: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.date(),
  comprobante: v.optional(v.string()),
});

export const NuevoGastoYProveedorSchema = v.object({
  concepto: v.pipe(v.string(), v.nonEmpty("El concepto es requerido")),
  metodo: v.enum(MetodoDeOperacion),
  proveedor: NuevoProveedorSchema,
  moneda: v.enum(Moneda),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  tasa: v.pipe(v.number(), v.integer(), v.minValue(1)),
  fecha: v.pipe(v.date()),
  comprobante: v.optional(v.string()),
});

type NuevoGastoYProveedor = v.InferOutput<typeof NuevoGastoYProveedorSchema>;
