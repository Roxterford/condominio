/* eslint-disable */
import * as types from './graphql';



/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n  query CuotaPage($cuota_id: String!) {\n    cuota: obtenerCuota(id: $cuota_id) {\n      __typename\n      ... on Cuota {\n        id\n        mes\n        anio\n        pagos {\n          monto_total\n          total\n        }\n      }\n      ... on CuotaEspecial {\n        detalles {\n          titulo\n        }\n      }\n    }\n  }\n": typeof types.CuotaPageDocument,
    "\n  query CuotasPage {\n    cuotas: obtenerCuotas {\n      data {\n        __typename\n        ... on CuotaRegular {\n          id\n          monto\n          mes\n          anio\n          registro\n        }\n\n        ... on CuotaEspecial {\n          id\n          monto\n          mes\n          anio\n          registro\n          detalles {\n            descripcion\n          }\n        }\n      }\n    }\n  }\n": typeof types.CuotasPageDocument,
    "\n  query RegistrarCuotaPage {\n    obtenerProveedores {\n      id\n      nombre\n    }\n    obtenerGastos {\n      data {\n        id\n        concepto\n        moneda\n        monto\n        fecha\n        proveedor {\n          id\n          nombre\n          rif\n          telefono\n          email\n        }\n      }\n    }\n  }\n": typeof types.RegistrarCuotaPageDocument,
    "\n  mutation RegistrarGasto($input: RegistrarGastoDTO!) {\n    registrarGasto(input: $input) {\n      id\n    }\n  }\n": typeof types.RegistrarGastoDocument,
    "\n  mutation RegistrarGastoYProveedor($input: RegistrarGastoYProveedorDTO!) {\n    registrarGastoYProveedor(input: $input) {\n      id\n      concepto\n    }\n  }\n": typeof types.RegistrarGastoYProveedorDocument,
};
const documents: Documents = {
    "\n  query CuotaPage($cuota_id: String!) {\n    cuota: obtenerCuota(id: $cuota_id) {\n      __typename\n      ... on Cuota {\n        id\n        mes\n        anio\n        pagos {\n          monto_total\n          total\n        }\n      }\n      ... on CuotaEspecial {\n        detalles {\n          titulo\n        }\n      }\n    }\n  }\n": types.CuotaPageDocument,
    "\n  query CuotasPage {\n    cuotas: obtenerCuotas {\n      data {\n        __typename\n        ... on CuotaRegular {\n          id\n          monto\n          mes\n          anio\n          registro\n        }\n\n        ... on CuotaEspecial {\n          id\n          monto\n          mes\n          anio\n          registro\n          detalles {\n            descripcion\n          }\n        }\n      }\n    }\n  }\n": types.CuotasPageDocument,
    "\n  query RegistrarCuotaPage {\n    obtenerProveedores {\n      id\n      nombre\n    }\n    obtenerGastos {\n      data {\n        id\n        concepto\n        moneda\n        monto\n        fecha\n        proveedor {\n          id\n          nombre\n          rif\n          telefono\n          email\n        }\n      }\n    }\n  }\n": types.RegistrarCuotaPageDocument,
    "\n  mutation RegistrarGasto($input: RegistrarGastoDTO!) {\n    registrarGasto(input: $input) {\n      id\n    }\n  }\n": types.RegistrarGastoDocument,
    "\n  mutation RegistrarGastoYProveedor($input: RegistrarGastoYProveedorDTO!) {\n    registrarGastoYProveedor(input: $input) {\n      id\n      concepto\n    }\n  }\n": types.RegistrarGastoYProveedorDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query CuotaPage($cuota_id: String!) {\n    cuota: obtenerCuota(id: $cuota_id) {\n      __typename\n      ... on Cuota {\n        id\n        mes\n        anio\n        pagos {\n          monto_total\n          total\n        }\n      }\n      ... on CuotaEspecial {\n        detalles {\n          titulo\n        }\n      }\n    }\n  }\n"): typeof import('./graphql').CuotaPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query CuotasPage {\n    cuotas: obtenerCuotas {\n      data {\n        __typename\n        ... on CuotaRegular {\n          id\n          monto\n          mes\n          anio\n          registro\n        }\n\n        ... on CuotaEspecial {\n          id\n          monto\n          mes\n          anio\n          registro\n          detalles {\n            descripcion\n          }\n        }\n      }\n    }\n  }\n"): typeof import('./graphql').CuotasPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query RegistrarCuotaPage {\n    obtenerProveedores {\n      id\n      nombre\n    }\n    obtenerGastos {\n      data {\n        id\n        concepto\n        moneda\n        monto\n        fecha\n        proveedor {\n          id\n          nombre\n          rif\n          telefono\n          email\n        }\n      }\n    }\n  }\n"): typeof import('./graphql').RegistrarCuotaPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation RegistrarGasto($input: RegistrarGastoDTO!) {\n    registrarGasto(input: $input) {\n      id\n    }\n  }\n"): typeof import('./graphql').RegistrarGastoDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation RegistrarGastoYProveedor($input: RegistrarGastoYProveedorDTO!) {\n    registrarGastoYProveedor(input: $input) {\n      id\n      concepto\n    }\n  }\n"): typeof import('./graphql').RegistrarGastoYProveedorDocument;


export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}
