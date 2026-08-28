import type { ReactNode } from "react";

interface StatCardProps {
  title: string;
  value: string | number;
  subtitle: string;
  icon: ReactNode;
  color: string;
}

export default function StatCard({
  title,
  value,
  subtitle,
  icon,
  color,
}: StatCardProps) {
  return (
    <div className="bg-white border border-gray-200 rounded-3xl p-6 flex items-start justify-between min-h-[125px]">

      {/* Left */}
      <div>

        <p className="text-sm text-gray-500">
          {title}
        </p>

        <h2 className="text-4xl font-bold text-gray-800 mt-3">
          {value}
        </h2>

        <p className="text-sm text-gray-400 mt-3">
          {subtitle}
        </p>

      </div>

      {/* Icon */}
      <div className={`w-14 h-14 rounded-2xl flex items-center justify-center ${color}`}>

        <div className="scale-110">
          {icon}
        </div>

      </div>

    </div>
  )
}