import {
  ChartPie,
  House,
  LayoutGrid,
  Newspaper,
  Plus,
  Settings,
} from "lucide-react";
import Link from "next/link";
import styles from "./layout.module.css";
import clsx from "clsx";

export default function PanelLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex">
      <Navbar />
      <main className="w-full py-5 px-10">{children}</main>
    </div>
  );
}

function Navbar() {
  return (
    <nav className="border-r pt-10 px-5">
      <ul className="grid gap-5 place-items-center">
        <li>
          <Link href="/pagos/registrar" className={styles.navlink}>
            <div className={clsx("bg-primary")}>
              <Plus className={clsx(styles.navlink__icon, "text-white")} />
            </div>
            Nuevo Pago
          </Link>
        </li>
        <li>
          <Link href="/" className={styles.navlink}>
            <div className={styles["navlink__icon-content"]}>
              <LayoutGrid className={styles.navlink__icon} />
            </div>
            Inicio
          </Link>
        </li>
        <li>
          <Link href="/villas" className={styles.navlink}>
            <div className={styles["navlink__icon-content"]}>
              <House className={styles.navlink__icon} />
            </div>
            Villas
          </Link>
        </li>
        <li>
          <Link href="/cuotas" className={styles.navlink}>
            <div className={styles["navlink__icon-content"]}>
              <Newspaper className={styles.navlink__icon} />
            </div>
            Cuotas
          </Link>
        </li>
        <li>
          <Link href="/reportes" className={styles.navlink}>
            <div className={styles["navlink__icon-content"]}>
              <ChartPie className={styles.navlink__icon} />
            </div>
            Reportes
          </Link>
        </li>
        <li>
          <Link href="/configuracion" className={styles.navlink}>
            <div className={styles["navlink__icon-content"]}>
              <Settings className={styles.navlink__icon} />
            </div>
            Configuración
          </Link>
        </li>
      </ul>
    </nav>
  );
}
