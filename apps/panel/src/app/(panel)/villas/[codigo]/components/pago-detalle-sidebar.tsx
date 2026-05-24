"use client"

import { format } from "date-fns"
import { es } from "date-fns/locale"
import {
  CreditCard,
  Hash,
  DollarSign,
  Calendar,
  ArrowLeftRight,
  PiggyBank,
  Target,
} from "lucide-react"

type PagoData = {
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

type Props = {
  pago: PagoData
}

export function PagoDetalleSidebar({ pago }: Props) {
  const fecha = pago.fecha
    ? format(new Date(pago.fecha), "PPP", { locale: es })
    : "—"

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-2xl border border-gray-200 p-6 space-y-5">
        <h3 className="text-lg font-semibold text-gray-800 flex items-center gap-2">
          <CreditCard size={20} className="text-blue-600" />
          Información del pago
        </h3>

        <div className="grid gap-4">
          <InfoRow
            icon={<Hash size={16} />}
            label="ID"
            value={pago.id}
          />
          <InfoRow
            icon={<DollarSign size={16} />}
            label="Monto total"
            value={`$${pago.total.toFixed(2)}`}
          />
          <InfoRow
            icon={<DollarSign size={16} />}
            label="Monto original"
            value={`${pago.monto.toFixed(2)} ${pago.moneda}`}
          />
          <InfoRow
            icon={<Calendar size={16} />}
            label="Fecha"
            value={fecha}
          />
          <InfoRow
            icon={<ArrowLeftRight size={16} />}
            label="Tasa de cambio"
            value={pago.tasa ? `Bs. ${pago.tasa.toFixed(2)}` : "N/A"}
          />
        </div>
      </div>

      <div className="bg-white rounded-2xl border border-gray-200 p-6 space-y-5">
        <h3 className="text-lg font-semibold text-gray-800 flex items-center gap-2">
          <PiggyBank size={20} className="text-blue-600" />
          Distribución
        </h3>

        <div className="grid gap-4">
          <InfoRow
            icon={<Target size={16} />}
            label="Destinado"
            value={`$${pago.destinado.toFixed(2)}`}
          />
          <InfoRow
            icon={<PiggyBank size={16} />}
            label="Disponible"
            value={`$${pago.disponible.toFixed(2)}`}
          />
        </div>
      </div>
    </div>
  )
}

function InfoRow({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: string
}) {
  return (
    <div className="flex items-center gap-3 pb-2 border-b border-gray-100 last:border-0">
      <span className="text-gray-400">{icon}</span>
      <span className="text-sm text-gray-500 flex-1">{label}</span>
      <span className="text-sm font-medium text-gray-800 text-right truncate max-w-[200px]">
        {value}
      </span>
    </div>
  )
}
