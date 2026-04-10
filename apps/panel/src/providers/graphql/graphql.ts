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
  DateTime: { input: any; output: any; }
};

export type BooleanCondition = {
  eq?: InputMaybe<Scalars['Boolean']['input']>;
};

export type Cuota = {
  actualizacion: Scalars['DateTime']['output'];
  anio: Scalars['Int']['output'];
  id: Scalars['ID']['output'];
  mes: Scalars['Int']['output'];
  monto: Scalars['Int']['output'];
  pagos: PagosCount;
  registro: Scalars['DateTime']['output'];
  villas: VillasCount;
};

export type CuotaEspecial = Cuota & {
  __typename?: 'CuotaEspecial';
  actualizacion: Scalars['DateTime']['output'];
  anio: Scalars['Int']['output'];
  detalles: Proyecto;
  id: Scalars['ID']['output'];
  mes: Scalars['Int']['output'];
  monto: Scalars['Int']['output'];
  pagos: PagosCount;
  registro: Scalars['DateTime']['output'];
  villas: VillasCount;
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
  mes: Scalars['Int']['output'];
  monto: Scalars['Int']['output'];
  pagos: PagosCount;
  registro: Scalars['DateTime']['output'];
  villas: VillasCount;
};

export type CuotaType = CuotaEspecial | CuotaRegular;

export enum EstadoDeProyecto {
  Activo = 'ACTIVO',
  Borrador = 'BORRADOR',
  Cerrado = 'CERRADO'
}

