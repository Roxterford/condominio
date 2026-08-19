"use client";

import { Button } from "@/components/ui/button";
import { useAppForm } from "@/hooks/useAppForm";
import { useOverlay } from "@/hooks/useOverlay";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { Proveedor, RegistrarCuotaDto } from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { Check, Loader2, Plus } from "lucide-react";
import { useEffect, useState, type SubmitEventHandler } from "react";
import { toast } from "sonner";
import { AgregarGastoOverlay } from "../agregar-gasto-overlay";
import { DesgloseDeGastoItem, DesgloseDeGastos } from "../desglose_de_gastos";
import { GastoSidebar, GastoSidebarData } from "../gasto_sidebar/gasto-sidebar";
import { RegistrarGastoOverlay } from "../registrar_gasto_overlay";
import {
  SeleccionarGastosOverlay,
  SeleccionarGastosOverlayProps,
} from "../seleccionar-gastos-overlay";
import { EstrategiaDeDistribucionField } from "./fields/estrategia-de-distribucion";
import { PeriodoField } from "./fields/periodo";
import { TipoDeCuotaField } from "./fields/tipo-de-cuota";
import { defaultValues, RegistrarCuotaFormSchema } from "./schema";

const RegistrarCuotaMutation = graphql(/* GraphQL */ `
  mutation RegistrarCuota($input: RegistrarCuotaDTO!) {
    registrarCuota(input: $input) {
      __typename
      ... on CuotaRegular {
        id
      }
      ... on CuotaEspecial {
        id
      }
    }
  }
`);

export interface RegistrarCuotaFormProps {
  proveedores: Pick<Proveedor, "id" | "nombre">[];
}

export function RegistrarCuotaForm({ proveedores }: RegistrarCuotaFormProps) {
  const router = useRouter();
  const agregarGastoOverlay = useOverlay();
  const registrarGastoOverlay = useOverlay();
  const seleccionarGastosOverlay = useOverlay();
  const gastoSidebar = useOverlay();
  const [selectedGasto, setSelectedGasto] = useState<GastoSidebarData | null>(
    null,
  );

  const [gastos, setGastos] = useState<DesgloseDeGastoItem[]>([]);

  const registrar = useMutation({
    mutationFn: (input: RegistrarCuotaDto) =>
      execute(RegistrarCuotaMutation, { input }),
  });

  const form = useAppForm({
    defaultValues,
    validators: {
      onChange: RegistrarCuotaFormSchema,
      onBlur: RegistrarCuotaFormSchema,
    },
    onSubmit: async ({ value }) => {
      const res = await registrar.mutateAsync({
        anio: value.anio_actual ? new Date().getFullYear() : value.anio,
        mes: value.mes,
        fecha_limite: value.fecha_limite,
        gastos: value.gastos.filter(Boolean),
        tipo: value.tipo,
      });

      if (res.errors?.length) {
        return toast.error(res.errors.at(0)?.message, {
          description: JSON.stringify(res.errors.at(0)?.locations, null, 4),
        });
      }

      toast.success("Cuota registrada con exito");
      const cuota = res.data?.registrarCuota;
      if (cuota) {
        router.push(`/cuotas/${cuota.id}`);
      }
    },
  });

  const handleGastosSelectos: SeleccionarGastosOverlayProps["onDone"] = (
    gastos,
  ) => {
    setGastos((prev) => [...prev, ...gastos]);
  };

  const handleDesglosePress = (gasto: DesgloseDeGastoItem) => {
    setSelectedGasto({
      id: gasto.operacion,
      concepto: gasto.concepto,
      monto_total: gasto.total,
      fecha: gasto.fecha,
      tasa: gasto.tasa,
      proveedor: gasto.proveedor,
    });
    gastoSidebar.open();
  };

  const handleSubmit: SubmitEventHandler = (event) => {
    event.preventDefault();
    event.stopPropagation();
    form.handleSubmit();
  };

  useEffect(() => {
    form.setFieldValue("gastos", (prev) => gastos.map((it) => it.operacion));
  }, [gastos]);

  return (
    <>
      <form.Subscribe>
        {(it) => <pre>{JSON.stringify(it.values, null, 4)}</pre>}
      </form.Subscribe>
      <form className="grid gap-5" onSubmit={handleSubmit}>
        <TipoDeCuotaField form={form} />
        <PeriodoField form={form} />

        <section>
          <div className="flex justify-between items-center">
            <h3>Desglose de gastos</h3>
            <Button
              type="button"
              variant="outline"
              onClick={agregarGastoOverlay.open}
            >
              <Plus /> Agregar Gasto
            </Button>
          </div>
          <DesgloseDeGastos
            data={gastos}
            onGastoPress={handleDesglosePress}
            onRemove={(g) =>
              setGastos((prev) =>
                prev.filter((it) => it.operacion != g.operacion),
              )
            }
          />
        </section>

        <EstrategiaDeDistribucionField form={form} />

        <div className="flex gap-2 justify-end">
          <Button type="button" variant="outline" onClick={() => form.reset()}>
            Cancelar
          </Button>
          <form.Subscribe
            selector={(state) => state.isValid}
            children={(isValid) => (
              <Button type="submit" disabled={registrar.isPending || !isValid}>
                Registrar
                {registrar.isPending ? (
                  <Loader2 className="animate-spin" />
                ) : (
                  <Check />
                )}
              </Button>
            )}
          />
        </div>
      </form>
      <GastoSidebar
        data={selectedGasto || undefined}
        {...gastoSidebar.overlayProps}
      />
      <AgregarGastoOverlay
        onSelect={(op) => {
          agregarGastoOverlay.close();
          if (op == "nuevo") return registrarGastoOverlay.open();
          if (op == "seleccionar") return seleccionarGastosOverlay.open();
        }}
        {...agregarGastoOverlay.overlayProps}
      />
      <SeleccionarGastosOverlay
        {...seleccionarGastosOverlay.overlayProps}
        omitIDs={form.getFieldValue("gastos")}
        onDone={(values) => {
          seleccionarGastosOverlay.close();
          handleGastosSelectos(values);
        }}
      />
      <RegistrarGastoOverlay
        proveedores={proveedores}
        {...registrarGastoOverlay.overlayProps}
      />
    </>
  );
}
