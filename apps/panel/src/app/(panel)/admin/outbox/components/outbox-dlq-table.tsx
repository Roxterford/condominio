import { RotateCcw } from "lucide-react";

import type { OutboxDlqEvent } from "@/lib/outbox-admin";
import { decodePayload } from "@/lib/outbox-admin";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface OutboxDlqTableProps {
  events: OutboxDlqEvent[];
  retryingId: string | null;
  onRetry: (id: string) => void;
}

function formatDate(value: string | null | undefined): string {
  if (!value) return "-";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString("es-VE");
}

export function OutboxDlqTable({
  events,
  retryingId,
  onRetry,
}: OutboxDlqTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">Evento</TableHead>
          <TableHead className="table__head">Falló el</TableHead>
          <TableHead className="table__head">Error</TableHead>
          <TableHead className="table__head text-center">Reintentos</TableHead>
          <TableHead className="table__head">Payload</TableHead>
          <TableHead className="table__head text-right">Acciones</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {events.length === 0 ? (
          <TableRow>
            <TableCell colSpan={6} className="text-center py-8 text-gray-500">
              La cola DLQ está vacía
            </TableCell>
          </TableRow>
        ) : (
          events.map((event) => (
            <TableRow key={event.ID}>
              <TableCell>
                <span className="font-medium">{event.EventName}</span>
                <div className="text-xs text-gray-400">{event.ID}</div>
              </TableCell>
              <TableCell>{formatDate(event.FailedAt)}</TableCell>
              <TableCell className="max-w-[220px]">
                <Badge
                  variant="destructive"
                  className="whitespace-normal h-auto py-1"
                >
                  {event.ErrorMessage || "Error desconocido"}
                </Badge>
              </TableCell>
              <TableCell className="text-center">
                <Badge variant="outline">{event.Retries}</Badge>
              </TableCell>
              <TableCell className="max-w-[220px]">
                <details className="group">
                  <summary className="cursor-pointer text-blue-600 hover:text-blue-700 text-xs">
                    Ver payload
                  </summary>
                  <pre className="mt-2 whitespace-pre-wrap rounded-md bg-gray-50 p-2 text-xs text-gray-700">
                    {decodePayload(event.Payload)}
                  </pre>
                </details>
              </TableCell>
              <TableCell className="text-right">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={retryingId !== null}
                  onClick={() => onRetry(event.ID)}
                >
                  <RotateCcw
                    className={
                      retryingId === event.ID ? "animate-spin" : undefined
                    }
                  />
                  Reintentar
                </Button>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}
