"use client";

import { useRef } from "react";
import { useDebounce } from "@/hooks/useDebounce";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemTitle,
} from "@/components/ui/item";

export interface VillaSelectorProps {
  items: Array<Record<"id" | "codigo", string>>;
  onDebounceChange?: (value: string) => void;
  onSelect?: (value: string | null) => void;
}

export function VillaSelector({
  items,
  onDebounceChange,
  onSelect,
}: VillaSelectorProps) {
  const lastNotifiedValueRef = useRef<string | null>(null);

  const debouncedChange = useDebounce((value: string) => {
    if (value === lastNotifiedValueRef.current) return;
    lastNotifiedValueRef.current = value;
    onDebounceChange?.(value);
  }, 300);

  return (
    <>
      <Combobox
        items={items}
        itemToStringLabel={(villa: (typeof items)[number]) => villa.codigo}
        onValueChange={(v) => {
          onSelect?.(v?.id ?? null);
        }}
      >
        <ComboboxInput
          placeholder="Buscar villa..."
          onChange={(e) => debouncedChange(e.target.value)}
        />
        <ComboboxContent>
          <ComboboxEmpty>Sin resultados.</ComboboxEmpty>
          <ComboboxList>
            {(villa: (typeof items)[number]) => (
              <ComboboxItem key={villa.codigo} value={villa}>
                <Item size="xs" className="p-0">
                  <ItemContent>
                    <ItemTitle className="whitespace-nowrap">
                      Villa {villa.codigo}
                    </ItemTitle>
                    <ItemDescription>{villa.id}</ItemDescription>
                  </ItemContent>
                </Item>
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </>
  );
}