export type Gasto = {
  __typename?: 'Gasto';
  concepto: Scalars['String']['output'];
  cuota?: Maybe<Scalars['ID']['output']>;
  descripcion?: Maybe<Scalars['String']['output']>;
  fecha: Scalars['DateTime']['output'];
  id: Scalars['ID']['output'];
  moneda: Scalars['String']['output'];
  monto: Scalars['Int']['output'];
  proveedor: Scalars['String']['output'];
  registrado_por: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  tasa: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type GastoWithProveedor = {
  __typename?: 'GastoWithProveedor';
  concepto: Scalars['String']['output'];
  cuota?: Maybe<Scalars['ID']['output']>;
  descripcion?: Maybe<Scalars['String']['output']>;
  fecha: Scalars['DateTime']['output'];
  id: Scalars['ID']['output'];
  moneda: Moneda;
  monto: Scalars['Int']['output'];
  proveedor: Proveedor;
  registrado_por: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  tasa: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

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

export enum MetodoDePago {
  Efectivo = 'Efectivo',
  PagoMovil = 'PagoMovil',
  Transferencia = 'Transferencia'
}

export enum Moneda {
  Usd = 'USD',
  Ved = 'VED'
}

export type Mutation = {
  __typename?: 'Mutation';
  _empty?: Maybe<Scalars['String']['output']>;
  login: LoginCredentialsDto;
  registrarGasto: Gasto;
  registrarGastoYProveedor: Gasto;
  registrarPago: Scalars['Boolean']['output'];
};


export type MutationLoginArgs = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};


export type MutationRegistrarGastoArgs = {
  input: RegistrarGastoDto;
};


export type MutationRegistrarGastoYProveedorArgs = {
  input: RegistrarGastoYProveedorDto;
};


export type MutationRegistrarPagoArgs = {
  input: RegistrarPagoDto;
};

export type ObtenerProveedoresDto = {
  and?: InputMaybe<Array<ObtenerProveedoresDto>>;
  id?: InputMaybe<StringCondition>;
  not?: InputMaybe<ObtenerProveedoresDto>;
  or?: InputMaybe<Array<ObtenerProveedoresDto>>;
};

export type Paginable = Gasto;

export type Paginated = {
  __typename?: 'Paginated';
  data: Array<Paginable>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type PaginatedCuota = {
  __typename?: 'PaginatedCuota';
  data: Array<CuotaType>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type PaginatedGasto = {
  __typename?: 'PaginatedGasto';
  data: Array<Gasto>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type PaginatedGastoWithProveedor = {
  __typename?: 'PaginatedGastoWithProveedor';
  data: Array<GastoWithProveedor>;
  limit: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
};

export type Paginator = {
  limit: Scalars['Int']['input'];
  page: Scalars['Int']['input'];
};

export type Pago = {
  __typename?: 'Pago';
  actualizacion: Scalars['DateTime']['output'];
  actualizado_por: Scalars['String']['output'];
  cuenta: Scalars['Int']['output'];
  destinado: Scalars['Int']['output'];
  fecha: Scalars['DateTime']['output'];
  id: Scalars['String']['output'];
  metodo: MetodoDePago;
  moneda: Moneda;
  monto: Scalars['Int']['output'];
  referencia: Scalars['String']['output'];
  registrado_por: Scalars['String']['output'];
  registro: Scalars['DateTime']['output'];
  tasa: Scalars['Int']['output'];
  villa: Scalars['Int']['output'];
};

export type PagosCount = {
  __typename?: 'PagosCount';
  monto_total: Scalars['Float']['output'];
  total: Scalars['Int']['output'];
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
  obtenerGastos: PaginatedGastoWithProveedor;
  obtenerProveedores: Array<Proveedor>;
  obtenerTasa: Tasa;
};


export type QueryObtenerCuotaArgs = {
  id: Scalars['String']['input'];
};


export type QueryObtenerCuotasArgs = {
  filter?: InputMaybe<CuotaFilter>;
  paginator?: InputMaybe<Paginator>;
};


export type QueryObtenerGastosArgs = {
  paginator?: InputMaybe<Paginator>;
};


export type QueryObtenerProveedoresArgs = {
  filter?: InputMaybe<ObtenerProveedoresDto>;
};

export type RegistrarGastoDto = {
  concepto: Scalars['String']['input'];
  fecha?: InputMaybe<Scalars['DateTime']['input']>;
  moneda: Moneda;
  monto: Scalars['Int']['input'];
  proveedor: Scalars['String']['input'];
};

export type RegistrarGastoYProveedorDto = {
  concepto: Scalars['String']['input'];
  fecha?: InputMaybe<Scalars['DateTime']['input']>;
  moneda: Moneda;
  monto: Scalars['Int']['input'];
  proveedor: RegistrarProveedorDto;
};

export type RegistrarPagoDto = {
  fecha?: InputMaybe<Scalars['DateTime']['input']>;
  metodo: MetodoDePago;
  moneda: Moneda;
  monto: Scalars['Int']['input'];
  referencia: Scalars['String']['input'];
  tasa: Scalars['Int']['input'];
  villa: Scalars['Int']['input'];
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

export type Tasa = {
  __typename?: 'Tasa';
  fecha: Scalars['String']['output'];
  fuente: Scalars['String']['output'];
  moneda: Scalars['String']['output'];
  tipo: Scalars['String']['output'];
  valor: Scalars['Int']['output'];
};

export enum TipoDeCuota {
  Especial = 'Especial',
  Regular = 'Regular',
  Semilla = 'Semilla'
}

export type VillasCount = {
  __typename?: 'VillasCount';
  total: Scalars['Int']['output'];
};

export type CuotaPageQueryVariables = Exact<{
  cuota_id: Scalars['String']['input'];
}>;


export type CuotaPageQuery = { __typename?: 'Query', cuota?:
    | { __typename: 'CuotaEspecial', id: string, mes: number, anio: number, pagos: { __typename?: 'PagosCount', monto_total: number, total: number }, detalles: { __typename?: 'Proyecto', titulo: string } }
    | { __typename: 'CuotaRegular', id: string, mes: number, anio: number, pagos: { __typename?: 'PagosCount', monto_total: number, total: number } }
   | null };

export type CuotasPageQueryVariables = Exact<{ [key: string]: never; }>;


export type CuotasPageQuery = { __typename?: 'Query', cuotas: { __typename?: 'PaginatedCuota', data: Array<
      | { __typename: 'CuotaEspecial', id: string, monto: number, mes: number, anio: number, registro: any, detalles: { __typename?: 'Proyecto', descripcion: string } }
      | { __typename: 'CuotaRegular', id: string, monto: number, mes: number, anio: number, registro: any }
    > } };

export type RegistrarCuotaPageQueryVariables = Exact<{ [key: string]: never; }>;


export type RegistrarCuotaPageQuery = { __typename?: 'Query', obtenerProveedores: Array<{ __typename?: 'Proveedor', id: string, nombre: string }>, obtenerGastos: { __typename?: 'PaginatedGastoWithProveedor', data: Array<{ __typename?: 'GastoWithProveedor', id: string, concepto: string, moneda: Moneda, monto: number, fecha: any, proveedor: { __typename?: 'Proveedor', id: string, nombre: string, rif: string, telefono?: string | null, email?: string | null } }> } };

export type RegistrarGastoMutationVariables = Exact<{
  input: RegistrarGastoDto;
}>;


export type RegistrarGastoMutation = { __typename?: 'Mutation', registrarGasto: { __typename?: 'Gasto', id: string } };

export type RegistrarGastoYProveedorMutationVariables = Exact<{
  input: RegistrarGastoYProveedorDto;
}>;


export type RegistrarGastoYProveedorMutation = { __typename?: 'Mutation', registrarGastoYProveedor: { __typename?: 'Gasto', id: string, concepto: string } };

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
      pagos {
        monto_total
        total
      }
    }
    ... on CuotaEspecial {
      detalles {
        titulo
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
      ... on CuotaRegular {
        id
        monto
        mes
        anio
        registro
      }
      ... on CuotaEspecial {
        id
        monto
        mes
        anio
        registro
        detalles {
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
  obtenerGastos {
    data {
      id
      concepto
      moneda
      monto
      fecha
      proveedor {
        id
        nombre
        rif
        telefono
        email
      }
    }
  }
}
    `) as unknown as TypedDocumentString<RegistrarCuotaPageQuery, RegistrarCuotaPageQueryVariables>;
export const RegistrarGastoDocument = new TypedDocumentString(`
    mutation RegistrarGasto($input: RegistrarGastoDTO!) {
  registrarGasto(input: $input) {
    id
  }
}
    `) as unknown as TypedDocumentString<RegistrarGastoMutation, RegistrarGastoMutationVariables>;
export const RegistrarGastoYProveedorDocument = new TypedDocumentString(`
    mutation RegistrarGastoYProveedor($input: RegistrarGastoYProveedorDTO!) {
  registrarGastoYProveedor(input: $input) {
    id
    concepto
  }
}
    `) as unknown as TypedDocumentString<RegistrarGastoYProveedorMutation, RegistrarGastoYProveedorMutationVariables>;