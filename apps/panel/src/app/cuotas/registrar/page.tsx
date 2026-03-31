import { RegistrarCuotaForm } from "@/features/administracion/components/registrar_cuota_form";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";

interface ProveedorFromApi {
  __typename?: "Proveedor";
  id: string;
  nombre?: string;
}

const ProveedoresQuery = graphql(`
  query Proveedores {
    obtenerProveedores {
      id
      nombre
    }
  }
`);

export default async function RegistrarCuotaPage() {
  const {
    data: { obtenerProveedores },
  } = await execute(ProveedoresQuery);

  console.log({ obtenerProveedores });
  const proveedores = obtenerProveedores.map((p) => ({
    id: p.id,
    nombre: p.nombre || "",
  }));

  return <RegistrarCuotaForm proveedores={proveedores} />;
}
