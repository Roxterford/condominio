import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { useQuery } from "@tanstack/react-query";

const PageQuery = graphql(/* GraphQL */ `
  query VillaPage($codigo: String!) {
    villa: obtenerUnidadPorCodigo(codigo: $codigo) {
      id
      titular_primario {
        __typename

        ... on Sujeto {
          id
          telefono
          email
        }
        ... on Ente {
          razon_social
          representante {
            nombres
            apellidos
          }
        }
        ... on Persona {
          nombres
          apellidos
        }
      }
      estado
      contacto {
        id
        nombres
      }
    }
  }
`);

export interface VillaPageProps {
  params: Promise<{
    codigo: string;
  }>;
}

export default async function VillaPage(props: VillaPageProps) {
  const params = await props.params;

  const {
    errors,
    data: { villa },
  } = await execute(PageQuery, {
    codigo: params.codigo,
  });

  if (errors) return <>ERROR</>;
  if (!villa) return <>NOT FOUND</>;

  let nombre = "Villa " + params.codigo;

  const { titular_primario } = villa;

  switch (titular_primario?.__typename) {
    case "Ente":
      nombre = titular_primario.razon_social;
      break;
    case "Persona":
      nombre = [
        titular_primario.nombres.split(" ").at(0),
        titular_primario.apellidos.split(" ").at(0),
      ].join(" ");
      break;
  }

  return (
    <>
      <h1>{nombre}</h1>
      <pre>{JSON.stringify(villa, null, 4)}</pre>
    </>
  );
}
