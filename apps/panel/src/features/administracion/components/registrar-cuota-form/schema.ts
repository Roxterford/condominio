import * as v from "valibot";
import {
  Mes,
  Movimiento,
  TipoDeCuota,
  Transaccion,
} from "@/providers/graphql/graphql";

export const EstrategiaDeDistribucion = [
  "lineal",
  "individual",
  "alicuota",
] as const;

export type EstrategiaDeDistribucion =
  (typeof EstrategiaDeDistribucion)[number];

export const MESES: { value: Mes; label: string }[] = [
  { value: Mes.Enero, label: "Enero" },
  { value: Mes.Febrero, label: "Febrero" },
  { value: Mes.Marzo, label: "Marzo" },
  { value: Mes.Abril, label: "Abril" },
  { value: Mes.Mayo, label: "Mayo" },
  { value: Mes.Junio, label: "Junio" },
  { value: Mes.Julio, label: "Julio" },
  { value: Mes.Agosto, label: "Agosto" },
  { value: Mes.Septiembre, label: "Septiembre" },
  { value: Mes.Octubre, label: "Octubre" },
  { value: Mes.Noviembre, label: "Noviembre" },
  { value: Mes.Diciembre, label: "Diciembre" },
];

export type GastoSeleccionado = Partial<
  Omit<Transaccion, "movimientos"> & {
    movimientos: Array<Partial<Movimiento>>;
  }
>;

export const RegistrarCuotaFormSchema = v.object({
  tipo: v.enum(TipoDeCuota),
  anio_actual: v.boolean(),
  anio: v.pipe(v.number(), v.minValue(2000), v.maxValue(2100)),
  mes: v.enum(Mes),
  fecha_emision: v.date(),
  fecha_limite: v.date(),
  estrategia: v.picklist(EstrategiaDeDistribucion),
  gastos: v.pipe(v.array(v.string()), v.minLength(1)),
});

export type RegistrarCuotaForm = v.InferOutput<
  typeof RegistrarCuotaFormSchema
>;

export const defaultValues: RegistrarCuotaForm = {
  tipo: TipoDeCuota.Regular,
  anio_actual: false,
  anio: new Date().getFullYear(),
  mes: Mes.Enero,
  fecha_emision: new Date(),
  fecha_limite: new Date(),
  estrategia: "lineal",
  gastos: [],
};
