import { AlertCircle, CheckCircle, House } from "lucide-react";
import styles from "./page.module.css";

export default function VillasPage() {
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
            <p className={styles.infobox__value}>460</p>
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
            <p className={styles.infobox__value}>460</p>
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
            <p className={styles.infobox__value}>460</p>
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
            <p className={styles.infobox__value}>460</p>
          </div>
        </li>
      </ul>
    </>
  );
}
