import * as v from 'valibot';
import { Mes } from './mes';

export enum TipoDeCuota {
  Regular = 'REGULAR',
  Especial = 'ESPECIAL',
  Semilla = 'SEMILLA',
}

const NuevaCuotaRegularSchema = v.object({
  tipo: v.literal(TipoDeCuota.Regular),
  periodo: v.object({
    actual: v.boolean(),
    anio: v.number(),
    mes: v.enum(Mes),
    fecha_emision: v.date(),
    fecha_limite: v.date(),
  }),
  notas_adicionales: v.optional(v.string()),
})

const NuevaCuotaEspecialSchema = v.object({
  tipo: v.literal(TipoDeCuota.Especial),
  notas_adicionales: v.optional(v.string()),
})

export const NuevaCuetoSchema = v.variant('tipo', [
  NuevaCuotaRegularSchema,
  NuevaCuotaEspecialSchema,
])


v.object({
  tipo: v.enum(TipoDeCuota),
  monto: v.number(),
  fecha: v.date(),
  descripcion: v.string(),
})




