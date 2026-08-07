/* eslint-disable */
import { DocumentTypeDecoration } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTime: { input: Date; output: Date; }
};

export type Abono = {
  __typename?: 'Abono';
  fecha: Scalars['DateTime']['output'];
  monto: Scalars['Float']['output'];
  pago: Scalars['ID']['output'];
};

export type BooleanCondition = {
  eq?: InputMaybe<Scalars['Boolean']['input']>;
};

export type Cuota = {
  actualizacion: Scalars['DateTime']['output'];
  anio: Scalars['Int']['output'];
  id: Scalars['ID']['output'];
  mes: Mes;
  monto: Scalars['Float']['output'];
  recaudacion: Recaudacion;
  registro: Scalars['DateTime']['output'];
};

export type CuotaEspecial = Cuota & {
  __typename?: 'CuotaEspecial';
  actualizacion: Scalars['DateTime']['output'];
  anio: Scalars['Int']['output'];
  detalles: Proyecto;
  id: Scalars['ID']['output'];
  mes: Mes;
  monto: Scalars['Float']['output'];
  recaudacion: Recaudacion;
  registro: Scalars['DateTime']['output'];
};

export type CuotaFilter = {
  and?: InputMaybe<Array<CuotaFilter>>;
  id?: InputMaybe<StringCondition>;
  monto?: InputMaybe<IntCondition>;
  not?: InputMaybe<CuotaFilter>;
  or?: InputMaybe<Array<CuotaFilter>>;
};

export type CuotaRegular = Cuota & {
  __typename?: 'CuotaRegular';
  actualizacion: Scalars['DateTime']['output'];
  anio: Scalars['Int']['output'];
  id: Scalars['ID']['output'];
  mes: Mes;
  monto: Scalars['Float']['output'];
  recaudacion: Recaudacion;
  registro: Scalars['DateTime']['output'];
};

export type CuotaType = CuotaEspecial | CuotaRegular;

export type Deuda = {
  __typename?: 'Deuda';
  abonos?: Maybe<Array<Abono>>;
  cuota: Scalars['ID']['output'];
  deuda: Scalars['Float']['output'];
  estado: EstadoDeDeuda;
  id: Scalars['String']['output'];
  monto: Scalars['Float']['output'];
  registro: Scalars['DateTime']['output'];
  unidad: Scalars['ID']['output'];
};

export type Ente = Sujeto & {
  __typename?: 'Ente';
  actualizacion: Scalars['DateTime']['output'];
  cedula: Scalars['String']['output'];
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  razon_social: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  representante: Persona;
  telefono: Scalars['String']['output'];
};

export enum EstadoDeDeuda {
  Abonada = 'ABONADA',
  Pendiente = 'PENDIENTE',
  Saldada = 'SALDADA'
}

export enum EstadoDeProyecto {
  Activo = 'ACTIVO',
  Borrador = 'BORRADOR',
  Cerrado = 'CERRADO'
}

export enum EstadoDeUnidad {
  Activa = 'ACTIVA',
  EnLitigio = 'EN_LITIGIO',
  Exenta = 'EXENTA',
  Inhabitada = 'INHABITADA',
  Preventa = 'PREVENTA',
  Suspendida = 'SUSPENDIDA'
}

export type IntCondition = {
  eq?: InputMaybe<Scalars['Int']['input']>;
  gt?: InputMaybe<Scalars['Int']['input']>;
  gte?: InputMaybe<Scalars['Int']['input']>;
  lt?: InputMaybe<Scalars['Int']['input']>;
  lte?: InputMaybe<Scalars['Int']['input']>;
};

export type LoginCredentialsDto = {
  __typename?: 'LoginCredentialsDTO';
  token: Scalars['String']['output'];
};

export enum Mes {
  Abril = 'ABRIL',
  Agosto = 'AGOSTO',
  Diciembre = 'DICIEMBRE',
  Enero = 'ENERO',
  Febrero = 'FEBRERO',
  Julio = 'JULIO',
  Junio = 'JUNIO',
  Marzo = 'MARZO',
  Mayo = 'MAYO',
  Noviembre = 'NOVIEMBRE',
  Octubre = 'OCTUBRE',
  Septiembre = 'SEPTIEMBRE'
}

