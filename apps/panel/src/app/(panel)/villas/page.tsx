import { VillasTable, VillasTableData } from "@/features/villas/components/villas_table";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { renderGraphql } from "@/providers/graphql/render";
import { AlertCircle, CheckCircle, House } from "lucide-react";
import styles from "./page.module.css";

const PageQuery = graphql(/* GraphQL */`
  query VillasPage {
    estadisticas: obtenerUnidadesEstadisticas {
      total_unidades
      unidades_activas
      unidades_con_pendientes
      unidades_inhabitadas
    }


    villas: obtenerUnidades(filter: { estado: { eq:"ACTIVA" } }) {
      data {
        codigo
        contacto {
          id
          email
          telefono
        }
        titular_primario {
          __typename
          ... on Sujeto {
            id
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

    const villas_table_data: VillasTableData[] = villas.data.map<VillasTableData>(
    (villa) => 
  {    


    const titular_primario = villa.titular_primario

    let nombre = ""
    
    switch (titular_primario?.__typename) {
      case "Persona": 
        nombre = titular_primario.nombres.split(" ").at(0) + " " + titular_primario.apellidos.split(" ").at(0)
        break;
      
      case  "Ente":
        nombre = titular_primario.razon_social
        break;
    }
      
      return {
      codigo: villa.codigo,
      propietario: {
        nombre
      },
      contacto: {
        email: villa.contacto?.email ?? "",
        telefono: villa.contacto?.telefono ?? ""
      },
      estado_pagos : "pendiente" // TODO: cambiar

    }}
  
  )

  return (
    <>
      <header>
        <h1>Villas</h1>
        <p className="page-description">
          Gestione la información de propietarios del condominio
        </p>
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
            <p className={styles.infobox__value}>{estadisticas?.total_unidades}</p>
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
            <p className={styles.infobox__value}>{estadisticas?.unidades_activas}</p>
          </div>
        </li>
        <div className={styles.infoboxes__divider}></div>
        <li className={styles.infobox}>
          <div
            className={[styles["infobox__icon-content"], "bg-yellow-100"].join(
              " ",
            )}
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
        <VillasTable data={villas_table_data} />
      </section>
    </>
  );
  });
}
