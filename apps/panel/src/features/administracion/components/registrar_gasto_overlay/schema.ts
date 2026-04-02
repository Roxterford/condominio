import * as v from "valibot";
import {
  NuevoGastoSchema,
  NuevoGastoYProveedorSchema,
} from "@/features/administracion/schemas/gasto.schema";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";

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
  provedor_registrado: true,
  concepto: "",
  proveedor: "",
  monto: 0,
  moneda: Moneda.USD,
  fecha: new Date(),
};