export enum MetodoDeTransaccion {
  Cheque = 'Cheque',
  Compensacion = 'Compensacion',
  Efectivo = 'Efectivo',
  PagoMovil = 'PagoMovil',
  TransferenciaInternacional = 'TransferenciaInternacional',
  TransferenciaNacional = 'TransferenciaNacional'
}

export enum Moneda {
  Usd = 'USD',
  Ved = 'VED'
}

export type Movimiento = {
  cuota?: Maybe<Scalars['ID']['output']>;
  id: Scalars['String']['output'];
  monto: Scalars['Float']['output'];
  tipo: TipoDeMovimiento;
};

export type MovimientoACondominio = Movimiento & {
  __typename?: 'MovimientoACondominio';
  cuota?: Maybe<Scalars['ID']['output']>;
  id: Scalars['String']['output'];
  monto: Scalars['Float']['output'];
  tipo: TipoDeMovimiento;
};

export type MovimientoAProveedor = Movimiento & {
  __typename?: 'MovimientoAProveedor';
  cuota?: Maybe<Scalars['ID']['output']>;
  id: Scalars['String']['output'];
  monto: Scalars['Float']['output'];
  proveedor: Proveedor;
  tipo: TipoDeMovimiento;
};

export type MovimientoAUnidad = Movimiento & {
  __typename?: 'MovimientoAUnidad';
  cuota?: Maybe<Scalars['ID']['output']>;
  id: Scalars['String']['output'];
  monto: Scalars['Float']['output'];
  tipo: TipoDeMovimiento;
  unidad: Unidad;
};

export type MovimientoType = MovimientoACondominio | MovimientoAProveedor | MovimientoAUnidad;

export type Mutation = {
  __typename?: 'Mutation';
  _empty?: Maybe<Scalars['String']['output']>;
  login: LoginCredentialsDto;
  registrarCuota: CuotaType;
  registrarGasto?: Maybe<Transaccion>;
};


export type MutationLoginArgs = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};


export type MutationRegistrarCuotaArgs = {
  input: RegistrarCuotaDto;
};


export type MutationRegistrarGastoArgs = {
  input: RegistrarGastoDto;
};

export type ObtenerProveedoresDto = {
  and?: InputMaybe<Array<ObtenerProveedoresDto>>;
  id?: InputMaybe<StringCondition>;
  not?: InputMaybe<ObtenerProveedoresDto>;
  or?: InputMaybe<Array<ObtenerProveedoresDto>>;
};

