"use client";

import { Button } from "@/components/ui/button";
import { useOverlay } from "@/hooks/useOverlay";
import { UserRoundPlus } from "lucide-react";
import { RegistrarTitularOverlay } from "./registrar-titular-overlay";

export function RegistrarTitularButton({
  variant = "default",
}: {
  variant?: "default" | "ghost";
}) {
  const registrarTitular = useOverlay({ closeOnDone: true });

  return (
    <>
      <Button variant={variant} onClick={registrarTitular.open}>
        <UserRoundPlus /> Registrar titular
      </Button>
      <RegistrarTitularOverlay {...registrarTitular.overlayProps} />
    </>
  );
}
