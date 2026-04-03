import {
  RegistrarCuotaForm,
  RegistrarCuotaFormProps,
} from "@/features/administracion/components/registrar_cuota_form";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";

const PageQuery = graphql(`
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
        proveedor
      }
    }
  }
`);

export default async function RegistrarCuotaPage() {
  const {
    data: { obtenerProveedores, obtenerGastos },
  } = await execute(PageQuery);

  const proveedores = obtenerProveedores.map((p) => ({
    id: p.id,
    nombre: p.nombre || "",
  }));

  const gastos: RegistrarCuotaFormProps["gastos"] = obtenerGastos.data.map(
    (g) => ({
      id: g.id,
      concepto: g.concepto,
      moneda: g.moneda,
      monto: g.monto,
      fecha: new Date(g.fecha),
      proveedor: g.proveedor,
    }),
  );

  return <RegistrarCuotaForm proveedores={proveedores} gastos={gastos} />;
}
