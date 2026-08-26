import type { OutboxPendingEvent } from "@/lib/outbox-admin";
import { decodePayload } from "@/lib/outbox-admin";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface OutboxPendingTableProps {
  events: OutboxPendingEvent[];
}

function formatDate(value: string | null | undefined): string {
  if (!value) return "-";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString("es-VE");
}

export function OutboxPendingTable({ events }: OutboxPendingTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="table__head">Evento</TableHead>
          <TableHead className="table__head">Correlation ID</TableHead>
          <TableHead className="table__head">Ocurrido</TableHead>
          <TableHead className="table__head">Creado</TableHead>
          <TableHead className="table__head text-center">Reintentos</TableHead>
          <TableHead className="table__head">Payload</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {events.length === 0 ? (
          <TableRow>
            <TableCell colSpan={6} className="text-center py-8 text-gray-500">
              No hay eventos pendientes
            </TableCell>
          </TableRow>
        ) : (
          events.map((event) => (
            <TableRow key={event.ID}>
              <TableCell>
                <span className="font-medium">{event.EventName}</span>
                <div className="text-xs text-gray-400">{event.ID}</div>
              </TableCell>
              <TableCell>
                {event.CorrelationID || (
                  <span className="text-gray-400">-</span>
                )}
              </TableCell>
              <TableCell>{formatDate(event.OccurredAt)}</TableCell>
              <TableCell>{formatDate(event.CreatedAt)}</TableCell>
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
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}
