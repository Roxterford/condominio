import { TipoDeCuota } from "@/providers/graphql/graphql";
import { Badge } from "./ui/badge";

const classmap: Record<TipoDeCuota, string> = {
  [TipoDeCuota.Especial]: "bg-purple-100 text-purple-700",
  [TipoDeCuota.Regular]: "bg-teal-100 text-teal-700",
  [TipoDeCuota.Semilla]: "bg-slate-100 text-slate-700",
};

interface TipoCuotaTagProps {
  type: TipoDeCuota;
}
export function TipoCuotaTag({ type }: TipoCuotaTagProps) {
  return <Badge className={classmap[type]}>{type}</Badge>;
}
