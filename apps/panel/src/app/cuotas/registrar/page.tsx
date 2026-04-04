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
      moneda: g.moneda || "USD",
      monto: g.monto,
      fecha: new Date(g.fecha),
      proveedor: {
        id: g.proveedor.id,
        nombre: g.proveedor.nombre,
        rif: g.proveedor.rif,
        telefono: g.proveedor.telefono || "",
        email: g.proveedor.email || "",
      },
    }),
  );

  return <RegistrarCuotaForm proveedores={proveedores} gastos={gastos} />;
}
