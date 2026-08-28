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
    "\n  query CuotaPage($cuota_id: String!) {\n    cuota: obtenerCuota(id: $cuota_id) {\n      __typename\n      ... on Cuota {\n        id\n        mes\n        anio\n        monto\n        gastos {\n          __typename\n          ... on Gasto {\n            operacion\n            concepto\n            moneda\n            monto\n            fecha\n            tasa\n            total\n          }\n          ... on GastoAProveedor {\n            proveedor {\n              id\n              nombre\n              rif\n              telefono\n              email\n            }\n          }\n        }\n        recaudacion {\n          unidades_aplicadas\n          unidades_solventes\n          monto_estimado\n          monto_recaudado\n        }\n      }\n      ... on CuotaEspecial {\n        detalles {\n          titulo\n          descripcion\n          justificacion\n        }\n      }\n    }\n    deudas: obtenerDeudas(filtro: { cuota: { eq: $cuota_id } }) {\n      data {\n        unidad {\n          codigo\n        }\n        titular {\n          display_name\n        }\n        deuda\n        estado\n      }\n    }\n  }\n": typeof types.CuotaPageDocument,
    "\n  query CuotasPage {\n    cuotas: obtenerCuotas {\n      data {\n        __typename\n        ... on Cuota {\n          id\n          monto\n          mes\n          anio\n          registro\n          recaudacion {\n            unidades_aplicadas\n            pagos_asociados\n          }\n        }\n\n        ... on CuotaEspecial {\n          detalles {\n            titulo\n            descripcion\n          }\n        }\n      }\n    }\n  }\n": typeof types.CuotasPageDocument,
    "\n  query RegistrarCuotaPage {\n    obtenerProveedores {\n      id\n      nombre\n    }\n  }\n": typeof types.RegistrarCuotaPageDocument,
    "\n  query DashboardPage {\n    proveedores: obtenerProveedores {\n      id\n      nombre\n    }\n  }\n": typeof types.DashboardPageDocument,
    "\n  mutation Login($email: String!, $pass: String!) {\n    login(email: $email, password: $pass) {\n      token\n    }\n  }\n": typeof types.LoginDocument,
    "\n  query RegistrarPagoPage($codigo_like: String!) {\n    unidades: obtenerUnidades(\n      filter: { codigo: { like: $codigo_like } }\n      paginator: { limit: 5, page: 1 }\n    ) {\n      data {\n        id\n        codigo\n      }\n    }\n  }\n": typeof types.RegistrarPagoPageDocument,
    "\n  query VillaPage($codigo: String!) {\n    unidad: obtenerUnidadPorCodigo(codigo: $codigo) {\n      codigo\n      estado\n      wallet\n      titular_primario {\n        __typename\n        ... on Sujeto {\n          id\n          cedula\n          display_name\n        }\n      }\n      titulares {\n        __typename\n\n        ... on Sujeto {\n          id\n          cedula\n          display_name\n          email\n        }\n      }\n    }\n    pagos: obtenerPagos(filtro: { unidad: { eq: $codigo } }) {\n      data {\n        operacion\n        concepto\n        monto\n        moneda\n      }\n    }\n    deudas: obtenerDeudas(filtro: { unidad: { eq: $codigo } }) {\n      data {\n        id\n        cuota {\n          __typename\n          ... on Deuda__Cuota {\n            id\n            nombre\n          }\n        }\n        deuda\n        monto\n        estado\n      }\n    }\n  }\n": typeof types.VillaPageDocument,
    "\n  query VillasPage {\n    estadisticas: obtenerUnidadesEstadisticas {\n      total_unidades\n      unidades_activas\n      unidades_con_pendientes\n      unidades_inhabitadas\n    }\n\n    villas: obtenerUnidades(filter: { estado: { eq: \"ACTIVA\" } }) {\n      data {\n        codigo\n        estado\n        wallet\n        contacto {\n          id\n          email\n          telefono\n        }\n        titular_primario {\n          __typename\n          ... on Sujeto {\n            id\n            cedula\n            display_name\n          }\n          ... on Persona {\n            nombres\n            apellidos\n          }\n          ... on Ente {\n            razon_social\n          }\n        }\n      }\n    }\n  }\n": typeof types.VillasPageDocument,
    "\n  query ObtenerPerodosDisponibles {\n    periodos: obtenerPeriodosDisponibles {\n      anio\n      mes\n    }\n  }\n": typeof types.ObtenerPerodosDisponiblesDocument,
    "\n  mutation RegistrarCuota($input: RegistrarCuotaDTO!) {\n    registrarCuota(input: $input) {\n      __typename\n      ... on CuotaRegular {\n        id\n      }\n      ... on CuotaEspecial {\n        id\n      }\n    }\n  }\n": typeof types.RegistrarCuotaDocument,
    "\n  mutation RegistrarGastoOverlay($input: RegistrarGastoDTO!) {\n    registrarGasto(input: $input) {\n      id\n      concepto\n    }\n  }\n": typeof types.RegistrarGastoOverlayDocument,
    "\n  query BuscarGastosHuerfanos($busqueda: String) {\n    gastos: obtenerGastos(\n      filter: { concepto: { like: $busqueda }, cuota: { eq: null } }\n    ) {\n      data {\n        __typename\n        ... on Gasto {\n          monto\n          total\n          operacion\n          concepto\n          metodo\n          moneda\n          tasa\n          fecha\n          metodo\n          registrado_por\n          registro\n        }\n        ... on GastoAProveedor {\n          proveedor {\n            id\n            nombre\n            rif\n            telefono\n            email\n            actualizado_en\n            creado_en\n            direccion\n          }\n        }\n      }\n    }\n  }\n": typeof types.BuscarGastosHuerfanosDocument,
};
const documents: Documents = {
    "\n  query CuotaPage($cuota_id: String!) {\n    cuota: obtenerCuota(id: $cuota_id) {\n      __typename\n      ... on Cuota {\n        id\n        mes\n        anio\n        monto\n        gastos {\n          __typename\n          ... on Gasto {\n            operacion\n            concepto\n            moneda\n            monto\n            fecha\n            tasa\n            total\n          }\n          ... on GastoAProveedor {\n            proveedor {\n              id\n              nombre\n              rif\n              telefono\n              email\n            }\n          }\n        }\n        recaudacion {\n          unidades_aplicadas\n          unidades_solventes\n          monto_estimado\n          monto_recaudado\n        }\n      }\n      ... on CuotaEspecial {\n        detalles {\n          titulo\n          descripcion\n          justificacion\n        }\n      }\n    }\n    deudas: obtenerDeudas(filtro: { cuota: { eq: $cuota_id } }) {\n      data {\n        unidad {\n          codigo\n        }\n        titular {\n          display_name\n        }\n        deuda\n        estado\n      }\n    }\n  }\n": types.CuotaPageDocument,
    "\n  query CuotasPage {\n    cuotas: obtenerCuotas {\n      data {\n        __typename\n        ... on Cuota {\n          id\n          monto\n          mes\n          anio\n          registro\n          recaudacion {\n            unidades_aplicadas\n            pagos_asociados\n          }\n        }\n\n        ... on CuotaEspecial {\n          detalles {\n            titulo\n            descripcion\n          }\n        }\n      }\n    }\n  }\n": types.CuotasPageDocument,
    "\n  query RegistrarCuotaPage {\n    obtenerProveedores {\n      id\n      nombre\n    }\n  }\n": types.RegistrarCuotaPageDocument,
    "\n  query DashboardPage {\n    proveedores: obtenerProveedores {\n      id\n      nombre\n    }\n  }\n": types.DashboardPageDocument,
    "\n  mutation Login($email: String!, $pass: String!) {\n    login(email: $email, password: $pass) {\n      token\n    }\n  }\n": types.LoginDocument,
    "\n  query RegistrarPagoPage($codigo_like: String!) {\n    unidades: obtenerUnidades(\n      filter: { codigo: { like: $codigo_like } }\n      paginator: { limit: 5, page: 1 }\n    ) {\n      data {\n        id\n        codigo\n      }\n    }\n  }\n": types.RegistrarPagoPageDocument,
    "\n  query VillaPage($codigo: String!) {\n    unidad: obtenerUnidadPorCodigo(codigo: $codigo) {\n      codigo\n      estado\n      wallet\n      titular_primario {\n        __typename\n        ... on Sujeto {\n          id\n          cedula\n          display_name\n        }\n      }\n      titulares {\n        __typename\n\n        ... on Sujeto {\n          id\n          cedula\n          display_name\n          email\n        }\n      }\n    }\n    pagos: obtenerPagos(filtro: { unidad: { eq: $codigo } }) {\n      data {\n        operacion\n        concepto\n        monto\n        moneda\n      }\n    }\n    deudas: obtenerDeudas(filtro: { unidad: { eq: $codigo } }) {\n      data {\n        id\n        cuota {\n          __typename\n          ... on Deuda__Cuota {\n            id\n            nombre\n          }\n        }\n        deuda\n        monto\n        estado\n      }\n    }\n  }\n": types.VillaPageDocument,
    "\n  query VillasPage {\n    estadisticas: obtenerUnidadesEstadisticas {\n      total_unidades\n      unidades_activas\n      unidades_con_pendientes\n      unidades_inhabitadas\n    }\n\n    villas: obtenerUnidades(filter: { estado: { eq: \"ACTIVA\" } }) {\n      data {\n        codigo\n        estado\n        wallet\n        contacto {\n          id\n          email\n          telefono\n        }\n        titular_primario {\n          __typename\n          ... on Sujeto {\n            id\n            cedula\n            display_name\n          }\n          ... on Persona {\n            nombres\n            apellidos\n          }\n          ... on Ente {\n            razon_social\n          }\n        }\n      }\n    }\n  }\n": types.VillasPageDocument,
    "\n  query ObtenerPerodosDisponibles {\n    periodos: obtenerPeriodosDisponibles {\n      anio\n      mes\n    }\n  }\n": types.ObtenerPerodosDisponiblesDocument,
    "\n  mutation RegistrarCuota($input: RegistrarCuotaDTO!) {\n    registrarCuota(input: $input) {\n      __typename\n      ... on CuotaRegular {\n        id\n      }\n      ... on CuotaEspecial {\n        id\n      }\n    }\n  }\n": types.RegistrarCuotaDocument,
    "\n  mutation RegistrarGastoOverlay($input: RegistrarGastoDTO!) {\n    registrarGasto(input: $input) {\n      id\n      concepto\n    }\n  }\n": types.RegistrarGastoOverlayDocument,
    "\n  query BuscarGastosHuerfanos($busqueda: String) {\n    gastos: obtenerGastos(\n      filter: { concepto: { like: $busqueda }, cuota: { eq: null } }\n    ) {\n      data {\n        __typename\n        ... on Gasto {\n          monto\n          total\n          operacion\n          concepto\n          metodo\n          moneda\n          tasa\n          fecha\n          metodo\n          registrado_por\n          registro\n        }\n        ... on GastoAProveedor {\n          proveedor {\n            id\n            nombre\n            rif\n            telefono\n            email\n            actualizado_en\n            creado_en\n            direccion\n          }\n        }\n      }\n    }\n  }\n": types.BuscarGastosHuerfanosDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query CuotaPage($cuota_id: String!) {\n    cuota: obtenerCuota(id: $cuota_id) {\n      __typename\n      ... on Cuota {\n        id\n        mes\n        anio\n        monto\n        gastos {\n          __typename\n          ... on Gasto {\n            operacion\n            concepto\n            moneda\n            monto\n            fecha\n            tasa\n            total\n          }\n          ... on GastoAProveedor {\n            proveedor {\n              id\n              nombre\n              rif\n              telefono\n              email\n            }\n          }\n        }\n        recaudacion {\n          unidades_aplicadas\n          unidades_solventes\n          monto_estimado\n          monto_recaudado\n        }\n      }\n      ... on CuotaEspecial {\n        detalles {\n          titulo\n          descripcion\n          justificacion\n        }\n      }\n    }\n    deudas: obtenerDeudas(filtro: { cuota: { eq: $cuota_id } }) {\n      data {\n        unidad {\n          codigo\n        }\n        titular {\n          display_name\n        }\n        deuda\n        estado\n      }\n    }\n  }\n"): typeof import('./graphql').CuotaPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query CuotasPage {\n    cuotas: obtenerCuotas {\n      data {\n        __typename\n        ... on Cuota {\n          id\n          monto\n          mes\n          anio\n          registro\n          recaudacion {\n            unidades_aplicadas\n            pagos_asociados\n          }\n        }\n\n        ... on CuotaEspecial {\n          detalles {\n            titulo\n            descripcion\n          }\n        }\n      }\n    }\n  }\n"): typeof import('./graphql').CuotasPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query RegistrarCuotaPage {\n    obtenerProveedores {\n      id\n      nombre\n    }\n  }\n"): typeof import('./graphql').RegistrarCuotaPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query DashboardPage {\n    proveedores: obtenerProveedores {\n      id\n      nombre\n    }\n  }\n"): typeof import('./graphql').DashboardPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation Login($email: String!, $pass: String!) {\n    login(email: $email, password: $pass) {\n      token\n    }\n  }\n"): typeof import('./graphql').LoginDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query RegistrarPagoPage($codigo_like: String!) {\n    unidades: obtenerUnidades(\n      filter: { codigo: { like: $codigo_like } }\n      paginator: { limit: 5, page: 1 }\n    ) {\n      data {\n        id\n        codigo\n      }\n    }\n  }\n"): typeof import('./graphql').RegistrarPagoPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query VillaPage($codigo: String!) {\n    unidad: obtenerUnidadPorCodigo(codigo: $codigo) {\n      codigo\n      estado\n      wallet\n      titular_primario {\n        __typename\n        ... on Sujeto {\n          id\n          cedula\n          display_name\n        }\n      }\n      titulares {\n        __typename\n\n        ... on Sujeto {\n          id\n          cedula\n          display_name\n          email\n        }\n      }\n    }\n    pagos: obtenerPagos(filtro: { unidad: { eq: $codigo } }) {\n      data {\n        operacion\n        concepto\n        monto\n        moneda\n      }\n    }\n    deudas: obtenerDeudas(filtro: { unidad: { eq: $codigo } }) {\n      data {\n        id\n        cuota {\n          __typename\n          ... on Deuda__Cuota {\n            id\n            nombre\n          }\n        }\n        deuda\n        monto\n        estado\n      }\n    }\n  }\n"): typeof import('./graphql').VillaPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query VillasPage {\n    estadisticas: obtenerUnidadesEstadisticas {\n      total_unidades\n      unidades_activas\n      unidades_con_pendientes\n      unidades_inhabitadas\n    }\n\n    villas: obtenerUnidades(filter: { estado: { eq: \"ACTIVA\" } }) {\n      data {\n        codigo\n        estado\n        wallet\n        contacto {\n          id\n          email\n          telefono\n        }\n        titular_primario {\n          __typename\n          ... on Sujeto {\n            id\n            cedula\n            display_name\n          }\n          ... on Persona {\n            nombres\n            apellidos\n          }\n          ... on Ente {\n            razon_social\n          }\n        }\n      }\n    }\n  }\n"): typeof import('./graphql').VillasPageDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query ObtenerPerodosDisponibles {\n    periodos: obtenerPeriodosDisponibles {\n      anio\n      mes\n    }\n  }\n"): typeof import('./graphql').ObtenerPerodosDisponiblesDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation RegistrarCuota($input: RegistrarCuotaDTO!) {\n    registrarCuota(input: $input) {\n      __typename\n      ... on CuotaRegular {\n        id\n      }\n      ... on CuotaEspecial {\n        id\n      }\n    }\n  }\n"): typeof import('./graphql').RegistrarCuotaDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation RegistrarGastoOverlay($input: RegistrarGastoDTO!) {\n    registrarGasto(input: $input) {\n      id\n      concepto\n    }\n  }\n"): typeof import('./graphql').RegistrarGastoOverlayDocument;
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query BuscarGastosHuerfanos($busqueda: String) {\n    gastos: obtenerGastos(\n      filter: { concepto: { like: $busqueda }, cuota: { eq: null } }\n    ) {\n      data {\n        __typename\n        ... on Gasto {\n          monto\n          total\n          operacion\n          concepto\n          metodo\n          moneda\n          tasa\n          fecha\n          metodo\n          registrado_por\n          registro\n        }\n        ... on GastoAProveedor {\n          proveedor {\n            id\n            nombre\n            rif\n            telefono\n            email\n            actualizado_en\n            creado_en\n            direccion\n          }\n        }\n      }\n    }\n  }\n"): typeof import('./graphql').BuscarGastosHuerfanosDocument;


export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}
