import { Badge } from "./ui/badge";

interface DeudaUnidadTagProps {
  pending: boolean;
}
export function DeudaUnidadTag({ pending }: DeudaUnidadTagProps) {
  if (pending)
    return (
      <Badge className="bg-yellow-100 text-yellow-700">Deuda pendiente</Badge>
    );
  return <Badge className="bg-green-100 text-green-700">Al día</Badge>;
}
