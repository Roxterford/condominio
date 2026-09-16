import type { ReactNode } from "react";
import { FlaskConical } from "lucide-react";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export type ResultRow = {
  key: string;
  groupLabel: string;
  icon: ReactNode;
  title: string;
  meta: string;
  typeLabel: string;
  mock?: boolean;
  mockHint?: string;
  trackSearch?: boolean;
  onClick: () => void;
};

const MOCK_HINT =
  "Dato de ejemplo: la búsqueda de propietarios aún está en desarrollo y no usa datos de la API.";

export function MockBadge({ hint = MOCK_HINT }: { hint?: string }) {
  return (
    <Tooltip>
      <TooltipTrigger render={<span className="inline-flex cursor-default" />}>
        <Badge
          variant="outline"
          className="cursor-default border-dashed text-muted-foreground"
        >
          <FlaskConical className="size-3" />
          Mock
        </Badge>
      </TooltipTrigger>
      <TooltipContent>{hint}</TooltipContent>
    </Tooltip>
  );
}

export function Highlighted({ text, query }: { text: string; query: string }) {
  const q = query.trim();
  if (!q) return <>{text}</>;
  const lower = text.toLocaleLowerCase();
  const ql = q.toLocaleLowerCase();
  const index = lower.indexOf(ql);
  if (index === -1) return <>{text}</>;
  return (
    <>
      {text.slice(0, index)}
      <mark className="rounded-[3px] bg-amber-200 px-0.5 text-inherit">
        {text.slice(index, index + q.length)}
      </mark>
      {text.slice(index + q.length)}
    </>
  );
}

type SearchResultRowProps = {
  row: ResultRow;
  query: string;
  active: boolean;
  onHover: () => void;
  onSelect: () => void;
};

export function SearchResultRow({
  row,
  query,
  active,
  onHover,
  onSelect,
}: SearchResultRowProps) {
  return (
    <button
      type="button"
      role="option"
      aria-selected={active}
      data-search-active={active ? "true" : undefined}
      onMouseMove={onHover}
      onClick={onSelect}
      className={cn(
        "flex w-full cursor-default items-center gap-3 rounded-lg px-2.5 py-2 text-left outline-none transition-colors",
        active ? "bg-muted text-foreground" : "hover:bg-muted/60",
      )}
    >
      <span className="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground [&_svg]:size-4">
        {row.icon}
      </span>
      <span className="min-w-0 flex-1">
        <span className="block truncate text-[13.5px] font-semibold">
          <Highlighted text={row.title} query={query} />
        </span>
        <span className="block truncate text-xs text-muted-foreground">
          {row.meta}
        </span>
      </span>
      <span className="flex shrink-0 items-center gap-1.5">
        {row.mock && <MockBadge hint={row.mockHint ?? MOCK_HINT} />}
        <span className="text-xs text-muted-foreground">{row.typeLabel}</span>
      </span>
    </button>
  );
}
