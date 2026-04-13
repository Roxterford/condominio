"use client";

import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  CuotasTable,
  CuotasTableProps,
  CuotasTableType,
} from "@/features/administracion/components/cuotas_table/cuotas_table";
import { useState } from "react";

export interface CuotasPageTapsProps {
  cuotas: CuotasTableProps["data"];
}

export function CuotasPageTaps({ cuotas }: CuotasPageTapsProps) {
  const [tap, setTap] = useState<TapValue>("todas");

  const map_tap_to_show: Record<TapValue, CuotasTableType> = {
    todas: "default",
    regulares: "regular",
    especiales: "especial",
  };

  return (
    <>
      <Taps onChangeAction={setTap} defaultValue={tap} />
      <CuotasTable data={cuotas} type={map_tap_to_show[tap]} />
    </>
  );
}

type TapValue = "todas" | "regulares" | "especiales";

export function Taps({
  onChangeAction,
  defaultValue = "todas",
}: {
  defaultValue?: TapValue;
  onChangeAction?: (value: TapValue) => void;
}) {
  return (
    <Tabs
      defaultValue={defaultValue}
      onValueChange={(v) => onChangeAction?.(v as TapValue)}
    >
      <TabsList variant="line">
        <TabsTrigger value="todas">Todas</TabsTrigger>
        <TabsTrigger value="regulares">Mensualidades</TabsTrigger>
        <TabsTrigger value="especiales">Cuotas Especiales</TabsTrigger>
      </TabsList>
    </Tabs>
  );
}
