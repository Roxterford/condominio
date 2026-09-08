import { useCallback, useState } from "react";

export interface OverlayProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  onDone?: () => void;
}

export interface UseOverlayOptions {
  initialState?: boolean;
  closeOnDone?: boolean;
  onDone?: () => void;
}

export function useOverlay(options?: UseOverlayOptions) {
  const _options = {
    initialState: false,
    closeOnDone: false,
    ...options,
  };

  const [isOpen, setIsOpen] = useState(_options.initialState);

  const open = useCallback(() => setIsOpen(true), []);
  const close = useCallback(() => setIsOpen(false), []);
  const toggle = useCallback(() => setIsOpen((prev) => !prev), []);

  return {
    isOpen,
    setIsOpen,
    open,
    close,
    toggle,
    overlayProps: {
      open: isOpen,
      onOpenChange: setIsOpen,
      onDone: _options.closeOnDone
        ? () => {
            _options.onDone?.();
            setIsOpen(false);
          }
        : _options.onDone,
    } as OverlayProps,
  };
}
