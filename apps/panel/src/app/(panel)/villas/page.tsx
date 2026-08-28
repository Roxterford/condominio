import { Button } from "@/components/ui/button";
import {
  VillasTable,
  VillasTableData,
} from "@/features/villas/components/villas_table";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import { AlertCircle, CheckCircle, House, Plus, UserPlus } from "lucide-react";
import Link from "next/link";
import styles from "./page.module.css";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

const PageQuery = graphql(/* GraphQL */ `
  query VillasPage {
    estadisticas: obtenerUnidadesEstadisticas {
      total_unidades
      unidades_activas
      unidades_con_pendientes
      unidades_inhabitadas
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
    if (!data.estadisticas || !data.villas) {
      return <div>Datos incompletos</div>;
    }

    const { estadisticas, villas } = data;

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
          <div className="flex gap-2">
            <Button asChild variant="outline">
              <Link href="/sujetos/registrar">
                <UserPlus /> Registrar propietario
              </Link>
            </Button>
            <Button asChild>
              <Link href="/unidades/registrar">
                <Plus /> Registrar unidad
              </Link>
            </Button>
          </div>
        </header>
        <ul className={[styles.infoboxes, "mt-10"].join(" ")}>
          <li className={styles.infobox}>
            <div
              className={[styles["infobox__icon-content"], "bg-blue-50"].join(
                " ",
              )}
            >
              <House />
            </div>
            <div>
              <h3 className={styles.infobox__title}>Total Villas</h3>
              <p className={styles.infobox__value}>
                {estadisticas?.total_unidades}
              </p>
            </div>
          </li>
          <div className={styles.infoboxes__divider}></div>
          <li className={styles.infobox}>
            <div
              className={[styles["infobox__icon-content"], "bg-green-100"].join(
                " ",
              )}
            >
              <CheckCircle className="text-green-700" />
            </div>
            <div>
              <h3 className={styles.infobox__title}>Villas Activas</h3>
              <p className={styles.infobox__value}>
                {estadisticas?.unidades_activas}
              </p>
            </div>
          </li>
          <div className={styles.infoboxes__divider}></div>
          <li className={styles.infobox}>
            <div
              className={[
                styles["infobox__icon-content"],
                "bg-yellow-100",
              ].join(" ")}
            >
              <AlertCircle className="text-yellow-700" />
            </div>
            <div>
              <h3 className={styles.infobox__title}>
                Villas con Pagos pendientes
              </h3>
              <p className={styles.infobox__value}>
                {estadisticas?.unidades_con_pendientes}
              </p>
            </div>
          </li>
          <div className={styles.infoboxes__divider}></div>
          <li className={styles.infobox}>
            <div
              className={[styles["infobox__icon-content"], "bg-red-100"].join(
                " ",
              )}
            >
              <AlertCircle className="text-red-700" />
            </div>
            <div>
              <h3 className={styles.infobox__title}>Villas Inhabitadas</h3>
              <p className={styles.infobox__value}>
                {estadisticas?.unidades_inhabitadas}
              </p>
            </div>
          </li>
        </ul>
        <section>
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
