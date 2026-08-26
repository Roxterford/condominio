import { EstadoDeUnidad } from "@/providers/graphql/graphql";
import { Badge } from "./ui/badge";

const classmap: Record<EstadoDeUnidad, string> = {
  [EstadoDeUnidad.Activa]: "bg-emerald-100 text-emerald-700",
  [EstadoDeUnidad.EnLitigio]: "bg-rose-100 text-rose-700",
  [EstadoDeUnidad.Exenta]: "bg-purple-100 text-purple-700",
  [EstadoDeUnidad.Inhabitada]: "bg-slate-100 text-slate-700",
  [EstadoDeUnidad.Preventa]: "bg-sky-100 text-sky-700",
  [EstadoDeUnidad.Suspendida]: "bg-amber-100 text-amber-700",
};

interface EstodoUnidadTagProps {
  state: EstadoDeUnidad;
}
export function EstadoUnidadTag({ state }: EstodoUnidadTagProps) {
  return <Badge className={classmap[state]}>{state}</Badge>;
}
