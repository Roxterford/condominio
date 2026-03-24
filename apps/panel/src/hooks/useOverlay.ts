import { OverlayProps } from "@/components/overlay/overlay";
import { useCallback, useState } from "react";

export const useOverlay = (initialState = false) => {
  const [isOpen, setIsOpen] = useState(initialState);

  const open = useCallback(() => setIsOpen(true), []);
  const close = useCallback(() => setIsOpen(false), []);
  const toggle = useCallback(() => setIsOpen((prev) => !prev), []);

  return {
    isOpen,
    setIsOpen, // Útil para vincularlo directamente a componentes de shadcn
    open,
    close,
    toggle,
    // Propiedades listas para esparcir en componentes tipo Dialog/Sheet
    overlayProps: {
      open: isOpen,
      onOpenChange: setIsOpen,
    } as OverlayProps,
  };
};
