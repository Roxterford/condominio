import type { LucideIcon } from "lucide-react";
import {
  ArrowLeftRight,
  CreditCardMinus,
  HandCoins,
  Home,
  LayoutDashboard,
  Newspaper,
  ReceiptText,
  Wrench,
} from "lucide-react";

export type AccionSistema = {
  id: string;
  titulo: string;
  descripcion: string;
  icono: LucideIcon;
  ruta: string;
  claves: string[];
  destacada?: boolean;
};

export const ACCIONES_DEL_SISTEMA: AccionSistema[] = [
  {
    id: "registrar-cuota",
    titulo: "Registrar una cuota",
    descripcion: "Nueva cuota mensual o especial",
    icono: ReceiptText,
    ruta: "/cuotas/registrar",
    destacada: true,
    claves: [
      "registrar cuota",
      "registrar una cuota",
      "nueva cuota",
      "emitir cuota",
      "mensualidad",
      "cuota especial",
      "cuota de administración",
      "prorrateo",
      "recalcular cuota",
      "monto de la cuota",
    ],
  },
  {
    id: "registrar-pago",
    titulo: "Registrar un pago",
    descripcion: "Abono o liquidación de una villa",
    icono: HandCoins,
    ruta: "/pagos/registrar",
    destacada: true,
    claves: [
      "registrar pago",
      "registrar un pago",
      "nuevo pago",
      "cobro",
      "abono",
      "cancelar cuota",
      "liquidar",
      "saldar",
      "transferencia",
      "efectivo",
    ],
  },
  {
    id: "registrar-gasto",
    titulo: "Registrar un gasto",
    descripcion: "Egreso a un proveedor",
    icono: CreditCardMinus,
    ruta: "/operaciones",
    destacada: true,
    claves: [
      "registrar gasto",
      "registrar un gasto",
      "nuevo gasto",
      "egreso",
      "débito",
      "pagar proveedor",
      "desembolso",
      "gasto común",
    ],
  },
  {
    id: "ver-cuotas",
    titulo: "Ver cuotas",
    descripcion: "Listado y recaudación",
    icono: Newspaper,
    ruta: "/cuotas",
    claves: [
      "ver cuotas",
      "listado de cuotas",
      "recaudación",
      "cuotas mensuales",
      "cuotas especiales",
    ],
  },
  {
    id: "ver-villas",
    titulo: "Ver villas",
    descripcion: "Unidades y propietarios",
    icono: Home,
    ruta: "/villas",
    claves: ["ver villas", "listado de villas", "unidades", "propietarios"],
  },
  {
    id: "ver-operaciones",
    titulo: "Ver operaciones",
    descripcion: "Pagos y gastos registrados",
    icono: ArrowLeftRight,
    ruta: "/operaciones",
    claves: [
      "ver operaciones",
      "movimientos",
      "finanzas",
      "transacciones",
      "historial",
    ],
  },
  {
    id: "ver-administracion",
    titulo: "Ver Administración",
    descripcion: "Outbox de eventos y DLQ",
    icono: Wrench,
    ruta: "/admin/outbox",
    claves: ["administración", "outbox", "eventos", "dlq", "cola de mensajes"],
  },
  {
    id: "ir-dashboard",
    titulo: "Ir al Dashboard",
    descripcion: "Resumen general del condominio",
    icono: LayoutDashboard,
    ruta: "/dashboard",
    claves: ["dashboard", "inicio", "resumen", "panel principal", "home"],
  },
];

export function accionesDestacadas(): AccionSistema[] {
  return ACCIONES_DEL_SISTEMA.filter((accion) => accion.destacada);
}

export function buscarAcciones(busqueda: string): AccionSistema[] {
  const consulta = normalizar(busqueda);
  if (!consulta) return [];

  const tokens = consulta.split(" ").filter(Boolean);

  const resultados = ACCIONES_DEL_SISTEMA.map((accion) => {
    const buscable = normalizar(
      [accion.titulo, accion.descripcion, ...accion.claves].join(" "),
    );
    const coincidenTodos = tokens.every((token) => buscable.includes(token));
    if (!coincidenTodos) return null;

    const titulo = normalizar(accion.titulo);
    let puntuacion = 1;
    if (titulo.includes(consulta)) puntuacion = 4;
    else if (titulo.split(" ").some((palabra) => palabra.startsWith(consulta)))
      puntuacion = 3;
    else if (buscable.includes(consulta)) puntuacion = 2;

    return { accion, puntuacion };
  });

  return resultados
    .filter(
      (r): r is { accion: AccionSistema; puntuacion: number } => r !== null,
    )
    .sort((a, b) => b.puntuacion - a.puntuacion)
    .map((resultado) => resultado.accion);
}

function normalizar(texto: string): string {
  return texto
    .toLocaleLowerCase("es")
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/\s+/g, " ")
    .trim();
}