export type PaginatedCuota = {
  __typename?: 'PaginatedCuota';
  data: Array<CuotaType>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type PaginatedDeuda = {
  __typename?: 'PaginatedDeuda';
  data: Array<Deuda>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type PaginatedTransaccion = {
  __typename?: 'PaginatedTransaccion';
  data: Array<Transaccion>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type PaginatedUnidad = {
  __typename?: 'PaginatedUnidad';
  data: Array<Unidad>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type Paginator = {
  limit: Scalars['Int']['input'];
  page: Scalars['Int']['input'];
};

export type Persona = Sujeto & {
  __typename?: 'Persona';
  actualizacion: Scalars['DateTime']['output'];
  apellidos: Scalars['String']['output'];
  cedula: Scalars['String']['output'];
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  nombres: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  telefono: Scalars['String']['output'];
};

export type Proveedor = {
  __typename?: 'Proveedor';
  actualizado_en: Scalars['DateTime']['output'];
  creado_en: Scalars['DateTime']['output'];
  direccion?: Maybe<Scalars['String']['output']>;
  email?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  nombre: Scalars['String']['output'];
  rif: Scalars['String']['output'];
  telefono?: Maybe<Scalars['String']['output']>;
};

export type Proyecto = {
  __typename?: 'Proyecto';
  actualizacion: Scalars['DateTime']['output'];
  descripcion: Scalars['String']['output'];
  estado: EstadoDeProyecto;
  fecha_limite: Scalars['DateTime']['output'];
  interes_por_mora: Scalars['Float']['output'];
  justificacion: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  titulo: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  _empty?: Maybe<Scalars['String']['output']>;
  obtenerCuota?: Maybe<CuotaType>;
  obtenerCuotas: PaginatedCuota;
  obtenerDeudasDeUnaUnidad: PaginatedDeuda;
  obtenerDeudasDeUnaUnidadPorCodigo: PaginatedDeuda;
  obtenerMovimientos: PaginatedTransaccion;
  obtenerProveedores: Array<Proveedor>;
  obtenerTasa: Tasa;
  obtenerUnidad?: Maybe<Unidad>;
  obtenerUnidadPorCodigo?: Maybe<Unidad>;
  obtenerUnidades?: Maybe<PaginatedUnidad>;
  obtenerUnidadesEstadisticas?: Maybe<UnidadesTotales>;
};


export type QueryObtenerCuotaArgs = {
  id: Scalars['String']['input'];
};


export type QueryObtenerCuotasArgs = {
  filter?: InputMaybe<CuotaFilter>;
  paginator?: InputMaybe<Paginator>;
};


export type QueryObtenerDeudasDeUnaUnidadArgs = {
  id: Scalars['ID']['input'];
  paginator?: InputMaybe<Paginator>;
};


export type QueryObtenerDeudasDeUnaUnidadPorCodigoArgs = {
  codigo: Scalars['String']['input'];
  paginator?: InputMaybe<Paginator>;
};


export type QueryObtenerMovimientosArgs = {
  filter?: InputMaybe<TransaccionFilter>;
  paginator?: InputMaybe<Paginator>;
  tipo?: InputMaybe<TipoDeMovimiento>;
};


export type QueryObtenerProveedoresArgs = {
  filter?: InputMaybe<ObtenerProveedoresDto>;
};


export type QueryObtenerTasaArgs = {
  anio?: InputMaybe<Scalars['Int']['input']>;
  dia?: InputMaybe<Scalars['Int']['input']>;
  mes?: InputMaybe<Mes>;
};


export type QueryObtenerUnidadArgs = {
  id: Scalars['ID']['input'];
};


export type QueryObtenerUnidadPorCodigoArgs = {
  codigo: Scalars['String']['input'];
};


export type QueryObtenerUnidadesArgs = {
  filter?: InputMaybe<UnidadFilter>;
  paginator?: InputMaybe<Paginator>;
};

export type Recaudacion = {
  __typename?: 'Recaudacion';
  moneda: Moneda;
  monto_estimado: Scalars['Float']['output'];
  monto_pendiente: Scalars['Float']['output'];
  monto_recaudado: Scalars['Float']['output'];
  pagos_asociados: Scalars['Int']['output'];
  unidades: Scalars['Int']['output'];
  unidades_aplicadas: Scalars['Int']['output'];
  unidades_pendientes: Scalars['Int']['output'];
  unidades_solventes: Scalars['Int']['output'];
};

export type RegistrarCuotaDto = {
  anio?: InputMaybe<Scalars['Int']['input']>;
  descripcion?: InputMaybe<Scalars['String']['input']>;
  fecha_limite?: InputMaybe<Scalars['DateTime']['input']>;
  justificacion?: InputMaybe<Scalars['String']['input']>;
  mes?: InputMaybe<Mes>;
  monto_total: Scalars['Int']['input'];
  tipo: TipoDeCuota;
  titulo?: InputMaybe<Scalars['String']['input']>;
};

export type RegistrarGastoDto = {
  consepto: Scalars['String']['input'];
  fecha?: InputMaybe<Scalars['DateTime']['input']>;
  metodo: MetodoDeTransaccion;
  moneda?: InputMaybe<Moneda>;
  monto: Scalars['Int']['input'];
  proveedor: Scalars['ID']['input'];
  referencia?: InputMaybe<Scalars['String']['input']>;
  tasa: Scalars['Int']['input'];
};

export type RegistrarProveedorDto = {
  direccion?: InputMaybe<Scalars['String']['input']>;
  email: Scalars['String']['input'];
  nombre: Scalars['String']['input'];
  rif: Scalars['String']['input'];
  telefono: Scalars['String']['input'];
};

export type StringCondition = {
  eq?: InputMaybe<Scalars['String']['input']>;
  in?: InputMaybe<Array<InputMaybe<Scalars['String']['input']>>>;
  like?: InputMaybe<Scalars['String']['input']>;
  regex?: InputMaybe<Scalars['String']['input']>;
};

export type Sujeto = {
  actualizacion: Scalars['DateTime']['output'];
  cedula: Scalars['String']['output'];
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  telefono: Scalars['String']['output'];
};

export type Tasa = {
  __typename?: 'Tasa';
  fecha: Scalars['String']['output'];
  fuente: Scalars['String']['output'];
  moneda: Scalars['String']['output'];
  tipo: Scalars['String']['output'];
  valor: Scalars['Float']['output'];
};

export enum TipoDeCuota {
  Especial = 'Especial',
  Regular = 'Regular',
  Semilla = 'Semilla'
}

export enum TipoDeMovimiento {
  Credito = 'Credito',
  Debito = 'Debito'
}

export type Titular = Ente | Persona;

export type Transaccion = {
  __typename?: 'Transaccion';
  concepto: Scalars['String']['output'];
  fecha: Scalars['DateTime']['output'];
  id: Scalars['String']['output'];
  metodo: MetodoDeTransaccion;
  moneda: Moneda;
  monto_total: Scalars['Float']['output'];
  movimientos: Array<MovimientoType>;
  registrado_por: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  tasa: Scalars['Float']['output'];
};

export type TransaccionFilter = {
  and?: InputMaybe<Array<TransaccionFilter>>;
  concepto?: InputMaybe<StringCondition>;
  cuota_id?: InputMaybe<StringCondition>;
  not?: InputMaybe<TransaccionFilter>;
  or?: InputMaybe<Array<TransaccionFilter>>;
};

export type Unidad = {
  __typename?: 'Unidad';
  codigo: Scalars['String']['output'];
  contacto?: Maybe<Persona>;
  deuda: Scalars['Float']['output'];
  estado: EstadoDeUnidad;
  id: Scalars['String']['output'];
  titular_primario?: Maybe<Titular>;
  wallet: Scalars['Float']['output'];
};

export type UnidadFilter = {
  and?: InputMaybe<Array<UnidadFilter>>;
  codigo?: InputMaybe<StringCondition>;
  estado?: InputMaybe<StringCondition>;
  id?: InputMaybe<StringCondition>;
  not?: InputMaybe<UnidadFilter>;
  or?: InputMaybe<Array<UnidadFilter>>;
};

export type UnidadesTotales = {
  __typename?: 'UnidadesTotales';
  total_asignado: Scalars['Float']['output'];
  total_pendiente: Scalars['Float']['output'];
  total_unidades: Scalars['Int']['output'];
  unidades_activas: Scalars['Int']['output'];
  unidades_con_pendientes: Scalars['Int']['output'];
  unidades_en_litigio: Scalars['Int']['output'];
  unidades_exentas: Scalars['Int']['output'];
  unidades_inhabitadas: Scalars['Int']['output'];
  unidades_preventa: Scalars['Int']['output'];
  unidades_solventes: Scalars['Int']['output'];
  unidades_suspendidas: Scalars['Int']['output'];
};

export type CuotaPageQueryVariables = Exact<{
  cuota_id: Scalars['String']['input'];
}>;


export type CuotaPageQuery = { __typename?: 'Query', cuota?:
    | { __typename: 'CuotaEspecial', id: string, mes: Mes, anio: number, monto: number, recaudacion: { __typename?: 'Recaudacion', unidades_aplicadas: number, unidades_solventes: number, monto_estimado: number, monto_recaudado: number }, detalles: { __typename?: 'Proyecto', titulo: string, descripcion: string, justificacion: string } }
    | { __typename: 'CuotaRegular', id: string, mes: Mes, anio: number, monto: number, recaudacion: { __typename?: 'Recaudacion', unidades_aplicadas: number, unidades_solventes: number, monto_estimado: number, monto_recaudado: number } }
   | null };

export type CuotasPageQueryVariables = Exact<{ [key: string]: never; }>;


export type CuotasPageQuery = { __typename?: 'Query', cuotas: { __typename?: 'PaginatedCuota', data: Array<
      | { __typename: 'CuotaEspecial', id: string, monto: number, mes: Mes, anio: number, registro: Date, recaudacion: { __typename?: 'Recaudacion', unidades_aplicadas: number, pagos_asociados: number }, detalles: { __typename?: 'Proyecto', titulo: string, descripcion: string } }
      | { __typename: 'CuotaRegular', id: string, monto: number, mes: Mes, anio: number, registro: Date, recaudacion: { __typename?: 'Recaudacion', unidades_aplicadas: number, pagos_asociados: number } }
    > } };

export type RegistrarCuotaPageQueryVariables = Exact<{ [key: string]: never; }>;


export type RegistrarCuotaPageQuery = { __typename?: 'Query', obtenerProveedores: Array<{ __typename?: 'Proveedor', id: string, nombre: string }> };

export type DashboardPageQueryVariables = Exact<{ [key: string]: never; }>;


export type DashboardPageQuery = { __typename?: 'Query', proveedores: Array<{ __typename?: 'Proveedor', id: string, nombre: string }> };

export type LoginMutationVariables = Exact<{
  email: Scalars['String']['input'];
  pass: Scalars['String']['input'];
}>;


export type LoginMutation = { __typename?: 'Mutation', login: { __typename?: 'LoginCredentialsDTO', token: string } };

export type RegistrarPagoPageQueryVariables = Exact<{
  codigo_like: Scalars['String']['input'];
}>;


export type RegistrarPagoPageQuery = { __typename?: 'Query', unidades?: { __typename?: 'PaginatedUnidad', data: Array<{ __typename?: 'Unidad', id: string, codigo: string }> } | null };

export type VillasPageQueryVariables = Exact<{ [key: string]: never; }>;


export type VillasPageQuery = { __typename?: 'Query', estadisticas?: { __typename?: 'UnidadesTotales', total_unidades: number, unidades_activas: number, unidades_con_pendientes: number, unidades_inhabitadas: number } | null, villas?: { __typename?: 'PaginatedUnidad', data: Array<{ __typename?: 'Unidad', codigo: string, contacto?: { __typename?: 'Persona', id: string, email: string, telefono: string } | null, titular_primario?:
        | { __typename: 'Ente', id: string, razon_social: string }
        | { __typename: 'Persona', id: string, nombres: string, apellidos: string }
       | null }> } | null };

export type RegistrarGastoOverlayMutationVariables = Exact<{
  input: RegistrarGastoDto;
}>;


export type RegistrarGastoOverlayMutation = { __typename?: 'Mutation', registrarGasto?: { __typename?: 'Transaccion', id: string, concepto: string } | null };

export type BuscarGastosHuerfanosQueryVariables = Exact<{
  busqueda?: InputMaybe<Scalars['String']['input']>;
}>;


export type BuscarGastosHuerfanosQuery = { __typename?: 'Query', gastos: { __typename?: 'PaginatedTransaccion', data: Array<{ __typename?: 'Transaccion', id: string, monto_total: number, concepto: string, metodo: MetodoDeTransaccion, moneda: Moneda, tasa: number, fecha: Date, movimientos: Array<
        | { __typename: 'MovimientoACondominio', id: string, monto: number, tipo: TipoDeMovimiento }
        | { __typename: 'MovimientoAProveedor', id: string, monto: number, tipo: TipoDeMovimiento, proveedor: { __typename?: 'Proveedor', id: string, nombre: string } }
        | { __typename: 'MovimientoAUnidad', id: string, monto: number, tipo: TipoDeMovimiento }
      > }> } };

export class TypedDocumentString<TResult, TVariables>
  extends String
  implements DocumentTypeDecoration<TResult, TVariables>
{
  __apiType?: NonNullable<DocumentTypeDecoration<TResult, TVariables>['__apiType']>;
  private value: string;
  public __meta__?: Record<string, any> | undefined;

  constructor(value: string, __meta__?: Record<string, any> | undefined) {
    super(value);
    this.value = value;
    this.__meta__ = __meta__;
  }

  override toString(): string & DocumentTypeDecoration<TResult, TVariables> {
    return this.value;
  }
}

export const CuotaPageDocument = new TypedDocumentString(`
    query CuotaPage($cuota_id: String!) {
  cuota: obtenerCuota(id: $cuota_id) {
    __typename
    ... on Cuota {
      id
      mes
      anio
      monto
      recaudacion {
        unidades_aplicadas
        unidades_solventes
        monto_estimado
        monto_recaudado
      }
    }
    ... on CuotaEspecial {
      detalles {
        titulo
        descripcion
        justificacion
      }
    }
  }
}
    `) as unknown as TypedDocumentString<CuotaPageQuery, CuotaPageQueryVariables>;
export const CuotasPageDocument = new TypedDocumentString(`
    query CuotasPage {
  cuotas: obtenerCuotas {
    data {
      __typename
      ... on Cuota {
        id
        monto
        mes
        anio
        registro
        recaudacion {
          unidades_aplicadas
          pagos_asociados
        }
      }
      ... on CuotaEspecial {
        detalles {
          titulo
          descripcion
        }
      }
    }
  }
}
    `) as unknown as TypedDocumentString<CuotasPageQuery, CuotasPageQueryVariables>;
export const RegistrarCuotaPageDocument = new TypedDocumentString(`
    query RegistrarCuotaPage {
  obtenerProveedores {
    id
    nombre
  }
}
    `) as unknown as TypedDocumentString<RegistrarCuotaPageQuery, RegistrarCuotaPageQueryVariables>;
export const DashboardPageDocument = new TypedDocumentString(`
    query DashboardPage {
  proveedores: obtenerProveedores {
    id
    nombre
  }
}
    `) as unknown as TypedDocumentString<DashboardPageQuery, DashboardPageQueryVariables>;
export const LoginDocument = new TypedDocumentString(`
    mutation Login($email: String!, $pass: String!) {
  login(email: $email, password: $pass) {
    token
  }
}
    `) as unknown as TypedDocumentString<LoginMutation, LoginMutationVariables>;
export const RegistrarPagoPageDocument = new TypedDocumentString(`
    query RegistrarPagoPage($codigo_like: String!) {
  unidades: obtenerUnidades(
    filter: {codigo: {like: $codigo_like}}
    paginator: {limit: 5, page: 1}
  ) {
    data {
      id
      codigo
    }
  }
}
    `) as unknown as TypedDocumentString<RegistrarPagoPageQuery, RegistrarPagoPageQueryVariables>;
export const VillasPageDocument = new TypedDocumentString(`
    query VillasPage {
  estadisticas: obtenerUnidadesEstadisticas {
    total_unidades
    unidades_activas
    unidades_con_pendientes
    unidades_inhabitadas
  }
  villas: obtenerUnidades(filter: {estado: {eq: "ACTIVA"}}) {
    data {
      codigo
      contacto {
        id
        email
        telefono
      }
      titular_primario {
        __typename
        ... on Sujeto {
          id
        }
        ... on Persona {
          nombres
          apellidos
        }
        ... on Ente {
          razon_social
        }
      }
    }
  }
}
    `) as unknown as TypedDocumentString<VillasPageQuery, VillasPageQueryVariables>;
export const RegistrarGastoOverlayDocument = new TypedDocumentString(`
    mutation RegistrarGastoOverlay($input: RegistrarGastoDTO!) {
  registrarGasto(input: $input) {
    id
    concepto
  }
}
    `) as unknown as TypedDocumentString<RegistrarGastoOverlayMutation, RegistrarGastoOverlayMutationVariables>;
export const BuscarGastosHuerfanosDocument = new TypedDocumentString(`
    query BuscarGastosHuerfanos($busqueda: String) {
  gastos: obtenerMovimientos(tipo: Debito, filter: {concepto: {like: $busqueda}}) {
    data {
      id
      monto_total
      concepto
      metodo
      moneda
      tasa
      fecha
      movimientos {
        __typename
        ... on Movimiento {
          id
          monto
          tipo
        }
        ... on MovimientoAProveedor {
          proveedor {
            id
            nombre
          }
        }
      }
    }
  }
}
    `) as unknown as TypedDocumentString<BuscarGastosHuerfanosQuery, BuscarGastosHuerfanosQueryVariables>;