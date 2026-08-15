import * as v from "valibot";
import {
  NuevoGastoSchema,
  NuevoGastoYProveedorSchema,
} from "@/features/administracion/schemas/nuevo-gasto.schema";
import { MetodoDeOperacion, Moneda } from "@/providers/graphql/graphql";

export const NuevoGastoFormSchema = v.variant("provedor_registrado", [
  v.object({
    provedor_registrado: v.literal(true),
    ...NuevoGastoSchema.entries,
  }),
  v.object({
    provedor_registrado: v.literal(false),
    ...NuevoGastoYProveedorSchema.entries,
  }),
]);

export type NuevoGastoForm = v.InferOutput<typeof NuevoGastoFormSchema>;

export const defaultValues: NuevoGastoForm = {
  metodo: MetodoDeOperacion.TransferenciaNacional,
  tasa: 0,
  provedor_registrado: true,
  concepto: "",
  proveedor: "",
  monto: 0,
  moneda: Moneda.Usd,
  fecha: new Date(),
};
