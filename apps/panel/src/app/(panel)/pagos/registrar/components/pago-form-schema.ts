import * as v from "valibot";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/providers/graphql/graphql";

export const NuevoPagoFormSchema = v.object({
  unidad: v.pipe(v.string(), v.nonEmpty()),
  fecha: v.nullish(v.date()),
  metodo: v.enum(MetodoDePago),
  referencia: v.pipe(v.string(), v.nonEmpty()),
  monto: v.pipe(v.number(), v.minValue(1), v.transform(redondearCentavos)),
  tasa: v.pipe(v.number(), v.minValue(0), v.transform(redondearCentavos)),
  moneda: v.enum(Moneda),
});

export type NuevoPagoForm = v.InferOutput<typeof NuevoPagoFormSchema>;

export const registrarPagoDefaultValues: NuevoPagoForm = {
  unidad: "",
  fecha: null,
  metodo: null as any,
  referencia: "",
  monto: 0,
  tasa: 0,
  moneda: Moneda.Usd,
};

function redondearCentavos(n: number): number {
  return Math.round(n * 100);
}
