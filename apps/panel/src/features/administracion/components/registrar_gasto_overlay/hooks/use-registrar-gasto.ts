import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { Moneda as MonedaGraphql } from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";
import { Moneda } from "@/features/administracion/schemas/moneda.schema";
import type { NuevoGasto, NuevoGastoYProveedor } from "@/features/administracion/schemas/gasto.schema";

const RegistrarGastoMutation = graphql(`
  mutation RegistrarGasto($input: RegistrarGastoDTO!) {
    registrarGasto(input: $input) {
      id
    }
  }
`);

const RegistrarGastoYProveedorMutation = graphql(`
  mutation RegistrarGastoYProveedor($input: RegistrarGastoYProveedorDTO!) {
    registrarGastoYProveedor(input: $input) {
      id
      concepto
    }
  }
`);

const toGraphqlMoneda = (moneda: Moneda) =>
  moneda === Moneda.USD ? MonedaGraphql.Usd : MonedaGraphql.Ved;

export function useRegistrarGasto() {
  const registrarGasto = useMutation({
    mutationFn: (data: NuevoGasto) => {
      return execute(RegistrarGastoMutation, {
        input: {
          concepto: data.concepto,
          proveedor: data.proveedor,
          monto: data.monto,
          fecha: data.fecha,
          moneda: toGraphqlMoneda(data.moneda),
        },
      });
    },
  });

  const registrarGastoYProveedor = useMutation({
    mutationFn: (data: NuevoGastoYProveedor) => {
      return execute(RegistrarGastoYProveedorMutation, {
        input: {
          concepto: data.concepto,
          proveedor: {
            nombre: data.proveedor.nombre,
            rif: data.proveedor.rif,
            telefono: data.proveedor.telefono,
            email: data.proveedor.email || "",
          },
          monto: data.monto,
          fecha: data.fecha,
          moneda: toGraphqlMoneda(data.moneda),
        },
      });
    },
  });

  return {
    registrarGasto,
    registrarGastoYProveedor,
  };
}
