"use client";

import { Button } from "@/components/ui/button";
import { RegistrarPagoOverlay } from "@/features/administracion/components/registrar-pago-overlay";
import { useOverlay } from "@/hooks/useOverlay";
import { CreditCardPlus } from "lucide-react";
import { VillaPageQuery } from "@/providers/graphql/graphql";

export function RegistrarPagoButton({
  unidad,
}: {
  unidad: VillaPageQuery["unidad"];
}) {
  const registrarPago = useOverlay();

  return (
    <>
      <Button onClick={registrarPago.open}>
        <CreditCardPlus /> Registrar pago
      </Button>
      <RegistrarPagoOverlay
        {...registrarPago.overlayProps}
        unidad={unidad}
        initialFocus="monto"
      />
    </>
  );
}
