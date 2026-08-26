import { EstadoDeDeuda } from "@/providers/graphql/graphql";
import { Badge } from "./ui/badge";

const classmap: Record<EstadoDeDeuda, string> = {
  [EstadoDeDeuda.Abonada]: "bg-yellow-100 text-yellow-700",
  [EstadoDeDeuda.Pendiente]: "bg-red-100 text-red-700",
  [EstadoDeDeuda.Saldada]: "bg-green-100 text-green-700",
};

interface EstodoDeudaTagProps {
  state: EstadoDeDeuda;
}
export function EstadoDeudaTag({ state }: EstodoDeudaTagProps) {
  return <Badge className={classmap[state]}>{state}</Badge>;
}
