import type { ReactNode } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

interface OutboxStatCardProps {
  title: string;
  value?: number;
  loading?: boolean;
  color?: string;
  icon?: ReactNode;
}

export function OutboxStatCard({
  title,
  value,
  loading,
  color = "bg-blue-100",
  icon,
}: OutboxStatCardProps) {
  return (
    <Card size="sm" className="gap-2">
      <CardContent className="flex items-center gap-4">
        <div
          className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl ${color}`}
        >
          {icon}
        </div>
        <div>
          <p className="text-sm text-gray-500">{title}</p>
          {loading ? (
            <Skeleton className="mt-1 h-8 w-16" />
          ) : (
            <p className="text-3xl font-bold text-gray-800">{value ?? "-"}</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
