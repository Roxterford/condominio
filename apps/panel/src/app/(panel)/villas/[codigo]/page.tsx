import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { useQuery } from "@tanstack/react-query";
import {
  CardSim,
  CreditCard,
  CreditCardIcon,
  Home,
  Mail,
  Phone,
} from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { EstadoDeDeuda } from "@/providers/graphql/graphql";
import Link from "next/link";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";

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

    deudas: obtenerDeudasDeUnaUnidadPorCodigo(codigo: $codigo) {
      total
      pages
      data {
        id
        estado
        cuota
        monto
        deuda
      }
    }

    pagos: obtenerPagos(filter: { unidad: { eq: $codigo } }) {
      total
      data {
        id
        unidad
        disponible
        unidad
        monto
        moneda
        tasa
        fecha
        total
        destinado
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
    data: { villa, deudas, pagos },
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
        <div className="mt-2 flex-1">
          <section className="flex items-center justify-between gap-5">
            <section>
              <div className="flex items-center gap-3">
                <h1>{nombre}</h1>
                <ul className="flex gap-2">
                  <li>
                    <Badge>Propietario</Badge>
                  </li>
                  <li>
                    {!villa.deuda ? (
                      <Badge className="text-green-600 bg-green-100">
                        Al día
                      </Badge>
                    ) : (
                      <Badge className="text-yellow-600 bg-yellow-100">
                        En Deuda
                      </Badge>
                    )}
                  </li>
                </ul>
              </div>
              <p className="flex items-center gap-3 text-gray-500 mt-1">
                <Home size={16} /> Villa {villa.codigo}
              </p>
            </section>

            <ul aria-label="Acciones">
              <li>
                <Button>
                  {" "}
                  <CreditCardIcon /> Registrar pago
                </Button>
              </li>
            </ul>
          </section>

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
                  <a
                    className="link"
                    href={`mailto:${villa.titular_primario.email}`}
                  >
                    {villa.titular_primario.email}
                  </a>
                </div>
              </li>
            </ul>
          )}
        </div>
      </header>
      <section>
        <ul aria-label="Información relevante">
          <li>
            <h4 className="text-sm text-gray-500 font-semibold flex items-center gap-1">
              Estado de cuenta
              <span
                title="El monto en rojo indica deuda pendiente; en verde indica saldo a favor"
                className="cursor-help text-gray-400 hover:text-gray-600"
                aria-label="El monto en rojo indica deuda pendiente; en verde indica saldo a favor"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  aria-hidden="true"
                >
                  <circle cx="12" cy="12" r="10" />
                  <line x1="12" y1="16" x2="12" y2="12" />
                  <line x1="12" y1="8" x2="12.01" y2="8" />
                </svg>
              </span>
            </h4>
            <p
              className={`font-medium ${villa.deuda > 0 ? "text-red-600" : "text-green-600"}`}
            >
              {villa.deuda > 0 ? "-" : "+"}${villa.deuda}
            </p>
          </li>
        </ul>
      </section>
      {/* <pre>{JSON.stringify(villa, null, 4)}</pre> */}

      <Tabs defaultValue="pagos">
        <TabsList variant="line">
          <TabsTrigger value="pagos">
            Pagos
            {pagos.total ? <Badge>{pagos.total}</Badge> : null}
          </TabsTrigger>
          <TabsTrigger value="deudas">
            Deudas
            {deudas.total ? (
              <Badge>{deudas.total.toLocaleString("es-VE")}</Badge>
            ) : null}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="pagos">
          <section>
            <h3>Historial de pagos</h3>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Monto</th>
                  <th>Fecha</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {pagos.data.map((pago) => (
                  <tr key={pago.id}>
                    <td className="p-10">
                      <button className="link">{pago.id}</button>
                    </td>
                    <td>${pago.total}</td>

                    <td className="p-10">{pago.fecha}</td>
                    <td className="p-10">
                      <Button variant="outline">Ver detalles</Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>
        </TabsContent>

        <TabsContent value="deudas">
          <section>
            <h3>Deuda Total: ${villa.deuda}</h3>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Cuota</th>
                  <th>Estado</th>
                  <th>Deuda</th>
                </tr>
              </thead>
              <tbody>
                {deudas.data.map((deuda) => (
                  <tr key={deuda.id}>
                    <td className="p-10">{deuda.id}</td>
                    <td className="p-10">
                      <Link href={"/cuotas/" + deuda.cuota} className="link">
                        {deuda.cuota}
                      </Link>
                    </td>
                    <td className="p-10">
                      {(() => {
                        switch (deuda.estado) {
                          case EstadoDeDeuda.Pendiente:
                            return (
                              <Badge className="bg-yellow-100 text-yellow-600">
                                Pendiente
                              </Badge>
                            );
                          case EstadoDeDeuda.Abonada:
                            return (
                              <Badge className="bg-blue-100 text-blue-600">
                                Abonada
                              </Badge>
                            );
                          case EstadoDeDeuda.Saldada:
                            return (
                              <Badge className="bg-green-100 text-green-600">
                                Saldada
                              </Badge>
                            );
                          default:
                            break;
                        }
                      })()}
                    </td>
                    <td className="p-10">
                      {deuda.deuda <= 0 ? "-" : "$" + deuda.deuda}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>
        </TabsContent>
      </Tabs>
    </>
  );
}
