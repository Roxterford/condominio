import { VillasTable } from "@/features/villas/components/villas_table";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { AlertCircle, CheckCircle, House } from "lucide-react";
import styles from "./page.module.css";

const PageQuery = graphql(`
  query VillasPage {
    villas: obtenerVillasTotales {
      total_villas
      villas_activas
      villas_con_pendientes
      villas_inhabitadas
    }
  }
`);

export default async function VillasPage() {
  const {
    data: { villas },
    errors,
  } = await execute(PageQuery);

  if (errors) return <pre>{JSON.stringify(errors, null, 4)}</pre>;

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
            <p className={styles.infobox__value}>{villas?.total_villas}</p>
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
            <p className={styles.infobox__value}>{villas?.villas_activas}</p>
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
              {villas?.villas_con_pendientes}
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
              {villas?.villas_inhabitadas}
            </p>
          </div>
        </li>
      </ul>
      <section>
        <VillasTable />
      </section>
    </>
  );
}
