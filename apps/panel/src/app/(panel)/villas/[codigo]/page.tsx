import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { useQuery } from "@tanstack/react-query";
import { Home, Mail, Phone } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

const PageQuery = graphql(/* GraphQL */ `
  query VillaPage($codigo: String!) {
    villa: obtenerUnidadPorCodigo(codigo: $codigo) {
      id
      codigo
      deuda
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
            id
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

  if (errors)
    return (
      <>
        <h1>ERROR</h1>
        <pre>{JSON.stringify(errors, null, 4)}</pre>
      </>
    );
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
      <header className="flex gap-5">
        <Avatar className="size-20">
          <AvatarImage alt="propietario/titular" />
          <AvatarFallback>CN</AvatarFallback>
        </Avatar>
        <div className="mt-2">
          <div className="flex gap-5 items-center">
            <h1>{nombre}</h1>
            <div className="flex gap-3">
              <Badge>Propietario</Badge>
              {!villa.deuda ? (
                <Badge className="text-green-600 bg-green-100">Al día</Badge>
              ) : (
                <Badge className="text-yellow-600 bg-yellow-100">
                  En Deuda
                </Badge>
              )}
            </div>
          </div>
          <p className="flex items-center gap-3 text-gray-500 mt-1">
            <Home size={16} /> Villa {villa.codigo}
          </p>

          {villa.titular_primario && (
            <ul className="flex gap-10 mt-5" aria-label="Contacto del titular">
              <li className="flex gap-3">
                <Phone size={18} className="text-gray-500" />
                <div>
                  <h4>Teléfono</h4>
                  <p>{villa.titular_primario.telefono}</p>
                </div>
              </li>
              <li className="flex gap-3">
                <Mail size={18} className="text-gray-500" />
                <div>
                  <h4>Email</h4>
                  <p>{villa.titular_primario.email}</p>
                </div>
              </li>
            </ul>
          )}
        </div>
      </header>
      <pre>{JSON.stringify(villa, null, 4)}</pre>
    </>
  );
}
