const tabs = [
  {
    name: "Vista General",
    active: true,
  },
  {
    name: "Pagos",
  },
  {
    name: "Estado de villas",
  },
]

export default function DashboardTabs() {
  return (
    <div className="mt-10 border-b border-gray-200">

      <div className="flex items-center overflow-x-auto gap-8">

        {tabs.map((tab, index) => (
          <button
            key={index}
            className={`pb-4 text-sm font-medium transition-all ${
              tab.active
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-gray-500 hover:text-gray-800"
            }`}
          >
            {tab.name}
          </button>
        ))}

      </div>

    </div>
  )
}