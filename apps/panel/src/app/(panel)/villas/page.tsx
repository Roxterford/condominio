import {
  VillasTable,
  VillasTableData,
} from "@/features/villas/components/villas_table";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import StatCard from "@/components/ui/StatCard";
import { money } from "@/lib/money-display";
import { Box, CircleAlert, CircleCheckBig, DollarSign } from "lucide-react";

const PageQuery = graphql(/* GraphQL */ `
  query VillasPage {
    resumen: obtenerResumenUnidades {
      total_unidades
      unidades_solventes
      unidades_con_pendientes
      total_pendiente
    }

    villas: obtenerUnidades(filter: { estado: { eq: "ACTIVA" } }) {
      data {
        codigo
        estado
        wallet
        contacto {
          id
          email
          telefono
        }
        titular_primario {
          __typename
          ... on Sujeto {
            id
            cedula
            display_name
          }
          ... on Persona {
            nombres
            apellidos
          }
          ... on Ente {
            razon_social
          }
        }
      }
    }
  }
`);

export default async function VillasPage() {
  return renderGraphql(await execute(PageQuery), (data) => {
    if (!data.resumen || !data.villas) {
      return <div>Datos incompletos</div>;
    }

    const { resumen, villas } = data;

    const villas_table_data: VillasTableData[] =
      villas.data.map<VillasTableData>((villa) => {
        const titular_primario = villa.titular_primario;

        return {
          codigo: villa.codigo,
          estado: villa.estado,
          wallet: villa.wallet,
          propietario: titular_primario ?? undefined,
          contacto: {
            email: villa.contacto?.email ?? "",
            telefono: villa.contacto?.telefono ?? "",
          },
          estado_pagos: "pendiente", // TODO: cambiar
        };
      });

    return (
      <>
        <header className="flex items-end justify-between">
          <div>
            <h1>Villas</h1>
            <p className="page-description">
              Gestione la información de propietarios del condominio
            </p>
          </div>
          <div className="flex gap-2"></div>
        </header>
        <ul className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6 mt-8">
          <StatCard
            color="bg-primary/10"
            icon={<Box className="text-primary" />}
            title="Unidades totales"
            value={resumen.total_unidades}
            subtitle=""
          ></StatCard>
          <StatCard
            color="bg-green-100"
            icon={<CircleCheckBig className="text-green-700" />}
            title="Solventes"
            value={resumen.unidades_solventes}
            subtitle=""
          ></StatCard>
          <StatCard
            color="bg-yellow-100"
            icon={<CircleAlert className="text-yellow-600" />}
            title="Pendientes de pago"
            value={resumen.unidades_con_pendientes}
            subtitle=""
          ></StatCard>
          <StatCard
            color="bg-rose-100"
            icon={<DollarSign className="text-rose-700" />}
            title="Deuda pendiente"
            value={money(resumen.total_pendiente)}
            subtitle=""
          ></StatCard>
        </ul>
        <section className="mt-5">
          <Tabs defaultValue="todas">
            <TabsList variant="line">
              <TabsTrigger value="todas">Todas</TabsTrigger>
              <TabsTrigger value="activas">Activas</TabsTrigger>
              <TabsTrigger value="deuda">Con deuda pendiente</TabsTrigger>
              <TabsTrigger value="inhabitadas">Inhabitadas</TabsTrigger>
            </TabsList>
            <TabsContent value="todas">
              <VillasTable data={villas_table_data} />
            </TabsContent>
          </Tabs>
        </section>
      </>
    );
  });
}
