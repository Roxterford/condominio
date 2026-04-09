export enum TipoDeCuota {
  REGULAR = "REGULAR",
  ESPECIAL = "ESPECIAL",
  SEMILLA = "SEMILLA",
}

export interface Cuota<T extends TipoDeCuota = TipoDeCuota> {
  id: string;
  monto: number;
  tipo: T;
  mes: number;
  anio: number;
  registro: Date;
  actualizacion: Date;
  detalles: T extends TipoDeCuota.ESPECIAL ? Proyecto : never;
}

export enum EstadoDeProyecto {
  BORRADOR = "BORRADOR",
  ACTIVO = "ACTIVO",
  CERRADO = "CERRADO",
}

export interface Proyecto {
  estado: EstadoDeProyecto;
  descripcion: string;
  justificacion: string;
  fecha_limite: Date;
  interes_por_mora: number;
  registro: Date;
  actualizacion: Date;
}
