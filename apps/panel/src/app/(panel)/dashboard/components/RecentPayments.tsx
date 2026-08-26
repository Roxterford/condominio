const payments = [
  {
    villa: "Villa 12",
    resident: "Carlos Mendoza",
    amount: "$350.00",
    type: "Mensualidad",
    status: "Completado",
  },
  {
    villa: "Villa 05",
    resident: "María Rodríguez",
    amount: "$350.00",
    type: "Mensualidad",
    status: "Completado",
  },
  {
    villa: "Villa 23",
    resident: "Juan Pérez",
    amount: "$500.00",
    type: "Especial",
    status: "Pendiente",
  },
  {
    villa: "Villa 08",
    resident: "Ana Gómez",
    amount: "$350.00",
    type: "Mensualidad",
    status: "Completado",
  },
]

export default function RecentPayments() {
  return (
    <div className="bg-white overflow-x-auto rounded-3xl border border-gray-200 p-6 h-full">

      {/* Header */}
      <div className="mb-6">

        <h2 className="text-2xl font-bold text-gray-800">
          Pagos Recientes
        </h2>

        <p className="text-gray-500 mt-1">
          Últimos pagos registrados en el sistema
        </p>

      </div>

      {/* Table Header */}
      <div className="grid grid-cols-5 text-sm text-gray-400 pb-4 border-b border-gray-100">

        <span>Villa</span>
        <span>Residente</span>
        <span>Monto</span>
        <span>Tipo</span>
        <span>Estado</span>

      </div>

      {/* Rows */}
      <div className="divide-y divide-gray-100">

        {payments.map((payment, index) => (
          <div
            key={index}
            className="grid grid-cols-5 items-center py-5 text-sm"
          >

            <div className="font-semibold text-gray-800">
              {payment.villa}
            </div>

            <div className="text-gray-600">
              {payment.resident}
            </div>

            <div className="font-semibold text-gray-800">
              {payment.amount}
            </div>

            <div>

              <span className="px-3 py-1 rounded-full text-xs bg-gray-100 text-gray-700">
                {payment.type}
              </span>

            </div>

            <div>

              <span
                className={`px-3 py-1 rounded-full text-xs ${
                  payment.status === "Completado"
                    ? "bg-green-100 text-green-700"
                    : "bg-yellow-100 text-yellow-700"
                }`}
              >
                {payment.status}
              </span>

            </div>

          </div>
        ))}

      </div>

    </div>
  )
}