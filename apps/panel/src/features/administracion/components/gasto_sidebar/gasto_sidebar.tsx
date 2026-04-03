import { OverlayProps } from "@/components/overlay";
import { Badge } from "@/components/ui/badge";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Gasto } from "../../schemas";
import styles from "./gasto_sidebar.module.css";

export interface GastoSidebar extends OverlayProps {
  gasto?: Gasto;
}

export function GastoSidebar({ open, onOpenChange, gasto }: GastoSidebar) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent>
        <SheetHeader className="px-8">
          <div className="flex gap-2 items-center">
            <SheetTitle className="font-semibold text-lg">
              Información de Pago
            </SheetTitle>
            <Badge variant="secondary">{gasto?.id || "nil"}</Badge>
          </div>
        </SheetHeader>
        <div className="overflow-auto bg-red-50 ">
          <section className={styles.section}>
            <div className={styles.infoboxes}>
              <div className={styles.infobox}>
                <span className={styles.infobox__title}>Total</span>
                <p className={`${styles.infobox__value} text-lg font-bold`}>
                  $ {gasto?.total || "0,00"}
                </p>
              </div>
              <div className={styles.infoboxes__separator}></div>
              <div className={styles.infobox}>
                <span className={styles.infobox__title}>Tasa</span>
                <p className={styles.infobox__value}>Bs. {gasto?.tasa}</p>
              </div>
              <div className={styles.infoboxes__separator}></div>
              <div className={styles.infobox}>
                <span className={styles.infobox__title}>Fecha</span>
                <p className={styles.infobox__value}>
                  {gasto?.fecha.toLocaleString("es")}
                </p>
              </div>
            </div>
          </section>
          <section className={styles.section}>
            <div className={styles.infobox}>
              <span className={styles.infobox__title}>Concepto</span>
              <p className={styles.infobox__value}>{gasto?.concepto}</p>
            </div>
            <div className={styles.infobox}>
              <span className={styles.infobox__title}>Monto</span>
              <p className={styles.infobox__value}>{gasto?.monto} </p>
            </div>
          </section>
          <section className={styles.section}>
            <div className={styles.infoboxes}>
              <div className={styles.infobox} data-type="sm">
                <span className={styles.infobox__title}>ID</span>
                <p className={styles.infobox__value}>{gasto?.id}</p>
              </div>
              <div className={styles.infobox} data-type="sm">
                <span className={styles.infobox__title}>Registro</span>
                <p className={styles.infobox__value}>{gasto?.id}</p>
              </div>
            </div>
          </section>
          <section className={styles.section}>
            <p className="font-medium">Proveedor</p>
            <div className="grid place-items-center select-none rounded-lg size-15 font-medium text-3xl bg-indigo-100 text-indigo-500">
              E
            </div>
            <table>
              <tbody>
                <tr>
                  <th>Nombre:</th>
                  <td>{gasto?.proveedor}</td>
                </tr>
                <tr>
                  <th>CI / RIF:</th>
                  <td>{gasto?.proveedor}</td>
                </tr>
                <tr>
                  <th>Teléfono</th>
                  <td>{gasto?.proveedor}</td>
                </tr>
                <tr>
                  <th>Correo:</th>
                  <td>{gasto?.proveedor}</td>
                </tr>
              </tbody>
            </table>
          </section>
        </div>
      </SheetContent>
    </Sheet>
  );
}
