import { useCallback, useEffect, useRef } from "react";

type Point = { clientX: number; clientY: number };

/**
 * Distingue un click simple de un doble click sin sacrificar la respuesta.
 *
 * El click simple ejecuta `onSingle` inmediatamente (el drawer abre al
 * instante). Para reconocer el doble click se registra un listener de `click`
 * global (fase captura) durante `interval` ms: el segundo click de un doble
 * click real aterriza sobre el overlay del drawer, por lo que hay que
 * capturarlo a nivel de documento y comparar su posicion con la del primero
 * (radio `hitRadius`) antes de ejecutar `onDouble`.
 */
export function useSingleDoubleClick<Arg>({
  onSingle,
  onDouble,
  interval = 300,
  hitRadius = 12,
}: {
  onSingle: (arg: Arg) => void;
  onDouble: (arg: Arg) => void;
  interval?: number;
  hitRadius?: number;
}): (arg: Arg, event: Point) => void {
  const onSingleRef = useRef(onSingle);
  const onDoubleRef = useRef(onDouble);

  useEffect(() => {
    onSingleRef.current = onSingle;
    onDoubleRef.current = onDouble;
  }, [onSingle, onDouble]);

  const controllerRef = useRef<{
    clearTimeout: () => void;
    removeListener: () => void;
  } | null>(null);

  const cancelPending = useCallback(() => {
    const controller = controllerRef.current;
    if (!controller) return;
    controller.clearTimeout();
    controller.removeListener();
    controllerRef.current = null;
  }, []);

  useEffect(() => cancelPending, [cancelPending]);

  return useCallback(
    (arg: Arg, event: Point) => {
      cancelPending();

      let timeoutId: ReturnType<typeof setTimeout> | null = null;
      let detachListener = () => {};

      const finish = () => {
        if (timeoutId !== null) clearTimeout(timeoutId);
        detachListener();
      };

      const onDocumentClick = (click: MouseEvent) => {
        if (!controllerRef.current) return;
        const distance = Math.hypot(
          click.clientX - event.clientX,
          click.clientY - event.clientY,
        );
        if (distance > hitRadius) return;
        finish();
        controllerRef.current = null;
        onDoubleRef.current(arg);
      };

      detachListener = () => {
        document.removeEventListener("click", onDocumentClick, true);
      };

      timeoutId = setTimeout(() => {
        controllerRef.current = null;
        finish();
      }, interval);

      document.addEventListener("click", onDocumentClick, true);

      controllerRef.current = {
        clearTimeout: () => {
          if (timeoutId !== null) clearTimeout(timeoutId);
        },
        removeListener: detachListener,
      };

      onSingleRef.current(arg);
    },
    [cancelPending, interval, hitRadius],
  );
}
