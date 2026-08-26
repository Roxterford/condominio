 "use client" 

import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts"

const data = [
  {
    month: "Enero",
    ingresos: 4000,
    gastos: 3200,
  },
  {
    month: "Febrero",
    ingresos: 4200,
    gastos: 3600,
  },
  {
    month: "Marzo",
    ingresos: 4100,
    gastos: 3500,
  },
  {
    month: "Abril",
    ingresos: 4500,
    gastos: 3900,
  },
  {
    month: "Mayo",
    ingresos: 4800,
    gastos: 4200,
  },
  {
    month: "Junio",
    ingresos: 5000,
    gastos: 4500,
  },
]

export default function FinancialChart() {
  return (
    <div className="bg-white rounded-3xl border border-gray-200 p-6 h-full">

      {/* Header */}
      <div className="mb-8">

        <h2 className="text-2xl font-bold text-gray-800">
          Resumen Financiero
        </h2>

        <p className="text-gray-500 mt-1">
          Ingresos vs Gastos de los últimos 6 meses
        </p>

      </div>

      {/* Chart */}
      <div className="h-[300px]">

        <ResponsiveContainer width="100%" height="100%">

          <BarChart data={data}>

            <CartesianGrid
              strokeDasharray="3 3"
              vertical={false}
              stroke="#E5E7EB"
            />

            <XAxis
              dataKey="month"
              tickLine={false}
              axisLine={false}
              tick={{ fill: "#6B7280", fontSize: 13 }}
            />

            <YAxis
              tickLine={false}
              axisLine={false}
              tick={{ fill: "#6B7280", fontSize: 13 }}
            />

            <Tooltip />

            <Legend />

            <Bar
              dataKey="ingresos"
              fill="#14B8A6"
              radius={[6, 6, 0, 0]}
            />

            <Bar
              dataKey="gastos"
              fill="#F97316"
              radius={[6, 6, 0, 0]}
            />

          </BarChart>

        </ResponsiveContainer>

      </div>

    </div>
  )
}