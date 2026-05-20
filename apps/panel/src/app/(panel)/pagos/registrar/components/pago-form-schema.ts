import * as v from "valibot";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";

export const NuevoPagoFormSchema = v.object({
  unidad: v.pipe(v.string(), v.nonEmpty()),
  fecha: v.nullish(v.date()),
  metodo: v.nullish(v.enum(MetodoDePago)),
  referencia: v.string(),
  monto: v.pipe(v.number(), v.integer(), v.minValue(1)),
  tasa: v.pipe(v.number(), v.integer(), v.minValue(1)),
  moneda: v.enum(Moneda),
});

export type NuevoPagoForm = v.InferOutput<typeof NuevoPagoFormSchema>;

export const registrarPagoDefaultValues: NuevoPagoForm = {
  unidad: "",
  fecha: null,
  metodo: null,
  referencia: "",
  monto: 0,
  tasa: 0,
  moneda: Moneda.USD,
};
