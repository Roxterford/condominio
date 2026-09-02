"use client";

import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { ComponentProps, useRef } from "react";

type MoneyInputProps = Omit<
  ComponentProps<typeof Input>,
  "type" | "value" | "defaultValue" | "onChange"
> & {
  /**
   * Valor real en la unidad menor de la moneda.
   * Ej.: 123456 representa 1.234,56.
   */
  value: number;
  onValueChange: (cents: number) => void;
  locale?: string;
  fractionDigits?: number;
};

function formatFromMinorUnits(
  valueInCents: number,
  locale: string,
  fractionDigits: number,
) {
  const divisor = 10 ** fractionDigits;

  return new Intl.NumberFormat(locale, {
    minimumFractionDigits: fractionDigits,
    maximumFractionDigits: fractionDigits,
  }).format(valueInCents / divisor);
}

function parseToMinorUnits(value: string) {
  // "1.234,56" -> "123456" -> 123456
  const digits = value.replace(/\D/g, "");

  if (!digits) return 0;

  const parsed = Number(digits);

  // Evita que un input excepcionalmente largo produzca un valor inseguro.
  return Number.isSafeInteger(parsed) ? parsed : 0;
}

export function MoneyInput({
  value: valueInCents,
  onValueChange,
  locale = "es-VE",
  fractionDigits = 2,
  className,
  onFocus,
  ...props
}: MoneyInputProps) {
  const inputRef = useRef<HTMLInputElement>(null);

  const displayValue = formatFromMinorUnits(
    valueInCents,
    locale,
    fractionDigits,
  );

  function moveCursorToEnd() {
    requestAnimationFrame(() => {
      const input = inputRef.current;
      if (!input) return;

      const end = input.value.length;
      input.setSelectionRange(end, end);
    });
  }

  return (
    <Input
      {...props}
      ref={inputRef}
      type="text"
      inputMode="numeric"
      autoComplete="off"
      value={displayValue}
      className={cn("text-right tabular-nums", className)}
      aria-label={props["aria-label"] ?? "Monto"}
      onFocus={(event) => {
        moveCursorToEnd();
        onFocus?.(event);
      }}
      onChange={(event) => {
        const cents = parseToMinorUnits(event.target.value);
        onValueChange(cents);
        moveCursorToEnd();
      }}
    />
  );
}
