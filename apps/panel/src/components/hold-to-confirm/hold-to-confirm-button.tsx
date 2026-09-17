"use client";

import { Button } from "@/components/ui/button";
import { cn } from "cn";
import { useCallback, useEffect, useRef, useState } from "react";

export interface HoldToConfirmButtonProps {
  onConfirm: () => void;
  duration?: number;
  disabled?: boolean;
  className?: string;
  children?: React.ReactNode;
}

export function HoldToConfirmButton({
  onConfirm,
  duration = 2000,
  disabled,
  className,
  children,
}: HoldToConfirmButtonProps) {
  const [progress, setProgress] = useState(0);
  const rafRef = useRef<number | null>(null);
  const startRef = useRef(0);

  const stop = useCallback(() => {
    if (rafRef.current) cancelAnimationFrame(rafRef.current);
    rafRef.current = null;
    setProgress(0);
  }, []);

  useEffect(() => () => stop(), [stop]);

  const start = useCallback(
    (event: React.PointerEvent<HTMLButtonElement>) => {
      if (event.button !== 0) return;
      if (rafRef.current) return;

      startRef.current = performance.now();

      const tick = (now: number) => {
        const pct = Math.min((now - startRef.current) / duration, 1);
        setProgress(pct);

        if (pct >= 1) {
          stop();
          onConfirm();
          return;
        }

        rafRef.current = requestAnimationFrame(tick);
      };

      rafRef.current = requestAnimationFrame(tick);
    },
    [duration, onConfirm, stop],
  );

  return (
    <Button
      type="button"
      disabled={disabled}
      className={cn(
        "relative select-none overflow-hidden disabled:pointer-events-none",
        className,
      )}
      onPointerDown={start}
      onPointerUp={stop}
      onPointerLeave={stop}
      onPointerCancel={stop}
    >
      <span
        aria-hidden
        className="pointer-events-none absolute inset-y-0 left-0 bg-white/25"
        style={{ width: `${progress * 100}%` }}
      />
      <span className="relative inline-flex items-center gap-1.5">
        {children}
      </span>
    </Button>
  );
}