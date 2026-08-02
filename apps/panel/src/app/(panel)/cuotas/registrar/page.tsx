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
  }
`);

export default async function RegistrarCuotaPage() {
  const result = await execute(PageQuery);

  if (result.errors || !result.data) {
    console.error(result)
    throw new Error("Error al cargar datos para registrar cuota");
  }

  const { obtenerProveedores } = result.data;

  const proveedores = obtenerProveedores.map((p) => ({
    id: p.id,
    nombre: p.nombre || "",
  }));

  const gastos: RegistrarCuotaFormProps["gastos"] = []
  return <RegistrarCuotaForm proveedores={proveedores} gastos={gastos} />;
}
