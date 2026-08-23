import { RegistrarCuotaForm } from "@/features/administracion/components/registrar-cuota-form";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";

const PageQuery = graphql(`
  query RegistrarCuotaPage {
    obtenerProveedores {
      id
      nombre
    }
  }
`);

export default async function RegistrarCuotaPage() {
  return renderGraphql(await execute(PageQuery), ({ obtenerProveedores }) => {
    const proveedores = obtenerProveedores.map((p) => ({
      id: p.id,
      nombre: p.nombre || "",
    }));

    return <RegistrarCuotaForm proveedores={proveedores} />;
  });
}
