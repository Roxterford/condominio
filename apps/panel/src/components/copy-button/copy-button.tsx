"use client";

import * as React from "react";
import { Check, Copy } from "lucide-react";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

type CopyButtonProps = {
  value: string;
  what?: string;
} & Omit<React.ComponentProps<"button">, "onClick">;

export function CopyButton({
  value,
  what = "Información",
  className,
  ...props
}: CopyButtonProps) {
  const [copied, setCopied] = React.useState(false);

  React.useEffect(() => {
    if (!copied) return;
    const timer = setTimeout(() => setCopied(false), 1500);
    return () => clearTimeout(timer);
  }, [copied]);

  const copiar = async () => {
    try {
      await navigator.clipboard.writeText(value);
      setCopied(true);
      toast.success("Copiado", { description: `${what} copiado` });
    } catch {
      toast.error("Error", {
        description: `No se pudo copiar ${what.toLowerCase()}`,
      });
    }
  };

  return (
    <button
      onClick={copiar}
      aria-label={`Copiar ${what.toLowerCase()}`}
      title={`Copiar ${what.toLowerCase()}`}
      className={cn(
        "shrink-0 text-muted-foreground transition-colors hover:text-foreground",
        className
      )}
      {...props}
    >
      {copied ? (
        <Check className="size-3.5 text-green-600" />
      ) : (
        <Copy className="size-3.5" />
      )}
    </button>
  );
}