"use client"

import { Button } from "@/components/ui/button"
import { useSidebar } from "@/contexts/sidebar-context"
import { PagoDetalleSidebar } from "./pago-detalle-sidebar"

type PagoRow = {
  id: string
  total: number
  fecha: any
  unidad: string
  disponible: number
  monto: number
  moneda: any
  tasa: number
  destinado: number
}

type PagosTableProps = {
  data: PagoRow[]
  total: number
}

export function PagosTable({ data }: PagosTableProps) {
  const { setView } = useSidebar()

  const verDetalle = (pago: PagoRow) => {
    setView({
      component: <PagoDetalleSidebar pago={pago} />,
      title: `Pago #${pago.id.slice(0, 8)}`,
    })
  }

  return (
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
          {data.map((pago) => (
            <tr key={pago.id}>
              <td className="p-10">
                <button className="link">{pago.id}</button>
              </td>
              <td>${pago.total}</td>
              <td className="p-10">{pago.fecha}</td>
              <td className="p-10">
                <Button variant="outline" onClick={() => verDetalle(pago)}>
                  Ver detalles
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
